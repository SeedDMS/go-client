/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
    "strconv"
    "seeddms.org/seeddms/apiclient"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List documents and folders",
	Long: `The ls command list all documents and sub folders of a
given folder.

The command takes one optional parameter containing the id of the
folder. It can be provided by the option --folder or as another
argument. If that id is not passed on the command line it will be
read from the configuration file.

Each line of the output consists of three columns:
1. a symbol for a document or folder
2. the id of the folder or document
3. the name of the folder or document
`,
	RunE: func(cmd *cobra.Command, args []string) error {
        rootid := 1
        if folderid != 0 {
            rootid = folderid
        } else if cfg.Ls.Folder != 0 {
            rootid = cfg.Ls.Folder
        }

        var err error
        if len(args) > 0 {
            rootid, err = strconv.Atoi(args[0])
        }
        c := apiclient.Connect(cfg.Url, cfg.ApiKey)

        _, err = c.Login(cfg.User, cfg.Password)
        if err != nil {
//            fmt.Printf("Failed to login: %s\n", err)
            return err
        }

        res, err := c.Children(rootid)
        if err != nil {
            return err
        }

        for _, row := range res.Data {
            if row.Objtype == "folder" {
                fmt.Printf("%s", "📁")
            } else {
                fmt.Printf("%s", "📄")
            }
            fmt.Printf(" %5d %s\n", row.Id, row.Name);
        }
        return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    lsCmd.PersistentFlags().IntVar(&folderid, "folder", 0, "Id of folder")
}
