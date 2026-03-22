import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { itemsApi } from '@/shared/api/client'
import type { AuditAction } from '@/shared/api/types'
import { cn } from '@/shared/ui/utils'

const actionColors: Record<AuditAction, string> = {
  created: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
  updated: 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
  disposed: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
}

interface Props {
  itemId: number
}

export function AuditLog({ itemId }: Props) {
  const { t } = useTranslation()

  const { data: events = [], isLoading } = useQuery({
    queryKey: ['audit', itemId],
    queryFn: () => itemsApi.getAuditLog(itemId),
  })

  if (isLoading) {
    return <p className="text-sm text-muted-foreground">{t('common.loading')}</p>
  }

  return (
    <div className="space-y-2">
      <h2 className="text-sm font-semibold">{t('items.audit_log')}</h2>

      {events.length === 0 ? (
        <p className="text-sm text-muted-foreground">{t('items.no_audit')}</p>
      ) : (
        <ol className="relative border-l border-border space-y-4 pl-4">
          {events.map((event) => (
            <li key={event.id} className="relative">
              {/* Timeline dot */}
              <div className="absolute -left-[1.1rem] top-1 w-3 h-3 rounded-full border-2 border-background bg-muted-foreground" />

              <div className="space-y-1">
                <div className="flex items-center gap-2 flex-wrap">
                  <span
                    className={cn(
                      'rounded-full px-2 py-0.5 text-xs font-medium',
                      actionColors[event.action],
                    )}
                  >
                    {t(`audit.action_${event.action}`)}
                  </span>
                  <span className="text-xs text-muted-foreground">
                    {t('audit.by', { name: event.actor?.full_name ?? `#${event.actor_id}` })}
                  </span>
                  <span className="text-xs text-muted-foreground ml-auto">
                    {new Date(event.created_at).toLocaleString()}
                  </span>
                </div>

                {event.tx_hash && (
                  <p className="text-xs text-muted-foreground font-mono break-all">
                    {t('audit.blockchain_hash')}: {event.tx_hash}
                  </p>
                )}
              </div>
            </li>
          ))}
        </ol>
      )}
    </div>
  )
}
