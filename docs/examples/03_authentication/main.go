// Package main - authentication demonstrates Before advice for auth/authz checks
package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/seyallius/gosaidno/v2/aspect"
	"github.com/seyallius/gosaidno/v2/aspect/wrap"
	"github.com/seyallius/gosaidno/v2/docs/examples/utils"
)

// -------------------------------------------- Auth System --------------------------------------------

type Session struct {
	UserID    string
	Role      string
	ExpiresAt time.Time
}

var sessions = map[string]*Session{
	"valid_token":   {UserID: "user_123", Role: "admin", ExpiresAt: time.Now().Add(1 * time.Hour)},
	"user_token":    {UserID: "user_456", Role: "user", ExpiresAt: time.Now().Add(1 * time.Hour)},
	"expired_token": {UserID: "user_789", Role: "user", ExpiresAt: time.Now().Add(-1 * time.Hour)},
}

func validateToken(token string) (*Session, error) {
	log.Printf("      🔍 [AUTH] Looking up token in session store")
	session, ok := sessions[token]
	if !ok {
		return nil, errors.New("invalid token")
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("token expired")
	}
	return session, nil
}

// -------------------------------------------- Setup --------------------------------------------

var registry = aspect.NewRegistry()

func setupAOP() {
	log.Println("=== Setting up Authentication AOP ===")

	registry.MustRegister("GetUserData")
	registry.MustRegister("DeleteUser")
	registry.MustRegister("UpdateSettings")

	// Authentication check (Before advice, priority 100)
	for _, fn := range []aspect.FuncKey{"GetUserData", "DeleteUser", "UpdateSettings"} {
		registry.MustAddAdvice(fn, aspect.Advice{
			Type:     aspect.Before,
			Priority: 100, // Highest - run first
			Handler: func(c *aspect.Context) error {
				utils.LogBefore(c, 100, "AUTHENTICATION")
				token := c.Args[0].(string)

				log.Printf("   🔐 [AUTH] Validating token: %s", token)
				session, err := validateToken(token)
				if err != nil {
					log.Printf("   ❌ [AUTH FAILED] %s - %v", c.FunctionName, err)
					return fmt.Errorf("authentication failed: %w", err)
				}

				// Store authenticated user in metadata
				c.Metadata["userID"] = session.UserID
				c.Metadata["role"] = session.Role
				log.Printf("   ✅ [AUTH SUCCESS] %s - user: %s, role: %s", c.FunctionName, session.UserID, session.Role)
				log.Printf("   💾 [METADATA] Stored user context for downstream advice")
				return nil
			},
		})
	}

	// Authorization check for DeleteUser (requires admin role)
	registry.MustAddAdvice("DeleteUser", aspect.Advice{
		Type:     aspect.Before,
		Priority: 90, // After authentication
		Handler: func(c *aspect.Context) error {
			utils.LogBefore(c, 90, "AUTHORIZATION")
			role := c.Metadata["role"].(string)
			userID := c.Metadata["userID"].(string)

			log.Printf("   🛡️  [AUTHZ] Checking if user %s has admin role (current: %s)", userID, role)
			if role != "admin" {
				log.Printf("   🚫 [AUTHZ FAILED] DeleteUser - user %s does not have admin role", userID)
				return errors.New("permission denied: admin role required")
			}

			log.Printf("   ✅ [AUTHZ SUCCESS] DeleteUser - admin access granted for user %s", userID)
			return nil
		},
	})

	// Audit logging (After advice)
	for _, fn := range []aspect.FuncKey{"DeleteUser", "UpdateSettings"} {
		registry.MustAddAdvice(fn, aspect.Advice{
			Type:     aspect.After,
			Priority: 100,
			Handler: func(c *aspect.Context) error {
				utils.LogAfter(c, 100, "AUDIT")
				userID, _ := c.Metadata["userID"].(string)
				status := "SUCCESS"
				if c.Error != nil {
					status = "FAILED"
				}

				log.Printf("   📋 [AUDIT] Function: %s", c.FunctionName)
				log.Printf("   👤 [AUDIT] User: %s", userID)
				log.Printf("   📊 [AUDIT] Status: %s", status)
				log.Printf("   🎯 [AUDIT] Args: %v", c.Args[1:])
				if c.Error != nil {
					log.Printf("   ❌ [AUDIT] Error: %v", c.Error)
				}
				log.Printf("   📝 [AUDIT] Audit trail recorded")
				return nil
			},
		})
	}

	// Success logging for GetUserData
	registry.MustAddAdvice("GetUserData", aspect.Advice{
		Type:     aspect.AfterReturning,
		Priority: 80,
		Handler: func(c *aspect.Context) error {
			utils.LogAfterReturning(c, 80, "SUCCESS LOG")
			userID := c.Metadata["userID"].(string)
			log.Printf("   📈 [METRICS] User %s successfully accessed data", userID)
			return nil
		},
	})

	log.Println("=== AOP Setup Complete ===\n")
}

// -------------------------------------------- Business Logic --------------------------------------------

func getUserDataImpl(token, userID string) (map[string]string, error) {
	log.Printf("   👨‍💼 [BUSINESS] getUserDataImpl executing for user: %s", userID)
	log.Printf("   💾 [BUSINESS] Fetching user data from database...")
	time.Sleep(25 * time.Millisecond)
	data := map[string]string{
		"id":    userID,
		"name":  "John Doe",
		"email": "john@example.com",
	}
	log.Printf("   ✅ [BUSINESS] Retrieved user data: %v", data)
	return data, nil
}

func deleteUserImpl(token, userID string) error {
	log.Printf("   🗑️  [BUSINESS] deleteUserImpl executing for user: %s", userID)
	log.Printf("   💾 [BUSINESS] Deleting user from database...")
	time.Sleep(50 * time.Millisecond)
	log.Printf("   ✅ [BUSINESS] User %s deleted successfully", userID)
	return nil
}

func updateSettingsImpl(token string, settings map[string]interface{}) error {
	log.Printf("   ⚙️  [BUSINESS] updateSettingsImpl executing with settings: %v", settings)
	log.Printf("   💾 [BUSINESS] Updating user settings in database...")
	time.Sleep(30 * time.Millisecond)
	log.Printf("   ✅ [BUSINESS] Settings updated successfully")
	return nil
}

// -------------------------------------------- Wrapped Functions --------------------------------------------

var (
	GetUserData    = wrap.Wrap2RE(registry, "GetUserData", getUserDataImpl)
	DeleteUser     = wrap.Wrap2E(registry, "DeleteUser", deleteUserImpl)
	UpdateSettings = wrap.Wrap2E(registry, "UpdateSettings", updateSettingsImpl)
)

// -------------------------------------------- Examples --------------------------------------------

func example1_SuccessfulAuth() {
	fmt.Println("\n========== Example 1: Successful Authentication ==========\n")

	data, err := GetUserData("valid_token", "user_123")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Got user data: %v\n", data)
}

func example2_InvalidToken() {
	fmt.Println("\n========== Example 2: Invalid Token ==========\n")

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("❌ Request rejected: %v\n", r)
		}
	}()

	_, _ = GetUserData("invalid_token", "user_123")
}

func example3_ExpiredToken() {
	fmt.Println("\n========== Example 3: Expired Token ==========\n")

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("❌ Request rejected: %v\n", r)
		}
	}()

	_, _ = GetUserData("expired_token", "user_789")
}

func example4_AuthorizationSuccess() {
	fmt.Println("\n========== Example 4: Authorization Success (Admin) ==========\n")

	err := DeleteUser("valid_token", "user_999")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Println("✅ User deleted successfully")
}

func example5_AuthorizationFailure() {
	fmt.Println("\n========== Example 5: Authorization Failure (Non-Admin) ==========\n")

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("❌ Request rejected: %v\n", r)
		}
	}()

	_ = DeleteUser("user_token", "user_999")
}

func example6_AuditLogging() {
	fmt.Println("\n========== Example 6: Audit Logging ==========\n")

	settings := map[string]interface{}{
		"theme":         "dark",
		"notifications": true,
	}

	err := UpdateSettings("valid_token", settings)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Println("✅ Settings updated (check audit log above)")
}

// -------------------------------------------- Main --------------------------------------------

func main() {
	setupAOP()

	example1_SuccessfulAuth()
	example2_InvalidToken()
	example3_ExpiredToken()
	example4_AuthorizationSuccess()
	example5_AuthorizationFailure()
	example6_AuditLogging()

	fmt.Println("\n========== Authentication Examples Complete ==========")
}
