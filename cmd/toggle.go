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

// toggleCmd represents the toggle command
var toggleCmd = &cobra.Command{
	Use:   "toggle <task-id>",
Short: "Toggle the completion status of a task",
Long: `Toggles the status of the specified task by its ID.

If the task is currently marked as incomplete, it will be marked as complete,
and vice versa. This is useful for quickly updating the task's status without
editing the full task details.

Requires a valid task ID as a positional argument. No flags are needed.

Examples:

  toggle 12345     # Toggles the status of the task with ID 12345
`,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("toggle called")
		idarg := args[0]


		uri:= "https://todo-api-s0tq.onrender.com/api/tasks/"

		putbody,_:=json.Marshal(map[string]interface{}{
		"types": []string{"toggle"},
		"taskId":idarg,
		"updates":map[string]interface{}{"toggle":true},
	})	

				reqBody := bytes.NewReader(putbody)


		req,err:=http.NewRequest("PUT",uri,reqBody)
		if err != nil {
		fmt.Println("Error creating request:", err)
			return
		}	

		req.Header.Set("Content-Type","application/json")

		
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
			fmt.Println("Toggle update failed:", msg)
			return
		}

		var res map[string]interface{}
		if err := json.Unmarshal(body, &res); err != nil {
			fmt.Println(" Failed to parse togggle response:", err)
			return
		}


		fmt.Println(" Toggled  task successfull")

	},
}

func init() {
	taskCmd.AddCommand(toggleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// toggleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// toggleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
