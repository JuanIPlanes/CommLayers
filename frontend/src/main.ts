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

const root = document.querySelector<HTMLDivElement>('#app')

if (!root) {
  throw new Error('Missing app root')
}

root.innerHTML = `
  <main class="shell">
    <header class="hero panel">
      <div class="hero-copy">
        <p class="eyebrow">CommLayers</p>
        <h1>First-wave execution console</h1>
        <p class="lede">This bootstrap turns the backend-first planning package into a runnable foundation: explicit first-wave contracts, deferred-wave holds, and visible timing behavior.</p>
      </div>
      <div class="hero-meta" id="hero-meta">Loading...</div>
    </header>

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
  const [bootstrap, firstWave, security, dataPlatform, benchmark, catalog, realtime, deferred, v2] = await Promise.all([
    fetchEnvelope<BootstrapData>(`${apiBase}/bootstrap`),
    fetchEnvelope<FirstWaveContract>(`${apiBase}/first-wave/contract`),
    fetchEnvelope<SecurityBootstrap>(`${apiBase}/security/bootstrap`),
    fetchEnvelope<DataPlatform>(`${apiBase}/data-platform`),
    fetchEnvelope<BenchmarkFramework>(`${apiBase}/benchmark-framework`),
    fetchEnvelope<CatalogData>(`${apiBase}/catalog`),
    fetchEnvelope<RealtimeComparisons>(`${apiBase}/comparisons/realtime`),
    fetchEnvelope<DeferredWavesData>(`${apiBase}/deferred-waves`),
    fetchEnvelope<V2Readiness>(`${apiBase}/v2-readiness`),
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

renderApp().catch((error) => {
  root.innerHTML = `<main class="shell"><section class="panel"><h1>Bootstrap failed</h1><pre>${error instanceof Error ? error.message : 'Unknown error'}</pre></section></main>`
})

startEvents()
wireAsyncDemo()
