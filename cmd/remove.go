/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Adityanath-madishetti/todo-cli/utils"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove [taskID]",
	Short: "to remove task based on id",
	Long: `Remove a task from your TODO list using its unique ID.

 examples:
  todo task remove 64c1f7a7bcecd0e0d1246af5
  todo task remove 1234abcd5678efgh90ijklmn

Make sure you are logged in before using this command. It uses your saved token for authentication.`,

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("remove called!!")

		/*
			Because Go only allows := when at least one variable on the left is new, and any reused variables must already exist in that same scope.
		*/


		// get task id as argument
			idarg := args[0]


		//here create a delete request



		uri:="https://todo-api-s0tq.onrender.com/api/tasks/"+idarg

		req,err:=http.NewRequest("DELETE",uri,nil)
		if err != nil {
		fmt.Println("Error creating request:", err)
			return
		}	

		err=utils.AddToken(req)
		if err!=nil{
			fmt.Println("error is :",err)
			return
		}
		
		resp, err := http.DefaultClient.Do(req)



		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}

		body, _ := io.ReadAll(resp.Body)

		defer resp.Body.Close()

		if resp.StatusCode!=http.StatusOK{
			var res map[string]interface{}
			json.Unmarshal(body, &res)
			msg, _ := res["message"].(string)
			fmt.Println("Removing failed:", msg)
			return
		}

		var res map[string]interface{}
		if err := json.Unmarshal(body, &res); err != nil {
			fmt.Println(" Failed to parse signup response:", err)
			return
		}


		fmt.Println(" removed  task successfull")

	},
}

func init() {
	taskCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
