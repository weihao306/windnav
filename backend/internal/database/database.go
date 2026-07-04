package database

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"windnav/backend/internal/config"
	"windnav/backend/internal/model"
)

func Open(cfg config.Config) (*gorm.DB, error) {
	gormConfig := &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)}

	switch cfg.DBDriver {
	case "sqlite":
		if err := os.MkdirAll(filepath.Dir(cfg.SQLitePath), 0o755); err != nil {
			return nil, fmt.Errorf("create sqlite directory: %w", err)
		}
		return gorm.Open(sqlite.Open(cfg.SQLitePath), gormConfig)
	case "postgres":
		return gorm.Open(postgres.Open(cfg.DatabaseURL), gormConfig)
	default:
		return nil, fmt.Errorf("unsupported database driver %q", cfg.DBDriver)
	}
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Site{},
		&model.Tag{},
		&model.Setting{},
		&model.SearchEngine{},
		&model.RefreshToken{},
	)
}

func Seed(db *gorm.DB, cfg config.Config) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	if err := db.Where("slug = ?", "common-tools").Delete(&model.Category{}).Error; err != nil {
		return fmt.Errorf("remove legacy category common-tools: %w", err)
	}

	var user model.User
	err = db.Where("username = ?", cfg.AdminUsername).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = model.User{
			Username:     cfg.AdminUsername,
			PasswordHash: string(passwordHash),
			DisplayName:  "Administrator",
			Role:         "admin",
			IsActive:     true,
		}
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("create admin user: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("load admin user: %w", err)
	}

	defaults := []model.Setting{
		{Key: "site_title", Value: "WindNav", ValueType: "string"},
		{Key: "site_subtitle", Value: "简单轻快的自建导航页", ValueType: "string"},
		{Key: "search_placeholder", Value: "搜索站点、标签或描述", ValueType: "string"},
		{Key: "default_theme", Value: "light", ValueType: "string"},
	}
	for _, setting := range defaults {
		var existing model.Setting
		err := db.Where("key = ?", setting.Key).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&setting).Error; err != nil {
				return fmt.Errorf("create setting %s: %w", setting.Key, err)
			}
		} else if err != nil {
			return fmt.Errorf("load setting %s: %w", setting.Key, err)
		}
	}

	engines := []model.SearchEngine{
		{Name: "百度", Slug: "baidu", SearchURL: "https://www.baidu.com/s?wd={query}", Icon: "百", SortOrder: 10, IsDefault: true, IsVisible: true},
		{Name: "Google", Slug: "google", SearchURL: "https://www.google.com/search?q={query}", Icon: "G", SortOrder: 20, IsDefault: false, IsVisible: true},
	}
	for _, engine := range engines {
		var existing model.SearchEngine
		err := db.Where("slug = ?", engine.Slug).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&engine).Error; err != nil {
				return fmt.Errorf("create search engine %s: %w", engine.Slug, err)
			}
		} else if err != nil {
			return fmt.Errorf("load search engine %s: %w", engine.Slug, err)
		}
	}

	return nil
}
