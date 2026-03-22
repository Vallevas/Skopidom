import { useTranslation } from 'react-i18next'
import { useAuth } from '@/app/auth-context'
import { useNavigate } from 'react-router-dom'

export function SettingsPage() {
  const { t, i18n } = useTranslation()
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  function handleLogout() {
    logout()
    navigate('/login', { replace: true })
  }

  return (
    <div className="p-4 md:p-6 max-w-md space-y-6">
      <h1 className="text-xl font-semibold">{t('settings.title')}</h1>

      {/* Profile */}
      <section className="rounded-lg border bg-card p-4 space-y-1">
        <p className="text-sm font-medium">{user?.full_name}</p>
        <p className="text-xs text-muted-foreground">{user?.email}</p>
        <p className="text-xs text-muted-foreground">
          {user?.role === 'admin' ? t('users.role_admin') : t('users.role_editor')}
        </p>
      </section>

      {/* Language */}
      <section className="rounded-lg border bg-card p-4 space-y-3">
        <p className="text-sm font-medium">{t('settings.language')}</p>
        <div className="flex gap-2">
          {(['ru', 'en'] as const).map((lang) => (
            <button
              key={lang}
              onClick={() => i18n.changeLanguage(lang)}
              className={`flex-1 rounded-md border px-3 py-2 text-sm transition-colors ${
                i18n.language === lang
                  ? 'bg-primary text-primary-foreground border-primary'
                  : 'hover:bg-accent'
              }`}
            >
              {t(`settings.lang_${lang}`)}
            </button>
          ))}
        </div>
      </section>

      {/* Logout */}
      <button
        onClick={handleLogout}
        className="w-full rounded-lg border border-destructive text-destructive px-4 py-2.5 text-sm font-medium hover:bg-destructive/10 transition-colors"
      >
        {t('nav.logout')}
      </button>
    </div>
  )
}
