package models

import (
	"sync"
	"time"
)

// ScanRequest 扫描请求
type ScanRequest struct {
	ID          string            `json:"id"`
	Targets     []string          `json:"targets"`
	Options     ScanOptions       `json:"options"`
	Timestamp   time.Time         `json:"timestamp"`
	UserAgent   string            `json:"user_agent"`
	Headers     map[string]string `json:"headers,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ScanOptions 扫描选项
type ScanOptions struct {
	Threads         int           `json:"threads"`
	Timeout         time.Duration `json:"timeout"`
	MaxRedirects    int           `json:"max_redirects"`
	RiskLevels      []int         `json:"risk_levels"`
	EnabledPlugins  []string      `json:"enabled_plugins,omitempty"`
	DisabledPlugins []string      `json:"disabled_plugins,omitempty"`
	ScanLevel       int           `json:"scan_level"`
	UserAgent       string        `json:"user_agent"`
	Proxy           ProxyConfig   `json:"proxy,omitempty"`
	OutputFormat    string        `json:"output_format"`
	Verbose         bool          `json:"verbose"`
	Debug           bool          `json:"debug"`
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	Enabled  bool   `json:"enabled"`
	Address  string `json:"address,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// ScanResult 扫描结果
type ScanResult struct {
	ID              string                    `json:"id"`
	StartTime       time.Time                 `json:"start_time"`
	EndTime         time.Time                 `json:"end_time"`
	Duration        time.Duration             `json:"duration"`
	Status          ScanStatus                `json:"status"`
	Targets         []*TargetResult           `json:"targets"`
	Vulnerabilities []*Vulnerability          `json:"vulnerabilities"`
	Statistics      *ScanStatistics           `json:"statistics"`
	Metadata        map[string]string         `json:"metadata,omitempty"`
	
	// 内部状态管理
	mutex           sync.RWMutex              `json:"-"`
}

// ScanStatus 扫描状态
type ScanStatus string

const (
	StatusPending   ScanStatus = "pending"
	StatusRunning   ScanStatus = "running"
	StatusCompleted ScanStatus = "completed"
	StatusFailed    ScanStatus = "failed"
	StatusCancelled ScanStatus = "cancelled"
)

// TargetResult 目标扫描结果
type TargetResult struct {
	URL          string            `json:"url"`
	Status       TargetStatus      `json:"status"`
	ResponseTime time.Duration     `json:"response_time"`
	StatusCode   int               `json:"status_code,omitempty"`
	ContentType  string            `json:"content_type,omitempty"`
	ContentSize  int64             `json:"content_size,omitempty"`
	Technologies []string          `json:"technologies,omitempty"`
	Errors       []string          `json:"errors,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// TargetStatus 目标状态
type TargetStatus string

const (
	TargetStatusPending   TargetStatus = "pending"
	TargetStatusScanning  TargetStatus = "scanning"
	TargetStatusCompleted TargetStatus = "completed"
	TargetStatusFailed    TargetStatus = "failed"
	TargetStatusSkipped   TargetStatus = "skipped"
)

// ScanStatistics 扫描统计信息
type ScanStatistics struct {
	TotalRequests        int                   `json:"total_requests"`
	TotalVulns           int                   `json:"total_vulnerabilities"`
	VulnsBySeverity      map[Severity]int      `json:"vulnerabilities_by_severity"`
	VulnsByCategory      map[Category]int      `json:"vulnerabilities_by_category"`
	VulnsByType          map[VulnType]int      `json:"vulnerabilities_by_type"`
	TargetsScanned       int                   `json:"targets_scanned"`
	TargetsWithVulns     int                   `json:"targets_with_vulnerabilities"`
	PluginsExecuted      map[string]int        `json:"plugins_executed,omitempty"`
	AvgResponseTime      time.Duration         `json:"avg_response_time"`
	ScanEfficiency       float64               `json:"scan_efficiency"`
	CoverageScore        float64               `json:"coverage_score"`
}

// NewScanResult 创建新的扫描结果
func NewScanResult(id string) *ScanResult {
	return &ScanResult{
		ID:              id,
		StartTime:       time.Now(),
		Status:          StatusPending,
		Targets:         make([]*TargetResult, 0),
		Vulnerabilities: make([]*Vulnerability, 0),
		Statistics: &ScanStatistics{
			VulnsBySeverity: make(map[Severity]int),
			VulnsByCategory: make(map[Category]int),
			VulnsByType:     make(map[VulnType]int),
			PluginsExecuted: make(map[string]int),
		},
		Metadata: make(map[string]string),
	}
}

// AddVulnerability 添加漏洞
func (sr *ScanResult) AddVulnerability(vuln *Vulnerability) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	
	sr.Vulnerabilities = append(sr.Vulnerabilities, vuln)
	sr.updateStatistics()
}

// AddTarget 添加目标
func (sr *ScanResult) AddTarget(target *TargetResult) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	
	sr.Targets = append(sr.Targets, target)
}

// UpdateTarget 更新目标状态
func (sr *ScanResult) UpdateTarget(url string, status TargetStatus) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	
	for _, target := range sr.Targets {
		if target.URL == url {
			target.Status = status
			break
		}
	}
}

// GetVulnerabilities 获取所有漏洞
func (sr *ScanResult) GetVulnerabilities() []*Vulnerability {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()
	
	// 返回副本以避免并发修改
	vulns := make([]*Vulnerability, len(sr.Vulnerabilities))
	copy(vulns, sr.Vulnerabilities)
	return vulns
}

// GetVulnerabilityCount 获取漏洞总数
func (sr *ScanResult) GetVulnerabilityCount() int {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()
	
	return len(sr.Vulnerabilities)
}

// GetVulnerabilitiesBySeverity 按严重程度获取漏洞
func (sr *ScanResult) GetVulnerabilitiesBySeverity(severity Severity) []*Vulnerability {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()
	
	var result []*Vulnerability
	for _, vuln := range sr.Vulnerabilities {
		if vuln.Severity == severity {
			result = append(result, vuln)
		}
	}
	return result
}

// HasVulnerabilities 检查是否存在漏洞
func (sr *ScanResult) HasVulnerabilities() bool {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()
	
	return len(sr.Vulnerabilities) > 0
}

// SetCompleted 设置扫描完成
func (sr *ScanResult) SetCompleted() {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	
	sr.EndTime = time.Now()
	sr.Duration = sr.EndTime.Sub(sr.StartTime)
	sr.Status = StatusCompleted
	sr.updateStatistics()
}

// SetFailed 设置扫描失败
func (sr *ScanResult) SetFailed() {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	
	sr.EndTime = time.Now()
	sr.Duration = sr.EndTime.Sub(sr.StartTime)
	sr.Status = StatusFailed
}

// SetRunning 设置扫描运行中
func (sr *ScanResult) SetRunning() {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()
	
	sr.Status = StatusRunning
}

// updateStatistics 更新统计信息
func (sr *ScanResult) updateStatistics() {
	stats := sr.Statistics
	
	// 更新漏洞总数
	stats.TotalVulns = len(sr.Vulnerabilities)
	
	// 重置计数器
	for k := range stats.VulnsBySeverity {
		delete(stats.VulnsBySeverity, k)
	}
	for k := range stats.VulnsByCategory {
		delete(stats.VulnsByCategory, k)
	}
	for k := range stats.VulnsByType {
		delete(stats.VulnsByType, k)
	}
	
	// 统计漏洞分布
	for _, vuln := range sr.Vulnerabilities {
		stats.VulnsBySeverity[vuln.Severity]++
		stats.VulnsByCategory[vuln.Category]++
		stats.VulnsByType[vuln.Type]++
	}
	
	// 统计目标相关信息
	stats.TargetsScanned = len(sr.Targets)
	
	vulnTargets := make(map[string]bool)
	for _, vuln := range sr.Vulnerabilities {
		vulnTargets[vuln.URL] = true
	}
	stats.TargetsWithVulns = len(vulnTargets)
	
	// 计算平均响应时间
	var totalResponseTime time.Duration
	var successfulRequests int
	for _, target := range sr.Targets {
		if target.Status == TargetStatusCompleted && target.ResponseTime > 0 {
			totalResponseTime += target.ResponseTime
			successfulRequests++
		}
	}
	if successfulRequests > 0 {
		stats.AvgResponseTime = totalResponseTime / time.Duration(successfulRequests)
	}
	
	// 计算扫描效率 (漏洞数/扫描时间)
	if sr.Duration > 0 {
		stats.ScanEfficiency = float64(stats.TotalVulns) / sr.Duration.Seconds()
	}
	
	// 计算覆盖率评分
	if stats.TargetsScanned > 0 {
		stats.CoverageScore = float64(stats.TargetsWithVulns) / float64(stats.TargetsScanned) * 100
		if stats.CoverageScore > 100 {
			stats.CoverageScore = 100
		}
	}
}

// GetSummary 获取扫描摘要
func (sr *ScanResult) GetSummary() map[string]interface{} {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()
	
	summary := map[string]interface{}{
		"total_vulnerabilities": sr.Statistics.TotalVulns,
		"targets_scanned":      sr.Statistics.TargetsScanned,
		"scan_duration":        sr.Duration.String(),
		"status":              string(sr.Status),
	}
	
	// 按严重程度统计
	for severity := SeverityInfo; severity <= SeverityCritical; severity++ {
		count := sr.Statistics.VulnsBySeverity[severity]
		summary[severity.String()] = count
	}
	
	return summary
}