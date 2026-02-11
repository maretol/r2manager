'use server'

import type { UploadActionResult, CreateDirectoryResult } from '@/types/upload'

const BASE_PATH = process.env.BASE_PATH || ''

const HOST = process.env.HOSTNAME
const PROTOCOL = process.env.PROTOCOL || 'http'
const PORT = process.env.PORT ? `:${process.env.PORT}` : ':3000'
const SERVER_URL = (HOST ? PROTOCOL + '://' + (HOST + PORT) : '') + BASE_PATH

type ObjectURLs = {
  objectUrl: string
  publicObjectUrl: string | null
}

export async function getObjectURLs(
  bucketName: string,
  key: string,
  publicBaseUrl: string
): Promise<ObjectURLs> {
  const objectUrl = `${BASE_PATH}/api/v1/buckets/${encodeURIComponent(bucketName)}/content/${encodeURIComponent(key)}`

  let publicObjectUrl: string | null = null
  if (publicBaseUrl.length > 0) {
    const baseUrl = publicBaseUrl.endsWith('/') ? publicBaseUrl.slice(0, -1) : publicBaseUrl
    publicObjectUrl = `${baseUrl}/${key}`
  }

  return { objectUrl, publicObjectUrl }
}

export async function uploadFile(formData: FormData): Promise<UploadActionResult> {
  const file = formData.get('file') as File
  const bucketName = formData.get('bucketName') as string
  const key = formData.get('key') as string
  const overwrite = formData.get('overwrite') === 'true'

  if (!file || !bucketName || !key) {
    return { success: false, error: '必須パラメータが不足しています' }
  }

  const uploadFormData = new FormData()
  uploadFormData.append('file', file)

  const encodedKey = key
    .split('/')
    .map(encodeURIComponent)
    .join('/')
  const query = overwrite ? '?overwrite=true' : ''
  const url = `${SERVER_URL}/api/v1/buckets/${encodeURIComponent(bucketName)}/objects/${encodedKey}${query}`

  const response = await fetch(url, {
    method: 'PUT',
    body: uploadFormData,
  })

  if (!response.ok) {
    const errorData = await response.json()
    return {
      success: false,
      error: errorData.error || 'アップロードに失敗しました',
      code: errorData.code,
    }
  }

  const data = await response.json()
  return { success: true, data }
}

export async function createDirectoryAction(
  bucketName: string,
  path: string
): Promise<CreateDirectoryResult> {
  if (!bucketName || !path) {
    return { success: false, error: '必須パラメータが不足しています' }
  }

  const url = `${SERVER_URL}/api/v1/buckets/${encodeURIComponent(bucketName)}/directories`

  const response = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ path }),
  })

  if (!response.ok) {
    const errorData = await response.json()
    return { success: false, error: errorData.error || 'ディレクトリの作成に失敗しました' }
  }

  const data = await response.json()
  return { success: true, data }
}
