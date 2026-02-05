import Link from 'next/link'
import { FileIcon, FolderIcon } from 'lucide-react'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import type { DisplayObject } from '@/types/object'
import { formatFileSize, formatDate } from '@/lib/object-utils'
import { cn } from '@/lib/utils'

type ObjectTableProps = {
  objects: DisplayObject[]
  bucketName: string
  prefix: string
  selectedObject?: DisplayObject | null
}

export function ObjectTable({ objects, bucketName, prefix, selectedObject }: ObjectTableProps) {
  const basePath = `/bucket/${encodeURIComponent(bucketName)}`

  const getFileHref = (obj: DisplayObject) => {
    const params = new URLSearchParams()
    if (prefix) params.set('prefix', prefix)
    if (selectedObject?.key !== obj.key) {
      params.set('selected', obj.key)
    }
    const query = params.toString()
    return `${basePath}${query ? `?${query}` : ''}`
  }

  if (objects.length === 0) {
    return (
      <div className="flex items-center justify-center py-12 text-muted-foreground">
        <p>No objects found</p>
      </div>
    )
  }

  return (
    <Table className="table-fixed">
      <TableHeader>
        <TableRow>
          <TableHead className="w-[60%]">Name</TableHead>
          <TableHead className="w-[20%]">Size</TableHead>
          <TableHead className="w-[20%]">Last Modified</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {objects.map((obj) => {
          const isSelected = !obj.isFolder && selectedObject?.key === obj.key
          if (obj.isFolder) {
            return (
              <TableRow key={obj.key}>
                <TableCell className="truncate">
                  <Link
                    href={`${basePath}?prefix=${encodeURIComponent(obj.key)}`}
                    className="flex items-center gap-2 hover:text-primary transition-colors"
                  >
                    <FolderIcon className="size-4 shrink-0 text-yellow-500" />
                    <span className="truncate">{obj.name}</span>
                  </Link>
                </TableCell>
                <TableCell className="text-muted-foreground">-</TableCell>
                <TableCell className="text-muted-foreground whitespace-nowrap">-</TableCell>
              </TableRow>
            )
          }

          const href = getFileHref(obj)
          return (
            <TableRow key={obj.key} className={cn(isSelected && 'bg-accent')}>
              <TableCell className="truncate p-0">
                <Link href={href} className="flex items-center gap-2 p-2 hover:text-primary transition-colors">
                  <FileIcon className="size-4 shrink-0 text-muted-foreground" />
                  <span className="truncate">{obj.name}</span>
                </Link>
              </TableCell>
              <TableCell className="text-muted-foreground p-0">
                <Link href={href} className="block p-2">
                  {formatFileSize(obj.size ?? 0)}
                </Link>
              </TableCell>
              <TableCell className="text-muted-foreground whitespace-nowrap p-0">
                <Link href={href} className="block p-2">
                  {obj.lastModified ? formatDate(obj.lastModified) : '-'}
                </Link>
              </TableCell>
            </TableRow>
          )
        })}
      </TableBody>
    </Table>
  )
}
