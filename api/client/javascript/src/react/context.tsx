'use client'

import { createContext, useContext } from 'react'
import type { MeterForge } from '../portal/index.js'

export * from '../portal/index.js'

export const MeterForgeContext = createContext<MeterForge | null>(null)

export type MeterForgeProviderProps = {
  children?: React.ReactNode
  value: MeterForge | null
}

export function MeterForgeProvider({ children, value }: MeterForgeProviderProps) {
  return (
    <MeterForgeContext.Provider value={value}>
      {children}
    </MeterForgeContext.Provider>
  )
}

export function useMeterForge() {
  const context = useContext(MeterForgeContext)
  if (typeof context === 'undefined') {
    throw new Error('useMeterForge must be used within a MeterForgeProvider')
  }

  return context
}
