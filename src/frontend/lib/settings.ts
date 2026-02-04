import type { AppSettings, BucketSettings } from '@/types/settings'
import { defaultAppSettings } from '@/types/settings'

export type { AppSettings, BucketSettings }

const SETTINGS_KEY = 'r2manager-settings'

export function loadSettings(): AppSettings {
  if (typeof window === 'undefined') {
    return defaultAppSettings
  }

  try {
    const stored = localStorage.getItem(SETTINGS_KEY)
    if (!stored) {
      return defaultAppSettings
    }
    return JSON.parse(stored) as AppSettings
  } catch {
    return defaultAppSettings
  }
}

export function saveSettings(settings: AppSettings): void {
  if (typeof window === 'undefined') {
    return
  }

  localStorage.setItem(SETTINGS_KEY, JSON.stringify(settings))
}

export function getBucketSettings(bucketName: string): BucketSettings | undefined {
  const settings = loadSettings()
  return settings.buckets[bucketName]
}

export function setBucketSettings(bucketName: string, bucketSettings: BucketSettings): void {
  const settings = loadSettings()
  settings.buckets[bucketName] = bucketSettings
  saveSettings(settings)
}

export function getPublicUrl(bucketName: string): string {
  const bucketSettings = getBucketSettings(bucketName)
  return bucketSettings?.publicUrl || ''
}

export function setPublicUrl(bucketName: string, publicUrl: string): void {
  const current = getBucketSettings(bucketName) || { publicUrl: '' }
  setBucketSettings(bucketName, { ...current, publicUrl })
}
