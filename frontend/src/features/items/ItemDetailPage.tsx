import { useState, useRef, useCallback, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import {
  ArrowLeft, Pencil, Trash2, Camera, X, ChevronDown, ChevronUp,
  Wrench, ChevronLeft, ChevronRight,
} from 'lucide-react'
import { itemsApi, buildingsApi, roomsApi, translateError } from '@/shared/api/client'
import type { ItemPhoto } from '@/shared/api/types'
import { useAuth } from '@/app/auth-context'
import { useToast } from '@/app/toast-context'
import { AuditLog } from './AuditLog'
import { cn } from '@/shared/ui/utils'

// ── Lightbox ──────────────────────────────────────────────────────────────────

function Lightbox({
  photos,
  startIndex,
  onClose,
}: {
  photos: ItemPhoto[]
  startIndex: number
  onClose: () => void
}) {
  const [index, setIndex] = useState(startIndex)

  const prev = useCallback(() =>
    setIndex((i) => (i - 1 + photos.length) % photos.length), [photos.length])

  const next = useCallback(() =>
    setIndex((i) => (i + 1) % photos.length), [photos.length])

  // Keyboard navigation.
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
      if (e.key === 'ArrowLeft') prev()
      if (e.key === 'ArrowRight') next()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose, prev, next])

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/90"
      onClick={onClose}
    >
      {/* Close */}
      <button
        className="absolute top-4 right-4 text-white/70 hover:text-white transition-colors"
        onClick={onClose}
      >
        <X size={28} />
      </button>

      {/* Prev */}
      {photos.length > 1 && (
        <button
          className="absolute left-4 top-1/2 -translate-y-1/2 text-white/70 hover:text-white transition-colors p-2"
          onClick={(e) => { e.stopPropagation(); prev() }}
        >
          <ChevronLeft size={36} />
        </button>
      )}

      {/* Image */}
      <img
        src={photos[index].url}
        alt=""
        className="max-h-[90vh] max-w-[90vw] rounded-lg object-contain"
        onClick={(e) => e.stopPropagation()}
      />

      {/* Next */}
      {photos.length > 1 && (
        <button
          className="absolute right-4 top-1/2 -translate-y-1/2 text-white/70 hover:text-white transition-colors p-2"
          onClick={(e) => { e.stopPropagation(); next() }}
        >
          <ChevronRight size={36} />
        </button>
      )}

      {/* Counter */}
      {photos.length > 1 && (
        <span className="absolute bottom-4 left-1/2 -translate-x-1/2 text-white/60 text-sm">
          {index + 1} / {photos.length}
        </span>
      )}
    </div>
  )
}

// ── ItemDetailPage ────────────────────────────────────────────────────────────

export function ItemDetailPage() {
  const { id } = useParams<{ id: string }>()
  const itemId = Number(id)
  const { t } = useTranslation()
  const { user } = useAuth()
  const toast = useToast()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const [editingDesc, setEditingDesc] = useState(false)
  const [description, setDescription] = useState('')
  const [showAudit, setShowAudit] = useState(false)
  const [confirmDispose, setConfirmDispose] = useState(false)
  const [movingRoom, setMovingRoom] = useState(false)
  const [selectedBuildingId, setSelectedBuildingId] = useState<number | undefined>()
  const [selectedRoomId, setSelectedRoomId] = useState<number | undefined>()
  const [lightboxIndex, setLightboxIndex] = useState<number | null>(null)
  const photoInputRef = useRef<HTMLInputElement>(null)

  const { data: item, isLoading } = useQuery({
    queryKey: ['item', itemId],
    queryFn: () => itemsApi.getById(itemId),
    enabled: !!itemId,
  })

  const { data: photos = [] } = useQuery({
    queryKey: ['photos', itemId],
    queryFn: () => itemsApi.listPhotos(itemId),
    enabled: !!itemId,
  })

  const { data: buildings = [] } = useQuery({
    queryKey: ['buildings'],
    queryFn: buildingsApi.list,
    enabled: movingRoom,
  })

  const { data: rooms = [] } = useQuery({
    queryKey: ['rooms', selectedBuildingId],
    queryFn: () => roomsApi.list(selectedBuildingId),
    enabled: movingRoom && !!selectedBuildingId,
  })

  const updateMutation = useMutation({
    mutationFn: (desc: string) => itemsApi.update(itemId, { description: desc }),
    onSuccess: (updated) => {
      queryClient.setQueryData(['item', itemId], updated)
      queryClient.invalidateQueries({ queryKey: ['items'] })
      setEditingDesc(false)
      toast.success(t('common.success'))
    },
    onError: (err) => toast.error(t(translateError((err as Error).message))),
  })

  const photoMutation = useMutation({
    mutationFn: (file: File) => itemsApi.uploadPhoto(itemId, file),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['photos', itemId] })
      toast.success(t('items.photo_uploaded'))
    },
    onError: (err) => toast.error(t(translateError((err as Error).message))),
  })

  const deletePhotoMutation = useMutation({
    mutationFn: (photoId: number) => itemsApi.deletePhoto(itemId, photoId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['photos', itemId] }),
    onError: (err) => toast.error(t(translateError((err as Error).message))),
  })

  const repairMutation = useMutation({
    mutationFn: () => itemsApi.toggleRepair(itemId),
    onSuccess: (updated) => {
      queryClient.setQueryData(['item', itemId], updated)
      queryClient.invalidateQueries({ queryKey: ['items'] })
      queryClient.invalidateQueries({ queryKey: ['audit', itemId] })
      toast.success(t('common.success'))
    },
    onError: (err) => toast.error(t(translateError((err as Error).message))),
  })

  const disposeMutation = useMutation({
    mutationFn: () => itemsApi.dispose(itemId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['items'] })
      navigate('/items', { replace: true })
    },
    onError: (err) => toast.error(t(translateError((err as Error).message))),
  })

  const moveMutation = useMutation({
    mutationFn: (roomId: number) => itemsApi.moveToRoom(itemId, { room_id: roomId }),
    onSuccess: (updated) => {
      queryClient.setQueryData(['item', itemId], updated)
      queryClient.invalidateQueries({ queryKey: ['items'] })
      queryClient.invalidateQueries({ queryKey: ['audit', itemId] })
      setMovingRoom(false)
      setSelectedBuildingId(undefined)
      setSelectedRoomId(undefined)
      toast.success(t('items.moved'))
    },
    onError: (err) => toast.error(t(translateError((err as Error).message))),
  })

  if (isLoading) {
    return <div className="p-6 text-center text-muted-foreground">{t('common.loading')}</div>
  }
  if (!item) return null

  const isAdmin = user?.role === 'admin'
  const canEdit = item.status !== 'disposed'
  // Dispose is available for active and in_repair items (admin only).
  const canDispose = isAdmin && item.status !== 'disposed'

  const statusBadge = {
    disposed: { label: t('items.disposed_badge'), cls: 'bg-destructive/10 text-destructive' },
    in_repair: {
      label: t('items.in_repair_badge'),
      cls: 'bg-amber-100 text-amber-800 dark:bg-amber-900/30 dark:text-amber-400',
    },
    active: null,
  }[item.status]

  return (
    <div className="p-4 md:p-6 max-w-2xl mx-auto space-y-5">
      {/* Back + action buttons */}
      <div className="flex items-center justify-between gap-2">
        <button
          onClick={() => navigate(-1)}
          className="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          <ArrowLeft size={16} />
          {t('common.back')}
        </button>

        {canEdit && (
          <div className="flex items-center gap-2">
            {/* Repair toggle */}
            <button
              onClick={() => repairMutation.mutate()}
              disabled={repairMutation.isPending}
              className={cn(
                'flex items-center gap-1.5 rounded-md border px-3 py-1.5 text-sm transition-colors disabled:opacity-50',
                item.status === 'in_repair'
                  ? 'border-teal-500 text-teal-700 hover:bg-teal-50 dark:text-teal-400 dark:hover:bg-teal-900/20'
                  : 'border-amber-500 text-amber-700 hover:bg-amber-50 dark:text-amber-400 dark:hover:bg-amber-900/20',
              )}
            >
              <Wrench size={15} />
              {item.status === 'in_repair'
                ? t('items.return_from_repair')
                : t('items.send_to_repair')}
            </button>

            {/* Dispose — active AND in_repair, admin only */}
            {canDispose && (
              <button
                onClick={() => setConfirmDispose(true)}
                className="flex items-center gap-1.5 rounded-md border border-destructive text-destructive px-3 py-1.5 text-sm hover:bg-destructive/10 transition-colors"
              >
                <Trash2 size={15} />
                {t('items.dispose')}
              </button>
            )}
          </div>
        )}
      </div>

      {/* Photos */}
      {canEdit ? (
        <div className="space-y-2">
          {photos.length > 0 ? (
            <div className="flex gap-2 overflow-x-auto pb-1">
              {photos.map((photo, idx) => (
                <div
                  key={photo.id}
                  className="relative shrink-0 w-40 h-28 rounded-lg overflow-hidden bg-muted group"
                >
                  <img
                    src={photo.url}
                    alt=""
                    className="w-full h-full object-cover cursor-zoom-in"
                    onClick={() => setLightboxIndex(idx)}
                  />
                  <button
                    onClick={() => deletePhotoMutation.mutate(photo.id)}
                    className="absolute top-1 right-1 rounded-full bg-black/60 text-white p-0.5 hover:bg-black/80 transition-colors"
                  >
                    <X size={12} />
                  </button>
                </div>
              ))}
              <button
                onClick={() => photoInputRef.current?.click()}
                className="shrink-0 w-28 h-28 rounded-lg border-2 border-dashed border-border flex flex-col items-center justify-center gap-1 text-muted-foreground hover:border-primary hover:text-primary transition-colors"
              >
                <Camera size={20} />
                <span className="text-xs text-center">{t('items.upload_photo')}</span>
              </button>
            </div>
          ) : (
            <div
              className="rounded-xl border-2 border-dashed border-border flex flex-col items-center justify-center gap-2 py-10 text-muted-foreground cursor-pointer hover:border-primary hover:text-primary transition-colors"
              onClick={() => photoInputRef.current?.click()}
            >
              <Camera size={28} />
              <span className="text-sm">{t('items.upload_photo')}</span>
            </div>
          )}
          <input
            ref={photoInputRef}
            type="file"
            accept=".jpg,.jpeg,.png,.webp"
            className="hidden"
            onChange={(e) => {
              const file = e.target.files?.[0]
              if (file) { photoMutation.mutate(file); e.target.value = '' }
            }}
          />
        </div>
      ) : photos.length > 0 ? (
        <div className="flex gap-2 overflow-x-auto pb-1">
          {photos.map((photo, idx) => (
            <div
              key={photo.id}
              className="shrink-0 w-40 h-28 rounded-lg overflow-hidden bg-muted opacity-70 cursor-zoom-in"
              onClick={() => setLightboxIndex(idx)}
            >
              <img src={photo.url} alt="" className="w-full h-full object-cover" />
            </div>
          ))}
        </div>
      ) : null}

      {/* Header */}
      <div className="flex items-start justify-between gap-2">
        <div>
          <h1 className="text-xl font-semibold">{item.name}</h1>
          <div className="mt-1 space-y-0.5">
            <p className="text-xs text-muted-foreground">{t('items.barcode')}</p>
            <p className="text-sm font-mono">{item.barcode}</p>
          </div>
          <div className="mt-1 space-y-0.5">
            <p className="text-xs text-muted-foreground">{t('items.inventory_number')}</p>
            <p className="text-sm font-mono">{item.inventory_number}</p>
          </div>
        </div>
        {statusBadge && (
          <span className={cn('rounded-full text-xs px-2 py-0.5 shrink-0', statusBadge.cls)}>
            {statusBadge.label}
          </span>
        )}
      </div>

      {/* Info grid */}
      <div className="grid grid-cols-2 gap-3">
        <InfoRow label={t('items.category')} value={item.category?.name} />
        <InfoRow
          label={t('items.room')}
          value={`${item.room?.building?.name} — ${item.room?.name}`}
        />
        <InfoRow label={t('items.created_by')} value={item.creator?.full_name} />
        <InfoRow
          label={t('items.created_at')}
          value={new Date(item.created_at).toLocaleDateString()}
        />
      </div>

      {/* Description */}
      <div className="space-y-1">
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground font-medium">
            {t('items.description')}
          </span>
          {canEdit && !editingDesc && (
            <button
              onClick={() => { setDescription(item.description); setEditingDesc(true) }}
              className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors"
            >
              <Pencil size={12} />
              {t('common.edit')}
            </button>
          )}
        </div>
        {editingDesc ? (
          <div className="space-y-2">
            <textarea
              rows={3}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring resize-none"
            />
            <div className="flex gap-2">
              <button
                onClick={() => setEditingDesc(false)}
                className="flex-1 rounded-md border px-3 py-1.5 text-sm hover:bg-accent transition-colors"
              >
                {t('common.cancel')}
              </button>
              <button
                onClick={() => updateMutation.mutate(description)}
                disabled={updateMutation.isPending}
                className="flex-1 rounded-md bg-primary text-primary-foreground px-3 py-1.5 text-sm font-medium hover:bg-primary/90 disabled:opacity-50 transition-colors"
              >
                {t('common.save')}
              </button>
            </div>
          </div>
        ) : (
          <p className={cn('text-sm', !item.description && 'text-muted-foreground italic')}>
            {item.description || '—'}
          </p>
        )}
      </div>

      {/* Move to room — only for mutable items */}
      {canEdit && (
        <div className="space-y-2">
          <button
            onClick={() => setMovingRoom((v) => !v)}
            className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            {movingRoom ? <ChevronUp size={14} /> : <ChevronDown size={14} />}
            {t('items.move_room')}
          </button>

          {movingRoom && (
            <div className="rounded-lg border p-3 space-y-2">
              <select
                value={selectedBuildingId ?? ''}
                onChange={(e) => {
                  setSelectedBuildingId(e.target.value ? Number(e.target.value) : undefined)
                  setSelectedRoomId(undefined)
                }}
                className="w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
              >
                <option value="">{t('items.select_building')}</option>
                {buildings.map((b) => (
                  <option key={b.id} value={b.id}>{b.name}</option>
                ))}
              </select>

              {selectedBuildingId && (
                <select
                  value={selectedRoomId ?? ''}
                  onChange={(e) =>
                    setSelectedRoomId(e.target.value ? Number(e.target.value) : undefined)
                  }
                  className="w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
                >
                  <option value="">{t('items.select_room')}</option>
                  {rooms.map((r) => (
                    <option key={r.id} value={r.id}>{r.name}</option>
                  ))}
                </select>
              )}

              <div className="flex gap-2">
                <button
                  onClick={() => {
                    setMovingRoom(false)
                    setSelectedBuildingId(undefined)
                    setSelectedRoomId(undefined)
                  }}
                  className="flex-1 rounded-md border px-3 py-1.5 text-sm hover:bg-accent transition-colors"
                >
                  {t('common.cancel')}
                </button>
                <button
                  onClick={() => selectedRoomId && moveMutation.mutate(selectedRoomId)}
                  disabled={!selectedRoomId || moveMutation.isPending}
                  className="flex-1 rounded-md bg-primary text-primary-foreground px-3 py-1.5 text-sm font-medium hover:bg-primary/90 disabled:opacity-50 transition-colors"
                >
                  {t('items.confirm_move')}
                </button>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Audit log */}
      <div className="border-t pt-4 space-y-3">
        <button
          onClick={() => setShowAudit((v) => !v)}
          className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          {showAudit ? <ChevronUp size={14} /> : <ChevronDown size={14} />}
          {t('items.audit_log')}
        </button>
        {showAudit && <AuditLog itemId={itemId} />}
      </div>

      {/* Dispose confirm dialog */}
      {confirmDispose && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-card rounded-xl p-6 max-w-sm w-full space-y-4">
            <p className="text-sm">{t('items.dispose_confirm')}</p>
            <div className="flex gap-2">
              <button
                onClick={() => setConfirmDispose(false)}
                className="flex-1 rounded-md border px-4 py-2 text-sm hover:bg-accent transition-colors"
              >
                {t('common.cancel')}
              </button>
              <button
                onClick={() => { setConfirmDispose(false); disposeMutation.mutate() }}
                className="flex-1 rounded-md bg-destructive text-destructive-foreground px-4 py-2 text-sm font-medium hover:bg-destructive/90 transition-colors"
              >
                {t('items.dispose')}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Photo lightbox */}
      {lightboxIndex !== null && (
        <Lightbox
          photos={photos}
          startIndex={lightboxIndex}
          onClose={() => setLightboxIndex(null)}
        />
      )}
    </div>
  )
}

function InfoRow({ label, value }: { label: string; value?: string | null }) {
  return (
    <div>
      <p className="text-xs text-muted-foreground">{label}</p>
      <p className="text-sm font-medium">{value ?? '—'}</p>
    </div>
  )
}
