import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { itemsApi } from '@/shared/api/client'
import type { AuditAction, MovePayload } from '@/shared/api/types'
import { cn } from '@/shared/ui/utils'

const actionColors: Record<AuditAction, string> = {
  created:              'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
  updated:              'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-400',
  disposed:             'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
  moved:                'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400',
  sent_to_repair:       'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400',
  returned_from_repair: 'bg-teal-100 text-teal-800 dark:bg-teal-900/30 dark:text-teal-400',
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
      {events.length === 0 ? (
        <p className="text-sm text-muted-foreground">{t('items.no_audit')}</p>
      ) : (
        <ol className="relative border-l border-border space-y-4 pl-4">
          {events.map((event) => {
            let moveData: MovePayload | null = null
            if (event.action === 'moved') {
              try {
                moveData = JSON.parse(event.payload) as MovePayload
              } catch {
                // malformed payload — show without details
              }
            }

            // Convert underscore key to translation key:
            // sent_to_repair → audit.action_sent_to_repair
            const actionKey = `audit.action_${event.action}`

            return (
              <li key={event.id} className="relative">
                <div className="absolute -left-[1.1rem] top-1 w-3 h-3 rounded-full border-2 border-background bg-muted-foreground" />

                <div className="space-y-1">
                  <div className="flex items-center gap-2 flex-wrap">
                    <span
                      className={cn(
                        'rounded-full px-2 py-0.5 text-xs font-medium',
                        actionColors[event.action],
                      )}
                    >
                      {t(actionKey)}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      {t('audit.by', {
                        name: event.actor?.full_name ?? `#${event.actor_id}`,
                      })}
                    </span>
                    <span className="text-xs text-muted-foreground ml-auto">
                      {new Date(event.created_at).toLocaleString()}
                    </span>
                  </div>

                  {moveData && (
                    <div className="text-xs text-muted-foreground space-y-0.5 pl-1">
                      <p>
                        {t('audit.moved_from', {
                          from_building: moveData.from_building_name,
                          from_room: moveData.from_room_name,
                        })}
                      </p>
                      <p>
                        {t('audit.moved_to', {
                          to_building: moveData.to_building_name,
                          to_room: moveData.to_room_name,
                        })}
                      </p>
                    </div>
                  )}

                  {event.tx_hash && (
                    <p className="text-xs text-muted-foreground font-mono break-all">
                      {t('audit.blockchain_hash')}: {event.tx_hash}
                    </p>
                  )}
                </div>
              </li>
            )
          })}
        </ol>
      )}
    </div>
  )
}
