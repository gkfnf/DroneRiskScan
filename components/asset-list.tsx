"use client"

import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { DrillIcon as Drone, Radio, Server, Wifi, Search } from "lucide-react"

const assets = [
  {
    id: 1,
    name: "DJI M300 RTK #001",
    type: "drone",
    ip: "192.168.1.100",
    status: "online",
    securityLevel: "危险",
  },
  {
    id: 2,
    name: "Ground Control Station",
    type: "control",
    ip: "192.168.1.50",
    status: "online",
    securityLevel: "安全",
  },
  {
    id: 3,
    name: "Video Stream Server",
    type: "server",
    ip: "192.168.1.200",
    status: "scanning",
    securityLevel: "危险",
  },
  {
    id: 4,
    name: "Network Router",
    type: "network",
    ip: "192.168.1.1",
    status: "offline",
    securityLevel: "高危险",
  },
]

export function AssetList() {
  const getIcon = (type: string) => {
    switch (type) {
      case "drone":
        return <Drone className="h-5 w-5" />
      case "control":
        return <Radio className="h-5 w-5" />
      case "server":
        return <Server className="h-5 w-5" />
      case "network":
        return <Wifi className="h-5 w-5" />
      default:
        return <Drone className="h-5 w-5" />
    }
  }

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "online":
        return <Badge className="bg-green-500">在线</Badge>
      case "offline":
        return <Badge variant="destructive">离线</Badge>
      case "scanning":
        return (
          <Badge variant="outline" className="border-amber-500 text-amber-500">
            扫描中
          </Badge>
        )
      default:
        return <Badge variant="outline">未知</Badge>
    }
  }

  const getSecurityBadge = (level: string) => {
    switch (level) {
      case "安全":
        return <Badge className="bg-green-500">安全</Badge>
      case "危险":
        return (
          <Badge variant="outline" className="border-amber-500 text-amber-500">
            危险
          </Badge>
        )
      case "高危险":
        return <Badge variant="destructive">高危险</Badge>
      default:
        return <Badge variant="outline">未知</Badge>
    }
  }

  return (
    <div className="space-y-2">
      {assets.map((asset) => (
        <div key={asset.id} className="flex items-center justify-between rounded-md border p-4">
          <div className="flex items-center gap-3">
            {getIcon(asset.type)}
            <div>
              <div className="font-medium">{asset.name}</div>
              <div className="text-sm text-muted-foreground">{asset.ip}</div>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <div className="flex gap-2">
              {getStatusBadge(asset.status)}
              {getSecurityBadge(asset.securityLevel)}
            </div>
            <Button variant="outline" size="icon" className="h-8 w-8 bg-transparent">
              <Search className="h-4 w-4" />
              <span className="sr-only">扫描</span>
            </Button>
          </div>
        </div>
      ))}
    </div>
  )
}
