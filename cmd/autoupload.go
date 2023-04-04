/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
    "log"
	"fmt"
    "os"
    "path"
    "time"
    "seeddms.org/seeddms/apiclient"
	"github.com/spf13/cobra"
    "github.com/0xAX/notificator"
    "github.com/fsnotify/fsnotify"
)

// uploadCmd represents the upload command
var autouploadCmd = &cobra.Command{
	Use:   "autoupload",
	Short: "Monitors local directory and uploads new files",
	Long: `This command checks a given local directory for new files until
the program is killed. Each new file is uploaded into SeedDMS.

This command could be used with scanner software like scanservjs to
monitor the output directory and save any new scan into SeedDMS.
`,
	Run: func(cmd *cobra.Command, args []string) {
        scandir := ""
        var err error
        if len(args) > 0 {
            scandir = args[0]
        } else {
            fmt.Println("Need to pass the directory to be monitored")
            return
        }
        c := apiclient.Connect(cfg.Url, cfg.ApiKey)

        _, err = c.Login(cfg.User, cfg.Password)
        if err != nil {
            fmt.Println("Failed to login")
            return
        }

        if comment == "" {
            comment = cfg.Autoupload.Comment
        }

        if folderid == 0 {
            folderid = cfg.Autoupload.Folder
        }
        if folderid == 0 {
            fmt.Println("Target folder not set")
            return
        }

        // The watch loop
        watcher, err := fsnotify.NewWatcher()
        if err != nil {
            fmt.Println(err)
            return
        }
        defer watcher.Close()

        done := make(chan bool)
        go func() {
            for {
                select {
                case event, ok := <-watcher.Events:
                    if !ok {
                        return
                    }
                    //log.Println("event:", event)
                    if event.Op&fsnotify.Create == fsnotify.Create {
                        //little pause to ensure the write operation is finished
                        time.Sleep(1 * time.Second)
                        log.Println("New file to upload:", path.Base(event.Name))

                        extraParams := map[string]string{
                            "name":        path.Base(event.Name),
                            "filename":    path.Base(event.Name),
                            "comment":     comment,
                        }

                        f, err := os.Open(event.Name)
                        if err != nil {
                            log.Println("Cannot open file:", path.Base(event.Name))
                            return
                        }
                        defer f.Close()

                        res, err := c.Upload(f, extraParams, folderid)
                        f.Close()
                        if err != nil {
                            log.Println(err)
                        } else {
                            log.Println(fmt.Sprintf("Document with Id=%d uploaded", res.Data.Id))
                            if keep != true {
                                log.Println("Removing file:", path.Base(event.Name))
                                e := os.Remove(event.Name)
                                if e != nil {
                                    log.Fatal(e)
                                }
                            } else {
                                log.Println("Keeping file:", path.Base(event.Name))
                            }
                            if !quiet {
                                var notify *notificator.Notificator
                                notify = notificator.New(notificator.Options{
                                    DefaultIcon: "/usr/share/seeddms-client/icons/warning.png",
                                    AppName:     "My test App",
                                })

                                notify.Push("Document uploaded", fmt.Sprintf("Document with Id=%d uploaded", res.Data.Id), "/usr/share/seeddms-client/icons/warning.png", notificator.UR_CRITICAL)
                            }
                        }
                    }
                case err, ok := <-watcher.Errors:
                    if !ok {
                        return
                    }
                    log.Println("error:", err)
                    return
                }
            }
        }()

        //err = watcher.Add("/home/docker_data/scanservjs/output")
        log.Println("Watching directory:", scandir)
        err = watcher.Add(scandir)
        if err != nil {
            log.Println("error: ", err)
            return
        }
        <-done

	},
}

func init() {
	rootCmd.AddCommand(autouploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    autouploadCmd.PersistentFlags().IntVar(&folderid, "folder", 0, "Id of parent folder")
    autouploadCmd.PersistentFlags().StringVar(&comment, "comment", "", "Comment of new document")
	autouploadCmd.PersistentFlags().BoolVar(&keep, "keep", false, "Keep the uploaded file")
}
