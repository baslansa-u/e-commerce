package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sq/ecommerce/pkg/response"
)

// AuthRequired ตรวจสอบ JWT token ก่อนเข้า route ที่ต้อง login
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ดึง token จาก header: "Authorization: Bearer <token>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "กรุณาเข้าสู่ระบบก่อน")
			c.Abort() // หยุด request ไม่ให้ผ่านไป handler
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "รูปแบบ token ไม่ถูกต้อง")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		secret := os.Getenv("JWT_SECRET")

		// Parse และตรวจสอบ token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "token ไม่ถูกต้องหรือหมดอายุ")
			c.Abort()
			return
		}

		// ดึง claims แล้วเก็บ user_id ไว้ใน context เพื่อใช้ใน handler
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "token ไม่ถูกต้อง")
			c.Abort()
			return
		}

		// เก็บ user_id ไว้ใน context (handler ดึงได้ด้วย c.GetUint("user_id"))
		userID := uint(claims["user_id"].(float64))
		c.Set("user_id", userID)

		c.Next() // ผ่านไป handler ต่อได้
	}
}
