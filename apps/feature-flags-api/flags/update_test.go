package flags

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateFlag(t *testing.T) {
	// Setup test flags in memory
	flagsCache["test"] = map[string]bool{
		"test_flag": false,
		"another_flag": true,
	}
	
	tests := []struct {
		name        string
		env         string
		key         string
		enabled     bool
		expectError bool
	}{
		{
			name:        "successful update to true",
			env:         "test",
			key:         "test_flag",
			enabled:     true,
			expectError: false,
		},
		{
			name:        "successful update to false",
			env:         "test", 
			key:         "another_flag",
			enabled:     false,
			expectError: false,
		},
		{
			name:        "flag not found",
			env:         "test",
			key:         "nonexistent",
			enabled:     true,
			expectError: true,
		},
		{
			name:        "env not found",
			env:         "nonexistent",
			key:         "test_flag",
			enabled:     true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateFlag(tt.env, tt.key, tt.enabled)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				
				// Verify the flag was actually updated
				value, exists, getErr := GetSingleFlag(tt.env, tt.key)
				require.NoError(t, getErr)
				require.True(t, exists)
				assert.Equal(t, tt.enabled, value)
			}
		})
	}
}
