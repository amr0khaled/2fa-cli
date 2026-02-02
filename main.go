package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"math"
	"os"
	"strconv"

	"amr.0x/2fa-cli/internal/crypt"
)

func hashKey(key string) []byte {
	data, err := hex.DecodeString(key)
	if err != nil {
		fmt.Println("Error hashKey")
		os.Exit(-1)
	}
	return data
}
func dt(hmac []byte) uint32 {
	lastIndex := len(hmac) - 1
	lowerByte := hmac[lastIndex]
	lowerNibble := lowerByte & 0xf
	offset := lowerNibble
	code := binary.BigEndian.Uint32(hmac[offset : offset+4])
	return code & 0x7FFFFFFF
}
func hexStr2Bytes(hexStr string) ([]byte, error) {
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	return hex.DecodeString(hexStr)
}

func main() {
	key := crypt.Gen()
	key, err := hexStr2Bytes("3132333435363738393031323334353637383930" + "313233343536373839303132")
	if err != nil {
		fmt.Println("Error cipher")
		os.Exit(-1)
	}
	t0 := int64(0)
	t := (int64(1111111109) - t0) / 30

	bT := strconv.FormatInt(t, 16)
	// binary.LittleEndian.PutUint64(bT, uint64(t))
	for len(bT) < 16 {
		bT = "0" + bT
	}
	fmt.Println("Time", t)
	fmt.Println("Time hex", bT)
	thex, err := hexStr2Bytes(bT)
	if err != nil {
		fmt.Println("Error DecodeString thex", err)
		os.Exit(-1)
	}
	hmac := crypt.R(thex, key)
	fmt.Println("hmac", string(hmac), hex.EncodeToString(hmac))
	fmt.Println("len hmac", len(hmac))

	sbits := dt(hmac)
	fmt.Println("sbits", sbits)
	// sdata, err := hex.DecodeString(sbits)
	// if err != nil {
	// 	fmt.Println("Error DecodeString", err)
	// 	os.Exit(-1)
	// }
	num := sbits % uint32(math.Pow10(8))
	result := strconv.Itoa(int(num))
	fmt.Println("num", num)
	for len(result) < 8 {
		result = "0" + result
	}

	fmt.Println("result", result)
}
