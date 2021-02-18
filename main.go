package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/betorvs/secretpublisher/config"
	_ "github.com/betorvs/secretpublisher/gateway/secret"
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
		secret := usecase.GenerateSecret(secretName)
		err := usecase.ManageSecret(secretName, secret)
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
		secret := usecase.GenerateSecret(secretName)
		err := usecase.CreateSecret(secretName, secret)
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
		secret := usecase.GenerateSecret(secretName)
		err := usecase.UpdateSecret(secretName, secret)
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
		res, err := usecase.CheckSecret(secretName, config.SecretNamespace)
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

var scanSecretsCmd = &cobra.Command{
	Use:   "scan-secrets",
	Short: "scan-secrets label=value",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("[ERROR] Need label=value")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		labels := args[0]
		res, err := usecase.ScanSecret(labels)
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(2)
		}
		fmt.Printf("%s", res)
	},
}

var scanCMCmd = &cobra.Command{
	Use:   "scan-configmaps",
	Short: "scan-configmaps label=value",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("[ERROR] Need label=value")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		labels := args[0]
		res, err := usecase.ScanConfigMap(labels)
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(2)
		}
		fmt.Printf("%s", res)
	},
}

var scanSecretsValuesCmd = &cobra.Command{
	Use:   "secret-subvalue",
	Short: "secret-subvalue label=value",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("[ERROR] Need label=value")
		}
		if !strings.Contains(config.MatchKey, ".") {
			return errors.New("--matchKey key.subkey")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		labels := args[0]
		res, err := usecase.ScanSubvalueSecret(labels)
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(2)
		}
		fmt.Printf("%s", res)
	},
}

func initCommands() {
	existCmd.Flags().StringVar(&config.SecretNamespace, "secretNamespace", os.Getenv("SECRET_NAMESPACE"), "Secret namespace in Kubernetes")
	existCmd.Flags().StringToStringVar(&config.StringData, "stringData", config.ParseStringData("data"), "map for stringData in secret, use: key=value")
	existCmd.Flags().StringToStringVar(&config.Labels, "labels", config.ParseStringData("labels"), "map for labels in secret, use: key=value")
	existCmd.Flags().StringToStringVar(&config.Annotations, "annotations", config.ParseStringData("annotations"), "map for annotations in secret, use: key=value")
	createCmd.Flags().StringVar(&config.SecretNamespace, "secretNamespace", os.Getenv("SECRET_NAMESPACE"), "Secret namespace in Kubernetes")
	createCmd.Flags().StringToStringVar(&config.StringData, "stringData", config.ParseStringData("data"), "map for stringData in secret, use: key=value")
	createCmd.Flags().StringToStringVar(&config.Labels, "labels", config.ParseStringData("labels"), "map for labels in secret, use: key=value")
	createCmd.Flags().StringToStringVar(&config.Annotations, "annotations", config.ParseStringData("annotations"), "map for annotations in secret, use: key=value")
	updateCmd.Flags().StringVar(&config.SecretNamespace, "secretNamespace", os.Getenv("SECRET_NAMESPACE"), "Secret namespace in Kubernetes")
	updateCmd.Flags().StringToStringVar(&config.StringData, "stringData", config.ParseStringData("data"), "map for stringData in secret, use: key=value")
	updateCmd.Flags().StringToStringVar(&config.Labels, "labels", config.ParseStringData("labels"), "map for labels in secret, use: key=value")
	updateCmd.Flags().StringToStringVar(&config.Annotations, "annotations", config.ParseStringData("annotations"), "map for annotations in secret, use: key=value")
	checkCmd.Flags().StringVar(&config.SecretNamespace, "secretNamespace", os.Getenv("SECRET_NAMESPACE"), "Secret namespace in Kubernetes")
	deleteCmd.Flags().StringVar(&config.SecretNamespace, "secretNamespace", os.Getenv("SECRET_NAMESPACE"), "Secret namespace in Kubernetes")
	scanSecretsCmd.Flags().StringVar(&config.SecretNamespace, "secretNamespace", os.Getenv("SECRET_NAMESPACE"), "Secret namespace in Kubernetes")
	scanSecretsCmd.Flags().StringVar(&config.DestinationNamespace, "destinationNamespace", os.Getenv("DESTINATION_NAMESPACE"), "Destination Secret namespace in Secret Receiver")
	scanSecretsCmd.Flags().StringVar(&config.NameSuffix, "nameSuffix", os.Getenv("NAME_SUFFIX"), "Destination Secret name suffix in Secret Receiver")
	scanCMCmd.Flags().StringVar(&config.SecretNamespace, "secretNamespace", os.Getenv("SECRET_NAMESPACE"), "Secret namespace in Kubernetes")
	scanCMCmd.Flags().StringVar(&config.DestinationNamespace, "destinationNamespace", os.Getenv("DESTINATION_NAMESPACE"), "Destination Secret namespace in Secret Receiver")
	scanCMCmd.Flags().StringVar(&config.NameSuffix, "nameSuffix", os.Getenv("NAME_SUFFIX"), "Destination Secret name suffix in Secret Receiver")
	scanSecretsValuesCmd.Flags().StringVar(&config.SecretNamespace, "secretNamespace", os.Getenv("SECRET_NAMESPACE"), "Secret namespace in Kubernetes")
	scanSecretsValuesCmd.Flags().StringVar(&config.DestinationNamespace, "destinationNamespace", os.Getenv("DESTINATION_NAMESPACE"), "Destination Secret namespace in Secret Receiver")
	scanSecretsValuesCmd.Flags().StringVar(&config.NameSuffix, "nameSuffix", os.Getenv("NAME_SUFFIX"), "Destination Secret name suffix in Secret Receiver from value in Secret")
	scanSecretsValuesCmd.Flags().StringVar(&config.MatchKey, "matchKey", os.Getenv("MATCH_KEY"), "Key inside Secret to be exported to Secret Receiver")
	scanSecretsValuesCmd.Flags().StringVar(&config.NewLabels, "newLabels", os.Getenv("NEW_LABELS"), "New Labels to be exported to Secret Receiver")
	scanSecretsValuesCmd.Flags().StringVar(&config.NewAnnotations, "newAnnotations", os.Getenv("NEW_ANNOTATIONS"), "New Annotations to be exported to Secret Receiver")
}

func main() {
	rootCmd := config.ConfigureRootCommand()
	initCommands()
	rootCmd.AddCommand(versionCmd, existCmd, createCmd, updateCmd, checkCmd, deleteCmd, scanSecretsCmd, scanCMCmd, scanSecretsValuesCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR]: %v\n", err)
		os.Exit(1)
	}
}
