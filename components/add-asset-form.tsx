"use client"

import type React from "react"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Plus } from "lucide-react"

export function AddAssetForm() {
  const [deviceName, setDeviceName] = useState("")
  const [deviceType, setDeviceType] = useState("")
  const [ipAddress, setIpAddress] = useState("")

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    // 处理表单提交
    console.log({ deviceName, deviceType, ipAddress })
    // 重置表单
    setDeviceName("")
    setDeviceType("")
    setIpAddress("")
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="deviceName">设备名称</Label>
        <Input
          id="deviceName"
          placeholder="输入设备名称"
          value={deviceName}
          onChange={(e) => setDeviceName(e.target.value)}
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="deviceType">设备类型</Label>
        <Select value={deviceType} onValueChange={setDeviceType}>
          <SelectTrigger id="deviceType">
            <SelectValue placeholder="选择设备类型" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="drone">无人机</SelectItem>
            <SelectItem value="control">控制站</SelectItem>
            <SelectItem value="server">服务器</SelectItem>
            <SelectItem value="network">网络设备</SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div className="space-y-2">
        <Label htmlFor="ipAddress">IP地址</Label>
        <Input
          id="ipAddress"
          placeholder="192.168.1.100"
          value={ipAddress}
          onChange={(e) => setIpAddress(e.target.value)}
        />
      </div>
      <Button type="submit" className="w-full">
        <Plus className="mr-2 h-4 w-4" />
        添加资产
      </Button>
    </form>
  )
}
