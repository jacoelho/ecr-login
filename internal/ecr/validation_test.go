package ecr_test

import (
	"testing"

	"github.com/jacoelho/ecr-login/internal/ecr"
)

func TestIsCredentialValid(t *testing.T) {
	tests := []struct {
		name       string
		credential ecr.Credential
		wantErr    bool
	}{
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name: "valid credentials",
			credential: ecr.Credential{
				Username:    "foo",
				Password:    "bar",
				RegistryURL: "https://example.org",
			},
			wantErr: false,
		},
		{
			name: "url without scheme",
			credential: ecr.Credential{
				Username:    "foo",
				Password:    "bar",
				RegistryURL: "example.org",
			},
			wantErr: false,
		},
		{
			name: "username with commands",
			credential: ecr.Credential{
				Username:    "foo; rm -fr *",
				Password:    "bar",
				RegistryURL: "https://example.org",
			},
			wantErr: true,
		},
		{
			name: "invalid registry",
			credential: ecr.Credential{
				Username:    "foo",
				Password:    "bar",
				RegistryURL: "//\\example.org",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := ecr.IsCredentialValid(tt.credential); (err != nil) != tt.wantErr {
				t.Errorf("IsCredentialValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
