import { Badge } from "@/components/ui/badge"
import { Progress } from "@/components/ui/progress"
import { Battery, Wifi, MapPin } from "lucide-react"

const drones = [
  {
    id: "DJI-001",
    status: "巡检中",
    battery: 78,
    location: "北区主干线",
    signal: "强",
  },
  {
    id: "DJI-002",
    status: "巡检中",
    battery: 65,
    location: "东区支线",
    signal: "中",
  },
  {
    id: "DJI-003",
    status: "充电中",
    battery: 32,
    location: "基地站",
    signal: "强",
  },
  {
    id: "DJI-004",
    status: "待命",
    battery: 100,
    location: "基地站",
    signal: "强",
  },
]

export function DroneStatus() {
  return (
    <div className="space-y-4">
      {drones.map((drone) => (
        <div key={drone.id} className="rounded-lg border p-4">
          <div className="flex items-center justify-between">
            <div className="font-medium">{drone.id}</div>
            <Badge
              variant={drone.status === "巡检中" ? "default" : drone.status === "充电中" ? "outline" : "secondary"}
            >
              {drone.status}
            </Badge>
          </div>
          <div className="mt-2 space-y-3">
            <div className="space-y-1">
              <div className="flex items-center justify-between text-sm">
                <div className="flex items-center gap-1">
                  <Battery className="h-3 w-3" />
                  <span>电池</span>
                </div>
                <span
                  className={`${
                    drone.battery < 20 ? "text-red-500" : drone.battery < 50 ? "text-amber-500" : "text-green-500"
                  }`}
                >
                  {drone.battery}%
                </span>
              </div>
              <Progress
                value={drone.battery}
                className={`h-1.5 ${
                  drone.battery < 20 ? "bg-red-100" : drone.battery < 50 ? "bg-amber-100" : "bg-green-100"
                }`}
                indicatorClassName={`${
                  drone.battery < 20 ? "bg-red-500" : drone.battery < 50 ? "bg-amber-500" : "bg-green-500"
                }`}
              />
            </div>
            <div className="flex items-center justify-between text-sm">
              <div className="flex items-center gap-1">
                <MapPin className="h-3 w-3" />
                <span>位置</span>
              </div>
              <span>{drone.location}</span>
            </div>
            <div className="flex items-center justify-between text-sm">
              <div className="flex items-center gap-1">
                <Wifi className="h-3 w-3" />
                <span>信号</span>
              </div>
              <span
                className={`${
                  drone.signal === "弱" ? "text-red-500" : drone.signal === "中" ? "text-amber-500" : "text-green-500"
                }`}
              >
                {drone.signal}
              </span>
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}
