from fastapi import FastAPI
from pydantic import BaseModel
import chromadb
from langchain_ollama import OllamaEmbeddings
import ollama
import json

app = FastAPI()

chroma_client = chromadb.PersistentClient(path="./chroma_db")
collection = chroma_client.get_or_create_collection(name="dnd_books")

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
    # Busca múltiplos contextos para cobrir mais habilidades
    search_queries = [
        f"{query.class_name} at-will powers level {query.level} D&D {query.edition}",
        f"{query.class_name} encounter powers level {query.level} D&D {query.edition}",
        f"{query.class_name} daily powers level {query.level} D&D {query.edition}",
        f"{query.class_name} utility powers level {query.level} D&D {query.edition}",
    ]
    
    embeddings = OllamaEmbeddings(model="nomic-embed-text")
    all_chunks = set()

    # Consulta cada livro indexado separadamente (em vez de um único top-k
    # global) para que livros menores, como "Poder Arcano", não sejam
    # abafados por livros muito maiores (ex: os 3 volumes do Livro do
    # Jogador) na busca semântica.
    books = INDEXED_BOOKS or [None]

    for q in search_queries:
        query_embedding = embeddings.embed_query(q)
        for book in books:
            results = collection.query(
                query_embeddings=[query_embedding],
                n_results=6,
                where={"book": book} if book else None,
            )
            if results['documents'][0]:
                for doc in results['documents'][0]:
                    all_chunks.add(doc)
    
    context = "\n\n".join(list(all_chunks))
    
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

    response = ollama.generate(model="llama3.2", prompt=prompt)
    raw = response.response
    
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
        
        return filtered
        
    except json.JSONDecodeError as e:
        return {"error": f"Erro ao parsear JSON: {str(e)}", "raw": json_str}

@app.get("/health")
async def health():
    return {"status": "ok", "message": "AI Service rodando!"}