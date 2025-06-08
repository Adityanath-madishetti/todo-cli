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

	"github.com/Adityanath-madishetti/todo-cli/utils"
	"github.com/spf13/cobra"
)

// signUpCmd represents the signUp command
var signUpCmd = &cobra.Command{
	Use:   "signup [--username <name>] [--password <passkey>]",
Short: "Register a new user into the TODO application",
Long: `Registers a new user account for the TODO CLI application.

You can supply the username and password using flags. If either flag is omitted,
the command will interactively prompt the user to enter the missing information.

The signup process sends your credentials to the backend API for registration.

Examples:

  signup --username alice --password mysecretpass
  signup                       # Prompts for username and password
`,

		Run: func(cmd *cobra.Command, args []string) {
		var username, password string

		// Step 1: Check positional args
		if len(args) >= 2 {
			username = args[0]
			password = args[1]
		} else {
			// Step 2: Check flags
			usernameFlag, _ := cmd.Flags().GetString("username")
			passwordFlag, _ := cmd.Flags().GetString("password")
			username = usernameFlag
			password = passwordFlag
		}

		// Step 3: Prompt if still empty
		if username == "" {
			username = utils.PromptInput("Choose a Username")
		}

		if password == "" {
			fmt.Printf("Choose a Password for user %s: ", username)
			hiddenPass, err := utils.SecurePasswordtInput()
			if err != nil {
				fmt.Println(" Error reading password:", err)
				return
			}
			password = hiddenPass
		}

		// Step 4: Send signup request
		postBody, _ := json.Marshal(map[string]string{
			"username": username,
			"password": password,
		})

		reqBody := bytes.NewReader(postBody)
		req, err := http.NewRequest("POST", "http://localhost:8080/api/auth/register", reqBody)
		if err != nil {
			fmt.Println(" Failed to create signup request:", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(" Failed to send signup request:", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			var res map[string]interface{}
			json.Unmarshal(body, &res)
			msg, _ := res["message"].(string)
			fmt.Println(" Signup failed:", msg)
			return
		}

		var res map[string]interface{}
		if err := json.Unmarshal(body, &res); err != nil {
			fmt.Println(" Failed to parse signup response:", err)
			return
		}

		userInfo, _ := res["userinfo"].(map[string]interface{})

		fmt.Println(" Signup successful! Welcome,", userInfo["name"])
		fmt.Println(" unique Id:", userInfo["userId"])
		fmt.Println(" You can now log in using:")
		fmt.Println("   todo login -u", username)

},

}

func init() {
	rootCmd.AddCommand(signUpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// signUpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// signUpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")


	signUpCmd.Flags().StringP("username","u","","this flag helps to provide username")
	signUpCmd.Flags().StringP("password","p","","this flag is to set password")

}
