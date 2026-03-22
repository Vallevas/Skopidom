import { useEffect, useRef, useState, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { BrowserMultiFormatReader, NotFoundException } from '@zxing/library'
import { itemsApi, ApiClientError } from '@/shared/api/client'

type ScanState =
  | { status: 'idle' }
  | { status: 'scanning' }
  | { status: 'found'; barcode: string }
  | { status: 'not_found'; barcode: string }
  | { status: 'error'; message: string }

export function ScannerPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const videoRef = useRef<HTMLVideoElement>(null)
  const readerRef = useRef<BrowserMultiFormatReader | null>(null)
  const [state, setState] = useState<ScanState>({ status: 'idle' })
  const [manual, setManual] = useState('')
  const [showManual, setShowManual] = useState(false)

  const handleBarcode = useCallback(
    async (barcode: string) => {
      setState({ status: 'found', barcode })
      try {
        const item = await itemsApi.getByBarcode(barcode)
        navigate(`/items/${item.id}`)
      } catch (err) {
        if (err instanceof ApiClientError && err.status === 404) {
          setState({ status: 'not_found', barcode })
        } else {
          setState({ status: 'error', message: t('common.unknown_error') })
        }
      }
    },
    [navigate, t],
  )

  useEffect(() => {
    const reader = new BrowserMultiFormatReader()
    readerRef.current = reader

    async function startScanner() {
      if (!videoRef.current) return
      try {
        setState({ status: 'scanning' })
        await reader.decodeFromVideoDevice(
          // ZXing expects string | null, not undefined
          null,
          videoRef.current,
          (result, err) => {
            if (result) {
              reader.reset()
              handleBarcode(result.getText())
            }
            if (err && !(err instanceof NotFoundException)) {
              console.error('ZXing error:', err)
            }
          },
        )
      } catch {
        setState({ status: 'error', message: t('scanner.camera_error') })
      }
    }

    startScanner()
    return () => { reader.reset() }
  }, [handleBarcode, t])

  function handleManualSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (manual.trim()) {
      readerRef.current?.reset()
      handleBarcode(manual.trim())
    }
  }

  function handleReset() {
    setState({ status: 'scanning' })
    setManual('')
  }

  return (
    <div className="flex flex-col h-full">
      <div className="p-4 border-b">
        <h1 className="text-lg font-semibold">{t('scanner.title')}</h1>
      </div>

      <div className="relative flex-1 bg-black overflow-hidden">
        <video ref={videoRef} className="w-full h-full object-cover" muted playsInline />

        {state.status === 'scanning' && (
          <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
            <div className="relative w-64 h-40">
              {(['tl', 'tr', 'bl', 'br'] as const).map((corner) => (
                <div
                  key={corner}
                  className={`absolute w-8 h-8 border-white border-2 ${
                    corner === 'tl' ? 'top-0 left-0 border-r-0 border-b-0' :
                    corner === 'tr' ? 'top-0 right-0 border-l-0 border-b-0' :
                    corner === 'bl' ? 'bottom-0 left-0 border-r-0 border-t-0' :
                    'bottom-0 right-0 border-l-0 border-t-0'
                  }`}
                />
              ))}
              <p className="absolute -bottom-8 left-0 right-0 text-center text-white text-sm">
                {t('scanner.instruction')}
              </p>
            </div>
          </div>
        )}

        {state.status !== 'scanning' && state.status !== 'idle' && (
          <div className="absolute inset-0 flex flex-col items-center justify-center bg-black/70 text-white p-6 text-center">
            {state.status === 'found' && (
              <p className="text-lg font-medium">{t('scanner.item_found')}</p>
            )}
            {state.status === 'not_found' && (
              <>
                <p className="text-lg font-medium text-red-400">{t('scanner.item_not_found')}</p>
                <p className="text-sm text-white/70 mt-1 font-mono">{state.barcode}</p>
                <button
                  onClick={handleReset}
                  className="mt-4 rounded-md bg-white/20 px-4 py-2 text-sm hover:bg-white/30 transition-colors"
                >
                  {t('scanner.scan')}
                </button>
              </>
            )}
            {state.status === 'error' && (
              <>
                <p className="text-red-400">{state.message}</p>
                <button
                  onClick={handleReset}
                  className="mt-4 rounded-md bg-white/20 px-4 py-2 text-sm hover:bg-white/30 transition-colors"
                >
                  {t('scanner.scan')}
                </button>
              </>
            )}
          </div>
        )}
      </div>

      <div className="p-4 border-t bg-background">
        {!showManual ? (
          <button
            onClick={() => setShowManual(true)}
            className="w-full text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            {t('scanner.enter_manually')}
          </button>
        ) : (
          <form onSubmit={handleManualSubmit} className="flex gap-2">
            <input
              type="text"
              value={manual}
              onChange={(e) => setManual(e.target.value)}
              placeholder={t('scanner.manual_barcode')}
              autoFocus
              className="flex-1 rounded-md border bg-background px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-ring"
            />
            <button
              type="submit"
              className="rounded-md bg-primary text-primary-foreground px-4 py-2 text-sm font-medium hover:bg-primary/90 transition-colors"
            >
              {t('scanner.scan')}
            </button>
          </form>
        )}
      </div>
    </div>
  )
}
