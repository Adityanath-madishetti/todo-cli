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
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
		Use:   "add <[title] [category] [priority] |  [-t=<title> -c=<category> -p=<priority>] >",
	Short: "Add a new task with title, category, and priority",
	Long: `Adds a new task to your TODO list by sending a POST request to the backend API.

You can provide the task details either as positional arguments or using flags:
  - Title: a short name or label for the task
  - Category: logical group like "work", "personal", etc.
  - Priority: a numeric value indicating task importance (e.g., 1.0, 2.5)

If all three arguments are not passed, the command will try to use flags:
  --title     string   task title
  --category  string   task category
  --priority  float32  task priority


Examples:

  add "Buy groceries" "personal" 2.0
  add --title="Fix bug" --category="work" --priority=1.5
`,
	
Run: func(cmd *cobra.Command, args []string) {
	fmt.Println(" Add task called!")

	var title, category string
	var priority float32

	if len(args) >= 3 {
		title = args[0]
		category = args[1]
		temp, err := strconv.ParseFloat(args[2], 32)
		if err != nil {
			fmt.Println(" Invalid priority value. Must be a number.")
			return
		}
		priority = float32(temp)
	} else {
		// collect from flags
		var err error
		title, _ = cmd.Flags().GetString("title")
		category, _ = cmd.Flags().GetString("category")
		priority, err = cmd.Flags().GetFloat32("priority")
		if err != nil {
			fmt.Println(" Error parsing priority flag:", err)
			return
		}
	}

	if title == "" || category == "" || priority == 0 {
		fmt.Println(" Title, category, and priority are required.")
		return
	}

	postBody, _ := json.Marshal(map[string]interface{}{
		"title":    title,
		"category": category,
		"priority": priority,
	})
	reqBody := bytes.NewReader(postBody)

	req, err := http.NewRequest("POST", "https://todo-api-s0tq.onrender.com/api/tasks/", reqBody)
	if err != nil {
		fmt.Println(" Request creation failed:", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	//------------------------------------
	// Read token from ~/.todo-session
	//------------------------------------
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(" Could not get home directory:", err)
		return
	}

	sessionFile := filepath.Join(home, ".todo-session")
	content, err := os.ReadFile(sessionFile)
	if err != nil {
		fmt.Println(" Failed to read session file:", err)
		return
	}
	token := strings.TrimSpace(string(content))
	if token == "" {
		fmt.Println(" Token is empty in session file.")
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	// Send HTTP request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(" Error sending request:", err)
		return
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)
	fmt.Println(" Received response!")

	if response.StatusCode != http.StatusOK {
		var res map[string]interface{}
		json.Unmarshal(body, &res)
		msg, _ := res["message"].(string)
		fmt.Println(" Adding task failed:", msg)
		return
	}

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println(" Failed to parse response:", err)
		return
	}

	fmt.Println(" Message from server:", res["message"])
},

}

func init() {
	taskCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addCmd.Flags().StringP("title","t","none","todo  add task -t=<title>")
	addCmd.Flags().StringP("category","c","none","todo  add task -c=<category>")
	addCmd.Flags().Float32P("priority","p",1,"todo  add task -p=<priorityvalue>")


}
