package cli

import (
	"fmt"

	"github.com/DioraFund/diora-go-blockchain/internal/keystore"
	"github.com/urfave/cli/v2"
)

// getWalletCommand returns the wallet command
func (c *CLI) getWalletCommand() *cli.Command {
	return &cli.Command{
		Name:  "wallet",
		Usage: "Manage Diora wallets and accounts",
		Subcommands: []*cli.Command{
			c.getWalletCreateCommand(),
			c.getWalletImportCommand(),
			c.getWalletListCommand(),
			c.getWalletBalanceCommand(),
			c.getWalletExportCommand(),
			c.getWalletDeleteCommand(),
		},
	}
}

// getWalletCreateCommand returns the wallet create command
func (c *CLI) getWalletCreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new wallet",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Wallet name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				Usage:    "Wallet password",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"d"},
				Usage:   "Wallet directory path",
				Value:   c.config.Keystore.Path,
			},
		},
		Action: func(ctx *cli.Context) error {
			name := ctx.String("name")
			password := ctx.String("password")
			path := ctx.String("path")

			c.logger.Info("Creating wallet", "name", name, "path", path)

			ks := keystore.NewKeystore(path)
			account, err := ks.NewAccount(name, password)
			if err != nil {
				return fmt.Errorf("failed to create wallet: %w", err)
			}

			fmt.Printf("✅ Wallet created successfully!\n")
			fmt.Printf("📝 Name: %s\n", name)
			fmt.Printf("📍 Address: %s\n", account.Address.Hex())
			fmt.Printf("🔑 Public Key: %s\n", account.PublicKey)
			fmt.Printf("📂 Path: %s\n", path)

			return nil
		},
	}
}

// getWalletImportCommand returns the wallet import command
func (c *CLI) getWalletImportCommand() *cli.Command {
	return &cli.Command{
		Name:  "import",
		Usage: "Import wallet from private key",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "private-key",
				Aliases:  []string{"key", "k"},
				Usage:    "Private key to import",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Wallet name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				Usage:    "Wallet password",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"d"},
				Usage:   "Wallet directory path",
				Value:   c.config.Keystore.Path,
			},
		},
		Action: func(ctx *cli.Context) error {
			privateKey := ctx.String("private-key")
			name := ctx.String("name")
			password := ctx.String("password")
			path := ctx.String("path")

			c.logger.Info("Importing wallet", "name", name, "path", path)

			ks := keystore.NewKeystore(path)
			account, err := ks.ImportPrivateKey(privateKey, name, password)
			if err != nil {
				return fmt.Errorf("failed to import wallet: %w", err)
			}

			fmt.Printf("✅ Wallet imported successfully!\n")
			fmt.Printf("📝 Name: %s\n", name)
			fmt.Printf("📍 Address: %s\n", account.Address.Hex())
			fmt.Printf("🔑 Public Key: %s\n", account.PublicKey)
			fmt.Printf("📂 Path: %s\n", path)

			return nil
		},
	}
}

// getWalletListCommand returns the wallet list command
func (c *CLI) getWalletListCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all wallets",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"d"},
				Usage:   "Wallet directory path",
				Value:   c.config.Keystore.Path,
			},
		},
		Action: func(ctx *cli.Context) error {
			path := ctx.String("path")

			c.logger.Info("Listing wallets", "path", path)

			ks := keystore.NewKeystore(path)
			accounts, err := ks.ListAccounts()
			if err != nil {
				return fmt.Errorf("failed to list wallets: %w", err)
			}

			if len(accounts) == 0 {
				fmt.Printf("📭 No wallets found in %s\n", path)
				return nil
			}

			fmt.Printf("📋 Found %d wallet(s):\n\n", len(accounts))
			for i, account := range accounts {
				fmt.Printf("%d. 📝 %s\n", i+1, account.Name)
				fmt.Printf("   📍 Address: %s\n", account.Address.Hex())
				fmt.Printf("   🔑 Public Key: %s\n", account.PublicKey)
				fmt.Printf("   📅 Created: %s\n", account.CreatedAt.Format("2006-01-02 15:04:05"))
				fmt.Println()
			}

			return nil
		},
	}
}

// getWalletBalanceCommand returns the wallet balance command
func (c *CLI) getWalletBalanceCommand() *cli.Command {
	return &cli.Command{
		Name:  "balance",
		Usage: "Get wallet balance",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Aliases:  []string{"addr", "a"},
				Usage:    "Wallet address",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "token",
				Aliases: []string{"t"},
				Usage:   "Token contract address (default: native DIO)",
				Value:   "",
			},
		},
		Action: func(ctx *cli.Context) error {
			address := ctx.String("address")
			token := ctx.String("token")

			c.logger.Info("Getting wallet balance", "address", address, "token", token)

			// TODO: Implement balance fetching logic
			fmt.Printf("💰 Balance for %s\n", address)
			if token == "" {
				fmt.Printf("🪙 Native DIO: 0.000000 DIO\n")
			} else {
				fmt.Printf("🪙 Token %s: 0.000000\n", token)
			}

			return nil
		},
	}
}

// getWalletExportCommand returns the wallet export command
func (c *CLI) getWalletExportCommand() *cli.Command {
	return &cli.Command{
		Name:  "export",
		Usage: "Export wallet private key",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Wallet name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				Usage:    "Wallet password",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"d"},
				Usage:   "Wallet directory path",
				Value:   c.config.Keystore.Path,
			},
		},
		Action: func(ctx *cli.Context) error {
			name := ctx.String("name")
			password := ctx.String("password")
			path := ctx.String("path")

			c.logger.Info("Exporting wallet", "name", name, "path", path)

			ks := keystore.NewKeystore(path)
			privateKey, err := ks.ExportPrivateKey(name, password)
			if err != nil {
				return fmt.Errorf("failed to export wallet: %w", err)
			}

			fmt.Printf("⚠️  WARNING: Keep your private key secure!\n")
			fmt.Printf("📝 Wallet: %s\n", name)
			fmt.Printf("🔑 Private Key: %s\n", privateKey)

			return nil
		},
	}
}

// getWalletDeleteCommand returns the wallet delete command
func (c *CLI) getWalletDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete wallet",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Wallet name",
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force delete without confirmation",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"d"},
				Usage:   "Wallet directory path",
				Value:   c.config.Keystore.Path,
			},
		},
		Action: func(ctx *cli.Context) error {
			name := ctx.String("name")
			force := ctx.Bool("force")
			path := ctx.String("path")

			c.logger.Info("Deleting wallet", "name", name, "path", path, "force", force)

			if !force {
				fmt.Printf("⚠️  Are you sure you want to delete wallet '%s'? This action cannot be undone.\n", name)
				fmt.Printf("Type 'yes' to confirm: ")
				var confirmation string
				fmt.Scanln(&confirmation)
				if confirmation != "yes" {
					fmt.Printf("❌ Wallet deletion cancelled.\n")
					return nil
				}
			}

			ks := keystore.NewKeystore(path)
			err := ks.DeleteAccount(name)
			if err != nil {
				return fmt.Errorf("failed to delete wallet: %w", err)
			}

			fmt.Printf("✅ Wallet '%s' deleted successfully.\n", name)

			return nil
		},
	}
}
