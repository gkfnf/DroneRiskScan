package services

import (
	"database/sql"
	"fmt"
	"time"

	"drone-security-scanner/internal/database"
	"drone-security-scanner/internal/models"
	"github.com/google/uuid"
)

type AssetService struct {
	db *database.DB
}

func NewAssetService(db *database.DB) *AssetService {
	return &AssetService{db: db}
}

func (s *AssetService) CreateAsset(asset *models.Asset) error {
	if asset.ID == "" {
		asset.ID = uuid.New().String()
	}
	asset.CreatedAt = time.Now()
	asset.UpdatedAt = time.Now()

	query := `
		INSERT INTO assets (id, name, type, ip_address, mac_address, location, zone, status, metadata, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query, asset.ID, asset.Name, asset.Type, asset.IPAddress, asset.MACAddress,
		asset.Location, asset.Zone, asset.Status, asset.Metadata, asset.CreatedAt, asset.UpdatedAt)

	return err
}

func (s *AssetService) GetAsset(id string) (*models.Asset, error) {
	query := `
		SELECT id, name, type, ip_address, mac_address, location, zone, status, 
		       COALESCE(last_seen, ''), metadata, created_at, updated_at
		FROM assets WHERE id = ?
	`

	var asset models.Asset
	var lastSeenStr string

	err := s.db.QueryRow(query, id).Scan(
		&asset.ID, &asset.Name, &asset.Type, &asset.IPAddress, &asset.MACAddress,
		&asset.Location, &asset.Zone, &asset.Status, &lastSeenStr, &asset.Metadata,
		&asset.CreatedAt, &asset.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("asset not found")
		}
		return nil, err
	}

	if lastSeenStr != "" {
		if lastSeen, err := time.Parse(time.RFC3339, lastSeenStr); err == nil {
			asset.LastSeen = lastSeen
		}
	}

	return &asset, nil
}

func (s *AssetService) ListAssets(assetType string, status string, limit, offset int) ([]*models.Asset, error) {
	query := `
		SELECT id, name, type, ip_address, mac_address, location, zone, status,
		       COALESCE(last_seen, ''), metadata, created_at, updated_at
		FROM assets
		WHERE 1=1
	`
	args := []interface{}{}

	if assetType != "" {
		query += " AND type = ?"
		args = append(args, assetType)
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

	var assets []*models.Asset
	for rows.Next() {
		var asset models.Asset
		var lastSeenStr string

		err := rows.Scan(
			&asset.ID, &asset.Name, &asset.Type, &asset.IPAddress, &asset.MACAddress,
			&asset.Location, &asset.Zone, &asset.Status, &lastSeenStr, &asset.Metadata,
			&asset.CreatedAt, &asset.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if lastSeenStr != "" {
			if lastSeen, err := time.Parse(time.RFC3339, lastSeenStr); err == nil {
				asset.LastSeen = lastSeen
			}
		}

		assets = append(assets, &asset)
	}

	return assets, nil
}

func (s *AssetService) UpdateAsset(asset *models.Asset) error {
	asset.UpdatedAt = time.Now()

	query := `
		UPDATE assets 
		SET name = ?, type = ?, ip_address = ?, mac_address = ?, location = ?, 
		    zone = ?, status = ?, metadata = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := s.db.Exec(query, asset.Name, asset.Type, asset.IPAddress, asset.MACAddress,
		asset.Location, asset.Zone, asset.Status, asset.Metadata, asset.UpdatedAt, asset.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset not found")
	}

	return nil
}

func (s *AssetService) DeleteAsset(id string) error {
	query := "DELETE FROM assets WHERE id = ?"

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset not found")
	}

	return nil
}

func (s *AssetService) UpdateAssetStatus(id, status string) error {
	query := "UPDATE assets SET status = ?, last_seen = ?, updated_at = ? WHERE id = ?"

	now := time.Now()
	result, err := s.db.Exec(query, status, now, now, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("asset not found")
	}

	return nil
}

func (s *AssetService) GetAssetsByType(assetType string) ([]*models.Asset, error) {
	return s.ListAssets(assetType, "", 0, 0)
}

func (s *AssetService) GetOnlineAssets() ([]*models.Asset, error) {
	return s.ListAssets("", "online", 0, 0)
}

func (s *AssetService) CountAssets() (int, error) {
	query := "SELECT COUNT(*) FROM assets"

	var count int
	err := s.db.QueryRow(query).Scan(&count)
	return count, err
}

func (s *AssetService) CountAssetsByStatus(status string) (int, error) {
	query := "SELECT COUNT(*) FROM assets WHERE status = ?"

	var count int
	err := s.db.QueryRow(query, status).Scan(&count)
	return count, err
}
