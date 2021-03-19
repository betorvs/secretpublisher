package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	// EncodingRequest string
	EncodingRequest string
	// ReceiverURL string
	ReceiverURL string
	// CommandTimeout string
	CommandTimeout string
	// SecretNamespace string
	SecretNamespace string
	// StringData map[string]string
	StringData map[string]string
	// Labels map[string]string
	Labels map[string]string
	// Annotations map[string]string
	Annotations map[string]string
	// PublisherTimeout time.Duration
	PublisherTimeout time.Duration
	// TestRun string
	TestRun string
	// LocalKubeconfig bool
	LocalKubeconfig bool
	// DestinationNamespace string
	DestinationNamespace string
	// NameSuffix string
	NameSuffix string
	// KeyNameSuffix string
	KeyNameSuffix string
	// MatchKey string
	MatchKey string
	// NewLabels string
	NewLabels string
	// NewAnnotations string
	NewAnnotations string
	// DisabledLabel string
	DisabledLabel string
	// MiddleName string
	MiddleName string
	// Debug bool
	Debug bool
)

// ParseStringData func
func ParseStringData(kind string) map[string]string {
	data := make(map[string]string)
	switch kind {
	case "data":
		if os.Getenv("STRING_DATA") != "" {
			data = ParseLabelsArg(os.Getenv("STRING_DATA"))
		}
	case "labels":
		if os.Getenv("LABELS") != "" {
			data = ParseLabelsArg(os.Getenv("LABELS"))
		}
	case "annotations":
		if os.Getenv("ANNOTATIONS") != "" {
			data = ParseLabelsArg(os.Getenv("ANNOTATIONS"))
		}
	}

	return data
}

// ParseLabelsArg func returns a map[string]string from a string
func ParseLabelsArg(labelArg string) map[string]string {
	labels := map[string]string{}
	if labelArg == "" {
		return labels
	}

	if strings.Contains(labelArg, ",") {
		pairs := strings.Split(labelArg, ",")
		for _, pair := range pairs {
			parts := strings.Split(pair, "=")
			if len(parts) == 2 {
				labels[parts[0]] = parts[1]
			}
		}
	} else {
		parts := strings.Split(labelArg, "=")
		labels[parts[0]] = parts[1]
	}

	return labels
}

// ConfigureRootCommand func
func ConfigureRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secretpublisher",
		Short: "\nSecret Publisher is a command line tool to interact with Secret Receiver",
		RunE:  run,
	}
	cmd.PersistentFlags().StringVar(&EncodingRequest, "encodingRequest", os.Getenv("ENCODING_REQUEST"), "use ENCODING_REQUEST environment variable")
	cmd.PersistentFlags().StringVar(&ReceiverURL, "receiverURL", os.Getenv("RECEIVER_URL"), "use RECEIVER_URL environment variable")
	cmd.PersistentFlags().StringVar(&TestRun, "testRun", "false", "use TESTRUN environment variable")
	cmd.PersistentFlags().BoolVar(&LocalKubeconfig, "localKubeconfig", false, "use local kubeconfig file")
	cmd.PersistentFlags().BoolVar(&Debug, "debug", false, "add --debug in the command")
	cmd.PersistentFlags().StringVar(&CommandTimeout, "commandTimeout", os.Getenv("COMMAND_TIMEOUT"), "use COMMAND_TIMEOUT environment variable")
	return cmd
}

// run func do everything
func run(cmd *cobra.Command, args []string) error {
	if ReceiverURL == "" {
		return fmt.Errorf("ReceiverURL is empty")
	}
	if EncodingRequest == "" {
		EncodingRequest = "disabled"
	}
	tmpTimeout, err := strconv.Atoi(CommandTimeout)
	if err != nil {
		tmpTimeout = 15
	}
	PublisherTimeout = time.Duration(tmpTimeout)
	// errManage := usecase.ManageSecret()
	// if errManage != nil {
	// 	return errManage
	// }
	return nil
}
