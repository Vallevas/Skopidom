import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQuery } from '@tanstack/react-query'
import { itemsApi } from '@/shared/api/client'
import type { Item } from '@/shared/api/types'
import { Search, X } from 'lucide-react'

interface Props {
  currentItemId: number
  onLinkSelect: (itemId: number) => void
}

export function ItemLinkSelector({ currentItemId, onLinkSelect }: Props) {
  const { t } = useTranslation()
  const [searchQuery, setSearchQuery] = useState('')
  const [isOpen, setIsOpen] = useState(false)

  const { data: allItems = [], isLoading } = useQuery({
    queryKey: ['items'],
    queryFn: () => itemsApi.list(),
    enabled: isOpen,
  })

  const filteredItems = allItems.filter(
    (item) =>
      item.id !== currentItemId &&
      (item.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        item.barcode.toLowerCase().includes(searchQuery.toLowerCase()) ||
        item.inventory_number.toLowerCase().includes(searchQuery.toLowerCase()))
  )

  const handleSelect = (item: Item) => {
    onLinkSelect(item.id)
    setSearchQuery('')
    setIsOpen(false)
  }

  if (!isOpen) {
    return (
      <button
        onClick={() => setIsOpen(true)}
        className="flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
      >
        <Search size={16} />
        {t('items.link_item')}
      </button>
    )
  }

  return (
    <div className="space-y-2">
      <div className="flex items-center gap-2">
        <div className="relative flex-1">
          <Search className="absolute left-2 top-1/2 -translate-y-1/2 text-muted-foreground" size={16} />
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder={t('items.search')}
            className="w-full rounded-md border bg-background pl-8 pr-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
            autoFocus
          />
        </div>
        <button
          onClick={() => {
            setIsOpen(false)
            setSearchQuery('')
          }}
          className="rounded-md p-1 hover:bg-accent transition-colors"
        >
          <X size={16} />
        </button>
      </div>

      {isLoading ? (
        <p className="text-sm text-muted-foreground">{t('common.loading')}</p>
      ) : filteredItems.length === 0 ? (
        <p className="text-sm text-muted-foreground">{t('items.no_items_found')}</p>
      ) : (
        <div className="max-h-48 overflow-y-auto space-y-1 rounded-md border">
          {filteredItems.map((item) => (
            <button
              key={item.id}
              onClick={() => handleSelect(item)}
              className="w-full text-left px-3 py-2 text-sm hover:bg-accent transition-colors flex items-center justify-between"
            >
              <div>
                <span className="font-medium">{item.name}</span>
                <span className="text-xs text-muted-foreground ml-2">
                  ({item.barcode})
                </span>
              </div>
              <span className="text-xs text-muted-foreground">
                {item.room?.name}
              </span>
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
