package utils

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"os"
)

func PromptInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func SecurePasswordtInput() (string,error){
	 fmt.Print("Enter Password: ")
	  bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	   if err != nil {
		return "",err
    }
	password := string(bytePassword)
	return password,nil
}

func AddToken(req *http.Request) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine home directory: %w", err)
	}

	sessionFile := filepath.Join(home, ".todo-session")

	content, err := os.ReadFile(sessionFile)
	if errors.Is(err, os.ErrNotExist) {
		// Don't print anything here — let the caller handle the "not logged in" case
		fmt.Println("here in token validation")
		return errors.New("not logged in")
	}
	if err != nil {
		return fmt.Errorf("could not read session file: %w", err)
	}

	token := strings.TrimSpace(string(content))
	if token == "" {
		return errors.New("not logged in: empty token")
	}

	req.Header.Add("Authorization", "Bearer "+token)
	return nil
}



// DisplayTaskList prints a list of tasks in a well-styled table
func DisplayTaskList(tasks []interface{}, cmd *cobra.Command, showCat, showPriority, showStatus, showAll bool) {
    tw := table.NewWriter()
    tw.SetOutputMirror(os.Stdout)

    cyan := color.New(color.FgCyan).SprintFunc()
    magenta := color.New(color.FgMagenta).SprintFunc()
    green := color.New(color.FgGreen).SprintFunc()
    red := color.New(color.FgRed).SprintFunc()
    yellow := color.New(color.FgYellow).Add(color.Bold).SprintFunc()

    header := table.Row{cyan("#"), cyan("ID"), cyan("Title")}
    if showAll || showCat {
        header = append(header, magenta("Category"))
    }
    if showAll || showStatus {
        header = append(header, magenta("Status"))
    }
    if showAll || showPriority {
        header = append(header, magenta("Priority"))
    }
    if showAll {
        header = append(header, magenta("Created At"), magenta("Completed At"))
    }

    tw.AppendHeader(header)
    tw.SetStyle(table.StyleColoredBright)
    tw.Style().Box.PaddingLeft = "  "
    tw.Style().Box.PaddingRight = "    "
    tw.SetColumnConfigs([]table.ColumnConfig{
        {Number: 1, WidthMin: 4},
		{Number: 2,WidthMin:15},
        {Number: 3, WidthMin: 25},
    })

    for i, t := range tasks {
		var emptyRow table.Row
        taskMap := t.(map[string]interface{})
        row := table.Row{i + 1, taskMap["taskId"], taskMap["title"]}

        if showAll || showCat {
            row = append(row, taskMap["category"])
			emptyRow=append(emptyRow, "")
        }
        if showAll || showStatus {
            status := taskMap["completed"]
            if status == true {
                row = append(row, green("✅"))
				emptyRow=append(emptyRow, "")

            } else {
                row = append(row, red("❌"))
				emptyRow=append(emptyRow, "")

            }
        }
        if showAll || showPriority {
            row = append(row, taskMap["priority"])
			emptyRow=append(emptyRow, "")

        }
        if showAll {
            row = append(row,
                formatTime(taskMap["creationTime"]),
                formatTime(taskMap["endTime"]),
				

            )

			emptyRow=append(emptyRow, "")
			emptyRow=append(emptyRow, "")
        }
        tw.AppendRow(row)
	   // Add a gap (empty row)
	   tw.AppendRow(emptyRow)

    }

    // Render the table (no footer)
    fmt.Println()
    tw.Render()
    fmt.Println()

now := time.Now().Format(time.RFC3339) // current time as RFC3339 string
footerRaw := "Generated on " + formatTime(now) // your function converts to IST
footerText := yellow(footerRaw)

    totalWidth := 80

    // Calculate left padding to center text
    rawLen := len(footerRaw)
    padLeft := (totalWidth - rawLen) / 2
    if padLeft < 0 {
        padLeft = 0
    }

    fmt.Println(strings.Repeat(" ", padLeft) + footerText)
}

// formatTime safely converts interface{} time values to readable string
func formatTime(t interface{}) string {
	if t == nil {
		return "-"
	}

	// Load IST timezone
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		// fallback in case loading location fails
		loc = time.FixedZone("IST", 5.5*3600)
	}

	switch v := t.(type) {
	case string:
		// Try to format ISO8601 time string
		parsed, err := time.Parse(time.RFC3339, v)
		if err == nil {
			// Convert to IST
			parsedIST := parsed.In(loc)
			return parsedIST.Format("02 Jan 2006 15:04")
		}
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}



func DisplayDescription(desc string) {
	wrapWidth := 70
	lines := strings.Split(text.WrapSoft(desc, wrapWidth), "\n")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Styling
	t.SetStyle(table.StyleLight)
	t.Style().Box = table.StyleBoxRounded
	t.Style().Options.DrawBorder = true
	t.Style().Options.SeparateRows = false

	// Dark black text on white background
t.Style().Color.Row = text.Colors{text.FgBlack, text.BgHiWhite, text.Bold}
	t.Style().Color.Border = text.Colors{text.FgCyan}

	for _, line := range lines {
		centered := centerLine(line, wrapWidth)
		t.AppendRow(table.Row{centered})
	}

	t.Render()
}

func centerLine(line string, width int) string {
	padding := width - len(line)
	if padding <= 0 {
		return line
	}
	left := padding / 2
	right := padding - left
	return strings.Repeat(" ", left) + line + strings.Repeat(" ", right)
}