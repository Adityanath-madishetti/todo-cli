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
var updateCmd = &cobra.Command{
	Use:   "set <taskid> [-t title] [-c category] [-p priority]",
	Short: "Update an existing task",
	Long: `Update a task by its ID. You can change its title, category, or priority:

Examples:
  todo-cli task set 123 -t "New Title"
  todo-cli task set 123 -c "Work"
  todo-cli task set 123 -p 2`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskId := args[0]
		if taskId == "" {
			fmt.Println("❌ taskId should not be empty!")
			cmd.Help()
			return
		}

		url := "https://todo-api-s0tq.onrender.com/api/tasks/"

		// Gather flags
		title, _ := cmd.Flags().GetString("title")
		category, _ := cmd.Flags().GetString("category")
		priority, _ := cmd.Flags().GetFloat32("priority")

		// Build request body
		body := make(map[string]interface{})
		var collection []string

		body["taskId"]=taskId

		updates := make(map[string]interface{})
		if title != "" {
			updates["title"] = title
			collection = append(collection, "title")
			
		}
		if category != "" {
			updates["category"] = category
			collection = append(collection, "category")

		}
		if priority !=-1 {

			if priority>3 || priority<1{
				fmt.Println("priority can be 1 or 2 or 3")
				return
			}

			updates["priority"] = priority
			collection = append(collection, "priority")
		}
		body["types"]=collection
		body["updates"]=updates

		// Nothing to update?
		if len(body) == 0 {
			fmt.Println("⚠️ Nothing to update. Provide at least one flag (-t, -c, or -p).")
			return
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			fmt.Println("❌ Failed to encode request body:", err)
			return
		}

		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Println("❌ Failed to create request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		
		err=utils.AddToken(req)
		if err!=nil{
			fmt.Println("error is :",err)
			return
		}
		

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			fmt.Println("❌ Request failed:", err)
			return
		}
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			var res map[string]interface{}
			
			json.Unmarshal(bodyBytes, &res)
			if msg, ok := res["message"].(string); ok {
				fmt.Println("❌ Error:", msg)
			} else {
				fmt.Println("❌ Failed with status:", resp.Status)
			}
			return
		}


		fmt.Println("✅ Task updated successfully.")
	},
}

func init() {
	taskCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringP("title", "t", "", "New title for the task")
	updateCmd.Flags().StringP("category", "c", "", "New category for the task")
	updateCmd.Flags().Float32P("priority", "p", -1, "New priority (e.g. 1, 2)")
}