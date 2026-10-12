package database

import (
	"fmt"
	"log"

	"samplegadai-backend/internal/config"
	"samplegadai-backend/internal/model"
	"samplegadai-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Printf("PostgreSQL connection failed: %v. Falling back to local SQLite database (samplegadai.db)...", err)
		
		sqliteDB, errSqlite := gorm.Open(sqlite.Open("samplegadai.db"), gormConfig)
		if errSqlite != nil {
			log.Fatalf("Failed to initialize SQLite fallback database: %v", errSqlite)
		}
		
		log.Println("SQLite fallback database initialized successfully")
		AutoMigrateAndSeed(sqliteDB)
		return sqliteDB
	}

	log.Println("PostgreSQL Database connection established successfully")
	return db
}

func RunMigrations(db *gorm.DB, migrationFilePath string) {
	log.Println("Initializing PostgreSQL schema and seeds via GORM AutoMigrate...")
	AutoMigrateAndSeed(db)
}

func AutoMigrateAndSeed(db *gorm.DB) {
	log.Println("Enabling PostgreSQL UUID extensions...")
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";")

	log.Println("Running GORM AutoMigrate for models...")
	err := db.AutoMigrate(
		&model.Branch{},
		&model.User{},
		&model.RefreshToken{},
		&model.PasswordReset{},
		&model.Customer{},
		&model.CustomerDocument{},
		&model.AssetCategory{},
		&model.CategoryChecklist{},
		&model.Asset{},
		&model.AssetPhoto{},
		&model.AssetInspection{},
		&model.MarketComparable{},
		&model.Appraisal{},
		&model.ApprovalRule{},
		&model.PawnTransaction{},
		&model.Disbursement{},
		&model.Payment{},
		&model.PawnExtension{},
		&model.AssetRelease{},
		&model.Inventory{},
		&model.AuditLog{},
	)
	if err != nil {
		log.Printf("AutoMigrate error: %v", err)
	} else {
		log.Println("GORM AutoMigrate completed successfully")
	}
	SeedInitialData(db)
}

func ResetAndSeed(db *gorm.DB) {
	log.Println("RESETTING DATABASE DARI 0: Dropping public schema CASCADE...")
	db.Exec("DROP SCHEMA IF EXISTS public CASCADE;")
	db.Exec("CREATE SCHEMA public;")
	db.Exec("GRANT ALL ON SCHEMA public TO postgres;")
	db.Exec("GRANT ALL ON SCHEMA public TO public;")
	log.Println("Schema public recreated clean. Re-running AutoMigrate & Seed...")
	AutoMigrateAndSeed(db)
}

func SeedInitialData(db *gorm.DB) {
	// Seed 5 Kota Padang Branches
	branchID1 := uuid.MustParse("00000000-0000-0000-0000-000000000001") // Nanggalo (Pusat Ops)
	branchID2 := uuid.MustParse("00000000-0000-0000-0000-000000000002") // Pasar Raya
	branchID3 := uuid.MustParse("00000000-0000-0000-0000-000000000003") // Azizi
	branchID4 := uuid.MustParse("00000000-0000-0000-0000-000000000004") // Ulak Karang
	branchID5 := uuid.MustParse("00000000-0000-0000-0000-000000000005") // Andalas

	branchesToSeed := []model.Branch{
		{
			ID:       branchID1,
			Code:     "BR-PDG-001",
			Name:     "Cabang Nanggalo",
			Address:  "Jl. Siteba No. 12, Nanggalo, Kota Padang (Sentral Gudang & Brankas Pusat)",
			Phone:    "0751-441122",
			IsActive: true,
		},
		{
			ID:       branchID2,
			Code:     "BR-PDG-002",
			Name:     "Cabang Pasar Raya",
			Address:  "Jl. M. Yamin No. 45, Pasar Raya, Kota Padang (Sentra Grosir & Emas)",
			Phone:    "0751-442233",
			IsActive: true,
		},
		{
			ID:       branchID3,
			Code:     "BR-PDG-003",
			Name:     "Cabang Azizi",
			Address:  "Jl. Dr. Sutomo No. 88, Azizi, Kota Padang (Sentra Gadget & Elektronik)",
			Phone:    "0751-443344",
			IsActive: true,
		},
		{
			ID:       branchID4,
			Code:     "BR-PDG-004",
			Name:     "Cabang Ulak Karang",
			Address:  "Jl. Khatib Sulaiman / Ulak Karang No. 33, Kota Padang",
			Phone:    "0751-445566",
			IsActive: true,
		},
		{
			ID:       branchID5,
			Code:     "BR-PDG-005",
			Name:     "Cabang Andalas",
			Address:  "Jl. Andalas No. 102, Simpang Haru - Andalas, Kota Padang",
			Phone:    "0751-447788",
			IsActive: true,
		},
	}

	for _, b := range branchesToSeed {
		var existing model.Branch
		if err := db.Where("code = ?", b.Code).First(&existing).Error; err != nil {
			db.Create(&b)
		} else {
			db.Model(&model.Branch{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
				"name":      b.Name,
				"address":   b.Address,
				"phone":     b.Phone,
				"is_active": true,
			})
		}
	}
	log.Println("Seeded 5 Kota Padang branches: Nanggalo, Pasar Raya, Azizi, Ulak Karang, Andalas")

	// Seed User Accounts
	hash, _ := utils.HashPassword("Admin123!")

	usersToSeed := []struct {
		ID       uuid.UUID
		Username string
		FullName string
		Email    string
		Role     string
		BranchID *uuid.UUID
	}{
		{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Username: "admin",
			FullName: "Super Administrator",
			Email:    "admin@samplegadai.com",
			Role:     "SUPER_ADMIN",
			BranchID: nil,
		},
		{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Username: "admin_pusat",
			FullName: "Head Office Super Admin (Pusat Padang)",
			Email:    "pusat@samplegadai.com",
			Role:     "SUPER_ADMIN",
			BranchID: nil,
		},
		{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000010"),
			Username: "admin_nanggalo",
			FullName: "Admin Cabang Nanggalo",
			Email:    "nanggalo@samplegadai.com",
			Role:     "MANAGER",
			BranchID: &branchID1,
		},
		{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000020"),
			Username: "admin_pasarraya",
			FullName: "Admin Cabang Pasar Raya",
			Email:    "pasarraya@samplegadai.com",
			Role:     "MANAGER",
			BranchID: &branchID2,
		},
		{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000030"),
			Username: "admin_azizi",
			FullName: "Admin Cabang Azizi",
			Email:    "azizi@samplegadai.com",
			Role:     "MANAGER",
			BranchID: &branchID3,
		},
		{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000040"),
			Username: "admin_ulakkarang",
			FullName: "Admin Cabang Ulak Karang",
			Email:    "ulakkarang@samplegadai.com",
			Role:     "MANAGER",
			BranchID: &branchID4,
		},
		{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000050"),
			Username: "admin_andalas",
			FullName: "Admin Cabang Andalas",
			Email:    "andalas@samplegadai.com",
			Role:     "MANAGER",
			BranchID: &branchID5,
		},
		{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000060"),
			Username: "penaksir_padang",
			FullName: "Penaksir Ahli Syariah (Certified)",
			Email:    "penaksir@samplegadai.com",
			Role:     "APPRAISER",
			BranchID: &branchID1,
		},
		{
			ID:       uuid.MustParse("00000000-0000-0000-0000-000000000070"),
			Username: "kasir_padang",
			FullName: "Kasir Operasional & Pencairan",
			Email:    "kasir@samplegadai.com",
			Role:     "CASHIER",
			BranchID: &branchID1,
		},
	}

	for _, u := range usersToSeed {
		var existing model.User
		if err := db.Where("username = ?", u.Username).First(&existing).Error; err != nil {
			userID := u.ID
			var idCheck model.User
			if db.Where("id = ?", userID).First(&idCheck).Error == nil {
				userID = uuid.New()
			}

			errCreate := db.Create(&model.User{
				ID:           userID,
				BranchID:     u.BranchID,
				Username:     u.Username,
				Email:        u.Email,
				PasswordHash: hash,
				FullName:     u.FullName,
				Role:         u.Role,
				IsActive:     true,
			}).Error
			if errCreate != nil {
				log.Printf("Error creating user %s: %v\n", u.Username, errCreate)
			} else {
				log.Printf("Seeded user: %s (%s / Admin123!)\n", u.FullName, u.Username)
			}
		} else {
			db.Model(&model.User{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
				"password_hash": hash,
				"full_name":     u.FullName,
				"role":          u.Role,
				"branch_id":     u.BranchID,
				"is_active":     true,
			})
			log.Printf("Updated user password hash & info for %s (%s)\n", u.FullName, u.Username)
		}
	}

	// Seed Asset Categories
	catID1 := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	catID2 := uuid.MustParse("00000000-0000-0000-0000-000000000004")
	catID3 := uuid.MustParse("00000000-0000-0000-0000-000000000005")
	catID4 := uuid.MustParse("00000000-0000-0000-0000-000000000006")

	var catCount int64
	db.Model(&model.AssetCategory{}).Count(&catCount)
	if catCount == 0 {
		db.Create(&model.AssetCategory{ID: catID1, Name: "Smartphone", LTVPercent: 70.00, IsActive: true})
		db.Create(&model.AssetCategory{ID: catID2, Name: "Laptop", LTVPercent: 65.00, IsActive: true})
		db.Create(&model.AssetCategory{ID: catID3, Name: "Kamera", LTVPercent: 60.00, IsActive: true})
		db.Create(&model.AssetCategory{ID: catID4, Name: "Emas", LTVPercent: 85.00, IsActive: true})
	}

	// Ensure default category checklists exist for each category
	categoryChecklistsNameMap := map[string][]string{
		"Smartphone": {"Layar & LCD", "Touchscreen", "Kamera Depan", "Kamera Belakang", "Speaker & Microphone", "Kesehatan Baterai", "Port Charging", "WiFi & Bluetooth", "Sinyal SIM", "Fingerprint / Face ID", "Kondisi Fisik Casing", "Kesesuaian IMEI"},
		"Laptop":     {"Display & LCD", "Keyboard & Trackpad", "Charger Original", "Daya Baterai", "Port USB / HDMI / Type-C", "Kekuatan Engsel / Hinge", "WiFi, Audio & Webcam", "Kondisi Fisik Casing"},
		"Kamera":     {"Sensor & Lensa Optik", "Flash & Autofocus System", "Tombol Shutter & Controller", "Layar LCD Flip / Viewfinder", "Slot Memory Card", "Baterai & Charger Original", "Kondisi Fisik Body"},
		"Emas":       {"Uji Kadar Karat (XRF/Batu Uji)", "Berat Timbangan Digital (Gram)", "Cap / Stamp Produsen", "Kejernihan & Bebas Cacat", "Kuitansi & Sertifikat Original"},
	}

	for catName, items := range categoryChecklistsNameMap {
		var cat model.AssetCategory
		if err := db.Where("name = ?", catName).First(&cat).Error; err == nil {
			var existingCount int64
			db.Model(&model.CategoryChecklist{}).Where("category_id = ?", cat.ID).Count(&existingCount)
			if existingCount == 0 {
				for _, item := range items {
					db.Create(&model.CategoryChecklist{
						CategoryID: cat.ID,
						ItemName:   item,
					})
				}
			}
		}
	}

	// Seed Approval Rules
	var ruleCount int64
	db.Model(&model.ApprovalRule{}).Count(&ruleCount)
	if ruleCount == 0 {
		db.Create(&model.ApprovalRule{MinAmount: 0, MaxAmount: 5000000, RequiredRole: "APPRAISER"})
		db.Create(&model.ApprovalRule{MinAmount: 5000001, MaxAmount: 20000000, RequiredRole: "MANAGER"})
		db.Create(&model.ApprovalRule{MinAmount: 20000001, MaxAmount: 999999999999, RequiredRole: "SUPER_ADMIN"})
	}

	log.Println("Database initialization completed: Clean state ready (0 customers, 0 pawns, 0 assets).")
}

