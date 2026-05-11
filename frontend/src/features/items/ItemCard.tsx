import { useTranslation } from 'react-i18next'
import { MapPin, Tag } from 'lucide-react'
import type { Item } from '@/shared/api/types'
import { cn } from '@/shared/ui/utils'

interface Props {
  item: Item
}

export function ItemCard({ item }: Props) {
  const { t } = useTranslation()
  const firstPhoto = item.photos?.[0]

  const statusBadge = {
    disposed: {
      label: t('items.disposed_badge'),
      cls: 'bg-destructive/10 text-destructive',
    },
    pending_disposal: {
      label: t('items.pending_disposal_badge'),
      cls: 'bg-orange-100 text-orange-800 dark:bg-orange-900/30 dark:text-orange-400',
    },
    in_repair: {
      label: t('items.in_repair_badge'),
      cls: 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400',
    },
    active: null,
  }[item.status]

  return (
    <div
      className={cn(
        'rounded-lg border bg-card p-4 space-y-3 hover:shadow-md transition-shadow',
        item.status === 'disposed' && 'opacity-60',
      )}
    >
      <div className="flex items-start justify-between gap-2">
        <div className="flex-1 min-w-0">
          <p className="font-medium truncate">{item.name}</p>
          <p className="text-xs text-muted-foreground font-mono mt-0.5">{item.barcode}</p>
        </div>
        {statusBadge && (
          <span className={cn('shrink-0 rounded-full text-xs px-2 py-0.5', statusBadge.cls)}>
            {statusBadge.label}
          </span>
        )}
      </div>

      <div className="space-y-1">
        <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
          <Tag size={12} />
          <span>{item.category?.name ?? '—'}</span>
        </div>
        <div className="flex items-center gap-1.5 text-xs text-muted-foreground">
          <MapPin size={12} />
          <span className="truncate">
            {item.room?.building?.name} — {item.room?.name}
          </span>
        </div>
      </div>

      {firstPhoto && (
        <img
          src={`${firstPhoto.mime_type};base64,${firstPhoto.base64_data}`}
          alt={item.name}
          className="w-full h-28 object-cover rounded-md"
        />
      )}
    </div>
  )
}
