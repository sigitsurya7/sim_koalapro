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
  Select,
  SelectItem,
  Switch,
} from '@heroui/react'
import { FiEdit2 } from 'react-icons/fi'
import { DataTable, DataTableColumn } from '@/components/data-table'
import { apiClient } from '@/lib/apiClient'

const isAdmin = () => {
  if (typeof window === 'undefined') return false
  const raw = localStorage.getItem('pm_user') ?? localStorage.getItem('pm-user')
  if (!raw) return false
  try {
    const parsed = JSON.parse(raw)
    return (parsed.role ?? parsed.role_name) === 'admin'
  } catch {
    return false
  }
}

type UserItem = {
  uid: string
  username: string
  role: string
  active: boolean
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
  data: UserItem[]
  pagination: Pagination
}

export default function UsersPage() {
  const [data, setData] = useState<UserItem[]>([])
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)
  const [isOpen, setIsOpen] = useState(false)
  const [editingUser, setEditingUser] = useState<UserItem | null>(null)
  const [form, setForm] = useState({ username: '', password: '', role: 'viewer', active: true })
  const [updatingIds, setUpdatingIds] = useState<Set<string>>(new Set())
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 20,
    pages: 1,
    total: 0,
  })
  const debounceRef = useRef<NodeJS.Timeout | null>(null)
  const [allowed, setAllowed] = useState(true)

  const loadData = async (page = pagination.page, limit = pagination.limit, query = search) => {
    setLoading(true)
    try {
      const res = await apiClient.get<ListResponse>('users', {
        params: {
          page,
          limit,
          search: query,
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
      console.error('Gagal memuat data user', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    const ok = isAdmin()
    setAllowed(ok)
    if (ok) {
      loadData()
    }
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

  const handleToggleActive = async (item: UserItem, nextValue: boolean) => {
    setUpdatingIds((prev) => {
      const next = new Set(prev)
      next.add(item.uid)
      return next
    })
    setData((prev) =>
      prev.map((row) => (row.uid === item.uid ? { ...row, active: nextValue } : row))
    )
    try {
      await apiClient.put(`users/${item.uid}`, { active: nextValue })
    } catch (err) {
      setData((prev) =>
        prev.map((row) => (row.uid === item.uid ? { ...row, active: item.active } : row))
      )
      console.error('Gagal mengubah status', err)
    } finally {
      setUpdatingIds((prev) => {
        const next = new Set(prev)
        next.delete(item.uid)
        return next
      })
    }
  }

  const openCreate = () => {
    setEditingUser(null)
    setForm({ username: '', password: '', role: 'viewer', active: true })
    setIsOpen(true)
  }

  const openEdit = (item: UserItem) => {
    setEditingUser(item)
    setForm({ username: item.username, password: '', role: item.role, active: item.active })
    setIsOpen(true)
  }

  const handleSave = async () => {
    if (!form.username || !form.role) return
    setSaving(true)
    try {
      if (editingUser) {
        const payload: Record<string, any> = {
          username: form.username,
          role: form.role,
          active: form.active,
        }
        if (form.password) {
          payload.password = form.password
        }
        await apiClient.put(`users/${editingUser.uid}`, payload)
      } else {
        await apiClient.post('users', {
          username: form.username,
          password: form.password,
          role: form.role,
        })
      }
      setIsOpen(false)
      setEditingUser(null)
      setForm({ username: '', password: '', role: 'viewer', active: true })
      loadData(1, pagination.limit, search)
    } catch (err) {
      console.error('Gagal menyimpan user', err)
    } finally {
      setSaving(false)
    }
  }

  const columns: DataTableColumn<UserItem>[] = [
    { key: 'username', label: 'Username' },
    { key: 'role', label: 'Role' },
    {
      key: 'active',
      label: 'Status',
      render: (item) => (
        <Switch
          size="sm"
          isSelected={item.active}
          isDisabled={updatingIds.has(item.uid)}
          onValueChange={(value) => handleToggleActive(item, value)}
        >
          {item.active ? 'Active' : 'Inactive'}
        </Switch>
      ),
    },
    {
      key: 'actions',
      label: 'Aksi',
      render: (item) => (
        <Button
          size="sm"
          variant="flat"
          startContent={<FiEdit2 />}
          onPress={() => openEdit(item)}
        >
          Edit
        </Button>
      ),
    },
  ]

  if (!allowed) {
    return (
      <section className="space-y-4">
        <Card className="border border-default-100/60 shadow-none">
          <CardHeader>
            <h1 className="text-lg font-semibold">Akses ditolak</h1>
          </CardHeader>
          <CardBody>
            <p className="text-sm text-default-500">Halaman ini hanya untuk admin.</p>
          </CardBody>
        </Card>
      </section>
    )
  }

  return (
    <section className="space-y-4">
      <div>
        <p className="text-sm text-default-500">Manajemen User</p>
        <h1 className="text-2xl font-bold">Users</h1>
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
        getRowId={(item) => item.uid}
        topRightSlot={
          <Button color="primary" size="sm" onPress={openCreate}>
            Tambah User
          </Button>
        }
      />

      <Modal isOpen={isOpen} onOpenChange={() => setIsOpen(false)} isDismissable={!saving}>
        <ModalContent>
          {() => (
            <>
              <ModalHeader className="flex flex-col gap-1">
                {editingUser ? 'Edit User' : 'Tambah User'}
              </ModalHeader>
              <ModalBody className="space-y-4">
                <Input
                  label="Username"
                  placeholder="username"
                  value={form.username}
                  onValueChange={(value) => setForm((prev) => ({ ...prev, username: value }))}
                />
                <Input
                  label="Password"
                  placeholder={editingUser ? 'Kosongkan jika tidak diubah' : 'password'}
                  type="password"
                  value={form.password}
                  onValueChange={(value) => setForm((prev) => ({ ...prev, password: value }))}
                />
                <Select
                  label="Role"
                  selectedKeys={[form.role]}
                  onChange={(e) => setForm((prev) => ({ ...prev, role: e.target.value }))}
                >
                  <SelectItem key="admin">admin</SelectItem>
                  <SelectItem key="viewer">viewer</SelectItem>
                </Select>
                {editingUser && (
                  <Switch
                    isSelected={form.active}
                    onValueChange={(value) => setForm((prev) => ({ ...prev, active: value }))}
                  >
                    Status aktif
                  </Switch>
                )}
              </ModalBody>
              <ModalFooter>
                <Button variant="flat" onPress={() => setIsOpen(false)} isDisabled={saving}>
                  Batal
                </Button>
                <Button color="primary" onPress={handleSave} isLoading={saving}>
                  Simpan
                </Button>
              </ModalFooter>
            </>
          )}
        </ModalContent>
      </Modal>
    </section>
  )
}
