# CommLayers ClickUp Staging Package

This folder is the local staging area for all CommLayers documents and task payloads that will later be synced to ClickUp.

Purpose:
- finalize backend-first project documents before any further ClickUp mutation
- stage bulk task data to avoid duplicate or rate-limited MCP actions
- preserve a clean local source of truth before MCP or Playwright-driven sync

Files:
- `01-prd-backend-first.md` - backend-first product requirements
- `02-tdd-backend-first.md` - backend-first technical design
- `03-fsd-sds-backend-first.md` - functional and software design details
- `04-clickup-task-streams.md` - stream-level task definitions for bulk creation
- `05-clickup-bulk-payload.md` - staged bulk payload structure for ClickUp creation/update
- `06-clickup-sync-rules.md` - operating rules for MCP and Playwright sync

ClickUp target mapping:
- Index doc/page: `8ckkc5r-7153` / `8ckkc5r-7193`
- PRD doc/page: `8ckkc5r-7193` / `8ckkc5r-7233`
- TDD doc/page: `8ckkc5r-7253` / `8ckkc5r-7293`
- FSD+SDS doc/page: `8ckkc5r-7273` / `8ckkc5r-7313`

Current architecture stance:
- v1 is backend-first and 2D desktop visualization only
- v2 adds architectural paradigms as a separate implementation track
- all backend communication families are in scope for v1 as working comparative implementations
- frontend is a governed visualization and comparison consumer of backend evidence

Non-negotiable operating rules:
1. Before any ClickUp MCP action, check whether the same or near-identical action is already pending.
2. If multiple similar actions exist, batch them into one bulk action instead of issuing separate MCP calls.
3. Prefer stream-level bulk task creation over one-off task creation.
4. If ClickUp auth/session or rate limits make MCP unsafe, open ClickUp via Playwright first and allow manual login before continuing.
5. No doc/task sync happens until these local markdown files are approved as the staging source of truth.
