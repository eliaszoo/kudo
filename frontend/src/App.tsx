import { useState, useEffect } from 'react'
import { Routes, Route } from 'react-router-dom'
import { Toaster } from 'sonner'
import Header from './components/Header'
import Dashboard from './pages/Dashboard'
import RewardTypes from './pages/RewardTypes'
import Transactions from './pages/Transactions'
import { useAuthStore } from './stores/authStore'

function App() {
  const { checkAuth } = useAuthStore()
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    checkAuth()
    setLoading(false)
  }, [checkAuth])

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-primary-500"></div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <main className="container mx-auto px-4 py-8">
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/reward-types" element={<RewardTypes />} />
          <Route path="/transactions" element={<Transactions />} />
        </Routes>
      </main>
      <Toaster position="top-right" />
    </div>
  )
}

export default App