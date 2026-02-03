// Dev: relative /api so browser hits same origin (localhost:3000), Vite proxy forwards to backend. No direct 127.0.0.1 — avoids 502/504 from system proxy.
// Prod: same origin (Caddy/nginx serves both).
const API_BASE = ''

export function apiUrl(path) {
  const p = path.startsWith('/') ? path : `/${path}`
  return `${API_BASE}${p}`
}
