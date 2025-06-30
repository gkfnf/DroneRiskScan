import { Progress } from "@/components/ui/progress"

const complianceData = [
  {
    standard: "等保2.0三级",
    compliance: 85,
    status: "合规",
  },
  {
    standard: "电力行业安全规范",
    compliance: 92,
    status: "合规",
  },
  {
    standard: "工控系统安全要求",
    compliance: 68,
    status: "部分合规",
  },
  {
    standard: "无人机安全标准",
    compliance: 75,
    status: "部分合规",
  },
]

export function ComplianceStatus() {
  return (
    <div className="space-y-4">
      {complianceData.map((item) => (
        <div key={item.standard} className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="font-medium text-sm">{item.standard}</span>
            <span className={`text-xs font-medium ${item.compliance >= 80 ? "text-green-600" : "text-amber-600"}`}>
              {item.status} ({item.compliance}%)
            </span>
          </div>
          <Progress
            value={item.compliance}
            className="h-2"
            indicatorClassName={item.compliance >= 80 ? "bg-green-600" : "bg-amber-600"}
          />
        </div>
      ))}
    </div>
  )
}
