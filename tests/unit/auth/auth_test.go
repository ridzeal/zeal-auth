package auth_test

import (
	"testing"
    "sso-backend/auth"
)

func TestGenerateToken(t *testing.T) {
    tests := []struct {
        name     string
        username string
        wantErr  bool
    }{
        // Existing cases
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
        
        // New cases
        {
            name:     "whitespace only username",
            username: "   ",
            wantErr:  true,
        },
        {
            name:     "username with leading/trailing spaces",
            username: " test@example.com ",
            wantErr:  false,
        },
        {
            name:     "username with newline characters",
            username: "test@example.com\n",
            wantErr:  true,
        },
        {
            name:     "extremely short username",
            username: "a@b.c",
            wantErr:  false,
        },
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
        {
            name:     "unicode characters",
            username: "用户@example.com",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            token, err := auth.GenerateToken(tt.username)
            
            if tt.wantErr {
                if err == nil {
                    t.Errorf("TestCase %s: expected error but got none", tt.name)
                }
                if token != "" {
                    t.Errorf("TestCase %s: expected empty token but got: %s", tt.name, token)
                }
            } else {
                if err != nil {
                    t.Errorf("TestCase %s: expected no error but got: %v", tt.name, err)
                }
                if token == "" {
                    t.Errorf("TestCase %s: expected non-empty token but got empty string", tt.name)
                }
            }
        })
    }
}