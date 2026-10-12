package dto

import (
	"time"

	"github.com/google/uuid"
)

// Auth DTOs
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	User         UserDTO   `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type UserDTO struct {
	ID        uuid.UUID  `json:"id"`
	BranchID  *uuid.UUID `json:"branch_id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	FullName  string     `json:"full_name"`
	Role      string     `json:"role"`
	IsActive  bool       `json:"is_active"`
	LastLogin *time.Time `json:"last_login_at,omitempty"`
}

type CreateUserRequest struct {
	BranchID *uuid.UUID `json:"branch_id"`
	Username string     `json:"username" binding:"required"`
	Email    string     `json:"email" binding:"required,email"`
	Password string     `json:"password" binding:"required,min=6"`
	FullName string     `json:"full_name" binding:"required"`
	Role     string     `json:"role" binding:"required"`
}

type UpdateUserRequest struct {
	BranchID *uuid.UUID `json:"branch_id"`
	Email    string     `json:"email"`
	FullName string     `json:"full_name"`
	Role     string     `json:"role"`
	IsActive *bool      `json:"is_active"`
}

// Branch DTOs
type BranchDTO struct {
	ID       uuid.UUID `json:"id"`
	Code     string    `json:"code"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	Phone    string    `json:"phone"`
	IsActive bool      `json:"is_active"`
}

type CreateBranchRequest struct {
	Code    string `json:"code" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

// Storage Presigned URL DTOs
type PresignedUploadRequest struct {
	FileName    string `json:"file_name" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	Folder      string `json:"folder"`
}

type PresignedUploadResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
}

// Customer DTOs
type CreateCustomerRequest struct {
	NIK                          string `json:"nik" binding:"required"`
	FullName                     string `json:"full_name" binding:"required"`
	BirthPlace                   string `json:"birth_place"`
	BirthDate                    string `json:"birth_date"` // YYYY-MM-DD
	Gender                       string `json:"gender"`
	Address                      string `json:"address"`
	Province                     string `json:"province"`
	City                         string `json:"city"`
	District                     string `json:"district"`
	Village                      string `json:"village"`
	Phone                        string `json:"phone" binding:"required"`
	WhatsApp                     string `json:"whatsapp"`
	Email                        string `json:"email"`
	Occupation                   string `json:"occupation"`
	MotherMaidenName             string `json:"mother_maiden_name"`
	EmergencyContactName         string `json:"emergency_contact_name"`
	EmergencyContactPhone        string `json:"emergency_contact_phone"`
	EmergencyContactRelationship string `json:"emergency_contact_relationship"`
}

type AddCustomerDocumentRequest struct {
	DocumentType string `json:"document_type" binding:"required"`
	ObjectKey    string `json:"object_key" binding:"required"`
}

type VerifyOCRRequest struct {
	NIK        string `json:"nik"`
	FullName   string `json:"full_name"`
	BirthPlace string `json:"birth_place"`
	BirthDate  string `json:"birth_date"`
	Gender     string `json:"gender"`
	Address    string `json:"address"`
	Status     string `json:"status" binding:"required"` // VERIFIED / REJECTED
}

// Asset DTOs
type CreateAssetRequest struct {
	CategoryID          uuid.UUID                `json:"category_id" binding:"required"`
	Subcategory         string                   `json:"subcategory"`
	Brand               string                   `json:"brand" binding:"required"`
	Model               string                   `json:"model" binding:"required"`
	Variant             string                   `json:"variant"`
	SerialNumber        string                   `json:"serial_number"`
	IMEI                string                   `json:"imei"`
	Year                int                      `json:"year"`
	Color               string                   `json:"color"`
	Specification       string                   `json:"specification"`
	Accessories         string                   `json:"accessories"`
	PhysicalCondition   string                   `json:"physical_condition" binding:"required"`
	FunctionalCondition string                   `json:"functional_condition" binding:"required"`
	PurchasePrice       int64                    `json:"purchase_price"`
	Notes               string                   `json:"notes"`
	Inspections         []InspectionItemRequest  `json:"inspections"`
	Photos              []AssetPhotoInputRequest `json:"photos"`
}

type InspectionItemRequest struct {
	ItemName string `json:"item_name" binding:"required"`
	Status   string `json:"status" binding:"required"` // GOOD, FAIR, POOR, BROKEN
	Notes    string `json:"notes"`
}

type AssetPhotoInputRequest struct {
	ObjectKey string `json:"object_key" binding:"required"`
	ImageType string `json:"image_type" binding:"required"`
	IsPrimary bool   `json:"is_primary"`
}

// Market Price Engine DTOs
type MarketPriceQuery struct {
	CategoryID    uuid.UUID `json:"category_id" form:"category_id"`
	CategoryName  string    `json:"category_name" form:"category_name"`
	Brand         string    `json:"brand" form:"brand" binding:"required"`
	Model         string    `json:"model" form:"model" binding:"required"`
	Variant       string    `json:"variant" form:"variant"`
	Specification string    `json:"specification" form:"specification"`
	Year          int       `json:"year" form:"year"`
	Condition     string    `json:"condition" form:"condition"`
}

type MarketPriceResult struct {
	RecommendedPrice int64                      `json:"recommended_price"`
	MinimumPrice     int64                      `json:"minimum_price"`
	MaximumPrice     int64                      `json:"maximum_price"`
	MedianPrice      int64                      `json:"median_price"`
	TrimmedMeanPrice int64                      `json:"trimmed_mean_price"`
	ConfidenceScore  int                        `json:"confidence_score"`
	ComparableCount  int                        `json:"comparable_count"`
	HistoricalPrice  int64                      `json:"historical_price"`
	Comparables      []MarketComparableResponse `json:"comparables"`
}

type MarketComparableResponse struct {
	Source          string    `json:"source"`
	ExternalID      string    `json:"external_id"`
	Title           string    `json:"title"`
	URL             string    `json:"url"`
	Price           int64     `json:"price"`
	Location        string    `json:"location"`
	Condition       string    `json:"condition"`
	SellerType      string    `json:"seller_type"`
	Brand           string    `json:"brand"`
	Model           string    `json:"model"`
	Variant         string    `json:"variant"`
	Specification   string    `json:"specification"`
	FetchedAt       time.Time `json:"fetched_at"`
	SimilarityScore float64   `json:"similarity_score"`
}

// Appraisal DTOs
type CreateAppraisalRequest struct {
	AssetID                    uuid.UUID `json:"asset_id" binding:"required"`
	MarketPrice                int64     `json:"market_price" binding:"required"`
	MarketPriceMin             int64     `json:"market_price_min"`
	MarketPriceMax             int64     `json:"market_price_max"`
	MarketConfidence           int       `json:"market_confidence"`
	HistoricalPrice            int64     `json:"historical_price"`
	ConditionAdjustmentPercent float64   `json:"condition_adjustment_percent"`
	AppraiserPrice             int64     `json:"appraiser_price" binding:"required"`
	AppraiserNotes             string    `json:"appraiser_notes"`
	LoanToValue                float64   `json:"loan_to_value" binding:"required"`
}

// Pawn Transaction DTOs
type CreatePawnTransactionRequest struct {
	CustomerID      uuid.UUID `json:"customer_id" binding:"required"`
	AppraisalID     uuid.UUID `json:"appraisal_id"`
	PrincipalAmount int64     `json:"principal_amount" binding:"required"`
	InterestRate    float64   `json:"interest_rate"`
	AdminFee        int64     `json:"administration_fee"`
	StorageFee      int64     `json:"storage_fee"`
	TenorDays       int       `json:"tenor_days"`
}

type ApprovalActionRequest struct {
	Action string `json:"action" binding:"required"` // APPROVE / REJECT
	Notes  string `json:"notes"`
}

type DisbursementRequest struct {
	Method          string `json:"method" binding:"required"` // CASH, BANK_TRANSFER, EWALLET
	BankName        string `json:"bank_name"`
	AccountNumber   string `json:"account_number"`
	ReferenceNumber string `json:"reference_number"`
}

// Payment DTOs
type PaymentCalculationResponse struct {
	OutstandingPrincipal int64 `json:"outstanding_principal"`
	AccruedInterest      int64 `json:"accrued_interest"`
	AdminFee             int64 `json:"admin_fee"`
	StorageFee           int64 `json:"storage_fee"`
	PenaltyFee           int64 `json:"penalty_fee"`
	DaysOverdue          int   `json:"days_overdue"`
	TotalPayoff          int64 `json:"total_payoff"`
}

type CreatePaymentRequest struct {
	PaymentType     string `json:"payment_type" binding:"required"` // INTEREST_PAYMENT, PARTIAL_PAYMENT, EXTENSION_PAYMENT, FULL_PAYOFF
	Amount          int64  `json:"amount" binding:"required"`
	PaymentMethod   string `json:"payment_method" binding:"required"`
	ReferenceNumber string `json:"reference_number"`
}

// Asset Release DTOs
type AssetReleaseRequest struct {
	ReceiverName     string `json:"receiver_name" binding:"required"`
	ReceiverIDNumber string `json:"receiver_id_number"`
	ProofObjectKey   string `json:"proof_object_key"`
	Notes            string `json:"notes"`
}

// Dashboard DTOs
type DashboardSummaryResponse struct {
	ActiveLoansCount      int64 `json:"active_loans_count"`
	TotalLoanOutstanding  int64 `json:"total_loan_outstanding"`
	TotalAssetMarketValue int64 `json:"total_asset_market_value"`
	TodayDisbursement     int64 `json:"today_disbursement"`
	TodayPayment          int64 `json:"today_payment"`
	InterestIncomeTotal   int64 `json:"interest_income_total"`
	OverdueCount          int64 `json:"overdue_count"`
	TotalCustomers        int64 `json:"total_customers"`
	StoredAssetsCount     int64 `json:"stored_assets_count"`
	CompletedLoansCount   int64 `json:"completed_loans_count"`
}

// Audit Log Query DTO
type AuditLogQuery struct {
	Entity string `json:"entity" form:"entity"`
	Action string `json:"action" form:"action"`
	Page   int    `json:"page" form:"page,default=1"`
	Limit  int    `json:"limit" form:"limit,default=20"`
}
