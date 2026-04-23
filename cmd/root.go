package cmd

import (
	"os"

	"ddnsd/internal/config"
	"ddnsd/internal/controller"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ddnsd",
		Short: "Dynamic DNS daemon to update DNS records when your IP address changes",
		Long:  `A provider-agnostic daemon that watches your IP address and updates DNS records when it changes.`,
		Run:   run,
	}

	cmd.Flags().StringP("config", "c", "", "path to config file (required)")
	err := cmd.MarkFlagRequired("config")
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal error")
		os.Exit(1)
	}

	return cmd
}

func run(cmd *cobra.Command, _ []string) {
	cfgFilePath := cmd.Flag("config").Value.String()
	cfg := initConfig(cfgFilePath)

	loglevel := viper.GetString("loglevel")
	ll, err := zerolog.ParseLevel(loglevel)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse loglevel")
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(ll)

	c := controller.NewController(cfg)
	err = c.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("fatal error")
		os.Exit(1)
	}
}

func initConfig(cfgFilePath string) *config.Config {
	viper.SetConfigFile(cfgFilePath)
	viper.SetEnvPrefix("DDNSD")
	viper.AutomaticEnv()

	viper.SetDefault("checkInterval", "10")
	viper.SetDefault("loglevel", "info")

	rawBytes, err := os.ReadFile(viper.ConfigFileUsed())
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal error reading config file")
		os.Exit(1)
	}

	var cfg config.Config
	err = yaml.Unmarshal(rawBytes, &cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal error unmarshalling config file")
		os.Exit(1)
	}

	return &cfg
}
