import { useEffect } from 'react'
import { useFamilyStore } from '../stores/familyStore'
import { useAuthStore } from '../stores/authStore'
import { Plus, TrendingUp, TrendingDown, DollarSign, Clock, Star } from 'lucide-react'
import { toast } from 'sonner'

export default function Dashboard() {
  const { 
    families, 
    users, 
    rewardTypes, 
    balances, 
    currentFamily,
    loading,
    fetchFamilies, 
    fetchUsers, 
    fetchRewardTypes,
    fetchBalance 
  } = useFamilyStore()
  
  const { isAuthenticated } = useAuthStore()

  useEffect(() => {
    if (isAuthenticated) {
      fetchFamilies()
    }
  }, [isAuthenticated, fetchFamilies])

  useEffect(() => {
    if (currentFamily) {
      fetchUsers(currentFamily.id)
      fetchRewardTypes(currentFamily.id)
    }
  }, [currentFamily, fetchUsers, fetchRewardTypes])

  useEffect(() => {
    if (users.length > 0 && rewardTypes.length > 0) {
      users.forEach(user => {
        if (user.role === 'child') {
          rewardTypes.forEach(type => {
            fetchBalance(user.id, type.id)
          })
        }
      })
    }
  }, [users, rewardTypes, fetchBalance])

  const getUnitIcon = (unitKind: string) => {
    switch (unitKind) {
      case 'money':
        return <DollarSign className="h-4 w-4" />
      case 'time':
        return <Clock className="h-4 w-4" />
      case 'points':
        return <Star className="h-4 w-4" />
      default:
        return <Star className="h-4 w-4" />
    }
  }

  const formatBalance = (balance: number, unitKind: string) => {
    if (unitKind === 'money') {
      return `¥${(balance / 100).toFixed(2)}`
    }
    return balance.toString()
  }

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-3xl font-bold text-gray-900 mb-4">请先登录</h1>
          <p className="text-gray-600">使用您的API Token进行身份验证</p>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">仪表板</h1>
          <p className="text-gray-600 mt-1">
            {currentFamily ? `当前家庭: ${currentFamily.name}` : '欢迎使用奖励系统'}
          </p>
        </div>
        <button
          onClick={() => toast.success('功能开发中...')}
          className="flex items-center space-x-2 bg-primary-500 text-white px-4 py-2 rounded-lg hover:bg-primary-600 transition-colors"
        >
          <Plus className="h-4 w-4" />
          <span>新增奖励</span>
        </button>
      </div>

      {/* Family Info */}
      {currentFamily && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">家庭信息</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="text-center">
              <div className="text-2xl font-bold text-primary-600">{users.length}</div>
              <div className="text-sm text-gray-600">家庭成员</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-success-600">{rewardTypes.length}</div>
              <div className="text-sm text-gray-600">奖励类型</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-warning-600">
                {users.filter(u => u.role === 'child').length}
              </div>
              <div className="text-sm text-gray-600">孩子数量</div>
            </div>
          </div>
        </div>
      )}

      {/* Children Balances */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {users.filter(user => user.role === 'child').map(child => (
          <div key={child.id} className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
              <span>{child.display_name}</span>
              <span className="ml-2 px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full">孩子</span>
            </h3>
            
            <div className="space-y-3">
              {rewardTypes.map(type => {
                const balance = balances[`${child.id}-${type.id}`]
                const balanceValue = balance?.balance || 0
                
                return (
                  <div key={type.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <div className="flex items-center space-x-3">
                      {getUnitIcon(type.unit_kind)}
                      <div>
                        <div className="font-medium text-gray-900">{type.name}</div>
                        <div className="text-sm text-gray-600">
                          {type.unit_kind === 'money' ? '元' : type.unit_label || '单位'}
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="text-lg font-bold text-gray-900">
                        {formatBalance(balanceValue, type.unit_kind)}
                      </div>
                      <div className="text-sm text-gray-600">可用余额</div>
                    </div>
                  </div>
                )
              })}
            </div>

            <div className="mt-4 flex space-x-2">
              <button
                onClick={() => toast.success('授予奖励功能开发中...')}
                className="flex items-center space-x-1 bg-success-500 text-white px-3 py-2 rounded text-sm hover:bg-success-600 transition-colors"
              >
                <TrendingUp className="h-3 w-3" />
                <span>授予</span>
              </button>
              <button
                onClick={() => toast.success('消费功能开发中...')}
                className="flex items-center space-x-1 bg-warning-500 text-white px-3 py-2 rounded text-sm hover:bg-warning-600 transition-colors"
              >
                <TrendingDown className="h-3 w-3" />
                <span>消费</span>
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* Quick Actions */}
      <div className="bg-white rounded-lg shadow p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">快捷操作</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <button
            onClick={() => toast.success('功能开发中...')}
            className="flex flex-col items-center p-4 border border-gray-200 rounded-lg hover:border-primary-500 hover:bg-primary-50 transition-colors"
          >
            <TrendingUp className="h-6 w-6 text-primary-500 mb-2" />
            <span className="text-sm font-medium">授予奖励</span>
          </button>
          <button
            onClick={() => toast.success('功能开发中...')}
            className="flex flex-col items-center p-4 border border-gray-200 rounded-lg hover:border-warning-500 hover:bg-warning-50 transition-colors"
          >
            <TrendingDown className="h-6 w-6 text-warning-500 mb-2" />
            <span className="text-sm font-medium">消费奖励</span>
          </button>
          <button
            onClick={() => toast.success('功能开发中...')}
            className="flex flex-col items-center p-4 border border-gray-200 rounded-lg hover:border-success-500 hover:bg-success-50 transition-colors"
          >
            <Star className="h-6 w-6 text-success-500 mb-2" />
            <span className="text-sm font-medium">新增类型</span>
          </button>
          <button
            onClick={() => toast.success('功能开发中...')}
            className="flex flex-col items-center p-4 border border-gray-200 rounded-lg hover:border-gray-400 hover:bg-gray-50 transition-colors"
          >
            <History className="h-6 w-6 text-gray-500 mb-2" />
            <span className="text-sm font-medium">查看记录</span>
          </button>
        </div>
      </div>
    </div>
  )
}