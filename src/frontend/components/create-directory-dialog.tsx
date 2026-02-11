'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { FolderPlus, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { createDirectoryAction } from '@/app/bucket/[name]/actions'
import { clearObjectsCache } from '@/lib/api'

type CreateDirectoryDialogProps = {
  bucketName: string
  prefix: string
}

export function CreateDirectoryDialog({ bucketName, prefix }: CreateDirectoryDialogProps) {
  const router = useRouter()
  const [open, setOpen] = useState(false)
  const [dirName, setDirName] = useState('')
  const [isCreating, setIsCreating] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleClose = () => {
    if (!isCreating) {
      setDirName('')
      setError(null)
      setOpen(false)
    }
  }

  const handleCreate = async () => {
    const trimmed = dirName.trim()
    if (!trimmed) {
      setError('フォルダ名を入力してください')
      return
    }

    if (trimmed.includes('/') || trimmed.includes('\\')) {
      setError('フォルダ名にスラッシュは使用できません')
      return
    }

    if (trimmed === '.' || trimmed === '..') {
      setError('無効なフォルダ名です')
      return
    }

    setIsCreating(true)
    setError(null)

    try {
      const path = prefix + trimmed
      const result = await createDirectoryAction(bucketName, path)

      if (result.success) {
        await clearObjectsCache(bucketName)
        router.refresh()
        setDirName('')
        setOpen(false)
      } else {
        setError(result.error)
      }
    } catch (err) {
      console.error('Failed to create directory:', err)
      setError('ディレクトリの作成に失敗しました')
    } finally {
      setIsCreating(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={(v) => (v ? setOpen(true) : handleClose())}>
      <DialogTrigger asChild>
        <Button variant="outline" className="cursor-pointer">
          <FolderPlus className="size-4 mr-1" />
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>フォルダ作成</DialogTitle>
          <DialogDescription>作成先: {prefix || '/'}</DialogDescription>
        </DialogHeader>

        <div className="space-y-2">
          <Input
            placeholder="フォルダ名"
            value={dirName}
            onChange={(e) => {
              setDirName(e.target.value)
              setError(null)
            }}
            onKeyDown={(e) => {
              if (e.key === 'Enter' && !isCreating) {
                handleCreate()
              }
            }}
            disabled={isCreating}
          />
          {error && <p className="text-sm text-destructive">{error}</p>}
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={handleClose} disabled={isCreating}>
            キャンセル
          </Button>
          <Button onClick={handleCreate} disabled={isCreating || !dirName.trim()} className="cursor-pointer">
            {isCreating ? (
              <>
                <Loader2 className="size-4 animate-spin mr-1" />
                作成中...
              </>
            ) : (
              '作成'
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
