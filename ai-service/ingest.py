import os
import pytesseract
from PIL import Image
import fitz  # PyMuPDF
import chromadb
from langchain_text_splitters import RecursiveCharacterTextSplitter
from langchain_ollama import OllamaEmbeddings

# Caminho do Tesseract
pytesseract.pytesseract.tesseract_cmd = r'C:\Program Files\Tesseract-OCR\tesseract.exe'

# Configuração do ChromaDB
chroma_client = chromadb.PersistentClient(path="./chroma_db")
collection = chroma_client.get_or_create_collection(
    name="dnd_books",
    metadata={"hnsw:space": "cosine"}
)

def extract_text_from_pdf(pdf_path: str) -> str:
    """Extrai texto de PDF escaneado usando OCR"""
    doc = fitz.open(pdf_path)
    full_text = ""
    
    print(f"Processando {pdf_path} — {len(doc)} páginas")
    
    for page_num, page in enumerate(doc):
        print(f"  Página {page_num + 1}/{len(doc)}...", end="\r")
        
        # Renderiza a página como imagem
        mat = fitz.Matrix(2, 2)  # zoom 2x para melhor OCR
        pix = page.get_pixmap(matrix=mat)
        img = Image.frombytes("RGB", [pix.width, pix.height], pix.samples)
        
        # OCR na imagem
        text = pytesseract.image_to_string(img, lang='por+eng')
        full_text += f"\n--- Página {page_num + 1} ---\n{text}"
    
    print(f"\n  ✅ {pdf_path} processado!")
    return full_text

def ingest_pdf(pdf_path: str, book_name: str):
    """Processa e indexa um PDF no ChromaDB"""
    print(f"\n📚 Iniciando ingestão: {book_name}")
    
    # Extrai texto
    text = extract_text_from_pdf(pdf_path)
    
    if not text.strip():
        print(f"  ⚠️ Nenhum texto extraído de {book_name}")
        return
    
    # Divide em chunks
    splitter = RecursiveCharacterTextSplitter(
        chunk_size=1000,
        chunk_overlap=200,
        separators=["\n\n", "\n", ".", " "]
    )
    chunks = splitter.split_text(text)
    print(f"  📄 {len(chunks)} chunks criados")
    
    # Gera embeddings e salva no ChromaDB
    embeddings = OllamaEmbeddings(model="nomic-embed-text")
    
    for i, chunk in enumerate(chunks):
        print(f"  Indexando chunk {i+1}/{len(chunks)}...", end="\r")
        
        embedding = embeddings.embed_query(chunk)
        
        collection.add(
            documents=[chunk],
            embeddings=[embedding],
            metadatas=[{"book": book_name, "chunk": i}],
            ids=[f"{book_name}_{i}"]
        )
    
    print(f"\n  ✅ {book_name} indexado com sucesso!")

if __name__ == "__main__":
    # Lista os PDFs na pasta books/
    books_dir = "./books"
    
    if not os.path.exists(books_dir):
        os.makedirs(books_dir)
        print("📁 Pasta 'books' criada!")
        print("   Coloque seus PDFs de D&D dentro de 'books/' e rode novamente.")
    else:
        pdf_files = [f for f in os.listdir(books_dir) if f.endswith('.pdf')]
        
        if not pdf_files:
            print("⚠️ Nenhum PDF encontrado em 'books/'")
            print("   Coloque seus PDFs de D&D dentro de 'books/' e rode novamente.")
        else:
            for pdf_file in pdf_files:
                pdf_path = os.path.join(books_dir, pdf_file)
                book_name = pdf_file.replace('.pdf', '')
                ingest_pdf(pdf_path, book_name)
            
            print("\n🎲 Todos os livros foram indexados com sucesso!")