'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { FolderIcon } from 'lucide-react'
import { SidebarMenuItem, SidebarMenuButton } from '@/components/ui/sidebar'

type BucketMenuItemProps = {
  bucketName: string
}

export function BucketMenuItem({ bucketName }: BucketMenuItemProps) {
  const pathname = usePathname()
  const href = `/bucket/${encodeURIComponent(bucketName)}`
  const isActive = pathname.startsWith(href)

  return (
    <SidebarMenuItem>
      <SidebarMenuButton asChild isActive={isActive} tooltip={bucketName}>
        <Link href={href}>
          <FolderIcon className="size-4" />
          <span>{bucketName}</span>
        </Link>
      </SidebarMenuButton>
    </SidebarMenuItem>
  )
}
