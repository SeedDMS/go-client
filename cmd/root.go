/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
    "seeddms.org/seeddms/client/config"
)

var (
    cfgFile string
    profile string
    url string
    apikey string
    user string
    password string
    quiet bool

    cfg *config.Config
)

// folderid is used by cmd 'ls', 'upload', and 'autoupload'
var (
    folderid int
)

var (
    version string="0.0.2"
    comment string
    keep bool
)

// rootCmd represents the base command. It also reads the configuration
// file
var rootCmd = &cobra.Command{
	Use:   "seeddms-client",
	Short: "Run various operations on a SeedDMS server",
	Long: `This program connects to a SeedDMS server and can list,
download and upload documents. It uses the RestAPI provided by SeedDMS.

The program can be fully configured with command line options but
usually it is much more convenient to write the basic configuration
into a file named .seeddms-client.yaml in your home directory. This
also prevents sensitive data like a password or api key to enter your
shell's history. Such a file may have different sections. Each
sections contains the configuration for a specific profile. Such a
profile usually addresses a single SeedDMS installation.
There must be at least one section named 'default'.  Further sections
may be named arbitrary. Use the flag --profile to select the
configuration of a specific section. Below is an example of the
'default' section.

default:
  Url: https://domain/restapi/index.php
  ApiKey: <your secret api key>
  Upload:
    Folder: 8345
  Ls:
    Folder: 1

It sets the url and credentials of your SeedDMS installation and
also different folders for the 'upload' and 'ls' commands.

With the above configuration the command

seeddms-client ls

will list the subfolders and documents of the folder with id 1
`,
    Version: version,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
    // Run before each command to load the configuration
    PersistentPreRun: func(cmd *cobra.Command, args []string) {

        // Take configuration from config file and set parameters if not
        // already set via command line
        var err error
        cfg, err = config.NewConfig(cfgFile, profile)
        if err != nil {
            cmd.Printf("Error reading configuration file: %s\n", err)
            os.Exit(1)
        }
        if url != "" {
            cfg.Url = url
        }
        if user != "" {
            cfg.User = user
        }
        if password != "" {
            cfg.Password = password
        }
        if apikey != "" {
            cfg.ApiKey = apikey
        }

        // Parse from parameters
        // insure all necessary parameters are defined
        var msg string
        if cfg.Url == "" {
            msg += "- your host URL\n"
        }
        if cfg.ApiKey == "" && cfg.User == "" {
            msg += " - the apikey or login of an existing user\n"
        }
        if cfg.ApiKey == "" && cfg.Password == "" {
            msg += " - the apikey or password of an existing user\n"
        }

        if len(msg) > 0 {
            cmd.Println("Could not run program, missing arguments:\n", msg,
                "\nYou might also directly provide the relative path to a configuration file. See in-line help for further details.")
            os.Exit(1)
        }

        return
    },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file, defaults to '~/.seeddms-client.yaml' if not set")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "select section in configuration, defaults to 'default' if not set")
    rootCmd.PersistentFlags().StringVar(&url, "url", "", "URL to REST Api")
    rootCmd.PersistentFlags().StringVarP(&user, "user", "u", "admin", "User name")
    rootCmd.PersistentFlags().StringVarP(&apikey, "apikey", "k", "", "Api key")
    rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Be quiet")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


