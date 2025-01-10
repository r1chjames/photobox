package auth

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestCreateHash(t *testing.T) {
	hashRX := TokenFormatRegEx()
	hash1, err := CreateHash("pa$$word", DefaultArgon2idHash())
	assert.NoError(t, err)

	assert.True(t, hashRX.MatchString(hash1), "Hash not in correct format")

	hash2, err := CreateHash("pa$$word", DefaultArgon2idHash())
	assert.NoError(t, err)

	assert.NotEqual(t, hash1, hash2, "hashes must be unique")
}

func TestComparePasswordAndHash(t *testing.T) {
	hash, err := CreateHash("pa$$word", DefaultArgon2idHash())
	assert.NoError(t, err)

	match, err := ComparePasswordAndHash("pa$$word", hash)
	assert.NoError(t, err)

	assert.True(t, match, "expected password and hash to match")

	match, err = ComparePasswordAndHash("otherPa$$word", hash)
	assert.NoError(t, err)

	assert.False(t, match, "expected password and hash to not match")
}

func TestDecodeHash(t *testing.T) {
	hash, err := CreateHash("pa$$word", DefaultArgon2idHash())
	assert.NoError(t, err)

	params, _, _, err := DecodeHash(hash)
	assert.NoError(t, err)
	if *params != *DefaultArgon2idHash() {
		t.Fatalf("expected %#v got %#v", *DefaultArgon2idHash(), *params)
	}
}

func TestStrictDecoding(t *testing.T) {
	// "bug" valid hash: $argon2id$v=19$m=65536,t=1,p=2$UDk0zEuIzbt0x3bwkf8Bgw$ihSfHWUJpTgDvNWiojrgcN4E0pJdUVmqCEdRZesx9tE
	ok, _, err := CheckHash("bug", "$argon2id$v=19$m=65536,t=1,p=2$UDk0zEuIzbt0x3bwkf8Bgw$ihSfHWUJpTgDvNWiojrgcN4E0pJdUVmqCEdRZesx9tE")
	assert.NoError(t, err)
	assert.True(t, ok, "expected password to match")

	// changed one last character of the hash
	ok, _, err = CheckHash("bug", "$argon2id$v=19$m=65536,t=1,p=2$UDk0zEuIzbt0x3bwkf8Bgw$ihSfHWUJpTgDvNWiojrgcN4E0pJdUVmqCEdRZesx9tF")
	assert.Error(t, err, "Hash validation should fail")

	if ok {
		t.Fatal("Hash validation should fail")
	}
}

func TestVariant(t *testing.T) {
	// Hash contains wrong variant
	_, _, err := CheckHash("pa$$word", "$argon2i$v=19$m=65536,t=1,p=2$mFe3kxhovyEByvwnUtr0ow$nU9AqnoPfzMOQhCHa9BDrQ+4bSfj69jgtvGu/2McCxU")
	if err != ErrIncompatibleVariant {
		t.Fatalf("expected error %s", ErrIncompatibleVariant)
	}
}

func TokenFormatRegEx() *regexp.Regexp {
	hashRX, err := regexp.Compile(`^\$argon2id\$v=19\$m=65536,t=1,p=[0-9]{1,4}\$[A-Za-z0-9+/]{22}\$[A-Za-z0-9+/]{43}$`)
	if err != nil {
		return nil
	}
	return hashRX
}
