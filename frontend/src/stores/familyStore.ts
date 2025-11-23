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
}

const API_BASE = '/api/v1'

const api = axios.create({
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
      // Mock data for now
      const mockFamilies: Family[] = [
        { id: 1, name: '我的家庭', created_at: new Date().toISOString() },
      ]
      set({ families: mockFamilies, currentFamily: mockFamilies[0] || null })
    } catch (error) {
      set({ error: '获取家庭信息失败' })
    } finally {
      set({ loading: false })
    }
  },

  fetchUsers: async (familyId: number) => {
    set({ loading: true, error: null })
    try {
      // Mock data for now
      const mockUsers: User[] = [
        { id: 1, family_id: familyId, role: 'guardian', display_name: '爸爸', is_active: true },
        { id: 2, family_id: familyId, role: 'child', display_name: '小明', is_active: true },
      ]
      set({ users: mockUsers })
    } catch (error) {
      set({ error: '获取用户信息失败' })
    } finally {
      set({ loading: false })
    }
  },

  fetchRewardTypes: async (familyId: number) => {
    set({ loading: true, error: null })
    try {
      // Mock data for now
      const mockRewardTypes: RewardType[] = [
        { id: 1, family_id: familyId, name: '零花钱', unit_kind: 'money', unit_label: '元' },
        { id: 2, family_id: familyId, name: '看电视时间', unit_kind: 'time', unit_label: '分钟' },
      ]
      set({ rewardTypes: mockRewardTypes })
    } catch (error) {
      set({ error: '获取奖励类型失败' })
    } finally {
      set({ loading: false })
    }
  },

  fetchBalance: async (childId: number, rewardTypeId: number) => {
    try {
      // Mock balance for now
      const balance = { balance: Math.floor(Math.random() * 1000) }
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
      // Mock creation
      const newFamily: Family = {
        id: Date.now(),
        name,
        created_at: new Date().toISOString(),
      }
      set(state => ({
        families: [...state.families, newFamily],
        currentFamily: newFamily,
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
      // Mock creation
      const newRewardType: RewardType = {
        id: Date.now(),
        family_id: data.family_id,
        name: data.name,
        unit_kind: data.unit_kind,
        unit_label: data.unit_label,
      }
      set(state => ({
        rewardTypes: [...state.rewardTypes, newRewardType],
      }))
    } catch (error) {
      set({ error: '创建奖励类型失败' })
    } finally {
      set({ loading: false })
    }
  },
}))