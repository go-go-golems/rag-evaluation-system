#!/usr/bin/env python3
"""Generate a Markdown inventory of the current pkg/widgetdsl module surface.

The script is intentionally lightweight: it parses the simple map literals and
moduleSpec entries in pkg/widgetdsl/module.go and applies a small hand-maintained
classification table. It is a planning aid for the Widget DSL v3 redesign, not a
Go parser.
"""
from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[6]
MODULE_GO = ROOT / "pkg/widgetdsl/module.go"
OUT = ROOT / "ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/04-widget-dsl-current-export-inventory.md"

MODULE_MAPS = {
    "ui.dsl": "uiHelpers",
    "data.dsl": "dataHelpers",
    "context_window.dsl": "contextWindowHelpers",
    "course.dsl": "courseHelpers",
    "cms.dsl": "cmsHelpers",
}

GENERIC_UI_HELPERS = {
    "appShell", "appNav", "button", "breadcrumbs", "emptyState", "fieldGrid",
    "markdownArticle", "meterBar", "pagination", "richArticle", "searchField",
    "tag", "tileGrid", "uploadDropArea", "caption", "checkList", "codeText",
    "dashboardGrid", "divider", "figureBlock", "formPanel", "formRow", "inline",
    "keyPointList", "keyValueStrip", "metadataGrid", "panel", "personSummary",
    "scrollRegion", "sectionBlock", "selectInput", "sidebarNav", "sidebarShell",
    "splitPane", "stack", "statusText", "stepList", "tabList", "textBlock",
    "textInput", "textareaInput",
}

ENGINE_HELPERS = {
    "dataTable",
    "contextBudgetBar", "contextGroupedStripDiagram", "contextLegend",
    "contextStackDiagram", "contextStripDiagram", "contextTreemap",
}

DOMAIN_HELPERS = {
    "anchoredCommentCard", "anchoredCommentRail", "annotationBadge",
    "annotationNoteCard", "annotationRailPanel", "contextDiagramPanel",
    "contextStyleSwatch", "contextTurnPagerPanel", "transcriptMessageCard",
    "transcriptReaderPanel", "transcriptRoleBadge", "transcriptSessionHeader",
    "transcriptWorkspacePanel", "contextUploadDropArea",
    "articleListPanel", "assetTile", "cmsShell", "contentStatusBadge",
    "markdownEditor", "mediaLibraryPanel", "mediaThumb",
    "contextStudioNavIcon", "courseLessonPanel", "courseSlidePanel",
    "courseStepNav", "courseStudioShell", "documentListPanel",
    "documentPreviewToolbar", "handoutDocumentShell", "slideShell",
}

GENERIC_ALIASES_IN_DOMAIN = {
    "breadcrumbs", "emptyState", "meterBar", "pagination", "searchField",
    "tag", "tileGrid", "markdownArticle", "richArticle",
}


def parse_helper_map(source: str, var_name: str) -> list[tuple[str, str]]:
    m = re.search(rf"var\s+{re.escape(var_name)}\s*=\s*map\[string\]string\{{(.*?)\n\}}", source, re.S)
    if not m:
        return []
    body = m.group(1)
    pairs: list[tuple[str, str]] = []
    for key, value in re.findall(r'"([^"]+)"\s*:\s*"([^"]+)"', body):
        pairs.append((key, value))
    return sorted(pairs)


def parse_recipes(source: str) -> dict[str, list[str]]:
    module_const_to_name = {
        "UIModuleName": "ui.dsl",
        "DataModuleName": "data.dsl",
        "DataV2ModuleName": "data.v2.dsl",
        "ContextWindowModuleName": "context_window.dsl",
        "CourseModuleName": "course.dsl",
        "CmsModuleName": "cms.dsl",
    }
    out: dict[str, list[str]] = {name: [] for name in module_const_to_name.values()}
    block = re.search(r"var\s+moduleSpecs\s*=\s*\[\]moduleSpec\{(.*?)\n\}\n\nvar moduleSpecsByName", source, re.S)
    if not block:
        return out
    for entry in re.findall(r"\{\s*name:\s*([A-Za-z0-9_]+),(.*?)(?=\n\s*\},)", block.group(1), re.S):
        const, body = entry
        mod = module_const_to_name.get(const)
        if not mod:
            continue
        r = re.search(r"recipes:\s*\[\]string\{([^}]*)\}", body)
        if r:
            out[mod] = re.findall(r'"([^"]+)"', r.group(1))
    return out


def classify(module: str, helper: str) -> str:
    if module == "ui.dsl" or helper in GENERIC_UI_HELPERS:
        if module != "ui.dsl" and helper in GENERIC_ALIASES_IN_DOMAIN:
            return "generic alias currently exported from a domain module"
        return "foundation/generic UI helper"
    if helper in ENGINE_HELPERS:
        return "engine-level helper"
    if helper in DOMAIN_HELPERS:
        return "domain component/helper"
    return "unclassified — review manually"


def main() -> None:
    source = MODULE_GO.read_text()
    recipes = parse_recipes(source)
    lines: list[str] = []
    lines.extend([
        "---",
        "Title: Widget DSL Current Export Inventory",
        "Ticket: RAGEVAL-SCHEDULE-WIDGETS",
        "Status: active",
        "Topics:",
        "    - ui-dsl",
        "    - widget-ir",
        "    - frontend-architecture",
        "DocType: reference",
        "Intent: long-term",
        "Owners: []",
        "RelatedFiles: []",
        "ExternalSources: []",
        "Summary: \"Generated inventory of the current pkg/widgetdsl module/helper/recipe surface used to plan Widget DSL v3.\"",
        "LastUpdated: 2026-07-07T16:55:00-04:00",
        "WhatFor: \"Use this to see what the old split DSL modules expose before porting functionality to widget.dsl.\"",
        "WhenToUse: \"Read during Widget DSL v3 Phase 0 inventory and whenever classifying old helpers as foundation, engine, domain, or compatibility aliases.\"",
        "---",
        "",
        "# Widget DSL Current Export Inventory",
        "",
        "This file is generated by `scripts/01-widget-dsl-export-inventory.py`. It inventories the current `pkg/widgetdsl/module.go` helper maps and recipe lists before the clean `widget.dsl` redesign begins.",
        "",
        "## Summary",
        "",
        "| Module | Helper count | Recipes | Notes |",
        "|---|---:|---|---|",
    ])
    for module, var_name in MODULE_MAPS.items():
        helpers = parse_helper_map(source, var_name)
        recipe_list = recipes.get(module, [])
        notes = "new v3 should preserve capability, not necessarily names"
        lines.append(f"| `{module}` | {len(helpers)} | {', '.join('`'+r+'`' for r in recipe_list) or '—'} | {notes} |")
    lines.append(f"| `data.v2.dsl` | 0 direct helper-map helpers | — | typed/fluent builder experiment installed separately |")
    lines.extend(["", "## Helper inventory", ""])
    for module, var_name in MODULE_MAPS.items():
        helpers = parse_helper_map(source, var_name)
        lines.extend([f"### `{module}`", "", "| Helper | Component type | Classification |", "|---|---|---|"])
        for helper, component in helpers:
            lines.append(f"| `{helper}` | `{component}` | {classify(module, helper)} |")
        recipe_list = recipes.get(module, [])
        if recipe_list:
            lines.extend(["", "Recipes:", ""])
            for recipe in recipe_list:
                lines.append(f"- `{module}.recipes.{recipe}(options)`")
        lines.append("")
    lines.extend([
        "## V3 implications",
        "",
        "- Generic helpers should live under `widget.dsl.ui` or slot helper `h`, not under every domain namespace.",
        "- Engine-level helpers should become typed engine builders such as `data.collection`, `data.matrix`, `time.week`, or explicit `raw.component` calls during experiments.",
        "- Domain component helpers and recipes should become domain views such as `cms.mediaLibrary`, `course.handouts`, and `context.workspace`.",
        "- `data.v2.dsl` is valuable implementation precedent, but v3 should expose `widget.dsl.data` rather than a separate public v2 module.",
        "- Current recipes are capability references; v3 should preserve their product behavior while changing names and authoring style freely.",
    ])
    OUT.write_text("\n".join(lines) + "\n")


if __name__ == "__main__":
    main()
