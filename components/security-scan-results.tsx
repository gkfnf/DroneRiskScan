"use client"

import { useState } from "react"
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs"
import { Progress } from "@/components/ui/progress"
import { Badge } from "@/components/ui/badge"
import { AlertTriangle, CheckCircle, Server, Database, Radio, Shield } from "lucide-react"

const vulnerabilities = [
  {
    id: 1,
    device: "DJI M300 RTK #001",
    category: "网络安全",
    issue: "固件版本过低",
    severity: "高危",
    details: "当前固件版本存在已知的远程代码执行漏洞",
    icon: <Server className="h-4 w-4" />,
  },
  {
    id: 2,
    device: "DJI M300 RTK #001",
    category: "射频安全",
    issue: "控制信号未加密",
    severity: "高危",
    details: "控制信号使用明文传输，易被截获和篡改",
    icon: <Radio className="h-4 w-4" />,
  },
  {
    id: 3,
    device: "Ground Control Station",
    category: "网络安全",
    issue: "弱密码策略",
    severity: "中危",
    details: "使用默认密码，未启用多因素认证",
    icon: <Server className="h-4 w-4" />,
  },
  {
    id: 4,
    device: "Video Stream Server",
    category: "数据安全",
    issue: "视频流未加密",
    severity: "高危",
    details: "视频数据明文传输，可被未授权方截获",
    icon: <Database className="h-4 w-4" />,
  },
  {
    id: 5,
    device: "Network Router",
    category: "网络安全",
    issue: "开放不必要端口",
    severity: "中危",
    details: "多个不必要的服务端口开放，增加攻击面",
    icon: <Server className="h-4 w-4" />,
  },
]

export function SecurityScanResults() {
  const [activeTab, setActiveTab] = useState("vulnerabilities")

  const getSeverityIcon = (severity: string) => {
    switch (severity) {
      case "高危":
        return <AlertTriangle className="h-4 w-4 text-red-500" />
      case "中危":
        return <AlertTriangle className="h-4 w-4 text-amber-500" />
      case "低危":
        return <AlertTriangle className="h-4 w-4 text-blue-500" />
      default:
        return <CheckCircle className="h-4 w-4 text-green-500" />
    }
  }

  const getSeverityBadge = (severity: string) => {
    switch (severity) {
      case "高危":
        return <Badge variant="destructive">高危</Badge>
      case "中危":
        return (
          <Badge variant="outline" className="border-amber-500 text-amber-500">
            中危
          </Badge>
        )
      case "低危":
        return (
          <Badge variant="outline" className="border-blue-500 text-blue-500">
            低危
          </Badge>
        )
      default:
        return <Badge className="bg-green-500">安全</Badge>
    }
  }

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
        <div className="rounded-md border p-4">
          <div className="flex items-center justify-between">
            <div className="text-sm font-medium">网络安全评分</div>
            <div className="text-2xl font-bold text-amber-500">72</div>
          </div>
          <Progress value={72} className="mt-2 h-2" indicatorClassName="bg-amber-500" />
          <p className="mt-1 text-xs text-muted-foreground">发现3个安全问题</p>
        </div>
        <div className="rounded-md border p-4">
          <div className="flex items-center justify-between">
            <div className="text-sm font-medium">数据安全评分</div>
            <div className="text-2xl font-bold text-red-500">65</div>
          </div>
          <Progress value={65} className="mt-2 h-2" indicatorClassName="bg-red-500" />
          <p className="mt-1 text-xs text-muted-foreground">发现1个安全问题</p>
        </div>
        <div className="rounded-md border p-4">
          <div className="flex items-center justify-between">
            <div className="text-sm font-medium">射频安全评分</div>
            <div className="text-2xl font-bold text-red-500">58</div>
          </div>
          <Progress value={58} className="mt-2 h-2" indicatorClassName="bg-red-500" />
          <p className="mt-1 text-xs text-muted-foreground">发现2个安全问题</p>
        </div>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="vulnerabilities">漏洞列表</TabsTrigger>
          <TabsTrigger value="compliance">合规检查</TabsTrigger>
          <TabsTrigger value="recommendations">修复建议</TabsTrigger>
        </TabsList>
        <TabsContent value="vulnerabilities" className="mt-4">
          <div className="rounded-md border">
            <div className="border-b p-4">
              <h3 className="font-medium">检测到的安全漏洞</h3>
              <p className="text-sm text-muted-foreground">最近一次扫描发现的安全问题</p>
            </div>
            <div className="divide-y">
              {vulnerabilities.map((vuln) => (
                <div key={vuln.id} className="p-4">
                  <div className="flex items-start justify-between">
                    <div>
                      <div className="flex items-center gap-2">
                        {getSeverityIcon(vuln.severity)}
                        <span className="font-medium">{vuln.issue}</span>
                        {getSeverityBadge(vuln.severity)}
                      </div>
                      <div className="mt-1 text-sm text-muted-foreground">{vuln.device}</div>
                      <div className="mt-2 text-sm">{vuln.details}</div>
                    </div>
                    <div className="flex items-center gap-1 rounded-full bg-muted px-2 py-1 text-xs">
                      {vuln.icon}
                      <span>{vuln.category}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </TabsContent>
        <TabsContent value="compliance" className="mt-4">
          <div className="rounded-md border p-6">
            <div className="flex flex-col items-center justify-center py-6">
              <Shield className="h-12 w-12 text-muted-foreground" />
              <h3 className="mt-4 text-lg font-medium">合规检查</h3>
              <p className="mt-2 text-center text-sm text-muted-foreground max-w-md">
                对照行业标准和法规要求进行合规性评估
              </p>
            </div>
          </div>
        </TabsContent>
        <TabsContent value="recommendations" className="mt-4">
          <div className="rounded-md border p-6">
            <div className="flex flex-col items-center justify-center py-6">
              <CheckCircle className="h-12 w-12 text-muted-foreground" />
              <h3 className="mt-4 text-lg font-medium">修复建议</h3>
              <p className="mt-2 text-center text-sm text-muted-foreground max-w-md">
                针对发现的安全问题提供具体修复方案和最佳实践建议
              </p>
            </div>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
