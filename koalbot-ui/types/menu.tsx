import { ReactNode } from 'react'
import { 
  FiHome,
  FiBarChart2,
  FiSettings,
  FiHelpCircle,
  FiCreditCard,
  FiDatabase,
  FiUser,
  FiLock
} from 'react-icons/fi'

export type NavItem = {
  title: string
  href?: string
  icon?: ReactNode
  badge?: ReactNode
  subItems?: NavItem[]
  adminOnly?: boolean
}

export const menuItems: NavItem[] = [
  {
    title: "Dashboard",
    icon: <FiHome size={20} />,
    href: "/dashboard",
    badge: null
  },
  {
    title: "Data Pengguna",
    icon: <FiDatabase size={20} />,
    subItems: [
      { title: "Stockity", href: "/master/member_stockity" },
      { title: "Binomo", href: "/master/member_binomo" },
      { title: "Olymptrade", href: "/master/member_olymptrade" },
    ]
  },
  {
    title: "Master Pengguna",
    icon: <FiUser size={20} />,
    href: "/users",
    adminOnly: true,
    badge: null
  },
]
