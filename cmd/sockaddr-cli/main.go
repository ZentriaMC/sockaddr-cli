package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/go-sockaddr/template"

	"github.com/ZentriaMC/sockaddr-cli/internal/core"
)

func main() {
	_ = core.Version
	runtime.GOMAXPROCS(1)

	args := os.Args
	if len(args) <= 1 {
		fmt.Fprintf(os.Stderr, "usage: %s [prog] [args...]\n", args[0])
		return
	}

	debug := os.Getenv("SOCKADDR_CLI_DEBUG") != ""
	dryRun := os.Getenv("SOCKADDR_CLI_DRY_RUN") != ""
	print0 := dryRun && os.Getenv("SOCKADDR_CLI_PRINT0") != ""

	if err := entrypoint(args[1:], debug, dryRun, print0); err != nil {
		fmt.Fprintf(os.Stderr, "unhandled error: %s\n", err)
		os.Exit(1)
	}
}

func entrypoint(args []string, debug, dryRun, print0 bool) (err error) {
	arg0 := args[0]

	var addrs sockaddr.IfAddrs
	if addrs, err = sockaddr.GetAllInterfaces(); err != nil {
		err = fmt.Errorf("failed to set up templating: %w", err)
		return
	}

	newArgs := make([]string, len(args))

	var processed string
	for idx, arg := range args {
		if processed, err = processArgument(addrs, arg); err != nil {
			err = fmt.Errorf("failed to process argument at idx %d: %w", idx, err)
			return
		}

		if debug {
			fmt.Fprintf(os.Stderr, "%s[%d] => %s\n", arg, idx, processed)
		}
		newArgs[idx] = processed
	}

	if dryRun {
		lastIdx := len(newArgs)-1
		for idx, value := range newArgs {
			fmt.Print(value)
			if idx != lastIdx {
				if print0 {
					fmt.Print("\u0000")
				} else {
					fmt.Print(" ")
				}
			} else if !print0 {
				fmt.Print("\n")
			}
		}
		return
	}

	var resolvedArg0 string
	if resolvedArg0, err = exec.LookPath(arg0); err != nil {
		err = fmt.Errorf("failed to look up '%s' in PATH: %w", arg0, err)
		return
	}

	if err = syscall.Exec(resolvedArg0, newArgs, os.Environ()); err != nil {
		err = fmt.Errorf("failed to execv: %w", err)
		return
	}

	return
}

func processArgument(addrs sockaddr.IfAddrs, argument string) (result string, err error) {
	result, err = template.ParseIfAddrs(argument, addrs)

	return
}
