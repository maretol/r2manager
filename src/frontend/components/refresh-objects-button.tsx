'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { RefreshCwIcon, Loader2Icon } from 'lucide-react'
import { clearObjectsCache } from '@/lib/api'
import { Button } from './ui/button'

type RefreshObjectsButtonProps = {
  bucketName: string
}

export function RefreshObjectsButton({ bucketName }: RefreshObjectsButtonProps) {
  const router = useRouter()
  const [loading, setLoading] = useState(false)

  const handleRefresh = async () => {
    setLoading(true)
    try {
      await clearObjectsCache(bucketName)
      router.refresh()
    } catch (error) {
      console.error('Failed to refresh objects:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Button
      variant="ghost"
      size="icon"
      onClick={handleRefresh}
      disabled={loading}
      title="Refresh objects"
      className="cursor-pointer"
    >
      {loading ? (
        <Loader2Icon className="size-4 animate-spin" />
      ) : (
        <RefreshCwIcon className="size-4" />
      )}
    </Button>
  )
}
