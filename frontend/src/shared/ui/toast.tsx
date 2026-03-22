import { useEffect } from 'react'
import { X } from 'lucide-react'
import { cn } from './utils'

export type ToastType = 'success' | 'error' | 'info'

export interface ToastItem {
  id: string
  message: string
  type: ToastType
}

interface Props {
  toasts: ToastItem[]
  onDismiss: (id: string) => void
}

const styles: Record<ToastType, string> = {
  success: 'bg-green-600 text-white',
  error: 'bg-destructive text-destructive-foreground',
  info: 'bg-primary text-primary-foreground',
}

function Toast({ toast, onDismiss }: { toast: ToastItem; onDismiss: () => void }) {
  useEffect(() => {
    const timer = setTimeout(onDismiss, 4000)
    return () => clearTimeout(timer)
  }, [onDismiss])

  return (
    <div
      className={cn(
        'flex items-center justify-between gap-3 rounded-lg px-4 py-3 shadow-lg text-sm max-w-sm w-full animate-in slide-in-from-bottom-2',
        styles[toast.type],
      )}
    >
      <span>{toast.message}</span>
      <button onClick={onDismiss} className="shrink-0 opacity-80 hover:opacity-100">
        <X size={14} />
      </button>
    </div>
  )
}

export function ToastContainer({ toasts, onDismiss }: Props) {
  if (toasts.length === 0) return null

  return (
    <div className="fixed bottom-24 md:bottom-6 left-1/2 -translate-x-1/2 z-[100] flex flex-col gap-2 items-center w-full px-4 pointer-events-none">
      {toasts.map((t) => (
        <div key={t.id} className="pointer-events-auto w-full max-w-sm">
          <Toast toast={t} onDismiss={() => onDismiss(t.id)} />
        </div>
      ))}
    </div>
  )
}
