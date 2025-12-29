'use client'

import { useEffect, useState } from 'react'
import { Button, Card, CardBody, CardHeader, Skeleton } from '@heroui/react'
import { FiRefreshCw, FiUserCheck, FiUserX, FiUsers } from 'react-icons/fi'
import { apiClient } from '@/lib/apiClient'

type SummaryResponse = {
  total: number
  active: number
  inactive: number
}

export default function HomePage() {
  const [summary, setSummary] = useState<SummaryResponse | null>(null)
  const [loading, setLoading] = useState(true)

  const loadData = async () => {
    setLoading(true)
    try {
      const res = await apiClient.get<SummaryResponse>('dashboard/summary')
      setSummary(res)
    } catch (err) {
      console.error('Gagal memuat dashboard', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [])

  const cards = [
    { label: 'Total pengguna', value: summary?.total ?? 0, icon: <FiUsers /> },
    { label: 'Aktif', value: summary?.active ?? 0, icon: <FiUserCheck /> },
    { label: 'Belum aktif', value: summary?.inactive ?? 0, icon: <FiUserX /> },
  ]

  return (
    <section className="space-y-6" title="Dashboard Overview">
      <div className="flex items-center justify-between gap-2">
        <div>
          <p className="text-sm text-default-500">Dashboard</p>
          <h1 className="text-2xl font-bold">Ringkasan pengguna</h1>
        </div>
        <Button size="sm" variant="flat" startContent={<FiRefreshCw />} onPress={loadData}>
          Refresh
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {cards.map((item) => (
          <Card key={item.label} shadow="md" className="border border-default-100/60 shadow-none">
            <CardHeader className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-default-100/60">{item.icon}</div>
              <div>
                <p className="text-sm text-default-500">{item.label}</p>
                {loading ? (
                  <Skeleton className="h-7 w-20 rounded-md" />
                ) : (
                  <p className="text-2xl font-semibold">{item.value}</p>
                )}
              </div>
            </CardHeader>
            <CardBody />
          </Card>
        ))}
      </div>
    </section>
  )
}
