<<<<<<< HEAD
"use client"

import { useState } from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Badge } from "@/components/ui/badge"
import { Progress } from "@/components/ui/progress"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Switch } from "@/components/ui/switch"
import { ScrollArea } from "@/components/ui/scroll-area"
import {
  Shield,
  Wifi,
  Server,
  Database,
  Zap,
  AlertTriangle,
  CheckCircle,
  Plus,
  Search,
  FileText,
  Activity,
  Radio,
  Satellite,
  Signal,
  Radar,
  WifiOff,
  Volume2,
  Eye,
  Play,
  Pause,
  RotateCcw,
  Download,
  Upload,
  Bluetooth,
} from "lucide-react"

interface Asset {
  id: string
  name: string
  type: "drone" | "gcs" | "server" | "network"
  ip: string
  status: "online" | "offline" | "scanning"
  lastScan?: string
  vulnerabilities?: number
}

interface ScanTask {
  id: string
  assetName: string
  progress: number
  status: "running" | "completed" | "failed"
  startTime: string
}

interface RFSignal {
  id: string
  frequency: number
  strength: number
  type: "GPS" | "RC" | "Video" | "Telemetry" | "WiFi" | "Bluetooth" | "Unknown"
  source?: string
  status: "normal" | "suspicious" | "threat"
  lastDetected: string
}

interface RFThreat {
  id: string
  type: "jamming" | "spoofing" | "interception" | "unauthorized"
  frequency: number
  severity: "low" | "medium" | "high" | "critical"
  description: string
  detectedAt: string
  affectedSystems: string[]
}

export default function DroneSecurityScanner() {
  const [assets, setAssets] = useState<Asset[]>([
    {
      id: "1",
      name: "DJI M300 RTK #001",
      type: "drone",
      ip: "192.168.1.100",
      status: "online",
      lastScan: "2024-01-15 14:30",
      vulnerabilities: 2,
    },
    {
      id: "2",
      name: "Ground Control Station",
      type: "gcs",
      ip: "192.168.1.50",
      status: "online",
      lastScan: "2024-01-15 12:15",
      vulnerabilities: 0,
    },
    {
      id: "3",
      name: "Video Stream Server",
      type: "server",
      ip: "192.168.1.200",
      status: "scanning",
      vulnerabilities: 1,
    },
    {
      id: "4",
      name: "Network Router",
      type: "network",
      ip: "192.168.1.1",
      status: "offline",
      lastScan: "2024-01-14 09:20",
      vulnerabilities: 3,
    },
  ])

  const [scanTasks, setScanTasks] = useState<ScanTask[]>([
    {
      id: "1",
      assetName: "Video Stream Server",
      progress: 65,
      status: "running",
      startTime: "2024-01-15 15:30",
    },
  ])

  const [rfSignals, setRfSignals] = useState<RFSignal[]>([
    {
      id: "1",
      frequency: 2400.5,
      strength: -45,
      type: "RC",
      source: "DJI Controller",
      status: "normal",
      lastDetected: "2024-01-15 15:30:25",
    },
    {
      id: "2",
      frequency: 5800.2,
      strength: -38,
      type: "Video",
      source: "DJI M300 RTK",
      status: "normal",
      lastDetected: "2024-01-15 15:30:24",
    },
    {
      id: "3",
      frequency: 1575.42,
      strength: -65,
      type: "GPS",
      source: "GPS Satellite",
      status: "suspicious",
      lastDetected: "2024-01-15 15:30:23",
    },
    {
      id: "4",
      frequency: 2437.0,
      strength: -72,
      type: "WiFi",
      source: "Unknown AP",
      status: "threat",
      lastDetected: "2024-01-15 15:30:20",
    },
  ])

  const [rfThreats, setRfThreats] = useState<RFThreat[]>([
    {
      id: "1",
      type: "spoofing",
      frequency: 1575.42,
      severity: "high",
      description: "检测到可疑的GPS信号，可能存在GPS欺骗攻击",
      detectedAt: "2024-01-15 15:25:30",
      affectedSystems: ["DJI M300 RTK #001", "DJI M300 RTK #002"],
    },
    {
      id: "2",
      type: "jamming",
      frequency: 2400.5,
      severity: "medium",
      description: "2.4GHz频段检测到强干扰信号，可能影响遥控通信",
      detectedAt: "2024-01-15 15:20:15",
      affectedSystems: ["Ground Control Station"],
    },
  ])

  const [isRfScanning, setIsRfScanning] = useState(false)
  const [rfScanProgress, setRfScanProgress] = useState(0)

  const [newAsset, setNewAsset] = useState({
    name: "",
    type: "drone" as Asset["type"],
    ip: "",
  })

  const getAssetIcon = (type: Asset["type"]) => {
    switch (type) {
      case "drone":
        return <Zap className="h-4 w-4" />
      case "gcs":
        return <Activity className="h-4 w-4" />
      case "server":
        return <Server className="h-4 w-4" />
      case "network":
        return <Wifi className="h-4 w-4" />
    }
  }

  const getStatusColor = (status: Asset["status"]) => {
    switch (status) {
      case "online":
        return "bg-green-500"
      case "offline":
        return "bg-red-500"
      case "scanning":
        return "bg-yellow-500"
    }
  }

  const getVulnerabilityBadge = (count?: number) => {
    if (!count) return <Badge variant="secondary">安全</Badge>
    if (count <= 2)
      return (
        <Badge variant="outline" className="text-yellow-600">
          低风险
        </Badge>
      )
    return <Badge variant="destructive">高风险</Badge>
  }

  const getRfSignalIcon = (type: RFSignal["type"]) => {
    switch (type) {
      case "GPS":
        return <Satellite className="h-4 w-4" />
      case "RC":
        return <Radio className="h-4 w-4" />
      case "Video":
        return <Eye className="h-4 w-4" />
      case "Telemetry":
        return <Activity className="h-4 w-4" />
      case "WiFi":
        return <Wifi className="h-4 w-4" />
      case "Bluetooth":
        return <Bluetooth className="h-4 w-4" />
      default:
        return <Signal className="h-4 w-4" />
    }
  }

  const getRfStatusColor = (status: RFSignal["status"]) => {
    switch (status) {
      case "normal":
        return "text-green-600"
      case "suspicious":
        return "text-yellow-600"
      case "threat":
        return "text-red-600"
    }
  }

  const getRfStatusBadge = (status: RFSignal["status"]) => {
    switch (status) {
      case "normal":
        return <Badge variant="secondary">正常</Badge>
      case "suspicious":
        return (
          <Badge variant="outline" className="text-yellow-600">
            可疑
          </Badge>
        )
      case "threat":
        return <Badge variant="destructive">威胁</Badge>
    }
  }

  const getThreatSeverityColor = (severity: RFThreat["severity"]) => {
    switch (severity) {
      case "low":
        return "bg-blue-500 text-white"
      case "medium":
        return "bg-yellow-500 text-white"
      case "high":
        return "bg-orange-500 text-white"
      case "critical":
        return "bg-red-500 text-white"
    }
  }

  const getThreatTypeText = (type: RFThreat["type"]) => {
    switch (type) {
      case "jamming":
        return "信号干扰"
      case "spoofing":
        return "信号欺骗"
      case "interception":
        return "信号截获"
      case "unauthorized":
        return "未授权接入"
    }
  }

  const addAsset = () => {
    if (newAsset.name && newAsset.ip) {
      const asset: Asset = {
        id: Date.now().toString(),
        ...newAsset,
        status: "offline",
      }
      setAssets([...assets, asset])
      setNewAsset({ name: "", type: "drone", ip: "" })
    }
  }

  const startScan = (assetId: string) => {
    const asset = assets.find((a) => a.id === assetId)
    if (asset) {
      setAssets(assets.map((a) => (a.id === assetId ? { ...a, status: "scanning" } : a)))

      const task: ScanTask = {
        id: Date.now().toString(),
        assetName: asset.name,
        progress: 0,
        status: "running",
        startTime: new Date().toLocaleString(),
      }
      setScanTasks([...scanTasks, task])
    }
  }

  const startRfScan = () => {
    setIsRfScanning(true)
    setRfScanProgress(0)

    // 模拟扫描进度
    const interval = setInterval(() => {
      setRfScanProgress((prev) => {
        if (prev >= 100) {
          clearInterval(interval)
          setIsRfScanning(false)
          return 100
        }
        return prev + 10
      })
    }, 500)
  }

  const stopRfScan = () => {
    setIsRfScanning(false)
    setRfScanProgress(0)
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center gap-3 mb-2">
            <Shield className="h-8 w-8 text-blue-600" />
            <h1 className="text-3xl font-bold text-gray-900">电力无人机安全扫描系统</h1>
          </div>
          <p className="text-gray-600">专业的无人机设备安全漏洞检测与评估平台</p>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">总资产数</p>
                  <p className="text-2xl font-bold text-gray-900">{assets.length}</p>
                </div>
                <Server className="h-8 w-8 text-blue-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">在线设备</p>
                  <p className="text-2xl font-bold text-green-600">
                    {assets.filter((a) => a.status === "online").length}
                  </p>
                </div>
                <CheckCircle className="h-8 w-8 text-green-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">射频威胁</p>
                  <p className="text-2xl font-bold text-red-600">{rfThreats.length}</p>
                </div>
                <Radio className="h-8 w-8 text-red-500" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">高危漏洞</p>
                  <p className="text-2xl font-bold text-red-600">
                    {assets.reduce((sum, a) => sum + (a.vulnerabilities || 0), 0)}
                  </p>
                </div>
                <AlertTriangle className="h-8 w-8 text-red-500" />
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Main Content */}
        <Tabs defaultValue="assets" className="space-y-6">
          <TabsList className="grid w-full grid-cols-5">
            <TabsTrigger value="assets">资产管理</TabsTrigger>
            <TabsTrigger value="scanning">扫描监控</TabsTrigger>
            <TabsTrigger value="rf-security">射频安全</TabsTrigger>
            <TabsTrigger value="reports">安全报告</TabsTrigger>
            <TabsTrigger value="settings">系统设置</TabsTrigger>
          </TabsList>

          {/* Assets Tab */}
          <TabsContent value="assets" className="space-y-6">
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              {/* Asset List */}
              <div className="lg:col-span-2">
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Database className="h-5 w-5" />
                      资产清单
                    </CardTitle>
                    <CardDescription>管理所有无人机设备和基础设施资产</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      {assets.map((asset) => (
                        <div key={asset.id} className="flex items-center justify-between p-4 border rounded-lg">
                          <div className="flex items-center gap-3">
                            {getAssetIcon(asset.type)}
                            <div>
                              <h3 className="font-medium">{asset.name}</h3>
                              <p className="text-sm text-gray-500">{asset.ip}</p>
                            </div>
                          </div>

                          <div className="flex items-center gap-3">
                            <div className="flex items-center gap-2">
                              <div className={`w-2 h-2 rounded-full ${getStatusColor(asset.status)}`} />
                              <span className="text-sm capitalize">{asset.status}</span>
                            </div>

                            {getVulnerabilityBadge(asset.vulnerabilities)}

                            <Button
                              size="sm"
                              onClick={() => startScan(asset.id)}
                              disabled={asset.status === "scanning"}
                            >
                              <Search className="h-4 w-4 mr-1" />
                              扫描
                            </Button>
                          </div>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </div>

              {/* Add Asset */}
              <div>
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <Plus className="h-5 w-5" />
                      添加资产
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div>
                      <Label htmlFor="asset-name">设备名称</Label>
                      <Input
                        id="asset-name"
                        value={newAsset.name}
                        onChange={(e) => setNewAsset({ ...newAsset, name: e.target.value })}
                        placeholder="输入设备名称"
                      />
                    </div>

                    <div>
                      <Label htmlFor="asset-type">设备类型</Label>
                      <select
                        className="w-full p-2 border rounded-md"
                        value={newAsset.type}
                        onChange={(e) => setNewAsset({ ...newAsset, type: e.target.value as Asset["type"] })}
                      >
                        <option value="drone">无人机</option>
                        <option value="gcs">地面控制站</option>
                        <option value="server">服务器</option>
                        <option value="network">网络设备</option>
                      </select>
                    </div>

                    <div>
                      <Label htmlFor="asset-ip">IP地址</Label>
                      <Input
                        id="asset-ip"
                        value={newAsset.ip}
                        onChange={(e) => setNewAsset({ ...newAsset, ip: e.target.value })}
                        placeholder="192.168.1.100"
                      />
                    </div>

                    <Button onClick={addAsset} className="w-full">
                      <Plus className="h-4 w-4 mr-2" />
                      添加资产
                    </Button>
                  </CardContent>
                </Card>
              </div>
            </div>
          </TabsContent>

          {/* Scanning Tab */}
          <TabsContent value="scanning" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Activity className="h-5 w-5" />
                  扫描任务监控
                </CardTitle>
                <CardDescription>实时监控正在进行的安全扫描任务</CardDescription>
              </CardHeader>
              <CardContent>
                {scanTasks.length > 0 ? (
                  <div className="space-y-4">
                    {scanTasks.map((task) => (
                      <div key={task.id} className="p-4 border rounded-lg">
                        <div className="flex items-center justify-between mb-2">
                          <h3 className="font-medium">{task.assetName}</h3>
                          <Badge variant={task.status === "running" ? "default" : "secondary"}>
                            {task.status === "running" ? "扫描中" : "已完成"}
                          </Badge>
                        </div>
                        <div className="space-y-2">
                          <div className="flex justify-between text-sm text-gray-600">
                            <span>进度: {task.progress}%</span>
                            <span>开始时间: {task.startTime}</span>
                          </div>
                          <Progress value={task.progress} className="w-full" />
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8 text-gray-500">
                    <Activity className="h-12 w-12 mx-auto mb-4 opacity-50" />
                    <p>当前没有正在进行的扫描任务</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>

          {/* RF Security Tab */}
          <TabsContent value="rf-security" className="space-y-6">
            {/* RF Scan Control */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Radar className="h-5 w-5" />
                  射频频谱扫描
                </CardTitle>
                <CardDescription>监控和分析无人机通信频段的射频信号</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex items-center gap-4 mb-4">
                  <Button
                    onClick={isRfScanning ? stopRfScan : startRfScan}
                    className={isRfScanning ? "bg-red-600 hover:bg-red-700" : ""}
                  >
                    {isRfScanning ? (
                      <>
                        <Pause className="h-4 w-4 mr-2" />
                        停止扫描
                      </>
                    ) : (
                      <>
                        <Play className="h-4 w-4 mr-2" />
                        开始扫描
                      </>
                    )}
                  </Button>
                  <Button variant="outline" onClick={() => setRfScanProgress(0)}>
                    <RotateCcw className="h-4 w-4 mr-2" />
                    重置
                  </Button>
                  <div className="flex items-center gap-2">
                    <Switch />
                    <Label>自动扫描</Label>
                  </div>
                </div>

                {isRfScanning && (
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>扫描进度: {rfScanProgress}%</span>
                      <span>频段: 2.4GHz - 5.8GHz</span>
                    </div>
                    <Progress value={rfScanProgress} className="w-full" />
                  </div>
                )}
              </CardContent>
            </Card>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* RF Signals */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Signal className="h-5 w-5" />
                    射频信号监测
                  </CardTitle>
                  <CardDescription>实时监测检测到的射频信号</CardDescription>
                </CardHeader>
                <CardContent>
                  <ScrollArea className="h-80">
                    <div className="space-y-3">
                      {rfSignals.map((signal) => (
                        <div key={signal.id} className="border rounded-lg p-3">
                          <div className="flex items-center justify-between mb-2">
                            <div className="flex items-center gap-2">
                              {getRfSignalIcon(signal.type)}
                              <span className="font-medium text-sm">{signal.type}</span>
                              {getRfStatusBadge(signal.status)}
                            </div>
                            <span className="text-xs text-gray-500">{signal.lastDetected}</span>
                          </div>

                          <div className="text-sm space-y-1">
                            <div className="flex justify-between">
                              <span className="text-gray-600">频率:</span>
                              <span className="font-mono">{signal.frequency} MHz</span>
                            </div>
                            <div className="flex justify-between">
                              <span className="text-gray-600">信号强度:</span>
                              <span className={`font-mono ${getRfStatusColor(signal.status)}`}>
                                {signal.strength} dBm
                              </span>
                            </div>
                            {signal.source && (
                              <div className="flex justify-between">
                                <span className="text-gray-600">信号源:</span>
                                <span className="text-xs">{signal.source}</span>
                              </div>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  </ScrollArea>
                </CardContent>
              </Card>

              {/* RF Threats */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <AlertTriangle className="h-5 w-5" />
                    射频安全威胁
                  </CardTitle>
                  <CardDescription>检测到的射频安全威胁和异常</CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {rfThreats.map((threat) => (
                      <div key={threat.id} className="border rounded-lg p-4">
                        <div className="flex items-start justify-between mb-3">
                          <div className="flex items-center gap-2">
                            <Badge className={getThreatSeverityColor(threat.severity)}>
                              {threat.severity.toUpperCase()}
                            </Badge>
                            <span className="font-medium text-sm">{getThreatTypeText(threat.type)}</span>
                          </div>
                          <span className="text-xs text-gray-500">{threat.detectedAt}</span>
                        </div>

                        <p className="text-sm text-gray-700 mb-3">{threat.description}</p>

                        <div className="text-sm space-y-2">
                          <div className="flex justify-between">
                            <span className="text-gray-600">频率:</span>
                            <span className="font-mono">{threat.frequency} MHz</span>
                          </div>
                          <div>
                            <span className="text-gray-600">受影响系统:</span>
                            <div className="mt-1 space-y-1">
                              {threat.affectedSystems.map((system, index) => (
                                <Badge key={index} variant="outline" className="text-xs mr-1">
                                  {system}
                                </Badge>
                              ))}
                            </div>
                          </div>
                        </div>

                        <div className="flex gap-2 mt-3">
                          <Button size="sm" variant="outline">
                            <Eye className="h-3 w-3 mr-1" />
                            详情
                          </Button>
                          <Button size="sm" variant="outline">
                            <WifiOff className="h-3 w-3 mr-1" />
                            屏蔽
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* RF Analysis Tools */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Volume2 className="h-5 w-5" />
                  射频分析工具
                </CardTitle>
                <CardDescription>专业的射频信号分析和检测工具</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                  <Button variant="outline" className="h-20 flex-col bg-transparent">
                    <Satellite className="h-6 w-6 mb-2" />
                    <span className="text-sm">GPS欺骗检测</span>
                  </Button>
                  <Button variant="outline" className="h-20 flex-col bg-transparent">
                    <Radio className="h-6 w-6 mb-2" />
                    <span className="text-sm">遥控信号分析</span>
                  </Button>
                  <Button variant="outline" className="h-20 flex-col bg-transparent">
                    <Wifi className="h-6 w-6 mb-2" />
                    <span className="text-sm">WiFi安全检测</span>
                  </Button>
                  <Button variant="outline" className="h-20 flex-col bg-transparent">
                    <Bluetooth className="h-6 w-6 mb-2" />
                    <span className="text-sm">蓝牙安全扫描</span>
                  </Button>
                </div>

                <div className="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <Label className="text-sm font-medium">扫描频率范围</Label>
                    <div className="flex gap-2 mt-1">
                      <Input placeholder="起始频率 (MHz)" className="text-sm" />
                      <Input placeholder="结束频率 (MHz)" className="text-sm" />
                    </div>
                  </div>
                  <div>
                    <Label className="text-sm font-medium">检测灵敏度</Label>
                    <select className="w-full p-2 border rounded-md mt-1 text-sm">
                      <option>高灵敏度 (-90 dBm)</option>
                      <option>标准灵敏度 (-80 dBm)</option>
                      <option>低灵敏度 (-70 dBm)</option>
                    </select>
                  </div>
                </div>

                <div className="flex gap-2 mt-4">
                  <Button>
                    <Download className="h-4 w-4 mr-2" />
                    导出频谱数据
                  </Button>
                  <Button variant="outline">
                    <Upload className="h-4 w-4 mr-2" />
                    导入配置
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Reports Tab */}
          <TabsContent value="reports" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <FileText className="h-5 w-5" />
                  安全报告
                </CardTitle>
                <CardDescription>查看和下载安全扫描报告</CardDescription>
              </CardHeader>
              <CardContent>
                <Alert>
                  <AlertTriangle className="h-4 w-4" />
                  <AlertDescription>
                    发现 {assets.reduce((sum, a) => sum + (a.vulnerabilities || 0), 0)} 个安全漏洞和 {rfThreats.length}{" "}
                    个射频威胁需要处理， 建议立即查看详细报告并制定修复计划。
                  </AlertDescription>
                </Alert>

                <div className="mt-6 space-y-4">
                  <div className="p-4 border rounded-lg">
                    <div className="flex items-center justify-between">
                      <div>
                        <h3 className="font-medium">网络基础设施安全报告</h3>
                        <p className="text-sm text-gray-500">2024-01-15 生成</p>
                      </div>
                      <Button variant="outline">
                        <FileText className="h-4 w-4 mr-2" />
                        下载报告
                      </Button>
                    </div>
                  </div>

                  <div className="p-4 border rounded-lg">
                    <div className="flex items-center justify-between">
                      <div>
                        <h3 className="font-medium">射频安全威胁分析报告</h3>
                        <p className="text-sm text-gray-500">2024-01-15 生成</p>
                      </div>
                      <Button variant="outline">
                        <FileText className="h-4 w-4 mr-2" />
                        下载报告
                      </Button>
                    </div>
                  </div>

                  <div className="p-4 border rounded-lg">
                    <div className="flex items-center justify-between">
                      <div>
                        <h3 className="font-medium">无人机设备安全评估</h3>
                        <p className="text-sm text-gray-500">2024-01-14 生成</p>
                      </div>
                      <Button variant="outline">
                        <FileText className="h-4 w-4 mr-2" />
                        下载报告
                      </Button>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Settings Tab */}
          <TabsContent value="settings" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>系统设置</CardTitle>
                <CardDescription>配置扫描参数和系统选项</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label>扫描深度</Label>
                  <select className="w-full p-2 border rounded-md mt-1">
                    <option>快速扫描</option>
                    <option>标准扫描</option>
                    <option>深度扫描</option>
                  </select>
                </div>

                <div>
                  <Label>报告格式</Label>
                  <select className="w-full p-2 border rounded-md mt-1">
                    <option>PDF</option>
                    <option>HTML</option>
                    <option>JSON</option>
                  </select>
                </div>

                <Button>保存设置</Button>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  )
=======
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
    Shield,
    Database,
    Radio,
    AlertTriangle,
    CheckCircle,
    Search,
    Plus,
    FileText,
    Users,
    Activity,
} from "lucide-react";
import { AssetList } from "@/components/asset-list";
import { AddAssetForm } from "@/components/add-asset-form";
import { SecurityScanResults } from "@/components/security-scan-results";
import { RfAnalysisPanel } from "@/components/rf-analysis-panel";
import { PowerIndustryRisks } from "@/components/power-industry-risks";
import { DroneDataManagement } from "@/components/drone-data-management";
import { TaskManagementSystem } from "@/components/task-management-system";

export default function Dashboard() {
    return (
        <div className="flex min-h-screen flex-col bg-white">
            <header className="border-b py-4">
                <div className="container">
                    <div className="flex items-center gap-2">
                        <Shield className="h-6 w-6 text-blue-600" />
                        <div>
                            <h1 className="text-xl font-bold">
                                电力无人机安全扫描系统
                            </h1>
                            <p className="text-sm text-muted-foreground">
                                专业的无人机设备安全漏洞检测与评估平台
                            </p>
                        </div>
                    </div>
                </div>
            </header>
            <main className="flex-1 py-6">
                <div className="container">
                    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
                        <Card>
                            <CardContent className="flex items-center justify-between p-6">
                                <div>
                                    <p className="text-sm text-muted-foreground">
                                        总资产数
                                    </p>
                                    <p className="text-3xl font-bold text-blue-600">
                                        4
                                    </p>
                                </div>
                                <div className="rounded-full bg-blue-100 p-2">
                                    <Database className="h-5 w-5 text-blue-600" />
                                </div>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardContent className="flex items-center justify-between p-6">
                                <div>
                                    <p className="text-sm text-muted-foreground">
                                        在线设备
                                    </p>
                                    <p className="text-3xl font-bold text-green-600">
                                        2
                                    </p>
                                </div>
                                <div className="rounded-full bg-green-100 p-2">
                                    <CheckCircle className="h-5 w-5 text-green-600" />
                                </div>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardContent className="flex items-center justify-between p-6">
                                <div>
                                    <p className="text-sm text-muted-foreground">
                                        射频威胁
                                    </p>
                                    <p className="text-3xl font-bold text-amber-600">
                                        2
                                    </p>
                                </div>
                                <div className="rounded-full bg-amber-100 p-2">
                                    <Radio className="h-5 w-5 text-amber-600" />
                                </div>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardContent className="flex items-center justify-between p-6">
                                <div>
                                    <p className="text-sm text-muted-foreground">
                                        高危漏洞
                                    </p>
                                    <p className="text-3xl font-bold text-red-600">
                                        6
                                    </p>
                                </div>
                                <div className="rounded-full bg-red-100 p-2">
                                    <AlertTriangle className="h-5 w-5 text-red-600" />
                                </div>
                            </CardContent>
                        </Card>
                    </div>

                    <div className="mt-6">
                        <Tabs defaultValue="assets" className="w-full">
                            <TabsList className="grid w-full grid-cols-8">
                                <TabsTrigger value="assets">
                                    资产管理
                                </TabsTrigger>
                                <TabsTrigger value="scan">扫描监控</TabsTrigger>
                                <TabsTrigger value="rf">射频安全</TabsTrigger>
                                <TabsTrigger value="risks">
                                    风险评估
                                </TabsTrigger>
                                <TabsTrigger value="data">数据管理</TabsTrigger>
                                <TabsTrigger value="tasks">
                                    任务管理
                                </TabsTrigger>
                                <TabsTrigger value="reports">
                                    安全报告
                                </TabsTrigger>
                                <TabsTrigger value="settings">
                                    系统设置
                                </TabsTrigger>
                            </TabsList>

                            <TabsContent value="assets" className="mt-6">
                                <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
                                    <div className="lg:col-span-2">
                                        <div className="rounded-md border">
                                            <div className="flex items-center gap-2 border-b p-4">
                                                <Database className="h-5 w-5" />
                                                <h2 className="font-medium">
                                                    资产清单
                                                </h2>
                                            </div>
                                            <div className="p-2">
                                                <p className="px-2 py-1 text-sm text-muted-foreground">
                                                    管理所有无人机设备和相关设备资产
                                                </p>
                                                <AssetList />
                                            </div>
                                        </div>
                                    </div>
                                    <div>
                                        <div className="rounded-md border">
                                            <div className="flex items-center gap-2 border-b p-4">
                                                <Plus className="h-5 w-5" />
                                                <h2 className="font-medium">
                                                    添加资产
                                                </h2>
                                            </div>
                                            <div className="p-4">
                                                <AddAssetForm />
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="scan" className="mt-6">
                                <div className="rounded-md border">
                                    <div className="flex items-center justify-between border-b p-4">
                                        <div className="flex items-center gap-2">
                                            <Search className="h-5 w-5" />
                                            <h2 className="font-medium">
                                                安全扫描结果
                                            </h2>
                                        </div>
                                        <Button>
                                            <Search className="mr-2 h-4 w-4" />
                                            开始扫描
                                        </Button>
                                    </div>
                                    <div className="p-4">
                                        <SecurityScanResults />
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="rf" className="mt-6">
                                <div className="rounded-md border">
                                    <div className="flex items-center justify-between border-b p-4">
                                        <div className="flex items-center gap-2">
                                            <Radio className="h-5 w-5" />
                                            <h2 className="font-medium">
                                                射频安全分析
                                            </h2>
                                        </div>
                                        <Button>
                                            <Radio className="mr-2 h-4 w-4" />
                                            开始射频分析
                                        </Button>
                                    </div>
                                    <div className="p-4">
                                        <RfAnalysisPanel />
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="risks" className="mt-6">
                                <div className="rounded-md border">
                                    <div className="flex items-center gap-2 border-b p-4">
                                        <AlertTriangle className="h-5 w-5" />
                                        <h2 className="font-medium">
                                            电力行业风险评估
                                        </h2>
                                    </div>
                                    <div className="p-4">
                                        <PowerIndustryRisks />
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="data" className="mt-6">
                                <div className="rounded-md border">
                                    <div className="flex items-center gap-2 border-b p-4">
                                        <Database className="h-5 w-5" />
                                        <h2 className="font-medium">
                                            无人机数据管理
                                        </h2>
                                    </div>
                                    <div className="p-4">
                                        <DroneDataManagement />
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="tasks" className="mt-6">
                                <div className="rounded-md border">
                                    <div className="flex items-center gap-2 border-b p-4">
                                        <FileText className="h-5 w-5" />
                                        <h2 className="font-medium">
                                            任务管理系统
                                        </h2>
                                    </div>
                                    <div className="p-4">
                                        <TaskManagementSystem />
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="reports" className="mt-6">
                                <div className="rounded-md border p-6">
                                    <div className="flex flex-col items-center justify-center py-10">
                                        <div className="rounded-full bg-blue-100 p-4">
                                            <Shield className="h-8 w-8 text-blue-600" />
                                        </div>
                                        <h3 className="mt-4 text-lg font-medium">
                                            安全报告
                                        </h3>
                                        <p className="mt-2 text-center text-sm text-muted-foreground max-w-md">
                                            生成详细的安全评估报告，包括网络安全、数据安全和射频安全风险分析
                                        </p>
                                        <Button className="mt-4">
                                            生成安全报告
                                        </Button>
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="settings" className="mt-6">
                                <div className="rounded-md border p-6">
                                    <div className="flex flex-col items-center justify-center py-10">
                                        <h3 className="text-lg font-medium">
                                            系统设置
                                        </h3>
                                        <p className="mt-2 text-center text-sm text-muted-foreground max-w-md">
                                            配置系统参数、扫描策略和安全评估标准
                                        </p>
                                        <Button
                                            variant="outline"
                                            className="mt-4 bg-transparent"
                                        >
                                            配置系统
                                        </Button>
                                    </div>
                                </div>
                            </TabsContent>
                        </Tabs>
                    </div>
                </div>
            </main>
            <footer className="border-t py-4">
                <div className="container">
                    <div className="flex items-center justify-between">
                        <p className="text-sm text-muted-foreground">
                            © 2025 电力无人机安全扫描系统
                        </p>
                        <p className="text-sm text-muted-foreground">
                            版本 1.0.0
                        </p>
                    </div>
                </div>
            </footer>
        </div>
    );
>>>>>>> 8a1195963e349af332e227e3f3e20bb08506797d
}
