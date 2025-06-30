"use client";

import { useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
    Shield,
    AlertTriangle,
    Database,
    Network,
    Monitor,
    Wifi,
    HardDrive,
    Cloud,
    Lock,
    Eye,
    Router,
    Smartphone,
    Plane,
    Building,
    TrendingUp,
    Activity,
    Server,
    Globe,
    Zap,
    MapPin,
    FileText,
} from "lucide-react";

// 基于架构图的实际风险分层评估
const ARCHITECTURE_LAYERS = {
    APPLICATION_LAYER: {
        name: "应用层",
        icon: <Monitor className="h-4 w-4" />,
        components: [
            "无人机业务应用系统",
            "统一技术中台",
            "业务中台",
            "电网GIS",
            "北斗综合服务平台",
            "人工智能平台",
        ],
        riskCategories: [
            "业务逻辑漏洞",
            "API安全",
            "数据处理安全",
            "平台集成风险",
        ],
    },
    NETWORK_LAYER: {
        name: "网络层",
        icon: <Network className="h-4 w-4" />,
        components: [
            "光纤/电力线专网",
            "无线专网(APN)",
            "互联网/WiFi/WAPI",
            "安全接入网关",
            "网闸",
            "防火墙",
        ],
        riskCategories: ["网络劫持", "中间人攻击", "DDoS攻击", "网络隔离失效"],
    },
    TERMINAL_LAYER: {
        name: "终端层",
        icon: <Smartphone className="h-4 w-4" />,
        components: [
            "移动终端(含APP、监控端)",
            "机场(含APP)",
            "无人机",
            "飞行控制系统",
        ],
        riskCategories: ["终端劫持", "恶意APP", "固件漏洞", "物理安全"],
    },
};

// 基于实际架构的安全区域划分
const SECURITY_ZONES = {
    MANAGEMENT_ZONE: {
        name: "管理信息大区",
        icon: <Building className="h-4 w-4" />,
        securityLevel: "III级",
        systems: [
            "综合管理系统",
            "运维管理系统",
            "监控分析平台",
            "GIS系统",
            "决策支持系统",
        ],
        threats: ["内部威胁", "权限滥用", "数据泄露", "系统入侵"],
    },
    PRODUCTION_ZONE: {
        name: "生产控制大区",
        icon: <Zap className="h-4 w-4" />,
        securityLevel: "II级",
        systems: [
            "无人机业务应用系统",
            "设备监控",
            "作业管控",
            "飞行调度",
            "数据采集处理",
        ],
        threats: ["生产中断", "控制系统攻击", "数据篡改", "设备故障"],
    },
    INTERNET_ZONE: {
        name: "互联网大区",
        icon: <Globe className="h-4 w-4" />,
        securityLevel: "IV级",
        systems: ["第三方服务对接", "云端数据同步", "远程运维", "移动应用服务"],
        threats: ["外部攻击", "恶意软件", "数据窃取", "服务劫持"],
    },
};

// 实际业务场景的风险评估模型
const BUSINESS_RISK_SCENARIOS = [
    {
        id: "BR001",
        name: "无人机业务应用系统API接口安全",
        layer: "APPLICATION_LAYER",
        zone: "PRODUCTION_ZONE",
        severity: "高",
        probability: "中",
        impact: "API接口未授权访问可能导致飞行任务被恶意篡改",
        businessImpact: "影响电力巡检作业安全，可能导致误判或漏检",
        mitigation: "实施API网关认证、OAuth2.0授权、接口加密传输",
        riskScore: 85,
        components: ["作业计划编制", "任务派发", "数据回传", "数据分析"],
        detectionMethods: ["API调用异常监控", "权限访问审计", "接口流量分析"],
    },
    {
        id: "BR002",
        name: "统一技术中台数据处理安全",
        layer: "APPLICATION_LAYER",
        zone: "MANAGEMENT_ZONE",
        severity: "高",
        probability: "中",
        impact: "中台数据处理漏洞可能导致多业务系统数据泄露",
        businessImpact: "影响电网GIS、人工智能平台等核心业务",
        mitigation: "数据脱敏处理、权限最小化、审计日志完整性",
        riskScore: 88,
        components: ["数据集成", "数据清洗", "数据分析", "结果分发"],
        detectionMethods: ["数据流向监控", "异常访问检测", "数据完整性校验"],
    },
    {
        id: "BR003",
        name: "光纤/电力线专网通信劫持",
        layer: "NETWORK_LAYER",
        zone: "PRODUCTION_ZONE",
        severity: "极高",
        probability: "低",
        impact: "专网通信被劫持可直接控制无人机飞行",
        businessImpact: "可能导致无人机坠毁、电力设施损坏、人员伤亡",
        mitigation: "专网加密、通信协议安全加固、异常检测",
        riskScore: 95,
        components: ["飞行控制指令", "遥测数据传输", "视频流传输", "状态监控"],
        detectionMethods: ["通信异常监控", "指令完整性校验", "加密强度检测"],
    },
    {
        id: "BR004",
        name: "安全接入网关配置漏洞",
        layer: "NETWORK_LAYER",
        zone: "MANAGEMENT_ZONE",
        severity: "高",
        probability: "中",
        impact: "网关配置错误可能导致网络边界失效",
        businessImpact: "管理信息大区与生产控制大区隔离失效",
        mitigation: "定期安全配置检查、访问控制策略优化、监控告警",
        riskScore: 78,
        components: ["访问控制策略", "流量过滤规则", "认证机制", "日志记录"],
        detectionMethods: ["配置合规检查", "流量异常分析", "访问行为监控"],
    },
    {
        id: "BR005",
        name: "移动终端APP安全漏洞",
        layer: "TERMINAL_LAYER",
        zone: "INTERNET_ZONE",
        severity: "中",
        probability: "高",
        impact: "移动APP漏洞可能导致作业数据泄露",
        businessImpact: "现场作业人员隐私泄露、作业计划被窃取",
        mitigation: "APP安全加固、代码混淆、运行时保护",
        riskScore: 72,
        components: ["用户认证", "数据存储", "通信加密", "权限控制"],
        detectionMethods: ["APP行为监控", "异常操作检测", "设备指纹识别"],
    },
    {
        id: "BR006",
        name: "无人机固件后门植入",
        layer: "TERMINAL_LAYER",
        zone: "PRODUCTION_ZONE",
        severity: "极高",
        probability: "低",
        impact: "固件后门可完全控制无人机，绕过所有安全措施",
        businessImpact: "严重威胁电力系统安全，可能被用于恶意攻击",
        mitigation: "固件签名验证、供应链安全管控、定期固件检查",
        riskScore: 92,
        components: ["飞行控制器", "导航系统", "通信模块", "载荷设备"],
        detectionMethods: ["固件完整性检测", "行为基线分析", "异常指令监控"],
    },
];

// 数据处理系统安全评估（非数据本身，而是处理数据的系统）
const DATA_PROCESSING_SYSTEMS = {
    COLLECTION_SYSTEM: {
        name: "数据采集系统",
        icon: <Database className="h-4 w-4" />,
        components: ["传感器接口", "数据预处理", "格式转换", "质量检查"],
        risks: ["数据注入攻击", "采集中断", "数据丢失", "格式篡改"],
        securityMeasures: ["输入验证", "备份机制", "完整性校验", "访问控制"],
    },
    TRANSMISSION_SYSTEM: {
        name: "数据传输系统",
        icon: <Network className="h-4 w-4" />,
        components: ["传输协议", "加密模块", "压缩算法", "错误纠正"],
        risks: ["中间人攻击", "数据窃听", "传输中断", "协议漏洞"],
        securityMeasures: ["端到端加密", "证书验证", "传输监控", "协议加固"],
    },
    STORAGE_SYSTEM: {
        name: "数据存储系统",
        icon: <HardDrive className="h-4 w-4" />,
        components: ["数据库", "文件系统", "备份存储", "归档系统"],
        risks: ["数据库注入", "权限滥用", "存储泄露", "备份失效"],
        securityMeasures: ["访问控制", "数据加密", "审计日志", "备份验证"],
    },
    ANALYSIS_SYSTEM: {
        name: "数据分析系统",
        icon: <TrendingUp className="h-4 w-4" />,
        components: ["算法引擎", "模型训练", "结果生成", "报告输出"],
        risks: ["算法投毒", "模型窃取", "结果篡改", "计算资源滥用"],
        securityMeasures: ["算法验证", "模型保护", "结果校验", "资源监控"],
    },
};

// 实时威胁监控数据
const THREAT_MONITORING = {
    currentThreatLevel: 72,
    activeThreats: [
        {
            id: "T001",
            type: "网络扫描",
            source: "外部IP",
            target: "无人机业务应用系统",
            severity: "中",
            timestamp: "2025-01-15 14:25:33",
            status: "监控中",
        },
        {
            id: "T002",
            type: "异常API调用",
            source: "移动终端",
            target: "统一技术中台",
            severity: "高",
            timestamp: "2025-01-15 14:18:22",
            status: "已阻断",
        },
        {
            id: "T003",
            type: "固件完整性异常",
            source: "无人机DJI-M300-003",
            target: "飞行控制系统",
            severity: "高",
            timestamp: "2025-01-15 13:55:14",
            status: "调查中",
        },
    ],
    riskTrends: [
        { time: "09:00", level: 45 },
        { time: "10:00", level: 52 },
        { time: "11:00", level: 58 },
        { time: "12:00", level: 65 },
        { time: "13:00", level: 72 },
        { time: "14:00", level: 68 },
    ],
};

export function PowerIndustryRisks() {
    const [selectedScenario, setSelectedScenario] = useState(
        BUSINESS_RISK_SCENARIOS[0],
    );
    const [selectedLayer, setSelectedLayer] = useState("APPLICATION_LAYER");
    const [selectedZone, setSelectedZone] = useState("PRODUCTION_ZONE");

    const getSeverityColor = (severity: string) => {
        switch (severity) {
            case "极高":
                return "bg-red-600";
            case "高":
                return "bg-orange-500";
            case "中":
                return "bg-yellow-500";
            case "低":
                return "bg-green-500";
            default:
                return "bg-gray-500";
        }
    };

    const getSeverityVariant = (severity: string) => {
        switch (severity) {
            case "极高":
                return "destructive";
            case "高":
                return "destructive";
            case "中":
                return "secondary";
            case "低":
                return "outline";
            default:
                return "outline";
        }
    };

    const getZoneSecurityLevel = (zone: string) => {
        return (
            SECURITY_ZONES[zone as keyof typeof SECURITY_ZONES]
                ?.securityLevel || "未知"
        );
    };

    const filteredScenarios = BUSINESS_RISK_SCENARIOS.filter((scenario) => {
        return (
            (selectedLayer === "ALL" || scenario.layer === selectedLayer) &&
            (selectedZone === "ALL" || scenario.zone === selectedZone)
        );
    });

    return (
        <div className="space-y-6">
            {/* 威胁态势总览 */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <Card>
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-muted-foreground">
                                    当前威胁等级
                                </p>
                                <p className="text-2xl font-bold text-orange-600">
                                    {THREAT_MONITORING.currentThreatLevel}
                                </p>
                            </div>
                            <div className="rounded-full bg-orange-100 p-2">
                                <AlertTriangle className="h-5 w-5 text-orange-600" />
                            </div>
                        </div>
                        <Progress
                            value={THREAT_MONITORING.currentThreatLevel}
                            className="mt-2"
                        />
                    </CardContent>
                </Card>

                <Card>
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-muted-foreground">
                                    活跃威胁
                                </p>
                                <p className="text-2xl font-bold text-red-600">
                                    {THREAT_MONITORING.activeThreats.length}
                                </p>
                            </div>
                            <div className="rounded-full bg-red-100 p-2">
                                <Eye className="h-5 w-5 text-red-600" />
                            </div>
                        </div>
                    </CardContent>
                </Card>

                <Card>
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-muted-foreground">
                                    风险场景
                                </p>
                                <p className="text-2xl font-bold text-purple-600">
                                    {BUSINESS_RISK_SCENARIOS.length}
                                </p>
                            </div>
                            <div className="rounded-full bg-purple-100 p-2">
                                <Shield className="h-5 w-5 text-purple-600" />
                            </div>
                        </div>
                    </CardContent>
                </Card>

                <Card>
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-muted-foreground">
                                    安全区域
                                </p>
                                <p className="text-2xl font-bold text-blue-600">
                                    {Object.keys(SECURITY_ZONES).length}
                                </p>
                            </div>
                            <div className="rounded-full bg-blue-100 p-2">
                                <Building className="h-5 w-5 text-blue-600" />
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </div>

            <Tabs defaultValue="architecture" className="w-full">
                <TabsList className="grid w-full grid-cols-5">
                    <TabsTrigger value="architecture">架构风险</TabsTrigger>
                    <TabsTrigger value="zones">安全区域</TabsTrigger>
                    <TabsTrigger value="data-systems">数据处理系统</TabsTrigger>
                    <TabsTrigger value="scenarios">风险场景</TabsTrigger>
                    <TabsTrigger value="monitoring">威胁监控</TabsTrigger>
                </TabsList>

                <TabsContent value="architecture" className="mt-6">
                    <div className="space-y-6">
                        {/* 架构分层风险评估 */}
                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <Network className="h-5 w-5" />
                                    无人机应用总体架构风险评估
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-4">
                                    {Object.entries(ARCHITECTURE_LAYERS).map(
                                        ([key, layer]) => (
                                            <div
                                                key={key}
                                                className="border rounded-lg p-4"
                                            >
                                                <div className="flex items-center gap-3 mb-3">
                                                    <div className="rounded-full bg-blue-100 p-2">
                                                        {layer.icon}
                                                    </div>
                                                    <div>
                                                        <h4 className="font-medium">
                                                            {layer.name}
                                                        </h4>
                                                        <p className="text-sm text-muted-foreground">
                                                            组件数量:{" "}
                                                            {
                                                                layer.components
                                                                    .length
                                                            }
                                                        </p>
                                                    </div>
                                                </div>

                                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                                    <div>
                                                        <h5 className="text-sm font-medium mb-2">
                                                            主要组件
                                                        </h5>
                                                        <div className="space-y-1">
                                                            {layer.components.map(
                                                                (
                                                                    component,
                                                                    index,
                                                                ) => (
                                                                    <Badge
                                                                        key={
                                                                            index
                                                                        }
                                                                        variant="outline"
                                                                        className="mr-2 mb-1"
                                                                    >
                                                                        {
                                                                            component
                                                                        }
                                                                    </Badge>
                                                                ),
                                                            )}
                                                        </div>
                                                    </div>
                                                    <div>
                                                        <h5 className="text-sm font-medium mb-2">
                                                            风险类别
                                                        </h5>
                                                        <div className="space-y-1">
                                                            {layer.riskCategories.map(
                                                                (
                                                                    risk,
                                                                    index,
                                                                ) => (
                                                                    <div
                                                                        key={
                                                                            index
                                                                        }
                                                                        className="text-sm text-red-600"
                                                                    >
                                                                        • {risk}
                                                                    </div>
                                                                ),
                                                            )}
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                        ),
                                    )}
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>

                <TabsContent value="zones" className="mt-6">
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                        {Object.entries(SECURITY_ZONES).map(([key, zone]) => (
                            <Card key={key}>
                                <CardHeader>
                                    <CardTitle className="flex items-center gap-2">
                                        {zone.icon}
                                        {zone.name}
                                    </CardTitle>
                                    <Badge variant="outline" className="w-fit">
                                        {zone.securityLevel}
                                    </Badge>
                                </CardHeader>
                                <CardContent>
                                    <div className="space-y-4">
                                        <div>
                                            <h5 className="text-sm font-medium mb-2">
                                                主要系统
                                            </h5>
                                            <div className="space-y-1">
                                                {zone.systems.map(
                                                    (system, index) => (
                                                        <div
                                                            key={index}
                                                            className="text-sm text-muted-foreground"
                                                        >
                                                            • {system}
                                                        </div>
                                                    ),
                                                )}
                                            </div>
                                        </div>

                                        <div>
                                            <h5 className="text-sm font-medium mb-2">
                                                主要威胁
                                            </h5>
                                            <div className="space-y-1">
                                                {zone.threats.map(
                                                    (threat, index) => (
                                                        <Badge
                                                            key={index}
                                                            variant="destructive"
                                                            className="mr-1 mb-1 text-xs"
                                                        >
                                                            {threat}
                                                        </Badge>
                                                    ),
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>
                        ))}
                    </div>
                </TabsContent>

                <TabsContent value="data-systems" className="mt-6">
                    <div className="space-y-6">
                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <Database className="h-5 w-5" />
                                    数据处理系统安全评估
                                </CardTitle>
                                <p className="text-sm text-muted-foreground">
                                    评估处理无人机采集数据的各个系统的安全风险
                                </p>
                            </CardHeader>
                            <CardContent>
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                    {Object.entries(
                                        DATA_PROCESSING_SYSTEMS,
                                    ).map(([key, system]) => (
                                        <div
                                            key={key}
                                            className="border rounded-lg p-4"
                                        >
                                            <div className="flex items-center gap-3 mb-3">
                                                <div className="rounded-full bg-green-100 p-2">
                                                    {system.icon}
                                                </div>
                                                <h4 className="font-medium">
                                                    {system.name}
                                                </h4>
                                            </div>

                                            <div className="space-y-3">
                                                <div>
                                                    <h5 className="text-sm font-medium mb-2">
                                                        系统组件
                                                    </h5>
                                                    <div className="flex flex-wrap gap-1">
                                                        {system.components.map(
                                                            (
                                                                component,
                                                                index,
                                                            ) => (
                                                                <Badge
                                                                    key={index}
                                                                    variant="outline"
                                                                    className="text-xs"
                                                                >
                                                                    {component}
                                                                </Badge>
                                                            ),
                                                        )}
                                                    </div>
                                                </div>

                                                <div>
                                                    <h5 className="text-sm font-medium mb-2">
                                                        安全风险
                                                    </h5>
                                                    <div className="space-y-1">
                                                        {system.risks.map(
                                                            (risk, index) => (
                                                                <div
                                                                    key={index}
                                                                    className="text-sm text-red-600"
                                                                >
                                                                    • {risk}
                                                                </div>
                                                            ),
                                                        )}
                                                    </div>
                                                </div>

                                                <div>
                                                    <h5 className="text-sm font-medium mb-2">
                                                        安全措施
                                                    </h5>
                                                    <div className="space-y-1">
                                                        {system.securityMeasures.map(
                                                            (
                                                                measure,
                                                                index,
                                                            ) => (
                                                                <div
                                                                    key={index}
                                                                    className="text-sm text-green-600"
                                                                >
                                                                    ✓ {measure}
                                                                </div>
                                                            ),
                                                        )}
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>

                <TabsContent value="scenarios" className="mt-6">
                    <div className="space-y-6">
                        {/* 筛选器 */}
                        <Card>
                            <CardContent className="p-4">
                                <div className="flex items-center gap-4">
                                    <div>
                                        <label className="text-sm font-medium">
                                            架构层次
                                        </label>
                                        <select
                                            className="ml-2 px-3 py-1 border rounded"
                                            value={selectedLayer}
                                            onChange={(e) =>
                                                setSelectedLayer(e.target.value)
                                            }
                                        >
                                            <option value="ALL">
                                                全部层次
                                            </option>
                                            {Object.entries(
                                                ARCHITECTURE_LAYERS,
                                            ).map(([key, layer]) => (
                                                <option key={key} value={key}>
                                                    {layer.name}
                                                </option>
                                            ))}
                                        </select>
                                    </div>
                                    <div>
                                        <label className="text-sm font-medium">
                                            安全区域
                                        </label>
                                        <select
                                            className="ml-2 px-3 py-1 border rounded"
                                            value={selectedZone}
                                            onChange={(e) =>
                                                setSelectedZone(e.target.value)
                                            }
                                        >
                                            <option value="ALL">
                                                全部区域
                                            </option>
                                            {Object.entries(SECURITY_ZONES).map(
                                                ([key, zone]) => (
                                                    <option
                                                        key={key}
                                                        value={key}
                                                    >
                                                        {zone.name}
                                                    </option>
                                                ),
                                            )}
                                        </select>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>

                        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                            {/* 风险场景列表 */}
                            <div className="space-y-3">
                                {filteredScenarios.map((scenario) => (
                                    <div
                                        key={scenario.id}
                                        className={`p-4 border rounded-lg cursor-pointer transition-colors ${
                                            selectedScenario.id === scenario.id
                                                ? "border-blue-500 bg-blue-50"
                                                : "hover:bg-gray-50"
                                        }`}
                                        onClick={() =>
                                            setSelectedScenario(scenario)
                                        }
                                    >
                                        <div className="flex items-start justify-between">
                                            <div className="flex-1">
                                                <div className="flex items-center gap-2 mb-2">
                                                    <span className="text-sm font-mono text-muted-foreground">
                                                        {scenario.id}
                                                    </span>
                                                    <Badge
                                                        variant={getSeverityVariant(
                                                            scenario.severity,
                                                        )}
                                                    >
                                                        {scenario.severity}
                                                    </Badge>
                                                    <Badge variant="outline">
                                                        {
                                                            ARCHITECTURE_LAYERS[
                                                                scenario.layer as keyof typeof ARCHITECTURE_LAYERS
                                                            ]?.name
                                                        }
                                                    </Badge>
                                                </div>
                                                <h4 className="font-medium mb-1">
                                                    {scenario.name}
                                                </h4>
                                                <p className="text-sm text-muted-foreground">
                                                    {scenario.impact}
                                                </p>
                                            </div>
                                            <div className="text-right">
                                                <div className="text-lg font-bold">
                                                    {scenario.riskScore}
                                                </div>
                                                <div className="text-xs text-muted-foreground">
                                                    风险分数
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>

                            {/* 场景详情 */}
                            <Card>
                                <CardHeader>
                                    <CardTitle className="text-lg">
                                        {selectedScenario.name}
                                    </CardTitle>
                                    <div className="flex items-center gap-2">
                                        <Badge
                                            variant={getSeverityVariant(
                                                selectedScenario.severity,
                                            )}
                                        >
                                            {selectedScenario.severity}
                                        </Badge>
                                        <Badge variant="outline">
                                            {getZoneSecurityLevel(
                                                selectedScenario.zone,
                                            )}
                                        </Badge>
                                    </div>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <div>
                                        <span className="text-sm text-muted-foreground">
                                            风险描述
                                        </span>
                                        <p className="mt-1 text-sm">
                                            {selectedScenario.impact}
                                        </p>
                                    </div>

                                    <div>
                                        <span className="text-sm text-muted-foreground">
                                            业务影响
                                        </span>
                                        <p className="mt-1 text-sm">
                                            {selectedScenario.businessImpact}
                                        </p>
                                    </div>

                                    <div>
                                        <span className="text-sm text-muted-foreground">
                                            缓解措施
                                        </span>
                                        <p className="mt-1 text-sm">
                                            {selectedScenario.mitigation}
                                        </p>
                                    </div>

                                    <div>
                                        <span className="text-sm text-muted-foreground">
                                            影响组件
                                        </span>
                                        <div className="mt-1 flex flex-wrap gap-1">
                                            {selectedScenario.components.map(
                                                (component, index) => (
                                                    <Badge
                                                        key={index}
                                                        variant="secondary"
                                                    >
                                                        {component}
                                                    </Badge>
                                                ),
                                            )}
                                        </div>
                                    </div>

                                    <div>
                                        <span className="text-sm text-muted-foreground">
                                            检测方法
                                        </span>
                                        <div className="mt-1 space-y-1">
                                            {selectedScenario.detectionMethods.map(
                                                (method, index) => (
                                                    <div
                                                        key={index}
                                                        className="text-sm text-green-600"
                                                    >
                                                        ✓ {method}
                                                    </div>
                                                ),
                                            )}
                                        </div>
                                    </div>

                                    <div>
                                        <span className="text-sm text-muted-foreground">
                                            风险评分
                                        </span>
                                        <div className="mt-1">
                                            <div className="flex items-center gap-2">
                                                <Progress
                                                    value={
                                                        selectedScenario.riskScore
                                                    }
                                                    className="flex-1"
                                                />
                                                <span className="font-bold">
                                                    {selectedScenario.riskScore}
                                                    /100
                                                </span>
                                            </div>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>
                        </div>
                    </div>
                </TabsContent>

                <TabsContent value="monitoring" className="mt-6">
                    <div className="space-y-6">
                        {/* 实时威胁监控 */}
                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <Activity className="h-5 w-5" />
                                    实时威胁监控
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-4">
                                    {THREAT_MONITORING.activeThreats.map(
                                        (threat) => (
                                            <div
                                                key={threat.id}
                                                className="flex items-center gap-3 p-3 border rounded-lg"
                                            >
                                                <div
                                                    className={`w-3 h-3 rounded-full ${
                                                        threat.severity === "高"
                                                            ? "bg-red-500"
                                                            : threat.severity ===
                                                                "中"
                                                              ? "bg-yellow-500"
                                                              : "bg-green-500"
                                                    }`}
                                                ></div>
                                                <div className="flex-1">
                                                    <div className="flex items-center gap-2 mb-1">
                                                        <span className="font-medium text-sm">
                                                            {threat.type}
                                                        </span>
                                                        <Badge
                                                            variant="outline"
                                                            className="text-xs"
                                                        >
                                                            {threat.id}
                                                        </Badge>
                                                    </div>
                                                    <p className="text-sm text-muted-foreground">
                                                        来源: {threat.source} →
                                                        目标: {threat.target}
                                                    </p>
                                                    <p className="text-xs text-muted-foreground">
                                                        {threat.timestamp}
                                                    </p>
                                                </div>
                                                <div className="text-right">
                                                    <Badge
                                                        variant={
                                                            threat.severity ===
                                                            "高"
                                                                ? "destructive"
                                                                : "secondary"
                                                        }
                                                    >
                                                        {threat.severity}
                                                    </Badge>
                                                    <div className="text-xs text-muted-foreground mt-1">
                                                        {threat.status}
                                                    </div>
                                                </div>
                                            </div>
                                        ),
                                    )}
                                </div>
                            </CardContent>
                        </Card>

                        {/* 威胁趋势图 */}
                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <TrendingUp className="h-5 w-5" />
                                    威胁等级趋势
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="h-64 flex items-end justify-between gap-2 px-4">
                                    {THREAT_MONITORING.riskTrends.map(
                                        (point, index) => (
                                            <div
                                                key={index}
                                                className="flex flex-col items-center"
                                            >
                                                <div
                                                    className="bg-blue-500 w-8 rounded-t"
                                                    style={{
                                                        height: `${(point.level / 100) * 200}px`,
                                                    }}
                                                ></div>
                                                <span className="text-xs mt-2">
                                                    {point.time}
                                                </span>
                                            </div>
                                        ),
                                    )}
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    );
}
