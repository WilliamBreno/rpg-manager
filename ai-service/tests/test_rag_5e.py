"""
Teste de integração da ingestão/consulta RAG de 5e (Fase 1 de TASKS_UI_E_FEATURES.md).

Escopo restrito (decisão do usuário): só os capítulos de Criação de
Personagens/Classes/Origens/Talentos/Magias do Livro do Jogador 2024 — não
inclui regras gerais de mestre, criaturas nem equipamento (esse último já tem
catálogo estruturado próprio, ver seed_items_5e.go).

Requer Ollama rodando localmente com o modelo nomic-embed-text e a base
`chroma_db_5e/` já populada via `python rag_5e.py`.
"""
import sys
import os

sys.path.insert(0, os.path.dirname(os.path.dirname(__file__)))

from rag_5e import query_5e, count


def test_base_5e_populada():
    assert count() > 0, "base 5e está vazia — rode `python rag_5e.py` antes do teste"


def test_query_pericias_bardo_cita_fonte():
    results = query_5e("quais são as perícias (skill proficiencies) que um bardo pode escolher")

    assert results, "nenhum resultado retornado para a query sobre perícias de bardo"

    livros = {r["livro"] for r in results}
    assert all(l.startswith("D&D 5.0 - livro-do-jogador-2024") for l in livros), \
        f"resultado citando fonte fora do escopo restrito: {livros}"

    textos = " ".join(r["texto"].lower() for r in results)
    assert "bardo" in textos, "trecho retornado não menciona 'bardo' — resultado pouco relevante"


def test_query_talentos_filtra_categoria():
    results = query_5e("quais são os pré-requisitos do talento Atirador de Elite", categoria="talentos")

    assert results, "nenhum resultado retornado para query filtrada por categoria 'talentos'"
    for r in results:
        assert r["categoria"] == "talentos"
        assert "Talentos" in r["livro"]


def test_query_magias_cita_fonte():
    results = query_5e("qual o efeito da magia Bola de Fogo", categoria="magias")

    assert results, "nenhum resultado retornado para query de magias"
    for r in results:
        assert r["categoria"] == "magias"
        assert "Magias" in r["livro"]
