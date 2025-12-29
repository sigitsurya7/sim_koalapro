"use client"

import { useState, useEffect, useRef } from 'react'
import { FiSearch, FiFile, FiUsers, FiSettings, FiHome, FiBarChart2, FiCalendar, FiX } from 'react-icons/fi'
import { useRouter } from 'next/navigation'
import { menuItems, NavItem } from '@/types/menu'

interface CommandItem {
  id: string
  title: string
  description: string
  icon: React.ReactNode
  action: () => void
  category: string
}

interface CommandPaletteProps {
  isOpen?: boolean
  onClose?: () => void
}

export default function CommandPalette({ isOpen: controlledOpen, onClose }: CommandPaletteProps) {
  const router = useRouter()
  const [internalOpen, setInternalOpen] = useState(false)
  const isControlled = controlledOpen !== undefined
  const isOpen = isControlled ? controlledOpen : internalOpen

  const closePalette = () => {
    if (isControlled) {
      onClose?.()
    } else {
      setInternalOpen(false)
    }
  }
  const [query, setQuery] = useState('')
  const [selectedIndex, setSelectedIndex] = useState(0)
  const inputRef = useRef<HTMLInputElement>(null)

  // Command items
  const flattenMenu = (items: NavItem[], parentTitle?: string): CommandItem[] => {
    return items.flatMap((item) => {
      const fullTitle = parentTitle ? `${parentTitle} / ${item.title}` : item.title
      const current: CommandItem[] = item.href ? [{
        id: item.href,
        title: fullTitle,
        description: item.href,
        icon: item.icon ?? <FiFile size={18} />,
        action: () => router.push(item.href!),
        category: 'Menu'
      }] : []
      const children = item.subItems ? flattenMenu(item.subItems, fullTitle) : []
      return [...current, ...children]
    })
  }

  const commands: CommandItem[] = flattenMenu(menuItems)

  // Filter commands based on query
  const filteredCommands = commands.filter(cmd =>
    cmd.title.toLowerCase().includes(query.toLowerCase()) ||
    cmd.description.toLowerCase().includes(query.toLowerCase()) ||
    cmd.category.toLowerCase().includes(query.toLowerCase())
  )

  // Group by category
  const groupedCommands = filteredCommands.reduce((groups, cmd) => {
    if (!groups[cmd.category]) {
      groups[cmd.category] = []
    }
    groups[cmd.category].push(cmd)
    return groups
  }, {} as Record<string, CommandItem[]>)

  // Handle keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Cmd+K or Ctrl+K
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault()
        if (isControlled) {
          onClose?.()
        } else {
          setInternalOpen(true)
        }
      }
      
      // Escape to close
      if (e.key === 'Escape' && isOpen) {
        closePalette()
      }
      
      // Arrow navigation when open
      if (isOpen) {
        if (e.key === 'ArrowDown') {
          e.preventDefault()
          setSelectedIndex(prev => 
            prev < filteredCommands.length - 1 ? prev + 1 : prev
          )
        }
        if (e.key === 'ArrowUp') {
          e.preventDefault()
          setSelectedIndex(prev => prev > 0 ? prev - 1 : prev)
        }
        if (e.key === 'Enter' && filteredCommands[selectedIndex]) {
          e.preventDefault()
          filteredCommands[selectedIndex].action()
          closePalette()
        }
      }
    }

    document.addEventListener('keydown', handleKeyDown)
    return () => document.removeEventListener('keydown', handleKeyDown)
  }, [isOpen, filteredCommands, selectedIndex])

  // Focus input when modal opens
  useEffect(() => {
    if (isOpen && inputRef.current) {
      setTimeout(() => inputRef.current?.focus(), 100)
      setQuery('')
      setSelectedIndex(0)
    }
  }, [isOpen])

  // Close when clicking outside
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      const target = e.target as HTMLElement
      if (isOpen && !target.closest('.command-palette')) {
        closePalette()
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [isOpen])

  const handleCommandClick = (action: () => void) => {
    action()
    closePalette()
  }

  if (!isOpen) return null

  return (
    <>
      {/* Overlay */}
      <div className="fixed inset-0 bg-black/50 backdrop-blur-sm z-55" />
      
      {/* Modal */}
      <div className="fixed inset-0 flex items-start justify-center pt-20 px-4 z-60">
        <div className="command-palette w-full max-w-2xl bg-white dark:bg-content1 rounded-xl shadow-2xl border border-gray-200 dark:border-gray-700 overflow-hidden">
          
          {/* Search Input */}
          <div className="p-4 border-b border-gray-200 dark:bg-content1">
            <div className="relative">
              <div className="absolute inset-y-0 left-3 flex items-center pointer-events-none">
                <FiSearch className="text-gray-400" size={20} />
              </div>
              <input
                ref={inputRef}
                type="text"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder="Type a command or search..."
                className="w-full pl-10 pr-10 py-3 bg-gray-50 dark:bg-gray-800 border-none rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
              />
              <div className="absolute inset-y-0 right-3 flex items-center">
                <kbd className="px-2 py-1 text-xs font-semibold text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 rounded border border-gray-300 dark:border-gray-600">
                  ESC
                </kbd>
              </div>
            </div>
          </div>

          {/* Results */}
          <div className="max-h-96 overflow-y-auto">
            {filteredCommands.length === 0 ? (
              <div className="p-8 text-center text-gray-500 dark:text-gray-400">
                <p>No results found for "{query}"</p>
                <p className="text-sm mt-2">Try a different search term</p>
              </div>
            ) : (
              Object.entries(groupedCommands).map(([category, items]) => (
                <div key={category} className="border-t border-gray-100 dark:border-gray-800 first:border-t-0">
                  <div className="px-4 py-2 bg-gray-50 dark:bg-gray-800">
                    <span className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                      {category}
                    </span>
                  </div>
                  <div className="py-1">
                    {items.map((cmd, index) => {
                      const globalIndex = filteredCommands.findIndex(c => c.id === cmd.id)
                      const isSelected = globalIndex === selectedIndex
                      
                      return (
                        <button
                          key={cmd.id}
                          onClick={() => handleCommandClick(cmd.action)}
                          className={`w-full flex items-center gap-3 px-4 py-3 text-left hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors ${
                            isSelected ? 'bg-gray-100 dark:bg-gray-800' : ''
                          }`}
                          onMouseEnter={() => setSelectedIndex(globalIndex)}
                        >
                          <div className={`p-2 rounded-lg ${
                            isSelected 
                              ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400' 
                              : 'bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400'
                          }`}>
                            {cmd.icon}
                          </div>
                          <div className="flex-1 min-w-0">
                            <div className="flex items-center gap-2">
                              <span className="font-medium text-gray-900 dark:text-white truncate">
                                {cmd.title}
                              </span>
                              {isSelected && (
                                <kbd className="ml-2 px-1.5 py-0.5 text-xs font-medium bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded border border-gray-300 dark:border-gray-600">
                                  ↵
                                </kbd>
                              )}
                            </div>
                            <p className="text-sm text-gray-500 dark:text-gray-400 truncate">
                              {cmd.description}
                            </p>
                          </div>
                        </button>
                      )
                    })}
                  </div>
                </div>
              ))
            )}
          </div>

          {/* Footer */}
          <div className="p-3 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-content1">
            <div className="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
              <div className="flex items-center gap-4">
                <div className="flex items-center gap-1">
                  <kbd className="px-1.5 py-0.5 rounded border border-gray-300 dark:border-gray-600">↑</kbd>
                  <kbd className="px-1.5 py-0.5 rounded border border-gray-300 dark:border-gray-600">↓</kbd>
                  <span>to navigate</span>
                </div>
                <div className="flex items-center gap-1">
                  <kbd className="px-1.5 py-0.5 rounded border border-gray-300 dark:border-gray-600">↵</kbd>
                  <span>to select</span>
                </div>
              </div>
              <div className="flex items-center gap-1">
                <kbd className="px-1.5 py-0.5 rounded border border-gray-300 dark:border-gray-600">ESC</kbd>
                <span>to close</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}
