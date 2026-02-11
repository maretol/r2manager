'use client'

import { useState, useRef, useCallback, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { Upload, X, FileIcon, CheckCircle2, AlertCircle, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { useFileUpload } from '@/hooks/use-file-upload'
import { clearObjectsCache } from '@/lib/api'
import { formatFileSize } from '@/lib/object-utils'
import type { UploadFileEntry } from '@/types/upload'

type UploadDialogProps = {
  bucketName: string
  prefix: string
}

export function UploadDialog({ bucketName, prefix }: UploadDialogProps) {
  const router = useRouter()
  const [open, setOpen] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [isDragOver, setIsDragOver] = useState(false)
  const [allowOverwrite, setAllowOverwrite] = useState(false)
  const [closingIn, setClosingIn] = useState<number | null>(null)
  const autoCloseTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const countdownRef = useRef<ReturnType<typeof setInterval> | null>(null)
  const clearFilesRef = useRef<() => void>(() => {})

  const clearAutoClose = useCallback(() => {
    if (autoCloseTimerRef.current) {
      clearTimeout(autoCloseTimerRef.current)
      autoCloseTimerRef.current = null
    }
    if (countdownRef.current) {
      clearInterval(countdownRef.current)
      countdownRef.current = null
    }
    setClosingIn(null)
  }, [])

  const { files, isPending, addFiles, removeFile, clearFiles, startUpload } = useFileUpload({
    bucketName,
    prefix,
    onAllCompleteAction: async (allSucceeded) => {
      await clearObjectsCache(bucketName)
      router.refresh()
      if (allSucceeded) {
        setClosingIn(3)
        countdownRef.current = setInterval(() => {
          setClosingIn((prev) => (prev !== null && prev > 1 ? prev - 1 : prev))
        }, 1000)
        autoCloseTimerRef.current = setTimeout(() => {
          clearAutoClose()
          clearFilesRef.current()
          setAllowOverwrite(false)
          setOpen(false)
        }, 3000)
      }
    },
  })

  useEffect(() => {
    clearFilesRef.current = clearFiles
  }, [clearFiles])

  const handleClose = () => {
    if (!isPending) {
      clearAutoClose()
      clearFiles()
      setAllowOverwrite(false)
      setOpen(false)
    }
  }

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(true)
  }, [])

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(false)
  }, [])

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault()
      setIsDragOver(false)
      const droppedFiles = Array.from(e.dataTransfer.files)
      if (droppedFiles.length > 0) {
        addFiles(droppedFiles)
      }
    },
    [addFiles],
  )

  const handleFileSelect = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const selected = Array.from(e.target.files ?? [])
      if (selected.length > 0) {
        addFiles(selected)
      }
      e.target.value = ''
    },
    [addFiles],
  )

  const pendingCount = files.filter((f) => f.status === 'pending').length
  const conflictCount = files.filter((f) => f.status === 'conflict').length
  const canUpload = (pendingCount > 0 || conflictCount > 0) && !isPending

  return (
    <Dialog open={open} onOpenChange={(v) => (v ? setOpen(true) : handleClose())}>
      <DialogTrigger asChild>
        <Button variant="outline" className="cursor-pointer">
          <Upload className="size-4 mr-1" />
          アップロード
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-xl" showCloseButton={!isPending}>
        <DialogHeader>
          <DialogTitle>ファイルアップロード</DialogTitle>
          <DialogDescription>アップロード先: {prefix || '/'}</DialogDescription>
        </DialogHeader>

        {/* ドロップゾーン */}
        <div
          className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors cursor-pointer ${
            isDragOver ? 'border-primary bg-primary/5' : 'border-muted-foreground/25'
          }`}
          onDragOver={handleDragOver}
          onDragEnter={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          onClick={() => fileInputRef.current?.click()}
        >
          <Upload className="size-8 mx-auto mb-2 text-muted-foreground" />
          <p className="text-sm text-muted-foreground">ファイルをドラッグ＆ドロップ、またはクリックして選択</p>
          <input ref={fileInputRef} type="file" multiple className="hidden" onChange={handleFileSelect} />
        </div>

        {/* 上書き許可チェックボックス */}
        <label className="flex items-center gap-2 cursor-pointer">
          <Checkbox
            checked={allowOverwrite}
            onCheckedChange={(checked) => setAllowOverwrite(checked === true)}
            disabled={isPending}
          />
          <span className="text-sm">同名ファイルの上書きを許可する</span>
        </label>

        {/* ファイルリスト */}
        {files.length > 0 && (
          <div className="max-h-60 overflow-y-auto space-y-2">
            {files.map((entry) => (
              <UploadFileItem key={entry.id} entry={entry} onRemove={removeFile} isPending={isPending} />
            ))}
          </div>
        )}

        <DialogFooter>
          {closingIn !== null && <p className="text-sm text-muted-foreground mr-auto">{closingIn}秒後に閉じます...</p>}
          <Button variant="outline" onClick={handleClose} disabled={isPending}>
            閉じる
          </Button>
          <Button onClick={() => startUpload(allowOverwrite)} disabled={!canUpload} className="cursor-pointer">
            {isPending ? (
              <>
                <Loader2 className="size-4 animate-spin mr-1" />
                アップロード中...
              </>
            ) : (
              <>アップロード</>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

function UploadFileItem({
  entry,
  onRemove,
  isPending,
}: {
  entry: UploadFileEntry
  onRemove: (id: string) => void
  isPending: boolean
}) {
  return (
    <div className="flex items-center gap-2 rounded-md border p-2">
      <div className="shrink-0">
        {entry.status === 'done' ? (
          <CheckCircle2 className="size-4 text-green-500" />
        ) : entry.status === 'error' ? (
          <AlertCircle className="size-4 text-destructive" />
        ) : entry.status === 'conflict' ? (
          <AlertCircle className="size-4 text-orange-500" />
        ) : entry.status === 'uploading' ? (
          <Loader2 className="size-4 animate-spin text-primary" />
        ) : (
          <FileIcon className="size-4 text-muted-foreground" />
        )}
      </div>

      <div className="flex-1 min-w-0">
        <div className="flex items-center justify-between gap-2">
          <span className="text-sm truncate">{entry.file.name}</span>
          <span className="text-xs text-muted-foreground shrink-0">{formatFileSize(entry.file.size)}</span>
        </div>

        {entry.status === 'uploading' && <p className="text-xs text-muted-foreground mt-1">アップロード中...</p>}

        {entry.status === 'error' && <p className="text-xs text-destructive mt-1">{entry.error}</p>}

        {entry.status === 'conflict' && <p className="text-xs text-orange-500 mt-1">{entry.error}</p>}
      </div>

      {entry.status === 'pending' && !isPending && (
        <Button
          variant="ghost"
          size="icon"
          className="size-6 shrink-0 cursor-pointer"
          onClick={() => onRemove(entry.id)}
        >
          <X className="size-3" />
        </Button>
      )}
    </div>
  )
}
