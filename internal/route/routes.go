package route

import (
	"samplegadai-backend/internal/config"
	"samplegadai-backend/internal/handler"
	"samplegadai-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, h *handler.Handler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS())

	// Health check endpoint for Docker & Container monitoring
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "samplegadai-backend",
			"version": "1.0.0",
		})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "healthy", "service": "samplegadai-backend"})
		})

		// Public Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.RefreshToken)
		}

		// Authenticated routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Storage Upload URL
			protected.POST("/storage/presigned-upload", h.GetPresignedUpload)

			// Customers
			cust := protected.Group("/customers")
			{
				cust.GET("", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "APPRAISER", "CASHIER", "MANAGER"), h.ListCustomers)
				cust.POST("", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "MANAGER"), h.CreateCustomer)
				cust.GET("/:id", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "APPRAISER", "CASHIER", "MANAGER"), h.GetCustomer)
				cust.POST("/:id/documents", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "MANAGER"), h.AddCustomerDocument)
				cust.POST("/documents/:doc_id/verify-ocr", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "MANAGER"), h.VerifyKTPOCR)
			}

			// Assets & Inspection Checklist
			assets := protected.Group("/assets")
			{
				assets.GET("", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "APPRAISER", "CASHIER", "MANAGER"), h.ListAssets)
				assets.POST("", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "APPRAISER", "MANAGER"), h.CreateAsset)
				assets.GET("/categories", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "APPRAISER", "CASHIER", "MANAGER"), h.ListCategories)
				assets.GET("/:id", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "APPRAISER", "CASHIER", "MANAGER"), h.GetAsset)
			}

			// Market Price Intelligence Engine & Appraisal
			protected.GET("/market-price/search", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "APPRAISER", "MANAGER"), h.SearchMarketPrice)
			protected.POST("/appraisals", middleware.RBACMiddleware("SUPER_ADMIN", "APPRAISER", "MANAGER"), h.CreateAppraisal)

			// Pawn Transactions & Approvals
			pawns := protected.Group("/pawn-transactions")
			{
				pawns.GET("", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "APPRAISER", "CASHIER", "MANAGER"), h.ListPawnTransactions)
				pawns.POST("", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "APPRAISER", "MANAGER"), h.CreatePawnTransaction)
				pawns.GET("/:id", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "APPRAISER", "CASHIER", "MANAGER"), h.GetPawnTransaction)
				pawns.POST("/:id/approve", middleware.RBACMiddleware("SUPER_ADMIN", "APPRAISER", "MANAGER"), h.ApprovePawnTransaction)
				pawns.POST("/:id/disburse", middleware.RBACMiddleware("SUPER_ADMIN", "CASHIER", "MANAGER"), h.DisburseLoan)
				pawns.GET("/:id/payment-calc", middleware.RBACMiddleware("SUPER_ADMIN", "CASHIER", "MANAGER"), h.CalculatePayment)
				pawns.POST("/:id/payments", middleware.RBACMiddleware("SUPER_ADMIN", "CASHIER", "MANAGER"), h.ProcessPayment)
				pawns.POST("/:id/release", middleware.RBACMiddleware("SUPER_ADMIN", "CASHIER", "MANAGER"), h.ReleaseAsset)
			}

			// Inventory & QR Scanning
			inventory := protected.Group("/inventory")
			{
				inventory.GET("", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "CASHIER", "MANAGER"), h.ListInventory)
				inventory.GET("/qr/:qr", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "CASHIER", "MANAGER"), h.GetInventoryByQR)
			}

			// Dashboard & Analytics
			protected.GET("/dashboard/summary", middleware.RBACMiddleware("SUPER_ADMIN", "ADMIN", "APPRAISER", "CASHIER", "MANAGER"), h.GetDashboardSummary)

			// Audit Logs
			protected.GET("/audit-logs", middleware.RBACMiddleware("SUPER_ADMIN", "MANAGER"), h.ListAuditLogs)
		}
	}

	return r
}
