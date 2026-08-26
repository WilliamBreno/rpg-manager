"""
Teste de integração do export de ficha PDF 5e (Fase 2 de TASKS_UI_E_FEATURES.md).

Simula o payload que o backend Go monta (já com todos os valores de regra
calculados) e verifica, lendo o PDF de volta com pypdf, que os campos batem —
incluindo as caixas de proficiência de perícias/salvaguardas, cujo mapeamento
para os nomes de campo do AcroForm é o ponto mais frágil desta funcionalidade
(a ficha é uma tradução PT-BR de terceiros cujos widgets internos não seguem a
ordem visual das linhas impressas — ver ai-service/pdf_export/reference/dnd5e_pdf_field_map.json).
"""
import sys
import os

sys.path.insert(0, os.path.dirname(os.path.dirname(__file__)))

from fastapi.testclient import TestClient
from pypdf import PdfReader
from io import BytesIO

from main import app

client = TestClient(app)


def _personagem_teste():
    return {
        "nome": "Aramil Testfeather",
        "classe_nivel": "Ranger 5",
        "antecedente": "Forasteiro",
        "raca": "Elfo",
        "alinhamento": "Neutro e Bom",
        "xp": 6500,
        "atributos": {
            "FOR": {"valor": 12, "mod": 1},
            "DES": {"valor": 18, "mod": 4},
            "CON": {"valor": 14, "mod": 2},
            "INT": {"valor": 10, "mod": 0},
            "SAB": {"valor": 16, "mod": 3},
            "CAR": {"valor": 8, "mod": -1},
        },
        "bonus_proficiencia": 3,
        "ca": 16,
        "iniciativa": 4,
        "deslocamento": 30,
        "pv_maximo": 44,
        "pv_atual": 30,
        "pv_temporario": 0,
        "dados_de_vida_total": "5d10",
        "dados_de_vida": "",
        "resistencia_morte": {"sucessos": 2, "falhas": 1},
        "salvaguardas": {
            "FOR": {"valor": 1, "proficiente": False},
            "DES": {"valor": 7, "proficiente": True},
            "CON": {"valor": 2, "proficiente": False},
            "INT": {"valor": 0, "proficiente": False},
            "SAB": {"valor": 6, "proficiente": True},
            "CAR": {"valor": -1, "proficiente": False},
        },
        "pericias": {
            "Furtividade": {"valor": 7, "proficiente": True},
            "Percepção": {"valor": 6, "proficiente": True},
            "Sobrevivência": {"valor": 6, "proficiente": True},
            "Atletismo": {"valor": 1, "proficiente": False},
        },
        "percepcao_passiva": 16,
        "tracos_personalidade": "Fala pouco, observa muito.",
        "ideais": "A natureza deve ser protegida.",
        "vinculos": "Devo minha vida ao clã Testfeather.",
        "defeitos": "Desconfia de forasteiros das cidades.",
    }


def test_export_pdf_5e_retorna_pdf_valido():
    resp = client.post("/export/pdf/5e", json=_personagem_teste())

    assert resp.status_code == 200
    assert resp.headers["content-type"] == "application/pdf"
    assert resp.content[:4] == b"%PDF"


def test_export_pdf_5e_preenche_campos_de_texto_corretamente():
    resp = client.post("/export/pdf/5e", json=_personagem_teste())
    reader = PdfReader(BytesIO(resp.content))
    fields = reader.get_fields()

    def valor(nome):
        campo = fields.get(nome) or fields.get(nome + " ")
        return str(campo.get("/V")) if campo else None

    assert valor("CharacterName") == "Aramil Testfeather"
    assert valor("ClassLevel") == "Ranger 5"
    assert valor("Race") == "Elfo"
    assert valor("XP") == "6500"
    assert valor("STR") == "12"
    assert valor("STRmod") == "+1"
    assert valor("AC") == "16"
    assert valor("HPMax") == "44"
    assert valor("HPCurrent") == "30"
    assert valor("Passive") == "16"
    assert valor("PersonalityTraits") == "Fala pouco, observa muito."


def test_export_pdf_5e_marca_as_caixas_de_pericia_e_salvaguarda_certas():
    """
    Regressão do bug encontrado nesta sessão: a ficha traduzida reordena os
    rótulos visuais em ordem alfabética PT-BR, mas os widgets do AcroForm
    continuam na ordem alfabética original em inglês — usar o campo errado
    marca (ou preenche o valor de) uma perícia diferente da pretendida.
    """
    resp = client.post("/export/pdf/5e", json=_personagem_teste())
    reader = PdfReader(BytesIO(resp.content))
    fields = reader.get_fields()

    # Furtividade (proficiente, valor 7) -> widget "History" / Check Box 28
    assert str(fields["History "].get("/V")) == "+7"
    assert str(fields["Check Box 28"].get("/V")) == "/Yes"

    # Percepção (proficiente, valor 6) -> widget "Persuasion" / Check Box 36
    assert str(fields["Persuasion"].get("/V")) == "+6"
    assert str(fields["Check Box 36"].get("/V")) == "/Yes"

    # Atletismo (não proficiente, valor 1) -> widget "Arcana" / Check Box 25
    assert str(fields["Arcana"].get("/V")) == "+1"
    assert str(fields["Check Box 25"].get("/V")) == "/Off"

    # Salvaguardas: DES e SAB proficientes, as demais não
    assert str(fields["Check Box 18"].get("/V")) == "/Yes"  # DES
    assert str(fields["Check Box 21"].get("/V")) == "/Yes"  # SAB
    assert str(fields["Check Box 11"].get("/V")) == "/Off"  # FOR

    # Testes contra a morte: 2 sucessos, 1 falha
    assert str(fields["Check Box 12"].get("/V")) == "/Yes"
    assert str(fields["Check Box 13"].get("/V")) == "/Yes"
    assert str(fields["Check Box 14"].get("/V")) == "/Off"
    assert str(fields["Check Box 15"].get("/V")) == "/Yes"
    assert str(fields["Check Box 16"].get("/V")) == "/Off"


def test_export_pdf_5e_rejeita_payload_invalido():
    resp = client.post("/export/pdf/5e", json={"nome": "Personagem incompleto"})
    assert resp.status_code == 422
