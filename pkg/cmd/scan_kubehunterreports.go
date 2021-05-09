package cmd

import (
	"context"

	"github.com/aquasecurity/starboard/pkg/generated/clientset/versioned"
	"github.com/aquasecurity/starboard/pkg/kubehunter"
	"github.com/aquasecurity/starboard/pkg/starboard"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

const (
	kubeHunterCmdShort = "Hunt for security weaknesses in your Kubernetes cluster"
)

func NewScanKubeHunterReportsCmd(cf *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubehunterreports",
		Short: kubeHunterCmdShort,
		RunE:  ScanKubeHunterReports(cf),
	}

	registerScannerOpts(cmd)

	return cmd
}

const (
	kubeHunterReportName = "cluster"
)

func ScanKubeHunterReports(cf *genericclioptions.ConfigFlags) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		kubeConfig, err := cf.ToRESTConfig()
		if err != nil {
			return err
		}
		kubeClientset, err := kubernetes.NewForConfig(kubeConfig)
		if err != nil {
			return err
		}
		opts, err := getScannerOpts(cmd)
		if err != nil {
			return err
		}
		config, err := starboard.NewConfigManager(kubeClientset, starboard.NamespaceName).Read(ctx)
		if err != nil {
			return err
		}
		report, err := kubehunter.NewScanner(starboard.NewScheme(), config, kubeClientset, opts).Scan(ctx)
		if err != nil {
			return err
		}
		starboardClientset, err := versioned.NewForConfig(kubeConfig)
		if err != nil {
			return err
		}
		return kubehunter.NewWriter(starboardClientset).Write(ctx, report, kubeHunterReportName)
	}
}
