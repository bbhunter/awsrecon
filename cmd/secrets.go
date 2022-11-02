package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hupe1980/awsrecon/pkg/config"
	"github.com/hupe1980/awsrecon/pkg/output"
	"github.com/hupe1980/awsrecon/pkg/recon"
	"github.com/spf13/cobra"
)

type secretsOptions struct {
	entropy              float64
	decrypt              bool
	verify               bool
	highEntropyThreshold float64
}

func newSecretsCmd(globalOpts *globalOptions) *cobra.Command {
	opts := &secretsOptions{}
	cmd := &cobra.Command{
		Use:           "secrets",
		Short:         "Enumerate secrets",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.NewConfig(globalOpts.profile, globalOpts.regions, globalOpts.userAgent, globalOpts.timeout)
			if err != nil {
				return err
			}

			recon := recon.NewSecretsRecon(cfg, func(o *recon.SecretsOptions) {
				o.Entropy = opts.entropy
				o.WithDecryption = opts.decrypt
			})

			secrets := recon.Run()

			output := output.NewTable([]string{
				"Service",
				"Region",
				"Type",
				"Name",
				"Value",
				"Entropy",
				"Hints",
			})

			sort.Slice(secrets, func(i, j int) bool {
				return secrets[i].AWSService < secrets[j].AWSService
			})

			for _, s := range secrets {
				output.Add([]string{
					s.AWSService,
					s.Region,
					s.Type,
					s.Name,
					s.Value,
					fmt.Sprintf("%f", s.Entropy),
					strings.Join(s.Hints, ","),
				})
			}

			output.Print()

			return nil
		},
	}

	cmd.Flags().Float64VarP(&opts.entropy, "entropy", "e", 0, "minimum entropy")
	cmd.Flags().BoolVarP(&opts.decrypt, "decrypt", "d", false, "decrypt secret")
	cmd.Flags().BoolVarP(&opts.verify, "verify", "", false, "verify secrets")
	cmd.Flags().Float64VarP(&opts.highEntropyThreshold, "high-entropy-threshold", "", 3.5, "high entropy threshold")

	return cmd
}