'use client'

import { useState } from 'react'
import { ObjectTable } from '@/components/object-table'
import { ObjectDetailPanel } from '@/components/object-detail-panel'
import type { DisplayObject } from '@/types/object'

type BucketContentProps = {
  objects: DisplayObject[]
  bucketName: string
}

export function BucketContent({ objects, bucketName }: BucketContentProps) {
  const [selectedObject, setSelectedObject] = useState<DisplayObject | null>(null)

  return (
    <div className="flex gap-4">
      <div className={selectedObject ? 'flex-1 min-w-0' : 'w-full'}>
        <ObjectTable
          objects={objects}
          bucketName={bucketName}
          selectedObject={selectedObject}
          onSelectObject={setSelectedObject}
        />
      </div>
      {selectedObject && (
        <div className="w-72 shrink-0 border rounded-lg bg-card">
          <ObjectDetailPanel object={selectedObject} bucketName={bucketName} onClose={() => setSelectedObject(null)} />
        </div>
      )}
    </div>
  )
}
