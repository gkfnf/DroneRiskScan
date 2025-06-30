package models

import (
	"time"
)

// Asset 资产模型
type Asset struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Type        string    `json:"type" db:"type"`
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	MACAddress  string    `json:"mac_address" db:"mac_address"`
	Location    string    `json:"location" db:"location"`
	Zone        string    `json:"zone" db:"zone"`
	Status      string    `json:"status" db:"status"`
	LastSeen    time.Time `json:"last_seen" db:"last_seen"`
	Metadata    string    `json:"metadata" db:"metadata"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ScanTask 扫描任务模型
type ScanTask struct {
	ID           string     `json:"id" db:"id"`
	AssetID      string     `json:"asset_id" db:"asset_id"`
	ScanType     string     `json:"scan_type" db:"scan_type"`
	Status       string     `json:"status" db:"status"`
	Progress     int        `json:"progress" db:"progress"`
	StartedAt    *time.Time `json:"started_at" db:"started_at"`
	CompletedAt  *time.Time `json:"completed_at" db:"completed_at"`
	ErrorMessage string     `json:"error_message" db:"error_message"`
	Config       string     `json:"config" db:"config"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}

// Vulnerability 漏洞模型
type Vulnerability struct {
	ID           string    `json:"id" db:"id"`
	AssetID      string    `json:"asset_id" db:"asset_id"`
	ScanTaskID   string    `json:"scan_task_id" db:"scan_task_id"`
	CVEID        string    `json:"cve_id" db:"cve_id"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	Severity     string    `json:"severity" db:"severity"`
	CVSSScore    float64   `json:"cvss_score" db:"cvss_score"`
	CVSSVector   string    `json:"cvss_vector" db:"cvss_vector"`
	Category     string    `json:"category" db:"category"`
	Status       string    `json:"status" db:"status"`
	DiscoveredAt time.Time `json:"discovered_at" db:"discovered_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Remediation  string    `json:"remediation" db:"remediation"`
	References   string    `json:"references" db:"references"`
}

// RFSignal 射频信号模型
type RFSignal struct {
	ID         string    `json:"id" db:"id"`
	Frequency  float64   `json:"frequency" db:"frequency"`
	Strength   float64   `json:"strength" db:"strength"`
	SignalType string    `json:"signal_type" db:"signal_type"`
	Source     string    `json:"source" db:"source"`
	Status     string    `json:"status" db:"status"`
	DetectedAt time.Time `json:"detected_at" db:"detected_at"`
	Location   string    `json:"location" db:"location"`
	Metadata   string    `json:"metadata" db:"metadata"`
}

// RFThreat 射频威胁模型
type RFThreat struct {
	ID              string    `json:"id" db:"id"`
	ThreatType      string    `json:"threat_type" db:"threat_type"`
	Frequency       float64   `json:"frequency" db:"frequency"`
	Severity        string    `json:"severity" db:"severity"`
	Description     string    `json:"description" db:"description"`
	DetectedAt      time.Time `json:"detected_at" db:"detected_at"`
	AffectedSystems string    `json:"affected_systems" db:"affected_systems"`
	Status          string    `json:"status" db:"status"`
	Mitigation      string    `json:"mitigation" db:"mitigation"`
}

// LocationRecord 位置记录模型
type LocationRecord struct {
	ID              string    `json:"id" db:"id"`
	VulnerabilityID string    `json:"vulnerability_id" db:"vulnerability_id"`
	Latitude        float64   `json:"latitude" db:"latitude"`
	Longitude       float64   `json:"longitude" db:"longitude"`
	Accuracy        float64   `json:"accuracy" db:"accuracy"`
	Altitude        float64   `json:"altitude" db:"altitude"`
	Address         string    `json:"address" db:"address"`
	Engineer        string    `json:"engineer" db:"engineer"`
	ProcessingTime  time.Time `json:"processing_time" db:"processing_time"`
	Notes           string    `json:"notes" db:"notes"`
	Photos          string    `json:"photos" db:"photos"`
}

// ScanConfig 扫描配置
type ScanConfig struct {
	Templates    []string          `json:"templates"`
	Targets      []string          `json:"targets"`
	Concurrency  int               `json:"concurrency"`
	RateLimit    int               `json:"rate_limit"`
	Timeout      int               `json:"timeout"`
	CustomParams map[string]string `json:"custom_params"`
}

// ScanResult 扫描结果
type ScanResult struct {
	TaskID          string          `json:"task_id"`
	AssetID         string          `json:"asset_id"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Summary         ScanSummary     `json:"summary"`
	Duration        time.Duration   `json:"duration"`
	Error           string          `json:"error,omitempty"`
}

// ScanSummary 扫描摘要
type ScanSummary struct {
	TotalChecks      int `json:"total_checks"`
	VulnerabilitiesFound int `json:"vulnerabilities_found"`
	CriticalCount    int `json:"critical_count"`
	HighCount        int `json:"high_count"`
	MediumCount      int `json:"medium_count"`
	LowCount         int `json:"low_count"`
}

// RFScanConfig 射频扫描配置
type RFScanConfig struct {
	FrequencyStart float64 `json:"frequency_start"`
	FrequencyEnd   float64 `json:"frequency_end"`
	SampleRate     int     `json:"sample_rate"`
	Gain           float64 `json:"gain"`
	Duration       int     `json:"duration"`
	Threshold      float64 `json:"threshold"`
}

// RFScanResult 射频扫描结果
type RFScanResult struct {
	Signals []RFSignal `json:"signals"`
	Threats []RFThreat `json:"threats"`
	Summary RFSummary  `json:"summary"`
}

// RFSummary 射频扫描摘要
type RFSummary struct {
	TotalSignals    int `json:"total_signals"`
	ThreatsDetected int `json:"threats_detected"`
	FrequencyRange  struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"frequency_range"`
	ScanDuration time.Duration `json:"scan_duration"`
}
