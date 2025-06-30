import { defineStore } from "pinia"
import { ref, computed } from "vue"
import type { Asset, ScanTask, Vulnerability, RFSignal, RFThreat } from "@/types"
import { api } from "@/api"

export const useSystemStore = defineStore("system", () => {
  // 状态
  const assets = ref<Asset[]>([])
  const scanTasks = ref<ScanTask[]>([])
  const vulnerabilities = ref<Vulnerability[]>([])
  const rfSignals = ref<RFSignal[]>([])
  const rfThreats = ref<RFThreat[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const onlineAssets = computed(() => assets.value.filter((asset) => asset.status === "online"))

  const criticalVulnerabilities = computed(() => vulnerabilities.value.filter((vuln) => vuln.severity === "critical"))

  const runningScanTasks = computed(() => scanTasks.value.filter((task) => task.status === "running"))

  const activeThreat = computed(() => rfThreats.value.filter((threat) => threat.status === "active"))

  // 统计数据
  const stats = computed(() => ({
    totalAssets: assets.value.length,
    onlineAssets: onlineAssets.value.length,
    totalVulnerabilities: vulnerabilities.value.length,
    criticalVulnerabilities: criticalVulnerabilities.value.length,
    highVulnerabilities: vulnerabilities.value.filter((v) => v.severity === "high").length,
    mediumVulnerabilities: vulnerabilities.value.filter((v) => v.severity === "medium").length,
    lowVulnerabilities: vulnerabilities.value.filter((v) => v.severity === "low").length,
    rfThreats: activeThreat.value.length,
    runningScanTasks: runningScanTasks.value.length,
  }))

  // 方法
  const initialize = async () => {
    try {
      isLoading.value = true
      await Promise.all([loadAssets(), loadScanTasks(), loadVulnerabilities(), loadRFSignals(), loadRFThreats()])
    } catch (err) {
      error.value = err instanceof Error ? err.message : "初始化失败"
    } finally {
      isLoading.value = false
    }
  }

  const loadAssets = async () => {
    const response = await api.assets.list()
    assets.value = response.data
  }

  const loadScanTasks = async () => {
    const response = await api.scans.list()
    scanTasks.value = response.data
  }

  const loadVulnerabilities = async () => {
    const response = await api.vulnerabilities.list()
    vulnerabilities.value = response.data
  }

  const loadRFSignals = async () => {
    const response = await api.rf.getSignals()
    rfSignals.value = response.data
  }

  const loadRFThreats = async () => {
    const response = await api.rf.getThreats()
    rfThreats.value = response.data
  }

  const createAsset = async (asset: Omit<Asset, "id" | "createdAt" | "updatedAt">) => {
    const response = await api.assets.create(asset)
    assets.value.push(response.data)
    return response.data
  }

  const updateAsset = async (id: string, updates: Partial<Asset>) => {
    const response = await api.assets.update(id, updates)
    const index = assets.value.findIndex((asset) => asset.id === id)
    if (index !== -1) {
      assets.value[index] = response.data
    }
    return response.data
  }

  const deleteAsset = async (id: string) => {
    await api.assets.delete(id)
    assets.value = assets.value.filter((asset) => asset.id !== id)
  }

  const startScan = async (assetId: string, scanType: string) => {
    const response = await api.scans.start(assetId, scanType)
    scanTasks.value.push(response.data)
    return response.data
  }

  const updateScanProgress = (taskId: string, progress: number, status?: string) => {
    const task = scanTasks.value.find((t) => t.id === taskId)
    if (task) {
      task.progress = progress
      if (status) task.status = status
    }
  }

  const clearError = () => {
    error.value = null
  }

  return {
    // 状态
    assets,
    scanTasks,
    vulnerabilities,
    rfSignals,
    rfThreats,
    isLoading,
    error,

    // 计算属性
    onlineAssets,
    criticalVulnerabilities,
    runningScanTasks,
    activeThreat,
    stats,

    // 方法
    initialize,
    loadAssets,
    loadScanTasks,
    loadVulnerabilities,
    loadRFSignals,
    loadRFThreats,
    createAsset,
    updateAsset,
    deleteAsset,
    startScan,
    updateScanProgress,
    clearError,
  }
})
