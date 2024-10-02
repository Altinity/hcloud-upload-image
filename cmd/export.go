package cmd

import (
	"fmt"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/spf13/cobra"

	"github.com/apricote/hcloud-upload-image/hcloudimages"
	"github.com/apricote/hcloud-upload-image/hcloudimages/contextlogger"
)

const (
	exportFlagImageID      = "image-id"
	exportFlagCmd          = "cmd"
	exportFlagArchitecture = "architecture"
	exportFlagServerType   = "server-type"
	exportFlagLabels       = "labels"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export --image-id=<ID> --cmd=<cmd> --architecture=<x86|arm>",
	Short: "Export the image by running a custom command in the rescue mode.",

	GroupID: "primary",

	PreRun: initClient,

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := contextlogger.From(ctx)

		imageID, _ := cmd.Flags().GetInt64(exportFlagImageID)
		exportCmd, _ := cmd.Flags().GetString(exportFlagCmd)
		architecture, _ := cmd.Flags().GetString(exportFlagArchitecture)
		serverType, _ := cmd.Flags().GetString(exportFlagServerType)
		labels, _ := cmd.Flags().GetStringToString(exportFlagLabels)

		options := hcloudimages.ExportOptions{
			ImageID: imageID,
			Cmd:     exportCmd,
			Labels:  labels,
		}

		if architecture != "" {
			options.Architecture = hcloud.Architecture(architecture)
		} else if serverType != "" {
			options.ServerType = &hcloud.ServerType{Name: serverType}
		}

		_, err := client.Export(ctx, options)
		if err != nil {
			return fmt.Errorf("failed to perform export operation: %w", err)
		}

		logger.InfoContext(ctx, "Successfully performed export operation")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().Int64(exportFlagImageID, 0, "Image ID")
	exportCmd.MarkFlagRequired(exportFlagImageID)

	exportCmd.Flags().String(exportFlagCmd, "", "Export command")
	exportCmd.MarkFlagRequired(exportFlagCmd)

	exportCmd.Flags().String(exportFlagArchitecture, "", "CPU architecture of the disk image [choices: x86, arm]")
	_ = exportCmd.RegisterFlagCompletionFunc(
		exportFlagArchitecture,
		cobra.FixedCompletions([]string{string(hcloud.ArchitectureX86), string(hcloud.ArchitectureARM)}, cobra.ShellCompDirectiveNoFileComp),
	)

	exportCmd.Flags().String(exportFlagServerType, "", "Explicitly use this server type to export the image. Mutually exclusive with --architecture.")

	// Only one of them needs to be set
	exportCmd.MarkFlagsOneRequired(exportFlagArchitecture, exportFlagServerType)
	exportCmd.MarkFlagsMutuallyExclusive(exportFlagArchitecture, exportFlagServerType)

	exportCmd.Flags().StringToString(exportFlagLabels, map[string]string{}, "Labels for the resulting image")
}
