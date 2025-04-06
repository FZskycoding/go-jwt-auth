package handlers

import (
	"jwt-practice/middleware"
	"jwt-practice/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// 註冊請求結構
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

// 登入請求結構
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}


// 登入處理
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求參數"})
		return
	}

	var user models.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用戶名或密碼錯誤"})
		return
	}

	// 驗證密碼
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用戶名或密碼錯誤"})
		return
	}

	// 生成 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 2).Unix(), // 2小時過期
	})

	tokenString, err := token.SignedString([]byte(middleware.SecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token生成失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

// 獲取用戶資料
func (h *AuthHandler) GetUserProfile(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授權"})
		return
	}

	var user models.User
	if err := h.db.First(&user, userId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用戶不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

// 獲取所有用戶資料
func (h *AuthHandler) GetAllProfile(c *gin.Context) {
    var users []models.User
    if err := h.db.Select("id, username, email").Find(&users).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "獲取用戶列表失敗"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "users": users,
    })
}

// 刪除用戶
func (h *AuthHandler) DeleteProfile(c *gin.Context) {
    // 從查詢參數中獲取要刪除的用戶ID
    userID := c.Query("id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "必須提供用戶ID"})
        return
    }

    // 查找並刪除用戶
    var user models.User
    if err := h.db.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "用戶不存在"})
        return
    }

    if err := h.db.Delete(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "刪除用戶失敗"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "用戶已成功刪除",
        "deleted_user": gin.H{
            "id": user.ID,
            "username": user.Username,
            "email": user.Email,
        },
    })
}
// 註冊處理
func (h *AuthHandler) Register(c *gin.Context) { //讓AuthHandler去呼叫Register
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求參數"})
		return
	}
	// 檢查用戶名是否已存在
	var existingUser models.User
	if err := h.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用戶名稱已被使用"})
		return
	}

	// 密碼加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密碼加密失敗"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	// 創建用戶
	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用戶創建失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "註冊成功",
		"username": req.Username})
}
