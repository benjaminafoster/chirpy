package auth

import (
	"testing"
	"time"
	"net/http"
	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
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
	tests := []struct{
		name string
		headers http.Header
		wantToken string
		wantErr bool
	}{
		{
			name: "Valid Token -- 12345",
			headers: map[string][]string{"Authorization": []string{"Bearer 12345",}},
			wantToken: "12345",
			wantErr: false,
		},
		{
			name: "Valid Token -- 67890",
			headers: map[string][]string{"Authorization": []string{"Bearer 67890",}},
			wantToken: "67890",
			wantErr: false,
		},
		{
			name: "Invalid Token -- Insufficient bearer token fields",
			headers: map[string][]string{"Authorization": []string{"Bearer",}},
			wantToken: "",
			wantErr: true,
		},
		{
			name: "Invalid Token -- 'Bearer' improperly formatted",
			headers: map[string][]string{"Authorization": []string{"bearer 67890",}},
			wantToken: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			_, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return 
			}
		})
	}
}