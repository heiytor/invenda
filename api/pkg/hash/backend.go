package hash

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type backend struct{}

var Backend Hash = &backend{}

func (b *backend) New(str string, opt *Options) string {
	salt, err := newSalt(opt.SaltLength)
	if err != nil {
		return ""
	}

	hash := argon2.IDKey([]byte(str), salt, opt.Iterations, opt.Memory, opt.Parallelism, opt.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, opt.Memory, opt.Iterations, opt.Parallelism, b64Salt, b64Hash)
}

func (b *backend) Compare(str, hash string) bool {
	opt, salt, decodedHash, err := decode(hash)
	if err != nil {
		return false
	}

	otherHash := argon2.IDKey([]byte(str), salt, opt.Iterations, opt.Memory, opt.Parallelism, opt.KeyLength)

	return subtle.ConstantTimeCompare(decodedHash, otherHash) == 1
}
