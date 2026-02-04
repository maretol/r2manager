'use client'

import Link from 'next/link'
import { FileIcon, FolderIcon } from 'lucide-react'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import type { DisplayObject } from '@/types/object'
import { formatFileSize, formatDate } from '@/lib/object-utils'

type ObjectTableProps = {
  objects: DisplayObject[]
  bucketName: string
}

export function ObjectTable({ objects, bucketName }: ObjectTableProps) {
  const basePath = `/bucket/${encodeURIComponent(bucketName)}`

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
        {objects.map((obj) => (
          <TableRow key={obj.key}>
            <TableCell className="truncate">
              {obj.isFolder ? (
                <Link
                  href={`${basePath}?prefix=${encodeURIComponent(obj.key)}`}
                  className="flex items-center gap-2 hover:text-primary transition-colors"
                >
                  <FolderIcon className="size-4 shrink-0 text-yellow-500" />
                  <span className="truncate">{obj.name}</span>
                </Link>
              ) : (
                <div className="flex items-center gap-2">
                  <FileIcon className="size-4 shrink-0 text-muted-foreground" />
                  <span className="truncate">{obj.name}</span>
                </div>
              )}
            </TableCell>
            <TableCell className="text-muted-foreground">
              {obj.isFolder ? '-' : formatFileSize(obj.size ?? 0)}
            </TableCell>
            <TableCell className="text-muted-foreground whitespace-nowrap">
              {obj.lastModified ? formatDate(obj.lastModified) : '-'}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}
