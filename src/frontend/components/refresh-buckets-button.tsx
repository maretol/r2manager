'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { RefreshCwIcon, Loader2Icon } from 'lucide-react'
import { clearBucketsCache } from '@/lib/api'

export function RefreshBucketsButton() {
  const router = useRouter()
  const [loading, setLoading] = useState(false)

  const handleRefresh = async () => {
    setLoading(true)
    try {
      await clearBucketsCache()
      router.refresh()
    } catch (error) {
      console.error('Failed to refresh buckets:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <button
      onClick={handleRefresh}
      disabled={loading}
      title="Refresh buckets"
      className="cursor-pointer"
    >
      {loading ? (
        <Loader2Icon className="size-4 animate-spin" />
      ) : (
        <RefreshCwIcon className="size-4" />
      )}
    </button>
  )
}
