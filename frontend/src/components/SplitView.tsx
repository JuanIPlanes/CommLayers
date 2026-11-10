import React, { useState, useRef, useEffect, useCallback } from 'react'
import { TreeGraphView } from './TreeGraphView'
import { RightPanel } from './RightPanel'
import { Header } from './Header'
import { GraphData, NodeId, ConnectionMode, ViewMode } from '../types'
import { mockGraphData, paradigmsGraphData } from '../data/graphData'

interface SplitViewProps {
  initialDividerPosition?: number
}

export const SplitView: React.FC<SplitViewProps> = ({
  initialDividerPosition = 0.4
}) => {
  const [viewMode, setViewMode] = useState<ViewMode>('tree')
  const [connectionMode, setConnectionMode] = useState<ConnectionMode>('technical')
  const [selectedNodeId, setSelectedNodeId] = useState<NodeId | null>(null)
  const [dividerPosition, setDividerPosition] = useState(initialDividerPosition)
  const [isDraggingDivider, setIsDraggingDivider] = useState(false)
  const [graphData, setGraphData] = useState<GraphData>(mockGraphData)
  
  const dividerRef = useRef<HTMLDivElement>(null)
  const leftPanelRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    setGraphData(connectionMode === 'technical' ? mockGraphData : paradigmsGraphData)
    setSelectedNodeId(null)
    setViewMode('tree')
  }, [connectionMode])

  const handleDividerMouseDown = useCallback((e: React.MouseEvent) => {
    e.preventDefault()
    setIsDraggingDivider(true)
    document.body.style.cursor = 'col-resize'
    document.body.style.userSelect = 'none'
  }, [])

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (!isDraggingDivider) return
      
      const container = dividerRef.current?.parentElement
      if (!container) return
      
      const rect = container.getBoundingClientRect()
      const newPosition = (e.clientX - rect.left) / rect.width
      
      const clampedPosition = Math.max(0.2, Math.min(0.8, newPosition))
      setDividerPosition(clampedPosition)
    }

    const handleMouseUp = () => {
      setIsDraggingDivider(false)
      document.body.style.cursor = ''
      document.body.style.userSelect = ''
    }

    if (isDraggingDivider) {
      document.addEventListener('mousemove', handleMouseMove)
      document.addEventListener('mouseup', handleMouseUp)
    }

    return () => {
      document.removeEventListener('mousemove', handleMouseMove)
      document.removeEventListener('mouseup', handleMouseUp)
    }
  }, [isDraggingDivider])

  const handleNodeClick = useCallback((nodeId: NodeId) => {
    const node = graphData.nodes.find(n => n.id === nodeId)
    if (!node || node.isPartial) return

    setSelectedNodeId(nodeId)
    setViewMode('split')
    
    const connectedEdges = graphData.edges.filter(e => e.source === nodeId || e.target === nodeId)
    const path = new Set<NodeId>([nodeId])
    connectedEdges.forEach(e => {
      path.add(e.source as NodeId)
      path.add(e.target as NodeId)
    })
    setHighlightedPath(Array.from(path))
  }, [graphData])

  const [highlightedPath, setHighlightedPath] = useState<NodeId[]>([])

  const handleNodeHover = useCallback((_nodeId: NodeId | null) => {
  }, [])

  const handleSearch = useCallback((_query: string) => {
  }, [])

  const handleModeChange = useCallback((mode: ConnectionMode) => {
    setConnectionMode(mode)
  }, [])

  const selectedNode = graphData.nodes.find(n => n.id === selectedNodeId)

  const isTreeMode = viewMode === 'tree'

  return (
    <div className="split-view" style={{ height: '100vh', width: '100vw' }}>
      <Header
        mode={connectionMode}
        onModeChange={handleModeChange}
        onSearch={handleSearch}
      />

      <main className="main-content" style={{ 
        height: `calc(100vh - 48px)`,
        display: 'flex',
        position: 'relative'
      }}>
        <div
          ref={leftPanelRef}
          className="left-panel"
          style={{
            width: isTreeMode ? '100%' : `${dividerPosition * 100}%`,
            height: '100%',
            position: 'relative',
            overflow: 'hidden',
            background: '#0a0a0a',
            transition: viewMode === 'split' ? 'width 0.3s cubic-bezier(0.16, 1, 0.3, 1)' : 'none'
          }}
        >
          <TreeGraphView
            graphData={graphData}
            selectedNodeId={selectedNodeId}
            onNodeClick={handleNodeClick}
            onNodeHover={handleNodeHover}
            isFullScreen={isTreeMode}
            highlightedPath={highlightedPath}
            connectionMode={connectionMode}
            width={isTreeMode ? window.innerWidth : window.innerWidth * dividerPosition}
            height={window.innerHeight - 48}
          />
        </div>

        {!isTreeMode && (
          <div
            ref={dividerRef}
            className="divider"
            style={{
              width: '1px',
              height: '100%',
              background: '#222222',
              cursor: 'col-resize',
              position: 'relative',
              zIndex: 10,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              touchAction: 'none'
            }}
            onMouseDown={handleDividerMouseDown}
            onTouchStart={(e) => handleDividerMouseDown(e as any)}
            role="separator"
            aria-label="Resize panels"
            aria-valuemin={20}
            aria-valuemax={80}
            aria-valuenow={Math.round(dividerPosition * 100)}
            tabIndex={0}
            onKeyDown={(e) => {
              if (e.key === 'ArrowLeft') {
                e.preventDefault()
                setDividerPosition(Math.max(0.2, dividerPosition - 0.02))
              } else if (e.key === 'ArrowRight') {
                e.preventDefault()
                setDividerPosition(Math.min(0.8, dividerPosition + 0.02))
              }
            }}
          >
            <div className="divider-handle" style={{
              width: '8px',
              height: '60px',
              background: isDraggingDivider ? '#444444' : '#333333',
              borderRadius: '4px',
              transition: 'background 0.15s ease, transform 0.15s ease',
              transform: isDraggingDivider ? 'scaleX(1.5)' : 'scaleX(1)',
              boxShadow: '0 0 8px rgba(255,255,255,0.05)'
            }} />
          </div>
        )}

        {!isTreeMode && selectedNode && (
          <div
            className="right-panel"
            style={{
              width: `${(1 - dividerPosition) * 100}%`,
              height: '100%',
              position: 'relative',
              overflow: 'hidden',
              background: '#0a0a0a',
              transition: 'width 0.3s cubic-bezier(0.16, 1, 0.3, 1)',
              display: 'flex',
              flexDirection: 'column'
            }}
          >
            <RightPanel
              node={selectedNode}
              graphData={graphData}
              highlightedPath={highlightedPath}
            />
          </div>
        )}
      </main>

      <div className="keyboard-hints" style={{
        position: 'fixed',
        bottom: '16px',
        right: '16px',
        display: 'flex',
        gap: '12px',
        opacity: 0.4,
        pointerEvents: 'none',
        zIndex: 100,
        fontSize: '11px',
        color: '#666'
      }}>
        <span><kbd style={{padding:'2px 6px',border:'1px solid #222',borderRadius:'4px',background:'rgba(255,255,255,0.05)'}}>Esc</kbd> Close detail</span>
        <span><kbd style={{padding:'2px 6px',border:'1px solid #222',borderRadius:'4px',background:'rgba(255,255,255,0.05)'}}>Tab</kbd> Navigate</span>
        <span><kbd style={{padding:'2px 6px',border:'1px solid #222',borderRadius:'4px',background:'rgba(255,255,255,0.05)'}}>←/→</kbd> Resize</span>
      </div>
    </div>
  )
}