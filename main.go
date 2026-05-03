package main

import (
	// c "amr.0x/2fa-cli/internal/cli"
	// "github.com/alecthomas/kong"
	"fmt"
	"os"

	// "os/exec"

	"amr.0x/2fa-cli/internal/res"
	"github.com/ProtonMail/gopenpgp/v3/constants"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/ProtonMail/gopenpgp/v3/profile"
)

func get(message string, value *string) {
	fmt.Print(message)
	fmt.Scanln(value)
}

const KEY_FILE = ".key"

func createKey() *crypto.Key {
	var name, email string
	get("Name: ", &name)
	get("Email: ", &email)

	pgp := crypto.PGPWithProfile(profile.Default())

	keyGen := pgp.KeyGeneration().AddUserId(name, email).New()
	key, err := keyGen.GenerateKeyWithSecurity(constants.HighSecurity)
	if err != nil {
		fmt.Println("Key Gen Error", err)
		os.Exit(-2)
		return nil
	}
	armour, err := key.Armor()
	if err != nil {
		fmt.Println("Key Gen Error", err)
		os.Exit(-2)
		return nil
	}
	if ok, err := res.Exists(KEY_FILE); err != nil {
		panic(err)
	} else if !ok {
		err := res.CreateFile(KEY_FILE)
		if err != nil {
			fmt.Println("Create key file error", err)
			os.Exit(-2)
			return nil
		}
	}
	err = res.WriteFile(KEY_FILE, armour)
	if err != nil {
		fmt.Println("Write key error", err)
		os.Exit(-2)
		return nil
	}
	return key
}
func getKey() (*crypto.Key, error) {
	var err error
	if exists, err := res.Exists(KEY_FILE); err != nil {
		return nil, err
	} else if !exists {
		return createKey(), nil
	}
	armoured, err := res.ReadFile(KEY_FILE)
	if err != nil {
		return nil, err
	}
	key, err := crypto.NewKeyFromArmored(*armoured)
	if err != nil {
		return nil, err
	}
	return key, nil
}
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// var cli c.Cli
	// ctx := kong.Parse(&cli,
	// 	kong.Name("2fa-cli"))
	// err := ctx.Run(&cli)
	// ctx.FatalIfErrorf(err)
	key, err := getKey()
	check(err)
	armour, err := key.Armor()
	check(err)
	fmt.Println("key", armour)
}
