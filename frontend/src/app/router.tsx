import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from './auth-context'
import { Layout } from './layout'

import { LoginPage } from '@/features/auth/LoginPage'
import { ItemsPage } from '@/features/items/ItemsPage'
import { ItemDetailPage } from '@/features/items/ItemDetailPage'
import { ScannerPage } from '@/features/scanner/ScannerPage'
import { UsersPage } from '@/features/users/UsersPage'
import { ManagePage } from '@/features/manage/ManagePage'
import { SettingsPage } from '@/features/settings/SettingsPage'

function RequireAuth({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth()
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />
}

function RequireAdmin({ children }: { children: React.ReactNode }) {
  const { user } = useAuth()
  return user?.role === 'admin' ? <>{children}</> : <Navigate to="/" replace />
}

export function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />

        <Route
          element={
            <RequireAuth>
              <Layout />
            </RequireAuth>
          }
        >
          <Route index element={<Navigate to="/items" replace />} />
          <Route path="/items" element={<ItemsPage />} />
          <Route path="/items/:id" element={<ItemDetailPage />} />
          <Route path="/scanner" element={<ScannerPage />} />
          <Route
            path="/manage"
            element={
              <RequireAdmin>
                <ManagePage />
              </RequireAdmin>
            }
          />
          <Route
            path="/users"
            element={
              <RequireAdmin>
                <UsersPage />
              </RequireAdmin>
            }
          />
          <Route path="/settings" element={<SettingsPage />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  )
}
