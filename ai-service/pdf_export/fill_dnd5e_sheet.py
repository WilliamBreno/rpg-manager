"""
Preenchimento da ficha oficial de personagem 5e (AcroForm) a partir dos dados
já calculados pelo backend Go (character_service.go / armor_service.go).

Este módulo não reimplementa nenhuma regra de D&D — todo o cálculo (modificadores,
bônus de proficiência, CA, PV, DVs, etc.) é feito no Go e chega aqui pronto. A
única responsabilidade deste módulo é mapear esses valores para os nomes de campo
do AcroForm, usando `reference/dnd5e_pdf_field_map.json` como fonte da verdade.
"""
import json
import os
from io import BytesIO

from pypdf import PdfReader, PdfWriter

REFERENCE_DIR = os.path.join(os.path.dirname(__file__), "reference")
FIELD_MAP_PATH = os.path.join(REFERENCE_DIR, "dnd5e_pdf_field_map.json")
SHEET_TEMPLATE_PATH = os.path.join(
    os.path.dirname(os.path.dirname(__file__)),
    "books",
    "D&D 5.0 - ficha-de-personagem-completavel-biblioteca-elfica.pdf",
)

with open(FIELD_MAP_PATH, encoding="utf-8") as fh:
    FIELD_MAP = json.load(fh)


def _signed(n: int) -> str:
    return f"+{n}" if n >= 0 else str(n)


def character_to_field_values(character: dict) -> dict:
    """
    Converte o payload calculado pelo Go (ver contrato abaixo) nos valores de
    campo do AcroForm (nome_do_campo -> valor de texto, ou nome_do_campo -> "/Yes"
    para checkboxes marcadas).

    Contrato esperado de `character` (todos os valores já calculados):
    {
      "nome": str, "classe_nivel": str, "antecedente": str, "raca": str,
      "alinhamento": str, "xp": int,
      "atributos": {"FOR": {"valor": int, "mod": int}, "DES": {...}, ... "CAR": {...}},
      "bonus_proficiencia": int, "ca": int, "iniciativa": int, "deslocamento": int,
      "pv_maximo": int, "pv_atual": int, "pv_temporario": int,
      "dados_de_vida_total": str, "dados_de_vida": str,
      "resistencia_morte": {"sucessos": int, "falhas": int},
      "salvaguardas": {"FOR": {"valor": int, "proficiente": bool}, ...},
      "pericias": {"Acrobacia": {"valor": int, "proficiente": bool}, ...},  # chave = nome_pt
      "percepcao_passiva": int,
      "tracos_personalidade": str, "ideais": str, "vinculos": str, "defeitos": str,
    }
    """
    fm = FIELD_MAP
    values: dict = {}

    ident = fm["identidade"]
    values[ident["nome_personagem"]] = character.get("nome", "")
    values[ident["classe_nivel"]] = character.get("classe_nivel", "")
    values[ident["antecedente"]] = character.get("antecedente", "")
    values[ident["raca"]] = character.get("raca", "")
    values[ident["alinhamento"]] = character.get("alinhamento", "")
    values[ident["xp"]] = str(character.get("xp", ""))

    for attr, campos in fm["atributos"].items():
        dados = character.get("atributos", {}).get(attr, {})
        values[campos["valor"]] = str(dados.get("valor", ""))
        values[campos["mod"]] = _signed(dados.get("mod", 0)) if "mod" in dados else ""

    combate = fm["combate"]
    values[combate["bonus_proficiencia"]] = _signed(character.get("bonus_proficiencia", 0))
    values[combate["ca"]] = str(character.get("ca", ""))
    values[combate["iniciativa"]] = _signed(character.get("iniciativa", 0))
    values[combate["deslocamento"]] = str(character.get("deslocamento", ""))
    values[combate["pv_maximo"]] = str(character.get("pv_maximo", ""))
    # PV atual: nunca sobrescrever com valor calculado — usar exatamente o que o
    # jogador tem registrado (bug conhecido #2 do projeto: PV manual não pode ser
    # sobrescrito por cálculo automático).
    values[combate["pv_atual"]] = str(character.get("pv_atual", ""))
    values[combate["pv_temporario"]] = str(character.get("pv_temporario", ""))
    values[combate["dados_de_vida_total"]] = str(character.get("dados_de_vida_total", ""))
    values[combate["dados_de_vida"]] = str(character.get("dados_de_vida", ""))

    morte = fm["resistencia_contra_morte"]
    sucessos = character.get("resistencia_morte", {}).get("sucessos", 0)
    falhas = character.get("resistencia_morte", {}).get("falhas", 0)
    for i, campo in enumerate(morte["sucessos"]):
        values[campo] = "/Yes" if i < sucessos else "/Off"
    for i, campo in enumerate(morte["falhas"]):
        values[campo] = "/Yes" if i < falhas else "/Off"

    for save in fm["salvaguardas"]:
        dados = character.get("salvaguardas", {}).get(save["atributo"], {})
        values[save["campo_valor"]] = _signed(dados.get("valor", 0))
        values[save["campo_proficiencia"]] = "/Yes" if dados.get("proficiente") else "/Off"

    for pericia in fm["pericias"]:
        dados = character.get("pericias", {}).get(pericia["nome_pt"], {})
        values[pericia["campo_valor"]] = _signed(dados.get("valor", 0))
        values[pericia["campo_proficiencia"]] = "/Yes" if dados.get("proficiente") else "/Off"

    values[fm["passiva"]["percepcao_passiva"]] = str(character.get("percepcao_passiva", ""))

    personalidade = fm["personalidade"]
    values[personalidade["tracos_personalidade"]] = character.get("tracos_personalidade", "") or ""
    values[personalidade["ideais"]] = character.get("ideais", "") or ""
    values[personalidade["vinculos"]] = character.get("vinculos", "") or ""
    values[personalidade["defeitos"]] = character.get("defeitos", "") or ""

    return values


def fill_sheet(character: dict) -> bytes:
    """Preenche a ficha oficial 5e com os dados do personagem e retorna os bytes do PDF."""
    values = character_to_field_values(character)

    reader = PdfReader(SHEET_TEMPLATE_PATH)
    writer = PdfWriter()
    writer.append(reader)

    for page in writer.pages:
        writer.update_page_form_field_values(page, values, auto_regenerate=False)

    writer.set_need_appearances_writer(True)

    buffer = BytesIO()
    writer.write(buffer)
    return buffer.getvalue()
