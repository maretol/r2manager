'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Download, Copy, Link, X, Check, Trash2, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import type { DisplayObject } from '@/types/object'
import { formatFileSize, formatDate } from '@/lib/object-utils'
import { getPublicUrl } from '@/lib/settings'
import { clearContentCache } from '@/lib/api'
import Image from 'next/image'

type ObjectDetailPanelProps = {
  object: DisplayObject
  bucketName: string
  prefix: string
}

const IMAGE_EXTENSIONS = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.svg', '.bmp', '.ico']

function isImageFile(filename: string): boolean {
  const lowerName = filename.toLowerCase()
  return IMAGE_EXTENSIONS.some((ext) => lowerName.endsWith(ext))
}

function getObjectUrl(bucketName: string, key: string): string {
  return `http://localhost:3000/api/v1/buckets/${encodeURIComponent(bucketName)}/content/${encodeURIComponent(key)}`
}

function getPublicObjectUrl(publicBaseUrl: string, key: string): string {
  const baseUrl = publicBaseUrl.endsWith('/') ? publicBaseUrl.slice(0, -1) : publicBaseUrl
  return `${baseUrl}/${key}`
}

export function ObjectDetailPanel({ object, bucketName, prefix }: ObjectDetailPanelProps) {
  const router = useRouter()
  const [copiedUrl, setCopiedUrl] = useState<'internal' | 'public' | null>(null)
  const [imageError, setImageError] = useState(false)
  const [clearingCache, setClearingCache] = useState(false)
  const [cacheCleared, setCacheCleared] = useState(false)

  const publicUrl = getPublicUrl(bucketName)
  const hasPublicUrl = publicUrl.length > 0
  const isImage = isImageFile(object.name)
  const objectUrl = getObjectUrl(bucketName, object.key)
  const publicObjectUrl = hasPublicUrl ? getPublicObjectUrl(publicUrl, object.key) : null

  const handleClose = () => {
    const basePath = `/bucket/${encodeURIComponent(bucketName)}`
    const params = new URLSearchParams()
    if (prefix) params.set('prefix', prefix)
    const query = params.toString()
    router.push(`${basePath}${query ? `?${query}` : ''}`)
  }

  const handleCopyUrl = async (url: string, type: 'internal' | 'public') => {
    try {
      await navigator.clipboard.writeText(url)
      setCopiedUrl(type)
      setTimeout(() => setCopiedUrl(null), 2000)
    } catch {
      console.error('Failed to copy URL')
    }
  }

  const handleDownload = () => {
    const link = document.createElement('a')
    link.href = objectUrl
    link.download = object.name
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  const handleClearCache = async () => {
    setClearingCache(true)
    setCacheCleared(false)
    try {
      await clearContentCache(bucketName, object.key)
      setCacheCleared(true)
      setTimeout(() => setCacheCleared(false), 2000)
    } catch (error) {
      console.error('Failed to clear cache:', error)
    } finally {
      setClearingCache(false)
    }
  }

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between p-4 border-b">
        <h3 className="font-semibold truncate">Details</h3>
        <Button variant="ghost" size="icon-xs" onClick={handleClose}>
          <X className="size-4" />
        </Button>
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {isImage && !imageError && (
          <div className="rounded-lg border bg-muted/30 p-2">
            <Image
              src={objectUrl}
              alt={object.name}
              className="w-full h-auto max-h-48 object-contain rounded"
              onError={() => setImageError(true)}
              width={'100'}
              height={'200'}
            />
          </div>
        )}

        <div className="space-y-3">
          <DetailRow label="Name" value={object.name} />
          <DetailRow label="Key" value={object.key} className="break-all" />
          <DetailRow label="Size" value={formatFileSize(object.size ?? 0)} />
          <DetailRow label="Last Modified" value={object.lastModified ? formatDate(object.lastModified) : '-'} />
          {object.etag && <DetailRow label="ETag" value={object.etag} className="break-all" />}
        </div>

        <Separator />

        <div className="space-y-2">
          <Button variant="outline" size="sm" className="w-full justify-start cursor-pointer" onClick={handleDownload}>
            <Download className="size-4" />
            Download
          </Button>

          <Button
            variant="outline"
            size="sm"
            className="w-full justify-start cursor-pointer"
            onClick={() => handleCopyUrl(objectUrl, 'internal')}
          >
            {copiedUrl === 'internal' ? <Check className="size-4" /> : <Copy className="size-4" />}
            {copiedUrl === 'internal' ? 'Copied!' : 'Copy URL'}
          </Button>

          {hasPublicUrl && publicObjectUrl && (
            <Button
              variant="outline"
              size="sm"
              className="w-full justify-start cursor-pointer"
              onClick={() => handleCopyUrl(publicObjectUrl, 'public')}
            >
              {copiedUrl === 'public' ? <Check className="size-4" /> : <Link className="size-4" />}
              {copiedUrl === 'public' ? 'Copied!' : 'Copy Public URL'}
            </Button>
          )}

          <Separator />

          <Button
            variant="outline"
            size="sm"
            className="w-full justify-start cursor-pointer"
            onClick={handleClearCache}
            disabled={clearingCache}
          >
            {clearingCache ? (
              <Loader2 className="size-4 animate-spin" />
            ) : cacheCleared ? (
              <Check className="size-4" />
            ) : (
              <Trash2 className="size-4" />
            )}
            {cacheCleared ? 'Cache Cleared!' : 'Clear Cache'}
          </Button>
        </div>
      </div>
    </div>
  )
}

type DetailRowProps = {
  label: string
  value: string
  className?: string
}

function DetailRow({ label, value, className }: DetailRowProps) {
  return (
    <div>
      <dt className="text-xs text-muted-foreground">{label}</dt>
      <dd className={`text-sm mt-0.5 ${className || ''}`}>{value}</dd>
    </div>
  )
}
