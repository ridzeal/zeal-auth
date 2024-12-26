package auth_test

import (
	"time"
	"strings"
	"testing"
	"github.com/dgrijalva/jwt-go"
	"sso-backend/auth"
)

// Helper function to validate token claims
func validateTokenClaims(t *testing.T, tokenString string, expectedUsername string) {
	token, err := auth.ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to get token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		t.Fatal("Username claim not found or invalid type")
	}

	if username != expectedUsername {
		t.Errorf("Expected username %s, got %s", expectedUsername, username)
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatal("Expiration claim not found or invalid type")
	}

	if time.Unix(int64(exp), 0).Before(time.Now()) {
		t.Error("Token is already expired")
	}
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
		expectedUsername string
	}{
		// Basic cases
		{
			name:     "valid email username",
			username: "test@example.com",
			wantErr:  false,
		},
		{
			name:     "empty username",
			username: "",
			wantErr:  true,
		},
		{
			name:     "very long username",
			username: "thisissuchaverylongusernamethatmightcauseissues@reallyreallylongdomain.com",
			wantErr:  false,
		},
		{
			name:     "special characters in username",
			username: "test.user+123@example.com",
			wantErr:  false,
		},

		// Whitespace cases
		{
			name:     "whitespace only username",
			username: "   ",
			wantErr:  true,
		},
		{
			name:     "username with leading/trailing spaces",
			username: " test@example.com ",
			wantErr:  false,
			expectedUsername: "test@example.com",
		},

		// Special character cases
		{
			name:     "username with newline characters",
			username: "test@example.com\n",
			wantErr:  true,
		},
		{
			name:     "username with control characters",
			username: "test\x00@example.com",
			wantErr:  true,
		},
		{
			name:     "username with tabs",
			username: "test\t@example.com",
			wantErr:  true,
		},
		{
			name:     "username with carriage return",
			username: "test\r@example.com",
			wantErr:  true,
		},

		// Length cases
		{
			name:     "extremely short username",
			username: "a@b.c",
			wantErr:  false,
		},
		{
			name:     "maximum length username",
			username: strings.Repeat("a", 254) + "@b.c",
			wantErr:  false,
		},

		// Security cases
		{
			name:     "SQL injection attempt",
			username: "'; DROP TABLE users; --",
			wantErr:  true,
		},
		{
			name:     "HTML injection attempt",
			username: "<script>alert('xss')</script>",
			wantErr:  true,
		},

		// Unicode cases
		{
			name:     "unicode characters",
			username: "用户@example.com",
			wantErr:  false,
		},
		{
			name:     "username with emoji",
			username: "test😊@example.com",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := auth.GenerateToken(tt.username)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if token != "" {
					t.Errorf("expected empty token but got: %s", token)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if token == "" {
					t.Errorf("expected non-empty token but got empty string")
				}
				
				expectedUsername := tt.expectedUsername
                if expectedUsername == "" {
                    expectedUsername = tt.username
                }
                
				// Validate token claims for successful cases
				if token != "" {
					validateTokenClaims(t, token, expectedUsername)
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name      string
		setupToken func() string
		wantErr   bool
		errorType string
	}{
		{
			name: "valid token",
			setupToken: func() string {
				token, _ := auth.GenerateToken("test@example.com")
				return token
			},
			wantErr: false,
		},
		{
			name: "empty token",
			setupToken: func() string {
				return ""
			},
			wantErr:   true,
			errorType: "token contains an invalid number of segments",
		},
		{
			name: "malformed token",
			setupToken: func() string {
				return "invalid.token.format"
			},
			wantErr:   true,
			errorType: "illegal base64",
		},
		{
			name: "expired token",
			setupToken: func() string {
				token := jwt.New(jwt.SigningMethodHS256)
				claims := token.Claims.(jwt.MapClaims)
				claims["username"] = "test@example.com"
				claims["exp"] = time.Now().Add(-time.Hour).Unix()
				tokenString, _ := token.SignedString([]byte("your-secret-key"))
				return tokenString
			},
			wantErr:   true,
			errorType: "Token is expired",
		},
		{
			name: "invalid signing method",
			setupToken: func() string {
				token := jwt.New(jwt.SigningMethodRS256)
				claims := token.Claims.(jwt.MapClaims)
				claims["username"] = "test@example.com"
				claims["exp"] = time.Now().Add(time.Hour).Unix()
				tokenString, _ := token.SignedString([]byte("wrong-key"))
				return tokenString
			},
			wantErr:   true,
			errorType: "invalid signing method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			validated, err := auth.ValidateToken(token)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else {
					var matchedError bool
                    expectedErrors := []string{
                        tt.errorType,
                        "invalid character",
                        "illegal base64",
                        "token contains an invalid number of segments",
                        "Token is expired",
                    }
                    
                    for _, expectedError := range expectedErrors {
                        if strings.Contains(err.Error(), expectedError) {
                            matchedError = true
                            break
                        }
                    }
                    
                    if !matchedError {
                        t.Errorf("expected one of the error types %v but got '%s'", 
                            expectedErrors, err.Error())
                    }
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
				if validated == nil {
					t.Errorf("expected valid token but got nil")
				}
				// For valid tokens, verify the claims
				if validated != nil {
					claims, ok := validated.Claims.(jwt.MapClaims)
					if !ok {
						t.Error("failed to get token claims")
					} else {
						if _, ok := claims["username"]; !ok {
							t.Error("username claim not found")
						}
						if _, ok := claims["exp"]; !ok {
							t.Error("exp claim not found")
						}
					}
				}
			}
		})
	}
}