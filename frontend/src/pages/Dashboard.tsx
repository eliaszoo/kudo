import { useEffect, useState } from 'react'
import { useFamilyStore } from '../stores/familyStore'
import { useAuthStore } from '../stores/authStore'
import { useNavigate } from 'react-router-dom'
import { Plus, TrendingUp, TrendingDown, DollarSign, Clock, Star, History, Trash2 } from 'lucide-react'
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
    fetchBalance,
    grantReward,
    spendReward,
    addChild,
    deleteChild,
  } = useFamilyStore()
  
  const { isAuthenticated } = useAuthStore()
  const navigate = useNavigate()

  const [showGrant, setShowGrant] = useState<{childId:number}|null>(null)
  const [showSpend, setShowSpend] = useState<{childId:number}|null>(null)
  const [form, setForm] = useState({ rewardTypeId: 0, value: 0, note: '' })
  const [showAddChild, setShowAddChild] = useState(false)
  const [childName, setChildName] = useState('')

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

  const submitGrant = async () => {
    if (!showGrant || !currentFamily) return
    const res = await grantReward(currentFamily.id, showGrant.childId, form.rewardTypeId || rewardTypes[0]?.id, form.value, form.note)
    if (res) {
      toast.success('授予成功')
      const rtId = form.rewardTypeId || rewardTypes[0]?.id
      if (rtId) fetchBalance(showGrant.childId, rtId)
      setShowGrant(null)
      setForm({ rewardTypeId: 0, value: 0, note: '' })
    } else {
      toast.error('授予失败')
    }
  }

  const submitSpend = async () => {
    if (!showSpend || !currentFamily) return
    const res = await spendReward(currentFamily.id, showSpend.childId, form.rewardTypeId || rewardTypes[0]?.id, form.value, form.note)
    if (res) {
      toast.success('消费成功')
      const rtId = form.rewardTypeId || rewardTypes[0]?.id
      if (rtId) fetchBalance(showSpend.childId, rtId)
      setShowSpend(null)
      setForm({ rewardTypeId: 0, value: 0, note: '' })
    } else {
      toast.error('消费失败')
    }
  }

  const submitAddChild = async () => {
    if (!currentFamily || !childName.trim()) {
      toast.error('请输入孩子名')
      return
    }
    await addChild(currentFamily.id, childName.trim())
    await fetchUsers(currentFamily.id)
    toast.success('添加孩子成功')
    setChildName('')
    setShowAddChild(false)
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
        <div className="flex items-center space-x-2">
          <button
            onClick={() => setShowAddChild(true)}
            className="flex items-center space-x-2 bg-primary-500 text-white px-4 py-2 rounded-lg hover:bg-primary-600 transition-colors"
          >
            <Plus className="h-4 w-4" />
            <span>添加孩子</span>
          </button>
        </div>
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
                onClick={() => setShowGrant({ childId: child.id })}
                className="flex items-center space-x-1 bg-success-500 text-white px-3 py-2 rounded text-sm hover:bg-success-600 transition-colors"
              >
                <TrendingUp className="h-3 w-3" />
                <span>授予</span>
              </button>
              <button
                onClick={() => setShowSpend({ childId: child.id })}
                className="flex items-center space-x-1 bg-warning-500 text-white px-3 py-2 rounded text-sm hover:bg-warning-600 transition-colors"
              >
                <TrendingDown className="h-3 w-3" />
                <span>消费</span>
              </button>
              <button
                onClick={async () => {
                  if (confirm('确认删除该孩子？这将清空其账户与交易。')) {
                    await deleteChild(child.id)
                    if (currentFamily) fetchUsers(currentFamily.id)
                    toast.success('已删除孩子')
                  }
                }}
                className="flex items-center space-x-1 bg-red-500 text-white px-3 py-2 rounded text-sm hover:bg-red-600 transition-colors"
              >
                <Trash2 className="h-3 w-3" />
                <span>删除</span>
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
            onClick={() => setShowGrant(users.find(u=>u.role==='child') ? { childId: users.find(u=>u.role==='child')!.id } : null)}
            className="flex flex-col items-center p-4 border border-gray-200 rounded-lg hover:border-primary-500 hover:bg-primary-50 transition-colors"
          >
            <TrendingUp className="h-6 w-6 text-primary-500 mb-2" />
            <span className="text-sm font-medium">授予奖励</span>
          </button>
          <button
            onClick={() => setShowSpend(users.find(u=>u.role==='child') ? { childId: users.find(u=>u.role==='child')!.id } : null)}
            className="flex flex-col items-center p-4 border border-gray-200 rounded-lg hover:border-warning-500 hover:bg-warning-50 transition-colors"
          >
            <TrendingDown className="h-6 w-6 text-warning-500 mb-2" />
            <span className="text-sm font-medium">消费奖励</span>
          </button>
          <button
            onClick={() => setShowAddChild(true)}
            className="flex flex-col items-center p-4 border border-gray-200 rounded-lg hover:border-success-500 hover:bg-success-50 transition-colors"
          >
            <Star className="h-6 w-6 text-success-500 mb-2" />
            <span className="text-sm font-medium">添加孩子</span>
          </button>
          <button
            onClick={() => {
              const firstChild = users.find(u=>u.role==='child')
              if (firstChild) navigate(`/transactions?child_id=${firstChild.id}`)
            }}
            className="flex flex-col items-center p-4 border border-gray-200 rounded-lg hover:border-gray-400 hover:bg-gray-50 transition-colors"
          >
            <History className="h-6 w-6 text-gray-500 mb-2" />
            <span className="text-sm font-medium">查看记录</span>
          </button>
        </div>
      </div>

      {showAddChild && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md mx-4">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">添加孩子</h2>
            <div className="space-y-4">
              <input className="w-full px-3 py-2 border border-gray-300 rounded-md" placeholder="孩子称呼，如：小明" value={childName} onChange={(e)=>setChildName(e.target.value)} />
              <div className="flex space-x-3">
                <button className="flex-1 px-4 py-2 border border-gray-300 rounded-md" onClick={()=>setShowAddChild(false)}>取消</button>
                <button className="flex-1 px-4 py-2 bg-primary-500 text-white rounded-md" onClick={submitAddChild}>添加</button>
              </div>
            </div>
          </div>
        </div>
      )}

      {(showGrant || showSpend) && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md mx-4">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">{showGrant ? '授予奖励' : '消费奖励'}</h2>
            <div className="space-y-4">
              <select className="w-full px-3 py-2 border border-gray-300 rounded-md" value={form.rewardTypeId} onChange={(e)=>setForm({...form, rewardTypeId: Number(e.target.value)})}>
                <option value={0}>选择奖励类型</option>
                {rewardTypes.map(rt => (
                  <option key={rt.id} value={rt.id}>{rt.name}</option>
                ))}
              </select>
              <input type="number" className="w-full px-3 py-2 border border-gray-300 rounded-md" placeholder="数量（单位与类型相关）" value={form.value} onChange={(e)=>setForm({...form, value: Number(e.target.value)})} />
              <input className="w-full px-3 py-2 border border-gray-300 rounded-md" placeholder="备注（可选）" value={form.note} onChange={(e)=>setForm({...form, note: e.target.value})} />
              <div className="flex space-x-3">
                <button className="flex-1 px-4 py-2 border border-gray-300 rounded-md" onClick={()=>{setShowGrant(null);setShowSpend(null)}}>取消</button>
                <button className="flex-1 px-4 py-2 bg-primary-500 text-white rounded-md" onClick={showGrant ? submitGrant : submitSpend}>{showGrant ? '授予' : '消费'}</button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}