import { ObjectTable } from '@/components/object-table'
import { ObjectDetailPanel } from '@/components/object-detail-panel'
import { AlertCircle } from 'lucide-react'
import type { DisplayObject } from '@/types/object'

type BucketContentProps = {
  objects: DisplayObject[]
  bucketName: string
  prefix: string
  selectedObject: DisplayObject | null
  selectedNotFound: boolean
  publicUrl: string
}

export function BucketContent({ objects, bucketName, prefix, selectedObject, selectedNotFound, publicUrl }: BucketContentProps) {
  return (
    <div className="flex gap-4">
      <div className="flex-1 min-w-0">
        <ObjectTable objects={objects} bucketName={bucketName} prefix={prefix} selectedObject={selectedObject} />
      </div>
      <div className="w-1/4 shrink-0 border rounded-lg bg-card">
        {selectedObject && (
          <ObjectDetailPanel object={selectedObject} bucketName={bucketName} prefix={prefix} publicUrl={publicUrl} />
        )}
        {selectedNotFound && (
          <div className="p-4">
            <div className="flex items-center gap-2 text-destructive">
              <AlertCircle className="size-4 shrink-0" />
              <p className="text-sm">指定されたオブジェクトが見つかりません</p>
            </div>
          </div>
        )}
        {!selectedObject && !selectedNotFound && (
          <div className="flex items-center justify-center h-full p-4">
            <p className="text-sm text-muted-foreground">オブジェクトを選択してください</p>
          </div>
        )}
      </div>
    </div>
  )
}
