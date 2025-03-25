package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "SecurePassword123",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false, // bcrypt will hash an empty string, though it's not recommended
		},
		{
			name:     "Long password close to limit",
			password: "ThisIsALongPasswordButStillValid1234567890abcdefghijklmno", // Less than 72 bytes
			wantErr:  false,
		},
		{
			name:     "Too long password",
			password: "ThisIsAVeryLongPasswordThatExceedsTheBcryptLimitOf72BytesAndShouldGenerateAnErrorWhenAttemptingToHashIt12345678901234567890",
			wantErr:  true, // bcrypt has a 72 byte limit
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Hash the password
			hash, err := HashPassword(tt.password)

			// Check if error was expected
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			// Verify no error occurred
			assert.NoError(t, err)

			// Check that hash is not empty
			assert.NotEmpty(t, hash)

			// Verify the hash is different from the original password
			assert.NotEqual(t, tt.password, hash)

			// Verify that hashing the same password again produces a different hash (due to salt)
			hash2, err := HashPassword(tt.password)
			assert.NoError(t, err)
			assert.NotEqual(t, hash, hash2, "Hashing the same password twice should produce different hashes")
		})
	}
}

func TestComparePasswords(t *testing.T) {
	tests := []struct {
		name           string
		password       string
		compareWith    string
		expectedResult bool
	}{
		{
			name:           "Correct password",
			password:       "SecurePassword123",
			compareWith:    "SecurePassword123",
			expectedResult: true,
		},
		{
			name:           "Incorrect password",
			password:       "SecurePassword123",
			compareWith:    "WrongPassword",
			expectedResult: false,
		},
		{
			name:           "Case sensitive comparison",
			password:       "SecurePassword123",
			compareWith:    "securepassword123",
			expectedResult: false,
		},
		{
			name:           "Empty compare password",
			password:       "SecurePassword123",
			compareWith:    "",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Hash the original password
			hash, err := HashPassword(tt.password)
			assert.NoError(t, err)

			// Compare with the test password
			result := ComparePasswords(hash, []byte(tt.compareWith))
			assert.Equal(t, tt.expectedResult, result)

			// Additional test for invalid hash
			if tt.expectedResult {
				// Test with invalid hash
				invalidHash := "invalid-hash-format"
				result = ComparePasswords(invalidHash, []byte(tt.compareWith))
				assert.False(t, result, "Invalid hash format should return false")
			}
		})
	}
}
