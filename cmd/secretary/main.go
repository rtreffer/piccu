package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"syscall"

	"github.com/rtreffer/piccu/pkg/secretary"
)

func main() {
	flags, plainVar, passVar, setVar, unsetVar := secretary.NewMultiFlagset()
	flag.Var(plainVar, "plain", "plain file secret to load")
	flag.Var(passVar, "pass", "pass secret name to load")
	flag.Var(setVar, "set", "set environment variables")
	flag.Var(unsetVar, "unset", "unset environmenr variables")
	passConfig := flag.String("pass.config", "", "pass config file to load")
	passStoreDir := flag.String("pass.store.dir", "", "pass store directory to use")

	showHelp := flag.Bool("help", false, "displays a help text")
	flag.BoolVar(showHelp, "h", false, "displays a help text")

	flag.Parse()

	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	if len(*flags) == 0 {
		fmt.Fprintln(os.Stderr, "no secrets backend specified - exiting")
		os.Exit(1)
	}

	keys := make([]string, 0, 32)

	if *passConfig != "" {
		secretary.PassSetConfig(*passConfig)
	}
	if *passStoreDir != "" {
		secretary.PassSetStoreDir(*passStoreDir)
	}

	for _, flagValue := range *flags {
		if flagValue.Type == secretary.LoadPlainFile {
			env, err := secretary.PlainLoadSecret(flagValue.Name)
			if err != nil {
				panic(err)
			}
			for k, v := range env {
				if err := os.Setenv(k, v); err != nil {
					panic(err)
				}
				keys = append(keys, k)
			}
			continue
		}

		if flagValue.Type == secretary.LoadPass {
			env, err := secretary.PassLoadSecret(flagValue.Name)
			if err != nil {
				panic(err)
			}
			if len(env) == 0 {
				fmt.Fprintln(os.Stderr, "no secrets loaded from", flagValue.Name)
				os.Exit(1)
			}
			for k, v := range env {
				if err := os.Setenv(k, v); err != nil {
					panic(err)
				}
				keys = append(keys, k)
			}
		}

		if flagValue.Type == secretary.Unset {
			os.Unsetenv(flagValue.Name)
		}

		if flagValue.Type == secretary.Set {
			parts := strings.SplitN(flagValue.Name, "=", 2)
			if len(parts) == 1 {
				os.Setenv(parts[0], "")
				continue
			}
			keys = append(keys, parts[0])
			os.Setenv(parts[0], parts[1])
		}
	}

	if oldKeys, _ := os.LookupEnv("SECRETARY_KEYS"); oldKeys != "" {
		keys = append(keys, strings.Split(oldKeys, " ")...)
	}

	// it is possible that keys exist multiple times or were deleted from the environment
	currentKeyMap := make(map[string]bool)
	for _, key := range keys {
		_, currentKeyMap[key] = os.LookupEnv(key)
	}
	keys = keys[:0]
	for k, v := range currentKeyMap {
		if v {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)
	if err := os.Setenv("SECRETARY_KEYS", strings.Join(keys, " ")); err != nil {
		panic(err)
	}

	remains := flag.Args()
	argv0, err := exec.LookPath(remains[0])
	if err != nil {
		argv0 = remains[0]
	}

	// if Exec succeeds then the new program takes over and panic is never called
	panic(syscall.Exec(argv0, remains, os.Environ()))
}
