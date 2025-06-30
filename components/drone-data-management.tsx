"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Input } from "@/components/ui/input";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import {
    Database,
    Download,
    Upload,
    Filter,
    Search,
    FileText,
    Image,
    Video,
    BarChart3,
    MapPin,
    Calendar,
    Clock,
    HardDrive,
    Wifi,
    Battery,
    Camera,
    Thermometer,
    Wind,
    Eye,
    AlertTriangle,
    CheckCircle,
    RefreshCw,
} from "lucide-react";

// Drone data types and structures
const DATA_TYPES = {
    TELEMETRY: {
        name: "遥测数据",
        icon: <BarChart3 className="h-4 w-4" />,
        color: "blue",
        fields: [
            "GPS坐标",
            "高度",
            "速度",
            "电池电量",
            "信号强度",
            "温度",
            "湿度",
            "风速",
        ],
    },
    IMAGERY: {
        name: "图像数据",
        icon: <Camera className="h-4 w-4" />,
        color: "green",
        fields: [
            "照片",
            "红外图像",
            "多光谱图像",
            "位置信息",
            "时间戳",
            "分辨率",
        ],
    },
    VIDEO: {
        name: "视频数据",
        icon: <Video className="h-4 w-4" />,
        color: "purple",
        fields: [
            "实时视频流",
            "录制视频",
            "编码格式",
            "码率",
            "帧率",
            "分辨率",
        ],
    },
    INSPECTION: {
        name: "巡检数据",
        icon: <Eye className="h-4 w-4" />,
        color: "orange",
        fields: [
            "设备状态",
            "缺陷检测",
            "热成像",
            "振动数据",
            "声学数据",
            "电磁场",
        ],
    },
    ENVIRONMENTAL: {
        name: "环境数据",
        icon: <Wind className="h-4 w-4" />,
        color: "teal",
        fields: ["气温", "湿度", "风速", "风向", "气压", "能见度", "降水量"],
    },
    SYSTEM: {
        name: "系统数据",
        icon: <HardDrive className="h-4 w-4" />,
        color: "gray",
        fields: [
            "系统日志",
            "错误报告",
            "性能指标",
            "网络状态",
            "存储使用",
            "CPU负载",
        ],
    },
};

// Mock drone data collection status
const DRONE_DATA_STATUS = [
    {
        droneId: "DJI-M300-001",
        droneName: "东区巡检无人机1号",
        status: "active",
        location: {
            lat: 39.9042,
            lng: 116.4074,
            address: "北京市朝阳区110kV变电站",
        },
        mission: "高压线路巡检",
        dataCollection: {
            telemetry: { status: "collecting", rate: "1Hz", size: "2.3MB" },
            imagery: { status: "collecting", rate: "2fps", size: "156MB" },
            video: { status: "recording", rate: "30fps", size: "1.2GB" },
            inspection: {
                status: "analyzing",
                rate: "continuous",
                size: "45MB",
            },
        },
        lastUpdate: "2025-01-15 14:23:45",
        batteryLevel: 87,
        signalStrength: 95,
        storageUsed: 68,
    },
    {
        droneId: "DJI-M300-002",
        droneName: "西区巡检无人机2号",
        status: "standby",
        location: {
            lat: 39.8942,
            lng: 116.3974,
            address: "北京市西城区220kV变电站",
        },
        mission: "待命中",
        dataCollection: {
            telemetry: { status: "idle", rate: "0Hz", size: "0MB" },
            imagery: { status: "idle", rate: "0fps", size: "0MB" },
            video: { status: "idle", rate: "0fps", size: "0MB" },
            inspection: { status: "idle", rate: "idle", size: "0MB" },
        },
        lastUpdate: "2025-01-15 13:45:12",
        batteryLevel: 95,
        signalStrength: 88,
        storageUsed: 12,
    },
    {
        droneId: "DJI-M300-003",
        droneName: "南区巡检无人机3号",
        status: "maintenance",
        location: { lat: 39.8742, lng: 116.4174, address: "维护基站" },
        mission: "设备维护",
        dataCollection: {
            telemetry: { status: "offline", rate: "0Hz", size: "0MB" },
            imagery: { status: "offline", rate: "0fps", size: "0MB" },
            video: { status: "offline", rate: "0fps", size: "0MB" },
            inspection: { status: "offline", rate: "offline", size: "0MB" },
        },
        lastUpdate: "2025-01-15 09:30:00",
        batteryLevel: 0,
        signalStrength: 0,
        storageUsed: 89,
    },
];

// Data analytics insights
const DATA_ANALYTICS = {
    totalDataCollected: "2.8TB",
    dailyGrowth: "156GB",
    dataTypes: {
        telemetry: { size: "145GB", percentage: 12 },
        imagery: { size: "892GB", percentage: 32 },
        video: { size: "1.2TB", percentage: 43 },
        inspection: { size: "234GB", percentage: 8 },
        environmental: { size: "89GB", percentage: 3 },
        system: { size: "56GB", percentage: 2 },
    },
    qualityMetrics: {
        dataIntegrity: 98.5,
        completeness: 94.2,
        accuracy: 96.8,
        timeliness: 99.1,
    },
    storageHealth: {
        total: "10TB",
        used: "2.8TB",
        available: "7.2TB",
        usage: 28,
    },
};

// Recent data processing activities
const RECENT_ACTIVITIES = [
    {
        id: 1,
        type: "data_processing",
        description: "东区110kV线路巡检数据分析完成",
        timestamp: "2025-01-15 14:20:00",
        status: "completed",
        details: "处理了156张高分辨率图像，检测到3个潜在缺陷",
    },
    {
        id: 2,
        type: "data_backup",
        description: "每日数据备份已启动",
        timestamp: "2025-01-15 14:00:00",
        status: "in_progress",
        details: "正在备份892GB图像数据到云存储",
    },
    {
        id: 3,
        type: "anomaly_detection",
        description: "检测到异常遥测数据",
        timestamp: "2025-01-15 13:45:00",
        status: "alert",
        details: "DJI-M300-001 GPS信号出现短暂中断",
    },
    {
        id: 4,
        type: "data_export",
        description: "巡检报告数据导出",
        timestamp: "2025-01-15 13:30:00",
        status: "completed",
        details: "导出了过去7天的巡检数据和分析结果",
    },
];

export function DroneDataManagement() {
    const [selectedDrone, setSelectedDrone] = useState(DRONE_DATA_STATUS[0]);
    const [dataFilter, setDataFilter] = useState("all");
    const [searchTerm, setSearchTerm] = useState("");
    const [isRefreshing, setIsRefreshing] = useState(false);

    const getStatusColor = (status: string) => {
        switch (status) {
            case "active":
                return "bg-green-500";
            case "standby":
                return "bg-yellow-500";
            case "maintenance":
                return "bg-red-500";
            case "collecting":
                return "bg-blue-500";
            case "recording":
                return "bg-purple-500";
            case "analyzing":
                return "bg-orange-500";
            case "idle":
                return "bg-gray-400";
            case "offline":
                return "bg-red-400";
            default:
                return "bg-gray-500";
        }
    };

    const getStatusIcon = (status: string) => {
        switch (status) {
            case "completed":
                return <CheckCircle className="h-4 w-4 text-green-500" />;
            case "in_progress":
                return (
                    <RefreshCw className="h-4 w-4 text-blue-500 animate-spin" />
                );
            case "alert":
                return <AlertTriangle className="h-4 w-4 text-red-500" />;
            default:
                return <Clock className="h-4 w-4 text-gray-500" />;
        }
    };

    const handleRefresh = () => {
        setIsRefreshing(true);
        setTimeout(() => setIsRefreshing(false), 2000);
    };

    return (
        <div className="space-y-6">
            {/* Data Overview Dashboard */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <Card>
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-muted-foreground">
                                    总数据量
                                </p>
                                <p className="text-2xl font-bold text-blue-600">
                                    {DATA_ANALYTICS.totalDataCollected}
                                </p>
                            </div>
                            <div className="rounded-full bg-blue-100 p-2">
                                <Database className="h-5 w-5 text-blue-600" />
                            </div>
                        </div>
                        <p className="text-xs text-muted-foreground mt-2">
                            今日增长: {DATA_ANALYTICS.dailyGrowth}
                        </p>
                    </CardContent>
                </Card>

                <Card>
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-muted-foreground">
                                    数据完整性
                                </p>
                                <p className="text-2xl font-bold text-green-600">
                                    {
                                        DATA_ANALYTICS.qualityMetrics
                                            .dataIntegrity
                                    }
                                    %
                                </p>
                            </div>
                            <div className="rounded-full bg-green-100 p-2">
                                <CheckCircle className="h-5 w-5 text-green-600" />
                            </div>
                        </div>
                        <Progress
                            value={DATA_ANALYTICS.qualityMetrics.dataIntegrity}
                            className="mt-2"
                        />
                    </CardContent>
                </Card>

                <Card>
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-muted-foreground">
                                    存储使用率
                                </p>
                                <p className="text-2xl font-bold text-purple-600">
                                    {DATA_ANALYTICS.storageHealth.usage}%
                                </p>
                            </div>
                            <div className="rounded-full bg-purple-100 p-2">
                                <HardDrive className="h-5 w-5 text-purple-600" />
                            </div>
                        </div>
                        <p className="text-xs text-muted-foreground mt-2">
                            可用: {DATA_ANALYTICS.storageHealth.available}
                        </p>
                    </CardContent>
                </Card>

                <Card>
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-muted-foreground">
                                    活跃无人机
                                </p>
                                <p className="text-2xl font-bold text-orange-600">
                                    {
                                        DRONE_DATA_STATUS.filter(
                                            (d) => d.status === "active",
                                        ).length
                                    }
                                </p>
                            </div>
                            <div className="rounded-full bg-orange-100 p-2">
                                <Eye className="h-5 w-5 text-orange-600" />
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </div>

            <Tabs defaultValue="overview" className="w-full">
                <TabsList className="grid w-full grid-cols-5">
                    <TabsTrigger value="overview">数据概览</TabsTrigger>
                    <TabsTrigger value="collection">数据采集</TabsTrigger>
                    <TabsTrigger value="analysis">数据分析</TabsTrigger>
                    <TabsTrigger value="storage">存储管理</TabsTrigger>
                    <TabsTrigger value="export">数据导出</TabsTrigger>
                </TabsList>

                <TabsContent value="overview" className="mt-6">
                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                        {/* Data Type Distribution */}
                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <BarChart3 className="h-5 w-5" />
                                    数据类型分布
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-4">
                                    {Object.entries(
                                        DATA_ANALYTICS.dataTypes,
                                    ).map(([key, data]) => {
                                        const typeInfo =
                                            DATA_TYPES[
                                                key.toUpperCase() as keyof typeof DATA_TYPES
                                            ];
                                        return (
                                            <div
                                                key={key}
                                                className="flex items-center gap-3"
                                            >
                                                <div className="rounded-full bg-blue-100 p-1">
                                                    {typeInfo?.icon}
                                                </div>
                                                <div className="flex-1">
                                                    <div className="flex items-center justify-between">
                                                        <span className="text-sm font-medium">
                                                            {typeInfo?.name}
                                                        </span>
                                                        <span className="text-sm text-muted-foreground">
                                                            {data.size}
                                                        </span>
                                                    </div>
                                                    <Progress
                                                        value={data.percentage}
                                                        className="mt-1"
                                                    />
                                                </div>
                                                <span className="text-sm text-muted-foreground">
                                                    {data.percentage}%
                                                </span>
                                            </div>
                                        );
                                    })}
                                </div>
                            </CardContent>
                        </Card>

                        {/* Data Quality Metrics */}
                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <CheckCircle className="h-5 w-5" />
                                    数据质量指标
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-4">
                                    {Object.entries(
                                        DATA_ANALYTICS.qualityMetrics,
                                    ).map(([key, value]) => {
                                        const labels = {
                                            dataIntegrity: "数据完整性",
                                            completeness: "数据完整度",
                                            accuracy: "数据准确性",
                                            timeliness: "数据及时性",
                                        };
                                        return (
                                            <div
                                                key={key}
                                                className="space-y-2"
                                            >
                                                <div className="flex items-center justify-between">
                                                    <span className="text-sm font-medium">
                                                        {
                                                            labels[
                                                                key as keyof typeof labels
                                                            ]
                                                        }
                                                    </span>
                                                    <span className="text-sm font-bold">
                                                        {value}%
                                                    </span>
                                                </div>
                                                <Progress value={value} />
                                            </div>
                                        );
                                    })}
                                </div>
                            </CardContent>
                        </Card>
                    </div>

                    {/* Recent Activities */}
                    <Card className="mt-6">
                        <CardHeader>
                            <div className="flex items-center justify-between">
                                <CardTitle className="flex items-center gap-2">
                                    <Clock className="h-5 w-5" />
                                    最近活动
                                </CardTitle>
                                <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={handleRefresh}
                                    disabled={isRefreshing}
                                >
                                    <RefreshCw
                                        className={`mr-2 h-4 w-4 ${
                                            isRefreshing ? "animate-spin" : ""
                                        }`}
                                    />
                                    刷新
                                </Button>
                            </div>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-3">
                                {RECENT_ACTIVITIES.map((activity) => (
                                    <div
                                        key={activity.id}
                                        className="flex items-start gap-3 p-3 border rounded-lg"
                                    >
                                        {getStatusIcon(activity.status)}
                                        <div className="flex-1">
                                            <p className="text-sm font-medium">
                                                {activity.description}
                                            </p>
                                            <p className="text-xs text-muted-foreground mt-1">
                                                {activity.details}
                                            </p>
                                            <p className="text-xs text-muted-foreground">
                                                {activity.timestamp}
                                            </p>
                                        </div>
                                        <Badge variant="outline">
                                            {activity.type.replace("_", " ")}
                                        </Badge>
                                    </div>
                                ))}
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="collection" className="mt-6">
                    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                        {/* Drone List */}
                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <Eye className="h-5 w-5" />
                                    无人机状态
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-3">
                                    {DRONE_DATA_STATUS.map((drone) => (
                                        <div
                                            key={drone.droneId}
                                            className={`p-3 border rounded-lg cursor-pointer transition-colors ${
                                                selectedDrone.droneId ===
                                                drone.droneId
                                                    ? "border-blue-500 bg-blue-50"
                                                    : "hover:bg-gray-50"
                                            }`}
                                            onClick={() =>
                                                setSelectedDrone(drone)
                                            }
                                        >
                                            <div className="flex items-center justify-between mb-2">
                                                <span className="font-medium text-sm">
                                                    {drone.droneName}
                                                </span>
                                                <div
                                                    className={`w-2 h-2 rounded-full ${getStatusColor(
                                                        drone.status,
                                                    )}`}
                                                ></div>
                                            </div>
                                            <p className="text-xs text-muted-foreground">
                                                {drone.droneId}
                                            </p>
                                            <p className="text-xs text-muted-foreground">
                                                {drone.mission}
                                            </p>
                                            <div className="flex items-center gap-2 mt-2">
                                                <div className="flex items-center gap-1">
                                                    <Battery className="h-3 w-3" />
                                                    <span className="text-xs">
                                                        {drone.batteryLevel}%
                                                    </span>
                                                </div>
                                                <div className="flex items-center gap-1">
                                                    <Wifi className="h-3 w-3" />
                                                    <span className="text-xs">
                                                        {drone.signalStrength}%
                                                    </span>
                                                </div>
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            </CardContent>
                        </Card>

                        {/* Selected Drone Details */}
                        <div className="lg:col-span-2 space-y-6">
                            <Card>
                                <CardHeader>
                                    <CardTitle>
                                        {selectedDrone.droneName}
                                    </CardTitle>
                                    <p className="text-sm text-muted-foreground">
                                        {selectedDrone.droneId}
                                    </p>
                                </CardHeader>
                                <CardContent>
                                    <div className="grid grid-cols-2 gap-4 mb-4">
                                        <div>
                                            <span className="text-sm text-muted-foreground">
                                                当前状态
                                            </span>
                                            <div className="flex items-center gap-2 mt-1">
                                                <div
                                                    className={`w-2 h-2 rounded-full ${getStatusColor(
                                                        selectedDrone.status,
                                                    )}`}
                                                ></div>
                                                <span className="text-sm font-medium">
                                                    {selectedDrone.status}
                                                </span>
                                            </div>
                                        </div>
                                        <div>
                                            <span className="text-sm text-muted-foreground">
                                                执行任务
                                            </span>
                                            <p className="text-sm font-medium mt-1">
                                                {selectedDrone.mission}
                                            </p>
                                        </div>
                                        <div>
                                            <span className="text-sm text-muted-foreground">
                                                位置
                                            </span>
                                            <p className="text-sm font-medium mt-1">
                                                {selectedDrone.location.address}
                                            </p>
                                        </div>
                                        <div>
                                            <span className="text-sm text-muted-foreground">
                                                最后更新
                                            </span>
                                            <p className="text-sm font-medium mt-1">
                                                {selectedDrone.lastUpdate}
                                            </p>
                                        </div>
                                    </div>

                                    <div className="grid grid-cols-3 gap-4 mb-4">
                                        <div className="text-center p-3 border rounded-lg">
                                            <Battery className="h-5 w-5 mx-auto mb-1" />
                                            <div className="text-lg font-bold">
                                                {selectedDrone.batteryLevel}%
                                            </div>
                                            <div className="text-xs text-muted-foreground">
                                                电池电量
                                            </div>
                                        </div>
                                        <div className="text-center p-3 border rounded-lg">
                                            <Wifi className="h-5 w-5 mx-auto mb-1" />
                                            <div className="text-lg font-bold">
                                                {selectedDrone.signalStrength}%
                                            </div>
                                            <div className="text-xs text-muted-foreground">
                                                信号强度
                                            </div>
                                        </div>
                                        <div className="text-center p-3 border rounded-lg">
                                            <HardDrive className="h-5 w-5 mx-auto mb-1" />
                                            <div className="text-lg font-bold">
                                                {selectedDrone.storageUsed}%
                                            </div>
                                            <div className="text-xs text-muted-foreground">
                                                存储使用
                                            </div>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>

                            {/* Data Collection Status */}
                            <Card>
                                <CardHeader>
                                    <CardTitle className="flex items-center gap-2">
                                        <Database className="h-5 w-5" />
                                        数据采集状态
                                    </CardTitle>
                                </CardHeader>
                                <CardContent>
                                    <div className="space-y-4">
                                        {Object.entries(
                                            selectedDrone.dataCollection,
                                        ).map(([type, data]) => {
                                            const typeInfo =
                                                DATA_TYPES[
                                                    type.toUpperCase() as keyof typeof DATA_TYPES
                                                ];
                                            return (
                                                <div
                                                    key={type}
                                                    className="flex items-center gap-3 p-3 border rounded-lg"
                                                >
                                                    <div className="rounded-full bg-blue-100 p-2">
                                                        {typeInfo?.icon}
                                                    </div>
                                                    <div className="flex-1">
                                                        <div className="flex items-center justify-between">
                                                            <span className="text-sm font-medium">
                                                                {typeInfo?.name}
                                                            </span>
                                                            <Badge
                                                                variant="outline"
                                                                className={`${getStatusColor(
                                                                    data.status,
                                                                )} text-white border-none`}
                                                            >
                                                                {data.status}
                                                            </Badge>
                                                        </div>
                                                        <div className="flex items-center gap-4 mt-1 text-xs text-muted-foreground">
                                                            <span>
                                                                采集率:{" "}
                                                                {data.rate}
                                                            </span>
                                                            <span>
                                                                数据量:{" "}
                                                                {data.size}
                                                            </span>
                                                        </div>
                                                    </div>
                                                </div>
                                            );
                                        })}
                                    </div>
                                </CardContent>
                            </Card>
                        </div>
                    </div>
                </TabsContent>

                <TabsContent value="analysis" className="mt-6">
                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <BarChart3 className="h-5 w-5" />
                                    数据分析工具
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-4">
                                    <div className="p-4 border rounded-lg">
                                        <h4 className="font-medium mb-2">
                                            智能缺陷检测
                                        </h4>
                                        <p className="text-sm text-muted-foreground mb-3">
                                            使用AI算法自动识别电力设备缺陷和异常
                                        </p>
                                        <Button className="w-full">
                                            启动缺陷检测
                                        </Button>
                                    </div>

                                    <div className="p-4 border rounded-lg">
                                        <h4 className="font-medium mb-2">
                                            热成像分析
                                        </h4>
                                        <p className="text-sm text-muted-foreground mb-3">
                                            分析红外热成像数据，检测设备过热问题
                                        </p>
                                        <Button
                                            className="w-full"
                                            variant="outline"
                                        >
                                            开始热成像分析
                                        </Button>
                                    </div>

                                    <div className="p-4 border rounded-lg">
                                        <h4 className="font-medium mb-2">
                                            趋势分析
                                        </h4>
                                        <p className="text-sm text-muted-foreground mb-3">
                                            分析历史数据，预测设备状态趋势
                                        </p>
                                        <Button
                                            className="w-full"
                                            variant="outline"
                                        >
                                            生成趋势报告
                                        </Button>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>

                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <AlertTriangle className="h-5 w-5" />
                                    异常检测结果
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-3">
                                    <div className="p-3 border rounded-lg border-red-200 bg-red-50">
                                        <div className="flex items-center gap-2 mb-2">
                                            <AlertTriangle className="h-4 w-4 text-red-500" />
                                            <span className="font-medium text-red-700">
                                                高温异常
                                            </span>
                                        </div>
                                        <p className="text-sm text-red-600">
                                            110kV变压器A相温度异常升高，建议立即检查
                                        </p>
                                        <p className="text-xs text-red-500 mt-1">
                                            检测时间: 2025-01-15 14:15
                                        </p>
                                    </div>

                                    <div className="p-3 border rounded-lg border-yellow-200 bg-yellow-50">
                                        <div className="flex items-center gap-2 mb-2">
                                            <AlertTriangle className="h-4 w-4 text-yellow-500" />
                                            <span className="font-medium text-yellow-700">
                                                电晕放电
                                            </span>
                                        </div>
                                        <p className="text-sm text-yellow-600">
                                            220kV线路绝缘子疑似电晕放电现象
                                        </p>
                                        <p className="text-xs text-yellow-500 mt-1">
                                            检测时间: 2025-01-15 13:42
                                        </p>
                                    </div>

                                    <div className="p-3 border rounded-lg border-blue-200 bg-blue-50">
                                        <div className="flex items-center gap-2 mb-2">
                                            <CheckCircle className="h-4 w-4 text-blue-500" />
                                            <span className="font-medium text-blue-700">
                                                设备正常
                                            </span>
                                        </div>
                                        <p className="text-sm text-blue-600">
                                            500kV主变压器运行状态正常
                                        </p>
                                        <p className="text-xs text-blue-500 mt-1">
                                            检测时间: 2025-01-15 13:30
                                        </p>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>

                <TabsContent value="storage" className="mt-6">
                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <HardDrive className="h-5 w-5" />
                                    存储统计
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="space-y-4">
                                    <div>
                                        <div className="flex items-center justify-between mb-2">
                                            <span className="text-sm font-medium">
                                                总存储空间
                                            </span>
                                            <span className="text-sm">
                                                {
                                                    DATA_ANALYTICS.storageHealth
                                                        .total
                                                }
                                            </span>
                                        </div>
                                        <Progress
                                            value={
                                                DATA_ANALYTICS.storageHealth
                                                    .usage
                                            }
                                        />
                                        <div className="flex items-center justify-between mt-1 text-xs text-muted-foreground">
                                            <span>
                                                已使用:{" "}
                                                {
                                                    DATA_ANALYTICS.storageHealth
                                                        .used
                                                }
                                            </span>
                                            <span>
                                                可用:{" "}
                                                {
                                                    DATA_ANALYTICS.storageHealth
                                                        .available
                                                }
                                            </span>
                                        </div>
                                    </div>

                                    <div className="grid grid-cols-2 gap-4">
                                        <div className="text-center p-3 border rounded-lg">
                                            <div className="text-lg font-bold text-blue-600">
                                                2.8TB
                                            </div>
                                            <div className="text-xs text-muted-foreground">
                                                本地存储
                                            </div>
                                        </div>
                                        <div className="text-center p-3 border rounded-lg">
                                            <div className="text-lg font-bold text-green-600">
                                                5.2TB
                                            </div>
                                            <div className="text-xs text-muted-foreground">
                                                云端备份
                                            </div>
                                        </div>
                                    </div>

                                    <div className="space-y-2">
                                        <Button className="w-full">
                                            <Upload className="mr-2 h-4 w-4" />
                                            开始云端备份
                                        </Button>
                                        <Button
                                            className="w-full"
                                            variant="outline"
                                        >
                                            <Download className="mr-2 h-4 w-4" />
                                            下载历史数据
                                        </Button>
                                    </div>
                                </div>
                            </CardContent>
                        </Card>

                        <Card>
                            <CardHeader>
                                <CardTitle className="flex items-center gap-2">
                                    <FileText className="h-5 w-5" />
                                    存储管理
                                </CardTitle>
                            </CardHeader>
                            <CardContent>
                                <div className="text-center py-10">
                                    <HardDrive className="h-12 w-12 text-muted-foreground mb-4 mx-auto" />
                                    <h3 className="font-medium">
                                        存储管理功能
                                    </h3>
                                    <p className="text-sm text-muted-foreground mt-1 max-w-md mx-auto">
                                        管理数据存储策略、备份计划和存储空间优化
                                    </p>
                                    <Button className="mt-4">
                                        配置存储策略
                                    </Button>
                                </div>
                            </CardContent>
                        </Card>
                    </div>
                </TabsContent>

                <TabsContent value="export" className="mt-6">
                    <Card>
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2">
                                <Download className="h-5 w-5" />
                                数据导出
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="text-center py-10">
                                <Download className="h-12 w-12 text-muted-foreground mb-4 mx-auto" />
                                <h3 className="font-medium">数据导出功能</h3>
                                <p className="text-sm text-muted-foreground mt-1 max-w-md mx-auto">
                                    导出巡检数据、分析报告和历史记录
                                </p>
                                <Button className="mt-4">开始数据导出</Button>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    );
}
