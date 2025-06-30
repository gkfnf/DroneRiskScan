import axios from "axios"
import type { Asset, ScanTask, Vulnerability, RFSignal, RFThreat, ApiResponse, PaginatedResponse } from "@/types"

// 创建 axios 实例
const apiClient = axios.create({
  baseURL: "/api/v1",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
})

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    // 添加认证token等
    const token = localStorage.getItem("auth_token")
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // 统一错误处理
    if (error.response?.status === 401) {
      // 处理未授权
      localStorage.removeItem("auth_token")
      window.location.href = "/login"
    }
    return Promise.reject(error)
  },
)

// API 接口定义
export const api = {
  // 资产管理
  assets: {
    list: (params?: { type?: string; status?: string; page?: number; limit?: number }) =>
      apiClient.get<PaginatedResponse<Asset>>("/assets", { params }),

    get: (id: string) => apiClient.get<ApiResponse<Asset>>(`/assets/${id}`),

    create: (asset: Omit<Asset, "id" | "createdAt" | "updatedAt">) =>
      apiClient.post<ApiResponse<Asset>>("/assets", asset),

    update: (id: string, updates: Partial<Asset>) => apiClient.put<ApiResponse<Asset>>(`/assets/${id}`, updates),

    delete: (id: string) => apiClient.delete<ApiResponse<void>>(`/assets/${id}`),

    updateStatus: (id: string, status: string) =>
      apiClient.patch<ApiResponse<Asset>>(`/assets/${id}/status`, { status }),
  },

  // 扫描管理
  scans: {
    list: (params?: { assetId?: string; status?: string; page?: number; limit?: number }) =>
      apiClient.get<PaginatedResponse<ScanTask>>("/scans", { params }),

    get: (id: string) => apiClient.get<ApiResponse<ScanTask>>(`/scans/${id}`),

    start: (assetId: string, scanType: string, config?: any) =>
      apiClient.post<ApiResponse<ScanTask>>("/scans/start", { assetId, scanType, config }),

    stop: (id: string) => apiClient.post<ApiResponse<void>>(`/scans/${id}/stop`),

    getProgress: (id: string) =>
      apiClient.get<ApiResponse<{ progress: number; status: string }>>(`/scans/${id}/progress`),
  },

  // 漏洞管理
  vulnerabilities: {
    list: (params?: {
      assetId?: string
      severity?: string
      status?: string
      page?: number
      limit?: number
    }) => apiClient.get<PaginatedResponse<Vulnerability>>("/vulnerabilities", { params }),

    get: (id: string) => apiClient.get<ApiResponse<Vulnerability>>(`/vulnerabilities/${id}`),

    update: (id: string, updates: Partial<Vulnerability>) =>
      apiClient.put<ApiResponse<Vulnerability>>(`/vulnerabilities/${id}`, updates),

    updateStatus: (id: string, status: string) =>
      apiClient.patch<ApiResponse<Vulnerability>>(`/vulnerabilities/${id}/status`, { status }),

    addNote: (id: string, note: string) => apiClient.post<ApiResponse<void>>(`/vulnerabilities/${id}/notes`, { note }),
  },

  // 射频安全
  rf: {
    getSignals: (params?: { frequency?: number; status?: string; limit?: number }) =>
      apiClient.get<PaginatedResponse<RFSignal>>("/rf/signals", { params }),

    getThreats: (params?: { severity?: string; status?: string; limit?: number }) =>
      apiClient.get<PaginatedResponse<RFThreat>>("/rf/threats", { params }),

    startScan: (config?: any) => apiClient.post<ApiResponse<{ taskId: string }>>("/rf/scan/start", config),

    stopScan: () => apiClient.post<ApiResponse<void>>("/rf/scan/stop"),

    getScanStatus: () => apiClient.get<ApiResponse<{ isScanning: boolean; progress: number }>>("/rf/scan/status"),
  },

  // 位置管理
  location: {
    recordProcessing: (data: {
      vulnerabilityId: string
      location: any
      address: string
      notes?: string
    }) => apiClient.post<ApiResponse<void>>("/location/record", data),

    getRecords: (params?: { vulnerabilityId?: string; page?: number; limit?: number }) =>
      apiClient.get<PaginatedResponse<any>>("/location/records", { params }),
  },

  // 报告
  reports: {
    generate: (type: string, params?: any) =>
      apiClient.post<ApiResponse<{ reportId: string }>>("/reports/generate", { type, params }),

    download: (reportId: string) => apiClient.get(`/reports/${reportId}/download`, { responseType: "blob" }),

    list: () => apiClient.get<PaginatedResponse<any>>("/reports"),
  },

  // 系统
  system: {
    health: () => apiClient.get<ApiResponse<{ status: string; timestamp: number }>>("/health"),

    stats: () => apiClient.get<ApiResponse<any>>("/system/stats"),

    config: () => apiClient.get<ApiResponse<any>>("/system/config"),

    updateConfig: (config: any) => apiClient.put<ApiResponse<void>>("/system/config", config),
  },
}

export default apiClient
