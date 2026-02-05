import { fetchObjects } from '@/lib/api'
import { parseObjectsToDisplay } from '@/lib/object-utils'
import { BreadcrumbNav } from '@/components/breadcrumb-nav'
import { BucketContent } from '@/components/bucket-content'
import { Card, CardContent, CardHeader } from '@/components/ui/card'

type BucketPageProps = {
  params: Promise<{ name: string }>
  searchParams: Promise<{ prefix?: string; selected?: string }>
}

export default async function BucketPage({ params, searchParams }: BucketPageProps) {
  const { name } = await params
  const { prefix = '', selected } = await searchParams
  const bucketName = decodeURIComponent(name)

  const response = await fetchObjects(bucketName, prefix)
  const displayObjects = parseObjectsToDisplay(response.objects, prefix)

  const selectedObject = selected ? (displayObjects.find((obj) => obj.key === selected) ?? null) : null
  const selectedNotFound = !!selected && !selectedObject

  return (
    <div className="flex flex-col gap-4 p-6">
      <Card className="w-full">
        <CardHeader className="pb-3">
          <BreadcrumbNav bucketName={bucketName} prefix={prefix} />
        </CardHeader>
        <CardContent>
          <BucketContent
            objects={displayObjects}
            bucketName={bucketName}
            prefix={prefix}
            selectedObject={selectedObject}
            selectedNotFound={selectedNotFound}
          />
        </CardContent>
      </Card>
    </div>
  )
}
