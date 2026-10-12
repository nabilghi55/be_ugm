package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"samplegadai-backend/internal/dto"
	"samplegadai-backend/internal/service"
	"samplegadai-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// AUTH HANDLERS
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid login request payload", err.Error())
		return
	}

	res, err := h.svc.Login(c.Request.Context(), req, c.ClientIP())
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", res)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Refresh token is required", err.Error())
		return
	}

	res, err := h.svc.RefreshTokenRotation(c.Request.Context(), req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", res)
}

// STORAGE PRESIGNED URL HANDLER
func (h *Handler) GetPresignedUpload(c *gin.Context) {
	var req dto.PresignedUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid upload request payload", err.Error())
		return
	}

	res, err := h.svc.GeneratePresignedUpload(c.Request.Context(), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate presigned upload URL", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Presigned upload URL generated successfully", res)
}

// CUSTOMER HANDLERS
func (h *Handler) CreateCustomer(c *gin.Context) {
	var req dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid customer payload", err.Error())
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	cust, err := h.svc.CreateCustomer(c.Request.Context(), req, userID, c.ClientIP())
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create customer", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Customer created successfully", cust)
}

func (h *Handler) GetCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid customer ID format", nil)
		return
	}

	cust, err := h.svc.GetCustomer(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Customer not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Customer fetched successfully", cust)
}

func (h *Handler) ListCustomers(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	customers, total, err := h.svc.ListCustomers(c.Request.Context(), search, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch customers", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Customers retrieved successfully", customers, page, limit, total)
}

func (h *Handler) AddCustomerDocument(c *gin.Context) {
	idStr := c.Param("id")
	customerID, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid customer ID format", nil)
		return
	}

	var req dto.AddCustomerDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid document payload", err.Error())
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	doc, err := h.svc.AddCustomerDocument(c.Request.Context(), customerID, req, userID, c.ClientIP())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save document metadata", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Customer document added successfully", doc)
}

func (h *Handler) VerifyKTPOCR(c *gin.Context) {
	docIDStr := c.Param("doc_id")
	docID, err := uuid.Parse(docIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid document ID format", nil)
		return
	}

	var req dto.VerifyOCRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid OCR verification payload", err.Error())
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	doc, err := h.svc.VerifyKTPOCR(c.Request.Context(), docID, req, userID, c.ClientIP())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to verify OCR document", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "KTP OCR verification updated", doc)
}

// ASSET HANDLERS
func (h *Handler) CreateAsset(c *gin.Context) {
	var req dto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid asset payload", err.Error())
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	asset, err := h.svc.CreateAsset(c.Request.Context(), req, userID, c.ClientIP())
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create asset", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Asset created successfully", asset)
}

func (h *Handler) GetAsset(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID format", nil)
		return
	}

	asset, err := h.svc.GetAsset(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Asset not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Asset retrieved successfully", asset)
}

func (h *Handler) ListAssets(c *gin.Context) {
	categoryID := c.Query("category_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	assets, total, err := h.svc.ListAssets(c.Request.Context(), categoryID, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch assets", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Assets retrieved successfully", assets, page, limit, total)
}

func (h *Handler) ListCategories(c *gin.Context) {
	categories, err := h.svc.ListCategories(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch categories", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Asset categories retrieved", categories)
}

// MARKET PRICE ENGINE & APPRAISAL HANDLERS
func (h *Handler) SearchMarketPrice(c *gin.Context) {
	var query dto.MarketPriceQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid market price search query", err.Error())
		return
	}

	result, err := h.svc.SearchMarketPrice(c.Request.Context(), query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to perform market price search", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Market price search successful", result)
}

func (h *Handler) CreateAppraisal(c *gin.Context) {
	var req dto.CreateAppraisalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid appraisal payload", err.Error())
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	appraisal, err := h.svc.CreateAppraisal(c.Request.Context(), req, userID, c.ClientIP())
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to record appraisal", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Appraisal completed successfully", appraisal)
}

// PAWN TRANSACTION HANDLERS
func (h *Handler) CreatePawnTransaction(c *gin.Context) {
	var req dto.CreatePawnTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Format pengajuan tidak valid: %v", err), err.Error())
		return
	}

	branchID := c.MustGet("branch_id").(uuid.UUID)
	userID := c.MustGet("user_id").(uuid.UUID)

	tx, err := h.svc.CreatePawnTransaction(c.Request.Context(), req, branchID, userID, c.ClientIP())
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Pawn transaction submitted for approval", tx)
}

func (h *Handler) ApprovePawnTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID format", nil)
		return
	}

	var req dto.ApprovalActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid approval payload", err.Error())
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	role := c.MustGet("role").(string)

	tx, err := h.svc.ApprovePawnTransaction(c.Request.Context(), id, req, userID, role, c.ClientIP())
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Pawn transaction approval updated", tx)
}

func (h *Handler) DisburseLoan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID format", nil)
		return
	}

	var req dto.DisbursementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid disbursement payload", err.Error())
		return
	}

	cashierID := c.MustGet("user_id").(uuid.UUID)
	tx, err := h.svc.DisburseLoan(c.Request.Context(), id, req, cashierID, c.ClientIP())
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Loan disbursed successfully and inventory updated", tx)
}

func (h *Handler) GetPawnTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID format", nil)
		return
	}

	pawn, err := h.svc.GetPawnTransaction(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Pawn transaction not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaction details retrieved", pawn)
}

func (h *Handler) ListPawnTransactions(c *gin.Context) {
	status := c.Query("status")
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	txs, total, err := h.svc.ListPawnTransactions(c.Request.Context(), status, search, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to list pawn transactions", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Pawn transactions retrieved", txs, page, limit, total)
}


// PAYMENT HANDLERS
func (h *Handler) CalculatePayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID format", nil)
		return
	}

	calc, err := h.svc.CalculatePayment(c.Request.Context(), id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment calculation complete", calc)
}

func (h *Handler) ProcessPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID format", nil)
		return
	}

	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid payment payload", err.Error())
		return
	}

	cashierID := c.MustGet("user_id").(uuid.UUID)
	payment, err := h.svc.ProcessPayment(c.Request.Context(), id, req, cashierID, c.ClientIP())
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Payment processed successfully", payment)
}

// ASSET RELEASE HANDLER
func (h *Handler) ReleaseAsset(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID format", nil)
		return
	}

	var req dto.AssetReleaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid asset release payload", err.Error())
		return
	}

	officerID := c.MustGet("user_id").(uuid.UUID)
	rel, err := h.svc.ReleaseAsset(c.Request.Context(), id, req, officerID, c.ClientIP())
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Asset released successfully to customer", rel)
}

// INVENTORY HANDLER
func (h *Handler) ListInventory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	items, total, err := h.svc.ListInventory(c.Request.Context(), page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to list inventory", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Inventory retrieved", items, page, limit, total)
}

func (h *Handler) GetInventoryByQR(c *gin.Context) {
	qr := c.Param("qr")
	item, err := h.svc.GetInventoryByQR(c.Request.Context(), qr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Inventory QR code not found", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Inventory details retrieved", item)
}

// DASHBOARD HANDLER
func (h *Handler) GetDashboardSummary(c *gin.Context) {
	summary, err := h.svc.GetDashboardSummary(c.Request.Context())
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to compute dashboard metrics", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Dashboard summary retrieved", summary)
}

// AUDIT LOG HANDLER
func (h *Handler) ListAuditLogs(c *gin.Context) {
	var query dto.AuditLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	logs, total, err := h.svc.ListAuditLogs(c.Request.Context(), query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch audit logs", err.Error())
		return
	}

	utils.PaginatedSuccessResponse(c, http.StatusOK, "Audit logs retrieved", logs, query.Page, query.Limit, total)
}
