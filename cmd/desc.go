/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Adityanath-madishetti/todo-cli/utils"
	"github.com/spf13/cobra"
)

// descCmd represents the desc command
var descCmd = &cobra.Command{
	Use:   "set-desc <taskid> [text | -i <file>]",
Short: "Set or update the description of a task",
Long: `Updates the description of a task by its task ID.

You can provide the description directly as a text argument,
or load it from a file using the '--input' or '-i' flag.

One of [text argument] or [-i file] is required.
This command sends a PUT request to the backend with the updated description.
Authentication is handled via a session token file (~/.todo-session).

Examples:

  set-desc 12345 "This task needs to be done before Friday."
  set-desc mytaskid -i description.txt
`,

	
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("desc called")

		taskId:=args[0] 

		if taskId==""{
			fmt.Println("taskId should not be empty!!")
			cmd.Help()
			return
		}

		inpFile,_:=cmd.Flags().GetString("input")
		var text string
		if inpFile!=""{  // flag is present then go for it

			//open file collect info and put in string 


			fd,err:=os.Open(inpFile)

			if err!=nil{
				fmt.Println(" error in openeing file :"+err.Error())
				return
			}	
				defer fd.Close()


			sc:=bufio.NewScanner(fd)	

			for sc.Scan(){
				text+=sc.Text()+"\n"
				
			}

				if err := sc.Err(); err != nil {
					fmt.Println("Error reading file:", err)
					return
				}
		}else{

			if len(args)>=2{
				text=args[1]
			}else{
				fmt.Println("text is not provided ")
				return
			}
		}

		//now encode text in body of put

		url:="http://localhost:8080/api/tasks/description/"+taskId

		
		jsonbytes,err:=json.Marshal(map[string]string{"text":text})

		if err!=nil{
			fmt.Println("error in marshaling te text ")
			return
		}

			bodyreader:=bytes.NewReader(jsonbytes)

		req,err:=http.NewRequest(http.MethodPut,url,bodyreader)
		if err != nil {
			fmt.Println(" Failed to create Task dexcription Update request:", err)
			return
		}	
		req.Header.Set("Content-Type","application/json")
		utils.AddToken(req)


		response,err:=http.DefaultClient.Do(req)

		if err!=nil{
			fmt.Println("error in sending teh request")
			return
		}

		defer response.Body.Close()

		bodyBytes, _ := io.ReadAll(response.Body)

		if response.StatusCode!=http.StatusOK{
			//unmarshal the body
			var res map[string]interface{}
			json.Unmarshal(bodyBytes, &res)
			msg, _ := res["message"].(string)
			fmt.Println("Task description update  failed: ", msg)
			return
		}

		var res map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &res); err != nil {
			fmt.Println(" Failed to parse signup response:", err)
			return
		}

		fmt.Println("succesfully changed the description :", res["message"])

	},
}

func init() {
	taskCmd.AddCommand(descCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// descCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// descCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	descCmd.Flags().StringP("input","i","",`setting the description to file eg: todo task set-desc -i="input.txt" `)


}
