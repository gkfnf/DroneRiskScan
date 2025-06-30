"use client"

import { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Alert, AlertDescription } from "@/components/ui/alert"
import { Switch } from "@/components/ui/switch"
import { Label } from "@/components/ui/label"
import {
  MapPin,
  Navigation,
  Satellite,
  AlertTriangle,
  CheckCircle,
  Target,
  Zap,
  Activity,
  Server,
  Shield,
  Camera,
  Save,
  RefreshCw,
} from "lucide-react"

interface LocationData {
  latitude: number
  longitude: number
  accuracy: number
  timestamp: number
}

export default function MobileLocationTracker() {
  const [currentLocation, setCurrentLocation] = useState<LocationData | null>(null)
  const [locationError, setLocationError] = useState<string | null>(null)
  const [isTracking, setIsTracking] = useState(false)
  const [autoTracking, setAutoTracking] = useState(false)

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
          timestamp: Date.now(),
        }

        setCurrentLocation(locationData)
        setLocationError(null)
      },
      (error) => {
        setLocationError("获取位置失败，请检查定位权限")
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 60000,
      },
    )
  }

  const formatCoordinate = (coord: number, type: "lat" | "lng") => {
    const direction = type === "lat" ? (coord >= 0 ? "N" : "S") : coord >= 0 ? "E" : "W"
    return `${Math.abs(coord).toFixed(6)}° ${direction}`
  }

  const recordCurrentLocation = () => {
    if (!currentLocation) {
      setLocationError("请先获取当前位置")
      return
    }

    // 记录处理位置的逻辑
    console.log("记录位置:", currentLocation)
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Mobile Header */}
      <div className="sticky top-0 z-50 bg-white border-b shadow-sm p-4">
        <h1 className="text-xl font-bold flex items-center gap-2">
          <MapPin className="h-6 w-6 text-blue-600" />
          位置追踪
        </h1>
      </div>

      {/* Current Location Card */}
      <div className="p-4">
        <Card className="mb-4">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <Navigation className="h-5 w-5" />
              当前位置
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center gap-4">
                <Button onClick={getCurrentLocation} disabled={isTracking} className="flex-1">
                  <Satellite className="h-4 w-4 mr-2" />
                  获取位置
                </Button>
                <div className="flex items-center gap-2">
                  <Switch checked={autoTracking} onCheckedChange={setAutoTracking} />
                  <Label className="text-sm">自动</Label>
                </div>
              </div>

              {locationError && (
                <Alert>
                  <AlertTriangle className="h-4 w-4" />
                  <AlertDescription className="text-sm">{locationError}</AlertDescription>
                </Alert>
              )}

              {currentLocation && (
                <div className="space-y-3">
                  <div className="flex items-center gap-2">
                    <CheckCircle className="h-4 w-4 text-green-500" />
                    <span className="text-sm text-green-600 font-medium">定位成功</span>
                  </div>

                  <div className="bg-gray-50 p-3 rounded-lg space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-gray-600">纬度:</span>
                      <span className="font-mono">{formatCoordinate(currentLocation.latitude, "lat")}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">经度:</span>
                      <span className="font-mono">{formatCoordinate(currentLocation.longitude, "lng")}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">精度:</span>
                      <span>±{currentLocation.accuracy.toFixed(0)}m</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">时间:</span>
                      <span>{new Date(currentLocation.timestamp).toLocaleTimeString()}</span>
                    </div>
                  </div>

                  <Button onClick={recordCurrentLocation} className="w-full">
                    <Save className="h-4 w-4 mr-2" />
                    记录处理位置
                  </Button>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        {/* Nearby Devices */}
        <Card className="mb-4">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <Target className="h-5 w-5" />
              附近设备
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="flex items-center justify-between p-3 border rounded-lg">
                <div className="flex items-center gap-3">
                  <Zap className="h-4 w-4 text-blue-500" />
                  <div>
                    <div className="font-medium text-sm">DJI M300 RTK #001</div>
                    <div className="text-xs text-gray-500">192.168.1.100</div>
                  </div>
                </div>
                <div className="text-right">
                  <Badge variant="outline" className="text-xs">
                    150m
                  </Badge>
                  <div className="w-2 h-2 bg-green-500 rounded-full mt-1 ml-auto" />
                </div>
              </div>

              <div className="flex items-center justify-between p-3 border rounded-lg">
                <div className="flex items-center gap-3">
                  <Shield className="h-4 w-4 text-purple-500" />
                  <div>
                    <div className="font-medium text-sm">无人机机库 #001</div>
                    <div className="text-xs text-gray-500">充电站</div>
                  </div>
                </div>
                <div className="text-right">
                  <Badge variant="outline" className="text-xs">
                    200m
                  </Badge>
                  <div className="w-2 h-2 bg-green-500 rounded-full mt-1 ml-auto" />
                </div>
              </div>

              <div className="flex items-center justify-between p-3 border rounded-lg">
                <div className="flex items-center gap-3">
                  <Server className="h-4 w-4 text-orange-500" />
                  <div>
                    <div className="font-medium text-sm">视频服务器</div>
                    <div className="text-xs text-gray-500">192.168.1.200</div>
                  </div>
                </div>
                <div className="text-right">
                  <Badge variant="outline" className="text-xs">
                    500m
                  </Badge>
                  <div className="w-2 h-2 bg-yellow-500 rounded-full mt-1 ml-auto" />
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Quick Actions */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">快速操作</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 gap-3">
              <Button variant="outline" className="h-16 flex-col bg-transparent">
                <Camera className="h-5 w-5 mb-1" />
                <span className="text-sm">拍照记录</span>
              </Button>
              <Button variant="outline" className="h-16 flex-col bg-transparent">
                <MapPin className="h-5 w-5 mb-1" />
                <span className="text-sm">标记位置</span>
              </Button>
              <Button variant="outline" className="h-16 flex-col bg-transparent">
                <Activity className="h-5 w-5 mb-1" />
                <span className="text-sm">开始巡检</span>
              </Button>
              <Button variant="outline" className="h-16 flex-col bg-transparent">
                <RefreshCw className="h-5 w-5 mb-1" />
                <span className="text-sm">同步数据</span>
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
