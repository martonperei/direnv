package cmd

import (
	"fmt"
	"os"
)

// CmdApplyExport is `direnv apply_export FILE`
var CmdApplyExport = &Cmd{
	Name:    "apply_export",
	Desc:    "Accepts a filename containing `direnv export` output and generates a series of bash export statements to apply the given env",
	Args:    []string{"SHELL", "TYPE", "FILE"},
	Private: true,
	Action:  actionSimple(cmdApplyExportAction),
}

func cmdApplyExportAction(env Env, args []string) (err error) {
	if len(args) < 4 {
		return fmt.Errorf("not enough arguments")
	}

	if len(args) > 4 {
		return fmt.Errorf("too many arguments")
	}	
	shell := DetectShell(args[1])
	if shell == nil {
		return fmt.Errorf("unknown target shell '%s'", args[1])
	}
	
	export_type := args[2];

	filename := args[3]

	dumped, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var dumpedEnv Env
	if (export_type == "json") {
		dumpedEnv, err = LoadEnvJSON(dumped);
		if err != nil {
			return err
		}
	} else {
		dumpedEnv, err = LoadEnv(string(dumped));
		if err != nil {
			return err
		}
	}

	exports := dumpedEnv.ToShell(shell)

	_, err = fmt.Println(exports)
	if err != nil {
		return err
	}

	return
}
