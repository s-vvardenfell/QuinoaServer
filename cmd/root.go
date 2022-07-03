package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/s-vvardenfell/QuinoaServer/config"
	"github.com/s-vvardenfell/QuinoaServer/generated"
	"github.com/s-vvardenfell/QuinoaServer/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cnfg config.Config

var rootCmd = &cobra.Command{
	Use:   "QuinoaServer",
	Short: "Server for Quinoa project that connects Parser, TgBot and Cache",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("QuinoaServar started!")
		grpcServ := grpc.NewServer()
		qServ := server.NewServer(cnfg)
		generated.RegisterMainServiceServer(grpcServ, qServ)

		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cnfg.ServerHost, cnfg.ServerPort))
		if err != nil {
			logrus.Fatalf("failed to listen: %v", err)
		}

		if cnfg.WithReflection {
			reflection.Register(grpcServ)
		}

		logrus.Info("Starting gRPC listener on port " + cnfg.ServerPort)
		if err := grpcServ.Serve(lis); err != nil {
			logrus.Fatalf("failed to serve: %v", err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is resources/config.yml)")
	// rootCmd.Flags().BoolP("debug", "d", false, "Runs in debug-mode")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		wd, err := os.Getwd()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(wd, "resources"))
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		cobra.CheckErr(err)
	}

	if err := viper.Unmarshal(&cnfg); err != nil {
		cobra.CheckErr(err)
	}

	if cnfg.Logrus.ToFile {
		if err := os.Mkdir(filepath.Dir(cnfg.Logrus.LogDir), 0644); err != nil && !errors.Is(err, os.ErrExist) {
			cobra.CheckErr(err)
		}

		file, err := os.OpenFile(cnfg.Logrus.LogDir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			logrus.SetOutput(file)
		} else {
			cobra.CheckErr(err)
		}
	}

	if cnfg.Logrus.ToJson {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	logrus.SetLevel(logrus.Level(cnfg.Logrus.LogLvl))
}
