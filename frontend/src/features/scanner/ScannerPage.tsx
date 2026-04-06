import { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { useQuery } from '@tanstack/react-query'
import { Html5Qrcode } from 'html5-qrcode'
import { Flashlight, FlashlightOff } from 'lucide-react'
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
  const scannerRef = useRef<Html5Qrcode | null>(null)
  const scannerDivRef = useRef<HTMLDivElement>(null)
  const videoTrackRef = useRef<MediaStreamTrack | null>(null)
  const torchOnRef = useRef(false)
  const [state, setState] = useState<ScanState>({ status: 'idle' })
  const [manual, setManual] = useState('')
  const [showManual, setShowManual] = useState(false)
  const [shouldScan, setShouldScan] = useState(true)
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [scannedBarcode, setScannedBarcode] = useState('')
  const [torchOn, setTorchOn] = useState(false)
  const [torchSupported, setTorchSupported] = useState(false)

  const { data: categories = [] } = useQuery({
    queryKey: ['categories'],
    queryFn: categoriesApi.list,
  })

  const { data: buildings = [] } = useQuery({
    queryKey: ['buildings'],
    queryFn: buildingsApi.list,
  })

  useEffect(() => {
    if (!shouldScan || !scannerDivRef.current) return

    const scannerId = 'barcode-scanner'
    const scanner = new Html5Qrcode(scannerId)
    scannerRef.current = scanner
    let isProcessing = false

    async function startScanner() {
      try {
        setState({ status: 'scanning' })
        
        await scanner.start(
          { 
            facingMode: 'environment'
          },
          {
            fps: 10,
            qrbox: { width: 300, height: 200 },
            aspectRatio: 1.777778,
            disableFlip: false
          },
          async (decodedText) => {
            if (isProcessing) return
            isProcessing = true
            
            setState({ status: 'found', barcode: decodedText })
            
            try {
              // Turn off torch BEFORE stopping scanner using ref
              if (torchOnRef.current && videoTrackRef.current) {
                await videoTrackRef.current.applyConstraints({
                  advanced: [{ torch: false } as any]
                }).catch(console.error)
                torchOnRef.current = false
                setTorchOn(false)
              }
              
              await scanner.stop()
              const item = await itemsApi.getByBarcode(decodedText)
              navigate(`/items/${item.id}`)
            } catch (err) {
              if (err instanceof ApiClientError && err.status === 404) {
                setState({ status: 'not_found', barcode: decodedText })
              } else {
                setState({ status: 'error', message: t('common.unknown_error') })
              }
            }
          },
          undefined
        )
        
        // Try to enable autofocus and check torch support after scanner starts
        const videoElement = document.querySelector('#barcode-scanner video') as HTMLVideoElement
        if (videoElement && videoElement.srcObject) {
          const stream = videoElement.srcObject as MediaStream
          const videoTrack = stream.getVideoTracks()[0]
          if (videoTrack) {
            videoTrackRef.current = videoTrack
            
            const capabilities = videoTrack.getCapabilities() as any
            
            // Check torch support
            if (capabilities.torch) {
              setTorchSupported(true)
            }
            
            // Enable autofocus
            if (capabilities.focusMode && capabilities.focusMode.includes('continuous')) {
              await videoTrack.applyConstraints({
                advanced: [{ focusMode: 'continuous' } as any]
              }).catch(() => {
                // Fallback to manual focus if continuous not supported
                videoTrack.applyConstraints({
                  advanced: [{ focusMode: 'manual' } as any]
                }).catch(console.error)
              })
            }
          }
        }
      } catch (err) {
        console.error('Scanner error:', err)
        setState({ status: 'error', message: t('scanner.camera_error') })
      }
    }

    startScanner()
    
    return () => {
      if (scanner.isScanning) {
        scanner.stop().catch(console.error)
      }
      videoTrackRef.current = null
      torchOnRef.current = false
      setTorchOn(false)
      setTorchSupported(false)
    }
  }, [shouldScan, navigate, t])

  function handleManualSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (manual.trim()) {
      setShouldScan(false)
      
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
      // Turn off torch before opening create dialog
      if (torchOnRef.current && videoTrackRef.current) {
        videoTrackRef.current.applyConstraints({
          advanced: [{ torch: false } as any]
        }).catch(console.error)
        torchOnRef.current = false
        setTorchOn(false)
      }
      setScannedBarcode(state.barcode)
      setCreateDialogOpen(true)
    }
  }

  function handleCreateDialogClose() {
    setCreateDialogOpen(false)
    setScannedBarcode('')
    handleReset()
  }

  async function toggleTorch() {
    if (!videoTrackRef.current || !torchSupported) return
    
    try {
      const newTorchState = !torchOn
      await videoTrackRef.current.applyConstraints({
        advanced: [{ torch: newTorchState } as any]
      })
      torchOnRef.current = newTorchState
      setTorchOn(newTorchState)
    } catch (err) {
      console.error('Torch toggle error:', err)
    }
  }

  return (
    <div className="flex flex-col h-full">
      <div className="p-4 border-b">
        <h1 className="text-lg font-semibold">{t('scanner.title')}</h1>
      </div>

      <div className="relative flex-1 bg-black overflow-hidden">
        <div 
          id="barcode-scanner" 
          ref={scannerDivRef}
          className="w-full h-full"
        />

        {state.status === 'scanning' && (
          <>
            <div className="absolute bottom-8 left-0 right-0 text-center pointer-events-none">
              <p className="text-white text-sm bg-black/50 inline-block px-4 py-2 rounded-full">
                {t('scanner.instruction')}
              </p>
            </div>
            
            {torchSupported && (
              <button
                onClick={toggleTorch}
                className="absolute top-4 right-4 p-3 rounded-full bg-black/50 text-white hover:bg-black/70 transition-colors"
                aria-label={torchOn ? 'Turn off flashlight' : 'Turn on flashlight'}
              >
                {torchOn ? <FlashlightOff size={24} /> : <Flashlight size={24} />}
              </button>
            )}
          </>
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
