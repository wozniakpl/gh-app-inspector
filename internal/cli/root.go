package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/wozniakpl/gh-app-inspector/internal/inspector"
	"github.com/wozniakpl/gh-app-inspector/internal/version"
)

type flags struct {
	appID          int64
	installationID int64
	pemPath        string
}

func NewRootCmd() *cobra.Command {
	var f flags
	cmd := &cobra.Command{
		Use:     "gh-app-inspector",
		Short:   "Inspect a GitHub App installation: permissions, repos, rate limit",
		Version: version.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := resolveEnvFallbacks(&f); err != nil {
				return err
			}
			if f.appID == 0 || f.installationID == 0 || f.pemPath == "" {
				return errors.New("missing required: --app-id, --installation-id, --pem (or GH_APP_ID / GH_INSTALLATION_ID / GH_APP_PEM)")
			}
			pem, err := os.ReadFile(f.pemPath)
			if err != nil {
				return fmt.Errorf("read pem: %w", err)
			}
			return inspector.Run(context.Background(), os.Stdout, inspector.Config{
				AppID:          f.appID,
				InstallationID: f.installationID,
				PEM:            pem,
			})
		},
	}
	cmd.Flags().Int64Var(&f.appID, "app-id", 0, "GitHub App ID (env: GH_APP_ID)")
	cmd.Flags().Int64Var(&f.installationID, "installation-id", 0, "Installation ID (env: GH_INSTALLATION_ID)")
	cmd.Flags().StringVar(&f.pemPath, "pem", "", "Path to private key PEM file (env: GH_APP_PEM)")
	return cmd
}

func resolveEnvFallbacks(f *flags) error {
	if f.appID == 0 {
		if v := os.Getenv("GH_APP_ID"); v != "" {
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf("GH_APP_ID: %w", err)
			}
			f.appID = n
		}
	}
	if f.installationID == 0 {
		if v := os.Getenv("GH_INSTALLATION_ID"); v != "" {
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf("GH_INSTALLATION_ID: %w", err)
			}
			f.installationID = n
		}
	}
	if f.pemPath == "" {
		f.pemPath = os.Getenv("GH_APP_PEM")
	}
	return nil
}
