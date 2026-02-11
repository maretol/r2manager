'use client'

import { useState, useCallback, useTransition } from 'react'
import type { UploadFileEntry, UploadFileStatus } from '@/types/upload'
import { uploadFile } from '@/app/bucket/[name]/actions'

type UseFileUploadOptions = {
  bucketName: string
  prefix: string
  onAllCompleteAction?: (allSucceeded: boolean) => void
}

export function useFileUpload({ bucketName, prefix, onAllCompleteAction }: UseFileUploadOptions) {
  const [files, setFiles] = useState<UploadFileEntry[]>([])
  const [isPending, startTransition] = useTransition()

  const addFiles = useCallback((newFiles: File[]) => {
    const entries: UploadFileEntry[] = newFiles.map((file) => ({
      id: crypto.randomUUID(),
      file,
      status: 'pending' as UploadFileStatus,
    }))
    setFiles((prev) => [...prev, ...entries])
  }, [])

  const removeFile = useCallback((id: string) => {
    setFiles((prev) => prev.filter((f) => f.id !== id))
  }, [])

  const clearFiles = useCallback(() => {
    setFiles([])
  }, [])

  const updateFile = useCallback((id: string, updates: Partial<UploadFileEntry>) => {
    setFiles((prev) => prev.map((f) => (f.id === id ? { ...f, ...updates } : f)))
  }, [])

  const startUpload = useCallback(
    (overwrite: boolean) => {
      const pendingFiles = files.filter((f) => f.status === 'pending' || f.status === 'conflict')

      startTransition(async () => {
        let hasAnyCompleted = false
        let allSucceeded = true

        for (const entry of pendingFiles) {
          const key = prefix + entry.file.name

          updateFile(entry.id, { status: 'uploading' })

          try {
            const formData = new FormData()
            formData.append('file', entry.file)
            formData.append('bucketName', bucketName)
            formData.append('key', key)
            if (overwrite) {
              formData.append('overwrite', 'true')
            }

            const result = await uploadFile(formData)

            if (result.success) {
              updateFile(entry.id, { status: 'done' })
              hasAnyCompleted = true
            } else if (result.code === 'CONFLICT') {
              allSucceeded = false
              updateFile(entry.id, {
                status: 'conflict',
                error: '同名のファイルが既に存在します',
              })
            } else {
              allSucceeded = false
              updateFile(entry.id, { status: 'error', error: result.error })
            }
          } catch (err) {
            allSucceeded = false
            console.error('Upload failed:', err)
            updateFile(entry.id, { status: 'error', error: 'アップロードに失敗しました' })
          }
        }

        if (hasAnyCompleted) {
          onAllCompleteAction?.(allSucceeded)
        }
      })
    },
    [files, bucketName, prefix, updateFile, onAllCompleteAction],
  )

  return {
    files,
    isPending,
    addFiles,
    removeFile,
    clearFiles,
    startUpload,
  }
}
