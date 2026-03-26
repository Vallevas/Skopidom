import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Plus, Pencil, Trash2 } from 'lucide-react'
import { buildingsApi, categoriesApi, roomsApi } from '@/shared/api/client'
import { useToast } from '@/app/toast-context'
import { cn } from '@/shared/ui/utils'

type Tab = 'categories' | 'buildings' | 'rooms'

export function ManagePage() {
  const { t } = useTranslation()
  const [tab, setTab] = useState<Tab>('categories')

  return (
    <div className="p-4 md:p-6 space-y-4">
      <h1 className="text-xl font-semibold">{t('nav.manage')}</h1>

      <div className="flex gap-1 border-b">
        {(['categories', 'buildings', 'rooms'] as Tab[]).map((t_) => (
          <button
            key={t_}
            onClick={() => setTab(t_)}
            className={cn(
              'px-4 py-2 text-sm font-medium border-b-2 -mb-px transition-colors',
              tab === t_
                ? 'border-primary text-primary'
                : 'border-transparent text-muted-foreground hover:text-foreground',
            )}
          >
            {t(`manage.${t_}`)}
          </button>
        ))}
      </div>

      {tab === 'categories' && <CategoriesTab />}
      {tab === 'buildings' && <BuildingsTab />}
      {tab === 'rooms' && <RoomsTab />}
    </div>
  )
}

function BuildingsTab() {
  const { t } = useTranslation()
  const toast = useToast()
  const queryClient = useQueryClient()
  const [form, setForm] = useState<{ id?: number; name: string; address: string } | null>(null)

  const { data: buildings = [] } = useQuery({
    queryKey: ['buildings'],
    queryFn: buildingsApi.list,
  })

  const saveMutation = useMutation({
    mutationFn: (data: { id?: number; name: string; address: string }) =>
      data.id
        ? buildingsApi.update(data.id, data.name, data.address)
        : buildingsApi.create(data.name, data.address),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['buildings'] })
      setForm(null)
      toast.success(t('common.success'))
    },
    onError: (err) => toast.error((err as Error).message),
  })

  const deleteMutation = useMutation({
    mutationFn: buildingsApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['buildings'] })
      toast.success(t('common.success'))
    },
    onError: (err) => toast.error((err as Error).message),
  })

  return (
    <div className="space-y-3">
      <button
        onClick={() => setForm({ name: '', address: '' })}
        className="flex items-center gap-1.5 rounded-md bg-primary text-primary-foreground px-3 py-1.5 text-sm font-medium hover:bg-primary/90 transition-colors"
      >
        <Plus size={15} />
        {t('manage.add_building')}
      </button>

      <div className="space-y-2">
        {buildings.map((b) => (
          <div key={b.id} className="flex items-center justify-between rounded-lg border bg-card px-4 py-3">
            <div>
              <p className="text-sm font-medium">{b.name}</p>
              <p className="text-xs text-muted-foreground">{b.address}</p>
            </div>
            <div className="flex gap-2">
              <button
                onClick={() => setForm({ id: b.id, name: b.name, address: b.address })}
                className="text-muted-foreground hover:text-foreground transition-colors"
              >
                <Pencil size={15} />
              </button>
              <button
                onClick={() => deleteMutation.mutate(b.id)}
                className="text-muted-foreground hover:text-destructive transition-colors"
              >
                <Trash2 size={15} />
              </button>
            </div>
          </div>
        ))}
      </div>

      {form !== null && (
        <InlineForm
          fields={[
            { key: 'name', label: t('manage.name'), value: form.name },
            { key: 'address', label: t('manage.address'), value: form.address },
          ]}
          onSave={(values) =>
            saveMutation.mutate({ id: form.id, name: values.name, address: values.address })
          }
          onCancel={() => setForm(null)}
          isPending={saveMutation.isPending}
        />
      )}
    </div>
  )
}

function CategoriesTab() {
  const { t } = useTranslation()
  const toast = useToast()
  const queryClient = useQueryClient()
  const [form, setForm] = useState<{ id?: number; name: string } | null>(null)

  const { data: categories = [] } = useQuery({
    queryKey: ['categories'],
    queryFn: categoriesApi.list,
  })

  const saveMutation = useMutation({
    mutationFn: (data: { id?: number; name: string }) =>
      data.id ? categoriesApi.update(data.id, data.name) : categoriesApi.create(data.name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] })
      setForm(null)
      toast.success(t('common.success'))
    },
    onError: (err) => toast.error((err as Error).message),
  })

  const deleteMutation = useMutation({
    mutationFn: categoriesApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] })
      toast.success(t('common.success'))
    },
    onError: (err) => toast.error((err as Error).message),
  })

  return (
    <div className="space-y-3">
      <button
        onClick={() => setForm({ name: '' })}
        className="flex items-center gap-1.5 rounded-md bg-primary text-primary-foreground px-3 py-1.5 text-sm font-medium hover:bg-primary/90 transition-colors"
      >
        <Plus size={15} />
        {t('manage.add_category')}
      </button>

      <div className="space-y-2">
        {categories.map((c) => (
          <div key={c.id} className="flex items-center justify-between rounded-lg border bg-card px-4 py-3">
            <p className="text-sm font-medium">{c.name}</p>
            <div className="flex gap-2">
              <button
                onClick={() => setForm({ id: c.id, name: c.name })}
                className="text-muted-foreground hover:text-foreground transition-colors"
              >
                <Pencil size={15} />
              </button>
              <button
                onClick={() => deleteMutation.mutate(c.id)}
                className="text-muted-foreground hover:text-destructive transition-colors"
              >
                <Trash2 size={15} />
              </button>
            </div>
          </div>
        ))}
      </div>

      {form !== null && (
        <InlineForm
          fields={[{ key: 'name', label: t('manage.name'), value: form.name }]}
          onSave={(values) => saveMutation.mutate({ id: form.id, name: values.name })}
          onCancel={() => setForm(null)}
          isPending={saveMutation.isPending}
        />
      )}
    </div>
  )
}

function RoomsTab() {
  const { t } = useTranslation()
  const toast = useToast()
  const queryClient = useQueryClient()
  const [selectedBuildingId, setSelectedBuildingId] = useState<number | undefined>()
  const [form, setForm] = useState<{ id?: number; name: string } | null>(null)

  const { data: buildings = [] } = useQuery({
    queryKey: ['buildings'],
    queryFn: buildingsApi.list,
  })

  const { data: rooms = [] } = useQuery({
    queryKey: ['rooms', selectedBuildingId],
    queryFn: () => roomsApi.list(selectedBuildingId),
    enabled: !!selectedBuildingId,
  })

  const saveMutation = useMutation({
    mutationFn: (data: { id?: number; name: string }) =>
      data.id && selectedBuildingId
        ? roomsApi.update(data.id, data.name, selectedBuildingId)
        : roomsApi.create(data.name, selectedBuildingId!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms', selectedBuildingId] })
      setForm(null)
      toast.success(t('common.success'))
    },
    onError: (err) => toast.error((err as Error).message),
  })

  const deleteMutation = useMutation({
    mutationFn: roomsApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rooms', selectedBuildingId] })
      toast.success(t('common.success'))
    },
    onError: (err) => toast.error((err as Error).message),
  })

  return (
    <div className="space-y-3">
      <select
        value={selectedBuildingId ?? ''}
        onChange={(e) => {
          setSelectedBuildingId(e.target.value ? Number(e.target.value) : undefined)
          setForm(null)
        }}
        className="w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
      >
        <option value="">{t('manage.select_building')}</option>
        {buildings.map((b) => (
          <option key={b.id} value={b.id}>{b.name}</option>
        ))}
      </select>

      {selectedBuildingId && (
        <>
          <button
            onClick={() => setForm({ name: '' })}
            className="flex items-center gap-1.5 rounded-md bg-primary text-primary-foreground px-3 py-1.5 text-sm font-medium hover:bg-primary/90 transition-colors"
          >
            <Plus size={15} />
            {t('manage.add_room')}
          </button>

          <div className="space-y-2">
            {rooms.map((r) => (
              <div key={r.id} className="flex items-center justify-between rounded-lg border bg-card px-4 py-3">
                <p className="text-sm font-medium">{r.name}</p>
                <div className="flex gap-2">
                  <button
                    onClick={() => setForm({ id: r.id, name: r.name })}
                    className="text-muted-foreground hover:text-foreground transition-colors"
                  >
                    <Pencil size={15} />
                  </button>
                  <button
                    onClick={() => deleteMutation.mutate(r.id)}
                    className="text-muted-foreground hover:text-destructive transition-colors"
                  >
                    <Trash2 size={15} />
                  </button>
                </div>
              </div>
            ))}
          </div>
        </>
      )}

      {form !== null && (
        <InlineForm
          fields={[{ key: 'name', label: t('manage.name'), value: form.name }]}
          onSave={(values) => saveMutation.mutate({ id: form.id, name: values.name })}
          onCancel={() => setForm(null)}
          isPending={saveMutation.isPending}
        />
      )}
    </div>
  )
}

function InlineForm({
  fields,
  onSave,
  onCancel,
  isPending,
}: {
  fields: { key: string; label: string; value: string }[]
  onSave: (values: Record<string, string>) => void
  onCancel: () => void
  isPending: boolean
}) {
  const { t } = useTranslation()
  const [values, setValues] = useState<Record<string, string>>(
    Object.fromEntries(fields.map((f) => [f.key, f.value])),
  )

  return (
    <div className="rounded-lg border bg-card p-4 space-y-3">
      {fields.map((f) => (
        <div key={f.key} className="space-y-1">
          <label className="text-sm font-medium">{f.label}</label>
          <input
            value={values[f.key]}
            onChange={(e) => setValues((v) => ({ ...v, [f.key]: e.target.value }))}
            className="w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
          />
        </div>
      ))}
      <div className="flex gap-2">
        <button
          onClick={onCancel}
          className="flex-1 rounded-md border px-3 py-1.5 text-sm hover:bg-accent transition-colors"
        >
          {t('common.cancel')}
        </button>
        <button
          onClick={() => onSave(values)}
          disabled={isPending}
          className="flex-1 rounded-md bg-primary text-primary-foreground px-3 py-1.5 text-sm font-medium hover:bg-primary/90 disabled:opacity-50 transition-colors"
        >
          {t('common.save')}
        </button>
      </div>
    </div>
  )
}
