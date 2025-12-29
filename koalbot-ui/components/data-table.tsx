"use client"

import { ReactNode } from 'react'
import { 
  Table,
  TableHeader,
  TableBody,
  TableColumn,
  TableRow,
  TableCell,
  Input,
  Select,
  SelectItem,
  Pagination,
  Chip
} from '@heroui/react'
import { IoSearchOutline } from 'react-icons/io5'

export type DataTableColumn<T> = {
  key: string
  label: string
  render?: (item: T) => ReactNode
}

export type DataTablePagination = {
  page: number
  limit: number
  pages: number
  total: number
}

interface DataTableProps<T> {
  columns: DataTableColumn<T>[]
  data: T[]
  pagination: DataTablePagination
  searchValue: string
  onSearchChange: (value: string) => void
  onPageChange: (page: number) => void
  onLimitChange: (limit: number) => void
  loading?: boolean
  getRowId?: (item: T) => string | number
  topRightSlot?: ReactNode
}

export function DataTable<T>({
  columns,
  data,
  pagination,
  searchValue,
  onSearchChange,
  onPageChange,
  onLimitChange,
  loading = false,
  getRowId,
  topRightSlot,
}: DataTableProps<T>) {
  const { page, pages, total, limit } = pagination

  const topContent = (
    <div className="flex justify-between items-center gap-2 flex-wrap">
      <div className="flex flex-col gap-2">
        <Input
          className="w-64"
          placeholder="Cari..."
          value={searchValue}
          onValueChange={onSearchChange}
          startContent={<IoSearchOutline />}
          size="sm"
          variant="bordered"
        />

        <span className='text-sm'>Total Data: {limit}</span>
      </div>
      
      <div className="flex flex-col gap-2">
        {topRightSlot}

        <div className='flex justify-end gap-2'>
          <span className='text-sm'>Limit:</span> 
          <select
            aria-label="Limit"
            className="w-max"
            onChange={(e) => onLimitChange(Number(e.target.value))}
          >
            {[10, 20, 50, 100].map((v) => (
              <option key={v} value={v}>{v}</option>
            ))}
          </select>
        </div>
      </div>
    </div>
  )

  const bottomContent = (
    <div className="flex items-center justify-between px-2 py-3">
      <p className="text-sm text-default-500">Halaman {page} dari {pages}</p>
      <Pagination
        page={page}
        total={pages}
        onChange={onPageChange}
        size="sm"
        showControls
      />
    </div>
  )

  return (
    <Table
      aria-label="Data table"
      topContent={topContent}
      bottomContent={bottomContent}
    >
      <TableHeader>
        {columns.map((col) => (
          <TableColumn key={col.key}>{col.label}</TableColumn>
        ))}
      </TableHeader>
      <TableBody loadingState={loading ? 'loading' : 'idle'} emptyContent="Tidak ada data">
        {data.map((item, idx) => {
          const rowKey = getRowId ? getRowId(item) : idx
          return (
            <TableRow key={rowKey}>
              {columns.map((col) => (
                <TableCell key={col.key}>
                  {col.render ? col.render(item) : (item as any)[col.key]}
                </TableCell>
              ))}
            </TableRow>
          )
        })}
      </TableBody>
    </Table>
  )
}
