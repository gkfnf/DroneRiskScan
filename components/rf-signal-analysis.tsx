"use client"

import { useState } from "react"
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs"
import { Button } from "@/components/ui/button"
import { Radio, Play, Pause, Download, AlertTriangle, Lock } from "lucide-react"

export function RfSignalAnalysis() {
  const [isCapturing, setIsCapturing] = useState(false)

  return (
    <div className="space-y-4">
      <Tabs defaultValue="spectrum">
        <div className="flex items-center justify-between">
          <TabsList>
            <TabsTrigger value="spectrum">频谱分析</TabsTrigger>
            <TabsTrigger value="demod">信号解调</TabsTrigger>
            <TabsTrigger value="protocol">协议分析</TabsTrigger>
          </TabsList>

          <div className="flex items-center gap-2">
            <Button
              variant={isCapturing ? "destructive" : "default"}
              size="sm"
              onClick={() => setIsCapturing(!isCapturing)}
            >
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
            <Button variant="outline" size="sm">
              <Download className="mr-2 h-4 w-4" />
              导出数据
            </Button>
          </div>
        </div>

        <TabsContent value="spectrum" className="mt-4">
          <div className="rounded-lg border">
            <div className="p-4 border-b">
              <h3 className="font-medium">无人机控制频段分析</h3>
              <p className="text-sm text-muted-foreground mt-1">使用HackRF捕获的2.4GHz和5.8GHz频段信号</p>
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
            <div className="p-4 border-t">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <h4 className="text-sm font-medium">检测到的信号</h4>
                  <ul className="mt-2 space-y-1 text-sm">
                    <li className="flex items-center gap-1">
                      <span className="w-2 h-2 rounded-full bg-red-500"></span>
                      <span>2.4GHz 控制信号 (未加密)</span>
                    </li>
                    <li className="flex items-center gap-1">
                      <span className="w-2 h-2 rounded-full bg-blue-500"></span>
                      <span>2.4GHz 遥测数据 (部分加密)</span>
                    </li>
                    <li className="flex items-center gap-1">
                      <span className="w-2 h-2 rounded-full bg-green-500"></span>
                      <span>5.8GHz 视频传输 (加密)</span>
                    </li>
                  </ul>
                </div>
                <div>
                  <h4 className="text-sm font-medium">安全评估</h4>
                  <ul className="mt-2 space-y-1 text-sm">
                    <li className="flex items-center gap-1 text-red-500">
                      <AlertTriangle className="h-3 w-3" />
                      <span>控制信号可被干扰</span>
                    </li>
                    <li className="flex items-center gap-1 text-amber-500">
                      <AlertTriangle className="h-3 w-3" />
                      <span>遥测数据加密强度不足</span>
                    </li>
                    <li className="flex items-center gap-1 text-green-500">
                      <span className="i-lucide-check-circle h-3 w-3"></span>
                      <span>视频传输加密安全</span>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </TabsContent>

        <TabsContent value="demod" className="mt-4">
          <div className="rounded-lg border p-4 h-[400px] flex items-center justify-center">
            <div className="text-center">
              <Radio className="h-12 w-12 text-muted-foreground mb-4 mx-auto" />
              <h3 className="font-medium">信号解调分析</h3>
              <p className="text-sm text-muted-foreground mt-1 max-w-md">
                此模块可对捕获的无人机信号进行解调分析，提取控制命令和数据包结构
              </p>
              <Button className="mt-4">
                <Play className="mr-2 h-4 w-4" />
                开始信号解调
              </Button>
            </div>
          </div>
        </TabsContent>

        <TabsContent value="protocol" className="mt-4">
          <div className="rounded-lg border p-4 h-[400px] flex items-center justify-center">
            <div className="text-center">
              <Lock className="h-12 w-12 text-muted-foreground mb-4 mx-auto" />
              <h3 className="font-medium">协议安全分析</h3>
              <p className="text-sm text-muted-foreground mt-1 max-w-md">
                此模块可分析无人机通信协议的安全性，检测加密强度、认证机制和协议漏洞
              </p>
              <Button className="mt-4">
                <Play className="mr-2 h-4 w-4" />
                开始协议分析
              </Button>
            </div>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
