package rout

import (
	"goshortly/util"
	"log"

	"golang.org/x/crypto/bcrypt"
	"goshortly/db"
	"goshortly/models"
)

// AutoMigrate ensures the AdminConfig schema is applied
func AutoMigrate() {
	modelsToMigrate := []interface{}{
		&models.AdminConfig{},
		&models.Invite{},
	}
	if err := db.DB.AutoMigrate(modelsToMigrate...); err != nil {
		log.Fatalf("❌ Failed to migrate models: %v", err)
	}
}

// Load returns the first AdminConfig record
func Load() (*models.AdminConfig, error) {
	cfg := &models.AdminConfig{}
	err := db.DB.First(cfg).Error
	return cfg, err
}

// Save persists changes to the AdminConfig
func Save(input *models.AdminConfig) error {
	var existing models.AdminConfig
	if err := db.DB.First(&existing, 1).Error; err != nil {
		return err
	}

	util.CopyFieldsExcept(&existing, input, []string{
		"AdminPassHash", "ID", "CreatedAt", "UpdatedAt", "DeletedAt",
	})
	existing.ID = 1
	return db.DB.Save(&existing).Error
}

// InitDefaults seeds a default config if none exists
func InitDefaults() {
	var count int64
	db.DB.Model(&models.AdminConfig{}).Count(&count)
	if count > 0 {
		return
	}

	log.Println("🟡 Seeding default AdminConfig...")

	passHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Failed to hash admin password: %v", err)
	}

	defaultCfg := models.AdminConfig{
		AdminPassHash:         string(passHash),
		BaseURL:               "http://localhost:8080",
		MaxURLLength:          10,
		BrandingName:          "GoShort.ly",
		DefaultUserLimit:      10,
		DefaultExpirationDays: 90,
		FailedLoginLimit:      5,
		SMTPSecurity:          "starttls",
		SMTPPort:              587,
		SMTPAuthMethod:        "login",
	}

	if err := db.DB.Create(&defaultCfg).Error; err != nil {
		log.Fatalf("❌ Failed to seed AdminConfig: %v", err)
	}

	log.Println("✅ AdminConfig seeded.")
}
