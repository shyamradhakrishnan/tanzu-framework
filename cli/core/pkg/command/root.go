// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aunum/log"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/vmware-tanzu/tanzu-framework/cli/core/pkg/cli"
	cliconfig "github.com/vmware-tanzu/tanzu-framework/cli/core/pkg/config"
	"github.com/vmware-tanzu/tanzu-framework/cli/core/pkg/pluginmanager"
	cliapi "github.com/vmware-tanzu/tanzu-framework/cli/runtime/apis/cli/v1alpha1"
	configapi "github.com/vmware-tanzu/tanzu-framework/cli/runtime/apis/config/v1alpha1"
	"github.com/vmware-tanzu/tanzu-framework/cli/runtime/component"
	"github.com/vmware-tanzu/tanzu-framework/cli/runtime/config"
)

// RootCmd is the core root Tanzu command
var RootCmd = &cobra.Command{
	Use: "tanzu",
	// Don't have Cobra print the error message, the CLI will
	// print it itself in a nicer format.
	SilenceErrors: true,
}

var (
	noInit      bool
	forceNoInit = "true" // a string variable so as to be overridable via linker flag
)

// NewRootCmd creates a root command.
func NewRootCmd() (*cobra.Command, error) {
	uFunc := cli.NewMainUsage().Func()
	RootCmd.SetUsageFunc(uFunc)
	k8sCmd.SetUsageFunc(uFunc)
	tmcCmd.SetUsageFunc(uFunc)

	ni := os.Getenv("TANZU_CLI_NO_INIT")
	if ni != "" || strings.EqualFold(forceNoInit, "true") {
		noInit = true
	}

	// configure defined environment variables under tanzu config file
	cliconfig.ConfigureEnvVariables()

	RootCmd.Short = component.Bold(`Tanzu CLI`)

	// TODO (pbarker): silencing usage for now as we are getting double usage from plugins on errors
	RootCmd.SilenceUsage = true

	RootCmd.AddCommand(
		pluginCmd,
		initCmd,
		updateCmd,
		versionCmd,
		completionCmd,
		configCmd,
		genAllDocsCmd,
	)

	// If the context and target feature is enabled, add the corresponding commands under root.
	if config.IsFeatureActivated(config.FeatureContextCommand) {
		RootCmd.AddCommand(
			contextCmd,
			k8sCmd,
			tmcCmd,
		)
		if err := addCtxPlugins(k8sCmd, configapi.CtxTypeK8s); err != nil {
			return nil, err
		}
		if err := addCtxPlugins(tmcCmd, configapi.CtxTypeTMC); err != nil {
			return nil, err
		}
	}

	plugins, err := getAvailablePlugins()
	if err != nil {
		return nil, err
	}

	if err = config.CopyLegacyConfigDir(); err != nil {
		return nil, fmt.Errorf("failed to copy legacy configuration directory to new location: %w", err)
	}

	// If context-aware-cli-for-plugins feature is not enabled
	// check that all plugins in the core distro are installed or do so.
	if !config.IsFeatureActivated(config.FeatureContextAwareCLIForPlugins) {
		plugins, err = checkAndInstallMissingPlugins(plugins)
		if err != nil {
			return nil, err
		}
	}

	for _, plugin := range plugins {
		RootCmd.AddCommand(cli.GetCmd(plugin))
	}

	duplicateAliasWarning()

	// Flag parsing must be deactivated because the root plugin won't know about all flags.
	RootCmd.DisableFlagParsing = true

	return RootCmd, nil
}

var k8sCmd = &cobra.Command{
	Use:     "kubernetes",
	Short:   "Tanzu CLI plugins that target a Kubernetes cluster",
	Aliases: []string{"k8s"},
	Annotations: map[string]string{
		"group": string(cliapi.TargetCmdGroup),
	},
}

var tmcCmd = &cobra.Command{
	Use:     "mission-control",
	Short:   "Tanzu CLI plugins that target a Tanzu Mission Control endpoint",
	Aliases: []string{"tmc"},
	Annotations: map[string]string{
		"group": string(cliapi.TargetCmdGroup),
	},
}

func addCtxPlugins(cmd *cobra.Command, ctxType configapi.ContextType) error {
	var ctxName string
	if ctx, _ := config.GetCurrentContext(ctxType); ctx != nil {
		ctxName = ctx.Name
	}

	ctxPlugins, standalonePlugins, err := pluginmanager.InstalledPlugins(ctxName)
	if err != nil {
		return fmt.Errorf("unable to find installed plugins: %w", err)
	}

	if ctxType == configapi.CtxTypeK8s {
		// Standalone plugins exist only for K8s context type.
		for i := range standalonePlugins {
			if standalonePlugins[i].Group == cliapi.SystemCmdGroup {
				// Do not include plugins from the system command group.
				continue
			}
			cmd.AddCommand(cli.GetCmd(&standalonePlugins[i]))
		}
	}

	for i := range ctxPlugins {
		cmd.AddCommand(cli.GetCmd(&ctxPlugins[i]))
	}
	return nil
}

func getAvailablePlugins() ([]*cliapi.PluginDescriptor, error) {
	plugins := make([]*cliapi.PluginDescriptor, 0)
	var err error

	if config.IsFeatureActivated(config.FeatureContextAwareCLIForPlugins) {
		currentServerName := ""

		server, err := config.GetCurrentServer()
		if err == nil && server != nil {
			currentServerName = server.Name
		}

		serverPlugin, standalonePlugins, err := pluginmanager.InstalledPlugins(currentServerName)
		if err != nil {
			return nil, fmt.Errorf("find installed plugins: %w", err)
		}

		//nolint:gocritic
		p := append(serverPlugin, standalonePlugins...)
		for i := range p {
			plugins = append(plugins, &p[i])
		}
	} else {
		// TODO: cli.ListPlugins is deprecated: Use pluginmanager.AvailablePluginsFromLocalSource or pluginmanager.AvailablePlugins instead
		//nolint:staticcheck
		plugins, err = cli.ListPlugins()
		if err != nil {
			return nil, fmt.Errorf("find available plugins: %w", err)
		}
	}
	return plugins, nil
}

func checkAndInstallMissingPlugins(plugins []*cliapi.PluginDescriptor) ([]*cliapi.PluginDescriptor, error) {
	// check that all plugins in the core distro are installed or do so.
	if !noInit && !cli.IsDistributionSatisfied(plugins) {
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		if err := s.Color("bgBlack", "bold", "fgWhite"); err != nil {
			return nil, err
		}
		s.Suffix = fmt.Sprintf(" %s", "initializing")
		s.Start()
		cfg, err := config.GetClientConfig()
		if err != nil {
			log.Fatal(err)
		}
		repos := cli.NewMultiRepo(cli.LoadRepositories(cfg)...)
		err = cli.EnsureDistro(repos)
		if err != nil {
			return nil, err
		}
		// TODO: cli.ListPlugins is deprecated: Use pluginmanager.AvailablePluginsFromLocalSource or pluginmanager.AvailablePlugins instead
		//nolint:staticcheck
		plugins, err = cli.ListPlugins()
		if err != nil {
			return nil, fmt.Errorf("find available plugins: %w", err)
		}
		s.Stop()
	}
	return plugins, nil
}

func duplicateAliasWarning() {
	var aliasMap = make(map[string][]string)
	for _, command := range RootCmd.Commands() {
		for _, alias := range command.Aliases {
			aliases, ok := aliasMap[alias]
			if !ok {
				aliasMap[alias] = []string{command.Name()}
			} else {
				aliasMap[alias] = append(aliases, command.Name())
			}
		}
	}

	for alias, plugins := range aliasMap {
		if len(plugins) > 1 {
			fmt.Fprintf(os.Stderr, "Warning, the alias %s is duplicated across plugins: %s\n\n", alias, strings.Join(plugins, ", "))
		}
	}
}

// Execute executes the CLI.
func Execute() error {
	root, err := NewRootCmd()
	if err != nil {
		return err
	}
	return root.Execute()
}
