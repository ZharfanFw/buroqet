package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"warehouse-service/internal/domain"
)

type warehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) domain.WarehouseRepository {
	return &warehouseRepository{db: db}
}

func (r *warehouseRepository) SavePackage(ctx context.Context, pkg *domain.Package) error {
	return r.db.WithContext(ctx).Create(pkg).Error
}

func (r *warehouseRepository) GetPackageByAWB(ctx context.Context, awb string) (*domain.Package, error) {
	var pkg domain.Package
	err := r.db.WithContext(ctx).Where("awb = ?", awb).First(&pkg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("package not found")
		}
		return nil, err
	}
	return &pkg, nil
}

func (r *warehouseRepository) GetAllPackages(ctx context.Context) ([]domain.Package, error) {
	var packages []domain.Package
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&packages).Error
	return packages, err
}

func (r *warehouseRepository) UpdatePackageStatus(ctx context.Context, awb string, status string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Package{}).
		Where("awb = ?", awb).
		Update("status", status).Error
}

func (r *warehouseRepository) CreateManifest(ctx context.Context, manifest *domain.Manifest) error {
	return r.db.WithContext(ctx).Create(manifest).Error
}

func (r *warehouseRepository) GetManifestByID(ctx context.Context, id string) (*domain.Manifest, error) {
	var manifest domain.Manifest
	err := r.db.WithContext(ctx).
		Preload("Packages").
		Where("id = ?", id).
		First(&manifest).Error
	if err != nil {
		return nil, err
	}
	return &manifest, nil
}

func (r *warehouseRepository) AddPackageToManifest(ctx context.Context, awb string, manifestID string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Package{}).
		Where("awb = ?", awb).
		Update("manifest_id", manifestID).Error
}

func (r *warehouseRepository) DispatchManifest(ctx context.Context, manifestID string) ([]string, error) {
	err := r.db.WithContext(ctx).
		Model(&domain.Manifest{}).
		Where("id = ?", manifestID).
		Update("status", "DISPATCHED").Error
	if err != nil {
		return nil, err
	}

	// Ambil semua AWB dalam manifest ini
	var packages []domain.Package
	err = r.db.WithContext(ctx).
		Where("manifest_id = ?", manifestID).
		Find(&packages).Error
	if err != nil {
		return nil, err
	}

	// Update semua package menjadi ON_TRANSIT
	err = r.db.WithContext(ctx).
		Model(&domain.Package{}).
		Where("manifest_id = ?", manifestID).
		Update("status", "ON_TRANSIT").Error
	if err != nil {
		return nil, err
	}

	// Kumpulkan semua AWB untuk dikembalikan
	awbs := make([]string, len(packages))
	for i, pkg := range packages {
		awbs[i] = pkg.AWB
	}
	return awbs, nil
}