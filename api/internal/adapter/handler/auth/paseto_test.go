package auth

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/r1chjames/photobox/api/internal/core/domain"
	"testing"
	"time"
)

var tokenService *PasetoMaker
var testUser *domain.User

func setup() {
	tokenService, _ = New("123456780abcdefghijklmnopqrstuvw", time.Minute)
	testUser = &domain.User{
		Username: "test1",
		Email:    "test1@example.com",
	}
}

func TestCreateInstanceWithInvalidSymmetricKey(t *testing.T) {
	_, err := New("123456780invalid", time.Minute)
	assert.Error(t, err, "SymmetricKey too short should be")
}

func TestCreateAndVerifyToken(t *testing.T) {
	setup()
	token, err := tokenService.CreateToken(testUser)
	assert.NoError(t, err)

	verifiedToken, err := tokenService.VerifyToken(token)
	assert.NoError(t, err)

	assert.Equal(t, testUser.Username, verifiedToken.Username)
}

func TestCreateAndVerifyInvalidToken(t *testing.T) {
	setup()
	token := "v2.local.dF3zt__uUQQx3b8Nnz51cITVQpB4ztEfQLyau4BPFLxqRp1JGws-invalid-invalid-invalid"

	_, err := tokenService.VerifyToken(token)
	assert.ErrorContains(t, err, "invalid token authentication")
}

func TestCreateAndVerifyExpiredToken(t *testing.T) {
	setup()
	token := "v2.local.dF3zt__uUQQx3b8Nnz51cITVQpB4ztEfQLyau4BPFLxqRp1JGws83NxA6kfHZP6o4Jx-LX7fZyQwCZt6PhKyc-WPd48e9ouNSlTMrx-IlQXC4VHKe82SJ7B1ZjvvPWRCtulCi-tfdSWDpxrIoZfvxJ28ULt19RPYIx-2AQq9_9QrPmEmOBp2C3qTi_1XfYK41sc_vVGQ09t9SgeMbS1Qtiz19UxuU9eQj50pZ1H9VsLBCELeeHwvUr7M_rAefdNxmaCuZQ3zBA.bnVsbA"

	_, err := tokenService.VerifyToken(token)
	assert.ErrorContains(t, err, "token has expired")
}
