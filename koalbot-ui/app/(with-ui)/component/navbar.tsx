"use client"

import { 
  Navbar as HeroNavbar, 
  NavbarContent, 
  NavbarBrand,
  Avatar,
  Dropdown,
  DropdownTrigger,
  DropdownMenu,
  DropdownItem,
  Button,
  Badge,
  Chip,
  Input,
  Kbd
} from "@heroui/react"
import { 
  FiBell, 
  FiSearch, 
  FiSettings,
  FiUser,
  FiLogOut,
  FiMoon,
  FiSun,
  FiMenu,
  FiX
} from 'react-icons/fi'
import { FaRegCircle, FaRegCircleDot } from "react-icons/fa6"
import { FaCashRegister } from "react-icons/fa";
import { useTheme } from 'next-themes'
import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { IoSearchOutline } from "react-icons/io5"
import toast from 'react-hot-toast'
import { apiClient } from '@/lib/apiClient'
import CommandPalette from '@/components/command-pallete'

interface NavbarProps {
  sidebarOpen: boolean
  toggleSidebar: () => void
}

type AuthUser = {
  username?: string
  role_name?: string
}

export default function Navbar({ sidebarOpen, toggleSidebar }: NavbarProps) {
  const [mounted, setMounted] = useState(false)
  const { theme, setTheme } = useTheme()
  const router = useRouter()
  const [loggingOut, setLoggingOut] = useState(false)
  const [cmdOpen, setCmdOpen] = useState(false)
  const [user, setUser] = useState<AuthUser>({})

  const handleLogout = async () => {
    if (loggingOut) return
    setLoggingOut(true)
    await toast.promise(
      (async () => {
        try {
          await apiClient.post('/logout')
        } catch (err) {
          // ignore API error, still clear client state
        }
        document.cookie = 'token=; path=/; max-age=0'
        localStorage.removeItem('pm-user')
        router.push('/login')
      })(),
      {
        loading: 'Logging out...',
        success: 'Logged out',
        error: 'Logout failed'
      }
    )
    setLoggingOut(false)
  }

  useEffect(() => {
    setMounted(true)
  }, [])

  useEffect(() => {
    if (typeof window === 'undefined') return
    const raw = localStorage.getItem('pm-user')
    if (raw) {
      try {
        const parsed = JSON.parse(raw)
        setUser({ username: parsed.username, role_name: parsed.role_name })
      } catch (e) {
        console.error('Failed parse user', e)
      }
    }
  }, [])

  return (
    <section className="mx-6 mt-2">
        <HeroNavbar 
        maxWidth="full"
        className="sticky top-0 rounded-lg shadow-sm left-0 right-0 dark:bg-content1 backdrop-blur-lg backdrop-saturate-150 z-50"
        style={{
            backdropFilter: 'blur(20px) saturate(180%)',
            WebkitBackdropFilter: 'blur(20px) saturate(180%)'
        }}
        >
        {/* Left Section */}
        <NavbarContent className="basis-1/5" justify="start">
          {
            !sidebarOpen && (
              <Button
                isIconOnly
                variant="light"
                size="sm"
                onPress={toggleSidebar}
                aria-label="Toggle sidebar"
              >
              {sidebarOpen ? <FaRegCircleDot size={14} /> : <FiMenu size={14} />}
              </Button>
            )
          }
          
            <Input
              className="w-max cursor-pointer"
              placeholder="Search Anything"
              startContent={<IoSearchOutline />}
              endContent={<Kbd keys={["command"]}>K</Kbd>}
              onClick={() => setCmdOpen(true)}
              readOnly
            />
        </NavbarContent>

        {/* Right Section */}
        <NavbarContent className="basis-1/5" justify="end">
            {/* Theme Toggle */}
            <div className="hidden sm:flex gap-2 items-center">
              <Button
                  isIconOnly
                  variant="light"
                  size="sm"
                  aria-label="Toggle theme"
                  className="rounded-full"
                  onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
                  disabled={!mounted}
              >
                  {mounted && theme === 'dark' ? (
                  <FiSun className="text-gray-600 dark:text-gray-300" size={18} />
                  ) : (
                  <FiMoon className="text-gray-600 dark:text-gray-300" size={18} />
                  )}
              </Button>
            </div>

            {/* User Profile */}
            <Dropdown placement="bottom-end">
            <DropdownTrigger>
                <Avatar
                isBordered
                size="sm"
                className="cursor-pointer transition-transform hover:scale-105"
                src="https://i.pravatar.cc/150?img=32"
                fallback={(user.username?.[0] || 'U').toUpperCase()}
                />
            </DropdownTrigger>
            <DropdownMenu aria-label="Profile Actions" onAction={(key) => {
              if (key === 'logout') {
                handleLogout()
              } else if (key === 'account') {
                router.push('/account')
              }
            }}>
                <DropdownItem key="profile" className="h-14 gap-2">
                <div className="flex flex-col">
                    <p className="font-semibold">{user.username || 'User'}</p>
                    <p className="text-gray-500 text-sm">{user.role_name || 'Role'}</p>
                </div>
                </DropdownItem>
                <DropdownItem key="account" startContent={<FiUser size={16} />}>
                  My Account
                </DropdownItem>
                <DropdownItem key="settings" startContent={<FiSettings size={16} />}>
                  Settings
                </DropdownItem>
                <DropdownItem 
                  key="logout" 
                  startContent={<FiLogOut size={16} />}
                  color="danger"
                  className="text-danger"
                >
                Log Out
                </DropdownItem>
            </DropdownMenu>
            </Dropdown>
        </NavbarContent>
        </HeroNavbar>
        <CommandPalette isOpen={cmdOpen} onClose={() => setCmdOpen(false)} />
    </section>
  )
}
