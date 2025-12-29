"use client"

import { useState } from 'react'
import Sidebar from './component/sidebar'
import Navbar from './component/navbar'
import CommandPalette from '@/components/command-pallete'

export default function MainLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const [sidebarOpen, setSidebarOpen] = useState(true)

  return (
    <div className="min-h-screen">
      <CommandPalette />
      {/* Container untuk Sidebar & Content */}
      <div className="flex">
        
        {/* SIDEBAR - Fixed Left */}
        <Sidebar
          isOpen={sidebarOpen} 
          onClose={() => setSidebarOpen(false)} 
        />
        
        {/* Container untuk Navbar & Main Content */}
        <div className={`flex-1 overflow-auto transition-all duration-300 ease-in-out ${
          sidebarOpen ? 'ml-64' : 'ml-0'
        }`}>
          
          {/* NAVBAR - Fixed di dalam container sebelah kanan */}
          <Navbar
            sidebarOpen={sidebarOpen} 
            toggleSidebar={() => setSidebarOpen(!sidebarOpen)} 
          />
          
          {/* MAIN CONTENT - Scrollable */}
          <main> {/* pt-16 untuk offset navbar */}
            <div className="p-6">
              {children}
            </div>
          </main>
        </div>
      </div>
    </div>
  )
}