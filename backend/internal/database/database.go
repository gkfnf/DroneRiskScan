package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func Initialize(dbPath string) (*DB, error) {
	// 确保数据库目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// 打开数据库连接
	sqlDB, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{sqlDB}

	// 创建表
	if err := db.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func (db *DB) createTables() error {
	queries := []string{
		// 资产表
		`CREATE TABLE IF NOT EXISTS assets (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			ip_address TEXT,
			mac_address TEXT,
			location TEXT,
			zone TEXT,
			status TEXT DEFAULT 'offline',
			last_seen DATETIME,
			metadata TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// 扫描任务表
		`CREATE TABLE IF NOT EXISTS scan_tasks (
			id TEXT PRIMARY KEY,
			asset_id TEXT,
			scan_type TEXT NOT NULL,
			status TEXT DEFAULT 'pending',
			progress INTEGER DEFAULT 0,
			started_at DATETIME,
			completed_at DATETIME,
			error_message TEXT,
			config TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (asset_id) REFERENCES assets(id)
		)`,

		// 漏洞表
		`CREATE TABLE IF NOT EXISTS vulnerabilities (
			id TEXT PRIMARY KEY,
			asset_id TEXT,
			scan_task_id TEXT,
			cve_id TEXT,
			title TEXT NOT NULL,
			description TEXT,
			severity TEXT NOT NULL,
			cvss_score REAL,
			cvss_vector TEXT,
			category TEXT,
			status TEXT DEFAULT 'open',
			discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			remediation TEXT,
			references TEXT,
			FOREIGN KEY (asset_id) REFERENCES assets(id),
			FOREIGN KEY (scan_task_id) REFERENCES scan_tasks(id)
		)`,

		// 射频信号表
		`CREATE TABLE IF NOT EXISTS rf_signals (
			id TEXT PRIMARY KEY,
			frequency REAL NOT NULL,
			strength REAL NOT NULL,
			signal_type TEXT,
			source TEXT,
			status TEXT DEFAULT 'normal',
			detected_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			location TEXT,
			metadata TEXT
		)`,

		// 射频威胁表
		`CREATE TABLE IF NOT EXISTS rf_threats (
			id TEXT PRIMARY KEY,
			threat_type TEXT NOT NULL,
			frequency REAL NOT NULL,
			severity TEXT NOT NULL,
			description TEXT,
			detected_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			affected_systems TEXT,
			status TEXT DEFAULT 'active',
			mitigation TEXT
		)`,

		// 位置记录表
		`CREATE TABLE IF NOT EXISTS location_records (
			id TEXT PRIMARY KEY,
			vulnerability_id TEXT,
			latitude REAL NOT NULL,
			longitude REAL NOT NULL,
			accuracy REAL,
			altitude REAL,
			address TEXT,
			engineer TEXT,
			processing_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			notes TEXT,
			photos TEXT,
			FOREIGN KEY (vulnerability_id) REFERENCES vulnerabilities(id)
		)`,

		// 系统配置表
		`CREATE TABLE IF NOT EXISTS system_config (
			key TEXT PRIMARY KEY,
			value TEXT,
			description TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// 审计日志表
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			action TEXT NOT NULL,
			resource_type TEXT,
			resource_id TEXT,
			details TEXT,
			ip_address TEXT,
			user_agent TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	// 创建索引
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_assets_type ON assets(type)",
		"CREATE INDEX IF NOT EXISTS idx_assets_status ON assets(status)",
		"CREATE INDEX IF NOT EXISTS idx_vulnerabilities_severity ON vulnerabilities(severity)",
		"CREATE INDEX IF NOT EXISTS idx_vulnerabilities_status ON vulnerabilities(status)",
		"CREATE INDEX IF NOT EXISTS idx_rf_signals_frequency ON rf_signals(frequency)",
		"CREATE INDEX IF NOT EXISTS idx_rf_threats_severity ON rf_threats(severity)",
		"CREATE INDEX IF NOT EXISTS idx_scan_tasks_status ON scan_tasks(status)",
	}

	for _, index := range indexes {
		if _, err := db.Exec(index); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}
