package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/betorvs/secretpublisher/config"
	_ "github.com/betorvs/secretpublisher/gateway"
	"github.com/betorvs/secretpublisher/usecase"
	"github.com/spf13/cobra"
)

var (
	// Version string
	Version = "development"
	// BuildInfo string
	BuildInfo string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of usernamectl",
	Long:  `All software has versions.`,
	Run: func(cmd *cobra.Command, args []string) {
		if BuildInfo != "" {
			fmt.Printf("secretpublisher command line tools version: %s, build: %s\n", Version, BuildInfo)
			os.Exit(0)
		}
		fmt.Printf("secretpublisher command line tools version: %s\n", Version)
		os.Exit(0)
	},
}

var existCmd = &cobra.Command{
	Use:   "exist",
	Short: "exist SECRET_NAME",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("[ERROR] Need at least secret name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		secretName := args[0]
		err := usecase.ManageSecret(secretName)
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(2)
		}
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create SECRET_NAME",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("[ERROR] Need at least secret name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		secretName := args[0]
		err := usecase.CreateSecret(secretName)
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(2)
		}
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update SECRET_NAME",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("[ERROR] Need at least secret name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		secretName := args[0]
		err := usecase.UpdateSecret(secretName)
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(2)
		}
	},
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check SECRET_NAME",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("[ERROR] Need at least secret name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		secretName := args[0]
		res, err := usecase.CheckSecret(secretName)
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(2)
		}
		fmt.Printf("%s", res)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete SECRET_NAME",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("Need at least secret name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		secretName := args[0]
		err := usecase.DeleteSecret(secretName)
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(2)
		}
	},
}

func initCommands() {
	existCmd.Flags().StringVar(&config.SecretNamespace, "secretNamespace", os.Getenv("SECRET_NAMESPACE"), "Secret namespace in Kubernetes")
	existCmd.Flags().StringToStringVar(&config.StringData, "stringData", config.ParseStringData("data"), "map for stringData in secret, use: key=value")
	existCmd.Flags().StringToStringVar(&config.Labels, "labels", config.ParseStringData("labels"), "map for labels in secret, use: key=value")
	createCmd.Flags().StringVar(&config.SecretNamespace, "secretNamespace", os.Getenv("SECRET_NAMESPACE"), "Secret namespace in Kubernetes")
	createCmd.Flags().StringToStringVar(&config.StringData, "stringData", config.ParseStringData("data"), "map for stringData in secret, use: key=value")
	createCmd.Flags().StringToStringVar(&config.Labels, "labels", config.ParseStringData("labels"), "map for labels in secret, use: key=value")
	updateCmd.Flags().StringVar(&config.SecretNamespace, "secretNamespace", os.Getenv("SECRET_NAMESPACE"), "Secret namespace in Kubernetes")
	updateCmd.Flags().StringToStringVar(&config.StringData, "stringData", config.ParseStringData("data"), "map for stringData in secret, use: key=value")
	updateCmd.Flags().StringToStringVar(&config.Labels, "labels", config.ParseStringData("labels"), "map for labels in secret, use: key=value")
	checkCmd.Flags().StringVar(&config.SecretNamespace, "secretNamespace", os.Getenv("SECRET_NAMESPACE"), "Secret namespace in Kubernetes")
	deleteCmd.Flags().StringVar(&config.SecretNamespace, "secretNamespace", os.Getenv("SECRET_NAMESPACE"), "Secret namespace in Kubernetes")
}

func main() {
	rootCmd := config.ConfigureRootCommand()
	initCommands()
	rootCmd.AddCommand(versionCmd, existCmd, createCmd, updateCmd, checkCmd, deleteCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR]: %v\n", err)
		os.Exit(1)
	}
}
