import axios from 'axios'
import type { ApiEnvelope } from './types'

export const authTokenKey = 'windnav_token'

export const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem(authTokenKey)
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const message = error.response?.data?.error?.message ?? '请求失败'
    return Promise.reject(new Error(message))
  },
)

export async function getData<T>(url: string, params?: Record<string, unknown>): Promise<T> {
  const response = await api.get<ApiEnvelope<T>>(url, { params })
  return response.data.data
}

export async function postData<T>(url: string, body?: unknown): Promise<T> {
  const response = await api.post<ApiEnvelope<T>>(url, body)
  return response.data.data
}

export async function putData<T>(url: string, body?: unknown): Promise<T> {
  const response = await api.put<ApiEnvelope<T>>(url, body)
  return response.data.data
}

export async function deleteData<T>(url: string): Promise<T> {
  const response = await api.delete<ApiEnvelope<T>>(url)
  return response.data.data
}
