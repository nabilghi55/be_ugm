package repository

import (
	"context"
	"time"

	"samplegadai-backend/internal/dto"
	"samplegadai-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

// User Repository
func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Branch").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Branch").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Branch").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) UpdateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *Repository) ListUsers(ctx context.Context, page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	offset := (page - 1) * limit

	r.db.WithContext(ctx).Model(&model.User{}).Count(&total)
	err := r.db.WithContext(ctx).Preload("Branch").Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

// Refresh Token Repository
func (r *Repository) SaveRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *Repository) GetRefreshToken(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := r.db.WithContext(ctx).Where("token_hash = ? AND is_revoked = false AND expires_at > ?", tokenHash, time.Now()).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, tokenHash, replacedBy string) error {
	return r.db.WithContext(ctx).Model(&model.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Updates(map[string]interface{}{
			"is_revoked":        true,
			"replaced_by_token": replacedBy,
		}).Error
}

// Password Reset Repository
func (r *Repository) SavePasswordReset(ctx context.Context, pr *model.PasswordReset) error {
	return r.db.WithContext(ctx).Create(pr).Error
}

func (r *Repository) GetPasswordReset(ctx context.Context, token string) (*model.PasswordReset, error) {
	var pr model.PasswordReset
	err := r.db.WithContext(ctx).Where("token = ? AND is_used = false AND expires_at > ?", token, time.Now()).First(&pr).Error
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *Repository) MarkPasswordResetUsed(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Model(&model.PasswordReset{}).Where("token = ?", token).Update("is_used", true).Error
}

// Branch Repository
func (r *Repository) ListBranches(ctx context.Context) ([]model.Branch, error) {
	var branches []model.Branch
	err := r.db.WithContext(ctx).Order("code ASC").Find(&branches).Error
	return branches, err
}

func (r *Repository) CreateBranch(ctx context.Context, branch *model.Branch) error {
	return r.db.WithContext(ctx).Create(branch).Error
}

// Customer Repository
func (r *Repository) CreateCustomer(ctx context.Context, customer *model.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *Repository) GetCustomerByID(ctx context.Context, id uuid.UUID) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.WithContext(ctx).Preload("Documents").Where("id = ?", id).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *Repository) ListCustomers(ctx context.Context, search string, page, limit int) ([]model.Customer, int64, error) {
	var customers []model.Customer
	var total int64
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Model(&model.Customer{})
	if search != "" {
		likeSearch := "%" + search + "%"
		query = query.Where("full_name ILIKE ? OR nik ILIKE ? OR customer_code ILIKE ? OR phone ILIKE ?", likeSearch, likeSearch, likeSearch, likeSearch)
	}

	query.Count(&total)
	err := query.Preload("Documents").Offset(offset).Limit(limit).Order("created_at DESC").Find(&customers).Error
	return customers, total, err
}

func (r *Repository) AddCustomerDocument(ctx context.Context, doc *model.CustomerDocument) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *Repository) UpdateCustomerDocument(ctx context.Context, doc *model.CustomerDocument) error {
	return r.db.WithContext(ctx).Save(doc).Error
}

func (r *Repository) GetCustomerDocumentByID(ctx context.Context, id uuid.UUID) (*model.CustomerDocument, error) {
	var doc model.CustomerDocument
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&doc).Error
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

// Asset Repository
func (r *Repository) CreateAsset(ctx context.Context, asset *model.Asset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *Repository) GetAssetByID(ctx context.Context, id uuid.UUID) (*model.Asset, error) {
	var asset model.Asset
	err := r.db.WithContext(ctx).Preload("Category").Preload("Photos").Preload("Inspections").Where("id = ?", id).First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *Repository) ListAssets(ctx context.Context, categoryID string, page, limit int) ([]model.Asset, int64, error) {
	var assets []model.Asset
	var total int64
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Model(&model.Asset{})
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	query.Count(&total)
	err := query.Preload("Category").Preload("Photos").Offset(offset).Limit(limit).Order("created_at DESC").Find(&assets).Error
	return assets, total, err
}

func (r *Repository) ListAssetCategories(ctx context.Context) ([]model.AssetCategory, error) {
	var categories []model.AssetCategory
	err := r.db.WithContext(ctx).Preload("Checklists").Order("name ASC").Find(&categories).Error
	return categories, err
}

// Appraisal Repository
func (r *Repository) CreateAppraisal(ctx context.Context, appraisal *model.Appraisal) error {
	return r.db.WithContext(ctx).Create(appraisal).Error
}

func (r *Repository) GetAppraisalByID(ctx context.Context, id uuid.UUID) (*model.Appraisal, error) {
	var appraisal model.Appraisal
	err := r.db.WithContext(ctx).Preload("Asset").Preload("Asset.Category").Preload("Asset.Photos").Preload("Asset.Inspections").Preload("Appraiser").Where("id = ?", id).First(&appraisal).Error
	if err != nil {
		return nil, err
	}
	return &appraisal, nil
}

func (r *Repository) GetLatestAppraisal(ctx context.Context) (*model.Appraisal, error) {
	var appraisal model.Appraisal
	err := r.db.WithContext(ctx).Preload("Asset").Preload("Asset.Category").Order("appraised_at DESC").First(&appraisal).Error
	if err != nil {
		return nil, err
	}
	return &appraisal, nil
}

// Approval Rules Repository
func (r *Repository) GetApprovalRuleForAmount(ctx context.Context, amount int64) (*model.ApprovalRule, error) {
	var rule model.ApprovalRule
	err := r.db.WithContext(ctx).Where("min_amount <= ? AND max_amount >= ?", amount, amount).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// Pawn Transaction Repository
func (r *Repository) CreatePawnTransaction(ctx context.Context, tx *model.PawnTransaction) error {
	return r.db.WithContext(ctx).Create(tx).Error
}

func (r *Repository) UpdatePawnTransaction(ctx context.Context, tx *model.PawnTransaction) error {
	return r.db.WithContext(ctx).Save(tx).Error
}

func (r *Repository) GetPawnTransactionByID(ctx context.Context, id uuid.UUID) (*model.PawnTransaction, error) {
	var tx model.PawnTransaction
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Branch").
		Preload("Appraisal").
		Preload("Appraisal.Asset").
		Preload("Appraisal.Asset.Category").
		Preload("Appraisal.Asset.Photos").
		Preload("Disbursement").
		Preload("Payments").
		Where("id = ?", id).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *Repository) ListPawnTransactions(ctx context.Context, status, search string, page, limit int) ([]model.PawnTransaction, int64, error) {
	var txs []model.PawnTransaction
	var total int64
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Model(&model.PawnTransaction{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		likeSearch := "%" + search + "%"
		query = query.Where("pawn_number ILIKE ?", likeSearch)
	}

	query.Count(&total)
	err := query.Preload("Customer").Preload("Branch").Preload("Appraisal").Preload("Appraisal.Asset").Offset(offset).Limit(limit).Order("created_at DESC").Find(&txs).Error
	return txs, total, err
}

// Disbursement Repository
func (r *Repository) CreateDisbursement(ctx context.Context, d *model.Disbursement) error {
	return r.db.WithContext(ctx).Create(d).Error
}

// Payment Repository
func (r *Repository) CreatePayment(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// Extension Repository
func (r *Repository) CreatePawnExtension(ctx context.Context, ext *model.PawnExtension) error {
	return r.db.WithContext(ctx).Create(ext).Error
}

// Asset Release Repository
func (r *Repository) CreateAssetRelease(ctx context.Context, rel *model.AssetRelease) error {
	return r.db.WithContext(ctx).Create(rel).Error
}

// Inventory Repository
func (r *Repository) CreateInventory(ctx context.Context, inv *model.Inventory) error {
	return r.db.WithContext(ctx).Create(inv).Error
}

func (r *Repository) ListInventories(ctx context.Context, page, limit int) ([]model.Inventory, int64, error) {
	var items []model.Inventory
	var total int64
	offset := (page - 1) * limit

	r.db.WithContext(ctx).Model(&model.Inventory{}).Count(&total)
	err := r.db.WithContext(ctx).Preload("Branch").Preload("Asset").Preload("PawnTransaction").Offset(offset).Limit(limit).Order("created_at DESC").Find(&items).Error
	return items, total, err
}

func (r *Repository) GetInventoryByQRCode(ctx context.Context, qr string) (*model.Inventory, error) {
	var inv model.Inventory
	err := r.db.WithContext(ctx).Preload("Branch").Preload("Asset").Preload("PawnTransaction").Where("qr_code = ?", qr).First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// Audit Log Repository
func (r *Repository) CreateAuditLog(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *Repository) ListAuditLogs(ctx context.Context, query dto.AuditLogQuery) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64
	offset := (query.Page - 1) * query.Limit

	q := r.db.WithContext(ctx).Model(&model.AuditLog{})
	if query.Entity != "" {
		q = q.Where("entity = ?", query.Entity)
	}
	if query.Action != "" {
		q = q.Where("action = ?", query.Action)
	}

	q.Count(&total)
	err := q.Preload("User").Offset(offset).Limit(query.Limit).Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}
