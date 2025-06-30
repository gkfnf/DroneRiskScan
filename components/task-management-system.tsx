"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import {
    Calendar,
    Clock,
    User,
    MapPin,
    CheckCircle,
    AlertTriangle,
    Play,
    Pause,
    Square,
    RotateCcw,
    Settings,
    Plus,
    Eye,
    Edit,
    Trash2,
    Download,
    Upload,
    Filter,
    Search,
    Bell,
    Camera,
    Radio,
    Database,
    Shield,
    Zap,
    Wind,
    Thermometer,
    Activity,
    FileText,
    Users,
    Target,
} from "lucide-react";

// Task types and priorities
const TASK_TYPES = {
    INSPECTION: {
        name: "设备巡检",
        icon: <Eye className="h-4 w-4" />,
        color: "blue",
        description: "电力设备例行巡检任务",
    },
    MAINTENANCE: {
        name: "设备维护",
        icon: <Settings className="h-4 w-4" />,
        color: "green",
        description: "设备维护和保养任务",
    },
    EMERGENCY: {
        name: "应急响应",
        icon: <AlertTriangle className="h-4 w-4" />,
        color: "red",
        description: "紧急故障响应任务",
    },
    SURVEILLANCE: {
        name: "安全监控",
        icon: <Shield className="h-4 w-4" />,
        color: "purple",
        description: "安全监控和巡查任务",
    },
};

const TASK_PRIORITIES = {
    CRITICAL: { name: "紧急", color: "red", value: 1 },
    HIGH: { name: "高", color: "orange", value: 2 },
    MEDIUM: { name: "中", color: "yellow", value: 3 },
    LOW: { name: "低", color: "green", value: 4 },
};

const TASK_STATUS = {
    PENDING: {
        name: "等待中",
        color: "gray",
        icon: <Clock className="h-4 w-4" />,
    },
    ASSIGNED: {
        name: "已分配",
        color: "blue",
        icon: <User className="h-4 w-4" />,
    },
    IN_PROGRESS: {
        name: "执行中",
        color: "orange",
        icon: <Play className="h-4 w-4" />,
    },
    PAUSED: {
        name: "已暂停",
        color: "yellow",
        icon: <Pause className="h-4 w-4" />,
    },
    COMPLETED: {
        name: "已完成",
        color: "green",
        icon: <CheckCircle className="h-4 w-4" />,
    },
    FAILED: {
        name: "失败",
        color: "red",
        icon: <AlertTriangle className="h-4 w-4" />,
    },
    CANCELLED: {
        name: "已取消",
        color: "gray",
        icon: <Square className="h-4 w-4" />,
    },
};

// Mock task data
const MOCK_TASKS = [
    {
        id: "TASK-2025-001",
        title: "东区110kV输电线路巡检",
        type: "INSPECTION",
        priority: "HIGH",
        status: "IN_PROGRESS",
        assignedTo: {
            drone: "DJI-M300-001",
            operator: "张工程师",
            dataHandler: "李分析师",
        },
        location: {
            name: "东区110kV变电站",
            coordinates: { lat: 39.9042, lng: 116.4074 },
            area: "朝阳区供电局",
        },
        schedule: {
            startTime: "2025-01-15 08:00:00",
            endTime: "2025-01-15 12:00:00",
            estimatedDuration: "4小时",
        },
        progress: 65,
        description:
            "对东区110kV输电线路进行例行巡检，重点检查绝缘子、导线、杆塔状态",
        createdAt: "2025-01-15 07:30:00",
        updatedAt: "2025-01-15 10:15:00",
    },
    {
        id: "TASK-2025-002",
        title: "西区220kV变电站设备检测",
        type: "MAINTENANCE",
        priority: "MEDIUM",
        status: "ASSIGNED",
        assignedTo: {
            drone: "DJI-M300-002",
            operator: "王工程师",
            dataHandler: "赵分析师",
        },
        location: {
            name: "西区220kV变电站",
            coordinates: { lat: 39.8942, lng: 116.3974 },
            area: "西城区供电局",
        },
        schedule: {
            startTime: "2025-01-15 14:00:00",
            endTime: "2025-01-15 17:00:00",
            estimatedDuration: "3小时",
        },
        progress: 0,
        description: "对220kV变电站主要设备进行预防性检测维护",
        createdAt: "2025-01-15 09:00:00",
        updatedAt: "2025-01-15 09:30:00",
    },
    {
        id: "TASK-2025-003",
        title: "南区高压线路应急检查",
        type: "EMERGENCY",
        priority: "CRITICAL",
        status: "PENDING",
        assignedTo: {
            drone: "待分配",
            operator: "待分配",
            dataHandler: "待分配",
        },
        location: {
            name: "南区500kV输电走廊",
            coordinates: { lat: 39.8542, lng: 116.4274 },
            area: "丰台区供电局",
        },
        schedule: {
            startTime: "2025-01-15 15:30:00",
            endTime: "2025-01-15 18:30:00",
            estimatedDuration: "3小时",
        },
        progress: 0,
        description:
            "紧急响应：南区500kV线路跳闸，需要立即进行无人机巡检确定故障点",
        createdAt: "2025-01-15 15:00:00",
        updatedAt: "2025-01-15 15:00:00",
    },
];

export function TaskManagementSystem() {
    const [tasks, setTasks] = useState(MOCK_TASKS);
    const [selectedTask, setSelectedTask] = useState(MOCK_TASKS[0]);
    const [taskFilter, setTaskFilter] = useState("all");
    const [statusFilter, setStatusFilter] = useState("all");
    const [searchTerm, setSearchTerm] = useState("");
    const [showCreateDialog, setShowCreateDialog] = useState(false);

    const getStatusBadgeVariant = (status: string) => {
        switch (status) {
            case "COMPLETED":
                return "default";
            case "IN_PROGRESS":
                return "secondary";
            case "FAILED":
                return "destructive";
            case "CRITICAL":
                return "destructive";
            default:
                return "outline";
        }
    };

    const getPriorityColor = (priority: string) => {
        return (
            TASK_PRIORITIES[priority as keyof typeof TASK_PRIORITIES]?.color ||
            "gray"
        );
    };

    const filteredTasks = tasks.filter((task) => {
        const matchesSearch =
            task.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
            task.description.toLowerCase().includes(searchTerm.toLowerCase());
        const matchesType = taskFilter === "all" || task.type === taskFilter;
        const matchesStatus =
            statusFilter === "all" || task.status === statusFilter;
        return matchesSearch && matchesType && matchesStatus;
    });

    const handleTaskAction = (taskId: string, action: string) => {
        setTasks(
            tasks.map((task) => {
                if (task.id === taskId) {
                    switch (action) {
                        case "start":
                            return {
                                ...task,
                                status: "IN_PROGRESS" as const,
                                updatedAt: new Date()
                                    .toISOString()
                                    .slice(0, 19)
                                    .replace("T", " "),
                            };
                        case "pause":
                            return {
                                ...task,
                                status: "PAUSED" as const,
                                updatedAt: new Date()
                                    .toISOString()
                                    .slice(0, 19)
                                    .replace("T", " "),
                            };
                        case "complete":
                            return {
                                ...task,
                                status: "COMPLETED" as const,
                                progress: 100,
                                updatedAt: new Date()
                                    .toISOString()
                                    .slice(0, 19)
                                    .replace("T", " "),
                            };
                        case "cancel":
                            return {
                                ...task,
                                status: "CANCELLED" as const,
                                updatedAt: new Date()
                                    .toISOString()
                                    .slice(0, 19)
                                    .replace("T", " "),
                            };
                        default:
                            return task;
                    }
                }
                return task;
            }),
        );
    };

    return (
        <div className="space-y-6">
            {/* Task Overview Dashboard */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <Card>
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-muted-foreground">
                                    总任务数
                                </p>
                                <p className="text-2xl font-bold text-blue-600">
                                    {tasks.length}
                                </p>
                            </div>
                            <div className="rounded-full bg-blue-100 p-2">
                                <FileText className="h-5 w-5 text-blue-600" />
                            </div>
                        </div>
                    </CardContent>
                </Card>

                <Card>
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-muted-foreground">
                                    进行中
                                </p>
                                <p className="text-2xl font-bold text-orange-600">
                                    {
                                        tasks.filter(
                                            (t) => t.status === "IN_PROGRESS",
                                        ).length
                                    }
                                </p>
                            </div>
                            <div className="rounded-full bg-orange-100 p-2">
                                <Play className="h-5 w-5 text-orange-600" />
                            </div>
                        </div>
                    </CardContent>
                </Card>

                <Card>
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-muted-foreground">
                                    已完成
                                </p>
                                <p className="text-2xl font-bold text-green-600">
                                    {
                                        tasks.filter(
                                            (t) => t.status === "COMPLETED",
                                        ).length
                                    }
                                </p>
                            </div>
                            <div className="rounded-full bg-green-100 p-2">
                                <CheckCircle className="h-5 w-5 text-green-600" />
                            </div>
                        </div>
                    </CardContent>
                </Card>

                <Card>
                    <CardContent className="p-4">
                        <div className="flex items-center justify-between">
                            <div>
                                <p className="text-sm text-muted-foreground">
                                    紧急任务
                                </p>
                                <p className="text-2xl font-bold text-red-600">
                                    {
                                        tasks.filter(
                                            (t) => t.priority === "CRITICAL",
                                        ).length
                                    }
                                </p>
                            </div>
                            <div className="rounded-full bg-red-100 p-2">
                                <AlertTriangle className="h-5 w-5 text-red-600" />
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </div>

            <Tabs defaultValue="tasks" className="w-full">
                <TabsList className="grid w-full grid-cols-4">
                    <TabsTrigger value="tasks">任务管理</TabsTrigger>
                    <TabsTrigger value="resources">资源调度</TabsTrigger>
                    <TabsTrigger value="monitoring">实时监控</TabsTrigger>
                    <TabsTrigger value="analytics">任务分析</TabsTrigger>
                </TabsList>

                <TabsContent value="tasks" className="mt-6">
                    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                        {/* Task List */}
                        <div className="lg:col-span-2 space-y-4">
                            <Card>
                                <CardHeader>
                                    <div className="flex items-center justify-between">
                                        <CardTitle className="flex items-center gap-2">
                                            <FileText className="h-5 w-5" />
                                            任务列表
                                        </CardTitle>
                                        <Button
                                            onClick={() =>
                                                setShowCreateDialog(true)
                                            }
                                        >
                                            <Plus className="mr-2 h-4 w-4" />
                                            创建任务
                                        </Button>
                                    </div>

                                    {/* Filters */}
                                    <div className="flex items-center gap-4 mt-4">
                                        <div className="flex-1">
                                            <Input
                                                placeholder="搜索任务..."
                                                value={searchTerm}
                                                onChange={(e) =>
                                                    setSearchTerm(
                                                        e.target.value,
                                                    )
                                                }
                                                className="max-w-sm"
                                            />
                                        </div>
                                        <Select
                                            value={taskFilter}
                                            onValueChange={setTaskFilter}
                                        >
                                            <SelectTrigger className="w-32">
                                                <SelectValue />
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="all">
                                                    所有类型
                                                </SelectItem>
                                                {Object.entries(TASK_TYPES).map(
                                                    ([key, type]) => (
                                                        <SelectItem
                                                            key={key}
                                                            value={key}
                                                        >
                                                            {type.name}
                                                        </SelectItem>
                                                    ),
                                                )}
                                            </SelectContent>
                                        </Select>
                                    </div>
                                </CardHeader>
                                <CardContent>
                                    <div className="space-y-3">
                                        {filteredTasks.map((task) => (
                                            <div
                                                key={task.id}
                                                className={`p-4 border rounded-lg cursor-pointer transition-colors ${
                                                    selectedTask.id === task.id
                                                        ? "border-blue-500 bg-blue-50"
                                                        : "hover:bg-gray-50"
                                                }`}
                                                onClick={() =>
                                                    setSelectedTask(task)
                                                }
                                            >
                                                <div className="flex items-start justify-between mb-2">
                                                    <div className="flex-1">
                                                        <div className="flex items-center gap-2 mb-1">
                                                            <span className="text-sm font-mono text-muted-foreground">
                                                                {task.id}
                                                            </span>
                                                            <Badge
                                                                variant={getStatusBadgeVariant(
                                                                    task.priority,
                                                                )}
                                                                className={`bg-${getPriorityColor(task.priority)}-500 text-white`}
                                                            >
                                                                {
                                                                    TASK_PRIORITIES[
                                                                        task.priority as keyof typeof TASK_PRIORITIES
                                                                    ]?.name
                                                                }
                                                            </Badge>
                                                            <Badge variant="outline">
                                                                {
                                                                    TASK_STATUS[
                                                                        task.status as keyof typeof TASK_STATUS
                                                                    ]?.name
                                                                }
                                                            </Badge>
                                                        </div>
                                                        <h4 className="font-medium">
                                                            {task.title}
                                                        </h4>
                                                        <p className="text-sm text-muted-foreground mt-1">
                                                            {task.description}
                                                        </p>
                                                        <div className="flex items-center gap-4 mt-2 text-xs text-muted-foreground">
                                                            <span className="flex items-center gap-1">
                                                                <MapPin className="h-3 w-3" />
                                                                {
                                                                    task
                                                                        .location
                                                                        .name
                                                                }
                                                            </span>
                                                            <span className="flex items-center gap-1">
                                                                <Clock className="h-3 w-3" />
                                                                {
                                                                    task
                                                                        .schedule
                                                                        .startTime
                                                                }
                                                            </span>
                                                        </div>
                                                    </div>
                                                    <div className="text-right">
                                                        <div className="text-lg font-bold">
                                                            {task.progress}%
                                                        </div>
                                                        <Progress
                                                            value={
                                                                task.progress
                                                            }
                                                            className="w-16 mt-1"
                                                        />
                                                    </div>
                                                </div>

                                                {/* Task Actions */}
                                                <div className="flex items-center gap-2 mt-3">
                                                    {task.status ===
                                                        "PENDING" && (
                                                        <Button
                                                            size="sm"
                                                            onClick={(e) => {
                                                                e.stopPropagation();
                                                                handleTaskAction(
                                                                    task.id,
                                                                    "start",
                                                                );
                                                            }}
                                                        >
                                                            <Play className="mr-1 h-3 w-3" />
                                                            开始
                                                        </Button>
                                                    )}
                                                    {task.status ===
                                                        "IN_PROGRESS" && (
                                                        <>
                                                            <Button
                                                                size="sm"
                                                                variant="outline"
                                                                onClick={(
                                                                    e,
                                                                ) => {
                                                                    e.stopPropagation();
                                                                    handleTaskAction(
                                                                        task.id,
                                                                        "pause",
                                                                    );
                                                                }}
                                                            >
                                                                <Pause className="mr-1 h-3 w-3" />
                                                                暂停
                                                            </Button>
                                                            <Button
                                                                size="sm"
                                                                onClick={(
                                                                    e,
                                                                ) => {
                                                                    e.stopPropagation();
                                                                    handleTaskAction(
                                                                        task.id,
                                                                        "complete",
                                                                    );
                                                                }}
                                                            >
                                                                <CheckCircle className="mr-1 h-3 w-3" />
                                                                完成
                                                            </Button>
                                                        </>
                                                    )}
                                                    {(task.status ===
                                                        "PENDING" ||
                                                        task.status ===
                                                            "ASSIGNED") && (
                                                        <Button
                                                            size="sm"
                                                            variant="destructive"
                                                            onClick={(e) => {
                                                                e.stopPropagation();
                                                                handleTaskAction(
                                                                    task.id,
                                                                    "cancel",
                                                                );
                                                            }}
                                                        >
                                                            <Square className="mr-1 h-3 w-3" />
                                                            取消
                                                        </Button>
                                                    )}
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                </CardContent>
                            </Card>
                        </div>

                        {/* Task Details */}
                        <div className="space-y-4">
                            <Card>
                                <CardHeader>
                                    <CardTitle className="text-lg">
                                        {selectedTask.title}
                                    </CardTitle>
                                    <div className="flex items-center gap-2">
                                        <Badge
                                            variant={getStatusBadgeVariant(
                                                selectedTask.priority,
                                            )}
                                            className={`bg-${getPriorityColor(selectedTask.priority)}-500 text-white`}
                                        >
                                            {
                                                TASK_PRIORITIES[
                                                    selectedTask.priority as keyof typeof TASK_PRIORITIES
                                                ]?.name
                                            }
                                        </Badge>
                                        <Badge variant="outline">
                                            {
                                                TASK_STATUS[
                                                    selectedTask.status as keyof typeof TASK_STATUS
                                                ]?.name
                                            }
                                        </Badge>
                                    </div>
                                </CardHeader>
                                <CardContent className="space-y-4">
                                    <div>
                                        <span className="text-sm text-muted-foreground">
                                            任务描述
                                        </span>
                                        <p className="mt-1 text-sm">
                                            {selectedTask.description}
                                        </p>
                                    </div>

                                    <div className="grid grid-cols-2 gap-4">
                                        <div>
                                            <span className="text-sm text-muted-foreground">
                                                执行地点
                                            </span>
                                            <p className="mt-1 text-sm font-medium">
                                                {selectedTask.location.name}
                                            </p>
                                        </div>
                                        <div>
                                            <span className="text-sm text-muted-foreground">
                                                执行进度
                                            </span>
                                            <div className="mt-1">
                                                <Progress
                                                    value={
                                                        selectedTask.progress
                                                    }
                                                />
                                                <p className="text-sm font-medium mt-1">
                                                    {selectedTask.progress}%
                                                </p>
                                            </div>
                                        </div>
                                    </div>

                                    <div>
                                        <span className="text-sm text-muted-foreground">
                                            分配资源
                                        </span>
                                        <div className="mt-1 space-y-1">
                                            <p className="text-sm">
                                                无人机:{" "}
                                                {selectedTask.assignedTo.drone}
                                            </p>
                                            <p className="text-sm">
                                                操作员:{" "}
                                                {
                                                    selectedTask.assignedTo
                                                        .operator
                                                }
                                            </p>
                                            <p className="text-sm">
                                                数据分析师:{" "}
                                                {
                                                    selectedTask.assignedTo
                                                        .dataHandler
                                                }
                                            </p>
                                        </div>
                                    </div>

                                    <div>
                                        <span className="text-sm text-muted-foreground">
                                            时间安排
                                        </span>
                                        <div className="mt-1 space-y-1">
                                            <p className="text-sm">
                                                开始时间:{" "}
                                                {
                                                    selectedTask.schedule
                                                        .startTime
                                                }
                                            </p>
                                            <p className="text-sm">
                                                结束时间:{" "}
                                                {selectedTask.schedule.endTime}
                                            </p>
                                            <p className="text-sm">
                                                预计时长:{" "}
                                                {
                                                    selectedTask.schedule
                                                        .estimatedDuration
                                                }
                                            </p>
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>
                        </div>
                    </div>
                </TabsContent>

                <TabsContent value="resources" className="mt-6">
                    <Card>
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2">
                                <Users className="h-5 w-5" />
                                资源调度管理
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="text-center py-10">
                                <Users className="h-12 w-12 text-muted-foreground mb-4 mx-auto" />
                                <h3 className="font-medium">资源调度功能</h3>
                                <p className="text-sm text-muted-foreground mt-1 max-w-md mx-auto">
                                    此模块用于管理和调度无人机、操作员和数据分析师等资源
                                </p>
                                <Button className="mt-4">配置资源调度</Button>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="monitoring" className="mt-6">
                    <Card>
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2">
                                <Activity className="h-5 w-5" />
                                实时监控
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="text-center py-10">
                                <Activity className="h-12 w-12 text-muted-foreground mb-4 mx-auto" />
                                <h3 className="font-medium">实时任务监控</h3>
                                <p className="text-sm text-muted-foreground mt-1 max-w-md mx-auto">
                                    实时监控所有正在执行的任务状态和进度
                                </p>
                                <Button className="mt-4">启动监控</Button>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>

                <TabsContent value="analytics" className="mt-6">
                    <Card>
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2">
                                <Target className="h-5 w-5" />
                                任务分析
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="text-center py-10">
                                <Target className="h-12 w-12 text-muted-foreground mb-4 mx-auto" />
                                <h3 className="font-medium">任务绩效分析</h3>
                                <p className="text-sm text-muted-foreground mt-1 max-w-md mx-auto">
                                    分析任务执行效率、资源利用率和完成质量
                                </p>
                                <Button className="mt-4">生成分析报告</Button>
                            </div>
                        </CardContent>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    );
}
