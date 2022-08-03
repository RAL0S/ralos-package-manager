package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func initAPM() error {
	if _, err := os.Stat(CONFIG_FILE_NAME); err == nil {
		return errors.New("cannot initialize, apm.toml already exists")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Failed to fetch working directory")
		return err
	}
	cfg := GetConfig()
	cfg.InstallPath = currentDir

	/*Packages: map[string]PackageInfo{
		"cutter": {
			Name: "cutter",
			Version: "2.1.0",
		},
		"ghidra": {
			Name: "ghidra",
			Version: "10.1.5",
		},
	},*/
	err = os.Mkdir("packages", 0755)
	if err != nil {
		log.Println("Failed to create packages directory")
		return err
	}

	return cfg.WriteToFile(CONFIG_FILE_NAME)
}

func listInstalledPackages() error {
	cfg := GetConfig()
	if err := cfg.LoadFromFile("apm.toml"); err != nil {
		return err
	}
	if len(cfg.Packages) == 0 {
		fmt.Println("[!] No packages are installed.")
	}
	fmt.Println("[+] Listing installed packages")
	for _, pkgInfo := range cfg.Packages {
		fmt.Printf("%s == %s\n", pkgInfo.Name, pkgInfo.Version)
	}
	return nil
}

func ensureInitialized() bool {
	if _, err := os.Stat(CONFIG_FILE_NAME); err != nil {
		fmt.Println("[!] Please initialize before running other commands!")
		return false
	}
	cfg := GetConfig()
	if err := cfg.LoadFromFile("apm.toml"); err != nil {
		fmt.Println("[!] Please initialize before running other commands!")
		return false
	}
	return true
}

func main() {
	app := &cli.App{
		Name:    "AttifyOS package manager (apm)",
		Version: APP_VERSION,
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"initialize", "initialise"},
				Usage:   "Initialize AttifyOS package manager",
				Action: func(ctx *cli.Context) error {
					return initAPM()
				},
			},
			{
				Name:  "list",
				Usage: "List installed packages",
				Action: func(ctx *cli.Context) error {
					if ensureInitialized() {
						return listInstalledPackages()
					}
					return nil
				},
			},
			{
				Name:  "install",
				Usage: "Install a package",
				Action: func(ctx *cli.Context) error {
					if ensureInitialized() {
						if ctx.NArg() == 0 {
							fmt.Println("Please specify the package to install!")
							return nil
						}
						targetPkg := ctx.Args().First()
						pm := PackageManager{}
						if pm.initialize(GetConfig()) {
							pm.installPackage(targetPkg)
						}
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
