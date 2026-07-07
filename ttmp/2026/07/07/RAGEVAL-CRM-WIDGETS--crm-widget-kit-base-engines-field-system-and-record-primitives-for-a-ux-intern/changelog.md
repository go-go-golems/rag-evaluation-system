# Changelog

## 2026-07-07

- Initial workspace created


## 2026-07-07

Created ticket + intern-facing CRM widget-kit design guide (design-doc/01) in textbook style: the engine+contract+preset pattern applied to CRM — domain model, the field system (FieldSpec + FieldRenderer contract, read/edit table), the engine catalog (BoardEngine kanban as the signature new engine, RecordShell, ActivityFeed, FieldRenderer/RecordFieldList, StatTile, FilterBar + reused MatrixGrid/CollectionPanel/SegmentedBar/MonthGrid), four core screens as ASCII+YAML compositions, IR/DSL wiring (crm.dsl), backend actions, and build order. Uploaded to reMarkable.


## 2026-07-07

Implemented the CRM widget kit end to end (M1-M5), mirroring the scheduling kit and following the existing decomposition. src/crm domain module; the field system (FieldRenderer + RecordFieldList engines, FieldSpec IR); BoardEngine kanban + DealCard; ActivityFeed timeline; RecordShell organism; StatTile; crm.ts presets (pipelineBoard/contactRecord/crmDashboard/tasksInbox). Generic engines register under data.dsl (no crm.dsl module); DealCard is a molecule; retro square style (no border-radius); extensive per-widget + WidgetRenderer stories. Verified all four core screens render via Storybook+Playwright (screenshots in various/screenshots). Typecheck + lefthook green on every commit.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/BoardEngine/BoardEngine.tsx — The kanban engine (drag-between-columns); BoardCardPayload contract
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/FieldRenderer/FieldRenderer.tsx — The field-system engine (read/edit per FieldType); FieldRenderPayload contract
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/crm/types.ts — CRM domain DTOs incl. FieldType/FieldValue/FieldDef — the shared vocabulary
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/widgets/defaultRegistry.ts — Registers the generic engines under dataWidgetRegistry (no crm.dsl module)
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/widgets/presets/crm.ts — CRM presets — pipelineBoard/contactRecord/crmDashboard/tasksInbox emit configured IR

