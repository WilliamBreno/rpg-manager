"""
Abstração da etapa de GERAÇÃO do RAG de skills — trocável via
RAG_LLM_PROVIDER=anthropic|ollama (padrão: anthropic).

Nota importante (achado da investigação, diferente do que se assumia antes de
implementar): o pipeline original não usava nenhum wrapper de LLM do
LangChain pra geração — `main.py` chamava `ollama.generate(...)` direto, via
o pacote `ollama`. LangChain só era usado pra EMBEDDING (`OllamaEmbeddings`),
não para geração. Não existia, portanto, um "wrapper do LLM" no LangChain
pra simplesmente trocar por `ChatAnthropic` — esta função substitui a chamada
direta ao `ollama.generate` por uma chamada direta ao SDK oficial da
Anthropic, mantendo o mesmo formato de entrada (um prompt de texto único) e
saída (uma string), sem introduzir uma dependência nova do LangChain pra
Anthropic que o resto do arquivo não usava.
"""
import os

RAG_LLM_PROVIDER = os.getenv("RAG_LLM_PROVIDER", "anthropic").strip().lower()

# Haiku — modelo mais rápido/barato da Anthropic, conforme pedido.
ANTHROPIC_MODEL = os.getenv("ANTHROPIC_MODEL", "claude-haiku-4-5-20251001")
ANTHROPIC_MAX_TOKENS = int(os.getenv("ANTHROPIC_MAX_TOKENS", "4096"))

OLLAMA_MODEL = os.getenv("OLLAMA_MODEL", "llama3.2")


def _generate_anthropic(prompt: str) -> str:
    from anthropic import Anthropic

    api_key = os.getenv("ANTHROPIC_API_KEY")
    if not api_key:
        raise RuntimeError(
            "ANTHROPIC_API_KEY não configurada. Defina a variável de ambiente "
            "(ex: no ai-service/.env) ou use RAG_LLM_PROVIDER=ollama pra rodar offline."
        )
    client = Anthropic(api_key=api_key)
    message = client.messages.create(
        model=ANTHROPIC_MODEL,
        max_tokens=ANTHROPIC_MAX_TOKENS,
        messages=[{"role": "user", "content": prompt}],
    )
    return "".join(block.text for block in message.content if block.type == "text")


def _generate_ollama(prompt: str) -> str:
    import ollama

    response = ollama.generate(model=OLLAMA_MODEL, prompt=prompt)
    return response.response


def generate_text(prompt: str) -> str:
    """Ponto único de geração — troca de provedor via RAG_LLM_PROVIDER, sem
    tocar em retriever/prompt (mesmo contrato: prompt string -> resposta string)."""
    if RAG_LLM_PROVIDER == "anthropic":
        return _generate_anthropic(prompt)
    elif RAG_LLM_PROVIDER == "ollama":
        return _generate_ollama(prompt)
    else:
        raise ValueError(f"RAG_LLM_PROVIDER inválido: {RAG_LLM_PROVIDER!r} (use 'anthropic' ou 'ollama')")
