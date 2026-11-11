---
id: DEC-024
title: Desktop Comparison UI Contract
status: proposed
date: 2026-11-11
---

## Context
Frontend comparative visualization (stream 86agbv0kg) requires formal component APIs for TreeGraphView, SplitView, RightPanel with FPS verification integration. This decision formalizes the desktop comparison UI contract.

## Decision
### Component Architecture
**Component Hierarchy**:
```
ComparisonDashboard (root)
├── TreeGraphView (left pane) — backend family topology
├── SplitView (resizable divider)
└── RightPanel (right pane) — detail view
    ├── MetricsPanel
    ├── EvidencePanel
    └── DelayTimelinePanel
```

### Component Contracts

#### 1. TreeGraphView
**Props**:
```typescript
interface TreeGraphViewProps {
  // Data
  topology: TopologyNode[];           // From /api/catalog + /api/deferred-waves
  selectedNodeId: string | null;      // Currently selected
  highlightPath: string[];            // Path from root to selected
  
  // Interaction
  onNodeSelect: (nodeId: string) => void;
  onNodeHover: (nodeId: string | null) => void;
  
  // Visual
  nodeRenderer: (node: TopologyNode) => ReactNode;
  edgeRenderer: (edge: TopologyEdge) => ReactNode;
  layout: 'hierarchical' | 'radial' | 'force';
  
  // Performance
  maxNodes: number;                   // Hard limit (default: 500)
  virtualize: boolean;                // Enable virtualization
}
```

**TopologyNode** (from `graphNodes` in `main.go`):
```typescript
interface TopologyNode {
  id: string;                         // e.g., "workflow_engine"
  label: string;                      // e.g., "Workflow engine"
  category: 'transport' | 'messaging' | 'sync' | 'backend-family' | 'frontend' | 'enrichment';
  status: 'bootstrap-ready' | 'bootstrap-demo' | 'planned' | 'deferred' | 'held';
  dependsOn: string[];                // Parent node IDs
  metrics?: NodeMetrics;              // Optional real-time metrics
  evidence?: EvidenceContract;        // From DEC-011
}
```

**Required Behaviors**:
- Render 500 nodes at 60fps (virtualized)
- Keyboard navigation (arrow keys, Enter to select)
- Search/filter by label, category, status
- Tooltip on hover with metrics + evidence summary
- Click selects node → updates RightPanel

#### 2. SplitView
**Props**:
```typescript
interface SplitViewProps {
  children: [ReactNode, ReactNode];   // [leftPane, rightPane]
  defaultSizes: [number, number];     // Percentage [30, 70]
  minSize: number;                    // Minimum pane % (default: 15)
  maxSize: number;                    // Maximum pane % (default: 85)
  onResize: (sizes: [number, number]) => void;
  direction: 'horizontal' | 'vertical';
}
```

**Required Behaviors**:
- Drag handle to resize (mouse/touch)
- Persist sizes to localStorage
- Keyboard shortcut: `Ctrl/Cmd + [` / `]` to adjust
- Minimum pane size enforced
- Smooth resize animation (150ms)

#### 3. RightPanel (with sub-panels)
**Props**:
```typescript
interface RightPanelProps {
  selectedNode: TopologyNode | null;
  tabs: ['metrics', 'evidence', 'delays'];
  activeTab: string;
  onTabChange: (tab: string) => void;
}
```

**MetricsPanel** (when `activeTab === 'metrics'`):
- Real-time metrics from `/api/benchmarks/metrics` + family-specific endpoints
- Sparkline charts for p50/p95/p99 over time
- Threshold indicators (green/yellow/red per DEC-011)
- Auto-refresh: 5s interval

**EvidencePanel** (when `activeTab === 'evidence'`):
- Evidence contract from selected node
- Fields: source, confidence, staleness, measured_at, delay_applied_ms
- Success/error code badge
- Link to raw evidence JSON

**DelayTimelinePanel** (when `activeTab === 'delays'`):
- Visible delay interactions (from DEC-026)
- Timeline of delayAppliedMs for recent operations
- Correlation with state transitions
- PKI observability delay metadata (from DEC-004)

### FPS Verification Contract (from 86agbv0pn - done)
**Target**: 60fps sustained
**Minimum**: 30fps (below = degradation alert)
**Measurement**:
```typescript
interface FPSMetrics {
  current: number;        // Current frame rate
  average: number;        // 30-second rolling average
  min: number;            // Minimum in window
  frameDrops: number;     // Frames > 16.67ms
  longFrames: number;     // Frames > 33.33ms (2 frames)
}
```

**Verification Endpoint**: `/api/frontend/fps` (client-reported)
```json
{
  "component": "TreeGraphView",
  "fps": { "current": 58, "average": 59.2, "min": 45, "frameDrops": 12, "longFrames": 3 },
  "timestamp": "2026-11-11T10:30:00Z",
  "sessionId": "sess-abc123"
}
```

**CI Gate**: Automated test runs comparison dashboard for 60s, fails if average < 55fps

### Data Flow Contract
1. **Initial Load**: `ComparisonDashboard` fetches `/api/catalog` + `/api/deferred-waves`
2. **Topology Build**: Transform to `TopologyNode[]` with `dependsOn` edges
3. **Selection**: User clicks node → `onNodeSelect` → `RightPanel` updates
4. **Metrics Poll**: `MetricsPanel` polls `/api/benchmarks/metrics` + family endpoints (5s)
5. **Evidence Fetch**: `EvidencePanel` fetches evidence for selected node (on select)
6. **Delay Timeline**: `DelayTimelinePanel` fetches `/api/async/demo` + PKI delays (on select)

### Accessibility Contract
- ARIA labels on all interactive elements
- Keyboard navigation: Tab order, arrow keys in TreeGraphView
- Focus visible on all interactive elements
- Color contrast: WCAG AA (4.5:1) for all text
- Screen reader: Live region for metric updates

### Responsive Breakpoints (Desktop Only)
| Breakpoint | Layout |
|------------|--------|
| ≥ 1440px | Full three-pane (TreeGraphView 30%, SplitView 5%, RightPanel 65%) |
| 1024-1439px | TreeGraphView collapsible (icon-only when collapsed) |
| < 1024px | Not supported (mobile is separate stream) |

## Consequences
- Component APIs are frozen: no prop changes without gate 86agbv0kg re-approval
- FPS verification is a hard gate: CI fails if average < 55fps
- Data flow is unidirectional: dashboard → panels (no panel → dashboard)
- Accessibility is non-negotiable: WCAG AA minimum
- Desktop only: mobile comparison is a separate future stream

## Acceptance Criteria
- [ ] All three components implement specified props exactly
- [ ] TreeGraphView renders 500 nodes at 60fps (virtualized)
- [ ] SplitView persists sizes, respects min/max
- [ ] RightPanel tabs switch correctly with data loading
- [ ] FPS verification test passes in CI (avg ≥ 55fps for 60s)
- [ ] Accessibility audit passes WCAG AA
- [ ] Data flow matches contract: no circular dependencies