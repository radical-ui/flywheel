package main

import (
	"fmt"
	"os"

	"github.com/radical-ui/flywheel/flutter"
	"github.com/spf13/cobra"
)

func main() {
	var logFile string

	rootCmd := &cobra.Command{
		Use:   "flywheel",
		Short: "An objection frontend, built with Flutter",
	}

	rootCmd.PersistentFlags().StringVarP(&logFile, "log-file", "l", "", "Write debug information to the specified logfile. Omitting or leaving empty will cause no logs to be written.")

	genCommand := &cobra.Command{
		Use:   "gen",
		Short: "Print the objects schema to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			runWithErrorHandling(logFile, runOptions{
				genBindings: true,
			})
		},
	}

	rootCmd.AddCommand(genCommand)

	var displayName string
	var url string

	previewCommand := &cobra.Command{
		Use:   "preview [bundle identifier]",
		Short: "Preview the application",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("a bundle identifier argument is required")
				os.Exit(1)
			}

			runWithErrorHandling(logFile, runOptions{
				genFlutter: &flutter.Options{
					DisplayName:      displayName,
					BundleIdentifier: args[0],
					Url:              url,
				},
				preview: true,
			})
		},
	}

	previewCommand.Flags().StringVar(&displayName, "display-name", "", "The application display name")
	previewCommand.Flags().StringVar(&url, "url", "ws://localhost:8000", "The objection server to connect to")

	rootCmd.AddCommand(previewCommand)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
