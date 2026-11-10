import React, { useEffect, useRef, useState, useCallback } from 'react'
import * as d3Force from 'd3-force'
import * as d3Selection from 'd3-selection'
import * as d3Drag from 'd3-drag'
import { GraphData, NodeData, EdgeData, NodeId } from '../types'

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

export const TreeGraphView: React.FC<TreeGraphViewProps> = ({
  graphData,
  selectedNodeId,
  onNodeClick,
  onNodeHover,
  isFullScreen,
  highlightedPath = [],
  width,
  height
}) => {
  const svgRef = useRef<SVGSVGElement>(null)
  const containerRef = useRef<HTMLDivElement>(null)
  const [simulation, setSimulation] = useState<d3Force.Simulation<SimNode, SimLink> | null>(null)
  const [dimensions, setDimensions] = useState({ width, height })
  const animationFrameRef = useRef<number>()
  const nodesRef = useRef<Map<NodeId, d3Selection.Selection<SVGGElement, SimNode, any, unknown>>>(new Map())

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
  }, [simulation, graphData.edges, selectedNodeId, highlightedPath, onNodeClick, onNodeHover])

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
        </div>
      )}
    </div>
  )
}