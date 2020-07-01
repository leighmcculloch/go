package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stellar/go/keypair"
	"golang.org/x/crypto/hkdf"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "stellar-hmac-gen [command]",
		Short: "Generate a Stellar key to be used as a signer for the account given a shared secret key and a nonce.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	hkdfCmd := &cobra.Command{
		Use:   "hkdf",
		Short: "Use HKDF",
	}
	hkdfCmd.Flags().String("secret", "", "The shared secret key")
	hkdfCmd.MarkFlagRequired("secret")
	hkdfCmd.Flags().String("account", "", "A Stellar account the key is being used with")
	hkdfCmd.MarkFlagRequired("account")
	hkdfCmd.Flags().String("nonce", "", "A nonce used to modify the generated")
	hkdfCmd.MarkFlagRequired("nonce")
	hkdfCmd.Run = func(cmd *cobra.Command, args []string) {
		secret, _ := cmd.Flags().GetString("secret")
		account, _ := cmd.Flags().GetString("account")
		nonce, _ := cmd.Flags().GetString("nonce")
		if secret == "" || account == "" || nonce == "" {
			cmd.Help()
			return
		}
		key, err := genHkdf(secret, account, nonce)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Key:", key.Seed())
	}
	rootCmd.AddCommand(hkdfCmd)

	hmacCmd := &cobra.Command{
		Use:   "hmac",
		Short: "Use HMAC-SHA256",
	}
	hmacCmd.Flags().String("secret", "", "The shared secret key")
	hmacCmd.MarkFlagRequired("secret")
	hmacCmd.Flags().String("account", "", "A Stellar account the key is being used with")
	hmacCmd.MarkFlagRequired("account")
	hmacCmd.Flags().String("nonce", "", "A nonce used to modify the generated")
	hmacCmd.MarkFlagRequired("nonce")
	hmacCmd.Run = func(cmd *cobra.Command, args []string) {
		secret, _ := cmd.Flags().GetString("secret")
		account, _ := cmd.Flags().GetString("account")
		nonce, _ := cmd.Flags().GetString("nonce")
		if secret == "" || account == "" || nonce == "" {
			cmd.Help()
			return
		}
		key, err := genHmac(secret, account, nonce)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Key:", key.Seed())
	}
	rootCmd.AddCommand(hmacCmd)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func genHmac(secret, account, nonce string) (*keypair.Full, error) {
	message := account + "," + nonce
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(message))
	keyBytes := hash.Sum(nil)
	keyBytes32 := [32]byte{}
	count := copy(keyBytes32[:], keyBytes)
	if count != 32 {
		return nil, fmt.Errorf("Error: unexpected key bytes length must be 32 got %d", count)
	}
	return keypair.FromRawSeed(keyBytes32)
}

func genHkdf(secret, account, nonce string) (*keypair.Full, error) {
	// TODO: The salt (nonce) should probably be 32 bytes.
	hash := hkdf.New(sha256.New, []byte(secret), []byte(nonce), []byte(account))
	keyBytes32 := [32]byte{}
	count, err := hash.Read(keyBytes32[:])
	if err != nil {
		return nil, err
	}
	if count != 32 {
		return nil, fmt.Errorf("Error: unexpected key bytes length must be 32 got %d", count)
	}
	return keypair.FromRawSeed(keyBytes32)
}
