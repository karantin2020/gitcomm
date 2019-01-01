package gitcomm

import (
	"fmt"
	"strings"

	survey "gopkg.in/AlecAivazis/survey.v1"
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
		Validate:  validator(50),
		Transform: survey.TransformString(linter),
	},
	{
		Name:      "body",
		Prompt:    &survey.Multiline{Message: "Type in the body"},
		Validate:  validator(72),
		Transform: survey.TransformString(linter),
	},
	{
		Name:      "footer",
		Prompt:    &survey.Multiline{Message: "Type in the footer"},
		Validate:  validator(24),
		Transform: survey.TransformString(linter),
	},
}

// Prompt function assignes user input to Message struct
func Prompt() Message {
	// Perform the questions
	err := survey.Ask(qs, &msg)
	CheckIfError(err)
	return msg
}

func linter(s string) string {
	// Remove all leading and trailing white spaces
	s = strings.TrimSpace(s)
	// Then strings.Title the first word in string
	flds := strings.Fields(s)
	flds[0] = strings.Title(flds[0])
	return strings.Join(flds, " ")
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
