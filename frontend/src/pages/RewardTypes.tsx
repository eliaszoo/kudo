import { useState, useEffect } from 'react'
import { useFamilyStore } from '../stores/familyStore'
import { Plus, Edit, Trash2, DollarSign, Clock, Star, Gift } from 'lucide-react'
import { toast } from 'sonner'

export default function RewardTypes() {
  const { rewardTypes, loading, fetchRewardTypes, createRewardType, currentFamily } = useFamilyStore()
  const [editTarget, setEditTarget] = useState<any|null>(null)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [formData, setFormData] = useState({
    name: '',
    unit_kind: 'money' as 'money' | 'time' | 'points' | 'custom',
    unit_label: '',
  })

  useEffect(() => {
    if (currentFamily?.id) {
      fetchRewardTypes(currentFamily.id)
    }
  }, [currentFamily, fetchRewardTypes])

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

  const getUnitLabel = (unitKind: string, customLabel?: string) => {
    switch (unitKind) {
      case 'money':
        return '元'
      case 'time':
        return '分钟'
      case 'points':
        return '积分'
      default:
        return customLabel || '单位'
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!formData.name.trim()) {
      toast.error('请输入奖励类型名称')
      return
    }

    try {
      await createRewardType({
        family_id: currentFamily?.id,
        name: formData.name,
        unit_kind: formData.unit_kind,
        unit_label: formData.unit_label,
      })
      
      toast.success('奖励类型创建成功')
      setShowCreateModal(false)
      setFormData({ name: '', unit_kind: 'money', unit_label: '' })
    } catch (error) {
      toast.error('创建奖励类型失败')
    }
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">奖励类型管理</h1>
          <p className="text-gray-600 mt-1">管理您的家庭奖励类型</p>
        </div>
        <button
          onClick={() => setShowCreateModal(true)}
          className="flex items-center space-x-2 bg-primary-500 text-white px-4 py-2 rounded-lg hover:bg-primary-600 transition-colors"
        >
          <Plus className="h-4 w-4" />
          <span>新增类型</span>
        </button>
      </div>

      {/* Reward Types List */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <h2 className="text-lg font-semibold text-gray-900">奖励类型列表</h2>
        </div>
        
        {loading ? (
          <div className="p-8 text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500 mx-auto"></div>
            <p className="mt-2 text-gray-600">加载中...</p>
          </div>
        ) : rewardTypes.length === 0 ? (
          <div className="p-8 text-center text-gray-500">
            <Gift className="h-12 w-12 mx-auto mb-4 text-gray-300" />
            <p>暂无奖励类型</p>
            <p className="text-sm mt-1">点击上方按钮创建第一个奖励类型</p>
          </div>
        ) : (
          <div className="divide-y divide-gray-200">
            {rewardTypes.map((type) => (
              <div key={type.id} className="px-6 py-4 flex items-center justify-between hover:bg-gray-50">
                <div className="flex items-center space-x-4">
                  <div className="flex items-center justify-center w-10 h-10 bg-primary-100 rounded-full">
                    {getUnitIcon(type.unit_kind)}
                  </div>
                  <div>
                    <h3 className="text-lg font-medium text-gray-900">{type.name}</h3>
                    <p className="text-sm text-gray-600">
                      类型: {type.unit_kind} | 单位: {getUnitLabel(type.unit_kind, type.unit_label)}
                    </p>
                  </div>
                </div>
                
                <div className="flex items-center space-x-2">
                  <button
                    onClick={() => openEdit(type)}
                    className="p-2 text-gray-400 hover:text-gray-600 transition-colors"
                  >
                    <Edit className="h-4 w-4" />
                  </button>
                  <button
                    onClick={() => toast.success('删除功能开发中...')}
                    className="p-2 text-gray-400 hover:text-red-600 transition-colors"
                  >
                    <Trash2 className="h-4 w-4" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Create Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md mx-4">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">创建奖励类型</h2>
            
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  类型名称
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                  placeholder="例如: 零花钱、看电视时间"
                  required
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  单位类型
                </label>
                <select
                  value={formData.unit_kind}
                  onChange={(e) => setFormData({ ...formData, unit_kind: e.target.value as any })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                >
                  <option value="money">货币（元）</option>
                  <option value="time">时间（分钟）</option>
                  <option value="points">积分</option>
                  <option value="custom">自定义</option>
                </select>
              </div>
              
              {formData.unit_kind === 'custom' && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    自定义单位
                  </label>
                  <input
                    type="text"
                    value={formData.unit_label}
                    onChange={(e) => setFormData({ ...formData, unit_label: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                    placeholder="例如: 星星、贴纸"
                  />
                </div>
              )}
              
              <div className="flex space-x-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowCreateModal(false)}
                  className="flex-1 px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 transition-colors"
                >
                  取消
                </button>
                <button
                  type="submit"
                  className="flex-1 px-4 py-2 bg-primary-500 text-white rounded-md hover:bg-primary-600 transition-colors"
                >
                  创建
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
  const [editForm, setEditForm] = useState({ name:'', unit_kind:'money' as 'money'|'time'|'points'|'custom', unit_label:'' })

  const openEdit = (type:any) => {
    setEditTarget(type)
    setEditForm({ name:type.name, unit_kind:type.unit_kind, unit_label:type.unit_label || '' })
  }

  const updateRewardType = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!editTarget) return
    try {
      const res = await api.patch(`/reward_types/${editTarget.id}`, { name: editForm.name, unit_kind: editForm.unit_kind, unit_label: editForm.unit_label })
      const updated = res.data.data
      const newList = rewardTypes.map(rt => rt.id === updated.id ? updated : rt)
      ;(useFamilyStore as any).setState({ rewardTypes: newList })
      toast.success('更新成功')
      setEditTarget(null)
    } catch (err) {
      toast.error('更新失败')
    }
  }
      {editTarget && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md mx-4">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">编辑奖励类型</h2>
            <form onSubmit={updateRewardType} className="space-y-4">
              <input className="w-full px-3 py-2 border border-gray-300 rounded-md" value={editForm.name} onChange={(e)=>setEditForm({...editForm, name:e.target.value})} />
              <select className="w-full px-3 py-2 border border-gray-300 rounded-md" value={editForm.unit_kind} onChange={(e)=>setEditForm({...editForm, unit_kind: e.target.value as any})}>
                <option value="money">货币（元）</option>
                <option value="time">时间（分钟）</option>
                <option value="points">积分</option>
                <option value="custom">自定义</option>
              </select>
              {editForm.unit_kind === 'custom' && (
                <input className="w-full px-3 py-2 border border-gray-300 rounded-md" placeholder="自定义单位" value={editForm.unit_label} onChange={(e)=>setEditForm({...editForm, unit_label: e.target.value})} />
              )}
              <div className="flex space-x-3">
                <button type="button" className="flex-1 px-4 py-2 border border-gray-300 rounded-md" onClick={()=>setEditTarget(null)}>取消</button>
                <button type="submit" className="flex-1 px-4 py-2 bg-primary-500 text-white rounded-md">保存</button>
              </div>
            </form>
          </div>
        </div>
      )}