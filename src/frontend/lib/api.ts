export type Bucket = {
  name: string
  creation_date: string
}

type GetBucketsResponse = {
  buckets: Bucket[]
}

type ErrorResponse = {
  error: string
}

// サーバサイドからの呼び出しでは自分のサーバのURLを使う必要がある
const HOST = process.env.HOSTNAME
const PROTOCOL = process.env.PROTOCOL || 'http'
const PORT = process.env.PORT ? `:${process.env.PORT}` : ':3000'
// HOSTが設定されている = サーバサイドからの呼び出しとして扱う
// クライアントサイドからの呼び出しでは相対パスで問題ないので空文字にする
const SERVER_URL = HOST ? PROTOCOL + '://' + (HOST + PORT) : ''

export async function fetchBuckets(): Promise<Bucket[]> {
  console.log('Fetching buckets from API...' + SERVER_URL)
  const response = await fetch(`${SERVER_URL}/api/v1/buckets`)
  if (!response.ok) {
    const errorData: ErrorResponse = await response.json()
    console.log(response.status)
    throw new Error(errorData.error || 'Failed to fetch buckets')
  }
  const data: GetBucketsResponse = await response.json()
  return data.buckets
}

import type { ListObjectsResponse, R2Object } from '@/types/object'

type GetObjectsResponse = {
  objects: R2Object[]
  prefix: string
  delimiter: string
  is_truncated: boolean
  next_continuation_token?: string
}

export async function fetchObjects(bucketName: string, prefix: string = ''): Promise<ListObjectsResponse> {
  const params = new URLSearchParams()
  if (prefix) {
    params.set('prefix', prefix)
  }
  params.set('delimiter', '/')

  const url = `${SERVER_URL}/api/v1/buckets/${encodeURIComponent(bucketName)}/objects?${params.toString()}`
  const response = await fetch(url)

  if (!response.ok) {
    const errorData: ErrorResponse = await response.json()
    throw new Error(errorData.error || 'Failed to fetch objects')
  }

  const data: GetObjectsResponse = await response.json()
  return {
    objects: data.objects,
    prefix: data.prefix,
    delimiter: data.delimiter,
    is_truncated: data.is_truncated,
    next_continuation_token: data.next_continuation_token,
  }
}

type ClearCacheResponse = {
  message: string
  deleted?: number
}

export async function clearBucketsCache(): Promise<ClearCacheResponse> {
  const response = await fetch(`${SERVER_URL}/api/v1/cache/api?type=buckets`, {
    method: 'DELETE',
  })

  if (!response.ok) {
    const errorData: ErrorResponse = await response.json()
    throw new Error(errorData.error || 'Failed to clear buckets cache')
  }

  return response.json()
}

export async function clearObjectsCache(bucketName: string): Promise<ClearCacheResponse> {
  const response = await fetch(`${SERVER_URL}/api/v1/cache/api?type=objects&bucket=${encodeURIComponent(bucketName)}`, {
    method: 'DELETE',
  })

  if (!response.ok) {
    const errorData: ErrorResponse = await response.json()
    throw new Error(errorData.error || 'Failed to clear objects cache')
  }

  return response.json()
}

export async function clearContentCache(bucketName: string, objectKey: string): Promise<ClearCacheResponse> {
  const params = new URLSearchParams({
    bucket: bucketName,
    key: objectKey,
  })

  const response = await fetch(`${SERVER_URL}/api/v1/cache/content?${params.toString()}`, {
    method: 'DELETE',
  })

  if (!response.ok) {
    const errorData: ErrorResponse = await response.json()
    throw new Error(errorData.error || 'Failed to clear content cache')
  }

  return response.json()
}
