package cmd

import (
	"fmt"
	"log"
)

// CmdRevert is `direnv revert $0`
var CmdRevert = &Cmd{
	Name:    "revert",
	Desc:    "unloads an .envrc or .env and prints the diff in terms of exports",
	Args:    []string{"SHELL"},
	Private: true,
	Action:  cmdWithWarnTimeout(actionWithConfig(revertCommand)),
}

func revertCommand(currentEnv Env, args []string, config *Config) (err error) {
	defer log.SetPrefix(log.Prefix())
	log.SetPrefix(log.Prefix() + "revert:")
	logDebug("start")

	var target string

	if len(args) > 1 {
		target = args[1]
	}

	shell := DetectShell(target)
	if shell == nil {
		return fmt.Errorf("unknown target shell '%s'", target)
	}

	logDebug("loading RCs")
	loadedRC := config.LoadedRC()
	toLoad := findEnvUp(config.WorkDir, config.LoadDotenv)

	if loadedRC == nil && toLoad == "" {
		return
	}

	logDebug("updating RC")
	log.SetPrefix(log.Prefix() + "update:")

	logDebug("Determining action:")
	logDebug("toLoad: %#v", toLoad)
	logDebug("loadedRC: %#v", loadedRC)

	var previousEnv, newEnv Env

	if previousEnv, err = config.Revert(currentEnv); err != nil {
		err = fmt.Errorf("Revert() failed: %w", err)
		logDebug("err: %v", err)
		return
	}

	newEnv = previousEnv.Copy()
	newEnv.CleanContext()


	if out := diffStatus(previousEnv.Diff(newEnv)); out != "" {
		logStatus(currentEnv, "export %s", out)
	}

	diffString := currentEnv.Diff(newEnv).ToShell(shell)
	logDebug("env diff %s", diffString)
	fmt.Print(diffString)
	return
}
