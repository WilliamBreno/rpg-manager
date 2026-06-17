import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
  timeout: 30000, // 30s — aguenta o cold start
  headers: {
    'Content-Type': 'application/json',
  },
})

// ── Toast de cold start ─────────────────────────────────────────────────────
const TOAST_ID = 'cold-start-toast'

function showWakingToast() {
  if (document.getElementById(TOAST_ID)) return
  const toast = document.createElement('div')
  toast.id = TOAST_ID
  toast.innerHTML = `
    <div style="display:flex;align-items:center;gap:10px">
      <div style="width:16px;height:16px;border:2px solid #facc15;border-top-color:transparent;border-radius:50%;animation:spin .8s linear infinite;flex-shrink:0"></div>
      <div>
        <p style="font-weight:600;font-size:14px;margin:0">Acordando o servidor...</p>
        <p style="font-size:12px;color:#9ca3af;margin:0">Pode demorar até 30s na primeira vez</p>
      </div>
    </div>
    <style>@keyframes spin{to{transform:rotate(360deg)}}</style>
  `
  Object.assign(toast.style, {
    position: 'fixed',
    bottom: '24px',
    left: '50%',
    transform: 'translateX(-50%)',
    background: '#1f2937',
    border: '1px solid #374151',
    borderLeft: '4px solid #facc15',
    color: '#fff',
    padding: '12px 20px',
    borderRadius: '10px',
    zIndex: '9999',
    boxShadow: '0 4px 24px rgba(0,0,0,0.4)',
    minWidth: '280px',
    fontFamily: 'inherit',
  })
  document.body.appendChild(toast)
}

function hideWakingToast() {
  document.getElementById(TOAST_ID)?.remove()
}

// ── Interceptors ────────────────────────────────────────────────────────────
let slowTimer: ReturnType<typeof setTimeout> | null = null

api.interceptors.request.use(config => {
  // Token
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`

  // Mostra toast se demorar mais de 4s
  slowTimer = setTimeout(showWakingToast, 4000)

  return config
})

api.interceptors.response.use(
  response => {
    if (slowTimer) { clearTimeout(slowTimer); slowTimer = null }
    hideWakingToast()
    return response
  },
  error => {
    if (slowTimer) { clearTimeout(slowTimer); slowTimer = null }
    hideWakingToast()

    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api