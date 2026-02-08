'use server'

import { bulkUpdateBucketSettings } from '@/lib/api'

export type SaveSettingsState = {
  success: boolean
  message: string
} | null

export async function saveBucketSettings(_prevState: SaveSettingsState, formData: FormData): Promise<SaveSettingsState> {
  const bucketNames = formData.getAll('bucket_name') as string[]

  const settings = bucketNames.map((name) => ({
    bucket_name: name,
    public_url: (formData.get(`public_url:${name}`) as string) || '',
  }))

  try {
    await bulkUpdateBucketSettings(settings)
    return { success: true, message: 'Settings saved successfully' }
  } catch (err) {
    return { success: false, message: err instanceof Error ? err.message : 'Failed to save settings' }
  }
}
