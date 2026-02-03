// 开发时默认走 Vite 代理 (/api -> localhost:8080)；若代理不生效，可设 VITE_API_URL=http://localhost:8080 直连后端
const API_BASE = import.meta.env.VITE_API_URL || ''

export function apiUrl(path) {
  const p = path.startsWith('/') ? path : `/${path}`
  return `${API_BASE}${p}`
}
