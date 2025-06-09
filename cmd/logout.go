/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from your current session",
	Long: `This command logs you out from the current session by removing the locally stored token.
No credentials will be stored after this.`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Could not determine home directory.")
			return
		}

		sessionFile := filepath.Join(home, ".todo-session")
		if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
			fmt.Println("You are already logged out.")
			return
		}

		err = os.Remove(sessionFile)
		if err != nil {
			fmt.Println("Error logging out:", err)
			return
		}

		fmt.Println("✅ Logged out successfully.")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
