import { Link, useLocation } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import { Home, Gift, History, LogOut } from 'lucide-react'

export default function Header() {
  const { isAuthenticated, logout } = useAuthStore()
  const location = useLocation()

  if (!isAuthenticated) {
    return null
  }

  const isActive = (path: string) => location.pathname === path

  return (
    <header className="bg-white shadow-sm border-b">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          <div className="flex items-center space-x-8">
            <h1 className="text-xl font-bold text-gray-900">奖励系统</h1>
            
            <nav className="flex space-x-6">
              <Link
                to="/"
                className={`flex items-center space-x-2 px-3 py-2 rounded-md text-sm font-medium ${
                  isActive('/') 
                    ? 'bg-primary-100 text-primary-700' 
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                <Home className="h-4 w-4" />
                <span>仪表板</span>
              </Link>
              
              <Link
                to="/reward-types"
                className={`flex items-center space-x-2 px-3 py-2 rounded-md text-sm font-medium ${
                  isActive('/reward-types') 
                    ? 'bg-primary-100 text-primary-700' 
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                <Gift className="h-4 w-4" />
                <span>奖励类型</span>
              </Link>
              
              <Link
                to="/transactions"
                className={`flex items-center space-x-2 px-3 py-2 rounded-md text-sm font-medium ${
                  isActive('/transactions') 
                    ? 'bg-primary-100 text-primary-700' 
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                <History className="h-4 w-4" />
                <span>交易记录</span>
              </Link>
            </nav>
          </div>
          
          <button
            onClick={logout}
            className="flex items-center space-x-2 px-3 py-2 rounded-md text-sm font-medium text-gray-600 hover:text-gray-900"
          >
            <LogOut className="h-4 w-4" />
            <span>退出</span>
          </button>
        </div>
      </div>
    </header>
  )
}