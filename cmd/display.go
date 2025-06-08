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

// displayCmd represents the display command
var displayCmd = &cobra.Command{
Use:   "display [taskid] ",
Short: "Display one or more tasks with optional filters and customizable output",
Long: `Displays task(s) from your TODO list. You can optionally provide a task ID
to fetch a specific task, or use filtering options to search by category,
title, priority, or status.

Additionally, you can control which fields to display using flags.

By default, all tasks are shown if no task ID or filter is specified.

Supported filters:
  --category     string    Filter by task category
  --title        string    Filter by task title
  --priority     float32   Filter by priority (e.g., 2.0)
  --status       string    Filter by completion status ("true" or "false")

Display controls:
  --show-category | -C      Show the category column
  --show-priorit |-P       Show the priority column
  --show-status  | -S        Show the status column
  --show-all     | -A        Show all columns

Examples:

  display                          # Show all tasks
  display 123                      # Show task with ID 123
  display --category=work         # Show all work-related tasks
  display --priority=1.0 --status=true --show-all
`,

	Run: func(cmd *cobra.Command, args []string) {
	fmt.Println("display called")

	// Get filter flags
	filterCat, _ := cmd.Flags().GetString("category")
	filterTitle, _ := cmd.Flags().GetString("title")
	filterPriority, _ := cmd.Flags().GetFloat32("priority")
	filterStatus, _ := cmd.Flags().GetString("status")

	// Get display flags
	showCat, _ := cmd.Flags().GetBool("show-category")
	showStatus, _ := cmd.Flags().GetBool("show-status")
	showPriority, _ := cmd.Flags().GetBool("show-priority")
	showAll, _ := cmd.Flags().GetBool("show-all")

	// Optional task ID argument
	var id string
	if len(args) >= 1 {
		id = args[0]

		fmt.Println("args: ",args)
	}

	if id != "" {
		fmt.Println("here it entered id branch.")

		// Create GET request manually
		req, err := http.NewRequest("GET", "https://todo-api-s0tq.onrender.com/api/tasks/"+id, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type","application/json")

		utils.AddToken(req)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error fetching task:", err)
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		fmt.Println("response received")

		if resp.StatusCode != http.StatusOK {
			var res map[string]interface{}
			json.Unmarshal(body, &res)
			msg, _ := res["message"].(string)
			fmt.Println("Error:", msg)
			return
		}

		// Parse and display the task


		var res map[string]interface{}
		json.Unmarshal(body, &res)

		task := res["task"]
		var senttask []interface{}
		senttask = append(senttask, task)

		utils.DisplayTaskList(senttask, cmd, showCat, showPriority, showStatus, showAll)
		return
	}

	// Construct the query string based on filter values
	query := ""

	addParam := func(key, value string) {
		if query != "" {
			query += "&"
		}
		query += fmt.Sprintf("%s=%s", key, value)
	}

	if filterCat != "" {
		addParam("category", filterCat)
	}
	if filterTitle != "" {
		addParam("title", filterTitle)
	}
	if filterPriority != -1 {
		// fmt.Println("priority: ",filterPriority)
		addParam("priority", fmt.Sprintf("%.2f", filterPriority))
	}
	if filterStatus=="true" || filterStatus=="false" {

		if filterStatus=="true"{
			addParam("status","true")
		}else{
			addParam("status","false")
		}
	}

	url := "https://todo-api-s0tq.onrender.com/api/tasks/filter"
	if query != "" {
		url += "?" + query
	}

	// Send GET request with query params

	// fmt.Println("url:" ,url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Request creation failed:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	utils.AddToken(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending filter request:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var res map[string]interface{}
		json.Unmarshal(body, &res)
		msg, _ := res["message"].(string)
		fmt.Println("Error:", msg)
		return
	}

	var res map[string]interface{}
	json.Unmarshal(body, &res)

	tasks, _ := res["tasks"].([]interface{})
	utils.DisplayTaskList(tasks, cmd, showCat, showPriority, showStatus, showAll)
},

}

func init() {


	
	taskCmd.AddCommand(displayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// displayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// displayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	displayCmd.Flags().StringP("category", "c", "", "Filter tasks by category")
	displayCmd.Flags().StringP("title", "t", "", "Filter tasks by title")
	displayCmd.Flags().Float32P("priority", "p",-1, "Filter tasks by priority")
	displayCmd.Flags().String("status", "s", "Filter by status: true, false or leave empty to ignore")

	displayCmd.Flags().BoolP("show-category", "C", false, "Show the category of each task in output")
	displayCmd.Flags().BoolP("show-status", "S", false, "Show the Status of each task in output")
	displayCmd.Flags().BoolP("show-priority", "P", false, "Show the priority of each task in output")
	displayCmd.Flags().BoolP("show-all", "A", false, "Show all details of each task")


	// irrelavant of flags set , title is always shown so no flag for show-title

}
