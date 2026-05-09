import { useTranslation } from 'react-i18next'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { itemsApi, translateError } from '@/shared/api/client'
import type { ItemRelation } from '@/shared/api/types'
import { useToast } from '@/app/toast-context'
import { Link, X } from 'lucide-react'
import { cn } from '@/shared/ui/utils'

interface Props {
  itemId: number
  relations: ItemRelation[]
  canEdit: boolean
}

export function LinkedItemsList({ itemId, relations, canEdit }: Props) {
  const { t } = useTranslation()
  const toast = useToast()
  const queryClient = useQueryClient()

  const unlinkMutation = useMutation({
    mutationFn: (relationId: number) => itemsApi.unlinkItems(relationId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['relations', itemId] })
      toast.success(t('common.success'))
    },
    onError: (err) => toast.error(t(translateError((err as Error).message))),
  })

  if (!relations || relations.length === 0) {
    return null
  }

  // Get the related item for each relation (the one that is not the current item)
  const linkedItems = relations.map((rel) => {
    const relatedItemId = rel.item_id_1 === itemId ? rel.item_id_2 : rel.item_id_1
    const relatedItem = rel.related_item
    return {
      id: relatedItemId,
      relationId: rel.id,
      name: relatedItem?.name || t('common.unknown'),
      barcode: relatedItem?.barcode || '',
      room: relatedItem?.room,
    }
  })

  return (
    <div className="space-y-2">
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <Link size={16} />
        <span className="font-medium">{t('items.linked_items')}</span>
      </div>
      <div className="space-y-1">
        {linkedItems.map((item) => (
          <div
            key={item.relationId}
            className="flex items-center justify-between gap-2 rounded-md border px-3 py-2"
          >
            <a
              href={`/items/${item.id}`}
              className="flex-1 text-sm hover:underline flex items-center gap-2"
            >
              <Link size={14} className="text-muted-foreground" />
              <span className="font-medium">{item.name}</span>
              <span className="text-xs text-muted-foreground">({item.barcode})</span>
              {item.room && (
                <span className="text-xs text-muted-foreground">— {item.room.name}</span>
              )}
            </a>
            {canEdit && (
              <button
                onClick={() => unlinkMutation.mutate(item.relationId)}
                disabled={unlinkMutation.isPending}
                className="rounded-md p-1 hover:bg-destructive/10 text-muted-foreground hover:text-destructive transition-colors disabled:opacity-50"
                title={t('items.unlink_item')}
              >
                <X size={14} />
              </button>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}
