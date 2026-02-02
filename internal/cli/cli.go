package cli

import (
	"encoding/base32"
	"encoding/hex"
	"errors"
	"fmt"

	"amr.0x/2fa-cli/internal/crypt"
	"amr.0x/2fa-cli/internal/totp"
)

type Context struct {
	Debug bool
}
type Cli struct {
	Hex    bool   `help:"if you want secret in hex value"`
	Base32 bool   `help:"if you want secret in base32 value default:true"`
	Secret string `help:"input a secret"`
	Digit  int    `help:"number of digits in OTP" default:"6"`
	Hash   string `help:"use one of those hash methods \"sha1\" \"sha256\" \"sha512\"" default:"sha256"`

	Get   GetEvent   `cmd:"" help:"Get OTP of stored secret"`
	New   NewEvent   `cmd:"" help:"Create a new secret"`
	Store StoreEvent `cmd:"" help:"Store a pre-generated secret"`

	hashMethod crypt.HashMethod
}

type Event interface {
	Run(*Cli) error
}
type GetEvent struct {
}
type NewEvent struct {
}
type StoreEvent struct {
}

func (event *GetEvent) Run(c *Cli) error {
	var otp string
	var secret string
	if c.Hex {
		secret = c.Secret
	} else {
		secretBytes, err := base32.StdEncoding.DecodeString(c.Secret)
		if err != nil {
			return err
		}
		secret = hex.EncodeToString(secretBytes)
	}
	err := totp.GenerateOTP(&otp, c.hashMethod, secret, c.Digit)
	if err != nil {
		return err
	}
	fmt.Println("OTP:", otp)
	return nil
}
func (event *NewEvent) Run(c *Cli) error {
	secretBytes, err := c.hashMethod.GenSecret()
	if err != nil {
		return err
	}
	var secret string
	if c.Hex {
		secret = hex.EncodeToString([]byte(secretBytes))
	} else {
		secret = base32.StdEncoding.EncodeToString(secretBytes)
	}
	fmt.Println("Secret:", secret)
	return nil
}
func (event *StoreEvent) Run(*Cli) error {

	return nil
}
func (c *Cli) Validate() error {
	if v, ok := crypt.Hashes[c.Hash]; !ok {
		return errors.New("The option is not valid. Choose one of (sha1, sha256, sha512)")
	} else {
		c.hashMethod = v
	}
	return nil
}
