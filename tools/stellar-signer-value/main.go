package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon/base"
	"github.com/stellar/go/support/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	exitCode := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(exitCode)
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	cmd := &cobra.Command{
		Use:   "stellar-signer-value",
		Short: "Get a sum of all the balances that a signer can sign for. Uses Horizon.",
	}
	cmd.SetArgs(args)
	cmd.SetOutput(stderr)

	horizonURL := horizonclient.DefaultPublicNetClient.HorizonURL
	cmd.Flags().StringVarP(&horizonURL, "horizon-url", "", horizonURL, "Horizon URL used for looking up account balances")
	cmd.MarkFlagRequired("horizon-url")

	signer := ""
	cmd.Flags().StringVarP(&signer, "signer", "s", signer, "Signer to get a sum of all balances")
	cmd.MarkFlagRequired("signer")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client := horizonclient.Client{HorizonURL: horizonURL}
		req := horizonclient.AccountsRequest{Signer: signer}

		accountCount := 0
		balances := map[base.Asset]int64{}

		accounts, err := client.Accounts(req)
		for err == nil && len(accounts.Embedded.Records) > 0 {
			for _, a := range accounts.Embedded.Records {
				for _, b := range a.Balances {
					floatAmount, err := strconv.ParseFloat(b.Balance, 64)
					if err != nil {
						return errors.Wrapf(err, "parsing balance %s for account %s", b.Balance, a.ID)
					}
					amount := int64(floatAmount * 10_000_000)
					balances[b.Asset] += amount
				}
				accountCount++
			}

			accounts, err = client.NextAccountsPage(accounts)
		}
		if err != nil {
			fmt.Println("Error:", err)
		}

		fmt.Printf("Accounts = %d\n", accountCount)
		for asset, total := range balances {
			whole := total / 10_000_000
			integral := total % 10_000_000
			wholeStr := message.NewPrinter(language.English).Sprint(whole)
			fmt.Printf("%s = %s.%d\n", assetName(asset), wholeStr, integral)
		}

		return nil
	}

	err := cmd.Execute()
	if err != nil {
		return 1
	}
	return 0
}

func assetName(asset base.Asset) string {
	if asset.Type == "native" {
		return "native"
	}
	return asset.Code + ":" + asset.Issuer
}
