import React, { useEffect, useState } from 'react'
import { createRoot } from 'react-dom/client'
import { Activity, AlertTriangle, Bot, Cloud, CreditCard, GitBranch, LayoutDashboard, Lock, Network, Rocket, Search, Shield, Server, Settings } from 'lucide-react'
import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'
import '../app/styles.css'

const API = import.meta.env.VITE_API_URL ?? ''

type Resource = {
  id: string
  name: string
  provider: string
  accountId: string
  region: string
  zone: string
  resourceType: string
  category: string
  status: string
  health: string
  costMonthly: number
  tags: Record<string, string>
  metadata: Record<string, unknown>
  updatedAt: string
}

type CloudAccount = { id: string; name: string; provider: string; accountId: string; defaultRegion: string; status: string }
type Alert = { id: string; severity: string; title: string; source: string; status: string }
type Finding = { id: string; severity: string; title: string; provider: string; explanation: string; fix: string }

function token() {
  return localStorage.getItem('infrasphere.token') ?? ''
}

async function api<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API}${path}`, {
    ...init,
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}`, ...(init?.headers ?? {}) }
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

function App() {
  const [authed, setAuthed] = useState(Boolean(token()))
  const [page, setPage] = useState('dashboard')
  if (!authed) return <Login onLogin={() => setAuthed(true)} />
  return (
    <Shell page={page} setPage={setPage}>
      {page === 'dashboard' && <Dashboard />}
      {page === 'connections' && <Connections />}
      {page === 'inventory' && <Inventory />}
      {page === 'topology' && <Topology />}
      {page === 'deploy' && <Deploy />}
      {page === 'observability' && <Observability />}
      {page === 'ai' && <Assistant />}
      {page === 'cost' && <Cost />}
      {page === 'security' && <Security />}
      {page === 'settings' && <SettingsPage />}
    </Shell>
  )
}

function Login({ onLogin }: { onLogin: () => void }) {
  const [email, setEmail] = useState('admin@infrasphere.local')
  const [password, setPassword] = useState('ChangeMe123!')
  const [error, setError] = useState('')
  async function submit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    try {
      const res = await fetch(`${API}/api/auth/login`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ email, password }) })
      if (!res.ok) throw new Error('Login failed')
      const data = await res.json()
      localStorage.setItem('infrasphere.token', data.token)
      onLogin()
    } catch (err) {
      setError(String(err))
    }
  }
  return (
    <main className="login">
      <form className="loginPanel" onSubmit={submit}>
        <div className="brand"><Cloud size={30} /> <span>InfraSphere</span></div>
        <h1>Control Plane</h1>
        <label>Email<input value={email} onChange={e => setEmail(e.target.value)} /></label>
        <label>Password<input type="password" value={password} onChange={e => setPassword(e.target.value)} /></label>
        {error && <p className="error">{error}</p>}
        <button><Lock size={16} /> Sign in</button>
      </form>
    </main>
  )
}

function Shell({ page, setPage, children }: { page: string; setPage: (p: string) => void; children: React.ReactNode }) {
  const nav = [
    ['dashboard', LayoutDashboard, 'Dashboard'], ['connections', Cloud, 'Connections'], ['inventory', Server, 'Inventory'],
    ['topology', GitBranch, 'Topology'], ['deploy', Rocket, 'Deploy'], ['observability', Activity, 'Observability'],
    ['ai', Bot, 'AI Assistant'], ['cost', CreditCard, 'Cost'], ['security', Shield, 'Security'], ['settings', Settings, 'Admin']
  ] as const
  return (
    <div className="app">
      <aside>
        <div className="logo"><Cloud /> InfraSphere</div>
        {nav.map(([id, Icon, label]) => <button className={page === id ? 'active' : ''} key={id} onClick={() => setPage(id)}><Icon size={17} />{label}</button>)}
      </aside>
      <section className="workspace">
        <header><div><b>Acme Platform Engineering</b><span>Unified hybrid and multi-cloud operations</span></div><button onClick={() => { localStorage.clear(); location.reload() }}>Sign out</button></header>
        {children}
      </section>
    </div>
  )
}

function useResources() {
  const [resources, setResources] = useState<Resource[]>([])
  useEffect(() => { api<Resource[]>('/api/inventory/resources').then(setResources).catch(console.error) }, [])
  return resources
}

function Dashboard() {
  const resources = useResources()
  const [alerts, setAlerts] = useState<Alert[]>([])
  useEffect(() => { api<Alert[]>('/api/observability/alerts').then(setAlerts).catch(console.error) }, [])
  const cost = resources.reduce((n, r) => n + r.costMonthly, 0)
  const chart = Object.entries(resources.reduce<Record<string, number>>((m, r) => ({ ...m, [r.provider]: (m[r.provider] ?? 0) + 1 }), {})).map(([provider, count]) => ({ provider, count }))
  return <Page title="Dashboard" subtitle="Operational truth across cloud, private cloud, cost, security, and reliability.">
    <div className="kpis">
      <Kpi label="Resources" value={resources.length} icon={<Server />} />
      <Kpi label="Monthly cost" value={`$${cost.toLocaleString()}`} icon={<CreditCard />} />
      <Kpi label="Critical alerts" value={alerts.filter(a => a.severity === 'critical').length} icon={<AlertTriangle />} />
      <Kpi label="Risk findings" value="2" icon={<Shield />} />
    </div>
    <div className="grid two">
      <Panel title="Resources by provider"><ResponsiveContainer width="100%" height={220}><BarChart data={chart}><CartesianGrid strokeDasharray="3 3" /><XAxis dataKey="provider" /><YAxis allowDecimals={false} /><Tooltip /><Bar dataKey="count" fill="#2563eb" radius={4} /></BarChart></ResponsiveContainer></Panel>
      <Panel title="Recent alerts">{alerts.map(a => <Row key={a.id} title={a.title} meta={`${a.severity} · ${a.source} · ${a.status}`} />)}</Panel>
    </div>
  </Page>
}

function Connections() {
  const [accounts, setAccounts] = useState<CloudAccount[]>([])
  useEffect(() => { api<CloudAccount[]>('/api/cloud-accounts').then(setAccounts).catch(console.error) }, [])
  return <Page title="Cloud Connections" subtitle="Register AWS, GCP, Azure, VCF, edge, and future providers behind one model.">
    <Panel title="Connected environments">{accounts.map(a => <Row key={a.id} title={a.name} meta={`${a.provider} · ${a.accountId} · ${a.defaultRegion} · ${a.status}`} />)}</Panel>
  </Page>
}

function Inventory() {
  const resources = useResources()
  const [query, setQuery] = useState('')
  const [selected, setSelected] = useState<Resource | null>(null)
  const filtered = resources.filter(r => `${r.name} ${r.provider} ${r.resourceType} ${r.region}`.toLowerCase().includes(query.toLowerCase()))
  return <Page title="Inventory" subtitle="Search and inspect normalized resources from every connected environment.">
    <div className="toolbar"><Search size={16} /><input placeholder="Search by name, provider, type, region, tag" value={query} onChange={e => setQuery(e.target.value)} /></div>
    <div className="table">
      <div className="thead"><span>Name</span><span>Provider</span><span>Type</span><span>Region</span><span>Health</span><span>Cost</span></div>
      {filtered.map(r => <button className="tr" key={r.id} onClick={() => setSelected(r)}><span>{r.name}</span><span className="badge">{r.provider}</span><span>{r.resourceType}</span><span>{r.region}</span><span className={r.health}>{r.health}</span><span>${r.costMonthly}</span></button>)}
    </div>
    {selected && <ResourceDetail resource={selected} onClose={() => setSelected(null)} />}
  </Page>
}

function ResourceDetail({ resource, onClose }: { resource: Resource; onClose: () => void }) {
  return <div className="drawer"><button className="close" onClick={onClose}>Close</button><h2>{resource.name}</h2><p>{resource.provider} · {resource.region} · {resource.status}</p>
    <div className="grid two"><Panel title="Metadata"><pre>{JSON.stringify(resource.metadata, null, 2)}</pre></Panel><Panel title="AI explanation"><p>This resource is part of the production payments path. Cost, health, security posture, and dependency relationships should be reviewed together before changes.</p></Panel></div>
  </div>
}

function Topology() {
  const resources = useResources()
  return <Page title="Network Topology" subtitle="Relationship graph placeholder for resource dependencies and workload blast radius.">
    <div className="topology">{resources.map(r => <div key={r.id} className={`node ${r.category.toLowerCase()}`}><Network size={16} />{r.name}<small>{r.category}</small></div>)}</div>
  </Page>
}

function Deploy() {
  const [status, setStatus] = useState('')
  async function create() {
    const d = await api('/api/deployments', { method: 'POST', body: JSON.stringify({ name: 'payments-api', provider: 'aws', target: 'prod-aws/us-west-2', workloadType: 'container', costEstimate: 482 }) })
    setStatus(`Created deployment ${(d as any).id} in DRAFT`)
  }
  return <Page title="Deployment Wizard" subtitle="Provider-neutral workload planning with approval gates before mutation.">
    <div className="wizard">{['Provider', 'Target', 'Region', 'Workload', 'Networking', 'Storage', 'Scaling', 'Observability', 'Cost', 'Review'].map((s, i) => <div key={s} className="step"><b>{i + 1}</b>{s}</div>)}</div>
    <Panel title="Example workload"><pre>{`apiVersion: infrasphere.io/v1
kind: Workload
metadata:
  name: payments-api
spec:
  provider: aws
  target:
    account: prod-aws
    region: us-west-2
    platform: eks
  container:
    image: ghcr.io/example/payments-api:v1.2.3
    port: 8080`}</pre><button onClick={create}><Rocket size={16} /> Create draft deployment</button>{status && <p>{status}</p>}</Panel>
  </Page>
}

function Observability() {
  const [alerts, setAlerts] = useState<Alert[]>([])
  useEffect(() => { api<Alert[]>('/api/observability/alerts').then(setAlerts).catch(console.error) }, [])
  return <Page title="Observability" subtitle="Metrics, logs, traces, alerts, SLOs, and deployment correlation."><Panel title="Alerts">{alerts.map(a => <Row key={a.id} title={a.title} meta={`${a.severity} · ${a.source}`} />)}</Panel></Page>
}

function Assistant() {
  const [message, setMessage] = useState('Which workloads are exposed to the internet?')
  const [answer, setAnswer] = useState('')
  async function ask() {
    const res = await api<{ answer: string; citations: string[] }>('/api/ai/chat', { method: 'POST', body: JSON.stringify({ message }) })
    setAnswer(`${res.answer}\n\nSources: ${res.citations.join(', ')}`)
  }
  return <Page title="AI Assistant" subtitle="Read-only infrastructure reasoning with explicit approval before actions."><Panel title="Ask InfraSphere"><div className="chat"><textarea value={message} onChange={e => setMessage(e.target.value)} /><button onClick={ask}><Bot size={16} /> Ask</button>{answer && <pre>{answer}</pre>}</div></Panel></Page>
}

function Cost() {
  const resources = useResources()
  return <Page title="Cost" subtitle="FinOps visibility by provider, resource, workload, account, and tag."><Panel title="Top cost drivers">{resources.sort((a, b) => b.costMonthly - a.costMonthly).map(r => <Row key={r.id} title={r.name} meta={`${r.provider} · $${r.costMonthly}/mo · ${r.category}`} />)}</Panel></Page>
}

function Security() {
  const [findings, setFindings] = useState<Finding[]>([])
  useEffect(() => { api<Finding[]>('/api/security/findings').then(setFindings).catch(console.error) }, [])
  return <Page title="Security Posture" subtitle="Find exposure, identity, encryption, region, and tagging risks across providers."><Panel title="Findings">{findings.map(f => <Row key={f.id} title={f.title} meta={`${f.severity} · ${f.provider} · ${f.fix}`} />)}</Panel></Page>
}

function SettingsPage() {
  return <Page title="Admin Settings" subtitle="RBAC, OIDC, API keys, audit logs, approvals, and tenant controls."><Panel title="Enterprise controls"><Row title="OIDC / OAuth2" meta="Configure issuer, client, groups, and role mapping." /><Row title="Approval workflows" meta="Require human approval for destructive production actions." /><Row title="Credential encryption" meta="Provider secrets are write-only and redacted from APIs." /></Panel></Page>
}

function Page({ title, subtitle, children }: { title: string; subtitle: string; children: React.ReactNode }) {
  return <main><div className="pageTitle"><h1>{title}</h1><p>{subtitle}</p></div>{children}</main>
}
function Kpi({ label, value, icon }: { label: string; value: React.ReactNode; icon: React.ReactNode }) { return <div className="kpi">{icon}<span>{label}</span><b>{value}</b></div> }
function Panel({ title, children }: { title: string; children: React.ReactNode }) { return <section className="panel"><h2>{title}</h2>{children}</section> }
function Row({ title, meta }: { title: string; meta: string }) { return <div className="row"><b>{title}</b><span>{meta}</span></div> }

createRoot(document.getElementById('root')!).render(<App />)
