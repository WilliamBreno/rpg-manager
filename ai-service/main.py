import time

from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException
from fastapi.responses import Response
from pydantic import BaseModel
import chromadb
from langchain_ollama import OllamaEmbeddings
import json

import rag_cache
from llm_provider import generate_text, RAG_LLM_PROVIDER
from pdf_export.fill_dnd5e_sheet import fill_sheet

load_dotenv()

app = FastAPI()

chroma_client = chromadb.PersistentClient(path="./chroma_db")
collection = chroma_client.get_or_create_collection(name="dnd_books")

# Reduz o contexto mandado pro LLM (causa raiz da lentidão, junto com o
# provedor) — ver CLAUDE.md pra medições antes/depois. Antes: 6 resultados
# por combinação (variante de busca × livro) sem teto nenhum de tamanho total
# nem ranking por relevância, gerando ~62KB de prompt.
N_RESULTS_PER_COMBO = 3     # era 6
TOP_N_CHUNKS = 8            # quantos chunks (dos ~28 combinações × N_RESULTS_PER_COMBO) sobrevivem, por relevância
MAX_CONTEXT_CHARS = 6000    # teto rígido de segurança, mesmo depois do corte acima

try:
    rag_cache.ensure_table()
except Exception as e:
    print(f"[CACHE] não foi possível preparar a tabela de cache no boot: {e}")

class SkillQuery(BaseModel):
    class_name: str
    edition: str
    level: int

def get_indexed_books() -> list[str]:
    """Lista os livros distintos indexados na coleção (metadado 'book' de ingest.py)."""
    result = collection.get(include=["metadatas"])
    return sorted({m["book"] for m in result["metadatas"] if m and "book" in m})

# Calculado uma vez no boot — usado para garantir que cada livro indexado
# (ex: um sourcebook menor como "Poder Arcano") tenha uma cota garantida de
# chunks na busca, em vez de competir por espaço no top-k global contra
# livros muito maiores (ex: os 3 volumes do Livro do Jogador).
INDEXED_BOOKS = get_indexed_books()

def search_relevant_chunks(query: str, n_results: int = 10) -> str:
    embeddings = OllamaEmbeddings(model="nomic-embed-text")
    query_embedding = embeddings.embed_query(query)

    results = collection.query(
        query_embeddings=[query_embedding],
        n_results=n_results
    )

    if not results['documents'][0]:
        return ""

    return "\n\n".join(results['documents'][0])

@app.post("/skills")
async def get_skills(query: SkillQuery):
    cache_key = rag_cache.make_cache_key(query.class_name, query.edition, query.level)
    cached = rag_cache.get_cached(cache_key)
    if cached is not None:
        print(f"[CACHE] hit para {cache_key!r} — pulando embed/busca/geração")
        return cached

    # Busca múltiplos contextos para cobrir mais habilidades
    search_queries = [
        f"{query.class_name} at-will powers level {query.level} D&D {query.edition}",
        f"{query.class_name} encounter powers level {query.level} D&D {query.edition}",
        f"{query.class_name} daily powers level {query.level} D&D {query.edition}",
        f"{query.class_name} utility powers level {query.level} D&D {query.edition}",
    ]

    embeddings = OllamaEmbeddings(model="nomic-embed-text")

    # Consulta cada livro indexado separadamente (em vez de um único top-k
    # global) para que livros menores, como "Poder Arcano", não sejam
    # abafados por livros muito maiores (ex: os 3 volumes do Livro do
    # Jogador) na busca semântica.
    books = INDEXED_BOOKS or [None]

    # --- INSTRUMENTAÇÃO DE TEMPO (mantida como observabilidade leve, não só diagnóstico pontual) ---
    t_embed_total = 0.0
    t_search_total = 0.0
    embed_calls = 0
    search_calls = 0

    # (distância, texto) de toda combinação variante×livro — a distância do
    # Chroma é "quanto menor, mais relevante" tanto pra espaço cosine quanto
    # L2, então dá pra rankear tudo junto no fim em vez de manter só um
    # set() deduplicado por texto exato sem noção de relevância (era o que
    # fazia o contexto virar ~62KB: tudo que qualquer uma das 4 variantes ×
    # todos os livros trouxesse ia pro prompt, sem corte).
    candidates: list[tuple[float, str]] = []

    for q in search_queries:
        t0 = time.perf_counter()
        query_embedding = embeddings.embed_query(q)
        t_embed_total += time.perf_counter() - t0
        embed_calls += 1
        for book in books:
            t0 = time.perf_counter()
            results = collection.query(
                query_embeddings=[query_embedding],
                n_results=N_RESULTS_PER_COMBO,
                where={"book": book} if book else None,
                include=["documents", "distances"],
            )
            t_search_total += time.perf_counter() - t0
            search_calls += 1
            docs = results['documents'][0] if results['documents'] else []
            dists = results['distances'][0] if results.get('distances') else [float("inf")] * len(docs)
            for doc, dist in zip(docs, dists):
                candidates.append((dist, doc))

    print(f"[TIMING] embeddings da pergunta: {t_embed_total:.3f}s total ({embed_calls} chamadas, {t_embed_total/embed_calls:.3f}s/cada)")
    print(f"[TIMING] busca vetorial no Chroma: {t_search_total:.3f}s total ({search_calls} chamadas, {t_search_total/search_calls:.3f}s/cada)")

    # Desduplica por texto mantendo a menor distância (mais relevante) de
    # cada chunk repetido entre combinações, depois ordena tudo por
    # relevância e mantém só os TOP_N_CHUNKS mais relevantes no total.
    best_by_text: dict[str, float] = {}
    for dist, text in candidates:
        if text not in best_by_text or dist < best_by_text[text]:
            best_by_text[text] = dist
    ranked = sorted(best_by_text.items(), key=lambda kv: kv[1])  # (texto, distância) crescente

    selected_chunks: list[str] = []
    total_chars = 0
    for text, _dist in ranked[:TOP_N_CHUNKS]:
        # Teto rígido de segurança: mesmo dentro do top-N, não deixa o
        # contexto total passar de MAX_CONTEXT_CHARS (protege contra um
        # cenário futuro com mais livros/variantes de busca).
        if selected_chunks and total_chars + len(text) > MAX_CONTEXT_CHARS:
            break
        selected_chunks.append(text)
        total_chars += len(text)

    context = "\n\n".join(selected_chunks)
    print(f"[TIMING] contexto final: {len(context)} chars ({len(selected_chunks)} chunks de {len(candidates)} candidatos brutos)")
    # --- FIM DA INSTRUMENTAÇÃO DE EMBED/BUSCA/CONTEXTO ---

    if not context:
        return {"error": "Nenhum contexto encontrado nos livros indexados"}

    prompt = f"""You are a D&D {query.edition} expert. Based ONLY on the rulebook content below, list ALL powers for the {query.class_name} class with min_level <= {query.level}.

IMPORTANT RULES:
- Only include powers from min_level 1 to {query.level}
- Use ONLY these power_type values: "at-will", "encounter", "daily", "utility"
- Utility powers start at level 2
- Include ALL at-will powers (they have no level restriction)
- Be thorough - include every power mentioned in the content

RULEBOOK CONTENT:
{context}

RESPOND WITH ONLY A JSON ARRAY. NO OTHER TEXT. Example:
[{{"name":"Power Name","power_type":"at-will","action_type":"Standard Action","range":"Melee weapon","target":"One creature","attack":"STR vs AC","hit":"1d8+STR modifier","miss":"","effect":"","min_level":1}}]

JSON array of ALL {query.class_name} powers up to level {query.level}:"""

    # --- INSTRUMENTAÇÃO DE TEMPO ---
    t0 = time.perf_counter()
    try:
        raw = generate_text(prompt)
    except Exception as e:
        return {"error": f"Falha ao gerar resposta via {RAG_LLM_PROVIDER}: {e}"}
    t_generate = time.perf_counter() - t0
    print(f"[TIMING] geração via {RAG_LLM_PROVIDER}: {t_generate:.3f}s (prompt ~{len(prompt)} chars, contexto ~{len(context)} chars)")
    # --- FIM DA INSTRUMENTAÇÃO DE GERAÇÃO ---

    raw = raw.replace('```json', '').replace('```', '').strip()

    start = raw.find('[')
    end = raw.rfind(']')

    if start == -1 or end == -1:
        return {"error": "IA não retornou JSON válido", "raw": raw}

    json_str = raw[start:end+1]

    try:
        skills = json.loads(json_str)

        type_map = {
            'at-will': 'at-will', 'at will': 'at-will', 'atwill': 'at-will',
            'cantrip': 'at-will', 'unlimited': 'at-will', 'bonus action': 'at-will',
            'encounter': 'encounter', 'per encounter': 'encounter',
            'daily': 'daily', 'per day': 'daily',
            'utility': 'utility',
        }

        filtered = []
        for skill in skills:
            pt = skill.get('power_type', '').lower().strip()
            skill['power_type'] = type_map.get(pt, 'at-will')

            min_level = skill.get('min_level', 1)
            try:
                min_level = int(min_level)
            except:
                min_level = 1

            skill['min_level'] = min_level

            # Filtra apenas habilidades do nível correto
            if min_level <= query.level:
                filtered.append(skill)

        rag_cache.set_cached(cache_key, filtered)
        return filtered

    except json.JSONDecodeError as e:
        return {"error": f"Erro ao parsear JSON: {str(e)}", "raw": json_str}

class AbilityScore(BaseModel):
    valor: int
    mod: int


class SaveOrSkill(BaseModel):
    valor: int
    proficiente: bool


class DeathSaves(BaseModel):
    sucessos: int = 0
    falhas: int = 0


class CharacterSheet5e(BaseModel):
    """
    Payload já totalmente calculado pelo backend Go (character_service.go /
    armor_service.go) — este serviço não reimplementa nenhuma regra de D&D,
    apenas mapeia os valores para os campos do AcroForm da ficha oficial.
    """
    nome: str
    classe_nivel: str
    antecedente: str = ""
    raca: str = ""
    alinhamento: str = ""
    xp: int = 0
    atributos: dict[str, AbilityScore]
    bonus_proficiencia: int
    ca: int
    iniciativa: int
    deslocamento: int
    pv_maximo: int
    pv_atual: int
    pv_temporario: int = 0
    dados_de_vida_total: str = ""
    dados_de_vida: str = ""
    resistencia_morte: DeathSaves = DeathSaves()
    salvaguardas: dict[str, SaveOrSkill]
    pericias: dict[str, SaveOrSkill] = {}
    percepcao_passiva: int = 10
    tracos_personalidade: str = ""
    ideais: str = ""
    vinculos: str = ""
    defeitos: str = ""


@app.post("/export/pdf/5e")
async def export_pdf_5e(character: CharacterSheet5e):
    try:
        pdf_bytes = fill_sheet(character.model_dump())
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Falha ao preencher a ficha: {e}")

    return Response(content=pdf_bytes, media_type="application/pdf")


@app.get("/health")
async def health():
    return {"status": "ok", "message": "AI Service rodando!"}


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)