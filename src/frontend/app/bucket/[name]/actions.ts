'use server'

const BASE_PATH = process.env.BASE_PATH || ''

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
