import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

export default function BucketLoading() {
  return (
    <div className="flex flex-col gap-4 p-6">
      <Card className="w-full max-w-5xl">
        <CardHeader className="pb-3">
          <Skeleton className="h-5 w-48" />
        </CardHeader>
        <CardContent>
          <div className="w-full">
            {/* Table header */}
            <div className="flex border-b py-3">
              <div className="w-[60%] px-4">
                <Skeleton className="h-4 w-16" />
              </div>
              <div className="w-[20%] px-4">
                <Skeleton className="h-4 w-10" />
              </div>
              <div className="w-[20%] px-4">
                <Skeleton className="h-4 w-24" />
              </div>
            </div>
            {/* Table rows */}
            {[...Array(5)].map((_, i) => (
              <div key={i} className="flex border-b py-3">
                <div className="w-[60%] px-4 flex items-center gap-2">
                  <Skeleton className="h-4 w-4 shrink-0" />
                  <Skeleton className="h-4 w-48" />
                </div>
                <div className="w-[20%] px-4">
                  <Skeleton className="h-4 w-16" />
                </div>
                <div className="w-[20%] px-4">
                  <Skeleton className="h-4 w-28" />
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
