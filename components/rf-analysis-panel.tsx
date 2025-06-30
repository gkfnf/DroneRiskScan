"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs"
import { Badge } from "@/components/ui/badge"
import { Radio, Play, Pause, Download, AlertTriangle, Wifi } from "lucide-react"

export function RfAnalysisPanel() {
  const [isCapturing, setIsCapturing] = useState(false)

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="font-medium">HackRF设备状态</h3>
          <p className="text-sm text-muted-foreground">射频安全测试设备</p>
        </div>
        <Badge className="bg-green-500">已连接</Badge>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div className="rounded-md border p-4">
          <h4 className="font-medium">设备信息</h4>
          <div className="mt-2 space-y-2 text-sm">
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">设备类型</span>
              <span>HackRF One</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">固件版本</span>
              <span>v2023.01.1</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">采样率</span>
              <span>20 MSPS</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">频率范围</span>
              <span>1 MHz - 6 GHz</span>
            </div>
          </div>
        </div>

        <div className="rounded-md border p-4">
          <h4 className="font-medium">扫描配置</h4>
          <div className="mt-2 space-y-2 text-sm">
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">目标频段</span>
              <span>2.4GHz / 5.8GHz</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">扫描模式</span>
              <span>被动监听</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground">分析工具</span>
              <span>GNU Radio</span>
            </div>
          </div>
        </div>
      </div>

      <div className="flex items-center justify-between">
        <Button variant={isCapturing ? "destructive" : "default"} onClick={() => setIsCapturing(!isCapturing)}>
          {isCapturing ? (
            <>
              <Pause className="mr-2 h-4 w-4" />
              停止捕获
            </>
          ) : (
            <>
              <Play className="mr-2 h-4 w-4" />
              开始捕获
            </>
          )}
        </Button>
        <Button variant="outline" disabled={!isCapturing}>
          <Download className="mr-2 h-4 w-4" />
          导出数据
        </Button>
      </div>

      <Tabs defaultValue="spectrum">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="spectrum">频谱分析</TabsTrigger>
          <TabsTrigger value="signals">信号检测</TabsTrigger>
          <TabsTrigger value="threats">威胁分析</TabsTrigger>
        </TabsList>
        <TabsContent value="spectrum" className="mt-4">
          <div className="rounded-md border">
            <div className="border-b p-4">
              <h3 className="font-medium">频谱分析</h3>
              <p className="text-sm text-muted-foreground">实时频谱监测</p>
            </div>
            <div className="h-[300px] bg-slate-900 relative">
              {/* 模拟频谱图 */}
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="text-white text-sm">
                  {isCapturing ? "正在捕获射频信号..." : "点击开始捕获按钮开始射频分析"}
                </div>
              </div>

              {/* 模拟频谱图上的标记 */}
              {isCapturing && (
                <>
                  <div className="absolute top-1/3 left-1/4 w-1 h-20 bg-red-500 opacity-70"></div>
                  <div className="absolute top-1/3 left-1/4 text-xs text-white">
                    <div className="bg-red-500 px-1 rounded">控制信号</div>
                  </div>

                  <div className="absolute top-1/4 left-2/3 w-1 h-16 bg-blue-500 opacity-70"></div>
                  <div className="absolute top-1/4 left-2/3 text-xs text-white">
                    <div className="bg-blue-500 px-1 rounded">遥测信号</div>
                  </div>

                  <div className="absolute top-1/2 left-1/2 w-1 h-24 bg-green-500 opacity-70"></div>
                  <div className="absolute top-1/2 left-1/2 text-xs text-white">
                    <div className="bg-green-500 px-1 rounded">视频传输</div>
                  </div>
                </>
              )}
            </div>
          </div>
        </TabsContent>
        <TabsContent value="signals" className="mt-4">
          <div className="rounded-md border">
            <div className="border-b p-4">
              <h3 className="font-medium">检测到的信号</h3>
              <p className="text-sm text-muted-foreground">无人机通信信号分析</p>
            </div>
            <div className="divide-y">
              <div className="p-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Radio className="h-4 w-4 text-red-500" />
                    <span className="font-medium">2.4GHz 控制信号</span>
                  </div>
                  <Badge variant="destructive">未加密</Badge>
                </div>
                <p className="mt-1 text-sm text-muted-foreground">DJI M300 RTK #001</p>
                <p className="mt-1 text-sm">控制信号使用明文传输，易被截获和篡改</p>
              </div>
              <div className="p-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Wifi className="h-4 w-4 text-amber-500" />
                    <span className="font-medium">2.4GHz 遥测数据</span>
                  </div>
                  <Badge variant="outline" className="border-amber-500 text-amber-500">
                    弱加密
                  </Badge>
                </div>
                <p className="mt-1 text-sm text-muted-foreground">DJI M300 RTK #001</p>
                <p className="mt-1 text-sm">使用弱加密算法，存在被破解风险</p>
              </div>
              <div className="p-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <Wifi className="h-4 w-4 text-green-500" />
                    <span className="font-medium">5.8GHz 视频传输</span>
                  </div>
                  <Badge className="bg-green-500">强加密</Badge>
                </div>
                <p className="mt-1 text-sm text-muted-foreground">DJI M300 RTK #001</p>
                <p className="mt-1 text-sm">使用AES-256加密，安全性较高</p>
              </div>
            </div>
          </div>
        </TabsContent>
        <TabsContent value="threats" className="mt-4">
          <div className="rounded-md border">
            <div className="border-b p-4">
              <h3 className="font-medium">射频威胁分析</h3>
              <p className="text-sm text-muted-foreground">潜在射频安全威胁</p>
            </div>
            <div className="divide-y">
              <div className="p-4">
                <div className="flex items-center gap-2">
                  <AlertTriangle className="h-4 w-4 text-red-500" />
                  <span className="font-medium">信号干扰风险</span>
                  <Badge variant="destructive">高危</Badge>
                </div>
                <p className="mt-1 text-sm text-muted-foreground">DJI M300 RTK #001</p>
                <p className="mt-1 text-sm">控制信号易受干扰，可被低成本设备干扰导致失控</p>
              </div>
              <div className="p-4">
                <div className="flex items-center gap-2">
                  <AlertTriangle className="h-4 w-4 text-red-500" />
                  <span className="font-medium">信号重放攻击</span>
                  <Badge variant="destructive">高危</Badge>
                </div>
                <p className="mt-1 text-sm text-muted-foreground">DJI M300 RTK #001</p>
                <p className="mt-1 text-sm">控制信号可被记录并重放，缺少时间戳或挑战-响应机制</p>
              </div>
              <div className="p-4">
                <div className="flex items-center gap-2">
                  <AlertTriangle className="h-4 w-4 text-amber-500" />
                  <span className="font-medium">GPS欺骗风险</span>
                  <Badge variant="outline" className="border-amber-500 text-amber-500">
                    中危
                  </Badge>
                </div>
                <p className="mt-1 text-sm text-muted-foreground">DJI M300 RTK #001</p>
                <p className="mt-1 text-sm">GPS信号可被伪造，可能导致无人机偏离预定航线</p>
              </div>
            </div>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
