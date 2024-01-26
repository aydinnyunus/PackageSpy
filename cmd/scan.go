/*
Copyright © 2024 Yunus AYDIN aydinnyunus@gmail.com
*/
package cmd

import (
	"fmt"
	"github.com/spf13/secretScanner/utils/npm"
	"github.com/spf13/secretScanner/utils/pypi"

	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println(cmd.Flag("username").Value)
		if cmd.Flags().Lookup("pypi").Changed {
			fmt.Println("pypi called")
			fmt.Println(cmd.Flag("username").Value.String())
			pypi.DownloadAllPyPIPackages(cmd.Flag("username").Value.String())
		} else if cmd.Flags().Lookup("npm").Changed {
			fmt.Println("npm called")
			fmt.Println(cmd.Flag("username").Value.String())
			npm.DownloadAllNpmPackages(cmd.Flag("username").Value.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	scanCmd.Flags().StringP("username", "u", "", "Username for Package Manager")
	scanCmd.Flags().BoolP("pypi", "p", false, "is PyPI")
	scanCmd.Flags().BoolP("npm", "n", false, "is NPM")
	/*
		scanCmd.Flags().BoolP("rubygems", "r", false, "is RubyGems")
	*/

}
