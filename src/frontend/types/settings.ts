export type BucketSettings = {
  publicUrl: string
}

export type AppSettings = {
  buckets: Record<string, BucketSettings>
}

export const defaultAppSettings: AppSettings = {
  buckets: {},
}
