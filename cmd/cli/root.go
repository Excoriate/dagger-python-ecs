package cli

import (
	"context"
	"fmt"
	"github.com/Excoriate/dagger-python-ecs/cmd/cli/aws"
	"github.com/Excoriate/dagger-python-ecs/cmd/cli/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	GlobalWorkingDirectory            string
	GlobalMountDir                    string
	GlobalTargetDir                   string
	GlobalTaskName                    string
	GlobalScanEnvVarKeys              []string
	GlobalEnvKeyValuePairsToSet       map[string]string
	GlobalCustomCommands              []string
	GlobalScanAWSKeys                 bool
	GlobalScanTFVars                  bool
	GlobalCustomCMDs                  []string
	GlobalDaggerInitClientWithWorkDir bool

	// Configuration file
	cfgFile string
)

var rootCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "stiletto",
	Long: `Stiletto is a command-line tool that helps automate the process of Continuous
Integration and Continuous Deployment (CD), using Dagger (dagger.io) as the main engine.
It provides options to manage various tasks such as building Docker images, running tests, linting, and deployment to AWS services.`,
	Example: `
  # DockerCmd pipeline with the specified working directory
  stiletto <command> --workdir /path/to/working/directory --task
  E.g.:
  pipeline docker --workdir /path/to/working/directory  --task=build`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		os.Exit(1)
	}
}

func addPersistentFlags() {
	rootCmd.PersistentFlags().StringVarP(&GlobalTaskName,
		"task",
		"t", "",
		"Name of the task to run. E.g.: build, test, lint, install")

	rootCmd.PersistentFlags().StringVarP(&GlobalWorkingDirectory,
		"work-dir",
		"w", "",
		"Work directory where the pipeline will be executed. If it's not set, "+
			"it'll use '.' value, which represents the current directory.")

	rootCmd.PersistentFlags().StringVarP(&GlobalTargetDir,
		"target-dir",
		"d", "",
		"Target directory represents the subdirectory within the mounted directory where the"+
			" actions (commands) will be executed.")

	rootCmd.PersistentFlags().StringVarP(&GlobalMountDir,
		"mount-dir",
		"m", "",
		"Mount directory represents what subdirectory within the working directory will be used"+
			" to mount into the container while it's performing its actions.")

	rootCmd.PersistentFlags().StringSliceVarP(&GlobalScanEnvVarKeys,
		"scan-env",
		"s", []string{},
		"List of environment variable keys that are already exported, that'll be scanned and set.")

	rootCmd.PersistentFlags().StringToStringVarP(&GlobalEnvKeyValuePairsToSet,
		"set-env",
		"e", map[string]string{},
		"List of environment variable key-value pairs to set.")

	rootCmd.PersistentFlags().StringSliceVarP(&GlobalCustomCommands,
		"commands",
		"c", []string{},
		"List of custom commands to run.")

	rootCmd.PersistentFlags().BoolVarP(&GlobalScanAWSKeys,
		"scan-aws-keys",
		"a", false,
		"Scan AWS keys and set them as environment variables.")

	rootCmd.PersistentFlags().BoolVarP(&GlobalScanTFVars,
		"scan-terraform-vars",
		"f", false,
		"Scan terraform exported environment variables and set it into the generated containers ("+
			"TG_VAR_).")

	rootCmd.PersistentFlags().StringSliceVarP(&GlobalCustomCMDs,
		"custom-cmds",
		"u", []string{},
		"List of custom commands to run.")

	rootCmd.PersistentFlags().BoolVarP(&GlobalDaggerInitClientWithWorkDir,
		"init-dagger-with-workdir",
		"", false,
		"Initialize Dagger client with the working directory.")

	_ = viper.BindPFlag("task", rootCmd.PersistentFlags().Lookup("task"))
	_ = viper.BindPFlag("work-dir", rootCmd.PersistentFlags().Lookup("work-dir"))
	_ = viper.BindPFlag("target-dir", rootCmd.PersistentFlags().Lookup("target-dir"))
	_ = viper.BindPFlag("mount-dir", rootCmd.PersistentFlags().Lookup("mount-dir"))
	_ = viper.BindPFlag("scan-env", rootCmd.PersistentFlags().Lookup("scan-env"))
	_ = viper.BindPFlag("set-env", rootCmd.PersistentFlags().Lookup("set-env"))
	_ = viper.BindPFlag("scan-aws-keys", rootCmd.PersistentFlags().Lookup("scan-aws-keys"))
	_ = viper.BindPFlag("scan-terraform-vars", rootCmd.PersistentFlags().Lookup("scan-terraform-vars"))
	_ = viper.BindPFlag("custom-cmds", rootCmd.PersistentFlags().Lookup("custom-cmds"))
	_ = viper.BindPFlag("init-dagger-with-workdir", rootCmd.PersistentFlags().Lookup(
		"init-dagger-with-workdir"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".k8sgpt.git" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".stiletto")

		_ = viper.SafeWriteConfig()
		//if err != nil {
		//	// Check if error relates to the file already exist.
		//	// If it does, then it's fine, otherwise, exit.
		//	if !os.IsExist(err) {
		//		fmt.Println(err)
		//		os.Exit(1)
		//	}
		//	//os.Exit(1)
		//}
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	addPersistentFlags()

	rootCmd.AddCommand(docker.DockerCmd)
	rootCmd.AddCommand(aws.Cmd)

	_ = rootCmd.MarkFlagRequired("task")
	_ = rootCmd.MarkFlagRequired("workdir")
}
