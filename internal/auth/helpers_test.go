package auth_test

import (
	"net/http"
	"testing"

	"github.com/GianniBuoni/chirpy/internal/auth"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		headers http.Header
		token   string
		want    string
		wantErr bool
	}{
		{
			name:    "Test valid",
			headers: http.Header{},
			token:   "faketoken",
			want:    "faketoken",
		},
		{
			name:    "Test invalid",
			headers: http.Header{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				tt.headers.Add("Authorization", "Bearer "+tt.token)
			}
			got, gotErr := auth.GetBearerToken(tt.headers)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBearerToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBearerToken() succeeded unexpectedly")
			}
			pass := tt.want == got
			if !pass {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
