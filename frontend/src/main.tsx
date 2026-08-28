import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import './index.css'
import App from './App.tsx'

// Aquece o backend (Render free tier "spin down" após inatividade — medido em
// produção: ~21s frio, ~1s quente, independente do banco, que já está no
// Neon) assim que o app carrega, em paralelo com o jogador lendo a tela de
// login/lendo a página — quando ele de fato submeter o primeiro request real,
// uma boa parte (às vezes todo) o cold start já aconteceu. Usa fetch direto
// (não a instância `api`/axios) de propósito: não deve disparar o toast de
// "Acordando o servidor" nem contar no contador de requisições em voo — é
// silencioso, e o toast continua existindo como rede de segurança pro caso
// desse aquecimento não terminar a tempo do primeiro request real.
const apiBase = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'
fetch(apiBase.replace(/\/api\/?$/, '') + '/health').catch(() => {})

const queryClient = new QueryClient()

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  </StrictMode>,
)