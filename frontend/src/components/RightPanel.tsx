import React, { useState, useEffect, useCallback, useRef } from 'react'
import { NodeData, GraphData, FlowDiagramData, VariantData, LiveDemoData, NodeId } from '../types'

interface RightPanelProps {
  node: NodeData
  graphData: GraphData
  highlightedPath: NodeId[]
}

export const RightPanel: React.FC<RightPanelProps> = ({
  node,
  graphData,
  highlightedPath
}) => {
  const [activeTab, setActiveTab] = useState<'flow' | 'variants' | 'demo'>('flow')
  const [carouselIndex, setCarouselIndex] = useState(0)

  useEffect(() => {
    const interval = setInterval(() => {
      setCarouselIndex(prev => (prev + 1) % 3)
    }, 5000)
    return () => clearInterval(interval)
  }, [])

  const handleCarouselNext = useCallback(() => {
    setCarouselIndex(prev => (prev + 1) % 3)
  }, [])

  const handleCarouselPrev = useCallback(() => {
    setCarouselIndex(prev => (prev - 1 + 3) % 3)
  }, [])

  const handleCarouselDotClick = useCallback((index: number) => {
    setCarouselIndex(index)
  }, [])

  const flowDiagram = node.flowDiagram
  const variants = node.variants || []
  const liveDemo = node.liveDemo

  const connectedNodes = graphData.nodes.filter(n => 
    highlightedPath.includes(n.id) && n.id !== node.id && !n.isPartial
  )

  return (
    <div className="right-panel-content" style={{ 
      height: '100%', 
      display: 'flex', 
      flexDirection: 'column',
      background: '#0a0a0a'
    }}>
      <div className="right-panel-top" style={{
        height: '35%',
        display: 'flex',
        flexDirection: 'column',
        borderBottom: '1px solid #222222',
        overflow: 'hidden'
      }}>
        <div className="tab-bar" style={{
          display: 'flex',
          gap: '4px',
          padding: '12px 16px 0',
          borderBottom: '1px solid #222222',
          marginBottom: '0',
          background: 'rgba(255,255,255,0.02)'
        }}>
          <button
            className={`tab-button ${activeTab === 'flow' ? 'active' : ''}`}
            onClick={() => setActiveTab('flow')}
            style={{
              padding: '8px 16px',
              border: 'none',
              background: activeTab === 'flow' ? 'rgba(255,255,255,0.08)' : 'transparent',
              color: activeTab === 'flow' ? '#ffffff' : '#888888',
              borderRadius: '6px 6px 0 0',
              fontSize: '12px',
              fontWeight: 500,
              cursor: 'pointer',
              transition: 'all 0.15s ease',
              borderBottom: activeTab === 'flow' ? '2px solid #ffffff' : '2px solid transparent',
              marginBottom: '-1px'
            }}
          >
            Architecture Flow
          </button>
          <button
            className={`tab-button ${activeTab === 'variants' ? 'active' : ''}`}
            onClick={() => setActiveTab('variants')}
            style={{
              padding: '8px 16px',
              border: 'none',
              background: activeTab === 'variants' ? 'rgba(255,255,255,0.08)' : 'transparent',
              color: activeTab === 'variants' ? '#ffffff' : '#888888',
              borderRadius: '6px 6px 0 0',
              fontSize: '12px',
              fontWeight: 500,
              cursor: 'pointer',
              transition: 'all 0.15s ease',
              borderBottom: activeTab === 'variants' ? '2px solid #ffffff' : '2px solid transparent',
              marginBottom: '-1px'
            }}
          >
            Variants
          </button>
          <button
            className={`tab-button ${activeTab === 'demo' ? 'active' : ''}`}
            onClick={() => setActiveTab('demo')}
            style={{
              padding: '8px 16px',
              border: 'none',
              background: activeTab === 'demo' ? 'rgba(255,255,255,0.08)' : 'transparent',
              color: activeTab === 'demo' ? '#ffffff' : '#888888',
              borderRadius: '6px 6px 0 0',
              fontSize: '12px',
              fontWeight: 500,
              cursor: 'pointer',
              transition: 'all 0.15s ease',
              borderBottom: activeTab === 'demo' ? '2px solid #ffffff' : '2px solid transparent',
              marginBottom: '-1px'
            }}
          >
            Live Demo
          </button>
        </div>

        <div className="tab-content" style={{ flex: 1, overflow: 'auto', padding: '16px' }}>
          {activeTab === 'flow' && flowDiagram && (
            <FlowDiagramView diagram={flowDiagram} />
          )}
          {activeTab === 'variants' && (
            <VariantsView variants={variants} nodeLabel={node.label} />
          )}
          {activeTab === 'demo' && liveDemo && (
            <LiveDemoView demo={liveDemo} nodeLabel={node.label} />
          )}
        </div>
      </div>

      <div className="right-panel-bottom" style={{
        height: '65%',
        display: 'flex',
        flexDirection: 'column',
        overflow: 'hidden'
      }}>
        <div className="node-details-section" style={{
          height: '40%',
          overflow: 'auto',
          padding: '16px',
          borderBottom: '1px solid #222222'
        }}>
          <NodeDetailsCard node={node} connectedNodes={connectedNodes} highlightedPath={highlightedPath} />
        </div>

        <div className="carousel-section" style={{
          height: '60%',
          display: 'flex',
          flexDirection: 'column',
          padding: '16px',
          overflow: 'hidden'
        }}>
          <Carousel
            currentIndex={carouselIndex}
            onNext={handleCarouselNext}
            onPrev={handleCarouselPrev}
            onDotClick={handleCarouselDotClick}
            node={node}
            variants={variants}
            liveDemo={liveDemo}
            flowDiagram={flowDiagram}
          />
        </div>
      </div>
    </div>
  )
}

const FlowDiagramView: React.FC<{ diagram: FlowDiagramData }> = ({ diagram }) => {
  const svgRef = useRef<SVGSVGElement>(null)
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 })

  useEffect(() => {
    const updateDimensions = () => {
      if (svgRef.current) {
        const rect = svgRef.current.getBoundingClientRect()
        setDimensions({ width: rect.width, height: rect.height })
      }
    }
    updateDimensions()
    window.addEventListener('resize', updateDimensions)
    return () => window.removeEventListener('resize', updateDimensions)
  }, [])

  const { width, height } = dimensions
  const phases = diagram.phases
  const phaseCount = phases.length
  const centerX = width / 2
  const centerY = height / 2
  const radius = Math.min(width, height) * 0.35
  const angleStep = (Math.PI * 2) / phaseCount
  const startAngle = -Math.PI / 2

  return (
    <div style={{ width: '100%', height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
      <svg
        ref={svgRef}
        width="100%"
        height="100%"
        viewBox={`0 0 ${width} ${height}`}
        preserveAspectRatio="xMidYMid meet"
        className="flow-diagram-svg"
      >
        <defs>
          <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
            <polygon points="0 0, 10 3.5, 0 7" fill="#ffffff" />
          </marker>
          <filter id="phase-glow" x="-50%" y="-50%" width="200%" height="200%">
            <feGaussianBlur stdDeviation="4" result="blur"/>
            <feMerge>
              <feMergeNode in="blur"/>
              <feMergeNode in="SourceGraphic"/>
            </feMerge>
          </filter>
        </defs>

        {phases.map((phase, i) => {
          const angle1 = startAngle + i * angleStep
          const angle2 = startAngle + ((i + 1) % phaseCount) * angleStep
          const x1 = centerX + Math.cos(angle1) * radius
          const y1 = centerY + Math.sin(angle1) * radius
          const x2 = centerX + Math.cos(angle2) * radius
          const y2 = centerY + Math.sin(angle2) * radius
          
          const isCompleted = phase.status === 'completed'
          const isActive = phase.status === 'active'
          
          return (
            <path
              key={phase.id}
              d={`M${x1},${y1} Q${centerX},${centerY} ${x2},${y2}`}
              stroke={isCompleted ? '#00c864' : (isActive ? '#ffffff' : '#333333')}
              strokeWidth={2}
              fill="none"
              strokeDasharray={isCompleted ? 'none' : '8 4'}
              strokeDashoffset={isActive ? 0 : 20}
              style={{
                animation: isActive ? 'flow-line 1.5s linear infinite' : 'none',
                filter: isActive ? 'url(#phase-glow)' : 'none'
              }}
              markerEnd="url(#arrowhead)"
              opacity={isCompleted ? 0.6 : (isActive ? 1 : 0.3)}
            />
          )
        })}

        {phases.map((phase, i) => {
          const angle = startAngle + i * angleStep
          const x = centerX + Math.cos(angle) * radius
          const y = centerY + Math.sin(angle) * radius
          const nodeRadius = 36
          
          const isCompleted = phase.status === 'completed'
          const isActive = phase.status === 'active'

          return (
            <g key={phase.id} transform={`translate(${x},${y})`}>
              <line
                x1={0}
                y1={0}
                x2={-Math.cos(angle) * (nodeRadius + 8)}
                y2={-Math.sin(angle) * (nodeRadius + 8)}
                stroke={isCompleted ? '#00c864' : (isActive ? '#ffffff' : '#222222')}
                strokeWidth={1}
                opacity={isCompleted ? 0.4 : (isActive ? 0.6 : 0.2)}
              />
              
              <circle
                r={nodeRadius}
                fill={isCompleted ? '#00c864' : (isActive ? '#ffffff' : 'transparent')}
                stroke={isCompleted ? '#00c864' : (isActive ? '#ffffff' : '#222222')}
                strokeWidth={isActive ? 2 : 1.5}
                filter={isActive ? 'url(#phase-glow)' : 'none'}
                style={{
                  animation: isActive ? 'pulse-glow 2s ease-in-out infinite' : 'none'
                }}
              />
              
              {isCompleted ? (
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#0a0a0a" strokeWidth="3">
                  <polyline points="20 6 9 17 4 12" className="checkmark" />
                </svg>
              ) : (
                <text
                  textAnchor="middle"
                  dominantBaseline="middle"
                  fontSize="14"
                  fontWeight={600}
                  fill={isActive ? '#0a0a0a' : '#888888'}
                >
                  {i + 1}
                </text>
              )}
              
              <text
                x={0}
                y={nodeRadius + 20}
                textAnchor="middle"
                fontSize="11"
                fontWeight={500}
                fill="#ffffff"
                style={{ maxWidth: '120px' }}
              >
                {phase.label}
              </text>
            </g>
          )
        })}

        <text
          x={centerX}
          y={centerY}
          textAnchor="middle"
          dominantBaseline="middle"
          fontSize="12"
          fontWeight={500}
          fill="#666666"
          style={{ pointerEvents: 'none' }}
        >
          {diagram.title}
        </text>
      </svg>
    </div>
  )
}

const VariantsView: React.FC<{ variants: VariantData[]; nodeLabel: string }> = ({ variants, nodeLabel }) => {
  if (variants.length === 0) {
    return (
      <div style={{ 
        height: '100%', 
        display: 'flex', 
        alignItems: 'center', 
        justifyContent: 'center',
        color: '#666666',
        fontSize: '14px'
      }}>
        No variants defined for {nodeLabel}
      </div>
    )
  }

  return (
    <div className="table-container" style={{ height: '100%' }}>
      <table className="table">
        <thead>
          <tr>
            <th style={{ width: '15%' }}>Variant</th>
            <th style={{ width: '25%' }}>Description</th>
            <th style={{ width: '20%' }}>Pros</th>
            <th style={{ width: '20%' }}>Cons</th>
            <th style={{ width: '20%' }}>Recommended For</th>
          </tr>
        </thead>
        <tbody>
          {variants.map((variant, idx) => (
            <tr key={variant.name} style={{ opacity: idx === 0 ? 1 : 0.7 }}>
              <td style={{ fontWeight: 600, color: '#ffffff' }}>{variant.name}</td>
              <td style={{ color: '#aaaaaa', fontSize: '11px' }}>{variant.description}</td>
              <td>
                <ul style={{ margin: 0, paddingLeft: '16px', fontSize: '11px' }}>
                  {variant.pros.slice(0, 2).map((pro, pIdx) => (
                    <li key={pIdx} style={{ color: '#00c864', margin: '2px 0' }}>{pro}</li>
                  ))}
                  {variant.pros.length > 2 && (
                    <li style={{ color: '#666', fontSize: '10px' }}>+{variant.pros.length - 2} more</li>
                  )}
                </ul>
              </td>
              <td>
                <ul style={{ margin: 0, paddingLeft: '16px', fontSize: '11px' }}>
                  {variant.cons.slice(0, 2).map((con, cIdx) => (
                    <li key={cIdx} style={{ color: '#ff4444', margin: '2px 0' }}>{con}</li>
                  ))}
                  {variant.cons.length > 2 && (
                    <li style={{ color: '#666', fontSize: '10px' }}>+{variant.cons.length - 2} more</li>
                  )}
                </ul>
              </td>
              <td>
                <div style={{ display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
                  {variant.recommendedFor.slice(0, 3).map((rec, rIdx) => (
                    <span key={rIdx} className="badge" style={{ fontSize: '10px' }}>{rec}</span>
                  ))}
                  {variant.recommendedFor.length > 3 && (
                    <span className="badge" style={{ fontSize: '10px', opacity: 0.6 }}>
                      +{variant.recommendedFor.length - 3}
                    </span>
                  )}
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

const LiveDemoView: React.FC<{ demo: LiveDemoData; nodeLabel: string }> = ({ demo, nodeLabel }) => {
  const [animatedValue, setAnimatedValue] = useState(0)

  useEffect(() => {
    let frame: number
    const animate = () => {
      setAnimatedValue(prev => (prev + 0.5) % 100)
      frame = requestAnimationFrame(animate)
    }
    frame = requestAnimationFrame(animate)
    return () => cancelAnimationFrame(frame)
  }, [])

  const statusColors = {
    connecting: '#ffb400',
    connected: '#00c864',
    disconnected: '#888888',
    error: '#ff4444'
  }

  return (
    <div style={{ height: '100%', display: 'flex', flexDirection: 'column', gap: '16px' }}>
      <div style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        padding: '12px 16px',
        background: 'rgba(255,255,255,0.03)',
        border: '1px solid #222222',
        borderRadius: '8px'
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <div
            className="indicator-dot"
            style={{
              width: '12px',
              height: '12px',
              background: statusColors[demo.connectionStatus],
              boxShadow: `0 0 12px ${statusColors[demo.connectionStatus]}`,
              animation: demo.connectionStatus === 'connected' ? 'pulse-glow 2s ease-in-out infinite' : 'none'
            }}
          />
          <div>
            <div style={{ fontWeight: 600, textTransform: 'capitalize' }}>
              {demo.connectionStatus}
            </div>
            <div style={{ fontSize: '11px', color: '#888' }}>{nodeLabel} Live Connection</div>
          </div>
        </div>
        <div style={{ display: 'flex', gap: '24px', fontSize: '13px' }}>
          <div>
            <span style={{ color: '#888' }}>Latency: </span>
            <span style={{ fontWeight: 600, color: '#00c864' }}>{demo.latencyMs}ms</span>
          </div>
          <div>
            <span style={{ color: '#888' }}>Throughput: </span>
            <span style={{ fontWeight: 600, color: '#ffffff' }}>{demo.messagesPerSecond.toLocaleString()}/s</span>
          </div>
        </div>
      </div>

      <div style={{ 
        flex: 1, 
        display: 'grid', 
        gridTemplateColumns: 'repeat(2, 1fr)', 
        gap: '12px',
        overflow: 'auto'
      }}>
        {demo.indicators.map((indicator) => (
          <div
            key={indicator.id}
            className="card"
            style={{
              padding: '16px',
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'center',
              transition: 'all 0.2s ease'
            }}
          >
            <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '8px' }}>
              <div
                className={`indicator-dot ${indicator.status} ${indicator.pulse ? 'pulse' : ''}`}
                style={{
                  width: '10px',
                  height: '10px',
                  animation: indicator.pulse ? 'pulse-glow 1.5s ease-in-out infinite' : 'none'
                }}
              />
              <span className="indicator-label">{indicator.label}</span>
            </div>
            <div className="indicator-value" style={{ fontSize: '24px', fontWeight: 700 }}>
              {indicator.value}
            </div>
          </div>
        ))}
      </div>

      <div style={{
        height: '80px',
        position: 'relative',
        background: 'rgba(255,255,255,0.02)',
        border: '1px solid #222222',
        borderRadius: '8px',
        overflow: 'hidden'
      }}>
        <svg width="100%" height="100%" style={{ display: 'block' }}>
          <defs>
            <linearGradient id="signalGradient" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset={`${animatedValue}%`} stopColor="#00c864" stopOpacity="1"/>
              <stop offset={`${(animatedValue + 10) % 100}%`} stopColor="#00c864" stopOpacity="0"/>
            </linearGradient>
          </defs>
          <rect width="100%" height="100%" fill="url(#signalGradient)" opacity="0.3"/>
          {[0, 20, 40, 60, 80].map((offset, idx) => (
            <circle
              key={idx}
              cx={`${(animatedValue + offset) % 100}%`}
              cy="50%"
              r="4"
              fill="#00c864"
              opacity="0.8"
              filter="url(#phase-glow)"
            />
          ))}
        </svg>
      </div>
    </div>
  )
}

const NodeDetailsCard: React.FC<{ 
  node: NodeData; 
  connectedNodes: NodeData[]
  highlightedPath: NodeId[]
}> = ({ node, connectedNodes, highlightedPath }) => {
  const details = node.details

  return (
    <div style={{ 
      display: 'grid', 
      gridTemplateColumns: 'repeat(2, 1fr)', 
      gap: '16px',
      height: '100%'
    }}>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        <DetailCard title="What" content={details.what} />
        <DetailCard title="How" content={details.how} />
        <DetailCard title="When" content={details.when} />
        <DetailCard title="Where" content={details.where} />
        <DetailCard title="Why" content={details.why} />
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
        <InfoCard title="Product Stage">
          <span className={`badge badge-${details.productStage}`}>
            {details.productStage.charAt(0).toUpperCase() + details.productStage.slice(1)}
          </span>
        </InfoCard>

        <InfoCard title="User/RPS Thresholds">
          <div style={{ color: '#ffffff', fontFamily: 'monospace', fontSize: '13px' }}>
            {details.userRPSThresholds}
          </div>
        </InfoCard>

        <InfoCard title="Live Cloud Costs">
          <div style={{ color: '#ffffff', fontFamily: 'monospace', fontSize: '13px' }}>
            {details.liveCloudCosts}
          </div>
        </InfoCard>

        {connectedNodes.length > 0 && (
          <InfoCard title="Connected Nodes" style={{ flex: 1 }}>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
              {connectedNodes.map(n => (
                <span
                  key={n.id}
                  className="badge"
                  style={{
                    background: highlightedPath.includes(n.id) ? 'rgba(255,255,255,0.15)' : 'rgba(255,255,255,0.05)',
                    borderColor: highlightedPath.includes(n.id) ? '#ffffff' : '#222222',
                    color: highlightedPath.includes(n.id) ? '#ffffff' : '#888888',
                    fontWeight: highlightedPath.includes(n.id) ? 600 : 500
                  }}
                >
                  {n.label}
                </span>
              ))}
            </div>
          </InfoCard>
        )}
      </div>
    </div>
  )
}

const DetailCard: React.FC<{ title: string; content: string }> = ({ title, content }) => (
  <div className="card" style={{ padding: '16px', flex: 1 }}>
    <div className="section-title">{title}</div>
    <div style={{ color: '#cccccc', lineHeight: 1.6, fontSize: '13px' }}>
      {content}
    </div>
  </div>
)

const InfoCard: React.FC<{ title: string; children: React.ReactNode; style?: React.CSSProperties }> = ({ 
  title, 
  children, 
  style 
}) => (
  <div className="card" style={{ padding: '16px', ...style }}>
    <div className="section-title">{title}</div>
    {children}
  </div>
)

interface CarouselProps {
  currentIndex: number
  onNext: () => void
  onPrev: () => void
  onDotClick: (index: number) => void
  node: NodeData
  variants: VariantData[]
  liveDemo: LiveDemoData | undefined
  flowDiagram: FlowDiagramData | undefined
}

const Carousel: React.FC<CarouselProps> = ({ 
  currentIndex, 
  onNext, 
  onPrev, 
  onDotClick,
  node,
  variants,
  liveDemo,
  flowDiagram
}) => {
  const carouselItems = [
    {
      title: 'Architecture Flow',
      render: () => flowDiagram ? <FlowDiagramView diagram={flowDiagram} /> : null
    },
    {
      title: 'Variants',
      render: () => variants.length > 0 ? <VariantsView variants={variants} nodeLabel={node.label} /> : null
    },
    {
      title: 'Live Demo',
      render: () => liveDemo ? <LiveDemoView demo={liveDemo} nodeLabel={node.label} /> : null
    }
  ].filter(item => item.render() !== null)

  const visibleItems = 3
  const totalItems = carouselItems.length
  
  if (totalItems === 0) {
    return (
      <div style={{ 
        flex: 1, 
        display: 'flex', 
        alignItems: 'center', 
        justifyContent: 'center',
        color: '#666666'
      }}>
        No carousel content available
      </div>
    )
  }

  return (
    <div style={{ flex: 1, display: 'flex', flexDirection: 'column', height: '100%' }}>
      <div style={{ 
        display: 'flex', 
        alignItems: 'center', 
        justifyContent: 'space-between',
        marginBottom: '12px',
        padding: '0 4px'
      }}>
        <h3 style={{ 
          fontSize: '14px', 
          fontWeight: 600, 
          color: '#ffffff',
          margin: 0 
        }}>
          {carouselItems[currentIndex].title}
        </h3>
        <div style={{ display: 'flex', gap: '8px' }}>
          <button
            className="carousel-arrow"
            onClick={onPrev}
            disabled={totalItems <= 1}
            aria-label="Previous"
            style={{ width: '32px', height: '32px' }}
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <polyline points="15 18 9 12 15 6"/>
            </svg>
          </button>
          <button
            className="carousel-arrow"
            onClick={onNext}
            disabled={totalItems <= 1}
            aria-label="Next"
            style={{ width: '32px', height: '32px' }}
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
          </button>
        </div>
      </div>

      <div style={{ 
        flex: 1, 
        overflow: 'hidden', 
        position: 'relative',
        borderRadius: '8px',
        background: 'rgba(255,255,255,0.02)',
        border: '1px solid #222222'
      }}>
        <div 
          style={{ 
            display: 'flex', 
            height: '100%', 
            transition: 'transform 0.4s cubic-bezier(0.16, 1, 0.3, 1)',
            transform: `translateX(-${currentIndex * (100 / visibleItems)}%)`
          }}
        >
          {carouselItems.map((item, _index) => (
            <div
              key={item.title}
              style={{ 
                width: `${100 / visibleItems}%`, 
                flexShrink: 0, 
                height: '100%',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                padding: '16px'
              }}
            >
              {item.render()}
            </div>
          ))}
        </div>
      </div>

      <div className="carousel-dots" style={{ marginTop: '12px' }}>
        {carouselItems.map((_, index) => (
          <button
            key={index}
            className={`carousel-dot ${index === currentIndex ? 'active' : ''}`}
            onClick={() => onDotClick(index)}
            aria-label={`Go to ${carouselItems[index].title}`}
            aria-selected={index === currentIndex}
          />
        ))}
      </div>
    </div>
  )
}