package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	pkg "github.com/IAmRiteshKoushik/pulse/pkg"
	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

func LoginUserCsrf(c *gin.Context) {

	token, err := pkg.NewCsrfToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token.",
		})
		return
	}

	// TODO: Setup cookie

	c.JSON(http.StatusOK, gin.H{
		"message": "Token generated successfully.",
		"token":   token,
	})
	return
}

func LoginUser(c *gin.Context) {

}

func RegisterUserAccountCsrf(c *gin.Context) {
	// Create CSRF token
	token, err := pkg.NewCsrfToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token.",
		})
		return
	}

	// Set CSRF cookie
	pkg.SetCsrfCookie(c, token, "/user/register")

	c.JSON(http.StatusOK, gin.H{
		"message": "Token generated successfully.",
		"token":   token,
	})
}

func RegisterUserAccount(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Fullname string `json:"fullname" binding:"required"`
		// PhoneNumber string `json:"phone_number" binding:"required,numeric,len=10"`
	}

	// Validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.RequestValidatorError(c, err)
		return
	}

	// Check if the user already exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	err := cmd.DBPool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM user_account WHERE email = $1)`, req.Email).Scan(&exists)

	if err != nil {
		pkg.DbError(c, err)
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User with this email already exists.",
		})
		return
	}

	// Generate OTP
	otp, err := pkg.GenerateOTP()
	if err != nil {
		cmd.Log.Error("[OTP-ERROR] Failed to generate OTP", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate OTP. Please try again later.",
		})
		return
	}

	// Store OTP in the database
	expiration := time.Now().Add(time.Minute * 5)
	_, err = cmd.DBPool.Exec(ctx, `INSERT INTO otps (id, email, full_name, otp, expires_at) VALUES ($1, $2, $3, $4, $5)`, uuid.New(), req.Email, req.Fullname, otp, expiration)

	if err != nil {
		pkg.DbError(c, err)
		return
	}

	// Send OTP to the user's email
	// go func() {
	// 	if err := pkg.SendOTPEmail(req.Email, otp); err != nil {
	// 		cmd.Log.Error("[EMAIL-ERROR] Failed to send OTP", err)
	// 	}
	// }()
	fmt.Println("OTP: ", otp);

	// Generate Temp Token with user data
	// additionalData := map[string]any{
	// 	"phone_number": req.PhoneNumber,
	// }
	// tempToken := pkg.CreateTempToken(req.Fullname, req.Email)
	// TODO: Set temp token to cookie
	c.SetCookie(
		"email",             // key
		req.Email,                // value
		900,                      // maxAge (15 minutes)
		"/user/register",         // path to be constructed for restriction
		cmd.EnvVars.Domain,       // domain
		cmd.EnvVars.CookieSecure, // secure
		true,                     // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message":    "OTP sent successfully.",
	})
}

func VerifyOtpCsrf(c *gin.Context) {
	// Create CSRF token
	token, err := pkg.NewCsrfToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token.",
		})
		return
	}

	// Set CSRF cookie
	pkg.SetCsrfCookie(c, token, "/user/register/otp")

	c.JSON(http.StatusOK, gin.H{
		"message": "Token generated successfully.",
		"token":   token,
	})
}

func VerifyOtp(c *gin.Context) {
	email, err := c.Cookie("email")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found in cookie"})
		return
	}

	// Parse the OTP from the JSON body
	var requestBody struct {
		Otp string `json:"otp"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}
	otp := requestBody.Otp
	if otp == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP is required"})
		return
	}

	// Retrieve the latest OTP for the given email from the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var latestOtp string
	var expiresAt time.Time
	var fullName string
		err = cmd.DBPool.QueryRow(ctx, `
		SELECT otp, expires_at, full_name 
		FROM otps 
		WHERE email = $1 
		ORDER BY expires_at DESC 
		LIMIT 1
	`, email).Scan(&latestOtp, &expiresAt, &fullName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve OTP"})
		return
	}

	// Check if the OTP is valid
	if otp != latestOtp {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	if time.Now().After(expiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP has expired"})
		return
	}

	// Delete all OTP records associated with the email after successful verification
	_, err = cmd.DBPool.Exec(ctx, `DELETE FROM otps WHERE email = $1`, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete OTP records"})
		return
	}

	userRole := true
	hostRole := false
	staffRole := false

	// Generate new Auth and Refresh tokens
	authToken := pkg.CreateAuthToken(fullName, email, userRole, hostRole, staffRole)
	refreshToken := pkg.CreateRefreshToken(fullName, email, userRole, hostRole, staffRole)

	// Create the user account in the database
	_, err = cmd.DBPool.Exec(ctx, `
		INSERT INTO user_account (id, email, full_name, refresh_token) 
		VALUES ($1, $2, $3, $4)
	`, uuid.New(), email, fullName, refreshToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user account"})
		return
	}

	// Set the tokens in the cookies using the cookie setters from pkg
	pkg.SetAuthCookie(c, authToken)
	pkg.SetRefreshCookie(c, refreshToken)

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "OTP verified and user account created successfully!"})
}

func ResendUserOtpCsrf(c *gin.Context) {

	token, err := pkg.NewCsrfToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token.",
		})
		return
	}

	// TODO: Setup cookie

	c.JSON(http.StatusOK, gin.H{
		"message": "Token generated successfully.",
		"token":   token,
	})
	return
}

func ResendUserOtp(c *gin.Context) {

}

func UserSession(c *gin.Context) {

}

func LogoutUser(c *gin.Context) {

}
