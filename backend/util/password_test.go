package util

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/guncv/Poll-Voting-Website/backend/constant"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func randomString(n int) string {
	var sb strings.Builder
	k := len(constant.Alphabet)

	for i := 0; i < n; i++ {
		c := constant.Alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func TestPassword(t *testing.T) {
	password := randomString(8)

	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	wrongPassword := randomString(8)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
