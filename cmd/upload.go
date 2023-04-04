/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
    "os"
    "path"
    "seeddms.org/seeddms/apiclient"
	"github.com/spf13/cobra"
    "github.com/0xAX/notificator"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload file and create new document",
	Long: `The upload command will upload any local file into SeedDMS and create
a new document.
`,
	Run: func(cmd *cobra.Command, args []string) {
        file := ""
        var err error
        if len(args) > 0 {
            file = path.Base(args[0])
        } else {
            fmt.Println("Need to pass a file to be uploaded")
            return
        }
        c := apiclient.Connect(cfg.Url, cfg.ApiKey)

        _, err = c.Login(cfg.User, cfg.Password)
        if err != nil {
            fmt.Println("Failed to login")
            return
        }
        if comment == "" {
            comment = cfg.Upload.Comment
        }
        extraParams := map[string]string{
            "name":        file,
            "filename":    file,
            "comment":     comment,
        }
        f, err := os.Open(args[0])
        if err != nil {
            fmt.Println("Cannot open file")
            return
        }
        defer f.Close()

        if folderid == 0 {
            folderid = cfg.Upload.Folder
        }
        if folderid == 0 {
            fmt.Println("Target folder not set")
            return
        }
        res, err := c.Upload(f, extraParams, folderid)
        if err != nil {
            fmt.Println(err)
            return
        }
        if !quiet {
            var notify *notificator.Notificator
            notify = notificator.New(notificator.Options{
                DefaultIcon: "/usr/share/seeddms-client/icons/warning.png",
                AppName:     "My test App",
            })

            notify.Push("Document uploaded", fmt.Sprintf("Document with Id=%d uploaded", res.Data.Id), "/usr/share/seeddms-client/icons/warning.png", notificator.UR_CRITICAL)
        }
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    uploadCmd.PersistentFlags().IntVar(&folderid, "folder", 0, "Id of parent folder")
    uploadCmd.PersistentFlags().StringVar(&comment, "comment", "", "Comment of new document")
}
