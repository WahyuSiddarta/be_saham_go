package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/WahyuSiddarta/be_saham_go/helper"
	"github.com/WahyuSiddarta/be_saham_go/middleware"
	"github.com/WahyuSiddarta/be_saham_go/models"
	"github.com/WahyuSiddarta/be_saham_go/validator"
	"github.com/labstack/echo/v4"
)

// AuthHandlers contains all authentication-related handlers
type AuthHandlers struct{}

// NewAuthHandlers creates a new instance of auth handlers
func NewAuthHandlers() *AuthHandlers {
	return &AuthHandlers{}
}

// convertPaymentData converts payment request data to model
func convertPaymentData(reqPaymentData *validator.PaymentDataRequest) (*models.PaymentData, error) {
	if reqPaymentData == nil {
		return nil, nil
	}

	paymentData := &models.PaymentData{
		OriginalPrice:  reqPaymentData.OriginalPrice,
		PaidPrice:      reqPaymentData.PaidPrice,
		DiscountAmount: reqPaymentData.DiscountAmount,
		DiscountReason: reqPaymentData.DiscountReason,
		PaymentMethod:  reqPaymentData.PaymentMethod,
		Notes:          reqPaymentData.Notes,
	}

	// Parse payment date if provided
	if reqPaymentData.PaymentDate != nil {
		paymentDate, err := time.Parse(time.RFC3339, *reqPaymentData.PaymentDate)
		if err != nil {
			return nil, err
		}
		paymentData.PaymentDate = &paymentDate
	}

	return paymentData, nil
}

// Login handles user authentication
func (h *AuthHandlers) Login(c echo.Context) error {
	// Get validated request from middleware
	req := validator.GetValidatedRequest(c).(*validator.LoginRequest)

	// Find user by email
	userRepo := models.NewUserRepository()
	user, err := userRepo.FindByEmail(req.Email)
	if err != nil {
		if Logger != nil {
			Logger.Error().Err(err).Str("email", req.Email).Msg("Gagal mencari pengguna saat login")
		}
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan server", "Kesalahan saat proses login")
	}

	if user == nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Email atau password tidak valid", "Pengguna tidak ditemukan")
	}

	// Validate password
	if err := userRepo.ValidatePassword(req.Password, user.Password); err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Email atau password tidak valid", "Validasi password gagal")
	}

	// Check if user account is active
	if user.Status != models.UserStatusActive {
		return helper.ErrorResponse(c, http.StatusForbidden, "Akses akun ditolak", "Status akun: "+string(user.Status))
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		if Logger != nil {
			Logger.Error().Err(err).Int("user_id", user.ID).Msg("Gagal membuat token")
		}
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan server", "Kesalahan saat membuat token autentikasi")
	}

	// Remove password from response
	user.Password = ""

	if Logger != nil {
		Logger.Info().Str("email", req.Email).Int("user_id", user.ID).Msg("Pengguna berhasil login")
	}

	return helper.JsonResponse(c, http.StatusOK, validator.LoginData{
		User:  user,
		Token: token,
	})
}

// Register handles user registration
func (h *AuthHandlers) Register(c echo.Context) error {
	// Get validated request from middleware
	req := validator.GetValidatedRequest(c).(*validator.RegisterRequest)

	// Check if user already exists
	userRepo := models.NewUserRepository()
	existingUser, err := userRepo.FindByEmail(req.Email)
	if err != nil {
		if Logger != nil {
			Logger.Error().Err(err).Str("email", req.Email).Msg("Gagal memeriksa pengguna saat registrasi")
		}
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan server", "Kesalahan saat proses registrasi")
	}

	if existingUser != nil {
		return helper.ErrorResponse(c, http.StatusConflict, "Email sudah terdaftar", "Email sudah digunakan oleh pengguna lain")
	}

	// Set default values if not provided
	status := models.UserStatusActive
	if req.Status != nil {
		status = *req.Status
	}

	userLevel := models.UserLevelFree
	if req.UserLevel != nil {
		userLevel = *req.UserLevel
	}

	createReq := &models.CreateUserRequest{
		Email:     req.Email,
		Password:  req.Password,
		Status:    status,
		UserLevel: userLevel,
	}

	// Create new user
	newUser, err := userRepo.Create(createReq)
	if err != nil {
		if Logger != nil {
			Logger.Error().Err(err).Str("email", req.Email).Msg("Gagal membuat pengguna saat registrasi")
		}
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan server", "Gagal membuat akun pengguna")
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(newUser.ID)
	if err != nil {
		if Logger != nil {
			Logger.Error().Err(err).Int("user_id", newUser.ID).Msg("Gagal membuat token untuk pengguna baru")
		}
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan server", "Error generating authentication token")
	}

	if Logger != nil {
		Logger.Info().Str("email", req.Email).Int("user_id", newUser.ID).Msg("Pengguna baru berhasil didaftarkan")
	}

	return helper.JsonResponse(c, http.StatusCreated, validator.RegisterData{
		User:  newUser,
		Token: token,
	})
}

// GetProfile returns the current user's profile
func (h *AuthHandlers) GetProfile(c echo.Context) error {
	// Get authenticated user from middleware
	authUser, err := middleware.GetAuthUser(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Autentikasi diperlukan", "Unable to retrieve user information")
	}

	return helper.JsonResponse(c, http.StatusOK, validator.ProfileData{
		User: authUser,
	})
}

// Logout handles user logout (token-based, client-side cleanup)
func (h *AuthHandlers) Logout(c echo.Context) error {
	// Get authenticated user for logging purposes
	authUser, err := middleware.GetAuthUser(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Autentikasi diperlukan", "Unable to retrieve user information")
	}

	if Logger != nil {
		Logger.Info().Str("email", authUser.Email).Int("user_id", authUser.ID).Msg("User logged out")
	}

	return helper.JsonResponse(c, http.StatusOK, map[string]interface{}{
		"message": "Logout berhasil. Silakan hapus token dari penyimpanan klien.",
	})
}

// UpdateUserLevel handles updating user subscription level (admin only)
func (h *AuthHandlers) UpdateUserLevel(c echo.Context) error {
	// Get user ID from path parameter
	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID pengguna tidak valid", "User ID must be a valid number")
	}

	// Get validated request from middleware
	req := validator.GetValidatedRequest(c).(*validator.UpdateUserLevelRequest)

	// Validate custom logic
	if errs := validator.CustomValidation(req); len(errs) > 0 {
		return helper.ErrorResponse(c, http.StatusBadRequest, "Kesalahan validasi", validator.ValidationErrorResponse{
			Success: false,
			Message: "Kesalahan validasi",
			Errors:  errs,
		})
	}

	// Get admin user ID for audit trail
	adminUser, err := middleware.GetAuthUser(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Autentikasi diperlukan", "Admin authentication required")
	}

	// Convert payment data if provided
	paymentData, err := convertPaymentData(req.PaymentData)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "Format tanggal pembayaran tidak valid", "Tanggal pembayaran harus dalam format ISO 8601")
	}

	// Update user level
	userRepo := models.NewUserRepository()
	result, err := userRepo.UpdateUserLevel(userID, req.UserLevel, paymentData, &adminUser.ID)
	if err != nil {
		if Logger != nil {
			Logger.Error().Err(err).Int("user_id", userID).Str("new_level", string(req.UserLevel)).Msg("Gagal memperbarui level pengguna")
		}
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui level pengguna", err.Error())
	}

	if Logger != nil {
		Logger.Info().Int("user_id", userID).Str("new_level", string(req.UserLevel)).Int("admin_id", adminUser.ID).Msg("Level pengguna diperbarui oleh admin")
	}

	return helper.JsonResponse(c, http.StatusOK, map[string]interface{}{
		"message": "Level pengguna berhasil diperbarui",
		"data":    result,
	})
}

// UpdateUserStatus handles updating user account status (admin only)
func (h *AuthHandlers) UpdateUserStatus(c echo.Context) error {
	// Get user ID from path parameter
	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusBadRequest, "ID pengguna tidak valid", "User ID must be a valid number")
	}

	// Get validated request from middleware
	req := validator.GetValidatedRequest(c).(*validator.UpdateUserStatusRequest)

	// Get admin user for logging
	adminUser, err := middleware.GetAuthUser(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Autentikasi diperlukan", "Admin authentication required")
	}

	// Update user status
	userRepo := models.NewUserRepository()
	updatedUser, err := userRepo.UpdateUserStatus(userID, req.Status)
	if err != nil {
		if Logger != nil {
			Logger.Error().Err(err).Int("user_id", userID).Str("new_status", string(req.Status)).Msg("Error updating user status")
		}
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Error updating user status", err.Error())
	}

	if Logger != nil {
		Logger.Info().Int("user_id", userID).Str("new_status", string(req.Status)).Int("admin_id", adminUser.ID).Msg("User status updated by admin")
	}

	return helper.JsonResponse(c, http.StatusOK, map[string]interface{}{
		"message": "Status pengguna berhasil diperbarui",
		"data": map[string]interface{}{
			"user": updatedUser,
		},
	})
}

// GetAllUsers returns paginated list of users with optional filters (admin only)
func (h *AuthHandlers) GetAllUsers(c echo.Context) error {
	// Get validated query parameters
	query := validator.GetValidatedQuery(c).(*validator.GetUsersQuery)

	// Set defaults
	page := query.Page
	if page <= 0 {
		page = 1
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}

	// Get users with filters
	userRepo := models.NewUserRepository()
	result, err := userRepo.GetAllUsers(page, limit, query.Status, query.UserLevel, query.EmailFilter)
	if err != nil {
		if Logger != nil {
			Logger.Error().Err(err).Msg("Error fetching users")
		}
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Error fetching users", err.Error())
	}

	return helper.JsonResponse(c, http.StatusOK, map[string]interface{}{
		"data": result,
	})
}

// DowngradeExpiredUsers downgrades users with expired premium subscriptions (admin/cron)
func (h *AuthHandlers) DowngradeExpiredUsers(c echo.Context) error {
	// Get admin user for logging
	adminUser, err := middleware.GetAuthUser(c)
	if err != nil {
		return helper.ErrorResponse(c, http.StatusUnauthorized, "Autentikasi diperlukan", "Admin authentication required")
	}

	// Downgrade expired users
	userRepo := models.NewUserRepository()
	result, err := userRepo.DowngradeExpiredUsers()
	if err != nil {
		if Logger != nil {
			Logger.Error().Err(err).Msg("Error downgrading expired users")
		}
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Error downgrading expired users", err.Error())
	}

	if Logger != nil {
		Logger.Info().Int("downgraded_count", result.DowngradedCount).Int("admin_id", adminUser.ID).Msg("Expired users downgraded")
	}

	return helper.JsonResponse(c, http.StatusOK, map[string]interface{}{
		"message": "Pengguna kedaluwarsa berhasil diturunkan tingkatnya",
		"data":    result,
	})
}

// GetExpiredUsers returns users with expired premium subscriptions (admin only)
func (h *AuthHandlers) GetExpiredUsers(c echo.Context) error {
	// Get expired users
	userRepo := models.NewUserRepository()
	expiredUsers, err := userRepo.GetExpiredUsers()
	if err != nil {
		if Logger != nil {
			Logger.Error().Err(err).Msg("Error fetching expired users")
		}
		return helper.ErrorResponse(c, http.StatusInternalServerError, "Error fetching expired users", err.Error())
	}

	return helper.JsonResponse(c, http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"expired_users": expiredUsers,
			"count":         len(expiredUsers),
		},
	})
}
