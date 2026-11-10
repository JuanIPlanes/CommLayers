import { GraphData, NodeId, EdgeId } from '../types'

const createNodeId = (id: string): NodeId => id as NodeId
const createEdgeId = (id: string): EdgeId => id as EdgeId

export const mockGraphData: GraphData = {
  nodes: [
    // Technical Connections Mode - Paradigms/Transports
    {
      id: createNodeId('request-response'),
      label: 'Request–Response',
      type: 'paradigm',
      description: 'Classic synchronous request-response pattern over HTTP/HTTPS. Client initiates, server responds, connection closes.',
      details: {
        what: 'Synchronous HTTP request-response cycle where client waits for server response before continuing.',
        how: 'Client opens TCP connection, sends HTTP request, server processes and returns response, connection terminates.',
        when: 'REST APIs, CRUD operations, form submissions, any operation requiring immediate response.',
        where: 'Web applications, mobile backends, microservices synchronous communication.',
        why: 'Simple mental model, universal tooling, built-in caching, retry semantics, firewall-friendly.',
        productStage: 'production',
        userRPSThresholds: '10K–100K RPS with horizontal scaling',
        liveCloudCosts: 'AWS/GCP $0.002/h per 1K RPS baseline'
      },
      flowDiagram: {
        title: 'Request–Response Flow',
        phases: [
          { id: '1', label: 'DNS Resolution', description: 'Resolve hostname to IP address', status: 'completed' },
          { id: '2', label: 'TCP Handshake', description: 'SYN → SYN-ACK → ACK (1 RTT)', status: 'completed' },
          { id: '3', label: 'TLS Handshake', description: 'ClientHello → ServerHello → Cert → Finished (1–2 RTT)', status: 'completed' },
          { id: '4', label: 'HTTP Request', description: 'GET/POST with headers + body', status: 'completed' },
          { id: '5', label: 'Server Processing', description: 'Route → Auth → Business Logic → DB', status: 'completed' },
          { id: '6', label: 'HTTP Response', description: 'Status + Headers + Body', status: 'completed' },
          { id: '7', label: 'Connection Close', description: 'FIN → FIN-ACK → ACK (or keep-alive)', status: 'completed' }
        ]
      },
      variants: [
        {
          name: 'HTTP/1.1',
          description: 'Text-based, head-of-line blocking, one request per connection',
          pros: ['Universal support', 'Simple debugging', 'Built-in caching'],
          cons: ['HOL blocking', 'Header overhead', 'No multiplexing'],
          recommendedFor: ['Simple APIs', 'Low concurrency', 'Legacy clients'],
          avoidWhen: ['High latency links', 'Many small requests', 'Real-time needs']
        },
        {
          name: 'HTTP/2',
          description: 'Binary framing, multiplexing, header compression, server push',
          pros: ['Multiplexing', 'HPACK compression', 'Priority streams'],
          cons: ['TCP HOL remains', 'Complex debugging', 'Priority issues'],
          recommendedFor: ['Modern web apps', 'Asset-heavy pages', 'API gateways'],
          avoidWhen: ['High packet loss', 'Simple use cases', 'Non-browser clients']
        },
        {
          name: 'HTTP/3 (QUIC)',
          description: 'UDP-based, 0-RTT, native multiplexing, connection migration',
          pros: ['No TCP HOL', '0-RTT resume', 'Connection migration'],
          cons: ['UDP blocking', 'CPU overhead', 'Newer ecosystem'],
          recommendedFor: ['Mobile networks', 'Global latency', 'Real-time APIs'],
          avoidWhen: ['Corporate firewalls', 'Legacy infra', 'CPU-constrained']
        }
      ],
      liveDemo: {
        connectionStatus: 'connected',
        latencyMs: 42,
        messagesPerSecond: 1250,
        indicators: [
          { id: '1', label: 'p50 Latency', value: '42ms', status: 'ok', pulse: true },
          { id: '2', label: 'p99 Latency', value: '187ms', status: 'ok' },
          { id: '3', label: 'Error Rate', value: '0.02%', status: 'ok' },
          { id: '4', label: 'Throughput', value: '1.25K req/s', status: 'ok', pulse: true }
        ]
      },
      position: { x: 200, y: 400 }
    },
    {
      id: createNodeId('websockets'),
      label: 'WebSockets',
      type: 'transport',
      description: 'Full-duplex bidirectional communication over a single TCP connection. Enables real-time server-to-client push.',
      details: {
        what: 'Persistent full-duplex TCP connection upgraded from HTTP, enabling bidirectional message flow without polling.',
        how: 'HTTP Upgrade handshake → WebSocket frames (text/binary) → Ping/Pong keepalive → Close handshake.',
        when: 'Real-time dashboards, collaborative editing, gaming, chat, live trading, IoT telemetry.',
        where: 'Browser clients, mobile apps, server-to-server real-time links, edge messaging.',
        why: 'True push, minimal overhead after handshake, binary support, native browser API.',
        productStage: 'production',
        userRPSThresholds: '50K–500K concurrent connections per node',
        liveCloudCosts: 'AWS/GCP $0.009/h per 10K concurrent connections'
      },
      flowDiagram: {
        title: 'WebSocket Connection Flow',
        phases: [
          { id: '1', label: 'HTTP Upgrade Request', description: 'GET /ws + Upgrade: websocket + Sec-WebSocket-Key', status: 'completed', animationDelay: 0 },
          { id: '2', label: 'Server Accepts Upgrade', description: '101 Switching Protocols + Sec-WebSocket-Accept', status: 'completed', animationDelay: 300 },
          { id: '3', label: 'TCP Duplex Established', description: 'Full-duplex channel open, bidirectional frames', status: 'completed', animationDelay: 600 },
          { id: '4', label: 'Bidirectional Messaging', description: 'Client ↔ Server: text/binary frames, ping/pong', status: 'active', animationDelay: 900 },
          { id: '5', label: 'Heartbeat & Keepalive', description: 'Ping/Pong every 30s, detect dead connections', status: 'pending', animationDelay: 1200 },
          { id: '6', label: 'Graceful Close', description: 'Close frame (1000) → Acknowledge → TCP FIN', status: 'pending', animationDelay: 1500 }
        ]
      },
      variants: [
        {
          name: 'Raw WebSocket',
          description: 'Native WebSocket API, manual reconnection, no fallback',
          pros: ['Zero dependencies', 'Full control', 'Binary support'],
          cons: ['No auto-reconnect', 'No fallback', 'Manual heartbeats'],
          recommendedFor: ['Simple real-time', 'Controlled environments', 'Custom protocols'],
          avoidWhen: ['Unreliable networks', 'Mobile clients', 'Production systems']
        },
        {
          name: 'Socket.io',
          description: 'Auto-reconnect, fallback transports, rooms, acknowledgments',
          pros: ['Auto-reconnect', 'HTTP long-poll fallback', 'Room/namespace'],
          cons: ['Larger bundle', 'Proprietary protocol', 'Server coupling'],
          recommendedFor: ['Rapid prototyping', 'Complex rooms', 'Mixed clients'],
          avoidWhen: ['Bundle size critical', 'Native WS needed', 'High scale']
        },
        {
          name: 'uWebSockets.js',
          description: 'High-performance C++ core, minimal overhead, native Node.js',
          pros: ['~10x faster', 'Tiny memory', 'Native performance'],
          cons: ['Native build', 'Less features', 'Smaller community'],
          recommendedFor: ['High-frequency trading', 'Gaming servers', 'Millions of conn'],
          avoidWhen: ['Rapid iteration', 'Feature-rich needs', 'Windows-only CI']
        }
      ],
      liveDemo: {
        connectionStatus: 'connected',
        latencyMs: 3,
        messagesPerSecond: 45000,
        indicators: [
          { id: '1', label: 'Round-trip', value: '3ms', status: 'ok', pulse: true },
          { id: '2', label: 'Connections', value: '47,231', status: 'ok', pulse: true },
          { id: '3', label: 'Msg/sec', value: '45.2K', status: 'ok', pulse: true },
          { id: '4', label: 'Memory', value: '234 MB', status: 'ok' }
        ]
      },
      position: { x: 600, y: 400 }
    },
    {
      id: createNodeId('realtime-sync-engine'),
      label: 'Realtime Sync Engine',
      type: 'pattern',
      description: 'Conflict-free replicated data types (CRDTs) and operational transformation for multi-user real-time synchronization.',
      details: {
        what: 'Distributed synchronization engine using CRDTs/OT to merge concurrent edits without central coordination.',
        how: 'Local-first mutations → Intent broadcast → Peer merge (CRDT/OT) → Convergence guarantee → Conflict resolution.',
        when: 'Collaborative editors (Figma, Notion), offline-first apps, multi-device sync, distributed teams.',
        where: 'Client-side (IndexedDB), edge workers, peer-to-peer (WebRTC), central relay servers.',
        why: 'Offline-capable, no single point of failure, eventual consistency, user intent preservation.',
        productStage: 'beta',
        userRPSThresholds: '10K–100K concurrent editing sessions',
        liveCloudCosts: 'AWS/GCP $0.015/h per 1K active sessions (includes relay)'
      },
      flowDiagram: {
        title: 'Realtime Sync Flow',
        phases: [
          { id: '1', label: 'Local Mutation', description: 'User edits → Local CRDT op generated', status: 'completed' },
          { id: '2', label: 'Optimistic UI', description: 'Immediate local reflect, no waiting', status: 'completed' },
          { id: '3', label: 'Broadcast Intent', description: 'WebSocket/WebRTC → Peers + Relay', status: 'active' },
          { id: '4', label: 'Peer Merge', description: 'Receive ops → CRDT merge → Converge', status: 'pending' },
          { id: '5', label: 'Ack & Confirm', description: 'Server acknowledges → Version vector', status: 'pending' },
          { id: '6', label: 'Conflict Resolution', description: 'Auto-merge or user prompt on semantic conflict', status: 'pending' }
        ]
      },
      variants: [
        {
          name: 'Yjs (CRDT)',
          description: 'Mature CRDT implementation, rich ecosystem, WebRTC support',
          pros: ['Battle-tested', 'Many providers', 'Offline-first', 'Rich text'],
          cons: ['Larger bundle', 'Learning curve', 'Memory growth'],
          recommendedFor: ['Rich text', 'Complex structures', 'P2P needed'],
          avoidWhen: ['Simple state', 'Bundle critical', 'Custom merge logic']
        },
        {
          name: 'Automerge',
          description: 'JSON-like CRDT, immutable, functional approach',
          pros: ['JSON-like', 'Immutable', 'Time-travel', 'WASM fast'],
          cons: ['Memory heavy', 'No rich text', 'Slower large docs'],
          recommendedFor: ['JSON documents', 'Audit trails', 'Functional apps'],
          avoidWhen: ['Rich text', 'High frequency', 'Mobile memory']
        },
        {
          name: 'Custom OT',
          description: 'Operational transformation tailored to domain',
          pros: ['Domain-optimized', 'Minimal overhead', 'Full control'],
          cons: ['Complex correctness', 'Centralized', 'Hard to test'],
          recommendedFor: ['Simple linear', 'Proven algorithms', 'Centralized OK'],
          avoidWhen: ['P2P needed', 'Complex structures', 'Team inexperience']
        }
      ],
      liveDemo: {
        connectionStatus: 'connected',
        latencyMs: 12,
        messagesPerSecond: 8500,
        indicators: [
          { id: '1', label: 'Sync Latency', value: '12ms', status: 'ok', pulse: true },
          { id: '2', label: 'Active Peers', value: '23', status: 'ok' },
          { id: '3', label: 'Ops/sec', value: '8.5K', status: 'ok', pulse: true },
          { id: '4', label: 'Conflicts', value: '0', status: 'ok' }
        ]
      },
      position: { x: 1000, y: 400 }
    },
    // Additional nodes for the graph
    {
      id: createNodeId('sse'),
      label: 'Server-Sent Events',
      type: 'transport',
      description: 'Unidirectional server-to-client streaming over HTTP. Simpler than WebSockets for server-push only.',
      details: {
        what: 'HTTP-based server push with automatic reconnection, text/event-stream format, browser-native EventSource API.',
        how: 'Client GET /events → Server holds connection → Streaming text/event-stream → Auto-reconnect on drop.',
        when: 'Notifications, live feeds, progress updates, one-way real-time, fallback for WebSockets.',
        where: 'Browser EventSource, polyfills for Node, edge functions, CDN-friendly.',
        why: 'Simpler than WS, HTTP/2 multiplexing, auto-reconnect, firewall-friendly, no binary.',
        productStage: 'production',
        userRPSThresholds: '100K–1M connections (stateless, HTTP/2)',
        liveCloudCosts: 'AWS/GCP $0.003/h per 10K connections'
      },
      flowDiagram: {
        title: 'SSE Flow',
        phases: [
          { id: '1', label: 'EventSource Connect', description: 'GET /events, Accept: text/event-stream', status: 'completed' },
          { id: '2', label: 'Server Streams', description: 'data: {...}\n\n (keep connection open)', status: 'active' },
          { id: '3', label: 'Auto-Reconnect', description: 'Browser handles retry, Last-Event-ID resume', status: 'pending' },
          { id: '4', label: 'Heartbeat', description: 'Comment lines : keepalive every 15s', status: 'pending' }
        ]
      },
      variants: [
        {
          name: 'Native EventSource',
          description: 'Browser built-in, zero dependencies, automatic reconnection',
          pros: ['Zero bytes', 'Auto-reconnect', 'Simple API'],
          cons: ['No binary', 'GET only', 'No custom headers', 'Connection limits'],
          recommendedFor: ['Simple notifications', 'Progress updates', 'Lightweight'],
          avoidWhen: ['Binary data', 'Auth headers', 'High frequency', 'Cross-origin complexity']
        },
        {
          name: 'Fetch + ReadableStream',
          description: 'Modern streaming fetch, full control, binary support',
          pros: ['Binary OK', 'Custom headers', 'AbortController', 'Stream API'],
          cons: ['Manual reconnect', 'More code', 'No IE11'],
          recommendedFor: ['Binary streams', 'Auth required', 'Custom retry', 'Modern only'],
          avoidWhen: ['Legacy support', 'Quick prototyping', 'Simple text events']
        }
      ],
      liveDemo: {
        connectionStatus: 'connected',
        latencyMs: 18,
        messagesPerSecond: 3200,
        indicators: [
          { id: '1', label: 'Latency', value: '18ms', status: 'ok' },
          { id: '2', label: 'Connections', value: '12,450', status: 'ok' },
          { id: '3', label: 'Events/sec', value: '3.2K', status: 'ok' },
          { id: '4', label: 'Reconnects', value: '0', status: 'ok' }
        ]
      },
      position: { x: 600, y: 200 }
    },
    {
      id: createNodeId('graphql-subscriptions'),
      label: 'GraphQL Subscriptions',
      type: 'protocol',
      description: 'Real-time GraphQL operations over WebSocket. Subscribe to mutations, receive pushed updates.',
      details: {
        what: 'GraphQL subscription operations delivered via WebSocket transport, integrated with schema and resolvers.',
        how: 'WS connection → GraphQL-WS protocol → Subscribe mutation → Server pushes on trigger → Client updates cache.',
        when: 'GraphQL backends, real-time UI updates, reactive data fetching, Apollo/Relay ecosystems.',
        where: 'Apollo Server, Yoga, Hasura, GraphQL Yoga, any GraphQL server with subscription support.',
        why: 'Schema-driven, typed subscriptions, cache integration, unified query/mutation/subscription.',
        productStage: 'production',
        userRPSThresholds: '10K–50K subscriptions per instance',
        liveCloudCosts: 'AWS/GCP $0.008/h per 1K active subscriptions'
      },
      flowDiagram: {
        title: 'GraphQL Subscription Flow',
        phases: [
          { id: '1', label: 'WS Connect', description: 'GraphQL-WS handshake, connectionInit', status: 'completed' },
          { id: '2', label: 'Subscribe Op', description: 'subscription { onMessage { id text } }', status: 'completed' },
          { id: '3', label: 'Server Registers', description: 'Resolver + PubSub → AsyncIterator', status: 'completed' },
          { id: '4', label: 'Mutation Fires', description: 'publish() → Iterator.next() → Push', status: 'active' },
          { id: '5', label: 'Client Receives', description: 'next: { data: {...} } → Cache update', status: 'pending' },
          { id: '6', label: 'Complete/Error', description: 'Complete → Unsubscribe, cleanup', status: 'pending' }
        ]
      },
      variants: [
        {
          name: 'graphql-ws (v5+)',
          description: 'Modern protocol, no legacy, TypeScript-first',
          pros: ['Active maintenance', 'Typed', 'Lite', 'No subscriptions-transport-ws'],
          cons: ['Newer', 'Breaking from v4'],
          recommendedFor: ['New projects', 'TypeScript', 'Apollo Client 3.8+'],
          avoidWhen: ['Legacy codebases', 'Old Apollo versions']
        },
        {
          name: 'subscriptions-transport-ws',
          description: 'Legacy protocol, wider server support',
          pros: ['Universal support', 'Older servers', 'Apollo v2/v3'],
          cons: ['Deprecated', 'No active maint', 'Heavier'],
          recommendedFor: ['Legacy migration', 'Hasura', 'Old tooling'],
          avoidWhen: ['Greenfield', 'Modern stack']
        }
      ],
      liveDemo: {
        connectionStatus: 'connected',
        latencyMs: 8,
        messagesPerSecond: 15000,
        indicators: [
          { id: '1', label: 'Sub Latency', value: '8ms', status: 'ok', pulse: true },
          { id: '2', label: 'Active Subs', value: '8,921', status: 'ok' },
          { id: '3', label: 'Publish/sec', value: '15K', status: 'ok', pulse: true },
          { id: '4', label: 'Cache Hits', value: '94%', status: 'ok' }
        ]
      },
      position: { x: 1000, y: 200 }
    },
    {
      id: createNodeId('grpc-streaming'),
      label: 'gRPC Streaming',
      type: 'protocol',
      description: 'HTTP/2-based bidirectional, client, or server streaming with protobuf serialization.',
      details: {
        what: 'Native HTTP/2 streaming RPCs: unary, server-streaming, client-streaming, bidirectional streaming.',
        how: 'HTTP/2 frames → Protobuf messages → Flow control (WINDOW_UPDATE) → Trailers for status.',
        when: 'Microservices, high-throughput internal APIs, polyglot environments, service mesh.',
        where: 'Kubernetes, service mesh (Istio/Linkerd), gRPC Gateway for HTTP/JSON, edge proxies.',
        why: 'Contract-first, efficient binary, flow control, codegen, interception, load balancing.',
        productStage: 'production',
        userRPSThresholds: '1M+ RPS per instance (bidirectional streaming)',
        liveCloudCosts: 'AWS/GCP $0.005/h per 10K RPS (internal)'
      },
      flowDiagram: {
        title: 'gRPC Bidirectional Streaming',
        phases: [
          { id: '1', label: 'HTTP/2 Connect', description: 'PRI + SETTINGS, stream ID allocation', status: 'completed' },
          { id: '2', label: 'Metadata Headers', description: ':method=POST, content-type=grpc', status: 'completed' },
          { id: '3', label: 'Stream Open', description: 'Client → Server: data frames (length-prefixed)', status: 'active' },
          { id: '4', label: 'Server Streams Back', description: 'Server → Client: interleaved data frames', status: 'active' },
          { id: '5', label: 'Flow Control', description: 'WINDOW_UPDATE per stream + connection', status: 'pending' },
          { id: '6', label: 'Trailers + Close', description: 'grpc-status, grpc-message, RST_STREAM', status: 'pending' }
        ]
      },
      variants: [
        {
          name: 'Unary',
          description: 'Single request → single response (standard RPC)',
          pros: ['Simple', 'Cached', 'Retriable', 'Load balanced'],
          cons: ['No streaming', 'RTT per call'],
          recommendedFor: ['CRUD', 'Queries', 'Commands'],
          avoidWhen: ['Streaming needed', 'Large payloads', 'Real-time']
        },
        {
          name: 'Server Streaming',
          description: 'Client sends request, server streams multiple responses',
          pros: ['Push from server', 'Flow controlled', 'Efficient'],
          cons: ['Client cannot stream', 'One-way'],
          recommendedFor: ['Large results', 'Progress', 'Feeds'],
          avoidWhen: ['Bidirectional needed', 'Client upload']
        },
        {
          name: 'Client Streaming',
          description: 'Client streams requests, server sends single response',
          pros: ['Upload large', 'Chunked send', 'Flow control'],
          cons: ['No server push', 'One response'],
          recommendedFor: ['Bulk insert', 'File upload', 'Log ingest'],
          avoidWhen: ['Server push', 'Interactive']
        },
        {
          name: 'Bidirectional',
          description: 'Both sides stream independently, full duplex',
          pros: ['True duplex', 'Independent', 'Low latency'],
          cons: ['Complex', 'Stateful', 'Harder LB'],
          recommendedFor: ['Chat', 'Gaming', 'Telemetry', 'Sync'],
          avoidWhen: ['Simple RPC', 'Stateless needed']
        }
      ],
      liveDemo: {
        connectionStatus: 'connected',
        latencyMs: 2,
        messagesPerSecond: 250000,
        indicators: [
          { id: '1', label: 'RTT', value: '2ms', status: 'ok', pulse: true },
          { id: '2', label: 'Streams', value: '1,234', status: 'ok' },
          { id: '3', label: 'Msg/sec', value: '250K', status: 'ok', pulse: true },
          { id: '4', label: 'Flow Ctrl', value: 'OK', status: 'ok' }
        ]
      },
      position: { x: 1400, y: 400 }
    },
    {
      id: createNodeId('mqtt'),
      label: 'MQTT',
      type: 'protocol',
      description: 'Lightweight pub/sub messaging for IoT. QoS 0/1/2, retained messages, last will, topic wildcards.',
      details: {
        what: 'Minimal pub/sub protocol over TCP (or WebSocket), designed for constrained devices and unreliable networks.',
        how: 'CONNECT → SUBSCRIBE (topic filters) → PUBLISH (QoS 0/1/2) → PINGREQ/PINGRESP → DISCONNECT.',
        when: 'IoT sensors, mobile push, telemetry, home automation, edge computing, satellite links.',
        where: 'Mosquitto, EMQX, VerneMQ, AWS IoT Core, Azure IoT Hub, HiveMQ, cloud brokers.',
        why: 'Tiny header (2 bytes), QoS guarantees, retained/last-will, topic hierarchy, massive scale.',
        productStage: 'production',
        userRPSThresholds: '10M+ messages/sec on clustered brokers',
        liveCloudCosts: 'AWS/GCP $0.004/h per 10K connected devices'
      },
      flowDiagram: {
        title: 'MQTT Pub/Sub Flow',
        phases: [
          { id: '1', label: 'CONNECT', description: 'ClientID, clean session, keepalive, will', status: 'completed' },
          { id: '2', label: 'CONNACK', description: 'Session present, return code', status: 'completed' },
          { id: '3', label: 'SUBSCRIBE', description: 'Topic filters + QoS (wildcards +/#)', status: 'completed' },
          { id: '4', label: 'PUBLISH', description: 'Topic + Payload + QoS + Retain', status: 'active' },
          { id: '5', label: 'QoS Flow', description: 'QoS1: PUBACK, QoS2: PUBREC→PUBREL→PUBCOMP', status: 'pending' },
          { id: '6', label: 'Retained/Will', description: 'Last known state, LWT on ungraceful disconnect', status: 'pending' }
        ]
      },
      variants: [
        {
          name: 'MQTT 3.1.1',
          description: 'Classic version, widely deployed, basic features',
          pros: ['Universal', 'Simple', 'Mature brokers'],
          cons: ['No shared subscriptions', 'Limited auth', 'No negative ack'],
          recommendedFor: ['Legacy devices', 'Simple IoT', 'Existing deployments'],
          avoidWhen: ['High scale', 'Advanced routing', 'Modern security']
        },
        {
          name: 'MQTT 5.0',
          description: 'Modern version, shared subs, properties, reason codes',
          pros: ['Shared subscriptions', 'User properties', 'Reason codes', 'Topic alias'],
          cons: ['Newer brokers', 'Slightly complex', 'Migration effort'],
          recommendedFor: ['New IoT', 'High scale', 'Advanced routing', 'K8s native'],
          avoidWhen: ['Legacy devices only', '3.1.1 locked']
        },
        {
          name: 'MQTT over WebSocket',
          description: 'MQTT framed in WebSocket for browser connectivity',
          pros: ['Browser native', 'Firewall friendly', 'Single port'],
          cons: ['WS overhead', 'Binary framing', 'Proxy issues'],
          recommendedFor: ['Web dashboards', 'Browser IoT', 'Hybrid apps'],
          avoidWhen: ['Native TCP OK', 'Ultra-constrained']
        }
      ],
      liveDemo: {
        connectionStatus: 'connected',
        latencyMs: 15,
        messagesPerSecond: 180000,
        indicators: [
          { id: '1', label: 'Latency', value: '15ms', status: 'ok' },
          { id: '2', label: 'Devices', value: '234K', status: 'ok', pulse: true },
          { id: '3', label: 'Msg/sec', value: '180K', status: 'ok', pulse: true },
          { id: '4', label: 'QoS2 Pending', value: '0', status: 'ok' }
        ]
      },
      position: { x: 200, y: 600 }
    },
    {
      id: createNodeId('webrtc-datachannel'),
      label: 'WebRTC DataChannel',
      type: 'transport',
      description: 'Peer-to-peer arbitrary data transfer with optional reliability, ordering, and congestion control.',
      details: {
        what: 'SCTP-over-DTLS-over-UDP data channels for P2P communication, configurable reliability/ordered delivery.',
        how: 'Offer/Answer (SDP) → ICE candidates → DTLS handshake → SCTP association → DataChannel open.',
        when: 'P2P file transfer, gaming, video chat data, decentralized apps, mesh networks.',
        where: 'Browsers (native), Node.js (wrtc), React Native, Flutter, Electron, native mobile.',
        why: 'True P2P, no server relay (after signaling), sub-100ms, congestion control, binary.',
        productStage: 'beta',
        userRPSThresholds: 'Mesh: 10–50 peers, SFU: 100–500, MCU: 10–20',
        liveCloudCosts: 'AWS/GCP $0.02/h per TURN relay GB (signaling negligible)'
      },
      flowDiagram: {
        title: 'WebRTC DataChannel Flow',
        phases: [
          { id: '1', label: 'Signaling', description: 'Offer/Answer exchange via WS/HTTP', status: 'completed' },
          { id: '2', label: 'ICE Gathering', description: 'Host/SRFLX/RELAY candidates', status: 'completed' },
          { id: '3', label: 'DTLS Handshake', description: 'Fingerprint verify, key derivation', status: 'completed' },
          { id: '4', label: 'SCTP Association', description: 'INIT/INIT-ACK, stream negotiation', status: 'active' },
          { id: '5', label: 'DataChannel Open', description: 'Binary/text, ordered/unordered, reliable/unreliable', status: 'active' },
          { id: '6', label: 'P2P Transfer', description: 'Direct UDP, congestion control (GCC)', status: 'pending' }
        ]
      },
      variants: [
        {
          name: 'Reliable Ordered',
          description: 'TCP-like: guaranteed delivery, in-order, retransmission',
          pros: ['Guaranteed', 'Ordered', 'Flow control'],
          cons: ['Head-of-line', 'Latency on loss', 'More overhead'],
          recommendedFor: ['File transfer', 'State sync', 'Critical data'],
          avoidWhen: ['Real-time gaming', 'Telemetry', 'Low latency critical']
        },
        {
          name: 'Unreliable Unordered',
          description: 'UDP-like: fire-and-forget, no retransmit, no ordering',
          pros: ['Lowest latency', 'No HOL', 'Minimal overhead'],
          cons: ['Loss possible', 'No ordering', 'App handles'],
          recommendedFor: ['Gaming', 'Voice/Video', 'Telemetry', 'Position updates'],
          avoidWhen: ['Reliability needed', 'File transfer', 'State sync']
        },
        {
          name: 'Partial Reliability',
          description: 'SCTP PR-SCTP: timed reliability, limited retransmissions',
          pros: ['Tunable', 'Best of both', 'Deadline-aware'],
          cons: ['Complex', 'Less support', 'SCTP only'],
          recommendedFor: ['Gaming with critical', 'Adaptive quality', 'Mixed traffic'],
          avoidWhen: ['Simple needs', 'Non-SCTP stacks']
        }
      ],
      liveDemo: {
        connectionStatus: 'connecting',
        latencyMs: 45,
        messagesPerSecond: 1200,
        indicators: [
          { id: '1', label: 'RTT', value: '45ms', status: 'ok', pulse: true },
          { id: '2', label: 'Peers', value: '3', status: 'warning' },
          { id: '3', label: 'DataChans', value: '2', status: 'ok' },
          { id: '4', label: 'Packet Loss', value: '0.3%', status: 'warning' }
        ]
      },
      position: { x: 1000, y: 600 }
    },
    {
      id: createNodeId('event-driven'),
      label: 'Event-Driven Architecture',
      type: 'paradigm',
      description: 'Decoupled services communicate via immutable events. Event sourcing, CQRS, saga patterns.',
      details: {
        what: 'Services emit/consume immutable events. Event store as source of truth. Projections build read models.',
        how: 'Command → Validation → Event → Event Store → Projections → Read Models → Subscribers react.',
        when: 'Complex domains, audit requirements, temporal queries, microservices decoupling, compliance.',
        where: 'Kafka, EventStoreDB, NATS JetStream, Redis Streams, AWS EventBridge, Azure Event Hubs.',
        why: 'Complete audit trail, time-travel debugging, replay capability, loose coupling, scalability.',
        productStage: 'production',
        userRPSThresholds: '1M+ events/sec (Kafka), 100K+ (EventStoreDB)',
        liveCloudCosts: 'AWS/GCP $0.012/h per 10K events/sec (Kafka managed)'
      },
      flowDiagram: {
        title: 'Event-Driven Flow',
        phases: [
          { id: '1', label: 'Command Received', description: 'API Gateway → Command Handler', status: 'completed' },
          { id: '2', label: 'Validation', description: 'Business rules, invariants checked', status: 'completed' },
          { id: '3', label: 'Event Emitted', description: 'Immutable fact appended to log', status: 'active' },
          { id: '4', label: 'Projections Update', description: 'Read models rebuilt asynchronously', status: 'pending' },
          { id: '5', label: 'Downstream React', description: 'Consumers process, emit new events', status: 'pending' },
          { id: '6', label: 'Saga Coordination', description: 'Compensating transactions on failure', status: 'pending' }
        ]
      },
      variants: [
        {
          name: 'Event Sourcing',
          description: 'State = fold(events). Full history, time travel, rebuild anytime',
          pros: ['Complete history', 'Time travel', 'Rebuild projections', 'Audit native'],
          cons: ['Event schema migration', 'Snapshot strategy', 'Learning curve', 'Storage growth'],
          recommendedFor: ['Financial', 'Compliance', 'Complex domains', 'Audit critical'],
          avoidWhen: ['Simple CRUD', 'Team inexperience', 'Low audit needs']
        },
        {
          name: 'CQRS',
          description: 'Separate write model (commands) from read model (queries)',
          pros: ['Optimized reads', 'Scalable independently', 'Different models', 'Polyglot persistence'],
          cons: ['Eventual consistency', 'Complexity', 'Sync logic', 'Dual write risk'],
          recommendedFor: ['High read/write ratio', 'Complex queries', 'Different access patterns'],
          avoidWhen: ['Simple domains', 'Strong consistency needed', 'Small team']
        },
        {
          name: 'Event Streaming (Kafka)',
          description: 'Log-centric, replayable, partitioned, ordered per partition',
          pros: ['Massive scale', 'Replay', 'Exactly-once', 'Ecosystem'],
          cons: ['Operational complexity', 'ZooKeeper/KRaft', 'Latency', 'Schema registry'],
          recommendedFor: ['High throughput', 'Event replay', 'Multi-consumer', 'Stream processing'],
          avoidWhen: ['Low volume', 'Simple RPC', 'Team Kafka experience']
        }
      ],
      liveDemo: {
        connectionStatus: 'connected',
        latencyMs: 5,
        messagesPerSecond: 500000,
        indicators: [
          { id: '1', label: 'Produce Latency', value: '5ms', status: 'ok', pulse: true },
          { id: '2', label: 'Events/sec', value: '500K', status: 'ok', pulse: true },
          { id: '3', label: 'Consumer Lag', value: '12', status: 'ok' },
          { id: '4', label: 'Partitions', value: '48', status: 'ok' }
        ]
      },
      position: { x: 1400, y: 200 }
    },
    {
      id: createNodeId('reactive-streams'),
      label: 'Reactive Streams',
      type: 'paradigm',
      description: 'Backpressure-aware async streams. Publisher/Subscriber/Processor contracts. Non-blocking, resilient.',
      details: {
        what: 'Standard for async stream processing with non-blocking backpressure. 4 interfaces: Publisher, Subscriber, Subscription, Processor.',
        how: 'Subscribe → onSubscribe(Subscription) → request(n) → onNext×n → onComplete/onError. Demand-driven flow.',
        when: 'High-throughput pipelines, backpressure critical, reactive databases, streaming ETL, resilience patterns.',
        where: 'Project Reactor, RxJava, Akka Streams, Mutiny, RSocket, Spring WebFlux, Vert.x.',
        why: 'Bounded memory, no thread blocking, composable operators, standard interop, cancellation.',
        productStage: 'production',
        userRPSThresholds: 'Millions of elements/sec with bounded memory',
        liveCloudCosts: 'AWS/GCP $0.006/h per 10K stream elements/sec'
      },
      flowDiagram: {
        title: 'Reactive Streams Flow',
        phases: [
          { id: '1', label: 'Source Emits', description: 'Flux/Mono/Observable creates stream', status: 'completed' },
          { id: '2', label: 'Subscribe', description: 'Downstream requests demand (request(n))', status: 'completed' },
          { id: '3', label: 'Backpressure', description: 'Upstream honors demand, buffers or drops', status: 'active' },
          { id: '4', label: 'Operators', description: 'map, filter, flatMap, buffer, window', status: 'active' },
          { id: '5', label: 'Side Effects', description: 'doOnNext, retry, timeout, circuit breaker', status: 'pending' },
          { id: '6', label: 'Terminal', description: 'onComplete / onError → cleanup', status: 'pending' }
        ]
      },
      variants: [
        {
          name: 'Project Reactor (Java/Kotlin)',
          description: 'Spring WebFlux native, rich operators, reactor-core',
          pros: ['Spring integration', 'Rich API', 'Mono/Flux', 'Context propagation'],
          cons: ['JVM only', 'Learning curve', 'Debugging async'],
          recommendedFor: ['Spring Boot', 'JVM microservices', 'Reactive DB drivers'],
          avoidWhen: ['Non-JVM', 'Simple async', 'Team unfamiliar']
        },
        {
          name: 'RxJS (TypeScript/JS)',
          description: 'Observables for frontend/backend, vast operator library',
          pros: ['Universal JS', 'Angular native', 'Huge ecosystem', 'Marble testing'],
          cons: ['Bundle size', 'Memory leaks', 'Scheduler complexity'],
          recommendedFor: ['Angular', 'Node.js streams', 'Complex UI state', 'Event composition'],
          avoidWhen: ['Simple promises', 'Bundle critical', 'No async complexity']
        },
        {
          name: 'RSocket',
          description: 'Application protocol with reactive streams semantics built-in',
          pros: ['Bi-directional', 'Resume', 'Flow control', 'Transport agnostic'],
          cons: ['Newer', 'Less tooling', 'Client support'],
          recommendedFor: ['Service mesh', 'Gateway', 'Mobile-server', 'Streaming RPC'],
          avoidWhen: ['Standard HTTP', 'REST needed', 'Simple CRUD']
        }
      ],
      liveDemo: {
        connectionStatus: 'connected',
        latencyMs: 3,
        messagesPerSecond: 1000000,
        indicators: [
          { id: '1', label: 'Throughput', value: '1M elem/s', status: 'ok', pulse: true },
          { id: '2', label: 'Backpressure', value: 'ACTIVE', status: 'ok', pulse: true },
          { id: '3', label: 'Buffer', value: '2.3%', status: 'ok' },
          { id: '4', label: 'Cancellations', value: '0', status: 'ok' }
        ]
      },
      position: { x: 1800, y: 400 }
    },
    // Overflow nodes (partial, no labels)
    {
      id: createNodeId('overflow-left-1'),
      label: '',
      type: 'concept',
      description: '',
      details: { what: '', how: '', when: '', where: '', why: '', productStage: 'production', userRPSThresholds: '', liveCloudCosts: '' },
      position: { x: -100, y: 400 },
      isPartial: 'overflow-left'
    },
    {
      id: createNodeId('overflow-left-2'),
      label: '',
      type: 'concept',
      description: '',
      details: { what: '', how: '', when: '', where: '', why: '', productStage: 'production', userRPSThresholds: '', liveCloudCosts: '' },
      position: { x: -100, y: 300 },
      isPartial: 'overflow-left'
    },
    {
      id: createNodeId('overflow-right-1'),
      label: '',
      type: 'concept',
      description: '',
      details: { what: '', how: '', when: '', where: '', why: '', productStage: 'production', userRPSThresholds: '', liveCloudCosts: '' },
      position: { x: 2200, y: 400 },
      isPartial: 'overflow-right'
    },
    {
      id: createNodeId('overflow-right-2'),
      label: '',
      type: 'concept',
      description: '',
      details: { what: '', how: '', when: '', where: '', why: '', productStage: 'production', userRPSThresholds: '', liveCloudCosts: '' },
      position: { x: 2200, y: 500 },
      isPartial: 'overflow-right'
    }
  ],
  edges: [
    { id: createEdgeId('e1'), source: createNodeId('request-response'), target: createNodeId('websockets'), type: 'evolution', opacity: 0.5 },
    { id: createEdgeId('e2'), source: createNodeId('websockets'), target: createNodeId('realtime-sync-engine'), type: 'connection', opacity: 0.5 },
    { id: createEdgeId('e3'), source: createNodeId('request-response'), target: createNodeId('sse'), type: 'variant', opacity: 0.25 },
    { id: createEdgeId('e4'), source: createNodeId('sse'), target: createNodeId('websockets'), type: 'variant', opacity: 0.25 },
    { id: createEdgeId('e5'), source: createNodeId('websockets'), target: createNodeId('graphql-subscriptions'), type: 'connection', opacity: 0.25 },
    { id: createEdgeId('e6'), source: createNodeId('websockets'), target: createNodeId('grpc-streaming'), type: 'connection', opacity: 0.25 },
    { id: createEdgeId('e7'), source: createNodeId('websockets'), target: createNodeId('mqtt'), type: 'connection', opacity: 0.25 },
    { id: createEdgeId('e8'), source: createNodeId('websockets'), target: createNodeId('webrtc-datachannel'), type: 'connection', opacity: 0.25 },
    { id: createEdgeId('e9'), source: createNodeId('realtime-sync-engine'), target: createNodeId('webrtc-datachannel'), type: 'connection', opacity: 0.25 },
    { id: createEdgeId('e10'), source: createNodeId('event-driven'), target: createNodeId('reactive-streams'), type: 'connection', opacity: 0.25 },
    { id: createEdgeId('e11'), source: createNodeId('grpc-streaming'), target: createNodeId('reactive-streams'), type: 'connection', opacity: 0.25 },
    { id: createEdgeId('e12'), source: createNodeId('mqtt'), target: createNodeId('event-driven'), type: 'connection', opacity: 0.25 },
  ]
}

export const paradigmsGraphData: GraphData = {
  nodes: [
    {
      id: createNodeId('event-driven'),
      label: 'Event-Driven',
      type: 'paradigm',
      description: 'Decoupled services communicate via immutable events. Event sourcing, CQRS, saga patterns.',
      details: {
        what: 'Services emit/consume immutable events. Event store as source of truth. Projections build read models.',
        how: 'Command → Validation → Event → Event Store → Projections → Read Models → Subscribers react.',
        when: 'Complex domains, audit requirements, temporal queries, microservices decoupling, compliance.',
        where: 'Kafka, EventStoreDB, NATS JetStream, Redis Streams, AWS EventBridge, Azure Event Hubs.',
        why: 'Complete audit trail, time-travel debugging, replay capability, loose coupling, scalability.',
        productStage: 'production',
        userRPSThresholds: '1M+ events/sec (Kafka), 100K+ (EventStoreDB)',
        liveCloudCosts: 'AWS/GCP $0.012/h per 10K events/sec (Kafka managed)'
      },
      flowDiagram: {
        title: 'Event-Driven Flow',
        phases: [
          { id: '1', label: 'Command Received', description: 'API Gateway → Command Handler', status: 'completed' },
          { id: '2', label: 'Validation', description: 'Business rules, invariants checked', status: 'completed' },
          { id: '3', label: 'Event Emitted', description: 'Immutable fact appended to log', status: 'active' },
          { id: '4', label: 'Projections Update', description: 'Read models rebuilt asynchronously', status: 'pending' },
          { id: '5', label: 'Downstream React', description: 'Consumers process, emit new events', status: 'pending' },
          { id: '6', label: 'Saga Coordination', description: 'Compensating transactions on failure', status: 'pending' }
        ]
      },
      variants: [
        {
          name: 'Event Sourcing',
          description: 'State = fold(events). Full history, time travel, rebuild anytime',
          pros: ['Complete history', 'Time travel', 'Rebuild projections', 'Audit native'],
          cons: ['Event schema migration', 'Snapshot strategy', 'Learning curve', 'Storage growth'],
          recommendedFor: ['Financial', 'Compliance', 'Complex domains', 'Audit critical'],
          avoidWhen: ['Simple CRUD', 'Team inexperience', 'Low audit needs']
        },
        {
          name: 'CQRS',
          description: 'Separate write model (commands) from read model (queries)',
          pros: ['Optimized reads', 'Scalable independently', 'Different models', 'Polyglot persistence'],
          cons: ['Eventual consistency', 'Complexity', 'Sync logic', 'Dual write risk'],
          recommendedFor: ['High read/write ratio', 'Complex queries', 'Different access patterns'],
          avoidWhen: ['Simple domains', 'Strong consistency needed', 'Small team']
        },
        {
          name: 'Event Streaming (Kafka)',
          description: 'Log-centric, replayable, partitioned, ordered per partition',
          pros: ['Massive scale', 'Replay', 'Exactly-once', 'Ecosystem'],
          cons: ['Operational complexity', 'ZooKeeper/KRaft', 'Latency', 'Schema registry'],
          recommendedFor: ['High throughput', 'Event replay', 'Multi-consumer', 'Stream processing'],
          avoidWhen: ['Low volume', 'Simple RPC', 'Team Kafka experience']
        }
      ],
      liveDemo: {
        connectionStatus: 'connected',
        latencyMs: 5,
        messagesPerSecond: 500000,
        indicators: [
          { id: '1', label: 'Produce Latency', value: '5ms', status: 'ok', pulse: true },
          { id: '2', label: 'Events/sec', value: '500K', status: 'ok', pulse: true },
          { id: '3', label: 'Consumer Lag', value: '12', status: 'ok' },
          { id: '4', label: 'Partitions', value: '48', status: 'ok' }
        ]
      },
      position: { x: 400, y: 300 }
    },
    {
      id: createNodeId('reactive-streams'),
      label: 'Reactive Streams',
      type: 'paradigm',
      description: 'Backpressure-aware async streams. Publisher/Subscriber/Processor contracts. Non-blocking, resilient.',
      details: {
        what: 'Standard for async stream processing with non-blocking backpressure. 4 interfaces: Publisher, Subscriber, Subscription, Processor.',
        how: 'Subscribe → onSubscribe(Subscription) → request(n) → onNext×n → onComplete/onError. Demand-driven flow.',
        when: 'High-throughput pipelines, backpressure critical, reactive databases, streaming ETL, resilience patterns.',
        where: 'Project Reactor, RxJava, Akka Streams, Mutiny, RSocket, Spring WebFlux, Vert.x.',
        why: 'Bounded memory, no thread blocking, composable operators, standard interop, cancellation.',
        productStage: 'production',
        userRPSThresholds: 'Millions of elements/sec with bounded memory',
        liveCloudCosts: 'AWS/GCP $0.006/h per 10K stream elements/sec'
      },
      flowDiagram: {
        title: 'Reactive Streams Flow',
        phases: [
          { id: '1', label: 'Source Emits', description: 'Flux/Mono/Observable creates stream', status: 'completed' },
          { id: '2', label: 'Subscribe', description: 'Downstream requests demand (request(n))', status: 'completed' },
          { id: '3', label: 'Backpressure', description: 'Upstream honors demand, buffers or drops', status: 'active' },
          { id: '4', label: 'Operators', description: 'map, filter, flatMap, buffer, window', status: 'active' },
          { id: '5', label: 'Side Effects', description: 'doOnNext, retry, timeout, circuit breaker', status: 'pending' },
          { id: '6', label: 'Terminal', description: 'onComplete / onError → cleanup', status: 'pending' }
        ]
      },
      variants: [
        {
          name: 'Project Reactor (Java/Kotlin)',
          description: 'Spring WebFlux native, rich operators, reactor-core',
          pros: ['Spring integration', 'Rich API', 'Mono/Flux', 'Context propagation'],
          cons: ['JVM only', 'Learning curve', 'Debugging async'],
          recommendedFor: ['Spring Boot', 'JVM microservices', 'Reactive DB drivers'],
          avoidWhen: ['Non-JVM', 'Simple async', 'Team unfamiliar']
        },
        {
          name: 'RxJS (TypeScript/JS)',
          description: 'Observables for frontend/backend, vast operator library',
          pros: ['Universal JS', 'Angular native', 'Huge ecosystem', 'Marble testing'],
          cons: ['Bundle size', 'Memory leaks', 'Scheduler complexity'],
          recommendedFor: ['Angular', 'Node.js streams', 'Complex UI state', 'Event composition'],
          avoidWhen: ['Simple promises', 'Bundle critical', 'No async complexity']
        },
        {
          name: 'RSocket',
          description: 'Application protocol with reactive streams semantics built-in',
          pros: ['Bi-directional', 'Resume', 'Flow control', 'Transport agnostic'],
          cons: ['Newer', 'Less tooling', 'Client support'],
          recommendedFor: ['Service mesh', 'Gateway', 'Mobile-server', 'Streaming RPC'],
          avoidWhen: ['Standard HTTP', 'REST needed', 'Simple CRUD']
        }
      ],
      liveDemo: {
        connectionStatus: 'connected',
        latencyMs: 3,
        messagesPerSecond: 1000000,
        indicators: [
          { id: '1', label: 'Throughput', value: '1M elem/s', status: 'ok', pulse: true },
          { id: '2', label: 'Backpressure', value: 'ACTIVE', status: 'ok', pulse: true },
          { id: '3', label: 'Buffer', value: '2.3%', status: 'ok' },
          { id: '4', label: 'Cancellations', value: '0', status: 'ok' }
        ]
      },
      position: { x: 800, y: 300 }
    },
    {
      id: createNodeId('local-first-crdt'),
      label: 'Local-First / CRDT',
      type: 'paradigm',
      description: 'Conflict-free replicated data types for offline-first, peer-to-peer synchronization.',
      details: {
        what: 'Data structures that automatically merge concurrent edits without coordination. Math-based convergence guarantees.',
        how: 'Local mutation → CRDT operation → Broadcast → Peer merge (commutative/associative/idempotent) → Convergence.',
        when: 'Collaborative apps, offline-first, P2P, multi-device sync, conflict-free merging.',
        where: 'Yjs, Automerge, Riak, Redis CRDT, custom implementations, WebRTC mesh.',
        why: 'Zero coordination, offline-capable, eventual consistency, user intent preservation, no central server.',
        productStage: 'beta',
        userRPSThresholds: '10K–100K concurrent sessions (depends on doc size)',
        liveCloudCosts: 'AWS/GCP $0.015/h per 1K active sessions (relay + storage)'
      },
      flowDiagram: {
        title: 'Local-First CRDT Flow',
        phases: [
          { id: '1', label: 'Local Edit', description: 'User types → CRDT op generated instantly', status: 'completed' },
          { id: '2', label: 'Optimistic UI', description: 'Immediate reflect, zero latency feel', status: 'completed' },
          { id: '3', label: 'Broadcast', description: 'WebSocket/WebRTC → All peers + relay', status: 'active' },
          { id: '4', label: 'Auto-Merge', description: 'Receive ops → Mathematical merge → Converge', status: 'pending' },
          { id: '5', label: 'Awareness', description: 'Cursors, selections, presence shared', status: 'pending' },
          { id: '6', label: 'Persistence', description: 'IndexedDB + Server backup → Durability', status: 'pending' }
        ]
      },
      variants: [
        {
          name: 'Yjs',
          description: 'Mature CRDT, rich text, many providers, WebRTC',
          pros: ['Battle-tested', 'Rich text', 'Many providers', 'Offline-first', 'Awareness'],
          cons: ['Bundle size', 'Memory growth', 'Learning curve'],
          recommendedFor: ['Rich text', 'Complex structures', 'P2P needed', 'Production apps'],
          avoidWhen: ['Simple state', 'Bundle critical', 'Custom merge']
        },
        {
          name: 'Automerge',
          description: 'JSON-like CRDT, immutable, functional, WASM accelerated',
          pros: ['JSON-like', 'Immutable', 'Time-travel', 'WASM fast', 'Document model'],
          cons: ['Memory heavy', 'No rich text', 'Slower large docs', 'No awareness'],
          recommendedFor: ['JSON docs', 'Audit trails', 'Functional apps', 'Time travel'],
          avoidWhen: ['Rich text', 'High frequency', 'Mobile memory', 'Real-time cursors']
        },
        {
          name: 'Riak DT / Redis CRDT',
          description: 'Server-side CRDTs, distributed database integrated',
          pros: ['Server-managed', 'Distributed', 'Operational simplicity', 'Strong consistency options'],
          cons: ['Vendor lock-in', 'Less flexible', 'No P2P', 'Schema limits'],
          recommendedFor: ['Server-centric', 'Counters/Sets/Maps', 'Operational ease', 'Existing Riak/Redis'],
          avoidWhen: ['P2P needed', 'Rich text', 'Custom types', 'Offline-first client']
        }
      ],
      liveDemo: {
        connectionStatus: 'connected',
        latencyMs: 12,
        messagesPerSecond: 8500,
        indicators: [
          { id: '1', label: 'Sync Latency', value: '12ms', status: 'ok', pulse: true },
          { id: '2', label: 'Active Peers', value: '23', status: 'ok' },
          { id: '3', label: 'Ops/sec', value: '8.5K', status: 'ok', pulse: true },
          { id: '4', label: 'Conflicts', value: '0', status: 'ok' }
        ]
      },
      position: { x: 1200, y: 300 }
    },
    {
      id: createNodeId('consensus-driven'),
      label: 'Consensus-Driven',
      type: 'paradigm',
      description: 'Distributed consensus protocols (Raft, Paxos) for linearizable state machine replication.',
      details: {
        what: 'State machine replicated across quorum. Leader election, log replication, safety + liveness guarantees.',
        how: 'Client → Leader → AppendEntries → Quorum ACK → Apply → Respond. Followers replicate, new leader elected on failure.',
        when: 'Strong consistency required, configuration management, coordination, leader election, distributed locks.',
        where: 'etcd, Consul, ZooKeeper, TiKV, CockroachDB, FoundationDB, custom Raft.',
        why: 'Linearizable reads/writes, automatic failover, proven algorithms, battle-tested implementations.',
        productStage: 'production',
        userRPSThresholds: '10K–100K writes/sec (depends on cluster size, network)',
        liveCloudCosts: 'AWS/GCP $0.02/h per 3-node cluster (managed etcd/Consul)'
      },
      flowDiagram: {
        title: 'Consensus (Raft) Flow',
        phases: [
          { id: '1', label: 'Leader Election', description: 'RequestVote, term increment, quorum', status: 'completed' },
          { id: '2', label: 'Client Request', description: 'Redirect to leader, append to log', status: 'completed' },
          { id: '3', label: 'Log Replication', description: 'AppendEntries to followers, parallel', status: 'active' },
          { id: '4', label: 'Quorum ACK', description: 'Majority acknowledges, commit index', status: 'pending' },
          { id: '5', label: 'Apply & Respond', description: 'State machine apply, client response', status: 'pending' },
          { id: '6', label: 'Follower Apply', description: 'Async apply, read index/lease read', status: 'pending' }
        ]
      },
      variants: [
        {
          name: 'Raft',
          description: 'Understandable consensus, leader-based, strong leader',
          pros: ['Understandable', 'Leader lease reads', 'Joint consensus', 'Mature impls'],
          cons: ['Leader bottleneck', 'Single leader', 'Cross-region latency'],
          recommendedFor: ['Config management', 'Service discovery', 'Metadata', 'Single DC'],
          avoidWhen: ['Multi-region writes', 'Highest throughput', 'Leaderless needed']
        },
        {
          name: 'Multi-Raft / Multi-Paxos',
          description: 'Sharded consensus groups, parallel leaders',
          pros: ['Horizontal scale', 'Parallel leaders', 'Tenant isolation', 'Hot spot split'],
          cons: ['Complex', 'Cross-shard tx', 'Operational burden', 'Rebalancing'],
          recommendedFor: ['Distributed SQL', 'Multi-tenant', 'High scale', 'Geo-distributed'],
          avoidWhen: ['Simple needs', 'Small cluster', 'Team inexperience']
        },
        {
          name: 'EPaxos / Flexible Paxos',
          description: 'Leaderless, commutative commands, fast path',
          pros: ['No leader', 'Fast path', 'Geo-friendly', 'Commutative fast'],
          cons: ['Complex', 'Less mature', 'Dependency tracking', 'Fewer impls'],
          recommendedFor: ['Geo-distributed', 'Low latency', 'Commutative ops', 'Research'],
          avoidWhen: ['Standard Raft OK', 'Team unfamiliar', 'Production critical']
        }
      ],
      liveDemo: {
        connectionStatus: 'connected',
        latencyMs: 8,
        messagesPerSecond: 25000,
        indicators: [
          { id: '1', label: 'Commit Latency', value: '8ms', status: 'ok', pulse: true },
          { id: '2', label: 'Leader', value: 'node-3', status: 'ok' },
          { id: '3', label: 'Writes/sec', value: '25K', status: 'ok', pulse: true },
          { id: '4', label: 'Term', value: '47', status: 'ok' }
        ]
      },
      position: { x: 400, y: 500 }
    },
    {
      id: createNodeId('planet-scale'),
      label: 'Planet-Scale Composite',
      type: 'paradigm',
      description: 'Composable building blocks: sharding, replication, consensus, edge caching, global routing.',
      details: {
        what: 'Architecture pattern composing multiple paradigms: sharded consensus, async replication, edge compute, global load balancing.',
        how: 'Global router → Regional shards (Raft) → Async cross-region replication → Edge cache → Client affinity.',
        when: 'Global user base, <100ms P99 worldwide, multi-region compliance, massive scale, high availability.',
        where: 'Spanner, CockroachDB, DynamoDB Global Tables, Cosmos DB, Fauna, custom composites.',
        why: 'Horizontal scale, geo-distribution, survival semantics, composable primitives, operational maturity.',
        productStage: 'production',
        userRPSThresholds: 'Millions of QPS globally, <50ms P50, <200ms P99',
        liveCloudCosts: 'AWS/GCP $0.15/h per 10K global QPS (managed planetary DB)'
      },
      flowDiagram: {
        title: 'Planet-Scale Composite Flow',
        phases: [
          { id: '1', label: 'Global Route', description: 'Anycast/DNS → Nearest region', status: 'completed' },
          { id: '2', label: 'Regional Shard', description: 'Raft group per shard, local quorum', status: 'active' },
          { id: '3', label: 'Local Commit', description: 'Sync replicate in region, acknowledge', status: 'pending' },
          { id: '4', label: 'Async Replication', description: 'Cross-region log shipping, eventual', status: 'pending' },
          { id: '5', label: 'Edge Cache', description: 'CDN/Workers cache reads, invalidate', status: 'pending' },
          { id: '6', label: 'Failover', description: 'Region loss → Promote replica, re-route', status: 'pending' }
        ]
      },
      variants: [
        {
          name: 'Spanner / CockroachDB',
          description: 'SQL, strong consistency, global transactions, TrueTime/HLC',
          pros: ['Strong consistency', 'Distributed SQL', 'Global transactions', 'Mature'],
          cons: ['Cost', 'Latency (Paxos)', 'Schema limits', 'Vendor/ops'],
          recommendedFor: ['Global SQL', 'Financial', 'Strong consistency', 'Multi-region ACID'],
          avoidWhen: ['Cost sensitive', 'Simple key-value', 'Team NoSQL']
        },
        {
          name: 'DynamoDB Global Tables / Cosmos DB',
          description: 'Multi-master, last-writer-wins, automatic replication',
          pros: ['Multi-master', 'Serverless', 'Auto-scale', 'Multi-model (Cosmos)'],
          cons: ['Eventual consistency', 'Conflict resolution', 'Cost at scale', 'Limited query'],
          recommendedFor: ['Multi-master', 'Serverless', 'Key-value', 'Global apps'],
          avoidWhen: ['Strong consistency', 'Complex queries', 'Relational needs']
        },
        {
          name: 'Custom Composite',
          description: 'Build from primitives: Raft shards + async replication + edge',
          pros: ['Full control', 'Cost optimized', 'Tailored', 'No vendor lock'],
          cons: ['Operational burden', 'Expertise needed', 'Time to market', 'Testing complexity'],
          recommendedFor: ['Unique requirements', 'Cost critical', 'Expert team', 'Long-term'],
          avoidWhen: ['Standard needs', 'Small team', 'Fast delivery', 'Managed preferred']
        }
      ],
      liveDemo: {
        connectionStatus: 'connected',
        latencyMs: 42,
        messagesPerSecond: 100000,
        indicators: [
          { id: '1', label: 'Global P50', value: '42ms', status: 'ok', pulse: true },
          { id: '2', label: 'Global P99', value: '187ms', status: 'ok' },
          { id: '3', label: 'Regions', value: '6 active', status: 'ok' },
          { id: '4', label: 'QPS', value: '1.2M', status: 'ok', pulse: true }
        ]
      },
      position: { x: 800, y: 500 }
    },
    // Overflow nodes for paradigms
    {
      id: createNodeId('overflow-left-1'),
      label: '',
      type: 'concept',
      description: '',
      details: { what: '', how: '', when: '', where: '', why: '', productStage: 'production', userRPSThresholds: '', liveCloudCosts: '' },
      position: { x: -100, y: 300 },
      isPartial: 'overflow-left'
    },
    {
      id: createNodeId('overflow-left-2'),
      label: '',
      type: 'concept',
      description: '',
      details: { what: '', how: '', when: '', where: '', why: '', productStage: 'production', userRPSThresholds: '', liveCloudCosts: '' },
      position: { x: -100, y: 500 },
      isPartial: 'overflow-left'
    },
    {
      id: createNodeId('overflow-right-1'),
      label: '',
      type: 'concept',
      description: '',
      details: { what: '', how: '', when: '', where: '', why: '', productStage: 'production', userRPSThresholds: '', liveCloudCosts: '' },
      position: { x: 1600, y: 300 },
      isPartial: 'overflow-right'
    },
    {
      id: createNodeId('overflow-right-2'),
      label: '',
      type: 'concept',
      description: '',
      details: { what: '', how: '', when: '', where: '', why: '', productStage: 'production', userRPSThresholds: '', liveCloudCosts: '' },
      position: { x: 1600, y: 500 },
      isPartial: 'overflow-right'
    }
  ],
  edges: [
    { id: createEdgeId('e1'), source: createNodeId('event-driven'), target: createNodeId('reactive-streams'), type: 'connection', opacity: 0.5 },
    { id: createEdgeId('e2'), source: createNodeId('reactive-streams'), target: createNodeId('local-first-crdt'), type: 'connection', opacity: 0.5 },
    { id: createEdgeId('e3'), source: createNodeId('event-driven'), target: createNodeId('consensus-driven'), type: 'connection', opacity: 0.25 },
    { id: createEdgeId('e4'), source: createNodeId('consensus-driven'), target: createNodeId('planet-scale'), type: 'connection', opacity: 0.25 },
    { id: createEdgeId('e5'), source: createNodeId('reactive-streams'), target: createNodeId('planet-scale'), type: 'connection', opacity: 0.25 },
    { id: createEdgeId('e6'), source: createNodeId('local-first-crdt'), target: createNodeId('planet-scale'), type: 'connection', opacity: 0.25 },
  ]
}