import { create } from 'zustand'
import axios from 'axios'

interface Family {
  id: number
  name: string
  created_at: string
}

interface User {
  id: number
  family_id: number
  role: 'guardian' | 'child'
  display_name: string
  wechat_openid?: string
  is_active: boolean
}

interface RewardType {
  id: number
  family_id: number
  name: string
  unit_kind: 'money' | 'time' | 'points' | 'custom'
  unit_label?: string
}

interface Balance {
  balance: number
}

interface FamilyStore {
  families: Family[]
  users: User[]
  rewardTypes: RewardType[]
  balances: Record<string, Balance>
  currentFamily: Family | null
  loading: boolean
  error: string | null
  fetchFamilies: () => Promise<void>
  fetchUsers: (familyId: number) => Promise<void>
  fetchRewardTypes: (familyId: number) => Promise<void>
  fetchBalance: (childId: number, rewardTypeId: number) => Promise<void>
  createFamily: (name: string) => Promise<void>
  createRewardType: (data: any) => Promise<void>
  addChild: (familyId: number, displayName: string) => Promise<void>
  deleteChild: (userId: number) => Promise<void>
  grantReward: (familyId: number, childId: number, rewardTypeId: number, value: number, note?: string) => Promise<{transaction_id:number,new_balance:number} | null>
  spendReward: (familyId: number, childId: number, rewardTypeId: number, value: number, note?: string) => Promise<{transaction_id:number,new_balance:number} | null>
}

const API_BASE = '/api/v1'

export const api = axios.create({
  baseURL: API_BASE,
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
    'Content-Type': 'application/json',
  },
})

export const useFamilyStore = create<FamilyStore>((set, get) => ({
  families: [],
  users: [],
  rewardTypes: [],
  balances: {},
  currentFamily: null,
  loading: false,
  error: null,

  fetchFamilies: async () => {
    set({ loading: true, error: null })
    try {
      const res = await api.get('/families')
      const families: Family[] = res.data.data || []
      let current = families[0] || null
      if (families.length > 0) {
        // Prefer a demo family with data
        const preferred = families.find(f => /张|李/.test(f.name)) || families[0]
        current = preferred
      }
      set({ families, currentFamily: current })
    } catch (error) {
      set({ error: '获取家庭信息失败' })
    } finally {
      set({ loading: false })
    }
  },

  fetchUsers: async (familyId: number) => {
    set({ loading: true, error: null })
    try {
      const res = await api.get('/users', { params: { family_id: familyId } })
      const users: User[] = res.data.data || []
      set({ users })
    } catch (error) {
      set({ error: '获取用户信息失败' })
    } finally {
      set({ loading: false })
    }
  },

  fetchRewardTypes: async (familyId: number) => {
    set({ loading: true, error: null })
    try {
      const res = await api.get('/reward_types', { params: { family_id: familyId } })
      const rewardTypes: RewardType[] = res.data.data || []
      set({ rewardTypes })
    } catch (error) {
      set({ error: '获取奖励类型失败' })
    } finally {
      set({ loading: false })
    }
  },

  fetchBalance: async (childId: number, rewardTypeId: number) => {
    try {
      const state = get()
      const familyId = state.currentFamily?.id
      if (!familyId) return
      const res = await api.get('/balances', { params: { family_id: familyId, child_id: childId, reward_type_id: rewardTypeId } })
      const balance = res.data.data || { balance: 0 }
      set(state => ({
        balances: {
          ...state.balances,
          [`${childId}-${rewardTypeId}`]: balance,
        },
      }))
    } catch (error) {
      console.error('获取余额失败:', error)
    }
  },

  createFamily: async (name: string) => {
    set({ loading: true, error: null })
    try {
      const res = await api.post('/families', { name })
      const created: Family = res.data.data
      set(state => ({
        families: [...state.families, created],
        currentFamily: created,
      }))
    } catch (error) {
      set({ error: '创建家庭失败' })
    } finally {
      set({ loading: false })
    }
  },

  createRewardType: async (data: any) => {
    set({ loading: true, error: null })
    try {
      const res = await api.post('/reward_types', data)
      const created: RewardType = res.data.data
      set(state => ({
        rewardTypes: [...state.rewardTypes, created],
      }))
    } catch (error) {
      set({ error: '创建奖励类型失败' })
    } finally {
      set({ loading: false })
    }
  },

  addChild: async (familyId: number, displayName: string) => {
    set({ loading: true, error: null })
    try {
      const res = await api.post('/users', { family_id: familyId, role: 'child', display_name: displayName })
      const user: User = res.data.data
      set(state => ({ users: [...state.users, user] }))
    } catch (error) {
      set({ error: '添加孩子失败' })
    } finally {
      set({ loading: false })
    }
  },

  deleteChild: async (userId: number) => {
    set({ loading: true, error: null })
    try {
      await api.delete(`/users/${userId}`)
      set(state => ({ users: state.users.filter(u => u.id !== userId) }))
    } catch (error) {
      set({ error: '删除孩子失败' })
    } finally {
      set({ loading: false })
    }
  },

  grantReward: async (familyId, childId, rewardTypeId, value, note) => {
    try {
      const res = await api.post('/rewards/grant', { family_id: familyId, child_id: childId, reward_type_id: rewardTypeId, value, note, idempotency_key: `grant-${childId}-${rewardTypeId}-${Date.now()}` })
      return res.data.data
    } catch (error) {
      set({ error: '授予奖励失败' })
      return null
    }
  },

  spendReward: async (familyId, childId, rewardTypeId, value, note) => {
    try {
      const res = await api.post('/rewards/spend', { family_id: familyId, child_id: childId, reward_type_id: rewardTypeId, value, note, idempotency_key: `spend-${childId}-${rewardTypeId}-${Date.now()}` })
      return res.data.data
    } catch (error) {
      set({ error: '消费奖励失败' })
      return null
    }
  },
}))