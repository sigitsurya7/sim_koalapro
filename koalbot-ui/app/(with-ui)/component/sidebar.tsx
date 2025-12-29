"use client"

import { ReactNode, useEffect, useRef, useState } from 'react'
import { 
  Button,
  Divider,
  Avatar
} from "@heroui/react"
import { 
  FiChevronRight,
  FiChevronDown,
  FiSettings
} from 'react-icons/fi'
import { usePathname, useRouter } from 'next/navigation'
import { motion, AnimatePresence } from 'framer-motion'
import { API_BASE } from '@/lib/apiClient'
import toast from 'react-hot-toast'
import { NavItem, menuItems } from '@/types/menu'

interface SidebarProps {
  isOpen: boolean
  onClose: () => void
}

type AuthUser = {
  username?: string
  role?: string
}

export default function Sidebar({ isOpen, onClose }: SidebarProps) {
  const pathname = usePathname()
  const router = useRouter()
  const [openKeys, setOpenKeys] = useState<Set<string>>(new Set())
  const [isMobile, setIsMobile] = useState(false)
  const [serverOnline, setServerOnline] = useState(true)
  const prevStatus = useRef(true)
  const reconnectTimer = useRef<NodeJS.Timeout | null>(null)
  const reconnectToastId = useRef<string | null>(null)
  const [user, setUser] = useState<AuthUser>({})

  useEffect(() => {
    const media = window.matchMedia('(min-width: 1024px)')
    const update = () => setIsMobile(!media.matches)
    update()
    media.addEventListener('change', update)
    return () => media.removeEventListener('change', update)
  }, [])

  const isActive = (href?: string) => {
    if (!href) return false
    return pathname === href || pathname.startsWith(`${href}/`)
  }

  const findTrail = (items: NavItem[], path: string): string[] => {
    for (const item of items) {
      const key = item.href ?? item.title
      if (item.href && (path === item.href || path.startsWith(`${item.href}/`))) {
        return [key]
      }
      if (item.subItems) {
        const childTrail = findTrail(item.subItems, path)
        if (childTrail.length) return [key, ...childTrail]
      }
    }
    return []
  }

  useEffect(() => {
    const trail = findTrail(menuItems, pathname)
    setOpenKeys(new Set(trail))
  }, [pathname])

  useEffect(() => {
    if (typeof window === 'undefined') return
    const raw = localStorage.getItem('pm_user') ?? localStorage.getItem('pm-user')
    if (raw) {
      try {
        const parsed = JSON.parse(raw)
        setUser({ username: parsed.username, role: parsed.role ?? parsed.role_name })
      } catch (e) {
        console.error('Failed parse user', e)
      }
    }
  }, [])

  useEffect(() => {
    if (!prevStatus.current && serverOnline) {
      location.reload()
    }
    prevStatus.current = serverOnline
  }, [serverOnline])

  const toggleKey = (key: string) => {
    setOpenKeys(prev => {
      const next = new Set(prev)
      if (next.has(key)) {
        next.delete(key)
      } else {
        next.add(key)
      }
      return next
    })
  }

  const handleNavigate = (href: string) => {
    router.push(href)
    if (isMobile) onClose()
  }

  const filterMenuItems = (items: NavItem[], isAdmin: boolean): NavItem[] => {
    return items.reduce<NavItem[]>((acc, item) => {
      if (item.adminOnly && !isAdmin) return acc
      if (item.subItems) {
        const filteredChildren = filterMenuItems(item.subItems, isAdmin)
        if (filteredChildren.length === 0 && !item.href) return acc
        acc.push({ ...item, subItems: filteredChildren })
        return acc
      }
      acc.push(item)
      return acc
    }, [])
  }

  const renderNavItem = (item: NavItem, depth = 0): ReactNode => {
    const key = item.href ?? item.title
    const hasChildren = !!item.subItems?.length
    const active = isActive(item.href)
    const open = openKeys.has(key)

    const handlePress = () => {
      if (hasChildren) toggleKey(key)
      if (item.href) handleNavigate(item.href)
    }

    const paddingLeft = depth ? 12 + depth * 8 : 0
    const iconColor = active || open ? 'px-2 text-blue-600 dark:text-blue-400' : 'px-2 text-gray-600 dark:text-gray-300'

    return (
      <div key={key} className="space-y-2">
        <Button
          fullWidth
          onPress={handlePress}
          variant={(active || open) ? "flat" : "light"}
          color={active ? "primary" : "default"}
          className={`justify-start h-12 px-4 rounded-xl transition-all ${
            active || open
              ? 'bg-blue-50 dark:bg-blue-900/20'
              : 'hover:bg-gray-100/50 dark:hover:bg-gray-800/50'
          }`}
          startContent={
            item.icon ? (
              <span className={iconColor}>
                {item.icon}
              </span>
            ) : null
          }
          endContent={
            hasChildren ? (
              <span className="ml-auto text-gray-500">
                {open ? <FiChevronDown size={16} /> : <FiChevronRight size={16} />}
              </span>
            ) : null
          }
          style={{ paddingLeft }}
        >
            {item.title}
        </Button>

        {hasChildren && (
          <div className={`${open ? 'block' : 'hidden'} space-y-1`} style={{ paddingLeft: paddingLeft + 12 }}>
            {item.subItems!.map((child) => renderNavItem(child, depth + 1))}
          </div>
        )}
      </div>
    )
  }

  return (
    <>
      {/* Mobile Overlay */}
      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 lg:hidden z-30"
            onClick={onClose}
          />
        )}
      </AnimatePresence>

      {/* Sidebar */}
      <motion.aside
        initial={{ x: -300 }}
        animate={{ x: isOpen ? 0 : -300 }}
        transition={{ type: "spring", stiffness: 300, damping: 30 }}
        className={`
          fixed left-0 top-0 h-screen w-64
          bg-background
          dark:bg-content1
          backdrop-blur-xl backdrop-saturate-150
          border-r border-gray-200/40
          shadow-xl
          z-40
          lg:translate-x-0
          flex flex-col
        `}
        style={{
          backdropFilter: 'blur(20px) saturate(180%)',
          WebkitBackdropFilter: 'blur(20px) saturate(180%)'
        }}
      >
        {/* Sidebar Header */}
        <div className="p-6 border-b border-gray-200/30 dark:border-gray-700/30">
          <div className="flex items-left flex-col">
            <div className="flex items-left space-x-3">
              <div className="h-10 w-10 rounded-xl bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center">
                <span className="text-white font-bold text-lg">KP</span>
              </div>
              <div className='text-left'>
                <h2 className="font-bold text-gray-800 dark:text-white">Koala Pro</h2>
                <div className="flex gap-2 items-center mt-2">
                  <span className={`h-2.5 w-2.5 rounded-full ${serverOnline ? 'bg-green-500' : 'bg-red-500'}`} />
                  <span className="text-xs text-gray-600 dark:text-gray-300">
                    {serverOnline ? 'Server connected' : 'Server offline'}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Sidebar Content */}
        <div className="flex-1 overflow-y-auto py-4 px-3">
          {/* Main Menu */}
          <nav className="space-y-1 mb-6">
            {filterMenuItems(menuItems, user.role === 'admin').map((item) => renderNavItem(item))}
          </nav>
        </div>

        {/* Footer */}
        <div className="p-4 border-t border-gray-200/30 dark:border-gray-700/30">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <Avatar
                size="sm"
                src="https://i.pravatar.cc/150?img=32"
                fallback={(user.username?.[0] || 'U').toUpperCase()}
              />
              <div>
                <p className="text-sm font-medium text-gray-800 dark:text-white">{user.username || 'User'}</p>
                <p className="text-xs text-gray-500 dark:text-gray-400">{user.role || 'Role'}</p>
              </div>
            </div>
            <Button
              isIconOnly
              variant="light"
              size="sm"
            >
              <FiSettings size={16} />
            </Button>
          </div>
        </div>
      </motion.aside>
    </>
  )
}
