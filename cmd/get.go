/*
Copyright © 2022 Uwe Steinmann <uwe@steinmann.cx>

*/
package cmd

import (
    "fmt"
    "strconv"
    "io"
    "os"
    "seeddms.org/seeddms/apiclient"
	"github.com/spf13/cobra"
)

var (
    documentid int
    filename string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Download a document",
	Long: `This command downloads the latest version of a document
and writes to local file. The filename is either specified with the
option --filename or is taken from the original filename of the
SeedDMS document version.

The command takes an optional parameter containing the id of the
document. It can be provided by the option --document or as another
argument. If that id is not passed on the command line it will be
read from the configuration file.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
        docid := 0
        var err error
        if documentid != 0 {
            docid = documentid
        }
        if len(args) > 0 {
            docid, err = strconv.Atoi(args[0])
        }
        c := apiclient.Connect(cfg.Url, cfg.ApiKey)

        _, err = c.Login(cfg.User, cfg.Password)
        if err != nil {
            fmt.Println("Failed to login")
            return err
        }

        res, err := c.Document(docid)
        if err != nil {
            fmt.Println("Failed to get metadata of document", docid)
            return err
        }
        fmt.Println(res.Data.Origfilename);

        body, err := c.Content(docid)
        if err != nil {
            fmt.Println("Failed to get content of document", docid)
            return err
        }

        if(filename == "") {
            filename = res.Data.Origfilename
        }
        destination, err := os.Create(filename)
        if err != nil {
            fmt.Println("Failed to create file")
            return err
        }
        nBytes, err := io.Copy(destination, body)
        if err != nil {
            fmt.Println("Failed to write file")
            return err
        }
        fmt.Println(nBytes, "Bytes written to", filename)

//        content, err := ioutil.ReadAll(body)
       // io.Copy(os.Stdout, body)
       return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
    getCmd.PersistentFlags().IntVar(&documentid, "document", 0, "Id of document")
    getCmd.PersistentFlags().StringVar(&filename, "filename", "", "Name of file to save document content")
}
