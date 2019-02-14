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
		"version",
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
		Transform: survey.TransformString(linterSubject),
	},
	{
		Name:      "body",
		Prompt:    &survey.Multiline{Message: "Type in the body"},
		Validate:  validator(320),
		Transform: survey.TransformString(linterBody),
	},
	{
		Name:      "foot",
		Prompt:    &survey.Multiline{Message: "Type in the foot"},
		Validate:  validator(50),
		Transform: survey.TransformString(linterFoot),
	},
}

var cs = &survey.Confirm{Message: "Is everything OK? Continue?"}
var ts = &survey.Select{
	Message: "Choose tag level",
	Options: []string{
		"patch",
		"minor",
		"major",
	},
	Default: "patch",
}

// Prompt function assignes user input to Message struct
func Prompt() Message {
	// Perform the questions
	err := survey.Ask(qs, &msg)
	checkInterrupt(err)
	Info("\nCommit message is:\n%s\n", msg.String())
	confirm := false
	err = survey.AskOne(cs, &confirm, nil)
	checkInterrupt(err)
	if !confirm {
		log.Printf("Commit flow interrupted by user\n")
		os.Exit(1)
	}
	return msg
}

// TagPrompt prompting tag version level to upgrade
func TagPrompt() string {
	level := "patch"
	err := survey.AskOne(ts, &level, nil)
	checkInterrupt(err)
	return level
}

// PromptComfirm is a common function to ask confirm before some action
func PromptComfirm(msg string) bool {
	confirm := false
	err := survey.AskOne(&survey.Confirm{Message: msg}, &confirm, nil)
	checkInterrupt(err)
	return confirm
}

func linterSubject(s string) string {
	// Remove all leading and trailing white spaces
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "...")
	// Then strings.Title the first word in string
	flds := strings.Fields(s)
	flds[0] = strings.Title(flds[0])
	return strings.Join(flds, " ")
}

func linterBody(s string) string {
	// Remove all leading and trailing white spaces
	s = strings.TrimSpace(s)
	// Split string to lines
	strs := strings.Split(s, "\n")
	// Then strings.Title the first word in string
	flds := strings.Fields(strs[0])
	flds[0] = strings.Title(flds[0])
	strs[0] = strings.Join(flds, " ")
	// Return glued lines
	return strings.Join(strs, "\n")
}

func linterFoot(s string) string {
	s = strings.TrimSpace(s)
	// Split string to lines
	strs := strings.Split(s, "\n")
	for i := len(strs); i > 0; i-- {
		if strings.HasPrefix(strs[i-1], "* ") {
			strs[i-1] = strings.TrimPrefix(strs[i-1], "* ")
		}
		strs[i-1] = linterSubject(strs[i-1])
		strs[i-1] = "* " + strs[i-1]
	}
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

func checkInterrupt(err error) {
	if err != nil {
		if err != terminal.InterruptErr {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
		}
		log.Printf("interrupted by user\n")
		os.Exit(1)
	}
}
