import { createRouter, createWebHistory } from "vue-router"
import type { RouteRecordRaw } from "vue-router"

const routes: RouteRecordRaw[] = [
  {
    path: "/",
    name: "Dashboard",
    component: () => import("@/views/Dashboard.vue"),
    meta: { title: "仪表板" },
  },
  {
    path: "/assets",
    name: "Assets",
    component: () => import("@/views/Assets.vue"),
    meta: { title: "资产管理" },
  },
  {
    path: "/scanning",
    name: "Scanning",
    component: () => import("@/views/Scanning.vue"),
    meta: { title: "扫描监控" },
  },
  {
    path: "/rf-security",
    name: "RFSecurity",
    component: () => import("@/views/RFSecurity.vue"),
    meta: { title: "射频安全" },
  },
  {
    path: "/vulnerabilities",
    name: "Vulnerabilities",
    component: () => import("@/views/Vulnerabilities.vue"),
    meta: { title: "漏洞管理" },
  },
  {
    path: "/vulnerability/:id",
    name: "VulnerabilityDetail",
    component: () => import("@/views/VulnerabilityDetail.vue"),
    meta: { title: "漏洞详情" },
  },
  {
    path: "/location",
    name: "Location",
    component: () => import("@/views/Location.vue"),
    meta: { title: "位置追踪" },
  },
  {
    path: "/reports",
    name: "Reports",
    component: () => import("@/views/Reports.vue"),
    meta: { title: "安全报告" },
  },
  {
    path: "/settings",
    name: "Settings",
    component: () => import("@/views/Settings.vue"),
    meta: { title: "系统设置" },
  },
  {
    path: "/mobile",
    children: [
      {
        path: "vulnerability",
        name: "MobileVulnerability",
        component: () => import("@/views/mobile/VulnerabilityList.vue"),
      },
      {
        path: "vulnerability/:id",
        name: "MobileVulnerabilityDetail",
        component: () => import("@/views/mobile/VulnerabilityDetail.vue"),
      },
      {
        path: "location",
        name: "MobileLocation",
        component: () => import("@/views/mobile/Location.vue"),
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  if (to.meta?.title) {
    document.title = `${to.meta.title} - 电力无人机安全扫描系统`
  }
  next()
})

export default router
