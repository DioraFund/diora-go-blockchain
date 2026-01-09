package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/DioraFund/diora-go-blockchain/internal/config"
	"github.com/DioraFund/diora-go-blockchain/internal/logger"
	"github.com/urfave/cli/v2"
)

// CLI represents the command line interface
type CLI struct {
	ctx    context.Context
	config *config.Config
	logger *logger.Logger
	app    *cli.App
}

// NewCLI creates a new CLI instance
func NewCLI(ctx context.Context, cfg *config.Config, log *logger.Logger) (*CLI, error) {
	cliApp := &CLI{
		ctx:    ctx,
		config: cfg,
		logger: log,
	}

	// Initialize CLI application
	app := &cli.App{
		Name:     "diora",
		Version:  "1.0.0",
		Usage:    "ABM Diora Blockchain Command Line Interface",
		Commands: cliApp.getCommands(),
		Flags:    cliApp.getGlobalFlags(),
		Before:   cliApp.beforeAction,
		After:    cliApp.afterAction,
	}

	cliApp.app = app
	return cliApp, nil
}

// Run executes the CLI application
func (c *CLI) Run() error {
	return c.app.Run(os.Args)
}

// getCommands returns all available CLI commands
func (c *CLI) getCommands() []*cli.Command {
	return []*cli.Command{
		c.getWalletCommand(),
		c.getAccountCommand(),
		c.getNetworkCommand(),
		c.getValidatorCommand(),
		c.getTransactionCommand(),
		c.getContractCommand(),
		c.getConfigCommand(),
		c.getMonitorCommand(),
		c.getDevCommand(),
	}
}

// getGlobalFlags returns global CLI flags
func (c *CLI) getGlobalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to configuration file",
			Value:   c.config.ConfigPath,
		},
		&cli.StringFlag{
			Name:    "rpc-url",
			Aliases: []string{"r"},
			Usage:   "RPC endpoint URL",
			Value:   c.config.RPC.URL,
		},
		&cli.StringFlag{
			Name:    "keystore",
			Aliases: []string{"k"},
			Usage:   "Keystore directory path",
			Value:   c.config.Keystore.Path,
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Enable verbose logging",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "no-color",
			Aliases: []string{"nc"},
			Usage:   "Disable colored output",
			Value:   false,
		},
	}
}

// beforeAction is called before each command execution
func (c *CLI) beforeAction(ctx *cli.Context) error {
	// Update configuration with CLI flags
	if ctx.IsSet("config") {
		c.config.ConfigPath = ctx.String("config")
	}
	if ctx.IsSet("rpc-url") {
		c.config.RPC.URL = ctx.String("rpc-url")
	}
	if ctx.IsSet("keystore") {
		c.config.Keystore.Path = ctx.String("keystore")
	}
	if ctx.IsSet("verbose") {
		c.config.Logging.Level = "debug"
	}

	// Reload configuration if needed
	if err := c.config.Reload(); err != nil {
		return fmt.Errorf("failed to reload configuration: %w", err)
	}

	c.logger.Info("Executing command", "command", ctx.Command.Name)
	return nil
}

// afterAction is called after each command execution
func (c *CLI) afterAction(ctx *cli.Context) error {
	c.logger.Info("Command completed", "command", ctx.Command.Name)
	return nil
}
