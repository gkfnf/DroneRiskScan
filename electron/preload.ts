import { contextBridge, ipcRenderer } from "electron"

// 暴露安全的API给渲染进程
contextBridge.exposeInMainWorld("electronAPI", {
  // 文件操作
  showSaveDialog: (options: any) => ipcRenderer.invoke("show-save-dialog", options),
  showOpenDialog: (options: any) => ipcRenderer.invoke("show-open-dialog", options),
  writeFile: (filePath: string, data: string) => ipcRenderer.invoke("write-file", filePath, data),
  readFile: (filePath: string) => ipcRenderer.invoke("read-file", filePath),

  // 系统信息
  getSystemInfo: () => ipcRenderer.invoke("get-system-info"),
  getNetworkInterfaces: () => ipcRenderer.invoke("get-network-interfaces"),
  getSerialPorts: () => ipcRenderer.invoke("get-serial-ports"),

  // 通知
  showNotification: (options: any) => ipcRenderer.invoke("show-notification", options),

  // 窗口控制
  minimizeWindow: () => ipcRenderer.invoke("minimize-window"),
  maximizeWindow: () => ipcRenderer.invoke("maximize-window"),
  closeWindow: () => ipcRenderer.invoke("close-window"),

  // 事件监听
  onNavigateTo: (callback: (path: string) => void) => {
    ipcRenderer.on("navigate-to", (event, path) => callback(path))
  },
  onCreateScanTask: (callback: () => void) => {
    ipcRenderer.on("create-scan-task", callback)
  },
  onImportAssets: (callback: (filePath: string) => void) => {
    ipcRenderer.on("import-assets", (event, filePath) => callback(filePath))
  },
  onExportReport: (callback: (filePath: string) => void) => {
    ipcRenderer.on("export-report", (event, filePath) => callback(filePath))
  },
  onStartFullScan: (callback: () => void) => {
    ipcRenderer.on("start-full-scan", callback)
  },
  onStartRfScan: (callback: () => void) => {
    ipcRenderer.on("start-rf-scan", callback)
  },
  onStopAllScans: (callback: () => void) => {
    ipcRenderer.on("stop-all-scans", callback)
  },

  // 移除事件监听
  removeAllListeners: (channel: string) => {
    ipcRenderer.removeAllListeners(channel)
  },
})

// 类型定义
declare global {
  interface Window {
    electronAPI: {
      showSaveDialog: (options: any) => Promise<any>
      showOpenDialog: (options: any) => Promise<any>
      writeFile: (filePath: string, data: string) => Promise<any>
      readFile: (filePath: string) => Promise<any>
      getSystemInfo: () => Promise<any>
      getNetworkInterfaces: () => Promise<any>
      getSerialPorts: () => Promise<any>
      showNotification: (options: any) => Promise<any>
      minimizeWindow: () => Promise<void>
      maximizeWindow: () => Promise<void>
      closeWindow: () => Promise<void>
      onNavigateTo: (callback: (path: string) => void) => void
      onCreateScanTask: (callback: () => void) => void
      onImportAssets: (callback: (filePath: string) => void) => void
      onExportReport: (callback: (filePath: string) => void) => void
      onStartFullScan: (callback: () => void) => void
      onStartRfScan: (callback: () => void) => void
      onStopAllScans: (callback: () => void) => void
      removeAllListeners: (channel: string) => void
    }
  }
}
