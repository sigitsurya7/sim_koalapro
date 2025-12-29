"use client"

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Button, Input } from "@heroui/react"
import { FiUser, FiLock, FiArrowRight } from 'react-icons/fi'
import { apiClient } from '@/lib/apiClient'
import toast from 'react-hot-toast'

export default function LoginPages() {
  const router = useRouter()
  const [formData, setFormData] = useState({
    username: '',
    password: ''
  })
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError(null)

    await toast.promise(
      (async () => {
        const res = await apiClient.post<{ token: string; message?: string; user?: any }>(
          '/login',
          { ...formData }
        )
        if (!res?.token) {
          throw new Error(res?.message || 'Token tidak ditemukan')
        }
        // Simpan token & user ke cookie/localStorage sebelum redirect
        document.cookie = `token=${res.token}; path=/; secure`
        localStorage.setItem('pm-user', JSON.stringify(res?.user ?? {}))
        const params = new URLSearchParams(window.location.search)
        const from = params.get('from')
        router.push(from || '/dashboard')
        return res
      })(),
      {
        loading: 'Memproses login...',
        success: 'Login berhasil',
        error: (err) => {
          const msg = err?.response?.data?.message || err?.message || 'Login gagal'
          setError(msg)
          return msg
        },
      }
    ).catch(() => {})

    setIsLoading(false)
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    })
  }

  return (
    <section className="flex justify-center items-center min-h-screen">
      
      {/* Login Card */}
      <div className="w-full max-w-sm">
        
        {/* Logo */}
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            Koala Pro
          </h1>
          <p className="text-gray-500 dark:text-gray-400 mt-1 text-sm">
            Masuk untuk melanjutkan
          </p>
        </div>

        {/* Login Form */}
        <form onSubmit={handleSubmit} className="space-y-4">
          {error && (
            <div className="rounded-lg border border-danger-200 bg-danger-50 px-3 py-2 text-sm text-danger-700">
              {error}
            </div>
          )}
          
          {/* Username */}
          <Input
            name="username"
            placeholder="Username"
            value={formData.username}
            onChange={handleChange}
            startContent={<FiUser className="text-gray-400" size={18} />}
            size="md"
            variant='bordered'
            radius="sm"
            required
            autoComplete="username"
          />

          {/* Password */}
          <Input
            name="password"
            type="password"
            placeholder="Password"
            value={formData.password}
            onChange={handleChange}
            startContent={<FiLock className="text-gray-400" size={18} />}
            size="md"
            variant='bordered'
            radius="sm"
            required
            autoComplete="current-password"
          />

          {/* Login Button */}
          <Button
            type="submit"
            isLoading={isLoading}
            fullWidth
            color='primary'
            endContent={!isLoading && <FiArrowRight size={18} />}
          >
            {isLoading ? 'Signing in...' : 'Sign In'}
          </Button>

        </form>
      </div>

    </section>
  )
}
