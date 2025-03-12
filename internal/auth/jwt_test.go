package auth_test

import (
	"testing"
	"time"

	"github.com/GianniBuoni/chirpy/internal/auth"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantErr     bool
	}{
		{
			name:        "Test for no error",
			userID:      uuid.New(),
			tokenSecret: "secret",
			expiresIn:   10 * time.Second,
			wantErr:     false,
		},
		{
			name:        "Test different parameters",
			userID:      uuid.New(),
			tokenSecret: "b35b083de3234751815553f0198b420d",
			expiresIn:   20 * time.Minute,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := auth.MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("MakeJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MakeJWT() succeeded unexpectedly")
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	// make test ids to use in testing
	testIds := []uuid.UUID{}
	for range 5 {
		id := uuid.New()
		testIds = append(testIds, id)
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userID       uuid.UUID
		tokenSecret  string
		tokenSecret2 string
		expiresIn    time.Duration
		waitToCheck  time.Duration
		want         uuid.UUID
		wantErr      bool
	}{
		{
			name:         "Test input match output",
			userID:       testIds[0],
			tokenSecret:  "secret",
			tokenSecret2: "secret",
			expiresIn:    20 * time.Minute,
			want:         testIds[0],
			wantErr:      false,
		},
		{
			name:         "Test input match output 2",
			userID:       testIds[1],
			tokenSecret:  "another secret",
			tokenSecret2: "another secret",
			expiresIn:    5 * time.Second,
			want:         testIds[1],
			wantErr:      false,
		},
		{
			name:         "Test expiration",
			userID:       testIds[2],
			tokenSecret:  "expiration",
			tokenSecret2: "expiration",
			expiresIn:    1 * time.Second,
			waitToCheck:  2 * time.Second,
			want:         uuid.UUID{},
			wantErr:      true,
		},
		{
			name:         "Test bad secret",
			userID:       testIds[2],
			tokenSecret:  "good",
			tokenSecret2: "bad",
			expiresIn:    1 * time.Second,
			want:         uuid.UUID{},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// make tokens to check
			tokenString, err := auth.MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if err != nil {
				t.Fatalf("MakeJWT() failed %v", err)
			}
			time.Sleep(tt.waitToCheck)
			got, gotErr := auth.ValidateJWT(tokenString, tt.tokenSecret2)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ValidateJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ValidateJWT() succeeded unexpectedly")
			}
			pass := tt.userID == got
			if !pass {
				t.Errorf("ValidateJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}
