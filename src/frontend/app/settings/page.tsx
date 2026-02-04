'use client'

import { useEffect, useState } from 'react'
import { SettingsIcon, SaveIcon, Loader2Icon } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { fetchBuckets, type Bucket } from '@/lib/api'
import { loadSettings, saveSettings, type AppSettings } from '@/lib/settings'

export default function SettingsPage() {
  const [buckets, setBuckets] = useState<Bucket[]>([])
  const [settings, setSettings] = useState<AppSettings | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [saveMessage, setSaveMessage] = useState<string | null>(null)

  useEffect(() => {
    async function load() {
      try {
        const [bucketsData, settingsData] = await Promise.all([fetchBuckets(), Promise.resolve(loadSettings())])
        setBuckets(bucketsData)
        setSettings(settingsData)
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load data')
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  const handlePublicUrlChange = (bucketName: string, publicUrl: string) => {
    if (!settings) return

    setSettings({
      ...settings,
      buckets: {
        ...settings.buckets,
        [bucketName]: {
          ...settings.buckets[bucketName],
          publicUrl,
        },
      },
    })
  }

  const handleSave = () => {
    if (!settings) return

    setSaving(true)
    setSaveMessage(null)

    try {
      saveSettings(settings)
      setSaveMessage('Settings saved successfully')
      setTimeout(() => setSaveMessage(null), 3000)
    } catch {
      setSaveMessage('Failed to save settings')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="flex flex-col gap-4 p-6">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <SettingsIcon className="size-5" />
              Settings
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-center py-12">
              <Loader2Icon className="size-6 animate-spin text-muted-foreground" />
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col gap-4 p-6">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <SettingsIcon className="size-5" />
              Settings
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-center py-12 text-destructive">
              <p>{error}</p>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-4 p-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <SettingsIcon className="size-5" />
            Settings
          </CardTitle>
          <CardDescription>Configure your R2 Contents Manager settings</CardDescription>
        </CardHeader>
        <CardContent>
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
                        value={settings?.buckets[bucket.name]?.publicUrl || ''}
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
        </CardContent>
      </Card>
    </div>
  )
}
