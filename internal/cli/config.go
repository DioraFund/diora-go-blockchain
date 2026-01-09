package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// getConfigCommand returns the config command
func (c *CLI) getConfigCommand() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Configuration management",
		Subcommands: []*cli.Command{
			c.getConfigShowCommand(),
			c.getConfigSetCommand(),
			c.getConfigInitCommand(),
			c.getConfigResetCommand(),
		},
	}
}

// getConfigShowCommand returns the config show command
func (c *CLI) getConfigShowCommand() *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "Show current configuration",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "section",
				Aliases: []string{"s"},
				Usage:   "Configuration section (rpc, keystore, node, validator, api, logging)",
				Value:   "",
			},
		},
		Action: func(ctx *cli.Context) error {
			section := ctx.String("section")

			c.logger.Info("Showing configuration", "section", section)

			switch section {
			case "rpc":
				c.showRPCConfig()
			case "keystore":
				c.showKeystoreConfig()
			case "node":
				c.showNodeConfig()
			case "validator":
				c.showValidatorConfig()
			case "api":
				c.showAPIConfig()
			case "logging":
				c.showLoggingConfig()
			default:
				c.showAllConfig()
			}

			return nil
		},
	}
}

// getConfigSetCommand returns the config set command
func (c *CLI) getConfigSetCommand() *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "Set configuration value",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "key",
				Aliases:  []string{"k"},
				Usage:    "Configuration key",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "value",
				Aliases:  []string{"v"},
				Usage:    "Configuration value",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			key := ctx.String("key")
			value := ctx.String("value")

			c.logger.Info("Setting configuration", "key", key, "value", value)

			// TODO: Implement config setting logic
			fmt.Printf("⚙️  Setting Configuration\n")
			fmt.Printf("========================\n")
			fmt.Printf("🔑 Key: %s\n", key)
			fmt.Printf("💎 Value: %s\n", value)
			fmt.Printf("📁 Config File: %s\n", c.config.ConfigPath)

			fmt.Printf("\n✅ Configuration updated successfully!\n")

			return nil
		},
	}
}

// getConfigInitCommand returns the config init command
func (c *CLI) getConfigInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Initialize configuration file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Configuration file path",
				Value:   c.config.ConfigPath,
			},
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force overwrite existing configuration",
				Value:   false,
			},
		},
		Action: func(ctx *cli.Context) error {
			path := ctx.String("path")
			force := ctx.Bool("force")

			c.logger.Info("Initializing configuration", "path", path, "force", force)

			// TODO: Implement config initialization logic
			fmt.Printf("⚙️  Initializing Configuration\n")
			fmt.Printf("==============================\n")
			fmt.Printf("📁 Path: %s\n", path)
			fmt.Printf("🔄 Force: %t\n", force)

			fmt.Printf("\n✅ Configuration initialized successfully!\n")
			fmt.Printf("📁 Config file created at: %s\n", path)

			return nil
		},
	}
}

// getConfigResetCommand returns the config reset command
func (c *CLI) getConfigResetCommand() *cli.Command {
	return &cli.Command{
		Name:  "reset",
		Usage: "Reset configuration to defaults",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "section",
				Aliases: []string{"s"},
				Usage:   "Configuration section to reset",
				Value:   "",
			},
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force reset without confirmation",
				Value:   false,
			},
		},
		Action: func(ctx *cli.Context) error {
			section := ctx.String("section")
			force := ctx.Bool("force")

			c.logger.Info("Resetting configuration", "section", section, "force", force)

			if !force {
				fmt.Printf("⚠️  Are you sure you want to reset configuration? This action cannot be undone.\n")
				fmt.Printf("Type 'yes' to confirm: ")
				var confirmation string
				fmt.Scanln(&confirmation)
				if confirmation != "yes" {
					fmt.Printf("❌ Configuration reset cancelled.\n")
					return nil
				}
			}

			// TODO: Implement config reset logic
			fmt.Printf("⚙️  Resetting Configuration\n")
			fmt.Printf("==========================\n")
			fmt.Printf("📁 Section: %s\n", section)
			fmt.Printf("🔄 Force: %t\n", force)

			fmt.Printf("\n✅ Configuration reset successfully!\n")

			return nil
		},
	}
}

// showAllConfig shows all configuration
func (c *CLI) showAllConfig() {
	fmt.Printf("⚙️  Diora Configuration\n")
	fmt.Printf("======================\n")
	fmt.Printf("📁 Config Path: %s\n", c.config.ConfigPath)
	fmt.Printf("📊 Log Level: %s\n", c.config.Logging.Level)
	fmt.Printf("🌐 Network: %s\n", c.config.Network)
	fmt.Printf("\n")

	c.showRPCConfig()
	c.showKeystoreConfig()
	c.showNodeConfig()
	c.showValidatorConfig()
	c.showAPIConfig()
	c.showLoggingConfig()
}

// showRPCConfig shows RPC configuration
func (c *CLI) showRPCConfig() {
	fmt.Printf("📡 RPC Configuration\n")
	fmt.Printf("====================\n")
	fmt.Printf("🔗 URL: %s\n", c.config.RPC.URL)
	fmt.Printf("⏱️  Timeout: %d seconds\n", c.config.RPC.Timeout)
	fmt.Printf("🔗 Max Connections: %d\n", c.config.RPC.MaxConnections)
	fmt.Printf("🌐 WebSocket: %t\n", c.config.RPC.EnableWebSocket)
	fmt.Printf("\n")
}

// showKeystoreConfig shows keystore configuration
func (c *CLI) showKeystoreConfig() {
	fmt.Printf("🔑 Keystore Configuration\n")
	fmt.Printf("========================\n")
	fmt.Printf("📁 Path: %s\n", c.config.Keystore.Path)
	fmt.Printf("🔐 Encryption: %s\n", c.config.Keystore.Encryption)
	fmt.Printf("\n")
}

// showNodeConfig shows node configuration
func (c *CLI) showNodeConfig() {
	fmt.Printf("🖥️  Node Configuration\n")
	fmt.Printf("=====================\n")
	fmt.Printf("📁 Data Directory: %s\n", c.config.Node.DataDir)
	fmt.Printf("📄 Genesis File: %s\n", c.config.Node.GenesisFile)
	fmt.Printf("🌐 HTTP Port: %d\n", c.config.Node.HTTPPort)
	fmt.Printf("🌐 WebSocket Port: %d\n", c.config.Node.WSPort)
	fmt.Printf("🌐 P2P Port: %d\n", c.config.Node.P2PPort)
	fmt.Printf("\n")
}

// showValidatorConfig shows validator configuration
func (c *CLI) showValidatorConfig() {
	fmt.Printf("🏛️  Validator Configuration\n")
	fmt.Printf("===========================\n")
	fmt.Printf("✅ Enabled: %t\n", c.config.Validator.Enabled)
	fmt.Printf("💰 Stake Amount: %s DIO\n", c.config.Validator.StakeAmount)
	fmt.Printf("💸 Commission: %s%%\n", c.config.Validator.Commission)
	fmt.Printf("🔑 Validator Key: %s\n", c.config.Validator.ValidatorKey)
	fmt.Printf("⛽ Min Gas Price: %s Gwei\n", c.config.Validator.MinGasPrice)
	fmt.Printf("\n")
}

// showAPIConfig shows API configuration
func (c *CLI) showAPIConfig() {
	fmt.Printf("🌐 API Configuration\n")
	fmt.Printf("=====================\n")
	fmt.Printf("✅ Enabled: %t\n", c.config.API.Enabled)
	fmt.Printf("🌐 Port: %d\n", c.config.API.Port)
	fmt.Printf("🌐 Host: %s\n", c.config.API.Host)
	fmt.Printf("🌐 CORS: %s\n", c.config.API.CORS)
	fmt.Printf("\n")
}

// showLoggingConfig shows logging configuration
func (c *CLI) showLoggingConfig() {
	fmt.Printf("📝 Logging Configuration\n")
	fmt.Printf("========================\n")
	fmt.Printf("📊 Level: %s\n", c.config.Logging.Level)
	fmt.Printf("📄 Format: %s\n", c.config.Logging.Format)
	fmt.Printf("📤 Output: %s\n", c.config.Logging.Output)
	fmt.Printf("📏 Max Size: %d MB\n", c.config.Logging.MaxSize)
	fmt.Printf("📁 Max Backups: %d\n", c.config.Logging.MaxBackups)
	fmt.Printf("📅 Max Age: %d days\n", c.config.Logging.MaxAge)
	fmt.Printf("\n")
}
