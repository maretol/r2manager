'use server'

const BASE_PATH = process.env.BASE_PATH || ''

export async function getObjectURL(bucketName: string, key: string): Promise<string> {
  return `${BASE_PATH}/api/v1/buckets/${encodeURIComponent(bucketName)}/content/${encodeURIComponent(key)}`
}
