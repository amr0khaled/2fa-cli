package main

import (
	"encoding/base32"
<<<<<<< Updated upstream
	"encoding/binary"
	"encoding/hex"
	"fmt"

	// "math"
	"os"
	// "strconv"

	//"time"

	c "amr.0x/totp-gen/internal/crypt"
	o "amr.0x/totp-gen/internal/totp"
=======
	"encoding/hex"
	"fmt"

	"os"

	"amr.0x/2fa-cli/internal/crypt"
	"amr.0x/2fa-cli/internal/totp"
>>>>>>> Stashed changes
)

func hexStr2Bytes(hexStr string) ([]byte, error) {
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	return hex.DecodeString(hexStr)
}

func main() {
<<<<<<< Updated upstream
	// key := crypt.Gen()
	// key, err := hexStr2Bytes("3132333435363738393031323334353637383930" + "313233343536373839303132")
	// if err != nil {
	//      fmt.Println("Error cipher")
	//      os.Exit(-1)
	// }
	// t0 := int64(0)
	// t := (int64(1111111109) - t0) / 30

	// bT := strconv.FormatInt(t, 16)
	// // binary.LittleEndian.PutUint64(bT, uint64(t))
	// for len(bT) < 16 {
	//      bT = "0" + bT
	// }
	// fmt.Println("Time", t)
	// fmt.Println("Time hex", bT)
	// thex, err := hexStr2Bytes(bT)
	// if err != nil {
	//      fmt.Println("Error DecodeString thex", err)
	//      os.Exit(-1)
	// }
	// hmac := crypt.R(thex, key)
	// fmt.Println("hmac", string(hmac), hex.EncodeToString(hmac))
	// fmt.Println("len hmac", len(hmac))

	// sbits := dt(hmac)
	// fmt.Println("sbits", sbits)
	// // sdata, err := hex.DecodeString(sbits)
	// // if err != nil {
	// //   fmt.Println("Error DecodeString", err)
	// //   os.Exit(-1)
	// // }
	// num := sbits % uint32(math.Pow10(8))
	// result := strconv.Itoa(int(num))
	// fmt.Println("num", num)
	// for len(result) < 8 {
	//      result = "0" + result
	// }

	// fmt.Println("result", result)
	hash := c.NewSha1()
=======
	hash := crypt.NewSha1()
>>>>>>> Stashed changes
	secretBytes, _ := hash.GenSecret()
	secret := hex.EncodeToString([]byte(secretBytes))
	secretBase32 := base32.StdEncoding.EncodeToString(secretBytes)
	fmt.Println("base32", secretBase32)
	secretBase32 = "MQDVCSI5MQO4XBPJEB6S7LKCYMPWTQHF"
	fmt.Println("base32", secretBase32)
<<<<<<< Updated upstream
	secretHex, _ := base32.StdEncoding.DecodeString(secretBase32)
	secret = hex.EncodeToString(secretHex)

	key := "12345678901234567890"
	secretHex, _ = hex.DecodeString(key)
	secret = hex.EncodeToString([]byte(key))
	fmt.Println("secret", secret)

	var otp string
	err := o.GenerateOTP(&otp, &hash, secret, 8)
	if err != nil {
		//	fmt.Println("Error", err)
=======

	secretHex, _ := base32.StdEncoding.DecodeString(secretBase32)
	secret = hex.EncodeToString(secretHex)
	fmt.Println("secret", secret)

	var otp string
	err := totp.GenerateOTP(&otp, &hash, secret, 6)
	if err != nil {
		fmt.Println("Error", err)
>>>>>>> Stashed changes
		os.Exit(-1)
	}
	fmt.Println("Geenrated", otp)

<<<<<<< Updated upstream
	if ok, _ := o.Validate(otp, secret); ok {
		//	fmt.Println("Valid")
	} else {
		//	fmt.Println("Not Valid")
	}

=======
	if ok, _ := totp.Validate(otp, secret); ok {
		fmt.Println("Valid")
	} else {
		fmt.Println("Not Valid")
	}
>>>>>>> Stashed changes
}
