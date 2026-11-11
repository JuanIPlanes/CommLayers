import React, { useEffect, useRef, useState, useCallback } from 'react'
import * as d3Force from 'd3-force'
import * as d3Selection from 'd3-selection'
import * as d3Drag from 'd3-drag'
import { GraphData, NodeData, EdgeData, NodeId } from '../types'

// Core Types
interface TreeGraphViewProps {
  graphData: GraphData
  selectedNodeId: NodeId | null
  onNodeClick: (nodeId: NodeId) => void
  onNodeHover: (nodeId: NodeId | null) => void
  isFullScreen: boolean
  highlightedPath?: NodeId[]
  connectionMode: 'technical' | 'paradigms'
  width: number
  height: number
}

const NODE_RADIUS = {
  focused: 56,
  normal: 36,
  small: 28,
  overflow: 36
}

interface SimNode extends NodeData {}
interface SimLink extends EdgeData {
  source: SimNode
  target: SimNode
}

// FPS Monitoring Types
interface FPSMetrics {
  current: number
  average: number
  min: number
  max: number
  frameCount: number
  lastUpdate: number
}

interface ComparisonVerificationResult {
  mode: 'technical' | 'paradigms'
  nodeCount: number
  edgeCount: number
  renderTimeMs: number
  fps: number
  passed: boolean
  thresholds: {
    minFPS: number
    maxRenderTimeMs: number
  }
}

// FPS Monitor Hook
const useFPSMonitor = (enabled: boolean = true) => {
  const metricsRef = useRef<FPSMetrics>({
    current: 0,
    average: 0,
    min: 60,
    max: 0,
    frameCount: 0,
    lastUpdate: performance.now(),
  })
  const frameTimesRef = useRef<number[]>([])
  const animationFrameRef = useRef<number>()
  const [fps, setFPS] = useState<number>(0)
  const [avgFPS, setAvgFPS] = useState<number>(0)

  const start = useCallback(() => {
    if (!enabled) return
    const tick = (now: number) => {
      const metrics = metricsRef.current
      const delta = now - metrics.lastUpdate
      metrics.frameCount++
      
      if (delta >= 1000) {
        const currentFPS = Math.round((metrics.frameCount * 1000) / delta)
        metrics.current = currentFPS
        metrics.min = Math.min(metrics.min, currentFPS)
        metrics.max = Math.max(metrics.max, currentFPS)
        
        frameTimesRef.current.push(currentFPS)
        if (frameTimesRef.current.length > 60) {
          frameTimesRef.current.shift()
        }
        
        const avg = frameTimesRef.current.reduce((a, b) => a + b, 0) / frameTimesRef.current.length
        metrics.average = Math.round(avg)
        
        setFPS(currentFPS)
        setAvgFPS(Math.round(avg))
        
        metrics.frameCount = 0
        metrics.lastUpdate = now
      }
      
      animationFrameRef.current = requestAnimationFrame(tick)
    }
    animationFrameRef.current = requestAnimationFrame(tick)
  }, [enabled])

  const stop = useCallback(() => {
    if (animationFrameRef.current) {
      cancelAnimationFrame(animationFrameRef.current)
      animationFrameRef.current = undefined
    }
  }, [])

  const reset = useCallback(() => {
    metricsRef.current = {
      current: 0,
      average: 0,
      min: 60,
      max: 0,
      frameCount: 0,
      lastUpdate: performance.now(),
    }
    frameTimesRef.current = []
    setFPS(0)
    setAvgFPS(0)
  }, [])

  const getMetrics = useCallback((): FPSMetrics => {
    return { ...metricsRef.current }
  }, [])

  useEffect(() => {
    if (enabled) {
      start()
    }
    return () => stop()
  }, [enabled, start, stop])

  return { fps, avgFPS, start, stop, reset, getMetrics }
}

// Comparison Verification Hook
const useComparisonVerification = () => {
  const verify = useCallback((
    mode: 'technical' | 'paradigms',
    nodeCount: number,
    edgeCount: number,
    renderTimeMs: number,
    currentFPS: number
  ): ComparisonVerificationResult => {
    const thresholds = {
      minFPS: 30,
      maxRenderTimeMs: 33.33, // ~30 FPS budget
    }
    
    const passed = currentFPS >= thresholds.minFPS && renderTimeMs <= thresholds.maxRenderTimeMs
    
    return {
      mode,
      nodeCount,
      edgeCount,
      renderTimeMs,
      fps: currentFPS,
      passed,
      thresholds,
    }
  }, [])

  return { verify }
}

export const TreeGraphView: React.FC<TreeGraphViewProps> = ({
  graphData,
  selectedNodeId,
  onNodeClick,
  onNodeHover,
  isFullScreen,
  highlightedPath = [],
  connectionMode,
  width,
  height
}) => {
  const svgRef = useRef<SVGSVGElement>(null)
  const containerRef = useRef<HTMLDivElement>(null)
  const [simulation, setSimulation] = useState<d3Force.Simulation<SimNode, SimLink> | null>(null)
  const [dimensions, setDimensions] = useState({ width, height })
  const animationFrameRef = useRef<number>()
  const nodesRef = useRef<Map<NodeId, d3Selection.Selection<SVGGElement, SimNode, any, unknown>>>(new Map())
  
  // FPS Monitoring
  const { fps, avgFPS, getMetrics } = useFPSMonitor(true)
  const { verify } = useComparisonVerification()
  const [comparisonResult, setComparisonResult] = useState<ComparisonVerificationResult | null>(null)
  const [showFPS, setShowFPS] = useState(false)
  const renderStartRef = useRef<number>(0)

  useEffect(() => {
    setDimensions({ width, height })
  }, [width, height])

  useEffect(() => {
    const container = containerRef.current
    if (!container) return

    const nodes: SimNode[] = graphData.nodes.map(n => ({ ...n })) as SimNode[]
    const edges: SimLink[] = graphData.edges.map(e => ({ ...e })) as SimLink[]

    nodes.forEach((node, i) => {
      if (!node.position) {
        const angle = (i / nodes.length) * Math.PI * 2
        const radius = Math.min(width, height) * 0.35
        node.position = {
          x: width / 2 + Math.cos(angle) * radius,
          y: height / 2 + Math.sin(angle) * radius
        }
      }
      if (node.position) {
        node.x = node.position.x
        node.y = node.position.y
      }
    })

    const sim = d3Force.forceSimulation<SimNode, SimLink>(nodes)
      .force('link', d3Force.forceLink<SimNode, SimLink>(edges)
        .id((d: SimNode) => d.id)
        .distance((d: SimLink) => {
          const source = d.source as SimNode
          const target = d.target as SimNode
          if (source.isFocused || target.isFocused) return 180
          return 140
        })
        .strength(0.3)
      )
      .force('charge', d3Force.forceManyBody<SimNode>().strength(-400))
      .force('center', d3Force.forceCenter<SimNode>(width / 2, height / 2))
      .force('collision', d3Force.forceCollide<SimNode>().radius((d: SimNode) => {
        const r = d.isFocused ? NODE_RADIUS.focused : NODE_RADIUS.normal
        return r + 20
      }).strength(0.8))
      .force('x', d3Force.forceX<SimNode>(width / 2).strength(0.02))
      .force('y', d3Force.forceY<SimNode>(height / 2).strength(0.02))
      .alphaDecay(0.02)
      .velocityDecay(0.4)

    setSimulation(sim)

    return () => {
      sim.stop()
    }
  }, [graphData.nodes.length, graphData.edges.length, width, height])

  const render = useCallback(() => {
    const svg = svgRef.current
    const sim = simulation
    if (!svg || !sim) return

    // Track render start time for performance measurement
    renderStartRef.current = performance.now()

    const nodes = sim.nodes() as SimNode[]
    const edges = graphData.edges as SimLink[]

    const edgeSelection = d3Selection.select(svg)
      .selectAll<SVGPathElement, SimLink>('.edge-path')
      .data(edges, (d: SimLink) => d.id)

    edgeSelection.join(
      enter => enter.append('path')
        .attr('class', 'edge-path')
        .attr('stroke', '#ffffff')
        .attr('stroke-width', 1)
        .attr('fill', 'none')
        .attr('stroke-linecap', 'round')
        .attr('opacity', (d: SimLink) => {
          const sourceNode = d.source as SimNode
          const targetNode = d.target as SimNode
          const sourceFocused = sourceNode?.isFocused
          const targetFocused = targetNode?.isFocused
          
          if (d.opacity !== undefined) return d.opacity
          if (sourceFocused || targetFocused) return 0.5
          return 0.25
        })
        .attr('stroke-dasharray', (d: SimLink) => {
          const sourceNode = d.source as SimNode
          const targetNode = d.target as SimNode
          if (sourceNode?.isFocused || targetNode?.isFocused) return '0'
          return '4 4'
        })
        .style('pointer-events', 'none'),
      update => update
        .attr('opacity', (d: SimLink) => {
          const sourceNode = d.source as SimNode
          const targetNode = d.target as SimNode
          const sourceFocused = sourceNode?.isFocused
          const targetFocused = targetNode?.isFocused
          
          if (d.opacity !== undefined) return d.opacity
          if (sourceFocused || targetFocused) return 0.5
          return 0.25
        })
        .attr('stroke-dasharray', (d: SimLink) => {
          const sourceNode = d.source as SimNode
          const targetNode = d.target as SimNode
          if (sourceNode?.isFocused || targetNode?.isFocused) return '0'
          return '4 4'
        })
    )

    const nodeSelection = d3Selection.select(svg)
      .selectAll<SVGGElement, SimNode>('.node-group')
      .data(nodes, (d: SimNode) => d.id)

    nodeSelection.join(
      enter => {
        const g = enter.append('g')
          .attr('class', 'node-group')
          .attr('transform', (d: SimNode) => `translate(${d.x || 0},${d.y || 0})`)
          .style('cursor', (d: SimNode) => d.isPartial ? 'default' : 'pointer')
          .on('click', (event: MouseEvent, d: SimNode) => {
            if (!d.isPartial) {
              event.stopPropagation()
              onNodeClick(d.id)
            }
          })
          .on('mouseenter', (_event: MouseEvent, d: SimNode) => {
            if (!d.isPartial) {
              onNodeHover(d.id)
            }
          })
          .on('mouseleave', () => {
            onNodeHover(null)
          })
          .call(d3Drag.drag<SVGGElement, SimNode>()
            .on('start', (_event: d3Drag.D3DragEvent<SVGGElement, SimNode, SimNode>, d: SimNode) => {
              if (!d.isPartial) {
                sim.alphaTarget(0.3).restart()
                d.fx = d.x ?? null
                d.fy = d.y ?? null
              }
            })
            .on('drag', (event: d3Drag.D3DragEvent<SVGGElement, SimNode, SimNode>, d: SimNode) => {
              if (!d.isPartial) {
                d.fx = event.x
                d.fy = event.y
              }
            })
            .on('end', (_event: d3Drag.D3DragEvent<SVGGElement, SimNode, SimNode>, d: SimNode) => {
              if (!d.isPartial) {
                sim.alphaTarget(0)
                d.fx = null
                d.fy = null
              }
            })
          )

        g.append('circle')
          .attr('class', 'node-circle')
          .attr('r', (d: SimNode) => {
            if (d.isFocused) return NODE_RADIUS.focused
            return NODE_RADIUS.normal
          })
          .attr('fill', (d: SimNode) => {
            if (d.isFocused) return '#ffffff'
            return 'transparent'
          })
          .attr('stroke', (d: SimNode) => {
            if (d.isFocused) return '#ffffff'
            return '#222222'
          })
          .attr('stroke-width', (d: SimNode) => d.isFocused ? 0 : 1.5)
          .attr('filter', (d: SimNode) => d.isFocused ? 'url(#glow)' : 'none')
          .style('opacity', (d: SimNode) => {
            if (d.isFocused) return 1
            if (d.isPartial) return 0.25
            const isHighlighted = highlightedPath.includes(d.id)
            return isHighlighted ? 0.5 : 0.25
          })

        g.append('circle')
          .attr('class', 'node-glow')
          .attr('r', (d: SimNode) => d.isFocused ? NODE_RADIUS.focused + 8 : 0)
          .attr('fill', 'none')
          .attr('stroke', '#ffffff')
          .attr('stroke-width', 2)
          .attr('opacity', (d: SimNode) => d.isFocused ? 0.6 : 0)
          .style('filter', 'blur(8px)')

        g.append('text')
          .attr('class', 'node-label')
          .attr('text-anchor', 'middle')
          .attr('dominant-baseline', 'middle')
          .attr('font-size', (d: SimNode) => {
            if (d.isFocused) return '16px'
            if (d.isPartial) return '0'
            return '12px'
          })
          .attr('font-weight', (d: SimNode) => d.isFocused ? '600' : '500')
          .attr('fill', (d: SimNode) => {
            if (d.isFocused) return '#0a0a0a'
            if (d.isPartial) return 'transparent'
            return '#ffffff'
          })
          .attr('opacity', (d: SimNode) => {
            if (d.isFocused) return 1
            if (d.isPartial) return 0
            const isHighlighted = highlightedPath.includes(d.id)
            return isHighlighted ? 0.5 : 0.25
          })
          .text((d: SimNode) => d.label)
          .style('pointer-events', 'none')
          .style('user-select', 'none')
          .each(function(this: SVGGElement, d: SimNode) {
            nodesRef.current.set(d.id, d3Selection.select(this))
          })

        return g
      },
      update => update
        .attr('transform', (d: SimNode) => `translate(${d.x || 0},${d.y || 0})`)
    )

    edgeSelection.attr('d', (d: SimLink) => {
      const source = d.source as SimNode
      const target = d.target as SimNode
      if (!source.x || !target.x || !source.y || !target.y) return ''
      
      const dx = target.x - source.x
      const dy = target.y - source.y
      const midX = (source.x + target.x) / 2
      const midY = (source.y + target.y) / 2
      const offsetX = -dy * 0.15
      const offsetY = dx * 0.15
      
      return `M${source.x},${source.y} Q${midX + offsetX},${midY + offsetY} ${target.x},${target.y}`
    })

    nodeSelection.attr('transform', (d: SimNode) => `translate(${d.x || 0},${d.y || 0})`)

    nodeSelection.each(function(d: SimNode) {
      if (d.isPartial === 'overflow-left' || d.isPartial === 'overflow-right') {
        const g = d3Selection.select(this)
        const circle = g.select('.node-circle')
        const label = g.select('.node-label')
        
        if (d.isPartial === 'overflow-left') {
          circle.attr('cx', 0).attr('r', NODE_RADIUS.normal)
          label.attr('x', 0).style('display', 'none')
        } else if (d.isPartial === 'overflow-right') {
          circle.attr('cx', 0).attr('r', NODE_RADIUS.normal)
          label.attr('x', 0).style('display', 'none')
        }
      }
    })

    // Calculate render time and run comparison verification
    const renderTimeMs = performance.now() - renderStartRef.current
    const metrics = getMetrics()
    const result = verify(connectionMode, graphData.nodes.length, graphData.edges.length, renderTimeMs, metrics.current)
    setComparisonResult(result)
  }, [simulation, graphData.edges, selectedNodeId, highlightedPath, onNodeClick, onNodeHover, connectionMode, getMetrics, verify])

  useEffect(() => {
    const sim = simulation
    if (!sim) return

    const tick = () => {
      render()
      animationFrameRef.current = requestAnimationFrame(tick)
    }

    sim.on('tick', () => {
    })

    animationFrameRef.current = requestAnimationFrame(tick)

    return () => {
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current)
      }
    }
  }, [simulation, render])

  useEffect(() => {
    if (!simulation) return

    const nodes = simulation.nodes() as SimNode[]
    nodes.forEach(node => {
      node.isFocused = node.id === selectedNodeId
    })

    simulation.alpha(0.3).restart()
  }, [selectedNodeId, simulation])

  const glowFilter = (
    <filter id="glow" x="-50%" y="-50%" width="200%" height="200%">
      <feGaussianBlur stdDeviation="6" result="blur"/>
      <feMerge>
        <feMergeNode in="blur"/>
        <feMergeNode in="SourceGraphic"/>
      </feMerge>
    </filter>
  )

  return (
    <div
      ref={containerRef}
      className="tree-graph-view"
      style={{ width: dimensions.width, height: dimensions.height }}
      role="application"
      aria-label="Communication layers graph visualization"
    >
      <svg
        ref={svgRef}
        width="100%"
        height="100%"
        style={{ display: 'block' }}
        aria-hidden="true"
      >
        <defs>
          {glowFilter}
          <clipPath id="panel-clip">
            <rect x="0" y="0" width={dimensions.width} height={dimensions.height} />
          </clipPath>
        </defs>
        <rect
          width="100%"
          height="100%"
          fill="#0a0a0a"
          clipPath="url(#panel-clip)"
        />
        <g className="edges-layer" clipPath="url(#panel-clip)" />
        <g className="nodes-layer" clipPath="url(#panel-clip)" />
      </svg>
      
      {isFullScreen && (
        <div className="graph-instructions" style={{
          position: 'absolute',
          bottom: '24px',
          left: '50%',
          transform: 'translateX(-50%)',
          display: 'flex',
          gap: '16px',
          color: '#666666',
          fontSize: '11px',
          pointerEvents: 'none',
          zIndex: 10
        }}>
          <span style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
            <kbd style={{ padding: '2px 6px', border: '1px solid #222', borderRadius: '4px', background: 'rgba(255,255,255,0.05)', fontFamily: 'inherit', fontWeight: 500 }}>Click</kbd> node to inspect
          </span>
          <span style={{ color: '#333' }}>·</span>
          <span style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
            <kbd style={{ padding: '2px 6px', border: '1px solid #222', borderRadius: '4px', background: 'rgba(255,255,255,0.05)', fontFamily: 'inherit', fontWeight: 500 }}>Drag</kbd> to reposition
          </span>
          <span style={{ color: '#333' }}>·</span>
          <span style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
            <kbd style={{ padding: '2px 6px', border: '1px solid #222', borderRadius: '4px', background: 'rgba(255,255,255,0.05)', fontFamily: 'inherit', fontWeight: 500 }}>Scroll</kbd> to zoom
          </span>
          <span style={{ color: '#333' }}>·</span>
          <span style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
            <kbd style={{ padding: '2px 6px', border: '1px solid #222', borderRadius: '4px', background: 'rgba(255,255,255,0.05)', fontFamily: 'inherit', fontWeight: 500 }}>F</kbd> toggle FPS
          </span>
        </div>
      )}

      {/* FPS Overlay */}
      {(showFPS || isFullScreen) && (
        <div style={{
          position: 'absolute',
          top: '12px',
          right: '12px',
          zIndex: 20,
          background: 'rgba(10,10,10,0.9)',
          border: '1px solid #222',
          borderRadius: '8px',
          padding: '12px 16px',
          fontFamily: 'monospace',
          fontSize: '11px',
          color: '#ffffff',
          minWidth: '180px',
          boxShadow: '0 4px 12px rgba(0,0,0,0.4)'
        }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
            <span style={{ color: '#888' }}>FPS Monitor</span>
            <button
              onClick={() => setShowFPS(false)}
              style={{
                background: 'none',
                border: 'none',
                color: '#666',
                cursor: 'pointer',
                fontSize: '14px',
                padding: '0',
                lineHeight: 1
              }}
              aria-label="Close FPS monitor"
            >
              ×
            </button>
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '8px' }}>
            <div>
              <div style={{ color: '#666' }}>Current</div>
              <div style={{ fontSize: '18px', fontWeight: 600, color: fps >= 50 ? '#00c864' : fps >= 30 ? '#ffb400' : '#ff4444' }}>
                {fps}
              </div>
            </div>
            <div>
              <div style={{ color: '#666' }}>Average</div>
              <div style={{ fontSize: '18px', fontWeight: 600, color: '#ffffff' }}>
                {avgFPS}
              </div>
            </div>
            <div>
              <div style={{ color: '#666' }}>Min</div>
              <div style={{ color: '#ff4444' }}>{getMetrics().min}</div>
            </div>
            <div>
              <div style={{ color: '#666' }}>Max</div>
              <div style={{ color: '#00c864' }}>{getMetrics().max}</div>
            </div>
          </div>
        </div>
      )}

      {/* Comparison Verification Indicator */}
      {comparisonResult && (
        <div style={{
          position: 'absolute',
          top: showFPS || isFullScreen ? '140px' : '12px',
          right: '12px',
          zIndex: 20,
          background: comparisonResult.passed ? 'rgba(0, 200, 100, 0.15)' : 'rgba(255, 68, 68, 0.15)',
          border: `1px solid ${comparisonResult.passed ? '#00c864' : '#ff4444'}`,
          borderRadius: '8px',
          padding: '10px 12px',
          fontFamily: 'monospace',
          fontSize: '10px',
          color: comparisonResult.passed ? '#00c864' : '#ff4444',
          minWidth: '160px',
          boxShadow: '0 4px 12px rgba(0,0,0,0.4)'
        }}>
          <div style={{ fontWeight: 600, marginBottom: '6px', textTransform: 'uppercase' }}>
            {comparisonResult.passed ? '✓ Verification Passed' : '✗ Verification Failed'}
          </div>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '4px', fontSize: '9px' }}>
            <span style={{ color: '#888' }}>Mode</span>
            <span>{comparisonResult.mode}</span>
            <span style={{ color: '#888' }}>Nodes</span>
            <span>{comparisonResult.nodeCount}</span>
            <span style={{ color: '#888' }}>Edges</span>
            <span>{comparisonResult.edgeCount}</span>
            <span style={{ color: '#888' }}>Render</span>
            <span>{comparisonResult.renderTimeMs.toFixed(2)}ms</span>
            <span style={{ color: '#888' }}>FPS</span>
            <span>{comparisonResult.fps}</span>
            <span style={{ color: '#888' }}>Min FPS</span>
            <span>{comparisonResult.thresholds.minFPS}</span>
          </div>
        </div>
      )}

      {/* Keyboard shortcut for FPS toggle */}
      <div
        onKeyDown={(e) => {
          if (e.key === 'f' || e.key === 'F') {
            e.preventDefault()
            setShowFPS(prev => !prev)
          }
        }}
        tabIndex={0}
        style={{ outline: 'none', position: 'absolute', opacity: 0, pointerEvents: 'none' }}
        aria-label="Press F to toggle FPS monitor"
      />
    </div>
  )
}