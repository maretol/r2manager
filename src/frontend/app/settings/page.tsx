import { SettingsIcon } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { fetchBuckets, fetchAllBucketSettings } from '@/lib/api'
import { SettingsForm } from '@/components/settings-form'

export default async function SettingsPage() {
  const buckets = await fetchBuckets()
  const settings = await fetchAllBucketSettings()

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
          <SettingsForm buckets={buckets} initialSettings={settings} />
        </CardContent>
      </Card>
    </div>
  )
}
