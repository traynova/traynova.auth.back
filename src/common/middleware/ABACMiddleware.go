package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ABACPolicy defines an Attribute-Based Access Control evaluator function
type ABACPolicy func(ctx *gin.Context, currentUserID uint, currentRoleID uint, targetResourceID uint) bool

// RequireResourceOwnership returns a Gin middleware that enforces ABAC / PBAC resource access controls.
// Rule:
// - RoleAdmin (4): Access granted to all resources.
// - Self-access: Access granted if currentUserID == targetResourceID.
// - Custom Policy: Evaluated if provided.
func RequireResourceOwnership(paramName string, customPolicy ...ABACPolicy) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
			c.Abort()
			return
		}

		currentUserID, ok := val.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "ID de usuario en contexto es inválido"})
			c.Abort()
			return
		}

		roleVal, _ := c.Get("role_id")
		currentRoleID, _ := roleVal.(uint)

		// Admin (Role 1 or 4 depending on setup) can access everything
		if currentRoleID == 1 || currentRoleID == 4 {
			c.Next()
			return
		}

		// Extract target resource ID from request param (e.g., :id)
		targetIDStr := c.Param(paramName)
		if targetIDStr != "" {
			targetID, err := strconv.ParseUint(targetIDStr, 10, 32)
			if err == nil {
				uintTargetID := uint(targetID)
				// Self ownership check
				if currentUserID == uintTargetID {
					c.Next()
					return
				}

				// Evaluate custom policy if passed (e.g., trainer-client relationship)
				if len(customPolicy) > 0 && customPolicy[0] != nil {
					if customPolicy[0](c, currentUserID, currentRoleID, uintTargetID) {
						c.Next()
						return
					}
				}
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado: No posee permisos sobre este recurso (ABAC Policy)"})
		c.Abort()
	}
}
