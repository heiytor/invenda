package hash_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/heiytor/invenda/api/pkg/hash"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/argon2"
)

func TestNew(t *testing.T) {
	cases := []struct {
		description string
		str         string
		opts        *hash.Options
	}{
		{
			description: "succeeds to hash the str",
			str:         "secret",
			opts:        hash.DefaultOptions,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			hash := hash.New(tc.str, tc.opts)
			require.NotZero(t, hash)

			parts := strings.Split(hash, "$")
			require.Equal(t, len(parts), 6)
			require.Equal(t, parts[1], "argon2id")
			require.Equal(t, parts[2], fmt.Sprintf("v=%d", argon2.Version))
			require.Equal(t, parts[3], fmt.Sprintf("m=%d,t=%d,p=%d", tc.opts.Memory, tc.opts.Iterations, tc.opts.Parallelism))
			require.NotZero(t, parts[4])
			require.NotZero(t, parts[5])
		})
	}
}

func TestCompare(t *testing.T) {
	cases := []struct {
		description string
		str         string
		hash        string
		expected    bool
	}{
		{
			description: "fails when str does not match with hash due to str",
			str:         "invalid",
			hash:        "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
			expected:    false,
		},
		{
			description: "fails when str does not match with hash due to params",
			str:         "secret",
			hash:        "$argon2id$v=20$m=65536,t=1,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg73fad",
			expected:    false,
		},
		{
			description: "fails when str does not match with hash due to salt",
			str:         "secret",
			hash:        "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo800wG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
			expected:    false,
		},
		{
			description: "fails when str does not match with hash due to hash",
			str:         "secret",
			hash:        "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg73fad",
			expected:    false,
		},
		{
			description: "succeeds when str matches with hash",
			str:         "secret",
			hash:        "$argon2id$v=19$m=65536,t=3,p=2$ODZA1eh/sfvUo6OMwG0HbA$lkg+QLLrFzYZkMcvZhoDtSgaHoAHKVnu2Ztexg9PUo0",
			expected:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			require.Equal(t, tc.expected, hash.Compare(tc.str, tc.hash))
		})
	}
}
