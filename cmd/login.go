/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Adityanath-madishetti/todo-cli/utils"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login [--username name ] [--password password]",
	Short: "Login to your TODO account",
	Long: `Use this command to log in to your TODO CLI.
You can pass username and password as arguments, flags, or interactively.`,
	Example: `  todo login
  todo login --username john --password secret
  todo login john secret`,
		// write ur logic here
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("login command called !!")
		var username, password string

		// First check positional args
		if len(args) >= 2 {
			username = args[0]
			password = args[1]
		} else {
			// Check flags
			usernameFlag, _ := cmd.Flags().GetString("username")
			passwordFlag, _ := cmd.Flags().GetString("password")
			username = usernameFlag
			password = passwordFlag
		}

		// Prompt if still empty
		if username == "" {
			username = utils.PromptInput("Enter Username")
		}

		if password == "" {
			fmt.Printf("Enter password for user %s: ", username)
			hiddenPass, err := utils.SecurePasswordtInput()
			if err != nil {
				fmt.Println("Error reading password:", err)
				return
			}
			password = hiddenPass
		}

		fmt.Println("\n username :=",username,"password:=",password)

		// Send request
		postBody, _ := json.Marshal(map[string]string{
			"username": username,
			"password": password,
		})

		reqBody := bytes.NewReader(postBody)
		// https://go-api-todo.onrender.com
		req, err := http.NewRequest("POST", "https://todo-api-s0tq.onrender.com/api/auth/login", reqBody)
		if err != nil {
			fmt.Println("Request creation failed:", err)
			return
		}



		req.Header.Add("Content-Type", "application/json")

		response, err := http.DefaultClient.Do(req)

		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer response.Body.Close()

		body, _ := io.ReadAll(response.Body)

		fmt.Println("recieved response!")

		if response.StatusCode != http.StatusOK {
			var res map[string]interface{}
			json.Unmarshal(body, &res)
			msg, _ := res["message"].(string)
			fmt.Println("Login failed:", msg)
			return
		}

		// Token extraction
		var res map[string]interface{}
		json.Unmarshal(body, &res)
		token, ok := res["token"].(string)
		if !ok || token == "" {
			fmt.Println("Invalid token received")
			return
		}

		// Save token to ~/.todo-session
		home, _ := os.UserHomeDir()
		sessionPath := filepath.Join(home, ".todo-session")
		err = os.WriteFile(sessionPath, []byte(token), 0600)
		if err != nil {
			fmt.Println("Failed to save session:", err)
			return
		}

		// fmt.Println("✅ Logged in successfully.,token: ",token)
		fmt.Println("✅ Logged in successfull")

	},
}







func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")


	// StringP(name string, shorthand string, value string, usage string) *string


	loginCmd.Flags().StringP("username","u","","this flag helps to provide username")
	loginCmd.Flags().StringP("password","p","","this flag is to set password")

}
