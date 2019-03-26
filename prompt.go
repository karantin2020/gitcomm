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
		"feat	[new feature]",
		"fix		[bug fix]",
		"docs	[changes to documentation]",
		"style	[format, missing semi colons, etc; no code change]",
		"refactor	[refactor production code]",
		"test	[add missing tests, refactor tests; no production code change]",
		"chore	[update grunt tasks etc; no production code change]",
		"version	[description of version upgrade]",
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
			Default:   "# If applied, this commit will\n",
			Formatter: linterBody,
			Validate: func(s string) error {
				if s == "" {
					return bb.NewValidationError("Body must not be empty string")
				}
				if len(s) > 320 {
					return bb.NewValidationError("Body cannot be longer than 320 characters")
				}
				ins := strings.Split(s, "\n")
				for i := range ins {
					if len(ins[i]) > 72 {
						return bb.NewValidationError("Body must be wraped at 72 characters")
					}
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
func Prompt() string {
	fillMessage(&msg)
	gitMsg := msg.String()
	Info("\nCommit message is:\n%s\n", gitMsg)
	var err error
	cp := bb.ConfirmPrompt{
		BasicPrompt: bb.BasicPrompt{
			Label:   "Is everything OK? Continue",
			Default: "N",
			NoIcons: true,
		},
		ConfirmOpt: "e",
	}
	c, err := cp.Run()
	checkConfirmStatus(c, err)
	if c == "E" {
		gitMsg, err = bb.Editor("", gitMsg)
		checkInterrupt(err)
		numlines := len(strings.Split(gitMsg, "\n")) + 2
		for ; numlines > 0; numlines-- {
			fmt.Print(bb.ClearUpLine())
		}
		Info("Commit message is:\n%s", gitMsg)
		checkConfirmStatus(bb.Confirm("Is everything OK? Continue", "N", true))
		return gitMsg
	}
	return gitMsg
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
		// if the line is commented with # at the start pass that line
		if ins[i][0] == '#' {
			continue
		}
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

func checkConfirmStatus(c string, err error) {
	checkInterrupt(err)
	if c == "N" {
		log.Printf("Commit flow interrupted by user\n")
		os.Exit(1)
	}
}
