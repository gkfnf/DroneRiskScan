"use client"

import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Monitor, HardDrive, Cpu, Wifi, Usb, Download, Upload, Bell, Minimize2, Maximize2, X, Info } from "lucide-react"

interface SystemInfo {
  platform: string
  arch: string
  version: string
  memory: number
  cpus: number
  userDataPath: string
  appVersion: string
}

interface NetworkInterface {
  name: string
  addresses: Array<{
    address: string
    family: string
    internal: boolean
  }>
}

export default function ElectronIntegration() {
  const [systemInfo, setSystemInfo] = useState<SystemInfo | null>(null)
  const [networkInterfaces, setNetworkInterfaces] = useState<NetworkInterface[]>([])
  const [serialPorts, setSerialPorts] = useState<any[]>([])
  const [isElectron, setIsElectron] = useState(false)

  useEffect(() => {
    // 检查是否在Electron环境中运行
    if (typeof window !== "undefined" && window.electronAPI) {
      setIsElectron(true)
      loadSystemInfo()
      loadNetworkInterfaces()
      loadSerialPorts()
      setupEventListeners()
    }

    return () => {
      if (window.electronAPI) {
        window.electronAPI.removeAllListeners("navigate-to")
        window.electronAPI.removeAllListeners("create-scan-task")
        window.electronAPI.removeAllListeners("import-assets")
        window.electronAPI.removeAllListeners("export-report")
        window.electronAPI.removeAllListeners("start-full-scan")
        window.electronAPI.removeAllListeners("start-rf-scan")
        window.electronAPI.removeAllListeners("stop-all-scans")
      }
    }
  }, [])

  const loadSystemInfo = async () => {
    try {
      const info = await window.electronAPI.getSystemInfo()
      setSystemInfo(info)
    } catch (error) {
      console.error("Failed to load system info:", error)
    }
  }

  const loadNetworkInterfaces = async () => {
    try {
      const result = await window.electronAPI.getNetworkInterfaces()
      if (result.success) {
        const interfaces = Object.entries(result.interfaces).map(([name, addresses]: [string, any]) => ({
          name,
          addresses: addresses.filter((addr: any) => !addr.internal),
        }))
        setNetworkInterfaces(interfaces)
      }
    } catch (error) {
      console.error("Failed to load network interfaces:", error)
    }
  }

  const loadSerialPorts = async () => {
    try {
      const result = await window.electronAPI.getSerialPorts()
      if (result.success) {
        setSerialPorts(result.ports)
      }
    } catch (error) {
      console.error("Failed to load serial ports:", error)
    }
  }

  const setupEventListeners = () => {
    window.electronAPI.onNavigateTo((path: string) => {
      console.log("Navigate to:", path)
      // 这里可以实现路由跳转逻辑
    })

    window.electronAPI.onCreateScanTask(() => {
      console.log("Create scan task triggered")
      // 实现创建扫描任务逻辑
    })

    window.electronAPI.onImportAssets((filePath: string) => {
      console.log("Import assets from:", filePath)
      handleImportAssets(filePath)
    })

    window.electronAPI.onExportReport((filePath: string) => {
      console.log("Export report to:", filePath)
      handleExportReport(filePath)
    })

    window.electronAPI.onStartFullScan(() => {
      console.log("Start full scan triggered")
      // 实现全面扫描逻辑
    })

    window.electronAPI.onStartRfScan(() => {
      console.log("Start RF scan triggered")
      // 实现射频扫描逻辑
    })

    window.electronAPI.onStopAllScans(() => {
      console.log("Stop all scans triggered")
      // 实现停止扫描逻辑
    })
  }

  const handleImportAssets = async (filePath: string) => {
    try {
      const result = await window.electronAPI.readFile(filePath)
      if (result.success) {
        console.log("File content:", result.data)
        // 处理导入的资产数据
        await window.electronAPI.showNotification({
          title: "导入成功",
          body: `已成功导入资产文件: ${filePath}`,
        })
      }
    } catch (error) {
      console.error("Failed to import assets:", error)
    }
  }

  const handleExportReport = async (filePath: string) => {
    try {
      const reportData = generateReportData()
      const result = await window.electronAPI.writeFile(filePath, reportData)
      if (result.success) {
        await window.electronAPI.showNotification({
          title: "导出成功",
          body: `报告已保存到: ${filePath}`,
        })
      }
    } catch (error) {
      console.error("Failed to export report:", error)
    }
  }

  const generateReportData = () => {
    // 生成报告数据
    const report = {
      timestamp: new Date().toISOString(),
      systemInfo,
      networkInterfaces,
      serialPorts,
      // 添加更多报告数据
    }
    return JSON.stringify(report, null, 2)
  }

  const handleWindowControl = async (action: "minimize" | "maximize" | "close") => {
    try {
      switch (action) {
        case "minimize":
          await window.electronAPI.minimizeWindow()
          break
        case "maximize":
          await window.electronAPI.maximizeWindow()
          break
        case "close":
          await window.electronAPI.closeWindow()
          break
      }
    } catch (error) {
      console.error("Window control error:", error)
    }
  }

  const formatBytes = (bytes: number) => {
    const sizes = ["Bytes", "KB", "MB", "GB", "TB"]
    if (bytes === 0) return "0 Bytes"
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return Math.round((bytes / Math.pow(1024, i)) * 100) / 100 + " " + sizes[i]
  }

  if (!isElectron) {
    return (
      <Alert>
        <Info className="h-4 w-4" />
        <AlertDescription>此功能仅在Electron桌面应用中可用。请下载并安装桌面版本以使用完整功能。</AlertDescription>
      </Alert>
    )
  }

  return (
    <div className="space-y-6">
      {/* Window Controls */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            <span className="flex items-center gap-2">
              <Monitor className="h-5 w-5" />
              窗口控制
            </span>
            <div className="flex gap-2">
              <Button size="sm" variant="outline" onClick={() => handleWindowControl("minimize")}>
                <Minimize2 className="h-4 w-4" />
              </Button>
              <Button size="sm" variant="outline" onClick={() => handleWindowControl("maximize")}>
                <Maximize2 className="h-4 w-4" />
              </Button>
              <Button size="sm" variant="outline" onClick={() => handleWindowControl("close")}>
                <X className="h-4 w-4" />
              </Button>
            </div>
          </CardTitle>
        </CardHeader>
      </Card>

      {/* System Information */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Cpu className="h-5 w-5" />
            系统信息
          </CardTitle>
        </CardHeader>
        <CardContent>
          {systemInfo && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-sm text-gray-600">操作系统:</span>
                  <Badge variant="outline">
                    {systemInfo.platform} {systemInfo.arch}
                  </Badge>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm text-gray-600">系统版本:</span>
                  <span className="text-sm">{systemInfo.version}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm text-gray-600">CPU核心:</span>
                  <span className="text-sm">{systemInfo.cpus} 核</span>
                </div>
              </div>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-sm text-gray-600">内存:</span>
                  <span className="text-sm">{formatBytes(systemInfo.memory)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm text-gray-600">应用版本:</span>
                  <span className="text-sm">{systemInfo.appVersion}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-sm text-gray-600">数据目录:</span>
                  <span className="text-xs text-gray-500 truncate max-w-32" title={systemInfo.userDataPath}>
                    {systemInfo.userDataPath}
                  </span>
                </div>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Network Interfaces */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Wifi className="h-5 w-5" />
            网络接口
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {networkInterfaces.map((iface, index) => (
              <div key={index} className="border rounded-lg p-3">
                <div className="flex items-center justify-between mb-2">
                  <span className="font-medium">{iface.name}</span>
                  <Badge variant="outline">{iface.addresses.length} 地址</Badge>
                </div>
                <div className="space-y-1">
                  {iface.addresses.map((addr, addrIndex) => (
                    <div key={addrIndex} className="flex justify-between text-sm">
                      <span className="text-gray-600">{addr.family}:</span>
                      <span className="font-mono">{addr.address}</span>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Serial Ports */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Usb className="h-5 w-5" />
            串口设备
          </CardTitle>
        </CardHeader>
        <CardContent>
          {serialPorts.length > 0 ? (
            <div className="space-y-3">
              {serialPorts.map((port, index) => (
                <div key={index} className="border rounded-lg p-3">
                  <div className="flex items-center justify-between mb-2">
                    <span className="font-medium">{port.path}</span>
                    <Badge variant="outline">串口</Badge>
                  </div>
                  <div className="text-sm space-y-1">
                    {port.manufacturer && (
                      <div className="flex justify-between">
                        <span className="text-gray-600">制造商:</span>
                        <span>{port.manufacturer}</span>
                      </div>
                    )}
                    {port.productId && (
                      <div className="flex justify-between">
                        <span className="text-gray-600">产品ID:</span>
                        <span className="font-mono">{port.productId}</span>
                      </div>
                    )}
                    {port.vendorId && (
                      <div className="flex justify-between">
                        <span className="text-gray-600">厂商ID:</span>
                        <span className="font-mono">{port.vendorId}</span>
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500">
              <Usb className="h-12 w-12 mx-auto mb-4 opacity-50" />
              <p>未检测到串口设备</p>
            </div>
          )}
        </CardContent>
      </Card>

      {/* File Operations */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <HardDrive className="h-5 w-5" />
            文件操作
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 gap-4">
            <Button variant="outline" className="h-20 flex-col bg-transparent">
              <Upload className="h-6 w-6 mb-2" />
              <span className="text-sm">导入资产</span>
            </Button>
            <Button variant="outline" className="h-20 flex-col bg-transparent">
              <Download className="h-6 w-6 mb-2" />
              <span className="text-sm">导出报告</span>
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Notifications */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Bell className="h-5 w-5" />
            通知测试
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Button
            onClick={() => {
              window.electronAPI.showNotification({
                title: "测试通知",
                body: "这是一个来自电力无人机安全扫描系统的测试通知",
              })
            }}
          >
            发送测试通知
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}
