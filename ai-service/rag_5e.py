"""
Ingestão e consulta RAG dos livros de regras de D&D 5e.

Pipeline separado e isolado do pipeline 4e (`ingest.py` / ChromaDB collection
`dnd_books`): usa extração direta de texto via PyMuPDF (os PDFs de 5e têm
camada de texto real, diferente da suposição de PDF escaneado que motivou o
OCR no pipeline 4e).

Armazenamento: NÃO usa ChromaDB (ver nota abaixo) — os embeddings ficam num
arquivo numpy (`chroma_db_5e/embeddings.npy`, já normalizados por linha) e os
metadados num JSON paralelo (`chroma_db_5e/metadata.json`), na mesma ordem.
Busca por similaridade é feita por força bruta (produto escalar contra a
matriz normalizada) — nessa escala (~10 mil vetores de 768 dimensões) isso é
instantâneo e exato, sem precisar de um índice aproximado.

Por que não ChromaDB: o `chromadb` 1.5.9 (bindings Rust) tem um bug real e
reproduzível de persistência do índice HNSW local no Windows — o
`PersistentClient` grava documentos/metadados no SQLite normalmente, mas o
segmento de índice vetorial nunca é de fato persistido em disco (só o
`index_metadata.pickle`, nunca os `.bin` reais), e toda query subsequente
falha com `InternalError: ... Error loading hnsw index`. Isso reproduziu de
forma idêntica em duas ingestões completas independentes (uma delas rodando
`client.close()` explicitamente ao final, o que não resolveu), sobreviveu a
`chroma vacuum` e a reiniciar via `chroma run` (servidor HTTP) apontando pro
mesmo diretório — não é causado por acesso concorrente. Ver issues públicas
do projeto chroma-core/chroma sobre HNSW não persistindo no Windows.
"""
import json
import os

import fitz  # PyMuPDF
import numpy as np
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_ollama import OllamaEmbeddings

BOOKS_DIR = os.path.join(os.path.dirname(__file__), "books")
STORE_DIR = os.path.join(os.path.dirname(__file__), "chroma_db_5e")
EMBEDDINGS_PATH = os.path.join(STORE_DIR, "embeddings.npy")
METADATA_PATH = os.path.join(STORE_DIR, "metadata.json")

# Escopo restrito (decisão do usuário): só os capítulos do Livro do Jogador 2024
# relevantes pra criação de personagem — classes, origens, talentos e magias.
# Não indexa regras gerais de mestre, criaturas nem equipamento (isso último já
# tem um catálogo estruturado próprio em seed_items_5e.go, não precisa de RAG).
# Faixas de página confirmadas via os cabeçalhos "CAPÍTULO N | ..." do PDF.
SCOPE = [
    # (arquivo_pdf, nome_livro, pagina_inicio, pagina_fim_inclusive, categoria)
    ("D&D 5.0 - livro-do-jogador-2024.pdf", "D&D 5.0 - livro-do-jogador-2024 (Criação de Personagens)", 39, 54, "classes"),
    ("D&D 5.0 - livro-do-jogador-2024.pdf", "D&D 5.0 - livro-do-jogador-2024 (Classes)", 55, 182, "classes"),
    ("D&D 5.0 - livro-do-jogador-2024.pdf", "D&D 5.0 - livro-do-jogador-2024 (Origens)", 183, 204, "classes"),
    ("D&D 5.0 - livro-do-jogador-2024.pdf", "D&D 5.0 - livro-do-jogador-2024 (Talentos)", 205, 218, "talentos"),
    ("D&D 5.0 - livro-do-jogador-2024.pdf", "D&D 5.0 - livro-do-jogador-2024 (Magias)", 241, 375, "magias"),
]


def _load_store() -> tuple[np.ndarray, list[dict]]:
    if not os.path.exists(EMBEDDINGS_PATH) or not os.path.exists(METADATA_PATH):
        return np.zeros((0, 768), dtype=np.float32), []
    embeddings = np.load(EMBEDDINGS_PATH)
    with open(METADATA_PATH, encoding="utf-8") as fh:
        metadatas = json.load(fh)
    return embeddings, metadatas


def _save_store(embeddings: np.ndarray, metadatas: list[dict]) -> None:
    os.makedirs(STORE_DIR, exist_ok=True)
    np.save(EMBEDDINGS_PATH, embeddings)
    with open(METADATA_PATH, "w", encoding="utf-8") as fh:
        json.dump(metadatas, fh, ensure_ascii=False)


def count() -> int:
    _, metadatas = _load_store()
    return len(metadatas)


def is_already_indexed(book_name: str) -> bool:
    _, metadatas = _load_store()
    return any(m["livro"] == book_name for m in metadatas)


def extract_pages(pdf_path: str, page_start: int | None = None, page_end: int | None = None) -> list[tuple[int, str]]:
    """Extrai texto de cada página via a camada de texto nativa do PDF (sem OCR).
    page_start/page_end são 1-based e inclusivos; None = do início/até o fim."""
    doc = fitz.open(pdf_path)
    pages = []
    for page_num, page in enumerate(doc, start=1):
        if page_start and page_num < page_start:
            continue
        if page_end and page_num > page_end:
            break
        text = page.get_text().strip()
        if text:
            pages.append((page_num, text))
    return pages


def _normalize(vectors: np.ndarray) -> np.ndarray:
    norms = np.linalg.norm(vectors, axis=1, keepdims=True)
    norms[norms == 0] = 1.0
    return vectors / norms


def ingest_book(
    pdf_path: str,
    book_name: str,
    categoria: str,
    embeddings_model: OllamaEmbeddings,
    page_start: int | None = None,
    page_end: int | None = None,
) -> int:
    """Processa e indexa uma faixa de páginas de um livro de 5e. Retorna nº de chunks indexados."""
    if is_already_indexed(book_name):
        print(f"⏭️  {book_name} já está indexado, pulando.")
        return 0

    print(f"\n📚 Iniciando ingestão: {book_name} (páginas {page_start or 1}-{page_end or '?'})")
    pages = extract_pages(pdf_path, page_start, page_end)
    print(f"  📄 {len(pages)} páginas com texto extraídas")

    splitter = RecursiveCharacterTextSplitter(
        chunk_size=1000,
        chunk_overlap=200,
        separators=["\n\n", "\n", ".", " "],
    )

    existing_embeddings, existing_metadatas = _load_store()
    new_vectors: list[list[float]] = []
    new_metadatas: list[dict] = []

    def checkpoint():
        if not new_vectors:
            return
        combined = np.vstack([existing_embeddings, np.array(new_vectors, dtype=np.float32)]) \
            if existing_embeddings.shape[0] else np.array(new_vectors, dtype=np.float32)
        _save_store(_normalize(combined), existing_metadatas + new_metadatas)

    total_chunks = 0
    last_page = pages[-1][0] if pages else 0
    for page_num, page_text in pages:
        chunks = splitter.split_text(page_text)
        for i, chunk in enumerate(chunks):
            embedding = embeddings_model.embed_query(chunk)
            new_vectors.append(embedding)
            new_metadatas.append({
                "documento": chunk,
                "livro": book_name,
                "edicao": "5e",
                "categoria": categoria,
                "pagina": page_num,
            })
            total_chunks += 1
        if total_chunks and total_chunks % 200 == 0:
            checkpoint()
        print(f"  Indexando... página {page_num}/{last_page} ({total_chunks} chunks até agora)", end="\r")

    checkpoint()
    print(f"\n  ✅ {book_name} indexado com sucesso! ({total_chunks} chunks)")
    return total_chunks


def ingest_all() -> None:
    embeddings_model = OllamaEmbeddings(model="nomic-embed-text")

    if not os.path.exists(BOOKS_DIR):
        print(f"⚠️ Pasta de livros não encontrada: {BOOKS_DIR}")
        return

    for pdf_file, book_name, page_start, page_end, categoria in SCOPE:
        pdf_path = os.path.join(BOOKS_DIR, pdf_file)
        if not os.path.exists(pdf_path):
            print(f"⚠️ {pdf_file} não encontrado em {BOOKS_DIR}, pulando.")
            continue
        ingest_book(pdf_path, book_name, categoria, embeddings_model, page_start, page_end)

    print("\n🎲 Ingestão 5e concluída!")


def query_5e(question: str, n_results: int = 5, categoria: str | None = None) -> list[dict]:
    """Consulta a base de conhecimento 5e. Retorna trechos com metadata de fonte (livro/página)."""
    embeddings, metadatas = _load_store()
    if embeddings.shape[0] == 0:
        return []

    embeddings_model = OllamaEmbeddings(model="nomic-embed-text")
    query_vec = np.array(embeddings_model.embed_query(question), dtype=np.float32)
    query_vec = query_vec / (np.linalg.norm(query_vec) or 1.0)

    if categoria:
        candidate_idx = [i for i, m in enumerate(metadatas) if m["categoria"] == categoria]
    else:
        candidate_idx = list(range(len(metadatas)))
    if not candidate_idx:
        return []

    scores = embeddings[candidate_idx] @ query_vec
    top_local = np.argsort(-scores)[:n_results]
    top_idx = [candidate_idx[i] for i in top_local]

    return [
        {
            "texto": metadatas[i]["documento"],
            "livro": metadatas[i]["livro"],
            "pagina": metadatas[i]["pagina"],
            "categoria": metadatas[i]["categoria"],
        }
        for i in top_idx
    ]


if __name__ == "__main__":
    ingest_all()
