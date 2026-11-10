export type NodeId = string & { readonly __brand: unique symbol }
export type EdgeId = string & { readonly __brand: unique symbol }

export interface Position {
  x: number
  y: number
}

// D3 SimulationNodeDatum requires these properties
export interface SimulationNodeDatum {
  index?: number
  x?: number
  y?: number
  vx?: number
  vy?: number
  fx?: number | null
  fy?: number | null
}

export interface NodeData extends SimulationNodeDatum {
  id: NodeId
  label: string
  type: 'paradigm' | 'transport' | 'protocol' | 'pattern' | 'infrastructure' | 'concept'
  description: string
  details: {
    what: string
    how: string
    when: string
    where: string
    why: string
    productStage: 'exploration' | 'prototype' | 'beta' | 'production' | 'deprecated'
    userRPSThresholds: string
    liveCloudCosts: string
  }
  flowDiagram?: FlowDiagramData
  variants?: VariantData[]
  liveDemo?: LiveDemoData
  position?: Position
  isFocused?: boolean
  isPartial?: 'left' | 'right' | 'overflow-left' | 'overflow-right'
}

export interface EdgeData {
  id: EdgeId
  source: NodeId | NodeData
  target: NodeId | NodeData
  type: 'connection' | 'dependency' | 'evolution' | 'variant'
  opacity?: number
}

export interface FlowDiagramData {
  title: string
  phases: FlowPhase[]
}

export interface FlowPhase {
  id: string
  label: string
  description: string
  status: 'pending' | 'active' | 'completed'
  animationDelay?: number
}

export interface VariantData {
  name: string
  description: string
  pros: string[]
  cons: string[]
  recommendedFor: string[]
  avoidWhen: string[]
}

export interface LiveDemoData {
  connectionStatus: 'connecting' | 'connected' | 'disconnected' | 'error'
  latencyMs: number
  messagesPerSecond: number
  indicators: LiveIndicator[]
}

export interface LiveIndicator {
  id: string
  label: string
  value: string
  status: 'ok' | 'warning' | 'error'
  pulse?: boolean
}

export interface GraphData {
  nodes: NodeData[]
  edges: EdgeData[]
}

export type ViewMode = 'tree' | 'split'

export type ConnectionMode = 'technical' | 'paradigms'

export interface AppState {
  viewMode: ViewMode
  connectionMode: ConnectionMode
  selectedNodeId: NodeId | null
  graphData: GraphData
  dividerPosition: number
  isDraggingDivider: boolean
  zoomLevel: number
  panOffset: Position
}