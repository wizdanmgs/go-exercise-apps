package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/spf13/cobra"
)

const (
	lower   = "abcdefghijklmnopqrstuvwxyz"
	upper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits  = "0123456789"
	symbols = "!@#$%^&*()-_=+[]{}<>?/|"
)

var (
	length     int
	useSymbols bool
)

var rootCmd = &cobra.Command{
	Use:   "passwordgen",
	Short: "A secure CLI password generator",
	RunE: func(cmd *cobra.Command, args []string) error {
		if length <= 0 {
			return fmt.Errorf("length must be greater than 0")
		}

		charset := buildCharset(useSymbols)

		password, err := generatePassword(length, charset)
		if err != nil {
			return err
		}

		fmt.Println(password)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&length, "length", "l", 16, "length of the password")
	rootCmd.Flags().BoolVarP(&useSymbols, "symbols", "s", true, "include symbols")
}

func buildCharset(includeSymbols bool) string {
	charset := lower + upper + digits
	if includeSymbols {
		charset += symbols
	}
	return charset
}

func generatePassword(length int, charset string) (string, error) {
	result := make([]byte, length)
	max := big.NewInt(int64(len(charset)))

	for i := range result {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		result[i] = charset[n.Int64()]
	}

	return string(result), nil
}
