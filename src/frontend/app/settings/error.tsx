'use client'

import { useEffect } from 'react'
import { AlertCircleIcon, RefreshCwIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

type ErrorProps = {
  error: Error & { digest?: string }
  reset: () => void
}

export default function SettingsError({ error, reset }: ErrorProps) {
  useEffect(() => {
    console.error('Settings page error:', error)
  }, [error])

  return (
    <div className="flex flex-col items-center justify-center p-6">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-destructive">
            <AlertCircleIcon className="size-5" />
            Error Loading Settings
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <p className="text-muted-foreground">
            {error.message || 'Failed to load settings. Please try again.'}
          </p>
          <Button onClick={reset} variant="outline" className="w-full cursor-pointer">
            <RefreshCwIcon className="size-4 mr-2" />
            Try Again
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}
