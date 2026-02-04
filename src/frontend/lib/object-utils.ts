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
  return date.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
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
