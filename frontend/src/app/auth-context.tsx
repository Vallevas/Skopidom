import {
  createContext,
  useContext,
  useState,
  useCallback,
  type ReactNode,
} from 'react'
import { tokenStorage } from '@/shared/api/client'
import type { User } from '@/shared/api/types'

const USER_KEY = 'skopidom_user'

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  login: (token: string, user: User) => void
  logout: () => void
}

const AuthContext = createContext<AuthState | null>(null)

// Read the full user object saved on login.
// JWT payload only carries uid + role, not full_name / email,
// so parsing the token on page load results in an incomplete user object.
function getStoredUser(): User | null {
  const raw = localStorage.getItem(USER_KEY)
  if (!raw) return null
  try {
    return JSON.parse(raw) as User
  } catch {
    return null
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(getStoredUser)

  const login = useCallback((token: string, user: User) => {
    tokenStorage.set(token)
    localStorage.setItem(USER_KEY, JSON.stringify(user))
    setUser(user)
  }, [])

  const logout = useCallback(() => {
    tokenStorage.clear()
    localStorage.removeItem(USER_KEY)
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
