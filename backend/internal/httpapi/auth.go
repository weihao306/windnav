package httpapi

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"windnav/backend/internal/config"
	"windnav/backend/internal/model"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func (s *Server) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "VALIDATION_ERROR", "请求参数无效", err.Error())
		return
	}

	var user model.User
	if err := s.db.Where("username = ? AND is_active = ?", req.Username, true).First(&user).Error; err != nil {
		fail(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "用户名或密码错误", nil)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		fail(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "用户名或密码错误", nil)
		return
	}

	now := time.Now()
	user.LastLoginAt = &now
	_ = s.db.Save(&user).Error

	token, err := signToken(s.cfg, user)
	if err != nil {
		fail(c, http.StatusInternalServerError, "TOKEN_ERROR", "生成登录令牌失败", nil)
		return
	}
	ok(c, gin.H{"token": token, "user": publicUser(user)}, nil)
}

func (s *Server) me(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "未登录", nil)
		return
	}
	ok(c, publicUser(user), nil)
}

func (s *Server) logout(c *gin.Context) {
	ok(c, gin.H{"ok": true}, nil)
}

func (s *Server) authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "未登录", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := Claims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
			return []byte(s.cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "登录已失效", nil)
			c.Abort()
			return
		}

		var user model.User
		err = s.db.Where("id = ? AND is_active = ?", claims.UserID, true).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "用户不存在或已停用", nil)
			c.Abort()
			return
		}
		if err != nil {
			fail(c, http.StatusInternalServerError, "SERVER_ERROR", "读取用户失败", nil)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func signToken(cfg config.Config, user model.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.Username,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JWTSecret))
}

func currentUser(c *gin.Context) (model.User, bool) {
	value, exists := c.Get("user")
	if !exists {
		return model.User{}, false
	}
	user, ok := value.(model.User)
	return user, ok
}

func publicUser(user model.User) gin.H {
	return gin.H{"id": user.ID, "username": user.Username, "displayName": user.DisplayName, "role": user.Role}
}
