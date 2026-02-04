import Link from 'next/link'
import { FolderIcon, SettingsIcon, ArrowRightIcon } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'

export default function Home() {
  return (
    <div className="flex min-h-screen items-center justify-center p-6">
      <Card className="w-full max-w-lg">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">R2 Contents Manager</CardTitle>
          <CardDescription>
            Manage your Cloudflare R2 storage buckets
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-start gap-4 rounded-lg border p-4">
            <FolderIcon className="size-8 text-muted-foreground mt-1" />
            <div className="flex-1">
              <h3 className="font-medium">Select a Bucket</h3>
              <p className="text-sm text-muted-foreground">
                Choose a bucket from the sidebar to browse its contents.
              </p>
            </div>
          </div>
          <div className="flex items-start gap-4 rounded-lg border p-4">
            <SettingsIcon className="size-8 text-muted-foreground mt-1" />
            <div className="flex-1">
              <h3 className="font-medium">Configure Settings</h3>
              <p className="text-sm text-muted-foreground mb-2">
                Set up your R2 credentials and preferences.
              </p>
              <Button asChild variant="outline" size="sm">
                <Link href="/settings">
                  Go to Settings
                  <ArrowRightIcon className="size-4 ml-1" />
                </Link>
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
