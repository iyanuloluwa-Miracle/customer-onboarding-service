package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-quickstart/internal/clients"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/services"
)

// CustomerController handles customer-related HTTP requests.
type CustomerController struct {
	service    *services.CustomerService
	bankClient clients.BankClientInterface
}

// NewCustomerController creates a new CustomerController.
func NewCustomerController(service *services.CustomerService, bankClient clients.BankClientInterface) *CustomerController {
	return &CustomerController{
		service:    service,
		bankClient: bankClient,
	}
}

// Onboard handles customer onboarding.
// @Summary Onboard a new customer
// @Description Register a new customer with phone, email, password, state and LGA. An OTP will be sent for phone verification.
// @Tags Customers
// @Accept json
// @Produce json
// @Param request body models.OnboardRequest true "Onboard Request"
// @Success 201 {object} models.APIResponse "OTP sent for verification"
// @Failure 400 {object} models.APIResponse "Invalid request or business rule violation"
// @Failure 500 {object} models.APIResponse "Internal server error"
// @Router /api/v1/customers/onboard [post]
func (c *CustomerController) Onboard(ctx *gin.Context) {
	var req models.OnboardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	if err := c.service.OnboardCustomer(req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Onboarding failed",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Customer registered. OTP sent to phone for verification.",
	})
}

// VerifyOTP handles OTP verification.
// @Summary Verify phone number with OTP
// @Description Complete onboarding by verifying the phone number with the OTP sent
// @Tags Customers
// @Accept json
// @Produce json
// @Param request body models.VerifyOTPRequest true "Verify OTP Request"
// @Success 200 {object} models.APIResponse "Phone verified successfully"
// @Failure 400 {object} models.APIResponse "Invalid or expired OTP"
// @Router /api/v1/customers/verify-otp [post]
func (c *CustomerController) VerifyOTP(ctx *gin.Context) {
	var req models.VerifyOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	if err := c.service.VerifyOTP(req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "OTP verification failed",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Phone number verified. Onboarding complete!",
	})
}

// ResendOTP handles resending OTP.
// @Summary Resend OTP
// @Description Resend OTP to the registered phone number
// @Tags Customers
// @Accept json
// @Produce json
// @Param request body models.ResendOTPRequest true "Resend OTP Request"
// @Success 200 {object} models.APIResponse "OTP resent successfully"
// @Failure 400 {object} models.APIResponse "Failed to resend OTP"
// @Router /api/v1/customers/resend-otp [post]
func (c *CustomerController) ResendOTP(ctx *gin.Context) {
	var req models.ResendOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request payload",
			Error:   err.Error(),
		})
		return
	}

	if err := c.service.ResendOTP(req.Phone); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Failed to resend OTP",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "OTP resent successfully",
	})
}

// GetAll handles fetching all customers.
// @Summary Get all customers
// @Description Retrieve all onboarded customers
// @Tags Customers
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.Customer} "List of customers"
// @Failure 500 {object} models.APIResponse "Internal server error"
// @Router /api/v1/customers [get]
func (c *CustomerController) GetAll(ctx *gin.Context) {
	customers, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch customers",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Customers retrieved successfully",
		Data:    customers,
	})
}

// GetBanks handles fetching banks from external API.
// @Summary Get all banks
// @Description Retrieve list of banks from external API
// @Tags Banks
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.Bank} "List of banks"
// @Failure 500 {object} models.APIResponse "Failed to fetch banks"
// @Router /api/v1/banks [get]
func (c *CustomerController) GetBanks(ctx *gin.Context) {
	banks, err := c.bankClient.GetBanks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch banks",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Banks retrieved successfully",
		Data:    banks,
	})
}
