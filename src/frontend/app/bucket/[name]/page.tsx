import { fetchObjects } from '@/lib/api'
import { parseObjectsToDisplay } from '@/lib/object-utils'
import { BreadcrumbNav } from '@/components/breadcrumb-nav'
import { ObjectTable } from '@/components/object-table'
import { Card, CardContent, CardHeader } from '@/components/ui/card'

type BucketPageProps = {
  params: Promise<{ name: string }>
  searchParams: Promise<{ prefix?: string }>
}

export default async function BucketPage({ params, searchParams }: BucketPageProps) {
  const { name } = await params
  const { prefix = '' } = await searchParams
  const bucketName = decodeURIComponent(name)

  const response = await fetchObjects(bucketName, prefix)
  const displayObjects = parseObjectsToDisplay(response.objects, prefix)

  return (
    <div className="flex flex-col gap-4 p-6">
      <Card className="w-full max-w-5xl">
        <CardHeader className="pb-3">
          <BreadcrumbNav bucketName={bucketName} prefix={prefix} />
        </CardHeader>
        <CardContent>
          <ObjectTable objects={displayObjects} bucketName={bucketName} />
        </CardContent>
      </Card>
    </div>
  )
}
