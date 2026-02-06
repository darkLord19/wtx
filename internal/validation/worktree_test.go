package validation

import (
	"strings"
	"testing"
)

func TestValidateName(t *testing.T) {
	v := NewWorktreeValidator()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid simple name",
			input:   "feature-auth",
			wantErr: false,
		},
		{
			name:    "valid with underscores",
			input:   "feature_user_auth",
			wantErr: false,
		},
		{
			name:    "valid with numbers",
			input:   "feature-123",
			wantErr: false,
		},
		{
			name:    "valid with forward slash",
			input:   "feature/auth",
			wantErr: false,
		},
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
			errMsg:  "cannot be empty",
		},
		{
			name:    "name with spaces",
			input:   "feature auth",
			wantErr: true,
			errMsg:  "can only contain",
		},
		{
			name:    "name too long",
			input:   strings.Repeat("a", 129),
			wantErr: true,
			errMsg:  "too long",
		},
		{
			name:    "dot name",
			input:   ".",
			wantErr: true,
			errMsg:  "cannot use",
		},
		{
			name:    "double dot name",
			input:   "..",
			wantErr: true,
			errMsg:  "cannot use",
		},
		{
			name:    "name with special chars",
			input:   "feature@auth",
			wantErr: true,
			errMsg:  "can only contain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateName(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateName() expected error for %q, got nil", tt.input)
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateName() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateName() unexpected error for %q: %v", tt.input, err)
				}
			}
		})
	}
}

func TestValidateBranchName(t *testing.T) {
	v := NewWorktreeValidator()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid branch name",
			input:   "feature/auth",
			wantErr: false,
		},
		{
			name:    "valid with hyphens",
			input:   "feature-auth",
			wantErr: false,
		},
		{
			name:    "empty is allowed",
			input:   "",
			wantErr: false,
		},
		{
			name:    "branch too long",
			input:   strings.Repeat("a", 256),
			wantErr: true,
			errMsg:  "too long",
		},
		{
			name:    "contains double dot",
			input:   "feature..auth",
			wantErr: true,
			errMsg:  "cannot contain",
		},
		{
			name:    "contains tilde",
			input:   "feature~1",
			wantErr: true,
			errMsg:  "cannot contain",
		},
		{
			name:    "contains caret",
			input:   "feature^auth",
			wantErr: true,
			errMsg:  "cannot contain",
		},
		{
			name:    "contains colon",
			input:   "feature:auth",
			wantErr: true,
			errMsg:  "cannot contain",
		},
		{
			name:    "contains space",
			input:   "feature auth",
			wantErr: true,
			errMsg:  "cannot contain",
		},
		{
			name:    "starts with slash",
			input:   "/feature",
			wantErr: true,
			errMsg:  "cannot start or end",
		},
		{
			name:    "ends with slash",
			input:   "feature/",
			wantErr: true,
			errMsg:  "cannot start or end",
		},
		{
			name:    "ends with .lock",
			input:   "feature.lock",
			wantErr: true,
			errMsg:  "cannot end with",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateBranchName(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateBranchName() expected error for %q, got nil", tt.input)
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateBranchName() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateBranchName() unexpected error for %q: %v", tt.input, err)
				}
			}
		})
	}
}

func TestValidatePath(t *testing.T) {
	v := NewWorktreeValidator()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid relative path",
			input:   "../worktrees",
			wantErr: false,
		},
		{
			name:    "valid absolute path",
			input:   "/home/user/worktrees",
			wantErr: false,
		},
		{
			name:    "empty path",
			input:   "",
			wantErr: true,
			errMsg:  "cannot be empty",
		},
		{
			name:    "path too long",
			input:   strings.Repeat("a", 4097),
			wantErr: true,
			errMsg:  "too long",
		},
		{
			name:    "path with spaces",
			input:   "../work trees",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidatePath(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePath() expected error for %q, got nil", tt.input)
					return
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidatePath() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePath() unexpected error for %q: %v", tt.input, err)
				}
			}
		})
	}
}

func TestSanitizeName(t *testing.T) {
	v := NewWorktreeValidator()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no changes needed",
			input: "feature-auth",
			want:  "feature-auth",
		},
		{
			name:  "replace spaces with hyphens",
			input: "feature auth",
			want:  "feature-auth",
		},
		{
			name:  "trim whitespace",
			input: "  feature-auth  ",
			want:  "feature-auth",
		},
		{
			name:  "collapse multiple hyphens",
			input: "feature--auth",
			want:  "feature-auth",
		},
		{
			name:  "complex sanitization",
			input: "  feature  --  auth  ",
			want:  "feature-auth",
		},
		{
			name:  "multiple spaces",
			input: "feature   auth   login",
			want:  "feature-auth-login",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := v.SanitizeName(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func BenchmarkValidateName(b *testing.B) {
	v := NewWorktreeValidator()
	name := "feature-user-authentication"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.ValidateName(name)
	}
}

func BenchmarkValidateBranchName(b *testing.B) {
	v := NewWorktreeValidator()
	branch := "feature/user-authentication"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.ValidateBranchName(branch)
	}
}

func BenchmarkSanitizeName(b *testing.B) {
	v := NewWorktreeValidator()
	name := "  feature  --  auth  "
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.SanitizeName(name)
	}
}
