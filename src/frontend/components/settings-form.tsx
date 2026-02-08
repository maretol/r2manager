'use client'

import { useState } from 'react'
import { SaveIcon, Loader2Icon } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { bulkUpdateBucketSettings, type Bucket, type BucketSettingsResponse } from '@/lib/api'

type SettingsFormProps = {
  buckets: Bucket[]
  initialSettings: BucketSettingsResponse[]
}

export function SettingsForm({ buckets, initialSettings }: SettingsFormProps) {
  const [publicUrls, setPublicUrls] = useState<Record<string, string>>(() => {
    const map: Record<string, string> = {}
    for (const s of initialSettings) {
      map[s.bucket_name] = s.public_url
    }
    return map
  })
  const [saving, setSaving] = useState(false)
  const [saveMessage, setSaveMessage] = useState<string | null>(null)

  const handlePublicUrlChange = (bucketName: string, value: string) => {
    setPublicUrls((prev) => ({ ...prev, [bucketName]: value }))
  }

  const handleSave = async () => {
    setSaving(true)
    setSaveMessage(null)

    try {
      const settings = buckets.map((bucket) => ({
        bucket_name: bucket.name,
        public_url: publicUrls[bucket.name] || '',
      }))
      await bulkUpdateBucketSettings(settings)
      setSaveMessage('Settings saved successfully')
      setTimeout(() => setSaveMessage(null), 3000)
    } catch {
      setSaveMessage('Failed to save settings')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-medium mb-4">Bucket Public URLs</h3>
        <p className="text-sm text-muted-foreground mb-4">
          Set the public URL for each bucket. This URL will be used to generate public links for objects.
        </p>

        {buckets.length === 0 ? (
          <p className="text-muted-foreground">No buckets found</p>
        ) : (
          <div className="space-y-4">
            {buckets.map((bucket) => (
              <div key={bucket.name} className="flex flex-col gap-2">
                <label htmlFor={`url-${bucket.name}`} className="text-sm font-medium">
                  {bucket.name}
                </label>
                <Input
                  id={`url-${bucket.name}`}
                  type="url"
                  placeholder="https://example.com"
                  value={publicUrls[bucket.name] || ''}
                  onChange={(e) => handlePublicUrlChange(bucket.name, e.target.value)}
                />
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="flex items-center gap-4">
        <Button onClick={handleSave} disabled={saving} className="cursor-pointer">
          {saving ? (
            <Loader2Icon className="size-4 animate-spin mr-2 pointer-events-auto" />
          ) : (
            <SaveIcon className="size-4 mr-2" />
          )}
          Save Settings
        </Button>
        {saveMessage && (
          <span className={`text-sm ${saveMessage.includes('success') ? 'text-green-600' : 'text-destructive'}`}>
            {saveMessage}
          </span>
        )}
      </div>
    </div>
  )
}
