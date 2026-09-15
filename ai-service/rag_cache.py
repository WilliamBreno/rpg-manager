"""
Cache simples de respostas do RAG de skills (4e), guardado no mesmo Postgres
(Neon) que o backend Go já usa — reaproveita as variáveis de ambiente
DB_HOST/DB_PORT/DB_USER/DB_PASSWORD/DB_NAME do padrão já estabelecido em
backend/.env (ver CLAUDE.md, seção "Backend (Go)"), só que lidas aqui via
ai-service/.env própria (arquivo novo, gitignored, não compartilhado com o
backend por simplicidade — aponte pras mesmas credenciais do Neon se quiser
usar o mesmo banco).

Chave de cache: como `/skills` não recebe uma "pergunta" em texto livre (só
class_name/edition/level), a normalização "lowercase + trim" pedida foi
adaptada para esses três campos combinados — é isso que varia de uma consulta
repetida pra outra.
"""
import json
import os

import psycopg2
from dotenv import load_dotenv

load_dotenv()

TABLE_NAME = "rag_skill_cache"


def _get_connection():
    return psycopg2.connect(
        host=os.getenv("DB_HOST"),
        port=os.getenv("DB_PORT", "5432"),
        user=os.getenv("DB_USER"),
        password=os.getenv("DB_PASSWORD"),
        dbname=os.getenv("DB_NAME"),
        # "prefer" funciona tanto contra um Postgres local sem SSL (setup de
        # dev atual, ver backend/.env) quanto contra Neon (que exige SSL) —
        # se um dia isso apontar pra Neon, defina DB_SSLMODE=require.
        sslmode=os.getenv("DB_SSLMODE", "prefer"),
    )


def ensure_table() -> None:
    """Idempotente — chamado no boot do serviço, mesmo espírito do AutoMigrate do Go."""
    with _get_connection() as conn:
        with conn.cursor() as cur:
            cur.execute(f"""
                CREATE TABLE IF NOT EXISTS {TABLE_NAME} (
                    cache_key TEXT PRIMARY KEY,
                    response TEXT NOT NULL,
                    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
                )
            """)
        conn.commit()


def make_cache_key(class_name: str, edition: str, level: int) -> str:
    """Normalização pedida (lowercase + trim de espaços) aplicada aos campos
    que de fato variam entre consultas equivalentes neste endpoint."""
    return f"{class_name.strip().lower()}|{edition.strip().lower()}|{level}"


def get_cached(cache_key: str):
    """Retorna o valor já decodificado (list/dict) ou None se não houver cache."""
    try:
        with _get_connection() as conn:
            with conn.cursor() as cur:
                cur.execute(f"SELECT response FROM {TABLE_NAME} WHERE cache_key = %s", (cache_key,))
                row = cur.fetchone()
                if row is None:
                    return None
                return json.loads(row[0])
    except Exception as e:
        # Falha de cache nunca deveria derrubar a consulta real — só loga e
        # segue como se fosse um cache-miss.
        print(f"[CACHE] erro ao ler cache ({cache_key}): {e}")
        return None


def set_cached(cache_key: str, value) -> None:
    try:
        payload = json.dumps(value, ensure_ascii=False)
        with _get_connection() as conn:
            with conn.cursor() as cur:
                cur.execute(
                    f"""
                    INSERT INTO {TABLE_NAME} (cache_key, response)
                    VALUES (%s, %s)
                    ON CONFLICT (cache_key) DO UPDATE SET response = EXCLUDED.response, created_at = now()
                    """,
                    (cache_key, payload),
                )
            conn.commit()
    except Exception as e:
        print(f"[CACHE] erro ao gravar cache ({cache_key}): {e}")
