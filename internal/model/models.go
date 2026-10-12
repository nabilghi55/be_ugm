package model

import (
	"time"

	"github.com/google/uuid"
)

type Branch struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Code      string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Address   string    `gorm:"type:text" json:"address"`
	Phone     string    `gorm:"type:varchar(30)" json:"phone"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	BranchID     *uuid.UUID `gorm:"type:uuid" json:"branch_id"`
	Branch       *Branch    `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Username     string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email        string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	FullName     string     `gorm:"type:varchar(100);not null" json:"full_name"`
	Role         string     `gorm:"type:varchar(30);not null" json:"role"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type RefreshToken struct {
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TokenHash        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"-"`
	ExpiresAt        time.Time `json:"expires_at"`
	IsRevoked        bool      `gorm:"default:false" json:"is_revoked"`
	ReplacedByToken  string    `gorm:"type:varchar(255)" json:"-"`
	CreatedAt        time.Time `json:"created_at"`
}

type PasswordReset struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Token     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	IsUsed    bool      `gorm:"default:false" json:"is_used"`
	CreatedAt time.Time `json:"created_at"`
}

type Customer struct {
	ID                            uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CustomerCode                  string             `gorm:"type:varchar(30);uniqueIndex;not null" json:"customer_code"`
	NIK                           string             `gorm:"type:varchar(20);uniqueIndex;not null" json:"nik"`
	FullName                      string             `gorm:"type:varchar(100);not null" json:"full_name"`
	BirthPlace                    string             `gorm:"type:varchar(50)" json:"birth_place"`
	BirthDate                     *time.Time         `json:"birth_date"`
	Gender                        string             `gorm:"type:varchar(10)" json:"gender"`
	Address                       string             `gorm:"type:text" json:"address"`
	Province                      string             `gorm:"type:varchar(50)" json:"province"`
	City                          string             `gorm:"type:varchar(50)" json:"city"`
	District                      string             `gorm:"type:varchar(50)" json:"district"`
	Village                       string             `gorm:"type:varchar(50)" json:"village"`
	Phone                         string             `gorm:"type:varchar(30);not null" json:"phone"`
	WhatsApp                      string             `gorm:"column:whatsapp;type:varchar(30)" json:"whatsapp"`
	Email                         string             `gorm:"type:varchar(100)" json:"email"`
	Occupation                    string             `gorm:"type:varchar(50)" json:"occupation"`
	MotherMaidenName              string             `gorm:"type:varchar(100)" json:"mother_maiden_name"`
	EmergencyContactName          string             `gorm:"type:varchar(100)" json:"emergency_contact_name"`
	EmergencyContactPhone         string             `gorm:"type:varchar(30)" json:"emergency_contact_phone"`
	EmergencyContactRelationship  string             `gorm:"type:varchar(50)" json:"emergency_contact_relationship"`
	Documents                     []CustomerDocument `gorm:"foreignKey:CustomerID" json:"documents,omitempty"`
	CreatedAt                     time.Time          `json:"created_at"`
	UpdatedAt                     time.Time          `json:"updated_at"`
}

type CustomerDocument struct {
	ID                 uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CustomerID         uuid.UUID  `gorm:"type:uuid;not null" json:"customer_id"`
	DocumentType       string     `gorm:"type:varchar(30);not null" json:"document_type"`
	ObjectKey          string     `gorm:"type:varchar(255);not null" json:"object_key"`
	ViewURL            string     `gorm:"-" json:"view_url,omitempty"`
	OCRData            string     `gorm:"type:jsonb" json:"ocr_data"`
	OCRConfidence      float64    `gorm:"type:numeric(5,2)" json:"ocr_confidence"`
	VerificationStatus string     `gorm:"type:varchar(20);default:'UNVERIFIED'" json:"verification_status"`
	VerifiedBy         *uuid.UUID `gorm:"type:uuid" json:"verified_by"`
	VerifiedAt         *time.Time `json:"verified_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type AssetCategory struct {
	ID         uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name       string              `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	LTVPercent float64             `gorm:"type:numeric(5,2);not null" json:"ltv_percent"`
	IsActive   bool                `gorm:"default:true" json:"is_active"`
	Checklists []CategoryChecklist `gorm:"foreignKey:CategoryID" json:"checklists,omitempty"`
	CreatedAt  time.Time           `json:"created_at"`
}

type CategoryChecklist struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CategoryID  uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	ItemName    string    `gorm:"type:varchar(100);not null" json:"item_name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Asset struct {
	ID                  uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	AssetCode           string             `gorm:"type:varchar(30);uniqueIndex;not null" json:"asset_code"`
	CategoryID          uuid.UUID          `gorm:"type:uuid;not null" json:"category_id"`
	Category            *AssetCategory     `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Subcategory         string             `gorm:"type:varchar(50)" json:"subcategory"`
	Brand               string             `gorm:"type:varchar(50);not null" json:"brand"`
	Model               string             `gorm:"type:varchar(100);not null" json:"model"`
	Variant             string             `gorm:"type:varchar(100)" json:"variant"`
	SerialNumber        string             `gorm:"type:varchar(100)" json:"serial_number"`
	IMEI                string             `gorm:"type:varchar(50)" json:"imei"`
	Year                int                `json:"year"`
	Color               string             `gorm:"type:varchar(30)" json:"color"`
	Specification       string             `gorm:"type:text" json:"specification"`
	Accessories         string             `gorm:"type:text" json:"accessories"`
	PhysicalCondition   string             `gorm:"type:varchar(20)" json:"physical_condition"`
	FunctionalCondition string             `gorm:"type:varchar(20)" json:"functional_condition"`
	PurchasePrice       int64              `json:"purchase_price"`
	Notes               string             `gorm:"type:text" json:"notes"`
	Photos              []AssetPhoto       `gorm:"foreignKey:AssetID" json:"photos,omitempty"`
	Inspections         []AssetInspection  `gorm:"foreignKey:AssetID" json:"inspections,omitempty"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
}

type AssetPhoto struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	AssetID   uuid.UUID  `gorm:"type:uuid;not null" json:"asset_id"`
	ObjectKey string     `gorm:"type:varchar(255);not null" json:"object_key"`
	ViewURL   string     `gorm:"-" json:"view_url,omitempty"`
	ImageType string     `gorm:"type:varchar(30);not null" json:"image_type"`
	IsPrimary bool       `gorm:"default:false" json:"is_primary"`
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
}

type AssetInspection struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	AssetID   uuid.UUID `gorm:"type:uuid;not null" json:"asset_id"`
	ItemName  string    `gorm:"type:varchar(100);not null" json:"item_name"`
	Status    string    `gorm:"type:varchar(20);not null" json:"status"`
	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

type MarketComparable struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Source          string    `gorm:"type:varchar(50);not null" json:"source"`
	ExternalID      string    `gorm:"type:varchar(100)" json:"external_id"`
	Title           string    `gorm:"type:varchar(255);not null" json:"title"`
	URL             string    `gorm:"type:text" json:"url"`
	Price           int64     `gorm:"not null" json:"price"`
	Location        string    `gorm:"type:varchar(100)" json:"location"`
	Condition       string    `gorm:"type:varchar(50)" json:"condition"`
	SellerType      string    `gorm:"type:varchar(50)" json:"seller_type"`
	Brand           string    `gorm:"type:varchar(50);not null" json:"brand"`
	Model           string    `gorm:"type:varchar(100);not null" json:"model"`
	Variant         string    `gorm:"type:varchar(100)" json:"variant"`
	Specification   string    `gorm:"type:text" json:"specification"`
	FetchedAt       time.Time `json:"fetched_at"`
	SimilarityScore float64   `gorm:"type:numeric(5,2)" json:"similarity_score"`
}

type Appraisal struct {
	ID                         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	AssetID                    uuid.UUID  `gorm:"type:uuid;not null" json:"asset_id"`
	Asset                      *Asset     `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	MarketPrice                int64      `gorm:"not null" json:"market_price"`
	MarketPriceMin             int64      `gorm:"not null" json:"market_price_min"`
	MarketPriceMax             int64      `gorm:"not null" json:"market_price_max"`
	MarketConfidence           int        `json:"market_confidence"`
	HistoricalPrice            int64      `json:"historical_price"`
	ConditionAdjustmentPercent float64    `gorm:"type:numeric(5,2)" json:"condition_adjustment_percent"`
	AdjustedMarketPrice        int64      `gorm:"not null" json:"adjusted_market_price"`
	AppraiserPrice             int64      `gorm:"not null" json:"appraiser_price"`
	AppraiserNotes             string     `gorm:"type:text" json:"appraiser_notes"`
	LoanToValue                float64    `gorm:"type:numeric(5,2);not null" json:"loan_to_value"`
	MaximumLoan                int64      `gorm:"not null" json:"maximum_loan"`
	RecommendedLoan            int64      `gorm:"not null" json:"recommended_loan"`
	AppraiserID                *uuid.UUID `gorm:"type:uuid" json:"appraiser_id"`
	Appraiser                  *User      `gorm:"foreignKey:AppraiserID" json:"appraiser,omitempty"`
	AppraisedAt                time.Time  `json:"appraised_at"`
}

type ApprovalRule struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	MinAmount    int64     `gorm:"not null" json:"min_amount"`
	MaxAmount    int64     `gorm:"not null" json:"max_amount"`
	RequiredRole string    `gorm:"type:varchar(30);not null" json:"required_role"`
	CreatedAt    time.Time `json:"created_at"`
}

type PawnTransaction struct {
	ID               uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PawnNumber       string        `gorm:"type:varchar(30);uniqueIndex;not null" json:"pawn_number"`
	CustomerID       uuid.UUID     `gorm:"type:uuid;not null" json:"customer_id"`
	Customer         *Customer     `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	BranchID         uuid.UUID     `gorm:"type:uuid;not null" json:"branch_id"`
	Branch           *Branch       `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	AppraisalID      uuid.UUID     `gorm:"type:uuid;not null" json:"appraisal_id"`
	Appraisal        *Appraisal    `gorm:"foreignKey:AppraisalID" json:"appraisal,omitempty"`
	PrincipalAmount  int64         `gorm:"not null" json:"principal_amount"`
	InterestRate     float64       `gorm:"type:numeric(5,2);not null" json:"interest_rate"`
	AdministrationFee int64        `gorm:"not null" json:"administration_fee"`
	StorageFee       int64         `gorm:"not null" json:"storage_fee"`
	PenaltyRate      float64       `gorm:"type:numeric(5,2);not null" json:"penalty_rate"`
	LoanDate         time.Time     `gorm:"type:date;not null" json:"loan_date"`
	DueDate          time.Time     `gorm:"type:date;not null" json:"due_date"`
	GracePeriod      int           `gorm:"default:7" json:"grace_period"`
	Status           string        `gorm:"type:varchar(30);not null;default:'DRAFT'" json:"status"`
	RequestedBy      *uuid.UUID    `gorm:"type:uuid" json:"requested_by"`
	ApprovedBy       *uuid.UUID    `gorm:"type:uuid" json:"approved_by"`
	ApprovedByUser   *User         `gorm:"foreignKey:ApprovedBy" json:"approved_by_user,omitempty"`
	ApprovedAt       *time.Time    `json:"approved_at"`
	ApprovalNotes    string        `gorm:"type:text" json:"approval_notes"`
	Disbursement     *Disbursement `gorm:"foreignKey:PawnTransactionID" json:"disbursement,omitempty"`
	Payments         []Payment     `gorm:"foreignKey:PawnTransactionID" json:"payments,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type Disbursement struct {
	ID                uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PawnTransactionID uuid.UUID  `gorm:"type:uuid;not null" json:"pawn_transaction_id"`
	Amount            int64      `gorm:"not null" json:"amount"`
	Method            string     `gorm:"type:varchar(30);not null" json:"method"`
	BankName          string     `gorm:"type:varchar(50)" json:"bank_name"`
	AccountNumber     string     `gorm:"type:varchar(50)" json:"account_number"`
	ReferenceNumber   string     `gorm:"type:varchar(100)" json:"reference_number"`
	ProcessedBy       *uuid.UUID `gorm:"type:uuid" json:"processed_by"`
	ProcessedAt       time.Time  `json:"processed_at"`
}

type Payment struct {
	ID                uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PawnTransactionID uuid.UUID  `gorm:"type:uuid;not null" json:"pawn_transaction_id"`
	PaymentNumber     string     `gorm:"type:varchar(30);uniqueIndex;not null" json:"payment_number"`
	PaymentType       string     `gorm:"type:varchar(30);not null" json:"payment_type"`
	PrincipalPaid     int64      `gorm:"default:0" json:"principal_paid"`
	InterestPaid      int64      `gorm:"default:0" json:"interest_paid"`
	AdminPaid         int64      `gorm:"default:0" json:"admin_paid"`
	StoragePaid       int64      `gorm:"default:0" json:"storage_paid"`
	PenaltyPaid       int64      `gorm:"default:0" json:"penalty_paid"`
	TotalAmount       int64      `gorm:"not null" json:"total_amount"`
	PaymentMethod     string     `gorm:"type:varchar(30);not null" json:"payment_method"`
	ReferenceNumber   string     `gorm:"type:varchar(100)" json:"reference_number"`
	ProcessedBy       *uuid.UUID `gorm:"type:uuid" json:"processed_by"`
	CreatedAt         time.Time  `json:"created_at"`
}

type PawnExtension struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PawnTransactionID uuid.UUID `gorm:"type:uuid;not null" json:"pawn_transaction_id"`
	PaymentID         uuid.UUID `gorm:"type:uuid;not null" json:"payment_id"`
	OldDueDate        time.Time `gorm:"type:date;not null" json:"old_due_date"`
	NewDueDate        time.Time `gorm:"type:date;not null" json:"new_due_date"`
	DaysExtended      int       `gorm:"default:30" json:"days_extended"`
	CreatedAt         time.Time `json:"created_at"`
}

type AssetRelease struct {
	ID                uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	PawnTransactionID uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"pawn_transaction_id"`
	ReceiverName      string     `gorm:"type:varchar(100);not null" json:"receiver_name"`
	ReceiverIDNumber  string     `gorm:"type:varchar(50)" json:"receiver_id_number"`
	ProofObjectKey    string     `gorm:"type:varchar(255)" json:"proof_object_key"`
	ProcessedBy       uuid.UUID  `gorm:"type:uuid;not null" json:"processed_by"`
	ReleasedAt        time.Time  `json:"released_at"`
	Notes             string     `gorm:"type:text" json:"notes"`
}

type Inventory struct {
	ID                uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	BranchID          uuid.UUID        `gorm:"type:uuid;not null" json:"branch_id"`
	Branch            *Branch          `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Warehouse         string           `gorm:"type:varchar(50);not null" json:"warehouse"`
	Rack              string           `gorm:"type:varchar(50);not null" json:"rack"`
	Box               string           `gorm:"type:varchar(50);not null" json:"box"`
	AssetID           uuid.UUID        `gorm:"type:uuid;uniqueIndex;not null" json:"asset_id"`
	Asset             *Asset           `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	PawnTransactionID *uuid.UUID       `gorm:"type:uuid" json:"pawn_transaction_id"`
	PawnTransaction   *PawnTransaction `gorm:"foreignKey:PawnTransactionID" json:"pawn_transaction,omitempty"`
	Status            string           `gorm:"type:varchar(30);default:'STORED'" json:"status"`
	QRCode            string           `gorm:"type:varchar(100);uniqueIndex;not null" json:"qr_code"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

type AuditLog struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID    *uuid.UUID `gorm:"type:uuid" json:"user_id"`
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action    string     `gorm:"type:varchar(50);not null" json:"action"`
	Entity    string     `gorm:"type:varchar(50);not null" json:"entity"`
	EntityID  string     `gorm:"type:varchar(100)" json:"entity_id"`
	OldValue  string     `gorm:"type:jsonb" json:"old_value"`
	NewValue  string     `gorm:"type:jsonb" json:"new_value"`
	IPAddress string     `gorm:"type:varchar(45)" json:"ip_address"`
	CreatedAt time.Time  `json:"created_at"`
}
