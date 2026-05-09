package res

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKeyFromBase32(t *testing.T) {
	pass := "JBSWY3DP"
	assert.Equal(t, &Key{
		raw:    []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
		armour: pass,
	}, NewKey([]byte{}, pass), "Armour is working")
}

func TestNewKeyFromSecret(t *testing.T) {
	bytes := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
	assert.Equal(t, &Key{
		raw:    bytes,
		armour: "JBSWY3DP",
	}, NewKey(bytes), "Bytes is working")
}

func TestNewAccount(t *testing.T) {

	pass := "JBSWY3DP"
	key := NewKey([]byte{}, pass)
	content := `

	email: test@test.com
	armour: JBSWY3DP
	digits: 6
	`
	assert.Equal(t, MfaAccount{
		Name:   "name",
		key:    key,
		Email:  "test@test.com",
		digits: 6,
	},
		*NewAccountFromFile("name", content),
		"reading  Mfa account")

}

func TestNewAccountFromFile(t *testing.T) {
	name := "test/key"
	pass := "JBSWY3DP"
	key := NewKey([]byte{}, pass)
	content := `
email: test@test.com
armour: JBSWY3DP
digits: 8
`
	mkdir(res.dir, name)
	WriteFile(Join(name, ".key"), content)

	m := readKeys()
	keyfile := m[name]
	assert.Equal(t, MfaAccount{
		Name:   "test/key",
		key:    key,
		Email:  "test@test.com",
		digits: 8,
	},
		*NewAccountFromFile(name, keyfile.Content),
		"reading  Mfa account")

	check(DeleteFile(name + "/.key"))
	check(DeleteFile(name))
	check(DeleteFile(strings.Split(name, "/")[0]))

}
