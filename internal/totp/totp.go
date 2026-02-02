package totp

import (
	"crypto/subtle"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"
	"strconv"
	"time"

	c "amr.0x/totp-gen/internal/crypt"
)

const T0 = 0
const Interval = 30

func hexStr2Bytes(hexStr string) ([]byte, error) {
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	return hex.DecodeString(hexStr)
}

func truncate(hash []byte) uint32 {
	lastIndex := len(hash) - 1
	lowerByte := hash[lastIndex]
	lowerNibble := lowerByte & 0xf
	offset := lowerNibble
	code := binary.BigEndian.Uint32(hash[offset : offset+4])
	return code & 0x7FFFFFFF
}

func prepareOTP(num uint32, digit int) string {
	fullPassword := num % uint32(math.Pow10(digit))
	otp := strconv.Itoa(int(fullPassword))
	for len(otp) < digit {
		otp = "0" + otp
	}
	return otp
}

var lenToHash = map[int]c.HashMethod{
	40:  &c.Sha1{},
	64:  &c.Sha256{},
	128: &c.Sha512{},
}

func equalityCheck(inputOTP string, realOTP string) bool {
	inputBytes := []byte(inputOTP)
	realBytes := []byte(realOTP)
	if subtle.ConstantTimeCompare(inputBytes, realBytes) == 1 {
		return true
	}
	return false
}

func Validate(otp string, secret string) (bool, error) {
	digit := len(otp)
	if digit < 6 {
		err := errors.New("OTP length shouldn't be less than 6 digits")
		return false, err
	}
	hash, ok := lenToHash[len(secret)]
	if !ok {
		err := errors.New("Secret is not valid; pass it as hex in form of string")
		return false, err
	}
	var generatedOTP string
	err := GenerateOTP(&generatedOTP, hash, secret, digit)
	if err != nil {
		return false, err
	}
	return equalityCheck(otp, string(generatedOTP)), nil
}
func GenerateOTP(otp *string, hash c.HashMethod, secret string, digit int) error {
	secretHex, err := hexStr2Bytes(secret)
	if err != nil {
		return err
	}
	date := time.Date(1970, time.January, 1, 0, 0, 59, 0, time.UTC)

	currentTime := date.Unix()
	time := (currentTime - T0) / Interval
	timeStr := strconv.FormatInt(time, 16)
	for len(timeStr) < 16 {
		timeStr = "0" + timeStr
	}
	timeHex, err := hexStr2Bytes(timeStr)
	if err != nil {
		return err
	}
	cipher := hash.Hmac(timeHex, secretHex)
	num := truncate(cipher)
	*otp = prepareOTP(num, digit)
	return nil
}
