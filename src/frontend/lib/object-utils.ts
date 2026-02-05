import type { R2Object, DisplayObject } from '@/types/object'

/**
 * R2オブジェクトの配列をDisplayObject配列に変換する
 * フォルダ（CommonPrefixes）とファイルを統合して表示用に整形
 */
export function parseObjectsToDisplay(objects: R2Object[], prefix: string): DisplayObject[] {
  const displayObjects: DisplayObject[] = []
  const seenFolders = new Set<string>()

  for (const obj of objects) {
    const keyWithoutPrefix = prefix ? obj.key.slice(prefix.length) : obj.key

    // キーがプレフィックスで終わる場合はスキップ（フォルダ自体のマーカー）
    if (keyWithoutPrefix === '') {
      continue
    }

    const slashIndex = keyWithoutPrefix.indexOf('/')

    if (slashIndex !== -1) {
      // フォルダ（サブディレクトリがある場合）
      const folderName = keyWithoutPrefix.slice(0, slashIndex)
      if (!seenFolders.has(folderName)) {
        seenFolders.add(folderName)
        displayObjects.push({
          name: folderName,
          key: prefix + folderName + '/',
          isFolder: true,
        })
      }
    } else {
      // ファイル
      displayObjects.push({
        name: keyWithoutPrefix,
        key: obj.key,
        isFolder: false,
        size: obj.size,
        lastModified: obj.last_modified,
        etag: obj.etag,
      })
    }
  }

  // フォルダを先に、その後ファイルを名前順でソート
  return displayObjects.sort((a, b) => {
    if (a.isFolder !== b.isFolder) {
      return a.isFolder ? -1 : 1
    }
    return a.name.localeCompare(b.name)
  })
}

/**
 * ファイルサイズを人間が読みやすい形式にフォーマット
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'

  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const k = 1024
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${units[i]}`
}

/**
 * 日付を表示用にフォーマット
 */
export function formatDate(dateString: string): string {
  const date = new Date(dateString)
  const isUTC = dateString.endsWith('Z') || dateString.endsWith('+00:00') || dateString.endsWith('+0000')

  if (isUTC) {
    // UTCの場合、+9時間でJSTに補正
    const jstDate = new Date(date.getTime() + 9 * 60 * 60 * 1000)
    const year = jstDate.getUTCFullYear()
    const month = String(jstDate.getUTCMonth() + 1).padStart(2, '0')
    const day = String(jstDate.getUTCDate()).padStart(2, '0')
    const hour = String(jstDate.getUTCHours()).padStart(2, '0')
    const minute = String(jstDate.getUTCMinutes()).padStart(2, '0')
    return `${year}/${month}/${day} ${hour}:${minute}`
  }

  // UTC以外の場合はそのままフォーマット
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  return `${year}/${month}/${day} ${hour}:${minute}`
}

/**
 * プレフィックスからパンくず用のパス配列を生成
 */
export function parsePrefixToBreadcrumbs(prefix: string): { name: string; prefix: string }[] {
  if (!prefix) return []

  const parts = prefix.split('/').filter(Boolean)
  const breadcrumbs: { name: string; prefix: string }[] = []

  let currentPrefix = ''
  for (const part of parts) {
    currentPrefix += part + '/'
    breadcrumbs.push({
      name: part,
      prefix: currentPrefix,
    })
  }

  return breadcrumbs
}
