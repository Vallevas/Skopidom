import {
  createContext,
  useContext,
  useState,
  useCallback,
  type ReactNode,
} from 'react'
import { tokenStorage } from '@/shared/api/client'
import type { User } from '@/shared/api/types'

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  login: (token: string, user: User) => void
  logout: () => void
}

const AuthContext = createContext<AuthState | null>(null)

// Parse user from existing token on page load.
function getStoredUser(): User | null {
  const token = tokenStorage.get()
  if (!token) return null
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    // Token carries uid and role — stored user is hydrated on first API call.
    // We keep a minimal object to avoid a round-trip on refresh.
    return payload as User
  } catch {
    return null
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(getStoredUser)

  const login = useCallback((token: string, user: User) => {
    tokenStorage.set(token)
    setUser(user)
  }, [])

  const logout = useCallback(() => {
    tokenStorage.clear()
    setUser(null)
  }, [])

  return (
    <AuthContext.Provider
      value={{ user, isAuthenticated: !!user, login, logout }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthState {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used inside AuthProvider')
  return ctx
}
