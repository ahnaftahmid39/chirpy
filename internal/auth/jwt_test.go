package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTFunctions(t *testing.T) {
	tests := []struct {
		name        string
		tokenSecret string
		expiresIn   time.Duration
		shouldPass  bool
	}{
		{"ValidToken", "supersecret", time.Minute * 5, true},
		{"ExpiredToken", "supersecret", -time.Minute, false},
		{"EmptySecret", "", time.Minute * 5, false},
		{"LongExpiry", "supersecret", time.Hour * 24, true},
	}

	passCount := 0
	failCount := 0

	for _, tc := range tests {
		userID := uuid.New()
		token, err := MakeJWT(userID, tc.tokenSecret, tc.expiresIn)

		if tc.shouldPass && err != nil {
			failCount++
			t.Errorf(`---------------------------------
Test Failed (Token Generation):
Inputs:     (userID: %v, tokenSecret: %q, expiresIn: %v)
Expecting:  Success
Actual:     Error: %v
`, userID, tc.tokenSecret, tc.expiresIn, err)
		} else if !tc.shouldPass && err == nil {
			failCount++
			t.Errorf(`---------------------------------
Test Failed (Expected Failure in Token Generation):
Inputs:     (userID: %v, tokenSecret: %q, expiresIn: %v)
Expecting:  Error
Actual:     Token generated
`, userID, tc.tokenSecret, tc.expiresIn)
		} else if tc.shouldPass {
			// Validate the token
			parsedUserID, err := ValidateJWT(token, tc.tokenSecret)

			if err != nil {
				failCount++
				t.Errorf(`---------------------------------
Test Failed (Token Validation):
Inputs:     (token: %q, tokenSecret: %q)
Expecting:  Valid userID
Actual:     Error: %v
`, token, tc.tokenSecret, err)
			} else if parsedUserID != userID {
				failCount++
				t.Errorf(`---------------------------------
Test Failed (UserID Mismatch):
Inputs:     (token: %q, tokenSecret: %q)
Expecting:  %v
Actual:     %v
`, token, tc.tokenSecret, userID, parsedUserID)
			} else {
				passCount++
				fmt.Printf(`---------------------------------
Test Passed:
Inputs:     (userID: %v, tokenSecret: %q, expiresIn: %v)
Generated Token: %q
Validated UserID: %v
`, userID, tc.tokenSecret, tc.expiresIn, token, parsedUserID)
			}
		} else {
			passCount++
			fmt.Printf(`---------------------------------
Test Passed (Expected Failure):
Inputs:     (userID: %v, tokenSecret: %q, expiresIn: %v)
Expected failure and got failure
`, userID, tc.tokenSecret, tc.expiresIn)
		}
	}

	fmt.Println("---------------------------------")
	fmt.Printf("%d passed, %d failed\n", passCount, failCount)
}

func TestValidateJWTFromSolution(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		authHeaderKey   string
		authHeaderValue string
		wantErr         bool
		expectedToken   string
	}{
		{
			authHeaderKey:   "Authorization",
			authHeaderValue: "",
			wantErr:         true,
			expectedToken:   "",
		},
		{
			authHeaderKey:   "Authorization",
			authHeaderValue: "Bearer",
			wantErr:         true,
			expectedToken:   "",
		},
		{
			authHeaderKey:   "Authorization",
			authHeaderValue: "Bearer abc",
			wantErr:         false,
			expectedToken:   "abc",
		},
		{
			authHeaderKey:   "auth",
			authHeaderValue: "Bearer abc",
			wantErr:         true,
			expectedToken:   "",
		},
		{
			authHeaderKey:   "Authorization",
			authHeaderValue: "bearer abc",
			wantErr:         true,
			expectedToken:   "",
		},
		{
			authHeaderKey:   "Authorization",
			authHeaderValue: "Bearer abc def",
			wantErr:         true,
			expectedToken:   "",
		},
	}

	failedCount := 0
	for _, tt := range tests {
		headers := http.Header{}
		headers.Add(tt.authHeaderKey, tt.authHeaderValue)

		token, err := GetBearerToken(headers)
		if (err != nil) != tt.wantErr {
			failedCount++
			t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
		}
		if token != tt.expectedToken {
			failedCount++
			t.Errorf("GetBearerToken() gotToken = %v, want %v", token, tt.expectedToken)
		}
	}

	t.Logf("Passed: %v, Failed: %v\n", len(tests)-failedCount, failedCount)
}
