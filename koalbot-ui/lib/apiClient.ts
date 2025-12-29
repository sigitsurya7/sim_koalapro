import axios, { AxiosRequestConfig, Method, AxiosError } from 'axios'
import toast from 'react-hot-toast'

type RequestOptions = {
  token?: string
  params?: Record<string, any>
  data?: any
  baseURL?: string
}

const getBrowserToken = () => {
  if (typeof document === 'undefined') return undefined
  const match = document.cookie.match(/(?:^|; )token=([^;]+)/)
  return match ? decodeURIComponent(match[1]) : undefined
}

export const API_BASE = process.env.NEXT_PUBLIC_API_URL || ''

const buildHeaders = (token?: string) => ({
  'Content-Type': 'application/json',
  ...(token ? { Authorization: `Bearer ${token}` } : {}),
})

const request = async <T>(method: Method, url: string, options: RequestOptions = {}) => {
  const { token, params, data, baseURL } = options
  const bearer = token ?? getBrowserToken()

  const config: AxiosRequestConfig = {
    method,
    url,
    baseURL: baseURL ?? process.env.NEXT_PUBLIC_API_URL,
    headers: buildHeaders(bearer),
    params,
    data,
  }

  try {
    const response = await axios.request<T>(config)
    return response.data
  } catch (err) {
    const error = err as AxiosError
    const status = error.response?.status
    if (status === 401 && typeof window !== 'undefined') {
      toast.error('Sesi berakhir, silakan login kembali')
      document.cookie = 'token=; path=/; max-age=0'
      localStorage.removeItem('pm-user')
      localStorage.removeItem('pm_user')
      window.location.href = '/login'
      return Promise.reject(error)
    }
    throw error
  }
}

export const apiClient = {
  get: <T>(url: string, options?: Omit<RequestOptions, 'data'>) => request<T>('GET', url, options),
  post: <T>(url: string, data?: any, options?: Omit<RequestOptions, 'data'>) => request<T>('POST', url, { ...options, data }),
  put: <T>(url: string, data?: any, options?: Omit<RequestOptions, 'data'>) => request<T>('PUT', url, { ...options, data }),
  delete: <T>(url: string, options?: Omit<RequestOptions, 'data'>) => request<T>('DELETE', url, options),
}

export type ApiClient = typeof apiClient
