# CommLayers ClickUp Sync Rules

## 1. Why This File Exists

ClickUp operations are rate-limited and can become inconsistent if the same action is sent multiple times through MCP. This file defines the mandatory sync discipline before any future ClickUp mutation.

## 2. Required Workflow

1. Approve local markdown staging files first.
2. Compare intended ClickUp mutations against already staged actions.
3. Merge duplicates and near-duplicates.
4. Prefer one bulk action per stream set over many single actions.
5. If MCP session/auth is risky, open ClickUp in Playwright and allow manual login before continuing.
6. Execute doc updates before task creation only if the update set is complete.
7. Execute task creation in one stream-level bulk pass.

## 3. Duplicate-Action Rule

Before any ClickUp MCP call, answer these questions:
- Is this action already represented in `05-clickup-bulk-payload.md`?
- Is there another pending action targeting the same doc, page, stream, or task family?
- Can these actions be merged into one bulk update or create call?

If yes, do not issue a separate MCP action.

## 4. MCP vs Playwright Policy

Use MCP when:
- auth is healthy
- bulk operations are ready
- rate limit risk is low

Use Playwright first when:
- manual login is required
- session instability or MCP auth ambiguity exists
- UI-based verification is safer than direct MCP mutation

## 5. Sync Order

1. Index and doc revisions
2. Stream-level bulk task creation
3. Task verification and dedupe check
4. Optional later subtask expansion

## 6. Safety Notes

- Do not mix one-off MCP task creation with planned bulk creation for the same stream.
- Do not update the same doc/page from multiple parallel MCP actions.
- If two pending actions touch the same ClickUp resource, consolidate first.
