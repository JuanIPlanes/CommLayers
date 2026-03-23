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
  description: string
  status: string
  dependencies: string[]
}

type BootstrapData = {
  summary: string
  runtime: {
    api: string
    frontend: string
    mandatoryInfra: string[]
    optionalInfra: string[]
  }
  auth: {
    identityProvider: string
    secretBroker: string
    currentState: string
  }
}

type CatalogEntry = {
  name: string
  status: string
  recommendedWhen: string[]
  avoidWhen: string[]
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

const app = document.querySelector<HTMLDivElement>('#app')

if (!app) {
  throw new Error('Missing app root')
}

app.innerHTML = `
  <main class="shell">
    <header class="hero">
      <p class="eyebrow">CommLayers</p>
      <h1>Backend-first bootstrap</h1>
      <p class="lede">A minimal first runnable version that exposes the first-wave foundation, SSE progress, and async demo behavior without pretending the full platform already exists.</p>
    </header>

    <section class="grid two-up">
      <article class="card" id="bootstrap-card"><h2>Bootstrap summary</h2><div class="body">Loading...</div></article>
      <article class="card" id="first-wave-card"><h2>First-wave streams</h2><div class="body">Loading...</div></article>
    </section>

    <section class="grid two-up">
      <article class="card" id="catalog-card"><h2>Catalog snapshot</h2><div class="body">Loading...</div></article>
      <article class="card" id="events-card"><h2>SSE progress demo</h2><div class="body"><ul id="events-list" class="stack"></ul></div></article>
    </section>

    <section class="grid two-up">
      <article class="card" id="async-card">
        <h2>Async demo</h2>
        <div class="body stack">
          <p>Start a simulated long-running operation with visible progress and delay metadata.</p>
          <button id="run-async" class="button">Run async demo</button>
          <pre id="async-output">Idle</pre>
        </div>
      </article>
      <article class="card" id="notes-card">
        <h2>What is deferred</h2>
        <div class="body">
          <ul class="stack compact">
            <li>Transport deepening beyond the bootstrap demo</li>
            <li>Messaging, sync, and projection execution</li>
            <li>Frontend comparative UI polish</li>
            <li>Pricing and v2 paradigm work</li>
          </ul>
        </div>
      </article>
    </section>
  </main>
`

async function fetchEnvelope<T>(path: string): Promise<Envelope<T>> {
  const response = await fetch(path)
  if (!response.ok) {
    throw new Error(`Request failed: ${response.status}`)
  }
  return response.json() as Promise<Envelope<T>>
}

async function renderBootstrap() {
  const bootstrap = await fetchEnvelope<BootstrapData>(`${apiBase}/bootstrap`)
  const firstWave = await fetchEnvelope<{ activeStreams: Stream[]; deferredStreams: Stream[]; executionOrder: string[] }>(`${apiBase}/first-wave`)
  const catalog = await fetchEnvelope<{ families: CatalogEntry[]; notes: string[] }>(`${apiBase}/catalog`)

  const bootstrapCard = document.querySelector('#bootstrap-card .body') as HTMLDivElement
  bootstrapCard.innerHTML = `
    <p>${bootstrap.data.summary}</p>
    <dl class="meta">
      <div><dt>API</dt><dd>${bootstrap.data.runtime.api}</dd></div>
      <div><dt>Frontend</dt><dd>${bootstrap.data.runtime.frontend}</dd></div>
      <div><dt>Mandatory infra</dt><dd>${bootstrap.data.runtime.mandatoryInfra.join(', ')}</dd></div>
      <div><dt>Optional infra</dt><dd>${bootstrap.data.runtime.optionalInfra.join(', ')}</dd></div>
      <div><dt>Identity</dt><dd>${bootstrap.data.auth.identityProvider}</dd></div>
      <div><dt>Secrets</dt><dd>${bootstrap.data.auth.secretBroker}</dd></div>
      <div><dt>Current state</dt><dd>${bootstrap.data.auth.currentState}</dd></div>
    </dl>
  `

  const firstWaveCard = document.querySelector('#first-wave-card .body') as HTMLDivElement
  firstWaveCard.innerHTML = `
    <p><strong>Execution order:</strong> ${firstWave.data.executionOrder.join(' -> ')}</p>
    <ul class="stack compact">
      ${firstWave.data.activeStreams
        .map(
          (stream) => `<li><strong>${stream.id}</strong> - ${stream.title}<br /><span>${stream.description}</span></li>`,
        )
        .join('')}
    </ul>
  `

  const catalogCard = document.querySelector('#catalog-card .body') as HTMLDivElement
  catalogCard.innerHTML = `
    <ul class="stack compact">
      ${catalog.data.families
        .map(
          (item) => `<li><strong>${item.name}</strong> <span class="chip">${item.status}</span><br />Use: ${item.recommendedWhen.join(', ')}<br />Avoid: ${item.avoidWhen.join(', ')}</li>`,
        )
        .join('')}
    </ul>
  `
}

function startEvents() {
  const list = document.querySelector('#events-list') as HTMLUListElement
  const source = new EventSource(`${eventBase}/api/events`)

  source.addEventListener('progress', (event) => {
    const item = document.createElement('li')
    item.textContent = event.data
    list.prepend(item)
  })

  source.onerror = () => {
    const item = document.createElement('li')
    item.textContent = 'Event stream ended.'
    list.prepend(item)
    source.close()
  }
}

function wireAsyncDemo() {
  const button = document.querySelector<HTMLButtonElement>('#run-async')
  const output = document.querySelector<HTMLPreElement>('#async-output')
  if (!button || !output) return

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

renderBootstrap().catch((error) => {
  app.innerHTML = `<main class="shell"><div class="card"><h1>Bootstrap failed</h1><pre>${error instanceof Error ? error.message : 'Unknown error'}</pre></div></main>`
})

startEvents()
wireAsyncDemo()
