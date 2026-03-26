import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { Plus, Pencil, Trash2 } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { usersApi, translateError } from '@/shared/api/client'
import type { User, UserRole } from '@/shared/api/types'
import { useAuth } from '@/app/auth-context'
import { useToast } from '@/app/toast-context'

export function UsersPage() {
  const { t } = useTranslation()
  const { user: me } = useAuth()
  const toast = useToast()
  const queryClient = useQueryClient()
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<User | null>(null)
  const [confirmDelete, setConfirmDelete] = useState<User | null>(null)

  const { data: users = [], isLoading } = useQuery({
    queryKey: ['users'],
    queryFn: usersApi.list,
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => usersApi.delete(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['users'] }),
    onError: (err) => toast.error(t(translateError((err as Error).message))),
  })

  return (
    <div className="p-4 md:p-6 space-y-4">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold">{t('users.title')}</h1>
        <button
          onClick={() => { setEditing(null); setFormOpen(true) }}
          className="flex items-center gap-1.5 rounded-md bg-primary text-primary-foreground px-3 py-1.5 text-sm font-medium hover:bg-primary/90 transition-colors"
        >
          <Plus size={16} />
          {t('users.add')}
        </button>
      </div>

      {isLoading ? (
        <p className="text-muted-foreground text-center py-8">{t('common.loading')}</p>
      ) : (
        <div className="space-y-2">
          {users.map((u) => (
            <div
              key={u.id}
              className="flex items-center justify-between rounded-lg border bg-card px-4 py-3"
            >
              <div>
                <p className="text-sm font-medium">{u.full_name}</p>
                <p className="text-xs text-muted-foreground">{u.email}</p>
              </div>
              <div className="flex items-center gap-3">
                <span className="text-xs rounded-full bg-secondary px-2 py-0.5">
                  {u.role === 'admin' ? t('users.role_admin') : t('users.role_editor')}
                </span>
                {u.id !== me?.id && (
                  <>
                    <button
                      onClick={() => { setEditing(u); setFormOpen(true) }}
                      className="text-muted-foreground hover:text-foreground transition-colors"
                    >
                      <Pencil size={15} />
                    </button>
                    <button
                      onClick={() => setConfirmDelete(u)}
                      className="text-muted-foreground hover:text-destructive transition-colors"
                    >
                      <Trash2 size={15} />
                    </button>
                  </>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {formOpen && (
        <UserFormDialog
          user={editing}
          onClose={() => setFormOpen(false)}
          onSuccess={() => {
            queryClient.invalidateQueries({ queryKey: ['users'] })
            setFormOpen(false)
          }}
        />
      )}

      {confirmDelete && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-card rounded-xl p-6 max-w-sm w-full space-y-4">
            <p className="text-sm">
              {t('users.delete_confirm', { name: confirmDelete.full_name })}
            </p>
            <div className="flex gap-2">
              <button
                onClick={() => setConfirmDelete(null)}
                className="flex-1 rounded-md border px-4 py-2 text-sm hover:bg-accent transition-colors"
              >
                {t('common.cancel')}
              </button>
              <button
                onClick={() => {
                  deleteMutation.mutate(confirmDelete.id)
                  setConfirmDelete(null)
                }}
                className="flex-1 rounded-md bg-destructive text-destructive-foreground px-4 py-2 text-sm font-medium hover:bg-destructive/90 transition-colors"
              >
                {t('common.delete')}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

const createSchema = z.object({
  full_name: z.string().min(1),
  email: z.string().email(),
  password: z.string().min(8),
  role: z.enum(['admin', 'editor']),
})

const editSchema = z.object({
  full_name: z.string().min(1),
  role: z.enum(['admin', 'editor']),
})

function UserFormDialog({
  user,
  onClose,
  onSuccess,
}: {
  user: User | null
  onClose: () => void
  onSuccess: () => void
}) {
  const { t } = useTranslation()
  const toast = useToast()
  const isEdit = !!user

  type FormData = z.infer<typeof createSchema>

  const { register, handleSubmit, formState: { errors, isSubmitting } } =
    useForm<FormData>({
      resolver: zodResolver(isEdit ? editSchema : createSchema),
      defaultValues: user
        ? { full_name: user.full_name, role: user.role as UserRole }
        : undefined,
    })

  async function onSubmit(data: FormData) {
    try {
      if (isEdit) {
        await usersApi.update(user!.id, { full_name: data.full_name, role: data.role })
      } else {
        await usersApi.create(data)
      }
      onSuccess()
    } catch (err) {
      toast.error(t(translateError((err as Error).message)))
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50">
      <div className="w-full sm:max-w-md bg-card rounded-t-xl sm:rounded-xl p-6 space-y-4">
        <h2 className="text-lg font-semibold">
          {isEdit ? t('common.edit') : t('users.add')}
        </h2>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-3">
          <div className="space-y-1">
            <label className="text-sm font-medium">{t('users.full_name')}</label>
            <input
              className="w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
              {...register('full_name')}
            />
            {errors.full_name && (
              <p className="text-xs text-destructive">{t('common.required')}</p>
            )}
          </div>

          {!isEdit && (
            <>
              <div className="space-y-1">
                <label className="text-sm font-medium">{t('users.email')}</label>
                <input
                  type="email"
                  className="w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
                  {...register('email')}
                />
                {errors.email && (
                  <p className="text-xs text-destructive">{t('common.required')}</p>
                )}
              </div>
              <div className="space-y-1">
                <label className="text-sm font-medium">{t('users.password')}</label>
                <input
                  type="password"
                  className="w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
                  {...register('password')}
                />
                {errors.password && (
                  <p className="text-xs text-destructive">{t('common.required')}</p>
                )}
              </div>
            </>
          )}

          <div className="space-y-1">
            <label className="text-sm font-medium">{t('users.role')}</label>
            <select
              className="w-full rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
              {...register('role')}
            >
              <option value="editor">{t('users.role_editor')}</option>
              <option value="admin">{t('users.role_admin')}</option>
            </select>
          </div>

          <div className="flex gap-2 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 rounded-md border px-4 py-2 text-sm hover:bg-accent transition-colors"
            >
              {t('common.cancel')}
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className="flex-1 rounded-md bg-primary text-primary-foreground px-4 py-2 text-sm font-medium hover:bg-primary/90 disabled:opacity-50 transition-colors"
            >
              {t('common.save')}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
