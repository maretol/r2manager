import { Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarHeader } from './ui/sidebar'

export function AppSideBar() {
  return (
    <Sidebar>
      <SidebarHeader />
      <SidebarContent>
        <SidebarGroup />
        <SidebarGroup />
      </SidebarContent>
      <SidebarFooter />
    </Sidebar>
  )
}
