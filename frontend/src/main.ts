import './styles.css'

type Envelope<T> = {
  data: T
  source: string
  confidence: string
  staleness: string
  measuredAt: string
}

type Stream = {
  id: string
  title: string
  role: string
  category: string
  description: string
  status: string
  dependencies: string[]
  deliverables: string[]
  acceptanceGates: string[]
  executionPackets: string[]
}

type BootstrapData = {
  summary: string
  phase: string
  repo: string
  runtime: {
    api: string
    frontend: string
    mandatoryInfra: string[]
    optionalInfra: string[]
    composeFiles: string[]
  }
  now: {
    activeStreams: string[]
    executionOrder: string[]
    goNoGoForImplementation: string
  }
  notYet: string[]
}

type CatalogData = {
  families: Array<{
    name: string
    category: string
    status: string
    recommendedWhen: string[]
    avoidWhen: string[]
    notes: string[]
  }>
  notes: string[]
}

type FirstWaveContract = {
  executionOrder: string[]
  streams: Stream[]
  rules: string[]
}

type SecurityBootstrap = {
  identity: {
    provider: string
    browserPattern: string
    accessTokenTarget: string
    refreshPolicy: string
    redirectPolicy: string
  }
  vault: {
    provider: string
    ownership: string[]
    pki: string[]
  }
  mtlsPhases: Array<{
    phase: number
    name: string
    goal: string
  }>
  currentState: string
}

type DataPlatform = {
  canonicalStores: Array<{
    name: string
    role: string
    status: string
    why: string
  }>
  optionalProfiles: Array<{
    name: string
    profile: string
    status: string
    why: string
  }>
  contracts: {
    entityOwnership: string[]
    schemaRules: string[]
    evidenceFields: string[]
  }
}

type BenchmarkFramework = {
  benchmarkMatrix: Array<{
    name: string
    measures: string[]
  }>
  telemetry: {
    logs: string[]
    metrics: string[]
    traces: string[]
  }
  closureRules: string[]
}

type DeferredWave = {
  name: string
  purpose: string
  streams: string[]
  focus: string[]
  reactivationRule: string[]
}

type DeferredWavesData = {
  waves: DeferredWave[]
  holdManifest: string
}

type V2Readiness = {
  status: string
  blockedBy: string[]
  unlockPath: string[]
}

type RealtimeComparisons = {
  comparisons: Array<{
    name: string
    recommendedFor: string[]
    avoidWhen: string[]
    latencyClass: string
    defaultDecision: string
    notes: string[]
  }>
  delayAppliedMs: number
  timeoutThresholdMs: number
}

type TransportSummary = {
  availableDemos: Array<{
    name: string
    endpoint: string
    status: string
  }>
  notes: string[]
}

type MessagingOverview = {
  families: Array<{
    name: string
    status: string
    useFor: string[]
  }>
  notes: string[]
}

type SyncOverview = {
  families: Array<{
    name: string
    status: string
    useFor: string[]
  }>
  notes: string[]
}

type ProjectionOverview = {
  families: Array<{
    name: string
    status: string
    useFor: string[]
  }>
  notes: string[]
}

type PricingOverview = {
  currency: string
  selectedFreshness: number
  prices: Array<{
    provider: string
    service: string
    region: string
    currency: string
    hourlyUsd: number
    freshnessS: number
    confidence: string
    staleness: string
    fallback: boolean
  }>
  notes: string[]
}

type V2Roadmap = {
  status: string
  paradigms: Array<{
    name: string
    status: string
    unlock: string
  }>
  v1Completed: string[]
  nextBeforeV2Implementation: string[]
}

type EventDrivenParadigm = {
  paradigm: string
  status: string
  commands: string[]
  events: Array<{
    id: string
    kind: string
    payload: Record<string, unknown>
    emittedAt: string
    source: string
  }>
  projections: {
    eventKinds: Record<string, number>
    sourceOfTruth: string
    transientBroadcast: string
  }
  notes: string[]
}

type ReactiveStreamDriven = {
  paradigm: string
  status: string
  streamSignals: {
    recentEventCount: number
    eventKinds: Record<string, number>
    latestEventAt: string
    transportSurface: string
  }
  backpressureHints: string[]
  notes: string[]
}

type LocalFirstCRDT = {
  paradigm: string
  status: string
  syncSignals: {
    recentSessionCount: number
    alignedSessions: number
    conflictSessions: number
    maxObservedLag: number
  }
  mergeRules: string[]
  notes: string[]
}

type ConsensusDriven = {
  paradigm: string
  status: string
  coordinationSignals: {
    candidateSessions: number
    quorumReady: number
    quorumBlocked: number
    rule: string
  }
  decisionRules: string[]
  notes: string[]
}

type PlanetScaleComposite = {
  paradigm: string
  status: string
  buildingBlocks: Array<{
    name: string
    role: string
  }>
  compositionRules: string[]
  notes: string[]
}

const uiMessages = {
  en: {
    title: 'First-wave execution console',
    lede: 'This bootstrap turns the backend-first planning package into a runnable foundation: explicit first-wave contracts, deferred-wave holds, and visible timing behavior.',
    localeLabel: 'Locale',
  },
  es: {
    title: 'Consola de ejecucion de primera ola',
    lede: 'Este bootstrap convierte el paquete de planificacion backend-first en una base ejecutable: contratos explicitos de primera ola, olas posteriores diferidas y comportamiento temporal visible.',
    localeLabel: 'Idioma',
  },
} as const

type Locale = keyof typeof uiMessages

type JobState = {
  id: string
  status: string
  progress: number
  currentStep: string
  delayAppliedMs: number
  timeoutThresholdMs: number
  timeline: string[]
}

const apiBase = (import.meta.env.VITE_API_BASE_URL as string | undefined) || '/api'
const eventBase = apiBase.replace(/\/api$/, '')
let currentSyncStatusUrl: string | null = null
let currentLocale: Locale = 'en'

const root = document.querySelector<HTMLDivElement>('#app')

if (!root) {
  throw new Error('Missing app root')
}

root.innerHTML = `
  <main class="shell">
    <header class="hero panel">
      <div class="hero-copy">
        <p class="eyebrow">CommLayers</p>
        <h1 id="hero-title">First-wave execution console</h1>
        <p class="lede" id="hero-lede">This bootstrap turns the backend-first planning package into a runnable foundation: explicit first-wave contracts, deferred-wave holds, and visible timing behavior.</p>
      </div>
      <div class="hero-meta" id="hero-meta">Loading...</div>
    </header>

    <section class="panel locale-panel">
      <div class="body locale-body">
        <label for="locale-select" id="locale-label">Locale</label>
        <select id="locale-select">
          <option value="en">English</option>
          <option value="es">Espanol</option>
        </select>
      </div>
    </section>

    <section class="grid overview-grid">
      <article class="panel" id="bootstrap-card"><h2>Bootstrap runtime</h2><div class="body">Loading...</div></article>
      <article class="panel" id="first-wave-card"><h2>First-wave order</h2><div class="body">Loading...</div></article>
    </section>

    <section class="panel" id="stream-contracts-card">
      <h2>First-wave stream contracts</h2>
      <div class="body stream-grid">Loading...</div>
    </section>

    <section class="grid deep-grid">
      <article class="panel" id="security-card"><h2>Security bootstrap</h2><div class="body">Loading...</div></article>
      <article class="panel" id="data-card"><h2>Canonical data platform</h2><div class="body">Loading...</div></article>
      <article class="panel" id="benchmark-card"><h2>Benchmark closure</h2><div class="body">Loading...</div></article>
    </section>

    <section class="grid demo-grid">
      <article class="panel" id="catalog-card"><h2>Catalog snapshot</h2><div class="body">Loading...</div></article>
      <article class="panel" id="realtime-card"><h2>Realtime guidance</h2><div class="body">Loading...</div></article>
      <article class="panel" id="transport-card"><h2>Transport demos</h2><div class="body">Loading...</div></article>
      <article class="panel" id="workflow-card"><h2>Messaging and workflow demos</h2><div class="body">Loading...</div></article>
      <article class="panel" id="sync-card"><h2>Sync and replication demos</h2><div class="body">Loading...</div></article>
      <article class="panel" id="projection-card"><h2>Projection demos</h2><div class="body">Loading...</div></article>
      <article class="panel" id="pricing-card"><h2>Pricing demos</h2><div class="body">Loading...</div></article>
      <article class="panel" id="v2-roadmap-card"><h2>V2 groundwork</h2><div class="body">Loading...</div></article>
      <article class="panel" id="v2-event-card"><h2>V2 event-driven slice</h2><div class="body">Loading...</div></article>
      <article class="panel" id="v2-reactive-card"><h2>V2 reactive stream slice</h2><div class="body">Loading...</div></article>
      <article class="panel" id="v2-localfirst-card"><h2>V2 local-first / CRDT slice</h2><div class="body">Loading...</div></article>
      <article class="panel" id="v2-consensus-card"><h2>V2 consensus-driven slice</h2><div class="body">Loading...</div></article>
      <article class="panel" id="v2-planet-card"><h2>V2 planet-scale composite slice</h2><div class="body">Loading...</div></article>
      <article class="panel" id="events-card"><h2>SSE progress demo</h2><div class="body"><ul id="events-list" class="stack compact"></ul></div></article>
      <article class="panel" id="async-card">
        <h2>Async visibility demo</h2>
        <div class="body stack compact">
          <p>Run a simulated long-running operation with visible progress and delay metadata.</p>
          <button id="run-async" class="button">Run async demo</button>
          <pre id="async-output">Idle</pre>
        </div>
      </article>
    </section>

    <section class="grid deferred-grid">
      <article class="panel" id="deferred-card"><h2>Deferred waves</h2><div class="body">Loading...</div></article>
      <article class="panel" id="v2-card"><h2>V2 readiness gate</h2><div class="body">Loading...</div></article>
    </section>
  </main>
`

function list(items: string[], className = 'stack compact'): string {
  return `<ul class="${className}">${items.map((item) => `<li>${item}</li>`).join('')}</ul>`
}

function apiPath(path: string): string {
  const separator = path.includes('?') ? '&' : '?'
  return `${path}${separator}lang=${currentLocale}`
}

function badge(value: string): string {
  return `<span class="badge badge-${value.replace(/[^a-z0-9]+/gi, '-').toLowerCase()}">${value}</span>`
}

function streamCard(stream: Stream): string {
  return `
    <article class="stream-card stream-${stream.status}">
      <header>
        <div>
          <p class="micro">${stream.id}</p>
          <h3>${stream.title}</h3>
        </div>
        ${badge(stream.status)}
      </header>
      <p class="muted">${stream.description}</p>
      <p><strong>Role:</strong> ${stream.role}</p>
      <p><strong>Dependencies:</strong> ${stream.dependencies.length ? stream.dependencies.join(', ') : 'none'}</p>
      <div class="mini-section">
        <h4>Deliverables</h4>
        ${list(stream.deliverables)}
      </div>
      <div class="mini-section">
        <h4>Acceptance gates</h4>
        ${list(stream.acceptanceGates)}
      </div>
      <div class="mini-section">
        <h4>Execution packets</h4>
        ${list(stream.executionPackets)}
      </div>
    </article>
  `
}

async function fetchEnvelope<T>(path: string): Promise<Envelope<T>> {
  const response = await fetch(path)
  if (!response.ok) {
    throw new Error(`Request failed: ${response.status}`)
  }
  return response.json() as Promise<Envelope<T>>
}

async function renderApp() {
  const copy = uiMessages[currentLocale]
  ;(document.querySelector('#hero-title') as HTMLHeadingElement).textContent = copy.title
  ;(document.querySelector('#hero-lede') as HTMLParagraphElement).textContent = copy.lede
  ;(document.querySelector('#locale-label') as HTMLLabelElement).textContent = copy.localeLabel
  const [bootstrap, firstWave, security, dataPlatform, benchmark, catalog, realtime, transports, messaging, sync, projections, pricing, deferred, v2, v2Roadmap, v2Event, v2Reactive, v2LocalFirst, v2Consensus, v2Planet] = await Promise.all([
    fetchEnvelope<BootstrapData>(apiPath(`${apiBase}/bootstrap`)),
    fetchEnvelope<FirstWaveContract>(apiPath(`${apiBase}/first-wave/contract`)),
    fetchEnvelope<SecurityBootstrap>(apiPath(`${apiBase}/security/bootstrap`)),
    fetchEnvelope<DataPlatform>(apiPath(`${apiBase}/data-platform`)),
    fetchEnvelope<BenchmarkFramework>(apiPath(`${apiBase}/benchmark-framework`)),
    fetchEnvelope<CatalogData>(apiPath(`${apiBase}/catalog`)),
    fetchEnvelope<RealtimeComparisons>(apiPath(`${apiBase}/comparisons/realtime`)),
    fetchEnvelope<TransportSummary>(apiPath(`${apiBase}/transports`)),
    fetchEnvelope<MessagingOverview>(apiPath(`${apiBase}/messaging`)),
    fetchEnvelope<SyncOverview>(apiPath(`${apiBase}/sync`)),
    fetchEnvelope<ProjectionOverview>(apiPath(`${apiBase}/projections`)),
    fetchEnvelope<PricingOverview>(apiPath(`${apiBase}/pricing`)),
    fetchEnvelope<DeferredWavesData>(apiPath(`${apiBase}/deferred-waves`)),
    fetchEnvelope<V2Readiness>(apiPath(`${apiBase}/v2-readiness`)),
    fetchEnvelope<V2Roadmap>(apiPath(`${apiBase}/v2/roadmap`)),
    fetchEnvelope<EventDrivenParadigm>(apiPath(`${apiBase}/v2/paradigms/event-driven`)),
    fetchEnvelope<ReactiveStreamDriven>(apiPath(`${apiBase}/v2/paradigms/reactive-stream-driven`)),
    fetchEnvelope<LocalFirstCRDT>(apiPath(`${apiBase}/v2/paradigms/local-first-crdt`)),
    fetchEnvelope<ConsensusDriven>(apiPath(`${apiBase}/v2/paradigms/consensus-driven`)),
    fetchEnvelope<PlanetScaleComposite>(apiPath(`${apiBase}/v2/paradigms/planet-scale-composite`)),
  ])

  const heroMeta = document.querySelector('#hero-meta') as HTMLDivElement
  heroMeta.innerHTML = `
    <div class="hero-stat"><span>Phase</span><strong>${bootstrap.data.phase}</strong></div>
    <div class="hero-stat"><span>Repo</span><a href="${bootstrap.data.repo}" target="_blank" rel="noreferrer">GitHub</a></div>
    <div class="hero-stat"><span>Execution</span><strong>${bootstrap.data.now.executionOrder.join(' -> ')}</strong></div>
  `

  const bootstrapCard = document.querySelector('#bootstrap-card .body') as HTMLDivElement
  bootstrapCard.innerHTML = `
    <p>${bootstrap.data.summary}</p>
    <dl class="meta-grid">
      <div><dt>API</dt><dd>${bootstrap.data.runtime.api}</dd></div>
      <div><dt>Frontend</dt><dd>${bootstrap.data.runtime.frontend}</dd></div>
      <div><dt>Mandatory infra</dt><dd>${bootstrap.data.runtime.mandatoryInfra.join(', ')}</dd></div>
      <div><dt>Optional infra</dt><dd>${bootstrap.data.runtime.optionalInfra.join(', ')}</dd></div>
      <div><dt>Compose files</dt><dd>${bootstrap.data.runtime.composeFiles.join(', ')}</dd></div>
    </dl>
    <div class="mini-section">
      <h4>Not yet implemented</h4>
      ${list(bootstrap.data.notYet)}
    </div>
  `

  const firstWaveCard = document.querySelector('#first-wave-card .body') as HTMLDivElement
  firstWaveCard.innerHTML = `
    <p><strong>Go / no-go:</strong> ${bootstrap.data.now.goNoGoForImplementation}</p>
    <p><strong>Execution order:</strong> ${firstWave.data.executionOrder.join(' -> ')}</p>
    <div class="mini-section">
      <h4>Wave rules</h4>
      ${list(firstWave.data.rules)}
    </div>
  `

  const contractsCard = document.querySelector('#stream-contracts-card .body') as HTMLDivElement
  contractsCard.innerHTML = firstWave.data.streams.map(streamCard).join('')

  const securityCard = document.querySelector('#security-card .body') as HTMLDivElement
  securityCard.innerHTML = `
    <p><strong>Browser pattern:</strong> ${security.data.identity.browserPattern}</p>
    <p><strong>Access tokens:</strong> ${security.data.identity.accessTokenTarget}</p>
    <p><strong>Refresh:</strong> ${security.data.identity.refreshPolicy}</p>
    <p><strong>Redirects:</strong> ${security.data.identity.redirectPolicy}</p>
    <div class="mini-section"><h4>Vault ownership</h4>${list(security.data.vault.ownership)}</div>
    <div class="mini-section"><h4>PKI baseline</h4>${list(security.data.vault.pki)}</div>
    <div class="mini-section">
      <h4>mTLS phases</h4>
      <ul class="stack compact">${security.data.mtlsPhases
        .map((phase) => `<li><strong>Phase ${phase.phase}:</strong> ${phase.name}<br /><span>${phase.goal}</span></li>`)
        .join('')}</ul>
    </div>
    <p class="muted">${security.data.currentState}</p>
  `

  const dataCard = document.querySelector('#data-card .body') as HTMLDivElement
  dataCard.innerHTML = `
    <div class="mini-section">
      <h4>Canonical stores</h4>
      <ul class="stack compact">${dataPlatform.data.canonicalStores
        .map((store) => `<li><strong>${store.name}</strong> ${badge(store.status)}<br />${store.role} - ${store.why}</li>`)
        .join('')}</ul>
    </div>
    <div class="mini-section">
      <h4>Optional profiles</h4>
      <ul class="stack compact">${dataPlatform.data.optionalProfiles
        .map((store) => `<li><strong>${store.name}</strong> <span class="chip">${store.profile}</span><br />${store.why}</li>`)
        .join('')}</ul>
    </div>
    <div class="mini-section"><h4>Entity ownership</h4>${list(dataPlatform.data.contracts.entityOwnership)}</div>
    <div class="mini-section"><h4>Schema rules</h4>${list(dataPlatform.data.contracts.schemaRules)}</div>
    <div class="mini-section"><h4>Evidence fields</h4>${list(dataPlatform.data.contracts.evidenceFields)}</div>
  `

  const benchmarkCard = document.querySelector('#benchmark-card .body') as HTMLDivElement
  benchmarkCard.innerHTML = `
    <div class="mini-section">
      <h4>Benchmark matrix</h4>
      <ul class="stack compact">${benchmark.data.benchmarkMatrix
        .map((item) => `<li><strong>${item.name}</strong><br />${item.measures.join(', ')}</li>`)
        .join('')}</ul>
    </div>
    <div class="mini-section"><h4>Logs</h4>${list(benchmark.data.telemetry.logs)}</div>
    <div class="mini-section"><h4>Metrics</h4>${list(benchmark.data.telemetry.metrics)}</div>
    <div class="mini-section"><h4>Traces</h4>${list(benchmark.data.telemetry.traces)}</div>
    <div class="mini-section"><h4>Closure rules</h4>${list(benchmark.data.closureRules)}</div>
  `

  const catalogCard = document.querySelector('#catalog-card .body') as HTMLDivElement
  catalogCard.innerHTML = `
    <ul class="stack compact">
      ${catalog.data.families
        .map(
          (item) => `<li><strong>${item.name}</strong> ${badge(item.status)}<br />Use: ${item.recommendedWhen.join(', ')}<br />Avoid: ${item.avoidWhen.join(', ')}<br /><span>${item.notes.join(' ')}</span></li>`,
        )
        .join('')}
    </ul>
  `

  const realtimeCard = document.querySelector('#realtime-card .body') as HTMLDivElement
  realtimeCard.innerHTML = `
    <p><strong>Delay applied:</strong> ${realtime.data.delayAppliedMs} ms</p>
    <p><strong>Timeout threshold:</strong> ${realtime.data.timeoutThresholdMs} ms</p>
    <ul class="stack compact">${realtime.data.comparisons
      .map(
        (item) => `<li><strong>${item.name}</strong> <span class="chip">${item.defaultDecision}</span><br />Use: ${item.recommendedFor.join(', ')}<br />Avoid: ${item.avoidWhen.join(', ')}<br />Latency: ${item.latencyClass}<br /><span>${item.notes.join(' ')}</span></li>`,
      )
      .join('')}</ul>
  `

  const transportCard = document.querySelector('#transport-card .body') as HTMLDivElement
  transportCard.innerHTML = `
    <p>Exercise the newly available transport demos through the frontend gateway.</p>
    <div class="mini-section"><h4>Available demos</h4>${list(
      transports.data.availableDemos.map((item) => `${item.name} -> ${item.endpoint} (${item.status})`),
    )}</div>
    <div class="mini-section"><h4>Notes</h4>${list(transports.data.notes)}</div>
    <div class="button-row">
      <button id="run-polling" class="button">Run polling demo</button>
      <button id="run-long-polling" class="button">Run long-polling demo</button>
      <button id="run-websocket" class="button">Run WebSocket demo</button>
    </div>
    <pre id="transport-output">Idle</pre>
  `

  const workflowCard = document.querySelector('#workflow-card .body') as HTMLDivElement
  workflowCard.innerHTML = `
    <p>Exercise the next sequential backend-family slice without introducing a real broker yet.</p>
    <div class="mini-section"><h4>Families</h4>${list(
      messaging.data.families.map((family) => `${family.name} (${family.status}) -> ${family.useFor.join(', ')}`),
    )}</div>
    <div class="mini-section"><h4>Notes</h4>${list(messaging.data.notes)}</div>
    <div class="button-row">
      <button id="run-workflow" class="button">Run workflow demo</button>
    </div>
    <pre id="workflow-output">Idle</pre>
  `

  const syncCard = document.querySelector('#sync-card .body') as HTMLDivElement
  syncCard.innerHTML = `
    <p>Model primary/replica drift and explicit replication before adding heavier distributed machinery.</p>
    <div class="mini-section"><h4>Families</h4>${list(
      sync.data.families.map((family) => `${family.name} (${family.status}) -> ${family.useFor.join(', ')}`),
    )}</div>
    <div class="mini-section"><h4>Notes</h4>${list(sync.data.notes)}</div>
    <div class="button-row">
      <button id="create-sync" class="button">Create sync session</button>
      <button id="mutate-primary" class="button">Mutate primary</button>
      <button id="mutate-replica" class="button">Mutate replica</button>
      <button id="replicate-sync" class="button">Replicate</button>
    </div>
    <pre id="sync-output">Idle</pre>
  `

  const projectionCard = document.querySelector('#projection-card .body') as HTMLDivElement
  projectionCard.innerHTML = `
    <p>Exercise the projection family without pretending real search, graph, or vector engines are wired yet.</p>
    <div class="mini-section"><h4>Families</h4>${list(
      projections.data.families.map((family) => `${family.name} (${family.status}) -> ${family.useFor.join(', ')}`),
    )}</div>
    <div class="mini-section"><h4>Notes</h4>${list(projections.data.notes)}</div>
    <div class="button-row">
      <button id="run-search" class="button">Search projection</button>
      <button id="run-graph" class="button">Graph projection</button>
      <button id="run-vector" class="button">Vector projection</button>
    </div>
    <pre id="projection-output">Idle</pre>
  `

  const pricingCard = document.querySelector('#pricing-card .body') as HTMLDivElement
  pricingCard.innerHTML = `
    <p>Expose the pricing contract honestly: USD only, selectable freshness, confidence/staleness, and fallback visibility.</p>
    <div class="mini-section"><h4>Notes</h4>${list(pricing.data.notes)}</div>
    <div class="button-row">
      <button id="pricing-5" class="button">5s freshness</button>
      <button id="pricing-30" class="button">30s freshness</button>
      <button id="pricing-60" class="button">60s freshness</button>
    </div>
    <pre id="pricing-output">${JSON.stringify(pricing.data, null, 2)}</pre>
  `

  const v2RoadmapCard = document.querySelector('#v2-roadmap-card .body') as HTMLDivElement
  v2RoadmapCard.innerHTML = `
    <p><strong>Status:</strong> ${v2Roadmap.data.status}</p>
    <div class="mini-section"><h4>Paradigms</h4>${list(
      v2Roadmap.data.paradigms.map((paradigm) => `${paradigm.name} (${paradigm.status}) -> ${paradigm.unlock}`),
    )}</div>
    <div class="mini-section"><h4>V1 completed</h4>${list(v2Roadmap.data.v1Completed)}</div>
    <div class="mini-section"><h4>Before implementation</h4>${list(v2Roadmap.data.nextBeforeV2Implementation)}</div>
  `

  const v2EventCard = document.querySelector('#v2-event-card .body') as HTMLDivElement
  v2EventCard.innerHTML = `
    <p><strong>Paradigm:</strong> ${v2Event.data.paradigm} (${v2Event.data.status})</p>
    <div class="mini-section"><h4>Commands</h4>${list(v2Event.data.commands)}</div>
    <div class="mini-section"><h4>Recent events</h4>${list(
      v2Event.data.events.map((event) => `${event.kind} @ ${event.emittedAt}`),
    )}</div>
    <div class="mini-section"><h4>Projection counts</h4>${list(
      Object.entries(v2Event.data.projections.eventKinds).map(([kind, count]) => `${kind}: ${count}`),
    )}</div>
    <div class="mini-section"><h4>Notes</h4>${list(v2Event.data.notes)}</div>
  `

  const v2ReactiveCard = document.querySelector('#v2-reactive-card .body') as HTMLDivElement
  v2ReactiveCard.innerHTML = `
    <p><strong>Paradigm:</strong> ${v2Reactive.data.paradigm} (${v2Reactive.data.status})</p>
    <div class="mini-section"><h4>Stream signals</h4>${list([
      `recent events: ${v2Reactive.data.streamSignals.recentEventCount}`,
      `latest event: ${v2Reactive.data.streamSignals.latestEventAt || 'n/a'}`,
      `transport: ${v2Reactive.data.streamSignals.transportSurface}`,
      ...Object.entries(v2Reactive.data.streamSignals.eventKinds).map(([kind, count]) => `${kind}: ${count}`),
    ])}</div>
    <div class="mini-section"><h4>Backpressure hints</h4>${list(v2Reactive.data.backpressureHints)}</div>
    <div class="mini-section"><h4>Notes</h4>${list(v2Reactive.data.notes)}</div>
  `

  const v2LocalFirstCard = document.querySelector('#v2-localfirst-card .body') as HTMLDivElement
  v2LocalFirstCard.innerHTML = `
    <p><strong>Paradigm:</strong> ${v2LocalFirst.data.paradigm} (${v2LocalFirst.data.status})</p>
    <div class="mini-section"><h4>Sync signals</h4>${list([
      `recent sessions: ${v2LocalFirst.data.syncSignals.recentSessionCount}`,
      `aligned sessions: ${v2LocalFirst.data.syncSignals.alignedSessions}`,
      `conflict sessions: ${v2LocalFirst.data.syncSignals.conflictSessions}`,
      `max observed lag: ${v2LocalFirst.data.syncSignals.maxObservedLag}`,
    ])}</div>
    <div class="mini-section"><h4>Merge rules</h4>${list(v2LocalFirst.data.mergeRules)}</div>
    <div class="mini-section"><h4>Notes</h4>${list(v2LocalFirst.data.notes)}</div>
  `

  const v2ConsensusCard = document.querySelector('#v2-consensus-card .body') as HTMLDivElement
  v2ConsensusCard.innerHTML = `
    <p><strong>Paradigm:</strong> ${v2Consensus.data.paradigm} (${v2Consensus.data.status})</p>
    <div class="mini-section"><h4>Coordination signals</h4>${list([
      `candidate sessions: ${v2Consensus.data.coordinationSignals.candidateSessions}`,
      `quorum ready: ${v2Consensus.data.coordinationSignals.quorumReady}`,
      `quorum blocked: ${v2Consensus.data.coordinationSignals.quorumBlocked}`,
      `rule: ${v2Consensus.data.coordinationSignals.rule}`,
    ])}</div>
    <div class="mini-section"><h4>Decision rules</h4>${list(v2Consensus.data.decisionRules)}</div>
    <div class="mini-section"><h4>Notes</h4>${list(v2Consensus.data.notes)}</div>
  `

  const v2PlanetCard = document.querySelector('#v2-planet-card .body') as HTMLDivElement
  v2PlanetCard.innerHTML = `
    <p><strong>Paradigm:</strong> ${v2Planet.data.paradigm} (${v2Planet.data.status})</p>
    <div class="mini-section"><h4>Building blocks</h4>${list(
      v2Planet.data.buildingBlocks.map((block) => `${block.name} -> ${block.role}`),
    )}</div>
    <div class="mini-section"><h4>Composition rules</h4>${list(v2Planet.data.compositionRules)}</div>
    <div class="mini-section"><h4>Notes</h4>${list(v2Planet.data.notes)}</div>
  `

  const deferredCard = document.querySelector('#deferred-card .body') as HTMLDivElement
  deferredCard.innerHTML = `
    <p><strong>Hold manifest:</strong> ${deferred.data.holdManifest}</p>
    <ul class="stack compact">${deferred.data.waves
      .map(
        (wave) => `<li><strong>${wave.name}</strong><br />${wave.purpose}<br />Streams: ${wave.streams.join(', ')}<br />Focus: ${wave.focus.join(', ')}<br />Reactivation: ${wave.reactivationRule.join(' ')}</li>`,
      )
      .join('')}</ul>
  `

  const v2Card = document.querySelector('#v2-card .body') as HTMLDivElement
  v2Card.innerHTML = `
    <p><strong>Status:</strong> ${badge(v2.data.status)}</p>
    <div class="mini-section"><h4>Blocked by</h4>${list(v2.data.blockedBy)}</div>
    <div class="mini-section"><h4>Unlock path</h4>${list(v2.data.unlockPath)}</div>
  `
}

async function runPolling(mode: 'polling' | 'long-polling') {
  const output = document.querySelector<HTMLPreElement>('#transport-output')
  if (!output) return

  output.textContent = `Starting ${mode} session...`
  const create = await fetch(apiPath(`${apiBase}/transports/polling?mode=${mode}`), { method: 'POST' })
  const created = (await create.json()) as Envelope<{ statusUrl: string; session: { recommendedPollMs: number } }>
  const states: unknown[] = [created.data.session]
  let completed = false

  while (!completed) {
    const response = await fetchEnvelope<{ status: string; timeline: string[]; step: number; totalSteps: number }>(
      apiPath(`${eventBase}${created.data.statusUrl}`),
    )
    states.push(response.data)
    output.textContent = JSON.stringify(states, null, 2)
    completed = response.data.status === 'completed'
    if (!completed && mode === 'polling') {
      await new Promise((resolve) => window.setTimeout(resolve, created.data.session.recommendedPollMs))
    }
  }
}

function runWebSocketDemo() {
  const output = document.querySelector<HTMLPreElement>('#transport-output')
  const button = document.querySelector<HTMLButtonElement>('#run-websocket')
  if (!output || !button) return

  button.disabled = true
  output.textContent = 'Opening WebSocket...'

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const socket = new WebSocket(`${protocol}//${window.location.host}/api/ws/demo`)
  const messages: unknown[] = []

  socket.onmessage = (event) => {
    try {
      messages.push(JSON.parse(event.data) as unknown)
    } catch {
      messages.push(event.data)
    }
    output.textContent = JSON.stringify(messages, null, 2)
  }

  socket.onerror = () => {
    output.textContent = 'WebSocket error'
    button.disabled = false
  }

  socket.onclose = () => {
    button.disabled = false
  }
}

async function runWorkflowDemo() {
  const output = document.querySelector<HTMLPreElement>('#workflow-output')
  if (!output) return

  output.textContent = 'Creating workflow run...'
  const create = await fetch(apiPath(`${apiBase}/messaging/workflows`), { method: 'POST' })
  const created = (await create.json()) as Envelope<{ statusUrl: string; workflow: unknown }>
  const states: unknown[] = [created.data.workflow]
  output.textContent = JSON.stringify(states, null, 2)

  let done = false
  while (!done) {
    const response = await fetchEnvelope<{ status: string }>(apiPath(`${eventBase}${created.data.statusUrl}`))
    states.push(response.data)
    output.textContent = JSON.stringify(states, null, 2)
    done = response.data.status === 'completed'
  }
}

async function updateSyncOutput(action: 'create' | 'mutate' | 'replicate', target?: 'primary' | 'replica') {
  const output = document.querySelector<HTMLPreElement>('#sync-output')
  if (!output) return

  if (action !== 'create' && !currentSyncStatusUrl) {
    output.textContent = 'Create a sync session first.'
    return
  }

  let response: Envelope<unknown>
  if (action === 'create') {
    const created = await fetch(apiPath(`${apiBase}/sync/sessions`), { method: 'POST' })
    response = (await created.json()) as Envelope<{ statusUrl: string }>
    currentSyncStatusUrl = `${eventBase}${(response.data as { statusUrl: string }).statusUrl}`
  } else if (action == 'mutate') {
    const mutated = await fetch(apiPath(`${currentSyncStatusUrl}/mutate?target=${target}`), { method: 'POST' })
    response = (await mutated.json()) as Envelope<unknown>
  } else {
    const replicated = await fetch(apiPath(`${currentSyncStatusUrl}/replicate`), { method: 'POST' })
    response = (await replicated.json()) as Envelope<unknown>
  }

  output.textContent = JSON.stringify(response.data, null, 2)
}

async function updateProjectionOutput(mode: 'search' | 'graph' | 'vector') {
  const output = document.querySelector<HTMLPreElement>('#projection-output')
  if (!output) return

  let path = apiPath(`${apiBase}/projections/graph`)
  if (mode === 'search') {
    path = apiPath(`${apiBase}/projections/search?q=workflow`)
  } else if (mode === 'vector') {
    path = apiPath(`${apiBase}/projections/vector?q=sync merge`)
  }

  output.textContent = 'Loading projection...'
  const response = await fetchEnvelope<unknown>(path)
  output.textContent = JSON.stringify(response.data, null, 2)
}

async function updatePricingOutput(freshness: 5 | 30 | 60) {
  const output = document.querySelector<HTMLPreElement>('#pricing-output')
  if (!output) return
  output.textContent = 'Loading pricing...'
  const response = await fetchEnvelope<PricingOverview>(apiPath(`${apiBase}/pricing?freshness=${freshness}`))
  output.textContent = JSON.stringify(response.data, null, 2)
}

function startEvents() {
  const listNode = document.querySelector('#events-list') as HTMLUListElement
  const source = new EventSource(`${eventBase}/api/events`)

  source.addEventListener('progress', (event) => {
    let label = event.data
    try {
      const parsed = JSON.parse(event.data) as { label?: string; delayAppliedMs?: number; step?: string }
      label = `${parsed.label ?? parsed.step ?? 'event'}${parsed.delayAppliedMs ? ` (${parsed.delayAppliedMs} ms delay)` : ''}`
    } catch {
      label = event.data
    }
    const item = document.createElement('li')
    item.textContent = label
    listNode.prepend(item)
  })

  source.onerror = () => {
    const item = document.createElement('li')
    item.textContent = 'Event stream ended.'
    listNode.prepend(item)
    source.close()
  }
}

function wireAsyncDemo() {
  const button = document.querySelector<HTMLButtonElement>('#run-async')
  const output = document.querySelector<HTMLPreElement>('#async-output')
  if (!button || !output) {
    return
  }

  button.addEventListener('click', async () => {
    button.disabled = true
    output.textContent = 'Submitting...'

    try {
      const response = await fetch(`${apiBase}/async/demo`, { method: 'POST' })
      const envelope = (await response.json()) as Envelope<{ statusUrl: string; state: JobState }>
      output.textContent = JSON.stringify(envelope.data.state, null, 2)

      const interval = window.setInterval(async () => {
        const state = await fetchEnvelope<JobState>(`${eventBase}${envelope.data.statusUrl}`)
        output.textContent = JSON.stringify(state.data, null, 2)
        if (state.data.status === 'succeeded' || state.data.status === 'failed') {
          window.clearInterval(interval)
          button.disabled = false
        }
      }, 1500)
    } catch (error) {
      output.textContent = error instanceof Error ? error.message : 'Unknown error'
      button.disabled = false
    }
  })
}

function wireTransportDemos() {
  const pollingButton = document.querySelector<HTMLButtonElement>('#run-polling')
  const longPollingButton = document.querySelector<HTMLButtonElement>('#run-long-polling')
  const websocketButton = document.querySelector<HTMLButtonElement>('#run-websocket')
  if (!pollingButton || !longPollingButton || !websocketButton) return

  pollingButton.addEventListener('click', async () => {
    pollingButton.disabled = true
    try {
      await runPolling('polling')
    } finally {
      pollingButton.disabled = false
    }
  })

  longPollingButton.addEventListener('click', async () => {
    longPollingButton.disabled = true
    try {
      await runPolling('long-polling')
    } finally {
      longPollingButton.disabled = false
    }
  })

  websocketButton.addEventListener('click', runWebSocketDemo)
}

function wireWorkflowDemo() {
  const button = document.querySelector<HTMLButtonElement>('#run-workflow')
  if (!button) return

  button.addEventListener('click', async () => {
    button.disabled = true
    try {
      await runWorkflowDemo()
    } finally {
      button.disabled = false
    }
  })
}

function wireSyncDemo() {
  const createButton = document.querySelector<HTMLButtonElement>('#create-sync')
  const mutatePrimaryButton = document.querySelector<HTMLButtonElement>('#mutate-primary')
  const mutateReplicaButton = document.querySelector<HTMLButtonElement>('#mutate-replica')
  const replicateButton = document.querySelector<HTMLButtonElement>('#replicate-sync')
  if (!createButton || !mutatePrimaryButton || !mutateReplicaButton || !replicateButton) return

  createButton.addEventListener('click', () => updateSyncOutput('create'))
  mutatePrimaryButton.addEventListener('click', () => updateSyncOutput('mutate', 'primary'))
  mutateReplicaButton.addEventListener('click', () => updateSyncOutput('mutate', 'replica'))
  replicateButton.addEventListener('click', () => updateSyncOutput('replicate'))
}

function wireLocaleSwitcher() {
  const select = document.querySelector<HTMLSelectElement>('#locale-select')
  if (!select) return
  select.value = currentLocale
  select.addEventListener('change', async () => {
    currentLocale = select.value as Locale
    await renderApp()
  })
}

function wireProjectionDemo() {
  const searchButton = document.querySelector<HTMLButtonElement>('#run-search')
  const graphButton = document.querySelector<HTMLButtonElement>('#run-graph')
  const vectorButton = document.querySelector<HTMLButtonElement>('#run-vector')
  if (!searchButton || !graphButton || !vectorButton) return

  searchButton.addEventListener('click', () => updateProjectionOutput('search'))
  graphButton.addEventListener('click', () => updateProjectionOutput('graph'))
  vectorButton.addEventListener('click', () => updateProjectionOutput('vector'))
}

function wirePricingDemo() {
  const pricing5 = document.querySelector<HTMLButtonElement>('#pricing-5')
  const pricing30 = document.querySelector<HTMLButtonElement>('#pricing-30')
  const pricing60 = document.querySelector<HTMLButtonElement>('#pricing-60')
  if (!pricing5 || !pricing30 || !pricing60) return

  pricing5.addEventListener('click', () => updatePricingOutput(5))
  pricing30.addEventListener('click', () => updatePricingOutput(30))
  pricing60.addEventListener('click', () => updatePricingOutput(60))
}

renderApp().catch((error) => {
  root.innerHTML = `<main class="shell"><section class="panel"><h1>Bootstrap failed</h1><pre>${error instanceof Error ? error.message : 'Unknown error'}</pre></section></main>`
})

startEvents()
wireAsyncDemo()
wireTransportDemos()
wireWorkflowDemo()
wireSyncDemo()
wireProjectionDemo()
wirePricingDemo()
wireLocaleSwitcher()
