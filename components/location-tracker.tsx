"use client"

import { useState, useEffect, useRef } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Label } from "@/components/ui/label"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Switch } from "@/components/ui/switch"
import {
  MapPin,
  Navigation,
  Satellite,
  AlertTriangle,
  CheckCircle,
  Clock,
  Target,
  Route,
  Zap,
  Activity,
  Server,
  Wifi,
  RefreshCw,
  Save,
  Eye,
  Settings,
  Shield,
  Camera,
} from "lucide-react"

interface LocationData {
  latitude: number
  longitude: number
  accuracy: number
  altitude?: number
  heading?: number
  speed?: number
  timestamp: number
}

interface DeviceLocation {
  id: string
  name: string
  type: "drone" | "gcs" | "server" | "network" | "hangar"
  location: LocationData
  status: "online" | "offline" | "maintenance"
  lastUpdate: string
  address?: string
  zone?: string
}

interface ProcessingRecord {
  id: string
  vulnerabilityId: string
  vulnerabilityTitle: string
  location: LocationData
  address: string
  processingTime: string
  engineer: string
  status: "in-progress" | "completed"
  notes?: string
  photos?: string[]
}

export default function LocationTracker() {
  const [currentLocation, setCurrentLocation] = useState<LocationData | null>(null)
  const [locationError, setLocationError] = useState<string | null>(null)
  const [isTracking, setIsTracking] = useState(false)
  const [locationHistory, setLocationHistory] = useState<LocationData[]>([])
  const [deviceLocations, setDeviceLocations] = useState<DeviceLocation[]>([
    {
      id: "drone-001",
      name: "DJI M300 RTK #001",
      type: "drone",
      location: {
        latitude: 39.9042,
        longitude: 116.4074,
        accuracy: 5,
        altitude: 120,
        timestamp: Date.now(),
      },
      status: "online",
      lastUpdate: "2024-01-15 14:30",
      address: "北京市朝阳区电力巡检站",
      zone: "华北电网",
    },
    {
      id: "hangar-001",
      name: "无人机机库 #001",
      type: "hangar",
      location: {
        latitude: 39.9052,
        longitude: 116.4084,
        accuracy: 2,
        timestamp: Date.now(),
      },
      status: "online",
      lastUpdate: "2024-01-15 14:25",
      address: "北京市朝阳区电力基地",
      zone: "华北电网",
    },
  ])
  const [processingRecords, setProcessingRecords] = useState<ProcessingRecord[]>([])
  const [autoTracking, setAutoTracking] = useState(false)
  const [geoFencing, setGeoFencing] = useState(true)

  const watchIdRef = useRef<number | null>(null)

  useEffect(() => {
    if (autoTracking) {
      startLocationTracking()
    } else {
      stopLocationTracking()
    }

    return () => {
      if (watchIdRef.current) {
        navigator.geolocation.clearWatch(watchIdRef.current)
      }
    }
  }, [autoTracking])

  const startLocationTracking = () => {
    if (!navigator.geolocation) {
      setLocationError("此设备不支持GPS定位")
      return
    }

    setIsTracking(true)
    setLocationError(null)

    const options = {
      enableHighAccuracy: true,
      timeout: 10000,
      maximumAge: 60000,
    }

    watchIdRef.current = navigator.geolocation.watchPosition(
      (position) => {
        const locationData: LocationData = {
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
          accuracy: position.coords.accuracy,
          altitude: position.coords.altitude || undefined,
          heading: position.coords.heading || undefined,
          speed: position.coords.speed || undefined,
          timestamp: Date.now(),
        }

        setCurrentLocation(locationData)
        setLocationHistory((prev) => [...prev.slice(-99), locationData]) // 保留最近100个位置点
        setLocationError(null)
      },
      (error) => {
        setLocationError(getLocationErrorMessage(error))
        setIsTracking(false)
      },
      options,
    )
  }

  const stopLocationTracking = () => {
    if (watchIdRef.current) {
      navigator.geolocation.clearWatch(watchIdRef.current)
      watchIdRef.current = null
    }
    setIsTracking(false)
  }

  const getCurrentLocation = () => {
    if (!navigator.geolocation) {
      setLocationError("此设备不支持GPS定位")
      return
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        const locationData: LocationData = {
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
          accuracy: position.coords.accuracy,
          altitude: position.coords.altitude || undefined,
          heading: position.coords.heading || undefined,
          speed: position.coords.speed || undefined,
          timestamp: Date.now(),
        }

        setCurrentLocation(locationData)
        setLocationError(null)
      },
      (error) => {
        setLocationError(getLocationErrorMessage(error))
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 60000,
      },
    )
  }

  const getLocationErrorMessage = (error: GeolocationPositionError) => {
    switch (error.code) {
      case error.PERMISSION_DENIED:
        return "用户拒绝了定位请求"
      case error.POSITION_UNAVAILABLE:
        return "位置信息不可用"
      case error.TIMEOUT:
        return "定位请求超时"
      default:
        return "获取位置时发生未知错误"
    }
  }

  const formatCoordinate = (coord: number, type: "lat" | "lng") => {
    const direction = type === "lat" ? (coord >= 0 ? "N" : "S") : coord >= 0 ? "E" : "W"
    return `${Math.abs(coord).toFixed(6)}° ${direction}`
  }

  const formatDistance = (meters: number) => {
    if (meters < 1000) {
      return `${meters.toFixed(0)}m`
    }
    return `${(meters / 1000).toFixed(2)}km`
  }

  const calculateDistance = (lat1: number, lon1: number, lat2: number, lon2: number) => {
    const R = 6371e3 // 地球半径（米）
    const φ1 = (lat1 * Math.PI) / 180
    const φ2 = (lat2 * Math.PI) / 180
    const Δφ = ((lat2 - lat1) * Math.PI) / 180
    const Δλ = ((lon2 - lon1) * Math.PI) / 180

    const a = Math.sin(Δφ / 2) * Math.sin(Δφ / 2) + Math.cos(φ1) * Math.cos(φ2) * Math.sin(Δλ / 2) * Math.sin(Δλ / 2)
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))

    return R * c
  }

  const recordProcessingLocation = (vulnerabilityId: string, vulnerabilityTitle: string) => {
    if (!currentLocation) {
      setLocationError("请先获取当前位置")
      return
    }

    const record: ProcessingRecord = {
      id: Date.now().toString(),
      vulnerabilityId,
      vulnerabilityTitle,
      location: currentLocation,
      address: "正在解析地址...", // 实际应用中会调用逆地理编码API
      processingTime: new Date().toLocaleString(),
      engineer: "当前用户", // 实际应用中从用户上下文获取
      status: "in-progress",
    }

    setProcessingRecords((prev) => [record, ...prev])
  }

  const getDeviceIcon = (type: DeviceLocation["type"]) => {
    switch (type) {
      case "drone":
        return <Zap className="h-4 w-4" />
      case "gcs":
        return <Activity className="h-4 w-4" />
      case "server":
        return <Server className="h-4 w-4" />
      case "network":
        return <Wifi className="h-4 w-4" />
      case "hangar":
        return <Shield className="h-4 w-4" />
    }
  }

  const getStatusColor = (status: DeviceLocation["status"]) => {
    switch (status) {
      case "online":
        return "bg-green-500"
      case "offline":
        return "bg-red-500"
      case "maintenance":
        return "bg-yellow-500"
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center gap-3 mb-2">
            <MapPin className="h-8 w-8 text-blue-600" />
            <h1 className="text-3xl font-bold text-gray-900">位置追踪与设备管理</h1>
          </div>
          <p className="text-gray-600">实时追踪处理位置和设备地理信息</p>
        </div>

        {/* Current Location Status */}
        <Card className="mb-6">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Navigation className="h-5 w-5" />
              当前位置状态
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <div className="flex items-center gap-4 mb-4">
                  <Button onClick={getCurrentLocation} disabled={isTracking}>
                    <Satellite className="h-4 w-4 mr-2" />
                    获取位置
                  </Button>
                  <div className="flex items-center gap-2">
                    <Switch checked={autoTracking} onCheckedChange={setAutoTracking} />
                    <Label>自动追踪</Label>
                  </div>
                </div>

                {locationError && (
                  <Alert className="mb-4">
                    <AlertTriangle className="h-4 w-4" />
                    <AlertDescription>{locationError}</AlertDescription>
                  </Alert>
                )}

                {currentLocation && (
                  <div className="space-y-2">
                    <div className="flex items-center gap-2">
                      <CheckCircle className="h-4 w-4 text-green-500" />
                      <span className="text-sm text-green-600">定位成功</span>
                    </div>
                    <div className="text-sm space-y-1">
                      <div>纬度: {formatCoordinate(currentLocation.latitude, "lat")}</div>
                      <div>经度: {formatCoordinate(currentLocation.longitude, "lng")}</div>
                      <div>精度: ±{currentLocation.accuracy.toFixed(0)}m</div>
                      {currentLocation.altitude && <div>海拔: {currentLocation.altitude.toFixed(0)}m</div>}
                      {currentLocation.speed && <div>速度: {(currentLocation.speed * 3.6).toFixed(1)}km/h</div>}
                    </div>
                  </div>
                )}
              </div>

              <div>
                <div className="flex items-center gap-4 mb-4">
                  <div className="flex items-center gap-2">
                    <Switch checked={geoFencing} onCheckedChange={setGeoFencing} />
                    <Label>地理围栏</Label>
                  </div>
                  <Button variant="outline" size="sm">
                    <Settings className="h-4 w-4 mr-2" />
                    设置
                  </Button>
                </div>

                <div className="text-sm space-y-2">
                  <div className="flex justify-between">
                    <span>追踪状态:</span>
                    <Badge variant={isTracking ? "default" : "secondary"}>{isTracking ? "追踪中" : "已停止"}</Badge>
                  </div>
                  <div className="flex justify-between">
                    <span>历史点数:</span>
                    <span>{locationHistory.length}</span>
                  </div>
                  <div className="flex justify-between">
                    <span>最后更新:</span>
                    <span>{currentLocation ? new Date(currentLocation.timestamp).toLocaleTimeString() : "无"}</span>
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Main Content */}
        <Tabs defaultValue="devices" className="space-y-6">
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="devices">设备位置</TabsTrigger>
            <TabsTrigger value="processing">处理记录</TabsTrigger>
            <TabsTrigger value="tracking">路径追踪</TabsTrigger>
            <TabsTrigger value="map">地图视图</TabsTrigger>
          </TabsList>

          {/* Devices Tab */}
          <TabsContent value="devices" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Target className="h-5 w-5" />
                  设备地理位置
                </CardTitle>
                <CardDescription>管理所有设备的地理位置信息</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {deviceLocations.map((device) => (
                    <div key={device.id} className="border rounded-lg p-4">
                      <div className="flex items-start justify-between">
                        <div className="flex items-start gap-3">
                          {getDeviceIcon(device.type)}
                          <div className="flex-1">
                            <div className="flex items-center gap-2 mb-1">
                              <h3 className="font-medium">{device.name}</h3>
                              <div className={`w-2 h-2 rounded-full ${getStatusColor(device.status)}`} />
                            </div>
                            <div className="text-sm text-gray-600 space-y-1">
                              <div>{device.address}</div>
                              <div className="flex items-center gap-4">
                                <span>
                                  {formatCoordinate(device.location.latitude, "lat")},{" "}
                                  {formatCoordinate(device.location.longitude, "lng")}
                                </span>
                                <span>精度: ±{device.location.accuracy}m</span>
                              </div>
                              <div className="flex items-center gap-4">
                                <span>区域: {device.zone}</span>
                                <span>更新: {device.lastUpdate}</span>
                              </div>
                            </div>
                          </div>
                        </div>

                        <div className="flex items-center gap-2">
                          {currentLocation && (
                            <Badge variant="outline" className="text-xs">
                              {formatDistance(
                                calculateDistance(
                                  currentLocation.latitude,
                                  currentLocation.longitude,
                                  device.location.latitude,
                                  device.location.longitude,
                                ),
                              )}
                            </Badge>
                          )}
                          <Button size="sm" variant="outline">
                            <Eye className="h-4 w-4 mr-1" />
                            查看
                          </Button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Processing Records Tab */}
          <TabsContent value="processing" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Clock className="h-5 w-5" />
                  处理位置记录
                </CardTitle>
                <CardDescription>记录漏洞处理的地理位置信息</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="mb-4">
                  <Button
                    onClick={() => recordProcessingLocation("vuln-001", "DJI 无人机飞控系统未授权访问漏洞")}
                    disabled={!currentLocation}
                  >
                    <MapPin className="h-4 w-4 mr-2" />
                    记录当前处理位置
                  </Button>
                </div>

                <div className="space-y-4">
                  {processingRecords.length > 0 ? (
                    processingRecords.map((record) => (
                      <div key={record.id} className="border rounded-lg p-4">
                        <div className="flex items-start justify-between mb-3">
                          <div>
                            <h3 className="font-medium mb-1">{record.vulnerabilityTitle}</h3>
                            <div className="text-sm text-gray-600">
                              <div>处理人员: {record.engineer}</div>
                              <div>处理时间: {record.processingTime}</div>
                            </div>
                          </div>
                          <Badge variant={record.status === "completed" ? "secondary" : "default"}>
                            {record.status === "completed" ? "已完成" : "处理中"}
                          </Badge>
                        </div>

                        <div className="text-sm space-y-1 mb-3">
                          <div className="flex items-center gap-2">
                            <MapPin className="h-4 w-4 text-gray-400" />
                            <span>
                              {formatCoordinate(record.location.latitude, "lat")},{" "}
                              {formatCoordinate(record.location.longitude, "lng")}
                            </span>
                          </div>
                          <div className="text-gray-600 ml-6">{record.address}</div>
                          <div className="text-gray-600 ml-6">精度: ±{record.location.accuracy}m</div>
                        </div>

                        <div className="flex items-center gap-2">
                          <Button size="sm" variant="outline">
                            <Camera className="h-4 w-4 mr-1" />
                            添加照片
                          </Button>
                          <Button size="sm" variant="outline">
                            <Save className="h-4 w-4 mr-1" />
                            保存备注
                          </Button>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-8 text-gray-500">
                      <MapPin className="h-12 w-12 mx-auto mb-4 opacity-50" />
                      <p>暂无处理位置记录</p>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Tracking Tab */}
          <TabsContent value="tracking" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Route className="h-5 w-5" />
                  路径追踪
                </CardTitle>
                <CardDescription>查看移动轨迹和路径分析</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
                  <div className="text-center">
                    <div className="text-2xl font-bold text-blue-600">{locationHistory.length}</div>
                    <div className="text-sm text-gray-600">追踪点数</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-green-600">
                      {locationHistory.length > 1
                        ? formatDistance(
                            locationHistory.reduce((total, point, index) => {
                              if (index === 0) return 0
                              const prev = locationHistory[index - 1]
                              return (
                                total +
                                calculateDistance(prev.latitude, prev.longitude, point.latitude, point.longitude)
                              )
                            }, 0),
                          )
                        : "0m"}
                    </div>
                    <div className="text-sm text-gray-600">总距离</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-purple-600">
                      {locationHistory.length > 0
                        ? Math.round((Date.now() - locationHistory[0].timestamp) / (1000 * 60)) + "分钟"
                        : "0分钟"}
                    </div>
                    <div className="text-sm text-gray-600">追踪时长</div>
                  </div>
                </div>

                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <h4 className="font-medium">最近轨迹点</h4>
                    <Button size="sm" variant="outline" onClick={() => setLocationHistory([])}>
                      <RefreshCw className="h-4 w-4 mr-1" />
                      清除历史
                    </Button>
                  </div>

                  <ScrollArea className="h-64 border rounded">
                    <div className="p-4 space-y-2">
                      {locationHistory
                        .slice(-20)
                        .reverse()
                        .map((point, index) => (
                          <div key={point.timestamp} className="flex items-center justify-between text-sm">
                            <div>
                              <span className="font-mono">
                                {formatCoordinate(point.latitude, "lat")}, {formatCoordinate(point.longitude, "lng")}
                              </span>
                            </div>
                            <div className="text-gray-500">{new Date(point.timestamp).toLocaleTimeString()}</div>
                          </div>
                        ))}
                    </div>
                  </ScrollArea>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Map Tab */}
          <TabsContent value="map" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <MapPin className="h-5 w-5" />
                  地图视图
                </CardTitle>
                <CardDescription>在地图上查看所有设备和位置信息</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="h-96 bg-gray-100 rounded-lg flex items-center justify-center">
                  <div className="text-center text-gray-500">
                    <MapPin className="h-12 w-12 mx-auto mb-4 opacity-50" />
                    <p className="mb-2">地图组件</p>
                    <p className="text-sm">
                      实际应用中会集成高德地图、百度地图或其他地图服务
                      <br />
                      显示设备位置、处理记录和移动轨迹
                    </p>
                  </div>
                </div>

                <div className="mt-4 grid grid-cols-2 md:grid-cols-4 gap-4">
                  <Button variant="outline" className="justify-start bg-transparent">
                    <Zap className="h-4 w-4 mr-2" />
                    无人机
                  </Button>
                  <Button variant="outline" className="justify-start bg-transparent">
                    <Shield className="h-4 w-4 mr-2" />
                    机库
                  </Button>
                  <Button variant="outline" className="justify-start bg-transparent">
                    <Server className="h-4 w-4 mr-2" />
                    服务器
                  </Button>
                  <Button variant="outline" className="justify-start bg-transparent">
                    <MapPin className="h-4 w-4 mr-2" />
                    处理点
                  </Button>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  )
}
