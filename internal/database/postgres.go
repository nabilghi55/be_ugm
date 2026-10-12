package database

import (
	"fmt"
	"log"
	"math"
	"time"

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

	// Seed Initial Sample Customer & Asset & Appraisal & Pawn Transaction
	var custCount int64
	db.Model(&model.Customer{}).Count(&custCount)
	if custCount == 0 {
		custID := uuid.MustParse("00000000-0000-0000-0000-000000000100")
		db.Create(&model.Customer{
			ID:           custID,
			CustomerCode: "CUST-260917001",
			NIK:          "3171012345670001",
			FullName:     "Ahmad Subagyo",
			Gender:       "MALE",
			Phone:        "081234567890",
			WhatsApp:     "081234567890",
			Address:      "Jl. Sudirman No. 45",
			Province:     "Sumatera Barat",
			City:         "Kota Padang",
			District:     "Padang Barat",
			Occupation:   "Karyawan Swasta",
		})

		catID1 := uuid.MustParse("00000000-0000-0000-0000-000000000003")
		assetID := uuid.MustParse("00000000-0000-0000-0000-000000000200")
		db.Create(&model.Asset{
			ID:                  assetID,
			AssetCode:           "AST-260917001",
			CategoryID:          catID1,
			Brand:               "Apple",
			Model:               "iPhone 15 Pro",
			Variant:             "256GB",
			IMEI:                "354812345678901",
			Year:                2024,
			PhysicalCondition:   "GOOD",
			FunctionalCondition: "GOOD",
		})

		appID := uuid.MustParse("00000000-0000-0000-0000-000000000007")
		userID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
		db.Create(&model.Appraisal{
			ID:                         appID,
			AssetID:                    assetID,
			MarketPrice:                14700000,
			MarketPriceMin:             14200000,
			MarketPriceMax:             15300000,
			MarketConfidence:           87,
			ConditionAdjustmentPercent: -5.00,
			AdjustedMarketPrice:        13965000,
			AppraiserPrice:             14000000,
			LoanToValue:                70.00,
			MaximumLoan:                9800000,
			RecommendedLoan:            9310000,
			AppraiserID:                &userID,
		})

		pawnID := uuid.MustParse("00000000-0000-0000-0000-000000000300")
		db.Create(&model.PawnTransaction{
			ID:                pawnID,
			PawnNumber:        "PWN-260917-001",
			CustomerID:        custID,
			BranchID:          branchID1,
			AppraisalID:       appID,
			PrincipalAmount:   7000000,
			InterestRate:      2.00,
			AdministrationFee: 50000,
			StorageFee:        0,
			PenaltyRate:       0.50,
			LoanDate:          time.Now(),
			DueDate:           time.Now().AddDate(0, 0, 30),
			Status:            "ACTIVE",
		})

		db.Create(&model.Inventory{
			BranchID:          branchID1,
			Warehouse:         "Gudang Utama",
			Rack:              "RAK-A1",
			Box:               "BOX-01",
			AssetID:           assetID,
			PawnTransactionID: &pawnID,
			Status:            "STORED",
			QRCode:            "QR-PWN-260917-001",
		})
	}

	// Seed 20 Kota Padang Customers
	customersToSeed := []struct {
		CustomerCode string
		NIK          string
		FullName     string
		Gender       string
		Phone        string
		WhatsApp     string
		Address      string
		City         string
		District     string
		Occupation   string
	}{
		{"CUST-PADANG-001", "1371011505850001", "Ahmad Subagyo", "MALE", "081266123401", "081266123401", "Jl. Prof. Dr. Hamka No. 12", "Kota Padang", "Koto Tangah", "Pedagang Emas & Kelontong"},
		{"CUST-PADANG-002", "1371025208900002", "Siti Rahmah Bagindo", "FEMALE", "081374567802", "081374567802", "Jl. Surau Gadang No. 45", "Kota Padang", "Nanggalo", "Pegawai Negeri Sipil (PNS)"},
		{"CUST-PADANG-003", "1371031011880003", "Budi Kurniawan Caniago", "MALE", "082169876503", "082169876503", "Jl. Pasar Raya I No. 8", "Kota Padang", "Padang Barat", "Pemilik Rumah Makan Padang"},
		{"CUST-PADANG-004", "1371044403920004", "Rina Marlina Chaniago", "FEMALE", "085274123904", "085274123904", "Jl. Azizi No. 19", "Kota Padang", "Kuranji", "Guru SMA Negeri Padang"},
		{"CUST-PADANG-005", "1371052107840005", "Dedi Putra Jambak", "MALE", "081267890105", "081267890105", "Jl. Limau Manis No. 88", "Kota Padang", "Pauh", "Dosen Universitas Andalas"},
		{"CUST-PADANG-006", "1371066512950006", "Fitriani Sikumbang", "FEMALE", "081987654306", "081987654306", "Jl. Tanjung Saba No. 22", "Kota Padang", "Lubuk Begalung", "Wiraswasta Konveksi Textile"},
		{"CUST-PADANG-007", "1371071409820007", "Hendri Tanjung", "MALE", "085363456707", "085363456707", "Jl. Raya Bungus No. 102", "Kota Padang", "Bungus Teluk Kabung", "Pengusaha Perikanan & Logistik"},
		{"CUST-PADANG-008", "1371085004910008", "Dewi Kartika Piliang", "FEMALE", "081275987608", "081275987608", "Jl. Sawahan No. 34", "Kota Padang", "Padang Timur", "Apoteker & Owner Farmasi"},
		{"CUST-PADANG-009", "1371090806870009", "Rahmat Hidayat Kutianyir", "MALE", "081363112209", "081363112209", "Jl. Air Tawar Barat No. 5", "Kota Padang", "Padang Utara", "Kontraktor Bangunan"},
		{"CUST-PADANG-010", "1371104910940010", "Nurhaliza Baheramsyah", "FEMALE", "082283445510", "082283445510", "Jl. Indarung No. 77", "Kota Padang", "Lubuk Kilangan", "Karyawan PT Semen Padang"},
		{"CUST-PADANG-011", "1371010302860011", "Eko Prasetyo Supriyadi", "MALE", "081261334411", "081261334411", "Jl. Lubuk Buaya No. 15", "Kota Padang", "Koto Tangah", "Manager Operasional Ritel"},
		{"CUST-PADANG-012", "1371026107930012", "Maya Anggraini Melayu", "FEMALE", "085265778812", "085265778812", "Jl. Siteba No. 29", "Kota Padang", "Nanggalo", "Desainer Pakaian Khas Minang"},
		{"CUST-PADANG-013", "1371031204890013", "Fadli Rabbani Malay", "MALE", "081371990013", "081371990013", "Jl. Olo Ladang No. 6", "Kota Padang", "Padang Barat", "Teknisi & Owner Service HP"},
		{"CUST-PADANG-014", "1371045508960014", "Yulia Sri Wahyuni", "FEMALE", "082388112214", "082388112214", "Jl. Ampang No. 18", "Kota Padang", "Kuranji", "Perawat RSUP Dr. M. Djamil"},
		{"CUST-PADANG-015", "1371051801830015", "Doni Alamsyah Pili", "MALE", "085376334415", "085376334415", "Jl. Pasar Baru No. 40", "Kota Padang", "Pauh", "Pengusaha Rental Mobil Padang"},
		{"CUST-PADANG-016", "1371064211970016", "Nia Ramadhani Supeno", "FEMALE", "081267556616", "081267556616", "Jl. Cua No. 11", "Kota Padang", "Lubuk Begalung", "Barista & Cafe Owner"},
		{"CUST-PADANG-017", "1371072905880017", "Arief Rahman Hakim", "MALE", "081364778817", "081364778817", "Jl. Seberang Palinggam No. 3", "Kota Padang", "Padang Selatan", "Arsitek & Interior Designer"},
		{"CUST-PADANG-018", "1371086309950018", "Tari Mulia Sagita", "FEMALE", "082170990018", "082170990018", "Jl. Marapalam No. 25", "Kota Padang", "Padang Timur", "Manager Bank Daerah Nagari"},
		{"CUST-PADANG-019", "1371091112810019", "Zulfikar Bagindo Sutan", "MALE", "085263113319", "085263113319", "Jl. Ulak Karang No. 90", "Kota Padang", "Padang Utara", "Tokoh Masyarakat & Pedagang Grosir"},
		{"CUST-PADANG-020", "1371105404980020", "Dina Larasati Gunawan", "FEMALE", "081268224420", "081268224420", "Jl. Batu Gadang No. 14", "Kota Padang", "Lubuk Kilangan", "Content Creator & Florist"},
	}

	for _, cData := range customersToSeed {
		var existing model.Customer
		if err := db.Where("nik = ?", cData.NIK).First(&existing).Error; err != nil {
			db.Create(&model.Customer{
				CustomerCode: cData.CustomerCode,
				NIK:          cData.NIK,
				FullName:     cData.FullName,
				Gender:       cData.Gender,
				Phone:        cData.Phone,
				WhatsApp:     cData.WhatsApp,
				Address:      cData.Address,
				Province:     "Sumatera Barat",
				City:         cData.City,
				District:     cData.District,
				Occupation:   cData.Occupation,
			})
		}
	}

	// Seed 18 Pawn Transactions & Assets
	var assetCount int64
	db.Model(&model.Asset{}).Count(&assetCount)
	if assetCount < 5 {
		var allCusts []model.Customer
		db.Order("created_at ASC").Find(&allCusts)

		var catEmas, catHP, catLaptop, catKamera model.AssetCategory
		db.Where("name = ?", "Emas").First(&catEmas)
		db.Where("name = ?", "Smartphone").First(&catHP)
		db.Where("name = ?", "Laptop").First(&catLaptop)
		db.Where("name = ?", "Kamera").First(&catKamera)

		itemsToSeed := []struct {
			Category    model.AssetCategory
			Brand       string
			Model       string
			Variant     string
			Serial      string
			MarketPrice float64
			LoanAmount  float64
			Status      string
			BranchID    uuid.UUID
			Warehouse   string
			Rack        string
			Box         string
		}{
			{catEmas, "Antam", "Kalung Emas 24K", "10 Gram", "SN-ANTAM-24K-10G", 13500000, 10000000, "ACTIVE", branchID1, "Gudang Utama Nanggalo", "RAK-EMAS-1", "BOX-01"},
			{catEmas, "UBS", "Cincin Emas 22K", "5.5 Gram", "SN-UBS-22K-55G", 6200000, 4800000, "ACTIVE", branchID2, "Gudang Pasar Raya", "RAK-EMAS-2", "BOX-05"},
			{catEmas, "Antam", "Gelang Emas 24K", "15 Gram", "SN-ANTAM-24K-15G", 20100000, 15000000, "ACTIVE", branchID3, "Gudang Azizi", "RAK-EMAS-1", "BOX-02"},
			{catEmas, "Antam", "Logam Mulia CertiCard", "50 Gram", "SN-ANTAM-50G-999", 67500000, 50000000, "ACTIVE", branchID1, "Brankas Utama Pusat", "RAK-BRANKAS-1", "BOX-SPEC"},
			{catHP, "Apple", "iPhone 14 Pro Max", "256GB Deep Purple", "354812099881001", 16500000, 11000000, "ACTIVE", branchID2, "Gudang Gadget Raya", "RAK-HP-1", "BOX-12"},
			{catHP, "Samsung", "Galaxy S24 Ultra 5G", "512GB Titanium Black", "354988771122334", 18900000, 12500000, "ACTIVE", branchID3, "Gudang Gadget Azizi", "RAK-HP-2", "BOX-08"},
			{catHP, "Apple", "iPhone 13 Pro", "128GB Sierra Blue", "354112233445566", 11200000, 7500000, "OVERDUE", branchID1, "Gudang Utama Nanggalo", "RAK-HP-3", "BOX-15"},
			{catHP, "Xiaomi", "Xiaomi 14 Ultra", "512GB Black Edition", "864112233445577", 13800000, 8800000, "ACTIVE", branchID2, "Gudang Pasar Raya", "RAK-HP-1", "BOX-04"},
			{catLaptop, "Apple", "MacBook Air M2", "8GB / 256GB Midnight", "C02G11223344", 14500000, 9500000, "ACTIVE", branchID3, "Gudang Laptop Azizi", "RAK-LAPTOP-1", "BOX-21"},
			{catLaptop, "Asus", "ROG Zephyrus G14", "RTX 4060 16GB", "SN-ROG-G14-2023", 22000000, 14000000, "EXTENDED", branchID1, "Gudang Utama Nanggalo", "RAK-LAPTOP-2", "BOX-09"},
			{catLaptop, "Lenovo", "ThinkPad X1 Carbon Gen 10", "i7 16GB 512GB", "SN-THINKPAD-X1-G10", 17500000, 11500000, "ACTIVE", branchID2, "Gudang Pasar Raya", "RAK-LAPTOP-1", "BOX-11"},
			{catLaptop, "Dell", "XPS 13 9315", "Intel i7 16GB SSD", "SN-DELL-XPS13-9315", 15000000, 9800000, "OVERDUE", branchID3, "Gudang Laptop Azizi", "RAK-LAPTOP-3", "BOX-03"},
			{catKamera, "Sony", "Alpha A7 IV Mirrorless", "Body Only Black", "SN-SONY-A7M4-998", 31000000, 20000000, "ACTIVE", branchID1, "Gudang Kamera Nanggalo", "RAK-OPTIC-1", "BOX-07"},
			{catKamera, "Canon", "EOS R6 Mark II", "Kit 24-105mm IS STM", "SN-CANON-R6M2-771", 38000000, 25000000, "ACTIVE", branchID2, "Gudang Kamera Raya", "RAK-OPTIC-2", "BOX-02"},
			{catKamera, "Fujifilm", "X-T5 Mirrorless", "Body Only Silver", "SN-FUJI-XT5-554", 24000000, 15500000, "EXTENDED", branchID3, "Gudang Optik Azizi", "RAK-OPTIC-1", "BOX-14"},
			{catHP, "Apple", "iPad Pro 11-inch M2", "256GB WiFi Space Gray", "DMPG88776655", 13500000, 8800000, "OVERDUE", branchID1, "Gudang Utama Nanggalo", "RAK-TABLET-1", "BOX-06"},
		}

		for idx, item := range itemsToSeed {
			if item.Category.ID == uuid.Nil {
				continue
			}
			cust := allCusts[idx%len(allCusts)]

			astID := uuid.New()
			astCode := fmt.Sprintf("AST-PADANG-%03d", idx+1)
			db.Create(&model.Asset{
				ID:                  astID,
				AssetCode:           astCode,
				CategoryID:          item.Category.ID,
				Brand:               item.Brand,
				Model:               item.Model,
				Variant:             item.Variant,
				IMEI:                item.Serial,
				SerialNumber:        item.Serial,
				Year:                2023 + (idx % 2),
				PhysicalCondition:   "GOOD",
				FunctionalCondition: "GOOD",
			})

			appID := uuid.New()
			recLoan := int64(math.Round(item.MarketPrice * 0.70))
			db.Create(&model.Appraisal{
				ID:                         appID,
				AssetID:                    astID,
				MarketPrice:                int64(item.MarketPrice),
				MarketPriceMin:             int64(item.MarketPrice * 0.95),
				MarketPriceMax:             int64(item.MarketPrice * 1.05),
				MarketConfidence:           90,
				ConditionAdjustmentPercent: 0,
				AdjustedMarketPrice:        int64(item.MarketPrice),
				AppraiserPrice:             int64(item.MarketPrice),
				LoanToValue:                70.00,
				MaximumLoan:                int64(item.MarketPrice * 0.70),
				RecommendedLoan:            recLoan,
			})

			pawnID := uuid.New()
			pawnNo := fmt.Sprintf("PWN-260917-%03d", idx+2)
			loanDate := time.Now().AddDate(0, 0, -(idx * 5))
			dueDate := loanDate.AddDate(0, 0, 30)

			db.Create(&model.PawnTransaction{
				ID:                pawnID,
				PawnNumber:        pawnNo,
				CustomerID:        cust.ID,
				BranchID:          item.BranchID,
				AppraisalID:       appID,
				PrincipalAmount:   int64(item.LoanAmount),
				InterestRate:      2.00,
				AdministrationFee: 50000,
				StorageFee:        0,
				PenaltyRate:       0.50,
				LoanDate:          loanDate,
				DueDate:           dueDate,
				Status:            item.Status,
			})

			db.Create(&model.Inventory{
				BranchID:          item.BranchID,
				Warehouse:         item.Warehouse,
				Rack:              item.Rack,
				Box:               item.Box,
				AssetID:           astID,
				PawnTransactionID: &pawnID,
				Status:            "STORED",
				QRCode:            fmt.Sprintf("QR-%s", pawnNo),
			})
		}
	}
}
