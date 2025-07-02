import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import App from './App'
import DashboardPage from './pages/Dashboard'
import { ThemeProvider } from './context/theme-provider'

createRoot(document.getElementById('root')!).render(
  <ThemeProvider defaultTheme="system" storageKey="duite-bot-theme">
    <StrictMode>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<App />} />
          <Route path="/dashboard" element={<DashboardPage />} />
        </Routes>
      </BrowserRouter>
    </StrictMode>
  </ThemeProvider>
)
