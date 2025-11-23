import { useState, useEffect } from 'react'
import { useFamilyStore } from '../stores/familyStore'
import { TrendingUp, TrendingDown, Calendar, User } from 'lucide-react'

interface Transaction {
  id: number
  type: 'credit' | 'debit'
  value: number
  note?: string
  created_at: string
  account_id: number
}

export default function Transactions() {
  const { users, rewardTypes, loading } = useFamilyStore()
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [filters, setFilters] = useState({
    child_id: '',
    reward_type_id: '',
    type: '',
  })

  useEffect(() => {
    // Mock transactions data
    const mockTransactions: Transaction[] = [
      {
        id: 1,
        type: 'credit',
        value: 10000,
        note: '完成作业',
        created_at: new Date().toISOString(),
        account_id: 1,
      },
      {
        id: 2,
        type: 'debit',
        value: 5000,
        note: '买文具',
        created_at: new Date(Date.now() - 86400000).toISOString(),
        account_id: 1,
      },
    ]
    setTransactions(mockTransactions)
  }, [])

  const getTypeIcon = (type: string) => {
    return type === 'credit' ? 
      <TrendingUp className="h-4 w-4 text-success-500" /> : 
      <TrendingDown className="h-4 w-4 text-warning-500" />
  }

  const getTypeLabel = (type: string) => {
    return type === 'credit' ? '授予' : '消费'
  }

  const formatValue = (value: number, unitKind: string) => {
    if (unitKind === 'money') {
      return `¥${(value / 100).toFixed(2)}`
    }
    return value.toString()
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const filteredTransactions = transactions.filter(tx => {
    if (filters.type && tx.type !== filters.type) return false
    return true
  })

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">交易记录</h1>
        <p className="text-gray-600 mt-1">查看和管理所有奖励交易记录</p>
      </div>

      {/* Filters */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">筛选条件</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              孩子
            </label>
            <select
              value={filters.child_id}
              onChange={(e) => setFilters({ ...filters, child_id: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              <option value="">全部孩子</option>
              {users.filter(u => u.role === 'child').map(child => (
                <option key={child.id} value={child.id}>{child.display_name}</option>
              ))}
            </select>
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              奖励类型
            </label>
            <select
              value={filters.reward_type_id}
              onChange={(e) => setFilters({ ...filters, reward_type_id: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              <option value="">全部类型</option>
              {rewardTypes.map(type => (
                <option key={type.id} value={type.id}>{type.name}</option>
              ))}
            </select>
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              交易类型
            </label>
            <select
              value={filters.type}
              onChange={(e) => setFilters({ ...filters, type: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
            >
              <option value="">全部类型</option>
              <option value="credit">授予</option>
              <option value="debit">消费</option>
            </select>
          </div>
        </div>
      </div>

      {/* Transactions List */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <h2 className="text-lg font-semibold text-gray-900">交易记录</h2>
          <p className="text-sm text-gray-600 mt-1">
            共 {filteredTransactions.length} 条记录
          </p>
        </div>
        
        {loading ? (
          <div className="p-8 text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500 mx-auto"></div>
            <p className="mt-2 text-gray-600">加载中...</p>
          </div>
        ) : filteredTransactions.length === 0 ? (
          <div className="p-8 text-center text-gray-500">
            <TrendingUp className="h-12 w-12 mx-auto mb-4 text-gray-300" />
            <p>暂无交易记录</p>
            <p className="text-sm mt-1">开始授予或消费奖励来创建记录</p>
          </div>
        ) : (
          <div className="divide-y divide-gray-200">
            {filteredTransactions.map((transaction) => {
              const child = users.find(u => u.role === 'child')
              const rewardType = rewardTypes[0] // Mock for now
              
              return (
                <div key={transaction.id} className="px-6 py-4 hover:bg-gray-50">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-4">
                      <div className="flex items-center justify-center w-10 h-10 bg-gray-100 rounded-full">
                        {getTypeIcon(transaction.type)}
                      </div>
                      <div>
                        <div className="flex items-center space-x-2">
                          <span className="font-medium text-gray-900">
                            {getTypeLabel(transaction.type)}
                          </span>
                          <span className="text-lg font-bold">
                            {formatValue(transaction.value, 'money')}
                          </span>
                        </div>
                        <div className="flex items-center space-x-4 text-sm text-gray-600 mt-1">
                          <div className="flex items-center space-x-1">
                            <User className="h-3 w-3" />
                            <span>{child?.display_name || '未知'}</span>
                          </div>
                          <div className="flex items-center space-x-1">
                            <span>{rewardType?.name || '未知类型'}</span>
                          </div>
                          {transaction.note && (
                            <div className="text-gray-500">
                              备注: {transaction.note}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                    
                    <div className="text-right">
                      <div className="text-sm text-gray-600 flex items-center">
                        <Calendar className="h-3 w-3 mr-1" />
                        {formatDate(transaction.created_at)}
                      </div>
                      <div className="text-xs text-gray-500 mt-1">
                        ID: {transaction.id}
                      </div>
                    </div>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}