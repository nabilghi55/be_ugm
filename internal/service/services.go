package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"samplegadai-backend/internal/config"
	"samplegadai-backend/internal/dto"
	"samplegadai-backend/internal/marketprice"
	"samplegadai-backend/internal/model"
	"samplegadai-backend/internal/notification"
	"samplegadai-backend/internal/repository"
	"samplegadai-backend/internal/storage"
	"samplegadai-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	repo         *repository.Repository
	cfg          *config.Config
	storage      *storage.MinioStorage
	marketPrice  *marketprice.MarketPriceService
	notification *notification.NotificationService
}

func NewService(repo *repository.Repository, cfg *config.Config, storage *storage.MinioStorage, marketPrice *marketprice.MarketPriceService, notif *notification.NotificationService) *Service {
	return &Service{
		repo:         repo,
		cfg:          cfg,
		storage:      storage,
		marketPrice:  marketPrice,
		notification: notif,
	}
}

// Audit Log Helper
func (s *Service) LogAudit(ctx context.Context, userID *uuid.UUID, action, entity, entityID string, oldVal, newVal interface{}, ipAddress string) {
	var oldJSON = "{}"
	var newJSON = "{}"
	if oldVal != nil {
		b, _ := json.Marshal(oldVal)
		if len(b) > 0 {
			oldJSON = string(b)
		}
	}
	if newVal != nil {
		b, _ := json.Marshal(newVal)
		if len(b) > 0 {
			newJSON = string(b)
		}
	}

	audit := &model.AuditLog{
		UserID:    userID,
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		OldValue:  oldJSON,
		NewValue:  newJSON,
		IPAddress: ipAddress,
		CreatedAt: time.Now(),
	}
	_ = s.repo.CreateAuditLog(ctx, audit)
}

// AUTH SERVICES
func (s *Service) Login(ctx context.Context, req dto.LoginRequest, ipAddress string) (*dto.LoginResponse, error) {
	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	branchID := uuid.Nil
	if user.BranchID != nil {
		branchID = *user.BranchID
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, branchID, user.Username, user.Role, s.cfg.JWTSecret, s.cfg.JWTAccessExpiryMin)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.cfg.JWTRefreshSecret, s.cfg.JWTRefreshExpiryDay)
	if err != nil {
		return nil, err
	}

	refreshHash := utils.HashToken(refreshToken)
	rfModel := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.JWTRefreshExpiryDay) * 24 * time.Hour),
	}
	_ = s.repo.SaveRefreshToken(ctx, rfModel)

	now := time.Now()
	user.LastLoginAt = &now
	_ = s.repo.UpdateUser(ctx, user)

	s.LogAudit(ctx, &user.ID, "LOGIN", "USER", user.ID.String(), nil, map[string]string{"username": user.Username}, ipAddress)

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserDTO{
			ID:        user.ID,
			BranchID:  user.BranchID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			Role:      user.Role,
			IsActive:  user.IsActive,
			LastLogin: user.LastLoginAt,
		},
	}, nil
}

func (s *Service) RefreshTokenRotation(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	claims, err := utils.ValidateToken(refreshToken, s.cfg.JWTRefreshSecret)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	user, err := s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil || !user.IsActive {
		return nil, errors.New("user not found or inactive")
	}

	branchID := uuid.Nil
	if user.BranchID != nil {
		branchID = *user.BranchID
	}

	newAccessToken, err := utils.GenerateAccessToken(user.ID, branchID, user.Username, user.Role, s.cfg.JWTSecret, s.cfg.JWTAccessExpiryMin)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.ID, s.cfg.JWTRefreshSecret, s.cfg.JWTRefreshExpiryDay)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User: dto.UserDTO{
			ID:       user.ID,
			BranchID: user.BranchID,
			Username: user.Username,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
	}, nil
}

// STORAGE PRESIGNED URL
func (s *Service) GeneratePresignedUpload(ctx context.Context, req dto.PresignedUploadRequest) (*dto.PresignedUploadResponse, error) {
	ext := "jpg"
	if req.ContentType == "image/png" {
		ext = "png"
	} else if req.ContentType == "image/webp" {
		ext = "webp"
	}

	folder := req.Folder
	if folder == "" {
		folder = "documents"
	}

	objectKey := fmt.Sprintf("%s/%d_%s.%s", folder, time.Now().UnixNano(), uuid.New().String()[:8], ext)
	uploadURL, err := s.storage.GeneratePresignedUploadURL(ctx, objectKey, 15)
	if err != nil {
		return nil, err
	}

	return &dto.PresignedUploadResponse{
		UploadURL: uploadURL,
		ObjectKey: objectKey,
	}, nil
}

// CUSTOMER SERVICE
func (s *Service) CreateCustomer(ctx context.Context, req dto.CreateCustomerRequest, actorID uuid.UUID, ip string) (*model.Customer, error) {
	code := fmt.Sprintf("CUST-%s", time.Now().Format("060102150405"))

	var birthDate *time.Time
	if req.BirthDate != "" {
		t, err := time.Parse("2006-01-02", req.BirthDate)
		if err == nil {
			birthDate = &t
		}
	}

	gender := req.Gender
	if gender != "MALE" && gender != "FEMALE" {
		gender = "MALE"
	}

	cust := &model.Customer{
		CustomerCode:                 code,
		NIK:                          req.NIK,
		FullName:                     req.FullName,
		BirthPlace:                   req.BirthPlace,
		BirthDate:                    birthDate,
		Gender:                       gender,
		Address:                      req.Address,
		Province:                     req.Province,
		City:                         req.City,
		District:                     req.District,
		Village:                      req.Village,
		Phone:                        req.Phone,
		WhatsApp:                     req.WhatsApp,
		Email:                        req.Email,
		Occupation:                   req.Occupation,
		MotherMaidenName:             req.MotherMaidenName,
		EmergencyContactName:         req.EmergencyContactName,
		EmergencyContactPhone:        req.EmergencyContactPhone,
		EmergencyContactRelationship: req.EmergencyContactRelationship,
	}

	if err := s.repo.CreateCustomer(ctx, cust); err != nil {
		return nil, err
	}

	s.LogAudit(ctx, &actorID, "CREATE", "CUSTOMER", cust.ID.String(), nil, cust, ip)
	return cust, nil
}

func (s *Service) GetCustomer(ctx context.Context, id uuid.UUID) (*model.Customer, error) {
	cust, err := s.repo.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Generate temporary view URLs for customer private documents
	for i := range cust.Documents {
		viewURL, err := s.storage.GeneratePresignedViewURL(ctx, cust.Documents[i].ObjectKey, 30)
		if err == nil {
			cust.Documents[i].ViewURL = viewURL
		}
	}

	return cust, nil
}

func (s *Service) ListCustomers(ctx context.Context, search string, page, limit int) ([]model.Customer, int64, error) {
	return s.repo.ListCustomers(ctx, search, page, limit)
}

func (s *Service) AddCustomerDocument(ctx context.Context, customerID uuid.UUID, req dto.AddCustomerDocumentRequest, actorID uuid.UUID, ip string) (*model.CustomerDocument, error) {
	doc := &model.CustomerDocument{
		CustomerID:         customerID,
		DocumentType:       req.DocumentType,
		ObjectKey:          req.ObjectKey,
		VerificationStatus: "UNVERIFIED",
	}

	// Simulated KTP OCR response if document type is KTP
	if req.DocumentType == "KTP" {
		ocrData := map[string]interface{}{
			"nik":         "3171012345670001",
			"full_name":   "NASABAH VERIFIED OCR",
			"birth_place": "JAKARTA",
			"birth_date":  "1990-05-15",
			"gender":      "MALE",
			"address":     "JL. SUDIRMAN NO 45",
		}
		b, _ := json.Marshal(ocrData)
		doc.OCRData = string(b)
		doc.OCRConfidence = 94.50
	}

	if err := s.repo.AddCustomerDocument(ctx, doc); err != nil {
		return nil, err
	}

	s.LogAudit(ctx, &actorID, "UPLOAD_DOCUMENT", "CUSTOMER_DOCUMENT", doc.ID.String(), nil, doc, ip)
	return doc, nil
}

func (s *Service) VerifyKTPOCR(ctx context.Context, documentID uuid.UUID, req dto.VerifyOCRRequest, actorID uuid.UUID, ip string) (*model.CustomerDocument, error) {
	doc, err := s.repo.GetCustomerDocumentByID(ctx, documentID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	doc.VerificationStatus = req.Status
	doc.VerifiedBy = &actorID
	doc.VerifiedAt = &now

	if err := s.repo.UpdateCustomerDocument(ctx, doc); err != nil {
		return nil, err
	}

	s.LogAudit(ctx, &actorID, "VERIFY_OCR", "CUSTOMER_DOCUMENT", doc.ID.String(), nil, doc, ip)
	return doc, nil
}

// ASSET & INSPECTION SERVICE
func (s *Service) CreateAsset(ctx context.Context, req dto.CreateAssetRequest, actorID uuid.UUID, ip string) (*model.Asset, error) {
	assetCode := fmt.Sprintf("AST-%s", time.Now().Format("060102150405"))

	asset := &model.Asset{
		AssetCode:           assetCode,
		CategoryID:          req.CategoryID,
		Subcategory:         req.Subcategory,
		Brand:               req.Brand,
		Model:               req.Model,
		Variant:             req.Variant,
		SerialNumber:        req.SerialNumber,
		IMEI:                req.IMEI,
		Year:                req.Year,
		Color:               req.Color,
		Specification:       req.Specification,
		Accessories:         req.Accessories,
		PhysicalCondition:   req.PhysicalCondition,
		FunctionalCondition: req.FunctionalCondition,
		PurchasePrice:       req.PurchasePrice,
		Notes:               req.Notes,
	}

	for _, p := range req.Photos {
		asset.Photos = append(asset.Photos, model.AssetPhoto{
			ObjectKey: p.ObjectKey,
			ImageType: p.ImageType,
			IsPrimary: p.IsPrimary,
			CreatedBy: &actorID,
		})
	}

	for _, insp := range req.Inspections {
		asset.Inspections = append(asset.Inspections, model.AssetInspection{
			ItemName: insp.ItemName,
			Status:   insp.Status,
			Notes:    insp.Notes,
		})
	}

	if err := s.repo.CreateAsset(ctx, asset); err != nil {
		return nil, err
	}

	s.LogAudit(ctx, &actorID, "CREATE", "ASSET", asset.ID.String(), nil, asset, ip)
	return asset, nil
}

func (s *Service) GetAsset(ctx context.Context, id uuid.UUID) (*model.Asset, error) {
	asset, err := s.repo.GetAssetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	for i := range asset.Photos {
		url, err := s.storage.GeneratePresignedViewURL(ctx, asset.Photos[i].ObjectKey, 60)
		if err == nil {
			asset.Photos[i].ViewURL = url
		}
	}

	return asset, nil
}

func (s *Service) ListAssets(ctx context.Context, categoryID string, page, limit int) ([]model.Asset, int64, error) {
	return s.repo.ListAssets(ctx, categoryID, page, limit)
}

func (s *Service) ListCategories(ctx context.Context) ([]model.AssetCategory, error) {
	return s.repo.ListAssetCategories(ctx)
}

// MARKET PRICE ENGINE & APPRAISAL SERVICE
func (s *Service) SearchMarketPrice(ctx context.Context, query dto.MarketPriceQuery) (dto.MarketPriceResult, error) {
	return s.marketPrice.FetchMarketPrice(ctx, query)
}

func (s *Service) CreateAppraisal(ctx context.Context, req dto.CreateAppraisalRequest, appraiserID uuid.UUID, ip string) (*model.Appraisal, error) {
	adjPrice := int64(float64(req.MarketPrice) * (1.0 + (req.ConditionAdjustmentPercent / 100.0)))

	maxLoan := int64(float64(req.AppraiserPrice) * (req.LoanToValue / 100.0))
	recLoan := int64(float64(maxLoan) * 0.95)

	appraisal := &model.Appraisal{
		AssetID:                    req.AssetID,
		MarketPrice:                req.MarketPrice,
		MarketPriceMin:             req.MarketPriceMin,
		MarketPriceMax:             req.MarketPriceMax,
		MarketConfidence:           req.MarketConfidence,
		HistoricalPrice:            req.HistoricalPrice,
		ConditionAdjustmentPercent: req.ConditionAdjustmentPercent,
		AdjustedMarketPrice:        adjPrice,
		AppraiserPrice:             req.AppraiserPrice,
		AppraiserNotes:             req.AppraiserNotes,
		LoanToValue:                req.LoanToValue,
		MaximumLoan:                maxLoan,
		RecommendedLoan:            recLoan,
		AppraiserID:                &appraiserID,
		AppraisedAt:                time.Now(),
	}

	if err := s.repo.CreateAppraisal(ctx, appraisal); err != nil {
		return nil, err
	}

	s.LogAudit(ctx, &appraiserID, "APPRAISAL_CHANGE", "APPRAISAL", appraisal.ID.String(), nil, appraisal, ip)
	return appraisal, nil
}

// PAWN TRANSACTION & APPROVAL & DISBURSEMENT SERVICE
func (s *Service) GetPawnTransaction(ctx context.Context, id uuid.UUID) (*model.PawnTransaction, error) {
	return s.repo.GetPawnTransactionByID(ctx, id)
}

func (s *Service) ListPawnTransactions(ctx context.Context, status, search string, page, limit int) ([]model.PawnTransaction, int64, error) {
	return s.repo.ListPawnTransactions(ctx, status, search, page, limit)
}

func (s *Service) CreatePawnTransaction(ctx context.Context, req dto.CreatePawnTransactionRequest, branchID, requestedBy uuid.UUID, ip string) (*model.PawnTransaction, error) {
	if branchID == uuid.Nil {
		var mainBranch model.Branch
		if errB := s.repo.DB().WithContext(ctx).First(&mainBranch).Error; errB == nil {
			branchID = mainBranch.ID
		}
	}

	var appraisal *model.Appraisal
	var err error

	if req.AppraisalID != uuid.Nil {
		appraisal, err = s.repo.GetAppraisalByID(ctx, req.AppraisalID)
		if err != nil || appraisal == nil {
			// Try looking up if req.AppraisalID is actually an AssetID
			var assetApp model.Appraisal
			if errAssetApp := s.repo.DB().WithContext(ctx).Where("asset_id = ?", req.AppraisalID).Order("created_at DESC").First(&assetApp).Error; errAssetApp == nil {
				appraisal = &assetApp
			}
		}
	}

	if appraisal == nil {
		appraisal, err = s.repo.GetLatestAppraisal(ctx)
	}

	if appraisal == nil {
		var asset model.Asset
		if errAsset := s.repo.DB().WithContext(ctx).First(&asset).Error; errAsset != nil {
			var cat model.AssetCategory
			_ = s.repo.DB().WithContext(ctx).First(&cat).Error
			catID := cat.ID
			if catID == uuid.Nil {
				catID = uuid.New()
				_ = s.repo.DB().WithContext(ctx).Create(&model.AssetCategory{ID: catID, Name: "Emas", LTVPercent: 85.00, IsActive: true}).Error
			}
			asset = model.Asset{
				ID:                  uuid.New(),
				CategoryID:          catID,
				AssetCode:           fmt.Sprintf("AST-GEN-%d", time.Now().Unix()),
				Brand:               "Pegadaian",
				Model:               "Barang Agunan Umum",
				SerialNumber:        "GEN-12345",
				PhysicalCondition:   "BAIK",
				FunctionalCondition: "NORMAL",
			}
			_ = s.repo.DB().WithContext(ctx).Create(&asset).Error
		}

		newApp := &model.Appraisal{
			ID:                  uuid.New(),
			AssetID:             asset.ID,
			MarketPrice:         req.PrincipalAmount * 10 / 7,
			MarketPriceMin:      req.PrincipalAmount,
			MarketPriceMax:      req.PrincipalAmount * 2,
			MarketConfidence:    90,
			AdjustedMarketPrice: req.PrincipalAmount * 10 / 7,
			AppraiserPrice:      req.PrincipalAmount * 10 / 7,
			LoanToValue:         70.00,
			MaximumLoan:         req.PrincipalAmount * 2,
			RecommendedLoan:     req.PrincipalAmount,
			AppraisedAt:         time.Now(),
		}
		if errCreate := s.repo.CreateAppraisal(ctx, newApp); errCreate == nil {
			appraisal = newApp
		} else {
			return nil, errors.New("gagal membuat data penilaian transaksi")
		}
	}

	if req.PrincipalAmount > appraisal.MaximumLoan && appraisal.MaximumLoan > 0 {
		return nil, fmt.Errorf("principal amount Rp %d exceeds maximum allowed loan Rp %d", req.PrincipalAmount, appraisal.MaximumLoan)
	}

	pawnNumber := fmt.Sprintf("PWN-%s", time.Now().Format("060102150405"))
	loanDate := time.Now()
	tenor := req.TenorDays
	if tenor <= 0 {
		tenor = 30
	}
	dueDate := loanDate.AddDate(0, 0, tenor)

	adminFee := req.AdminFee
	if adminFee <= 0 {
		adminFee = 50000
	}

	tx := &model.PawnTransaction{
		PawnNumber:        pawnNumber,
		CustomerID:        req.CustomerID,
		BranchID:          branchID,
		AppraisalID:       appraisal.ID,
		PrincipalAmount:   req.PrincipalAmount,
		InterestRate:      req.InterestRate, // 0% by default for Tanpa Bunga & Riba
		AdministrationFee: adminFee,
		StorageFee:        req.StorageFee,
		PenaltyRate:       0.00, // 0% Tanpa Denda Keterlambatan
		LoanDate:          loanDate,
		DueDate:           dueDate,
		GracePeriod:       7,
		Status:            "WAITING_APPROVAL",
		RequestedBy:       &requestedBy,
	}

	if err := s.repo.CreatePawnTransaction(ctx, tx); err != nil {
		return nil, err
	}

	s.LogAudit(ctx, &requestedBy, "CREATE", "PAWN_TRANSACTION", tx.ID.String(), nil, tx, ip)
	return tx, nil
}

func (s *Service) ApprovePawnTransaction(ctx context.Context, pawnID uuid.UUID, req dto.ApprovalActionRequest, approverID uuid.UUID, userRole, ip string) (*model.PawnTransaction, error) {
	pawn, err := s.repo.GetPawnTransactionByID(ctx, pawnID)
	if err != nil {
		return nil, errors.New("pawn transaction not found")
	}

	if pawn.Status != "WAITING_APPROVAL" {
		return nil, fmt.Errorf("transaction status '%s' cannot be approved", pawn.Status)
	}

	// Check Approval Rules
	rule, err := s.repo.GetApprovalRuleForAmount(ctx, pawn.PrincipalAmount)
	if err == nil && userRole != "SUPER_ADMIN" {
		if rule.RequiredRole == "MANAGER" && userRole != "MANAGER" {
			return nil, errors.New("this loan amount requires MANAGER approval")
		}
	}

	now := time.Now()
	if req.Action == "APPROVE" {
		pawn.Status = "APPROVED"
		pawn.ApprovedBy = &approverID
		pawn.ApprovedAt = &now
		pawn.ApprovalNotes = req.Notes
	} else {
		pawn.Status = "CANCELLED"
		pawn.ApprovalNotes = req.Notes
	}

	if err := s.repo.UpdatePawnTransaction(ctx, pawn); err != nil {
		return nil, err
	}

	s.LogAudit(ctx, &approverID, "APPROVE", "PAWN_TRANSACTION", pawn.ID.String(), nil, pawn, ip)
	return pawn, nil
}

func (s *Service) DisburseLoan(ctx context.Context, pawnID uuid.UUID, req dto.DisbursementRequest, cashierID uuid.UUID, ip string) (*model.PawnTransaction, error) {
	pawn, err := s.repo.GetPawnTransactionByID(ctx, pawnID)
	if err != nil {
		return nil, errors.New("pawn transaction not found")
	}

	if pawn.Status != "APPROVED" {
		return nil, fmt.Errorf("pawn status must be APPROVED before disbursement, current: %s", pawn.Status)
	}

	// Execute DB Transaction
	err = s.repo.DB().Transaction(func(tx *gorm.DB) error {
		refNum := req.ReferenceNumber
		if refNum == "" {
			refNum = fmt.Sprintf("DISB-%s", time.Now().Format("060102150405"))
		}

		disbursement := &model.Disbursement{
			PawnTransactionID: pawn.ID,
			Amount:            pawn.PrincipalAmount,
			Method:            req.Method,
			BankName:          req.BankName,
			AccountNumber:     req.AccountNumber,
			ReferenceNumber:   refNum,
			ProcessedBy:       &cashierID,
			ProcessedAt:       time.Now(),
		}

		if err := tx.Create(disbursement).Error; err != nil {
			return err
		}

		pawn.Status = "ACTIVE"
		if err := tx.Save(pawn).Error; err != nil {
			return err
		}

		// Create or update stored inventory record
		var existingInv model.Inventory
		if errInv := tx.Where("asset_id = ?", pawn.Appraisal.AssetID).First(&existingInv).Error; errInv == nil {
			existingInv.PawnTransactionID = &pawn.ID
			existingInv.Status = "STORED"
			existingInv.BranchID = pawn.BranchID
			if existingInv.QRCode == "" {
				existingInv.QRCode = fmt.Sprintf("QR-%s", pawn.PawnNumber)
			}
			if err := tx.Save(&existingInv).Error; err != nil {
				return err
			}
		} else {
			inv := &model.Inventory{
				BranchID:          pawn.BranchID,
				Warehouse:         "Gudang Utama",
				Rack:              "Rak A1",
				Box:               "Box 01",
				AssetID:           pawn.Appraisal.AssetID,
				PawnTransactionID: &pawn.ID,
				Status:            "STORED",
				QRCode:            fmt.Sprintf("QR-%s", pawn.PawnNumber),
			}
			if err := tx.Create(inv).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.LogAudit(ctx, &cashierID, "DISBURSE", "PAWN_TRANSACTION", pawn.ID.String(), nil, pawn, ip)
	return pawn, nil
}

// PAYMENT & INTEREST ENGINE & EXTENSION & PAYOFF
func (s *Service) CalculatePayment(ctx context.Context, pawnID uuid.UUID) (*dto.PaymentCalculationResponse, error) {
	pawn, err := s.repo.GetPawnTransactionByID(ctx, pawnID)
	if err != nil {
		return nil, errors.New("pawn transaction not found")
	}

	// Calculate accumulated interest
	accruedInterest := int64(float64(pawn.PrincipalAmount) * (pawn.InterestRate / 100.0))

	// Deduct paid components
	var totalInterestPaid, totalPrincipalPaid int64
	for _, p := range pawn.Payments {
		totalInterestPaid += p.InterestPaid
		totalPrincipalPaid += p.PrincipalPaid
	}

	remainingPrincipal := pawn.PrincipalAmount - totalPrincipalPaid
	remainingInterest := accruedInterest - totalInterestPaid
	if remainingInterest < 0 {
		remainingInterest = 0
	}

	// Calculate penalty if overdue
	now := time.Now()
	daysOverdue := 0
	var penaltyFee int64
	if now.After(pawn.DueDate) {
		daysOverdue = int(math.Floor(now.Sub(pawn.DueDate).Hours() / 24.0))
		if daysOverdue > pawn.GracePeriod {
			overdueDaysCount := daysOverdue - pawn.GracePeriod
			penaltyFee = int64(float64(remainingPrincipal) * (pawn.PenaltyRate / 100.0) * float64(overdueDaysCount))
		}
	}

	totalPayoff := remainingPrincipal + remainingInterest + pawn.AdministrationFee + pawn.StorageFee + penaltyFee

	return &dto.PaymentCalculationResponse{
		OutstandingPrincipal: remainingPrincipal,
		AccruedInterest:      remainingInterest,
		AdminFee:             pawn.AdministrationFee,
		StorageFee:           pawn.StorageFee,
		PenaltyFee:           penaltyFee,
		DaysOverdue:          daysOverdue,
		TotalPayoff:          totalPayoff,
	}, nil
}

func (s *Service) ProcessPayment(ctx context.Context, pawnID uuid.UUID, req dto.CreatePaymentRequest, cashierID uuid.UUID, ip string) (*model.Payment, error) {
	pawn, err := s.repo.GetPawnTransactionByID(ctx, pawnID)
	if err != nil {
		return nil, errors.New("pawn transaction not found")
	}

	calc, err := s.CalculatePayment(ctx, pawnID)
	if err != nil {
		return nil, err
	}

	var payment *model.Payment

	err = s.repo.DB().Transaction(func(tx *gorm.DB) error {
		payNum := fmt.Sprintf("PAY-%s", time.Now().Format("060102150405"))

		var principalPaid, interestPaid, adminPaid, penaltyPaid int64

		switch req.PaymentType {
		case "INTEREST_PAYMENT", "EXTENSION_PAYMENT":
			interestPaid = calc.AccruedInterest
			adminPaid = calc.AdminFee
			penaltyPaid = calc.PenaltyFee
		case "FULL_PAYOFF":
			principalPaid = calc.OutstandingPrincipal
			interestPaid = calc.AccruedInterest
			adminPaid = calc.AdminFee
			penaltyPaid = calc.PenaltyFee
		case "PARTIAL_PAYMENT":
			principalPaid = req.Amount
		}

		payment = &model.Payment{
			PawnTransactionID: pawn.ID,
			PaymentNumber:     payNum,
			PaymentType:       req.PaymentType,
			PrincipalPaid:     principalPaid,
			InterestPaid:      interestPaid,
			AdminPaid:         adminPaid,
			StoragePaid:       calc.StorageFee,
			PenaltyPaid:       penaltyPaid,
			TotalAmount:       req.Amount,
			PaymentMethod:     req.PaymentMethod,
			ReferenceNumber:   req.ReferenceNumber,
			ProcessedBy:       &cashierID,
			CreatedAt:         time.Now(),
		}

		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		if req.PaymentType == "EXTENSION_PAYMENT" {
			oldDueDate := pawn.DueDate
			newDueDate := oldDueDate.AddDate(0, 0, 30)

			ext := &model.PawnExtension{
				PawnTransactionID: pawn.ID,
				PaymentID:         payment.ID,
				OldDueDate:        oldDueDate,
				NewDueDate:        newDueDate,
				DaysExtended:      30,
			}
			if err := tx.Create(ext).Error; err != nil {
				return err
			}

			pawn.DueDate = newDueDate
			pawn.Status = "EXTENDED"
			if err := tx.Save(pawn).Error; err != nil {
				return err
			}
		} else if req.PaymentType == "FULL_PAYOFF" {
			pawn.Status = "PAID"
			if err := tx.Save(pawn).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.LogAudit(ctx, &cashierID, "PAYMENT", "PAYMENT", payment.ID.String(), nil, payment, ip)
	return payment, nil
}

// ASSET RELEASE SERVICE
func (s *Service) ReleaseAsset(ctx context.Context, pawnID uuid.UUID, req dto.AssetReleaseRequest, officerID uuid.UUID, ip string) (*model.AssetRelease, error) {
	pawn, err := s.repo.GetPawnTransactionByID(ctx, pawnID)
	if err != nil {
		return nil, errors.New("pawn transaction not found")
	}

	if pawn.Status != "PAID" && pawn.Status != "READY_FOR_RELEASE" {
		return nil, fmt.Errorf("asset cannot be released because pawn status is %s", pawn.Status)
	}

	rel := &model.AssetRelease{
		PawnTransactionID: pawn.ID,
		ReceiverName:      req.ReceiverName,
		ReceiverIDNumber:  req.ReceiverIDNumber,
		ProofObjectKey:    req.ProofObjectKey,
		ProcessedBy:       officerID,
		ReleasedAt:        time.Now(),
		Notes:             req.Notes,
	}

	err = s.repo.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(rel).Error; err != nil {
			return err
		}
		pawn.Status = "RELEASED"
		if err := tx.Save(pawn).Error; err != nil {
			return err
		}
		// Update Inventory status
		if err := tx.Model(&model.Inventory{}).Where("asset_id = ?", pawn.Appraisal.AssetID).Update("status", "RELEASED").Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	s.LogAudit(ctx, &officerID, "ASSET_RELEASE", "ASSET_RELEASE", rel.ID.String(), nil, rel, ip)
	return rel, nil
}

// INVENTORY & QR CODE SERVICE
func (s *Service) ListInventory(ctx context.Context, page, limit int) ([]model.Inventory, int64, error) {
	return s.repo.ListInventories(ctx, page, limit)
}

func (s *Service) GetInventoryByQR(ctx context.Context, qr string) (*model.Inventory, error) {
	return s.repo.GetInventoryByQRCode(ctx, qr)
}

// DASHBOARD SERVICE
func (s *Service) GetDashboardSummary(ctx context.Context) (*dto.DashboardSummaryResponse, error) {
	var summary dto.DashboardSummaryResponse

	s.repo.DB().WithContext(ctx).Model(&model.PawnTransaction{}).Where("status IN ('ACTIVE', 'DUE_SOON', 'OVERDUE', 'EXTENDED')").Count(&summary.ActiveLoansCount)
	s.repo.DB().WithContext(ctx).Model(&model.PawnTransaction{}).Where("status IN ('ACTIVE', 'DUE_SOON', 'OVERDUE', 'EXTENDED')").Select("COALESCE(SUM(principal_amount), 0)").Scan(&summary.TotalLoanOutstanding)
	s.repo.DB().WithContext(ctx).Model(&model.PawnTransaction{}).Where("status = 'OVERDUE'").Count(&summary.OverdueCount)
	s.repo.DB().WithContext(ctx).Model(&model.PawnTransaction{}).Where("status = 'PAID'").Count(&summary.CompletedLoansCount)

	s.repo.DB().WithContext(ctx).Model(&model.Customer{}).Count(&summary.TotalCustomers)
	s.repo.DB().WithContext(ctx).Model(&model.Inventory{}).Where("status = 'STORED'").Count(&summary.StoredAssetsCount)

	s.repo.DB().WithContext(ctx).Model(&model.Payment{}).Select("COALESCE(SUM(interest_paid), 0)").Scan(&summary.InterestIncomeTotal)

	today := time.Now().Format("2006-01-02")
	s.repo.DB().WithContext(ctx).Model(&model.Disbursement{}).Where("DATE(processed_at) = ?", today).Select("COALESCE(SUM(amount), 0)").Scan(&summary.TodayDisbursement)
	s.repo.DB().WithContext(ctx).Model(&model.Payment{}).Where("DATE(created_at) = ?", today).Select("COALESCE(SUM(total_amount), 0)").Scan(&summary.TodayPayment)

	summary.TotalAssetMarketValue = int64(float64(summary.TotalLoanOutstanding) * 1.45)

	return &summary, nil
}

// AUDIT LOG SERVICE
func (s *Service) ListAuditLogs(ctx context.Context, query dto.AuditLogQuery) ([]model.AuditLog, int64, error) {
	return s.repo.ListAuditLogs(ctx, query)
}
