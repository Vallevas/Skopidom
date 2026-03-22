import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { Package, ScanLine, Users, Settings, LogOut, Wrench } from 'lucide-react'
import { useAuth } from './auth-context'
import { cn } from '@/shared/ui/utils'

export function Layout() {
  const { t } = useTranslation()
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const isAdmin = user?.role === 'admin'

  function handleLogout() {
    logout()
    navigate('/login', { replace: true })
  }

  // Desktop nav: no scanner (desktop users don't need camera scanning)
  const desktopNav = [
    { to: '/items', icon: Package, label: t('nav.items') },
    ...(isAdmin ? [{ to: '/manage', icon: Wrench, label: t('nav.manage') }] : []),
    ...(isAdmin ? [{ to: '/users', icon: Users, label: t('nav.users') }] : []),
    { to: '/settings', icon: Settings, label: t('nav.settings') },
  ]

  // Mobile nav: scanner included, manage merged into settings
  const mobileNav = [
    { to: '/items', icon: Package, label: t('nav.items') },
    { to: '/scanner', icon: ScanLine, label: t('nav.scanner') },
    ...(isAdmin ? [{ to: '/manage', icon: Wrench, label: t('nav.manage') }] : []),
    ...(isAdmin ? [{ to: '/users', icon: Users, label: t('nav.users') }] : []),
    { to: '/settings', icon: Settings, label: t('nav.settings') },
  ]

  return (
    <div className="flex h-screen bg-background">
      {/* Sidebar (desktop) */}
      <aside className="hidden md:flex md:w-56 md:flex-col border-r">
        <div className="px-4 py-5 border-b">
          <span className="font-semibold text-lg">{t('app.name')}</span>
          <p className="text-xs text-muted-foreground">{t('app.tagline')}</p>
        </div>

        <nav className="flex-1 px-2 py-3 space-y-1">
          {desktopNav.map(({ to, icon: Icon, label }) => (
            <NavLink
              key={to}
              to={to}
              className={({ isActive }) =>
                cn(
                  'flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors',
                  isActive
                    ? 'bg-primary text-primary-foreground'
                    : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground',
                )
              }
            >
              <Icon size={16} />
              {label}
            </NavLink>
          ))}
        </nav>

        <div className="px-2 py-3 border-t">
          <div className="px-3 py-2 text-xs text-muted-foreground truncate">{user?.email}</div>
          <button
            onClick={handleLogout}
            className="flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors"
          >
            <LogOut size={16} />
            {t('nav.logout')}
          </button>
        </div>
      </aside>

      {/* Main content */}
      <div className="flex flex-1 flex-col overflow-hidden">
        <main className="flex-1 overflow-y-auto pb-20 md:pb-0">
          <Outlet />
        </main>

        {/* Bottom nav (mobile only) */}
        <nav className="md:hidden fixed bottom-0 inset-x-0 border-t bg-background safe-bottom z-50">
          <div className="flex justify-around">
            {mobileNav.map(({ to, icon: Icon, label }) => (
              <NavLink
                key={to}
                to={to}
                className={({ isActive }) =>
                  cn(
                    'flex flex-1 flex-col items-center gap-0.5 py-3 text-xs transition-colors',
                    isActive ? 'text-primary' : 'text-muted-foreground',
                  )
                }
              >
                <Icon size={20} />
                <span>{label}</span>
              </NavLink>
            ))}
          </div>
        </nav>
      </div>
    </div>
  )
}
