export type UploadFileStatus = 'pending' | 'uploading' | 'done' | 'error' | 'conflict'

export type UploadFileEntry = {
  id: string
  file: File
  status: UploadFileStatus
  error?: string
}

export type UploadResult = {
  key: string
  size: number
  etag: string
}

export type UploadActionResult =
  | {
      success: true
      data: UploadResult
    }
  | {
      success: false
      error: string
      code?: string
    }

export type CreateDirectoryResult =
  | {
      success: true
      data: UploadResult
    }
  | {
      success: false
      error: string
    }
