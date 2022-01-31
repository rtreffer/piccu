package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rtreffer/piccu/pkg/cicci"
	"github.com/rtreffer/piccu/pkg/flags"
	"github.com/rtreffer/piccu/pkg/ioutils"
	"github.com/rtreffer/piccu/pkg/piccu"
	"github.com/rtreffer/piccu/pkg/secretary"
)

func main() {
	showHelp := flag.Bool("help", false, "displays a help text")
	flag.BoolVar(showHelp, "h", false, "displays a help text")

	fileFlags, plainVar, passVar, setVar, unsetVar := secretary.NewMultiFlagset()
	flag.Var(plainVar, "plain", "plain environment file to load")
	flag.Var(passVar, "pass", "pass secret name to load")
	flag.Var(setVar, "set", "set environment variables")
	flag.Var(unsetVar, "unset", "unset environmenr variables")
	passConfig := flag.String("pass.config", "", "pass config file to load")
	passStoreDir := flag.String("pass.store.dir", "", "pass store directory to use")

	release := flag.String("ubuntu", "focal:arm64", "ubuntu release to use (supported releases: "+strings.Join(piccu.GetImageNames(), ",")+")")
	output := flag.String("output", "disk.img", "output image")

	injectBootFile := make(flags.StringArray, 0)
	flag.Var(&injectBootFile, "boot.firmware.file", "inject the give file under /boot/firmware (e.g. meta-data)")

	flag.Parse()

	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	// resolve the image and load it
	image, found := piccu.ImagesByKey()[*release]
	if !found {
		fmt.Fprintln(os.Stderr, "could not find image", *release)
		os.Exit(1)
	}
	cached, err := piccu.Fetch(image, "", 7*24*time.Hour)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not download", *release, err)
		os.Exit(2)
	}
	_ = cached

	// load secrets

	if *passConfig != "" {
		secretary.PassSetConfig(*passConfig)
	}
	if *passStoreDir != "" {
		secretary.PassSetStoreDir(*passStoreDir)
	}

	secretKeys := make(map[string]string)
	for _, flagValue := range *fileFlags {
		if flagValue.Type == secretary.LoadPlainFile {
			env, err := secretary.PlainLoadSecret(flagValue.Name)
			if err != nil {
				panic(err)
			}
			for k, v := range env {
				secretKeys[k] = v
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
				secretKeys[k] = v
			}
		}

		if flagValue.Type == secretary.Unset {
			delete(secretKeys, flagValue.Name)
		}

		if flagValue.Type == secretary.Set {
			parts := strings.SplitN(flagValue.Name, "=", 2)
			if len(parts) == 1 {
				secretKeys[parts[0]] = ""
				continue
			}
			secretKeys[parts[0]] = parts[1]
		}
	}

	// build the cloud-config
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}
	files, err := cicci.CollectFiles(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't find files to merge:", err)
		os.Exit(1)
	}
	cloudConfigArchive := ""
	if len(files) != 0 {
		expanded, err := files.LoadAndExpand(secretKeys)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't load/expand files:", err)
			os.Exit(3)
		}
		errors := expanded.Validate()
		for _, err := range errors {
			fmt.Fprintln(os.Stderr, "WARNING:", err)
		}
		cloudConfigArchive, err = cicci.CreateMultipartArchive(expanded)
		if err != nil {
			fmt.Fprintln(os.Stderr, "can't create multipart archive:", err)
			os.Exit(4)
		}
	}

	// copy the file to the output
	if stat, err := os.Stat(*output); err == nil {
		if stat.Mode().IsRegular() {
			os.Remove(*output)
		}
	}

	err = ioutils.Copy(cached, *output)
	if err != nil {
		os.Remove(*output)
		panic(err)
	}

	// inject the cloud-config
	fmt.Println("modifying", *output)

	img, err := piccu.OpenImage(*output)
	if err != nil {
		os.Remove(*output)
		panic(err)
	}
	if cloudConfigArchive != "" {
		userData, err := piccu.GzipString(cloudConfigArchive)
		if err != nil {
			os.Remove(*output)
			panic(err)
		}
		fmt.Println("adding user-data")
		if err := img.InjectFile("user-data", userData); err != nil {
			os.Remove(*output)
			panic(err)
		}
	}

	for _, bootfile := range injectBootFile {
		data, err := os.ReadFile(bootfile)
		if err != nil {
			os.Remove(*output)
			panic(err)
		}
		name := filepath.Base(bootfile)
		fmt.Println("adding", name)
		if err := img.InjectFile(name, data); err != nil {
			os.Remove(*output)
			panic(err)
		}
	}

	fmt.Println("syncing", *output)
	img.Close()
}
