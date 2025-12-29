'use client'

import { useEffect, useRef, useState } from 'react'
import {
  Button,
  Card,
  CardBody,
  CardHeader,
  Input,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  Switch,
} from '@heroui/react'
import { FiTrash2 } from 'react-icons/fi'
import { DataTable, DataTableColumn } from '@/components/data-table'
import { apiClient } from '@/lib/apiClient'

const JENIS = 'stockity'

type MasterPengguna = {
  id: number
  uuid: string
  id_pengguna: number
  telegram?: string | null
  jenis: string
  active: boolean
  created_at: string
  updated_at?: string | null
}

type Pagination = {
  limit: number
  page: number
  pages: number
  search: string
  total: number
}

type ListResponse = {
  data: MasterPengguna[]
  pagination: Pagination
}

export default function MemberStockityPage() {
  const [data, setData] = useState<MasterPengguna[]>([])
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [deletingUser, setDeletingUser] = useState<MasterPengguna | null>(null)
  const [isOpen, setIsOpen] = useState(false)
  const [form, setForm] = useState({ id_pengguna: '', telegram: '', active: false })
  const [updatingIds, setUpdatingIds] = useState<Set<number>>(new Set())
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 20,
    pages: 1,
    total: 0,
  })
  const debounceRef = useRef<NodeJS.Timeout | null>(null)

  const loadData = async (page = pagination.page, limit = pagination.limit, query = search) => {
    setLoading(true)
    try {
      const res = await apiClient.get<ListResponse>('master-pengguna', {
        params: {
          page,
          limit,
          search: query,
          jenis: JENIS,
        },
      })
      setData(res.data || [])
      setPagination({
        page: res.pagination.page,
        limit: res.pagination.limit,
        pages: res.pagination.pages,
        total: res.pagination.total,
      })
    } catch (err) {
      console.error('Gagal memuat data pengguna', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const handleSearchChange = (value: string) => {
    setSearch(value)
    if (debounceRef.current) clearTimeout(debounceRef.current)
    debounceRef.current = setTimeout(() => {
      loadData(1, pagination.limit, value)
    }, 300)
  }

  const handlePageChange = (page: number) => {
    loadData(page, pagination.limit, search)
  }

  const handleLimitChange = (limit: number) => {
    loadData(1, limit, search)
  }

  const handleToggleActive = async (item: MasterPengguna, nextValue: boolean) => {
    setUpdatingIds((prev) => {
      const next = new Set(prev)
      next.add(item.id)
      return next
    })
    setData((prev) =>
      prev.map((row) => (row.id === item.id ? { ...row, active: nextValue } : row))
    )
    try {
      await apiClient.put(`master-pengguna/${item.id}`, { active: nextValue })
    } catch (err) {
      setData((prev) =>
        prev.map((row) => (row.id === item.id ? { ...row, active: item.active } : row))
      )
      console.error('Gagal mengubah status', err)
    } finally {
      setUpdatingIds((prev) => {
        const next = new Set(prev)
        next.delete(item.id)
        return next
      })
    }
  }

  const handleCreate = async () => {
    const id = Number(form.id_pengguna)
    if (!id) {
      return
    }
    setSaving(true)
    try {
      await apiClient.post('master-pengguna', {
        id_pengguna: id,
        telegram: form.telegram || undefined,
        jenis: JENIS,
        active: form.active,
      })
      setIsOpen(false)
      setForm({ id_pengguna: '', telegram: '', active: false })
      loadData(1, pagination.limit, search)
    } catch (err) {
      console.error('Gagal menambah pengguna', err)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (item: MasterPengguna) => {
    setUpdatingIds((prev) => {
      const next = new Set(prev)
      next.add(item.id)
      return next
    })
    try {
      await apiClient.delete(`master-pengguna/${item.id}`)
      setData((prev) => prev.filter((row) => row.id !== item.id))
    } catch (err) {
      console.error('Gagal menghapus pengguna', err)
    } finally {
      setUpdatingIds((prev) => {
        const next = new Set(prev)
        next.delete(item.id)
        setDeletingUser(null)
        return next
      })
    }
  }

  const columns: DataTableColumn<MasterPengguna>[] = [
    { key: 'id', label: 'ID' },
    { key: 'id_pengguna', label: 'ID Pengguna' },
    { key: 'telegram', label: 'Telegram', render: (item) => item.telegram ?? '-' },
    {
      key: 'active',
      label: 'Status',
      render: (item) => (
        <Switch
          size="sm"
          isSelected={item.active}
          isDisabled={updatingIds.has(item.id)}
          onValueChange={(value) => handleToggleActive(item, value)}
        >
          {item.active ? 'Active' : 'Inactive'}
        </Switch>
      ),
    },
    {
      key: 'created_at',
      label: 'Dibuat pada',
      render: (item) => new Date(item.created_at).toLocaleString('id-ID'),
    },
    {
      key: 'actions',
      label: 'Aksi',
      render: (item) => (
        <Button
          startContent={<FiTrash2 />}
          size="sm"
          variant="flat"
          color="danger"
          isDisabled={updatingIds.has(item.id)}
          onPress={() => setDeletingUser(item)}
        >
          Hapus
        </Button>
      ),
    }
  ]

  return (
    <section className="space-y-4">
      <div>
        <p className="text-sm text-default-500">Data Pengguna</p>
        <h1 className="text-2xl font-bold">Member Stockity</h1>
      </div>

      <DataTable
        columns={columns}
        data={data}
        pagination={{
          page: pagination.page,
          pages: pagination.pages,
          limit: pagination.limit,
          total: pagination.total,
        }}
        searchValue={search}
        onSearchChange={handleSearchChange}
        onPageChange={handlePageChange}
        onLimitChange={handleLimitChange}
        loading={loading}
        getRowId={(item) => item.id}
        topRightSlot={
          <Button color="primary" size="sm" onPress={() => setIsOpen(true)}>
            Tambah Pengguna
          </Button>
        }
      />

      <Modal isOpen={isOpen} onOpenChange={() => setIsOpen(false)} isDismissable={!saving}>
        <ModalContent>
          {() => (
            <>
              <ModalHeader className="flex flex-col gap-1">Tambah Pengguna</ModalHeader>
              <ModalBody className="space-y-4">
                <Input
                  label="ID Pengguna"
                  placeholder="Masukkan ID Pengguna"
                  type="number"
                  value={form.id_pengguna}
                  onValueChange={(value) => setForm((prev) => ({ ...prev, id_pengguna: value }))}
                />
                <Input
                  label="Telegram"
                  placeholder="@username"
                  value={form.telegram}
                  onValueChange={(value) => setForm((prev) => ({ ...prev, telegram: value }))}
                />
                <Switch
                  isSelected={form.active}
                  onValueChange={(value) => setForm((prev) => ({ ...prev, active: value }))}
                >
                  Status aktif
                </Switch>
              </ModalBody>
              <ModalFooter>
                <Button variant="flat" onPress={() => setIsOpen(false)} isDisabled={saving}>
                  Batal
                </Button>
                <Button color="primary" onPress={handleCreate} isLoading={saving}>
                  Simpan
                </Button>
              </ModalFooter>
            </>
          )}
        </ModalContent>
      </Modal>

      <Modal isOpen={!!deletingUser} onOpenChange={() => setDeletingUser(null)} isDismissable={!saving}>
        <ModalContent>
          {() => (
            <>
              <ModalHeader className="flex flex-col gap-1">Konfirmasi Hapus</ModalHeader>
              <ModalBody>
                <p>Apakah anda yakin ingin menghapus user ?</p>
                {deletingUser && (
                  <p className="text-sm text-default-500">ID Pengguna: {deletingUser.id_pengguna}</p>
                )}
              </ModalBody>
              <ModalFooter>
                <Button variant="flat" onPress={() => setDeletingUser(null)}>
                  Batal
                </Button>
                <Button
                  color="danger"
                  isLoading={deletingUser ? updatingIds.has(deletingUser.id) : false}
                  onPress={() => deletingUser && handleDelete(deletingUser)}
                >
                  Hapus
                </Button>
              </ModalFooter>
            </>
          )}
        </ModalContent>
      </Modal>
    </section>
  )
}
