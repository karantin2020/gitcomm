package main

import (
	"io/ioutil"
	"log"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/karantin2020/gitcomm"
)

func main() {
	app := cli.App("gitcomm", "Automate git commit messaging")
	app.Version("V version", "gitcomm 0.1.1")

	app.Spec = "[-Avs]"

	var (
		// declare the -r flag as a boolean flag
		addFiles = app.BoolOpt("A addAll", false, "Adds, modifies, and removes index entries"+
			"to match the working tree. Evals `git add -A`")
		verbose = app.BoolOpt("v verbose", false, "Switch log output")
		show    = app.BoolOpt("s show", false, "Show last commit or not")
	)

	// Specify the action to execute when the app is invoked correctly
	app.Action = func() {
		if !*verbose {
			log.SetFlags(0)
			log.SetOutput(ioutil.Discard)
		}
		if !gitcomm.CheckForUncommited() {
			log.Printf("nothing to commit, working tree clean\n")
			return
		}
		log.Printf("there are new changes in working directory\n")
		msg := gitcomm.Prompt().String()
		gitcomm.GitExec(*addFiles, *show, msg)
	}

	// Invoke the app passing in os.Args
	app.Run(os.Args)
}
