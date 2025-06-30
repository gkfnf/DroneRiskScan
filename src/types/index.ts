export interface Asset {
  id: string
  name: string
  type: "drone" | "gcs" | "server" | "network" | "hangar"
  ipAddress: string
  macAddress?: string
  location?: string
  zone?: string
  status: "online" | "offline" | "scanning" | "maintenance"
  lastSeen?: Date
  metadata?: Record<string, any>
  createdAt: Date
  updatedAt: Date
}

export interface ScanTask {
  id: string
  assetId: string
  scanType: string
  status: "pending" | "running" | "completed" | "failed"
  progress: number
  startedAt?: Date
  completedAt?: Date
  errorMessage?: string
  config?: ScanConfig
  createdAt: Date
}

export interface Vulnerability {
  id: string
  assetId: string
  scanTaskId?: string
  cveId?: string
  title: string
  description: string
  severity: "critical" | "high" | "medium" | "low"
  cvssScore?: number
  cvssVector?: string
  category: string
  status: "open" | "in-progress" | "resolved" | "false-positive"
  discoveredAt: Date
  updatedAt: Date
  remediation?: string
  references?: string[]
}

export interface RFSignal {
  id: string
  frequency: number
  strength: number
  signalType: "GPS" | "RC" | "Video" | "Telemetry" | "WiFi" | "Bluetooth" | "Unknown"
  source?: string
  status: "normal" | "suspicious" | "threat"
  detectedAt: Date
  location?: string
  metadata?: Record<string, any>
}

export interface RFThreat {
  id: string
  threatType: "jamming" | "spoofing" | "interception" | "unauthorized"
  frequency: number
  severity: "low" | "medium" | "high" | "critical"
  description: string
  detectedAt: Date
  affectedSystems: string[]
  status: "active" | "mitigated" | "resolved"
  mitigation?: string
}

export interface ScanConfig {
  templates?: string[]
  targets?: string[]
  concurrency?: number
  rateLimit?: number
  timeout?: number
  customParams?: Record<string, string>
}

export interface LocationData {
  latitude: number
  longitude: number
  accuracy: number
  altitude?: number
  heading?: number
  speed?: number
  timestamp: number
}

export interface LocationRecord {
  id: string
  vulnerabilityId: string
  location: LocationData
  address: string
  engineer: string
  processingTime: Date
  notes?: string
  photos?: string[]
}

export interface SystemStats {
  totalAssets: number
  onlineAssets: number
  totalVulnerabilities: number
  criticalVulnerabilities: number
  highVulnerabilities: number
  mediumVulnerabilities: number
  lowVulnerabilities: number
  rfThreats: number
  runningScanTasks: number
}

export interface ApiResponse<T> {
  data: T
  message?: string
  success: boolean
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
  hasNext: boolean
  hasPrev: boolean
}
