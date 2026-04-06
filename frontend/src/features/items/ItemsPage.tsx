import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { Link } from 'react-router-dom'
import { Plus, Search } from 'lucide-react'
import { itemsApi, categoriesApi, buildingsApi, roomsApi } from '@/shared/api/client'
import type { ItemFilter, ItemStatus } from '@/shared/api/types'
import { ItemCard } from './ItemCard'
import { CreateItemDialog } from './CreateItemDialog'

export function ItemsPage() {
  const { t } = useTranslation()
  const [search, setSearch] = useState('')
  const [filter, setFilter] = useState<ItemFilter>({ status: 'active' })
  const [selectedBuildingId, setSelectedBuildingId] = useState<number | undefined>()
  const [createOpen, setCreateOpen] = useState(false)

  const { data: items = [], isLoading } = useQuery({
    queryKey: ['items', filter],
    queryFn: () => itemsApi.list(filter),
  })

  const { data: categories = [] } = useQuery({
    queryKey: ['categories'],
    queryFn: categoriesApi.list,
  })

  const { data: buildings = [] } = useQuery({
    queryKey: ['buildings'],
    queryFn: buildingsApi.list,
  })

  const { data: rooms = [] } = useQuery({
    queryKey: ['rooms', selectedBuildingId],
    queryFn: () => roomsApi.list(selectedBuildingId),
    enabled: !!selectedBuildingId,
  })

  function handleBuildingChange(buildingId: number | undefined) {
    setSelectedBuildingId(buildingId)
    setFilter((f) => ({ ...f, room_id: undefined }))
  }

  const filtered = items.filter((item) => {
    // Search filter
    if (search) {
      const q = search.toLowerCase()
      if (
        !item.name.toLowerCase().includes(q) &&
        !item.barcode.toLowerCase().includes(q)
      ) {
        return false
      }
    }
    
    // Building filter (when building selected but no specific room)
    if (selectedBuildingId && !filter.room_id) {
      if (item.room?.building_id !== selectedBuildingId) {
        return false
      }
    }
    
    return true
  })

  return (
    <div className="p-4 md:p-6 space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold">{t('items.title')}</h1>
        <button
          onClick={() => setCreateOpen(true)}
          className="flex items-center gap-1.5 rounded-md bg-primary text-primary-foreground px-3 py-1.5 text-sm font-medium hover:bg-primary/90 transition-colors"
        >
          <Plus size={16} />
          {t('items.add')}
        </button>
      </div>

      <div className="relative">
        <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
        <input
          type="search"
          placeholder={t('items.search')}
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="w-full rounded-md border bg-background pl-9 pr-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
        />
      </div>

      <div className="flex flex-wrap gap-2">
        <select
          value={filter.status ?? 'active'}
          onChange={(e) =>
            setFilter((f) => ({ ...f, status: e.target.value as ItemStatus }))
          }
          className="rounded-md border bg-background px-2 py-1.5 text-sm outline-none focus:ring-2 focus:ring-ring"
        >
          <option value="active">{t('items.status_active')}</option>
          <option value="in_repair">{t('items.status_in_repair')}</option>
          <option value="disposed">{t('items.status_disposed')}</option>
        </select>

        <select
          value={filter.category_id ?? ''}
          onChange={(e) =>
            setFilter((f) => ({
              ...f,
              category_id: e.target.value ? Number(e.target.value) : undefined,
            }))
          }
          className="rounded-md border bg-background px-2 py-1.5 text-sm outline-none focus:ring-2 focus:ring-ring"
        >
          <option value="">{t('items.all_categories')}</option>
          {categories.map((c) => (
            <option key={c.id} value={c.id}>{c.name}</option>
          ))}
        </select>

        <select
          value={selectedBuildingId ?? ''}
          onChange={(e) =>
            handleBuildingChange(e.target.value ? Number(e.target.value) : undefined)
          }
          className="rounded-md border bg-background px-2 py-1.5 text-sm outline-none focus:ring-2 focus:ring-ring"
        >
          <option value="">{t('items.all_buildings')}</option>
          {buildings.map((b) => (
            <option key={b.id} value={b.id}>{b.name}</option>
          ))}
        </select>

        {selectedBuildingId && (
          <select
            value={filter.room_id ?? ''}
            onChange={(e) =>
              setFilter((f) => ({
                ...f,
                room_id: e.target.value ? Number(e.target.value) : undefined,
              }))
            }
            className="rounded-md border bg-background px-2 py-1.5 text-sm outline-none focus:ring-2 focus:ring-ring"
          >
            <option value="">{t('items.all_rooms')}</option>
            {rooms.map((r) => (
              <option key={r.id} value={r.id}>{r.name}</option>
            ))}
          </select>
        )}
      </div>

      {isLoading ? (
        <div className="text-center py-12 text-muted-foreground">{t('common.loading')}</div>
      ) : filtered.length === 0 ? (
        <div className="text-center py-12 text-muted-foreground">{t('items.empty')}</div>
      ) : (
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {filtered.map((item) => (
            <Link key={item.id} to={`/items/${item.id}`}>
              <ItemCard item={item} />
            </Link>
          ))}
        </div>
      )}

      <CreateItemDialog
        open={createOpen}
        onClose={() => setCreateOpen(false)}
        categories={categories}
        buildings={buildings}
      />
    </div>
  )
}
