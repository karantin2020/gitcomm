package gitcomm

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	bb "github.com/karantin2020/promptui"
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

func fillMessage(msg *Message) {
	var err error
	msg.Type, err = bb.PromptAfterSelect("Choose a type(<scope>)", types)
	checkInterrupt(err)
	p := bb.Prompt{
		BasicPrompt: bb.BasicPrompt{
			Label:     "Type in the subject",
			Formatter: linterSubject,
			Validate: func(s string) error {
				if s == "" {
					return bb.NewValidationError("Subject must not be empty string")
				}
				if len(s) > 72 {
					return bb.NewValidationError("Subject cannot be longer than 72 characters")
				}
				return nil
			},
		},
	}
	msg.Subject, err = p.Run()
	checkInterrupt(err)
	mlBody := bb.MultilinePrompt{
		BasicPrompt: bb.BasicPrompt{
			Label:     "Type in the body",
			Default:   msg.Subject,
			Formatter: linterBody,
			Validate: func(s string) error {
				if s == "" {
					return bb.NewValidationError("Body must not be empty string")
				}
				if len(s) > 320 {
					return bb.NewValidationError("Body cannot be longer than 320 characters")
				}
				return nil
			},
		},
	}
	msg.Body, err = mlBody.Run()
	checkInterrupt(err)
	mlFoot := bb.MultilinePrompt{
		BasicPrompt: bb.BasicPrompt{
			Label:     "Type in the foot",
			Formatter: linterFoot,
			Validate: func(s string) error {
				if s == "" {
					return bb.NewValidationError("Foot must not be empty string")
				}
				if len(s) > 50 {
					return bb.NewValidationError("Foot cannot be longer than 50 characters")
				}
				return nil
			},
		},
	}
	msg.Foot, err = mlFoot.Run()
	checkInterrupt(err)
}

// Prompt function assignes user input to Message struct
func Prompt() Message {
	fillMessage(&msg)
	Info("\nCommit message is:\n%s\n", msg.String())
	c, err := bb.Confirm("Is everything OK? Continue", "N", true)
	checkInterrupt(err)
	if c == "N" {
		log.Printf("Commit flow interrupted by user\n")
		os.Exit(1)
	}
	return msg
}

// TagPrompt prompting tag version level to upgrade
func TagPrompt() string {
	s := bb.Select{
		Label: "Choose tag level",
		Items: []string{
			"patch",
			"minor",
			"major",
		},
		Default: 0,
	}
	_, level, err := s.Run()
	checkInterrupt(err)
	return level
}

// PromptConfirm is a common function to ask confirm before some action
func PromptConfirm(msg string) bool {
	c, err := bb.Confirm(msg, "N", false)
	checkInterrupt(err)
	if c == "N" {
		return false
	}
	return true
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
	var upl = func(sl string) string {
		rs := []rune(sl)
		if len(rs) > 0 {
			rs[0] = unicode.ToUpper(rs[0])
		}
		return string(rs)
	}
	out := []string{}
	ins := strings.Split(s, "\n")
	for i := range ins {
		out = append(out, upl(ins[i]))
	}
	return strings.Join(out, "\n")
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
		if err != bb.ErrInterrupt {
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
		}
		log.Printf("interrupted by user\n")
		os.Exit(1)
	}
}
