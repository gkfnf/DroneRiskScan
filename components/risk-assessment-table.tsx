import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { AlertTriangle, ArrowRight, Server, Database, Radio } from "lucide-react"

const riskData = [
  {
    id: 1,
    risk: "无人机控制信号未加密",
    category: "射频安全",
    severity: "高危",
    impact: "可能导致无人机被劫持或控制",
    icon: <Radio className="h-4 w-4" />,
  },
  {
    id: 2,
    risk: "Web后台存在SQL注入漏洞",
    category: "网络安全",
    severity: "高危",
    impact: "可能导致数据库被非法访问或数据泄露",
    icon: <Server className="h-4 w-4" />,
  },
  {
    id: 3,
    risk: "任务下发系统缺少访问控制",
    category: "网络安全",
    severity: "高危",
    impact: "未授权用户可能下发恶意任务",
    icon: <Server className="h-4 w-4" />,
  },
  {
    id: 4,
    risk: "无人机遥测数据未加密传输",
    category: "数据安全",
    severity: "中危",
    impact: "敏感数据可能被截获",
    icon: <Database className="h-4 w-4" />,
  },
  {
    id: 5,
    risk: "无人机易受GPS欺骗攻击",
    category: "射频安全",
    severity: "中危",
    impact: "可能导致无人机偏离预定航线",
    icon: <Radio className="h-4 w-4" />,
  },
]

export function RiskAssessmentTable() {
  return (
    <div className="overflow-hidden rounded-md border">
      <table className="w-full">
        <thead>
          <tr className="bg-muted/50">
            <th className="whitespace-nowrap px-4 py-3 text-left text-sm font-medium">风险</th>
            <th className="whitespace-nowrap px-4 py-3 text-left text-sm font-medium">类别</th>
            <th className="whitespace-nowrap px-4 py-3 text-left text-sm font-medium">严重程度</th>
            <th className="whitespace-nowrap px-4 py-3 text-left text-sm font-medium hidden md:table-cell">潜在影响</th>
            <th className="whitespace-nowrap px-4 py-3 text-right text-sm font-medium"></th>
          </tr>
        </thead>
        <tbody>
          {riskData.map((risk, index) => (
            <tr key={risk.id} className={index % 2 === 0 ? "bg-background" : "bg-muted/30"}>
              <td className="px-4 py-3 text-sm">
                <div className="flex items-center gap-2">
                  <AlertTriangle
                    className={`h-4 w-4 ${risk.severity === "高危" ? "text-red-500" : "text-amber-500"}`}
                  />
                  {risk.risk}
                </div>
              </td>
              <td className="px-4 py-3 text-sm">
                <div className="flex items-center gap-1">
                  {risk.icon}
                  <span>{risk.category}</span>
                </div>
              </td>
              <td className="px-4 py-3 text-sm">
                <Badge variant={risk.severity === "高危" ? "destructive" : "default"}>{risk.severity}</Badge>
              </td>
              <td className="px-4 py-3 text-sm text-muted-foreground hidden md:table-cell">{risk.impact}</td>
              <td className="px-4 py-3 text-sm text-right">
                <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                  <ArrowRight className="h-4 w-4" />
                  <span className="sr-only">查看详情</span>
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
