import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { itemsApi, roomsApi } from '@/shared/api/client'
import type { Building, Category } from '@/shared/api/types'

const schema = z.object({
  barcode: z.string().min(1),
  inventory_number: z.string().min(1),
  name: z.string().min(1),
  category_id: z.coerce.number().min(1),
  building_id: z.coerce.number().min(1),
  room_id: z.coerce.number().min(1),
  description: z.string().optional(),
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

  const mutation = useMutation({
    mutationFn: (data: Omit<FormData, 'building_id'>) => itemsApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['items'] })
      reset()
      setSelectedBuildingId(undefined)
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
