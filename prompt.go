package gitcomm

import (
	"fmt"
	"log"
	"os"
	"strings"

	survey "gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

var (
	msg = Message{
		Type:    "feat",
		Subject: "",
		Body:    "",
		Foot:    "",
	}
	types = []string{
		"feat",
		"fix",
		"docs",
		"style",
		"refactor",
		"test",
		"chore",
	}
)

// the questions to ask
var qs = []*survey.Question{
	{
		Name: "type",
		Prompt: &survey.Select{
			Message: "Choose a type(<scope>)",
			Options: types,
			Default: "feat",
		},
	},
	{
		Name:      "subject",
		Prompt:    &survey.Input{Message: "Type in the subject"},
		Validate:  validator(72),
		Transform: survey.TransformString(linter),
	},
	{
		Name:      "body",
		Prompt:    &survey.Multiline{Message: "Type in the body"},
		Validate:  validator(320),
		Transform: survey.TransformString(linter),
	},
	{
		Name:      "foot",
		Prompt:    &survey.Multiline{Message: "Type in the foot"},
		Validate:  validator(50),
		Transform: survey.TransformString(linter),
	},
}

var cs = &survey.Confirm{Message: "Is everything OK? Continue?"}

// Prompt function assignes user input to Message struct
func Prompt() Message {
	// Perform the questions
	err := survey.Ask(qs, &msg)
	if err != nil {
		if err != terminal.InterruptErr {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
		}
		os.Exit(1)
	}
	Info("\nCommit message is:\n%s\n", msg.String())
	confirm := false
	survey.AskOne(cs, &confirm, nil)
	if !confirm {
		log.Printf("Commit flow breaked by user\n")
		os.Exit(1)
	}
	return msg
}

func linter(s string) string {
	// Remove all leading and trailing white spaces
	s = strings.TrimSpace(s)
	// Split string to lines
	strs := strings.Split(s, "\n")
	strs[0] = strings.TrimSuffix(strs[0], "...")
	// Then strings.Title the first word in string
	flds := strings.Fields(strs[0])
	flds[0] = strings.Title(flds[0])
	strs[0] = strings.Join(flds, " ")
	// Return glued lines
	return strings.Join(strs, "\n")
}

func validator(n int) func(val interface{}) error {
	return func(val interface{}) error {
		// since we are validating an Input, the assertion will always succeed
		if str, ok := val.(string); !ok || str == "" || len(str) > n {
			return fmt.Errorf("This response cannot be longer than %d characters", n)
		}
		return nil
	}
}
