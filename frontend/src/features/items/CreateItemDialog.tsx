import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { itemsApi, roomsApi } from '@/shared/api/client'
import type { Building, Category, Item } from '@/shared/api/types'
import { Search, X } from 'lucide-react'

const schema = z.object({
  barcode: z.string().min(1),
  inventory_number: z.string().min(1),
  name: z.string().min(1),
  category_id: z.coerce.number().min(1),
  building_id: z.coerce.number().min(1),
  room_id: z.coerce.number().min(1),
  description: z.string().optional(),
  linked_item_id: z.coerce.number().min(1).optional().or(z.literal(0)),
})

type FormData = z.infer<typeof schema>

interface Props {
  open: boolean
  onClose: () => void
  categories: Category[]
  buildings: Building[]
  initialBarcode?: string
}

export function CreateItemDialog({ open, onClose, categories, buildings, initialBarcode }: Props) {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const [selectedBuildingId, setSelectedBuildingId] = useState<number | undefined>()
  const [linkedItemSearchOpen, setLinkedItemSearchOpen] = useState(false)
  const [linkedItemQuery, setLinkedItemQuery] = useState('')
  const [selectedLinkedItem, setSelectedLinkedItem] = useState<Item | null>(null)

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    formState: { errors },
  } = useForm<FormData>({ resolver: zodResolver(schema) })

  // Rooms loaded only after building is selected.
  const { data: rooms = [] } = useQuery({
    queryKey: ['rooms', selectedBuildingId],
    queryFn: () => roomsApi.list(selectedBuildingId),
    enabled: !!selectedBuildingId,
  })

  const { data: allItems = [] } = useQuery({
    queryKey: ['items'],
    queryFn: () => itemsApi.list(),
    enabled: linkedItemSearchOpen,
  })

  const filteredItems = allItems.filter(
    (item) =>
      item.name.toLowerCase().includes(linkedItemQuery.toLowerCase()) ||
      item.barcode.toLowerCase().includes(linkedItemQuery.toLowerCase()) ||
      item.inventory_number.toLowerCase().includes(linkedItemQuery.toLowerCase())
  )

  const mutation = useMutation({
    mutationFn: async (data: Omit<FormData, 'building_id' | 'linked_item_id'> & { linked_item_id?: number }) => {
      const newItem = await itemsApi.create(data)
      // If a linked item was selected, create the relation after item creation
      if (data.linked_item_id && data.linked_item_id > 0) {
        await itemsApi.linkItems({ item_id_1: newItem.id, item_id_2: data.linked_item_id })
      }
      return newItem
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['items'] })
      reset()
      setSelectedBuildingId(undefined)
      setSelectedLinkedItem(null)
      setLinkedItemSearchOpen(false)
      onClose()
    },
  })

  // Set initial barcode when dialog opens
  useEffect(() => {
    if (open && initialBarcode) {
      setValue('barcode', initialBarcode)
    }
  }, [open, initialBarcode, setValue])

  function onSubmit({ building_id: _b, ...data }: FormData) {
    mutation.mutate(data)
  }

  function handleClose() {
    reset()
    setSelectedBuildingId(undefined)
    setSelectedLinkedItem(null)
    setLinkedItemSearchOpen(false)
    setLinkedItemQuery('')
    onClose()
  }

  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 overflow-y-auto">
      <div className="w-full sm:max-w-md bg-card rounded-t-xl sm:rounded-xl shadow-lg my-auto">
        <div className="p-6 space-y-4 max-h-[90vh] overflow-y-auto">
          <h2 className="text-lg font-semibold">{t('items.add')}</h2>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-3">
            <Field label={t('items.barcode')} error={!!errors.barcode}>
              <input className={inputCls} {...register('barcode')} placeholder="SN-2024-001" />
            </Field>

            <Field label={t('items.inventory_number')} error={!!errors.inventory_number}>
              <input className={inputCls} {...register('inventory_number')} placeholder="INV-2024-001" />
            </Field>

            <Field label={t('items.name')} error={!!errors.name}>
              <input className={inputCls} {...register('name')} />
            </Field>

            <Field label={t('items.category')} error={!!errors.category_id}>
              <select className={inputCls} {...register('category_id')}>
                <option value="">—</option>
                {categories.map((c) => (
                  <option key={c.id} value={c.id}>{c.name}</option>
                ))}
              </select>
            </Field>

            {/* Building first */}
            <Field label={t('items.building')} error={!!errors.building_id}>
              <select
                className={inputCls}
                {...register('building_id')}
                onChange={(e) => {
                  const val = Number(e.target.value)
                  setSelectedBuildingId(val || undefined)
                  setValue('room_id', 0) // reset room on building change
                }}
              >
                <option value="">—</option>
                {buildings.map((b) => (
                  <option key={b.id} value={b.id}>{b.name}</option>
                ))}
              </select>
            </Field>

            {/* Room — only after building */}
            <Field label={t('items.room')} error={!!errors.room_id}>
              <select
                className={inputCls}
                {...register('room_id')}
                disabled={!selectedBuildingId}
              >
                <option value="">
                  {selectedBuildingId ? '—' : t('items.select_building_first')}
                </option>
                {rooms.map((r) => (
                  <option key={r.id} value={r.id}>{r.name}</option>
                ))}
              </select>
            </Field>

            <Field label={t('items.description')}>
              <textarea className={inputCls + ' resize-none'} rows={2} {...register('description')} />
            </Field>

            {/* Linked item selector */}
            <Field label={t('items.linked_items')}>
              {!linkedItemSearchOpen ? (
                <div className="space-y-2">
                  {selectedLinkedItem ? (
                    <div className="flex items-center justify-between gap-2 rounded-md border px-3 py-2">
                      <div className="flex-1">
                        <span className="font-medium text-sm">{selectedLinkedItem.name}</span>
                        <span className="text-xs text-muted-foreground ml-2">
                          ({selectedLinkedItem.barcode})
                        </span>
                      </div>
                      <button
                        type="button"
                        onClick={() => {
                          setSelectedLinkedItem(null)
                          setValue('linked_item_id', 0)
                        }}
                        className="rounded-md p-1 hover:bg-accent transition-colors"
                      >
                        <X size={14} />
                      </button>
                    </div>
                  ) : (
                    <button
                      type="button"
                      onClick={() => setLinkedItemSearchOpen(true)}
                      className="flex items-center gap-2 w-full rounded-md border px-3 py-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
                    >
                      <Search size={16} />
                      {t('items.link_item')}
                    </button>
                  )}
                </div>
              ) : (
                <div className="space-y-2">
                  <div className="flex items-center gap-2">
                    <div className="relative flex-1">
                      <Search className="absolute left-2 top-1/2 -translate-y-1/2 text-muted-foreground" size={16} />
                      <input
                        type="text"
                        value={linkedItemQuery}
                        onChange={(e) => setLinkedItemQuery(e.target.value)}
                        placeholder={t('items.search')}
                        className="w-full rounded-md border bg-background pl-8 pr-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
                        autoFocus
                      />
                    </div>
                    <button
                      type="button"
                      onClick={() => {
                        setLinkedItemSearchOpen(false)
                        setLinkedItemQuery('')
                      }}
                      className="rounded-md p-1 hover:bg-accent transition-colors"
                    >
                      <X size={16} />
                    </button>
                  </div>

                  <div className="max-h-48 overflow-y-auto space-y-1 rounded-md border">
                    {filteredItems.length === 0 ? (
                      <p className="text-sm text-muted-foreground p-2">{t('items.no_items_found')}</p>
                    ) : (
                      filteredItems.map((item) => (
                        <button
                          key={item.id}
                          type="button"
                          onClick={() => {
                            setSelectedLinkedItem(item)
                            setValue('linked_item_id', item.id)
                            setLinkedItemSearchOpen(false)
                            setLinkedItemQuery('')
                          }}
                          className="w-full text-left px-3 py-2 text-sm hover:bg-accent transition-colors"
                        >
                          <span className="font-medium">{item.name}</span>
                          <span className="text-xs text-muted-foreground ml-2">
                            ({item.barcode})
                          </span>
                        </button>
                      ))
                    )}
                  </div>
                </div>
              )}
            </Field>

            {mutation.error && (
              <p className="text-sm text-destructive">{(mutation.error as Error).message}</p>
            )}

            <div className="flex gap-2 pt-2">
              <button
                type="button"
                onClick={handleClose}
                className="flex-1 rounded-md border px-4 py-2 text-sm hover:bg-accent transition-colors"
              >
                {t('common.cancel')}
              </button>
              <button
                type="submit"
                disabled={mutation.isPending}
                className="flex-1 rounded-md bg-primary text-primary-foreground px-4 py-2 text-sm font-medium hover:bg-primary/90 disabled:opacity-50 transition-colors"
              >
                {t('common.save')}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}

const inputCls =
  'w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring disabled:opacity-50'

function Field({
  label,
  error,
  children,
}: {
  label: string
  error?: boolean
  children: React.ReactNode
}) {
  const { t } = useTranslation()
  return (
    <div className="space-y-1">
      <label className="text-sm font-medium">{label}</label>
      {children}
      {error && <p className="text-xs text-destructive">{t('common.required')}</p>}
    </div>
  )
}
