import { create } from 'zustand'

interface AuthState {
  token: string | null
  isAuthenticated: boolean
  login: (token: string) => void
  logout: () => void
  checkAuth: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem('token') || 'demo-token',
  isAuthenticated: true,
  login: (token: string) => {
    localStorage.setItem('token', token)
    set({ token, isAuthenticated: true })
  },
  logout: () => {
    localStorage.removeItem('token')
    set({ token: null, isAuthenticated: false })
  },
  checkAuth: () => {
    const token = localStorage.getItem('token') || 'demo-token'
    localStorage.setItem('token', token)
    set({ token, isAuthenticated: true })
  },
}))