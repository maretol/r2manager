import Link from 'next/link'
import { AlertCircleIcon, SettingsIcon } from 'lucide-react'
import { fetchBuckets } from '@/lib/api'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
} from './ui/sidebar'
import { BucketMenuItem } from './bucket-menu-item'
import { RefreshBucketsButton } from './refresh-buckets-button'

export async function AppSideBar() {
  let buckets: Awaited<ReturnType<typeof fetchBuckets>> | null = null
  try {
    buckets = await fetchBuckets()
  } catch {
    // fallback below
  }

  return (
    <Sidebar>
      <SidebarHeader>
        <Link href="/" className="flex items-center gap-2 px-2 py-1">
          <h1 className="text-lg font-semibold">R2 Contents Manager</h1>
        </Link>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <div className="flex items-center justify-between pr-2">
            <SidebarGroupLabel>Buckets</SidebarGroupLabel>
            <RefreshBucketsButton />
          </div>
          {buckets ? (
            <SidebarMenu>
              {buckets.map((bucket) => (
                <BucketMenuItem key={bucket.name} bucketName={bucket.name} />
              ))}
            </SidebarMenu>
          ) : (
            <div className="flex flex-col items-center gap-2 px-4 py-6 text-sm text-muted-foreground">
              <AlertCircleIcon className="size-5 text-destructive" />
              <p className="text-center">バケット一覧の取得に失敗しました</p>
              <p className="text-center text-xs">リロードするか、設定を確認してください</p>
            </div>
          )}
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton asChild tooltip="Settings">
              <Link href="/settings">
                <SettingsIcon className="size-4" />
                <span>Settings</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}
