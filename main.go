package main

import (
	// "encoding/base32"
	// "encoding/hex"
	// "fmt"

	// "os"

	c "amr.0x/2fa-cli/internal/cli"
	"github.com/alecthomas/kong"
	// "amr.0x/2fa-cli/internal/crypt"
	// "amr.0x/2fa-cli/internal/totp"
)

func main() {
	var cli c.Cli
	ctx := kong.Parse(&cli,
		kong.Name("2fa-cli"))
	err := ctx.Run(&cli)
	ctx.FatalIfErrorf(err)
	// config := *cli.GetConf()
	// fmt.Println("config", config)
	// hash := crypt.NewSha1()
	// secretBytes, _ := hash.GenSecret()
	// secret := hex.EncodeToString([]byte(secretBytes))
	// secretBase32 := base32.StdEncoding.EncodeToString(secretBytes)
	// fmt.Println("base32", secretBase32)
	// secretBase32 = "MQDVCSI5MQO4XBPJEB6S7LKCYMPWTQHF"
	// fmt.Println("base32", secretBase32)

	// secretHex, _ := base32.StdEncoding.DecodeString(secretBase32)
	// secret = hex.EncodeToString(secretHex)
	// fmt.Println("secret", secret)

	// var otp string
	// err := totp.GenerateOTP(&otp, &hash, secret, 6)
	// if err != nil {
	// 	fmt.Println("Error", err)
	// 	os.Exit(-1)
	// }
	// fmt.Println("Geenrated", otp)

	// if ok, _ := totp.Validate(otp, secret); ok {
	// 	fmt.Println("Valid")
	// } else {
	// 	fmt.Println("Not Valid")
	// }
}
