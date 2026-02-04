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

export async function fetchBuckets(): Promise<Bucket[]> {
  const response = await fetch('http://localhost:3000/api/v1/buckets')
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

  const url = `http://localhost:3000/api/v1/buckets/${encodeURIComponent(bucketName)}/objects?${params.toString()}`
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
