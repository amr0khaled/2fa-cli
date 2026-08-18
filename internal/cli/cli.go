// Package cli
package cli

import (
	"bytes"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"

	"amr.0x/2fa-cli/internal/crypt"
	"amr.0x/2fa-cli/internal/res"
	"amr.0x/2fa-cli/internal/totp"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
)

func check(err error) {
	if err != nil {
		panic("Encountered an issue: " + err.Error())
	}
}

var Root = &cobra.Command{
	Use:   "2fa",
	Short: "A 2FA Authenticator",
}

type CLI struct {
	Secret     *res.Key
	HashMethod crypt.HashMethod
	IsHex      bool
	Digits     int8
	Otp        string
}

var cli *CLI = &CLI{}

var initCmd = &cobra.Command{}

var getCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Get OTP of stored secret",
	Aliases: []string{
		"g",
	},
	ValidArgs: []cobra.Completion{
		"name",
	},
	Example: "  2fa get email@google",
	// 2fa get name/test
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("not enough arguments")
		}
		name := args[0]
		account := res.Accounts[name]
		if account == nil {
			return fmt.Errorf("%s is not found in store", name)
		}
		digit := account.GetDigits()
		digitFlag, err := cmd.Flags().GetInt8("digits")
		if err != nil {
			return err
		}
		if digitFlag != digit {
			digit = digitFlag
		}
		secret := account.GetKey()
		cli = &CLI{
			Secret:     secret,
			IsHex:      false,
			Digits:     digit,
			HashMethod: crypt.NewSha256(),
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var otp string
		err := totp.GenerateOTP(&otp, cli.HashMethod, cli.Secret.GetRaw(), int(cli.Digits))
		if err != nil {
			return err
		}
		cli.Otp = otp
		cpy, err := cmd.Flags().GetBool("copy")
		if err != nil {
			return err
		}
		if cpy {
			if os.Getenv("WAYLAND_DISPLAY") != "" {
				path, err := exec.LookPath("wl-copy")
				if err != nil {
					return fmt.Errorf("wl-clipboard utility not found: install 'wl-clipboard' via your package manager")
				}
				com := exec.Command(path)
				com.Stdin = bytes.NewBufferString(cli.Otp)
				return com.Run()
			} else {
				if err := clipboard.Init(); err != nil {
					return err
				}
				_ = clipboard.Write(clipboard.FmtText, []byte(cli.Otp))
			}
			fmt.Println("Copied OTP.")
		} else {
			fmt.Printf("OTP: %s", cli.Otp)
		}
		return nil
	},
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate new secret",
	// 2fa new
	PreRunE: func(cmd *cobra.Command, args []string) error {
		cli = &CLI{}
		isHex, err := cmd.Flags().GetBool("hex")
		if err != nil {
			return err
		}
		cli.IsHex = isHex
		hash, err := cmd.Flags().GetString("hash")
		if err != nil {
			return err
		}
		cli.HashMethod = crypt.Hashes[hash]
		if cli.HashMethod == nil {
			cli.HashMethod = crypt.NewSha256()
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		secretBytes, err := cli.HashMethod.GenSecret()
		if err != nil {
			return err
		}
		var armour string
		if cli.IsHex {
			armour = hex.EncodeToString([]byte(secretBytes))
		} else {
			armour = base32.StdEncoding.EncodeToString(secretBytes)
		}
		fmt.Println("New key:\n\t", armour)
		return nil
	},
}

var storeCmd = &cobra.Command{
	Use:   "store <name> <secret>",
	Short: "Store new secret of an account",
	Aliases: []string{
		"s",
	},
	// 2fa store name/test [secret]
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("not enough arguments")
		}
		secret := args[1]
		if secret == "" {
			return fmt.Errorf("error: Secret is not provided")
		}
		isHex, err := cmd.Flags().GetBool("hex")
		if err != nil {
			return err
		}
		hash, err := cmd.Flags().GetString("hash")
		if err != nil {
			return err
		}
		cli = &CLI{}
		if hash != "" {
			switch hash {
			case "sha1":
				cli.HashMethod = crypt.NewSha1()
			case "sha256":
				cli.HashMethod = crypt.NewSha256()
			case "sha512":
				cli.HashMethod = crypt.NewSha512()
			default:
				return fmt.Errorf("error: %s is not suitable\n    use one of those hash methods \"sha1\" \"sha256\" \"sha512\"", hash)
			}
		} else {
			cli.HashMethod = crypt.NewSha256()
		}
		if err != nil {
			return err
		}
		if isHex {
			secretBytes, err := hex.DecodeString(secret)
			if err != nil {
				return err
			}
			cli.Secret = res.NewKey(secretBytes, "")
		} else {
			cli.Secret = res.NewKey([]byte{}, secret)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		digits, err := cmd.Flags().GetInt8("digits")
		if err != nil {
			return err
		}
		account := res.NewMFaAccount(name, cli.Secret, digits)
		return account.Store()
	},
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("completion called")
	},
}

func init() {
	getCmd.Flags().BoolP("copy", "c", false, "directly copy to clipboard")
	Root.PersistentFlags().Int8P("digits", "d", 6, "define number of digits needed in OTP")
	Root.PersistentFlags().String("hash", "sha256", "define the hash method (sha1, sha256, sha512)")
	Root.PersistentFlags().Bool("hex", false, "define the input secret are in hex, otherwise it is in base32")
	Root.AddCommand(getCmd)
	Root.AddCommand(genCmd)
	Root.AddCommand(storeCmd)
	Root.AddCommand(completionCmd)
}
