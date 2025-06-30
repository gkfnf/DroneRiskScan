import { Badge } from "@/components/ui/badge"
import { AlertTriangle, Zap, ArrowRight } from "lucide-react"
import { Button } from "@/components/ui/button"

const alerts = [
  {
    id: 1,
    location: "北区主干线 A-42",
    issue: "绝缘子损坏",
    severity: "高",
    timestamp: "今天 09:45",
    coordinates: "39.9042° N, 116.4074° E",
  },
  {
    id: 2,
    location: "东区支线 B-15",
    issue: "导线松弛",
    severity: "中",
    timestamp: "今天 11:23",
    coordinates: "39.9142° N, 116.4174° E",
  },
  {
    id: 3,
    location: "南区变电站附近",
    issue: "杆塔倾斜",
    severity: "高",
    timestamp: "昨天 16:30",
    coordinates: "39.8942° N, 116.3974° E",
  },
  {
    id: 4,
    location: "西区支线 C-08",
    issue: "植被过近",
    severity: "低",
    timestamp: "昨天 14:15",
    coordinates: "39.9242° N, 116.3874° E",
  },
]

export function RecentAlerts() {
  return (
    <div className="space-y-4">
      {alerts.map((alert) => (
        <div key={alert.id} className="flex items-start justify-between rounded-lg border p-4">
          <div className="space-y-1">
            <div className="flex items-center gap-2">
              <AlertTriangle
                className={`h-4 w-4 ${
                  alert.severity === "高"
                    ? "text-red-500"
                    : alert.severity === "中"
                      ? "text-amber-500"
                      : "text-blue-500"
                }`}
              />
              <span className="font-medium">{alert.issue}</span>
              <Badge
                variant={alert.severity === "高" ? "destructive" : alert.severity === "中" ? "default" : "secondary"}
              >
                {alert.severity}优先级
              </Badge>
            </div>
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Zap className="h-3 w-3" />
              <span>{alert.location}</span>
            </div>
            <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground">
              <span>{alert.timestamp}</span>
              <span>{alert.coordinates}</span>
            </div>
          </div>
          <Button variant="ghost" size="icon" className="h-8 w-8">
            <ArrowRight className="h-4 w-4" />
            <span className="sr-only">查看详情</span>
          </Button>
        </div>
      ))}
      <div className="flex justify-center">
        <Button variant="outline" size="sm">
          查看所有告警
        </Button>
      </div>
    </div>
  )
}
