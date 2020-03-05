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
	// PublisherTimeout time.Duration
	PublisherTimeout time.Duration
	// TestRun string
	TestRun string
	// Debug bool
	Debug bool
)

// ParseStringData func
func ParseStringData() map[string]string {
	data := make(map[string]string)
	if os.Getenv("STRING_DATA") != "" {
		values := strings.Split(os.Getenv("STRING_DATA"), ",")
		for _, e := range values {
			parts := strings.Split(e, "=")
			data[parts[0]] = parts[1]
		}
	}
	return data
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
	cmd.PersistentFlags().BoolVar(&Debug, "debug", false, "add --debug in the command")
	cmd.PersistentFlags().StringVar(&CommandTimeout, "commandTimeout", os.Getenv("COMMAND_TIMEOUT"), "use COMMAND_TIMEOUT environment variable")
	// cmd.PersistentFlags().StringToStringVar(&StringData, "stringData", getEnv(), "use STRING_DATA environment variable like: name=path_to_file,")
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
