import Link from 'next/link'
import { ChevronRightIcon, HomeIcon } from 'lucide-react'
import { parsePrefixToBreadcrumbs } from '@/lib/object-utils'

type BreadcrumbNavProps = {
  bucketName: string
  prefix: string
}

export function BreadcrumbNav({ bucketName, prefix }: BreadcrumbNavProps) {
  const breadcrumbs = parsePrefixToBreadcrumbs(prefix)
  const basePath = `/bucket/${encodeURIComponent(bucketName)}`

  return (
    <nav className="flex items-center gap-1 text-sm text-muted-foreground">
      <Link
        href={basePath}
        className="flex items-center gap-1 hover:text-foreground transition-colors"
      >
        <HomeIcon className="size-4" />
        <span>{bucketName}</span>
      </Link>
      {breadcrumbs.map((crumb, index) => (
        <span key={crumb.prefix} className="flex items-center gap-1">
          <ChevronRightIcon className="size-4" />
          {index === breadcrumbs.length - 1 ? (
            <span className="text-foreground font-medium">{crumb.name}</span>
          ) : (
            <Link
              href={`${basePath}?prefix=${encodeURIComponent(crumb.prefix)}`}
              className="hover:text-foreground transition-colors"
            >
              {crumb.name}
            </Link>
          )}
        </span>
      ))}
    </nav>
  )
}
