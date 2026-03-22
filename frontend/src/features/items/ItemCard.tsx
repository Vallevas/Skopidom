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
          <p className="text-xs text-muted-foreground font-mono mt-0.5">
            {item.barcode}
          </p>
        </div>
        {item.status === 'disposed' && (
          <span className="shrink-0 rounded-full bg-destructive/10 text-destructive text-xs px-2 py-0.5">
            {t('items.disposed_badge')}
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
          src={firstPhoto.url}
          alt={item.name}
          className="w-full h-28 object-cover rounded-md"
        />
      )}
    </div>
  )
}
