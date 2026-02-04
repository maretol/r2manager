export type R2Object = {
  key: string
  size: number
  last_modified: string
  etag: string
}

export type ListObjectsResponse = {
  objects: R2Object[]
  prefix: string
  delimiter: string
  is_truncated: boolean
  next_continuation_token?: string
}

export type DisplayObject = {
  name: string
  key: string
  isFolder: boolean
  size?: number
  lastModified?: string
}
