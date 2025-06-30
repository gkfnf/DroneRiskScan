import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
    Shield,
    Database,
    Radio,
    AlertTriangle,
    CheckCircle,
    Search,
    Plus,
    FileText,
    Users,
    Activity,
} from "lucide-react";
import { AssetList } from "@/components/asset-list";
import { AddAssetForm } from "@/components/add-asset-form";
import { SecurityScanResults } from "@/components/security-scan-results";
import { RfAnalysisPanel } from "@/components/rf-analysis-panel";
import { PowerIndustryRisks } from "@/components/power-industry-risks";
import { DroneDataManagement } from "@/components/drone-data-management";
import { TaskManagementSystem } from "@/components/task-management-system";

export default function Dashboard() {
    return (
        <div className="flex min-h-screen flex-col bg-white">
            <header className="border-b py-4">
                <div className="container">
                    <div className="flex items-center gap-2">
                        <Shield className="h-6 w-6 text-blue-600" />
                        <div>
                            <h1 className="text-xl font-bold">
                                电力无人机安全扫描系统
                            </h1>
                            <p className="text-sm text-muted-foreground">
                                专业的无人机设备安全漏洞检测与评估平台
                            </p>
                        </div>
                    </div>
                </div>
            </header>
            <main className="flex-1 py-6">
                <div className="container">
                    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
                        <Card>
                            <CardContent className="flex items-center justify-between p-6">
                                <div>
                                    <p className="text-sm text-muted-foreground">
                                        总资产数
                                    </p>
                                    <p className="text-3xl font-bold text-blue-600">
                                        4
                                    </p>
                                </div>
                                <div className="rounded-full bg-blue-100 p-2">
                                    <Database className="h-5 w-5 text-blue-600" />
                                </div>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardContent className="flex items-center justify-between p-6">
                                <div>
                                    <p className="text-sm text-muted-foreground">
                                        在线设备
                                    </p>
                                    <p className="text-3xl font-bold text-green-600">
                                        2
                                    </p>
                                </div>
                                <div className="rounded-full bg-green-100 p-2">
                                    <CheckCircle className="h-5 w-5 text-green-600" />
                                </div>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardContent className="flex items-center justify-between p-6">
                                <div>
                                    <p className="text-sm text-muted-foreground">
                                        射频威胁
                                    </p>
                                    <p className="text-3xl font-bold text-amber-600">
                                        2
                                    </p>
                                </div>
                                <div className="rounded-full bg-amber-100 p-2">
                                    <Radio className="h-5 w-5 text-amber-600" />
                                </div>
                            </CardContent>
                        </Card>
                        <Card>
                            <CardContent className="flex items-center justify-between p-6">
                                <div>
                                    <p className="text-sm text-muted-foreground">
                                        高危漏洞
                                    </p>
                                    <p className="text-3xl font-bold text-red-600">
                                        6
                                    </p>
                                </div>
                                <div className="rounded-full bg-red-100 p-2">
                                    <AlertTriangle className="h-5 w-5 text-red-600" />
                                </div>
                            </CardContent>
                        </Card>
                    </div>

                    <div className="mt-6">
                        <Tabs defaultValue="assets" className="w-full">
                            <TabsList className="grid w-full grid-cols-8">
                                <TabsTrigger value="assets">
                                    资产管理
                                </TabsTrigger>
                                <TabsTrigger value="scan">扫描监控</TabsTrigger>
                                <TabsTrigger value="rf">射频安全</TabsTrigger>
                                <TabsTrigger value="risks">
                                    风险评估
                                </TabsTrigger>
                                <TabsTrigger value="data">数据管理</TabsTrigger>
                                <TabsTrigger value="tasks">
                                    任务管理
                                </TabsTrigger>
                                <TabsTrigger value="reports">
                                    安全报告
                                </TabsTrigger>
                                <TabsTrigger value="settings">
                                    系统设置
                                </TabsTrigger>
                            </TabsList>

                            <TabsContent value="assets" className="mt-6">
                                <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
                                    <div className="lg:col-span-2">
                                        <div className="rounded-md border">
                                            <div className="flex items-center gap-2 border-b p-4">
                                                <Database className="h-5 w-5" />
                                                <h2 className="font-medium">
                                                    资产清单
                                                </h2>
                                            </div>
                                            <div className="p-2">
                                                <p className="px-2 py-1 text-sm text-muted-foreground">
                                                    管理所有无人机设备和相关设备资产
                                                </p>
                                                <AssetList />
                                            </div>
                                        </div>
                                    </div>
                                    <div>
                                        <div className="rounded-md border">
                                            <div className="flex items-center gap-2 border-b p-4">
                                                <Plus className="h-5 w-5" />
                                                <h2 className="font-medium">
                                                    添加资产
                                                </h2>
                                            </div>
                                            <div className="p-4">
                                                <AddAssetForm />
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="scan" className="mt-6">
                                <div className="rounded-md border">
                                    <div className="flex items-center justify-between border-b p-4">
                                        <div className="flex items-center gap-2">
                                            <Search className="h-5 w-5" />
                                            <h2 className="font-medium">
                                                安全扫描结果
                                            </h2>
                                        </div>
                                        <Button>
                                            <Search className="mr-2 h-4 w-4" />
                                            开始扫描
                                        </Button>
                                    </div>
                                    <div className="p-4">
                                        <SecurityScanResults />
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="rf" className="mt-6">
                                <div className="rounded-md border">
                                    <div className="flex items-center justify-between border-b p-4">
                                        <div className="flex items-center gap-2">
                                            <Radio className="h-5 w-5" />
                                            <h2 className="font-medium">
                                                射频安全分析
                                            </h2>
                                        </div>
                                        <Button>
                                            <Radio className="mr-2 h-4 w-4" />
                                            开始射频分析
                                        </Button>
                                    </div>
                                    <div className="p-4">
                                        <RfAnalysisPanel />
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="risks" className="mt-6">
                                <div className="rounded-md border">
                                    <div className="flex items-center gap-2 border-b p-4">
                                        <AlertTriangle className="h-5 w-5" />
                                        <h2 className="font-medium">
                                            电力行业风险评估
                                        </h2>
                                    </div>
                                    <div className="p-4">
                                        <PowerIndustryRisks />
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="data" className="mt-6">
                                <div className="rounded-md border">
                                    <div className="flex items-center gap-2 border-b p-4">
                                        <Database className="h-5 w-5" />
                                        <h2 className="font-medium">
                                            无人机数据管理
                                        </h2>
                                    </div>
                                    <div className="p-4">
                                        <DroneDataManagement />
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="tasks" className="mt-6">
                                <div className="rounded-md border">
                                    <div className="flex items-center gap-2 border-b p-4">
                                        <FileText className="h-5 w-5" />
                                        <h2 className="font-medium">
                                            任务管理系统
                                        </h2>
                                    </div>
                                    <div className="p-4">
                                        <TaskManagementSystem />
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="reports" className="mt-6">
                                <div className="rounded-md border p-6">
                                    <div className="flex flex-col items-center justify-center py-10">
                                        <div className="rounded-full bg-blue-100 p-4">
                                            <Shield className="h-8 w-8 text-blue-600" />
                                        </div>
                                        <h3 className="mt-4 text-lg font-medium">
                                            安全报告
                                        </h3>
                                        <p className="mt-2 text-center text-sm text-muted-foreground max-w-md">
                                            生成详细的安全评估报告，包括网络安全、数据安全和射频安全风险分析
                                        </p>
                                        <Button className="mt-4">
                                            生成安全报告
                                        </Button>
                                    </div>
                                </div>
                            </TabsContent>

                            <TabsContent value="settings" className="mt-6">
                                <div className="rounded-md border p-6">
                                    <div className="flex flex-col items-center justify-center py-10">
                                        <h3 className="text-lg font-medium">
                                            系统设置
                                        </h3>
                                        <p className="mt-2 text-center text-sm text-muted-foreground max-w-md">
                                            配置系统参数、扫描策略和安全评估标准
                                        </p>
                                        <Button
                                            variant="outline"
                                            className="mt-4 bg-transparent"
                                        >
                                            配置系统
                                        </Button>
                                    </div>
                                </div>
                            </TabsContent>
                        </Tabs>
                    </div>
                </div>
            </main>
            <footer className="border-t py-4">
                <div className="container">
                    <div className="flex items-center justify-between">
                        <p className="text-sm text-muted-foreground">
                            © 2025 电力无人机安全扫描系统
                        </p>
                        <p className="text-sm text-muted-foreground">
                            版本 1.0.0
                        </p>
                    </div>
                </div>
            </footer>
        </div>
    );
}
