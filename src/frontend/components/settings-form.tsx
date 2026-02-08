'use client'

import { useActionState, useEffect, useState } from 'react'
import { SaveIcon, Loader2Icon } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { type Bucket, type BucketSettingsResponse } from '@/lib/api'
import { saveBucketSettings } from '@/app/settings/actions'

type SettingsFormProps = {
  buckets: Bucket[]
  initialSettings: BucketSettingsResponse[]
}

export function SettingsForm({ buckets, initialSettings }: SettingsFormProps) {
  const [state, formAction, isPending] = useActionState(saveBucketSettings, null)
  const [lastState, setLastState] = useState(state)
  const [showMessage, setShowMessage] = useState(false)
  const [publicUrls, setPublicUrls] = useState<Record<string, string>>(() => {
    const map: Record<string, string> = {}
    for (const s of initialSettings) {
      map[s.bucket_name] = s.public_url
    }
    return map
  })

  if (state !== lastState) {
    setLastState(state)
    setShowMessage(state !== null)
  }

  useEffect(() => {
    if (!showMessage) return
    const timer = setTimeout(() => setShowMessage(false), 3000)
    return () => clearTimeout(timer)
  }, [showMessage])

  return (
    <form action={formAction} className="space-y-6">
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
                <input type="hidden" name="bucket_name" value={bucket.name} />
                <label htmlFor={`url-${bucket.name}`} className="text-sm font-medium">
                  {bucket.name}
                </label>
                <Input
                  id={`url-${bucket.name}`}
                  name={`public_url:${bucket.name}`}
                  type="url"
                  placeholder="https://example.com"
                  value={publicUrls[bucket.name] || ''}
                  onChange={(e) => setPublicUrls((prev) => ({ ...prev, [bucket.name]: e.target.value }))}
                />
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="flex items-center gap-4">
        <Button type="submit" disabled={isPending} className="cursor-pointer">
          {isPending ? (
            <Loader2Icon className="size-4 animate-spin mr-2 pointer-events-auto" />
          ) : (
            <SaveIcon className="size-4 mr-2" />
          )}
          Save Settings
        </Button>
        {showMessage && state && (
          <span className={`text-sm ${state.success ? 'text-green-600' : 'text-destructive'}`}>
            {state.message}
          </span>
        )}
      </div>
    </form>
  )
}
