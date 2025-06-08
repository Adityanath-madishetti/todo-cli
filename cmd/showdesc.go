/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Adityanath-madishetti/todo-cli/utils"
	"github.com/spf13/cobra"
)

// showdescCmd represents the showdesc command
var showdescCmd = &cobra.Command{
	Use:   "show-desc <taskid> [-f <file>]",
	Short: "Show the detailed description of a specific task",
	Long: `Fetch and display the description of a task using its task ID.

You can choose to print the description to the terminal or save it to a file
by providing the optional '--output' (or '-f') flag. The command will always
print the server message (success or error), and if successful, the description
text will be either shown or saved based on your choice.

Examples:

  show-desc 12345
  show-desc 12345 -f taskdesc.txt
  show-desc mytaskid --output=result.txt`,

	Args: cobra.ExactArgs(1),

	
	Run: func(cmd *cobra.Command, args []string) {
	fmt.Println("showdesc called")

	if len(args) < 1 {
		fmt.Println("taskId is required as an argument!")
		cmd.Help()
		return
	}

	taskId := args[0]
	if taskId == "" {
		fmt.Println("taskId should not be empty!")
		cmd.Help()
		return
	}

	opfile, err := cmd.Flags().GetString("output")
	if err != nil {
		fmt.Println("Error getting output flag:", err)
		return
	}

	// If output file is specified, test if it's creatable before sending the request
	var file *os.File
	if opfile != "" {
		file, err = os.Create(opfile)
		if err != nil {
			fmt.Println("Error: Cannot create output file:", err)
			return
		}
		defer file.Close()
	}

		url := "https://todo-api-s0tq.onrender.com/api/tasks/description/"+taskId

	
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		utils.AddToken(req)


		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error sending GET request:", err)
			return
		}
		defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Parse response as JSON
	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return
	}

	// Always show msg
	// if msg, ok := res["message"]; ok {
	// 	fmt.Println("Server message:", msg)
	// } else {
	// 	fmt.Println("Warning: 'msg' field not found in response")
	// }

	if resp.StatusCode == http.StatusOK {
		text, ok := res["text"].(string)
		if !ok {
			fmt.Println("No valid 'text' field found in response")
			return
		}

		if opfile != "" {
			if _, err := file.WriteString(text); err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
			fmt.Println("Description written to", opfile)
		} else {
			if(text!=""){
				utils.DisplayDescription(text)
			}else{
				utils.DisplayDescription("No Description Available to this task")
			}
		}
	} else {
		fmt.Printf("Request failed :%s ",res["message"])
	}
},

}

func init() {
	taskCmd.AddCommand(showdescCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showdescCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showdescCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")


		showdescCmd.Flags().StringP("output","o","",`setting the description to file eg: todo task show-desc af123.. -o="output.txt" `)

}
