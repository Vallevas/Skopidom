import { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useQuery } from '@tanstack/react-query'
import { BrowserMultiFormatReader, NotFoundException } from '@zxing/library'
import { itemsApi, categoriesApi, buildingsApi, ApiClientError } from '@/shared/api/client'
import { CreateItemDialog } from '@/features/items/CreateItemDialog'

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
  const [shouldScan, setShouldScan] = useState(true)
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [scannedBarcode, setScannedBarcode] = useState('')

  const { data: categories = [] } = useQuery({
    queryKey: ['categories'],
    queryFn: categoriesApi.list,
  })

  const { data: buildings = [] } = useQuery({
    queryKey: ['buildings'],
    queryFn: buildingsApi.list,
  })

  useEffect(() => {
    if (!shouldScan) return

    const reader = new BrowserMultiFormatReader()
    readerRef.current = reader

    async function startScanner() {
      if (!videoRef.current) return
      try {
        setState({ status: 'scanning' })
        await reader.decodeFromVideoDevice(
          null,
          videoRef.current,
          async (result, err) => {
            if (result) {
              const barcode = result.getText()
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
    return () => { 
      reader.reset()
    }
  }, [shouldScan, navigate, t])

  function handleManualSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (manual.trim()) {
      setShouldScan(false)
      readerRef.current?.reset()
      
      const barcode = manual.trim()
      setState({ status: 'found', barcode })
      
      itemsApi.getByBarcode(barcode)
        .then(item => navigate(`/items/${item.id}`))
        .catch(err => {
          if (err instanceof ApiClientError && err.status === 404) {
            setState({ status: 'not_found', barcode })
          } else {
            setState({ status: 'error', message: t('common.unknown_error') })
          }
        })
    }
  }

  function handleReset() {
    setState({ status: 'idle' })
    setManual('')
    setShouldScan(false)
    setTimeout(() => setShouldScan(true), 100)
  }

  function handleCreateItem() {
    if (state.status === 'not_found') {
      setScannedBarcode(state.barcode)
      setCreateDialogOpen(true)
    }
  }

  function handleCreateDialogClose() {
    setCreateDialogOpen(false)
    setScannedBarcode('')
    handleReset()
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
                <div className="flex gap-2 mt-4">
                  <button
                    onClick={handleCreateItem}
                    className="rounded-md bg-primary text-primary-foreground px-4 py-2 text-sm font-medium hover:bg-primary/90 transition-colors"
                  >
                    {t('scanner.create_item')}
                  </button>
                  <button
                    onClick={handleReset}
                    className="rounded-md bg-white/20 px-4 py-2 text-sm hover:bg-white/30 transition-colors"
                  >
                    {t('scanner.scan_again')}
                  </button>
                </div>
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

      <CreateItemDialog
        open={createDialogOpen}
        onClose={handleCreateDialogClose}
        categories={categories}
        buildings={buildings}
        initialBarcode={scannedBarcode}
      />
    </div>
  )
}
