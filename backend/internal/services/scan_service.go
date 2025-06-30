package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"drone-security-scanner/internal/database"
	"drone-security-scanner/internal/models"
	"github.com/google/uuid"
)

type ScanService struct {
	db *database.DB
}

func NewScanService(db *database.DB) *ScanService {
	return &ScanService{db: db}
}

func (s *ScanService) CreateScanTask(task *models.ScanTask) error {
	if task.ID == "" {
		task.ID = uuid.New().String()
	}
	task.CreatedAt = time.Now()

	query := `
		INSERT INTO scan_tasks (id, asset_id, scan_type, status, progress, config, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query, task.ID, task.AssetID, task.ScanType, task.Status,
		task.Progress, task.Config, task.CreatedAt)

	return err
}

func (s *ScanService) GetScanTask(id string) (*models.ScanTask, error) {
	query := `
		SELECT id, asset_id, scan_type, status, progress, 
		       started_at, completed_at, error_message, config, created_at
		FROM scan_tasks WHERE id = ?
	`

	var task models.ScanTask
	var startedAt, completedAt sql.NullTime

	err := s.db.QueryRow(query, id).Scan(
		&task.ID, &task.AssetID, &task.ScanType, &task.Status, &task.Progress,
		&startedAt, &completedAt, &task.ErrorMessage, &task.Config, &task.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("scan task not found")
		}
		return nil, err
	}

	if startedAt.Valid {
		task.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}

	return &task, nil
}

func (s *ScanService) ListScanTasks(assetID string, status string, limit, offset int) ([]*models.ScanTask, error) {
	query := `
		SELECT id, asset_id, scan_type, status, progress,
		       started_at, completed_at, error_message, config, created_at
		FROM scan_tasks
		WHERE 1=1
	`
	args := []interface{}{}

	if assetID != "" {
		query += " AND asset_id = ?"
		args = append(args, assetID)
	}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	if offset > 0 {
		query += " OFFSET ?"
		args = append(args, offset)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.ScanTask
	for rows.Next() {
		var task models.ScanTask
		var startedAt, completedAt sql.NullTime

		err := rows.Scan(
			&task.ID, &task.AssetID, &task.ScanType, &task.Status, &task.Progress,
			&startedAt, &completedAt, &task.ErrorMessage, &task.Config, &task.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if startedAt.Valid {
			task.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (s *ScanService) UpdateScanTaskStatus(id, status string, progress int) error {
	query := "UPDATE scan_tasks SET status = ?, progress = ? WHERE id = ?"

	result, err := s.db.Exec(query, status, progress, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("scan task not found")
	}

	return nil
}

func (s *ScanService) StartScanTask(id string) error {
	now := time.Now()
	query := "UPDATE scan_tasks SET status = 'running', started_at = ? WHERE id = ?"

	result, err := s.db.Exec(query, now, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("scan task not found")
	}

	return nil
}

func (s *ScanService) CompleteScanTask(id string, errorMessage string) error {
	now := time.Now()
	status := "completed"
	if errorMessage != "" {
		status = "failed"
	}

	query := "UPDATE scan_tasks SET status = ?, completed_at = ?, error_message = ? WHERE id = ?"

	result, err := s.db.Exec(query, status, now, errorMessage, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("scan task not found")
	}

	return nil
}

func (s *ScanService) GetRunningScanTasks() ([]*models.ScanTask, error) {
	return s.ListScanTasks("", "running", 0, 0)
}

func (s *ScanService) CreateScanTaskWithConfig(assetID, scanType string, config *models.ScanConfig) (*models.ScanTask, error) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	task := &models.ScanTask{
		AssetID:  assetID,
		ScanType: scanType,
		Status:   "pending",
		Progress: 0,
		Config:   string(configJSON),
	}

	if err := s.CreateScanTask(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *ScanService) GetScanTaskConfig(id string) (*models.ScanConfig, error) {
	task, err := s.GetScanTask(id)
	if err != nil {
		return nil, err
	}

	var config models.ScanConfig
	if err := json.Unmarshal([]byte(task.Config), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
