package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	if _, err := os.Stat("packages"); err == nil {
		return errors.New("cannot initialize, packages directory already exists")
	}

	cfg := GetConfig()
	cfg.InstallPath = currentDir

	err = os.Mkdir("packages", 0755)
	if err != nil {
		log.Println("Failed to create packages directory")
		return err
	}

	err = cfg.WriteToFile(CONFIG_FILE_NAME)
	if err != nil {
		return err
	}
	log.Println("Package manager initialized, created apm.toml")
	return nil
}

func listInstalledPackages() error {
	cfg := GetConfig()
	fmt.Println("[+] Listing installed packages")
	if len(cfg.Packages) == 0 {
		fmt.Println("[!] No packages are installed.")
	}
	for _, pkgInfo := range cfg.Packages {
		fmt.Printf("%s == %s\n", pkgInfo.Name, pkgInfo.Version)
	}
	return nil
}

func ensureInitialized() bool {
	exePath, err := os.Executable()
	if err != nil {
		log.Println("Failed to get executable path")
		return false
	}

	configFilePath := filepath.Join(filepath.Dir(exePath), CONFIG_FILE_NAME)
	if _, err := os.Stat(configFilePath); err != nil {
		fmt.Println("[!] Please initialize before running other commands!")
		return false
	}
	cfg := GetConfig()
	if err := cfg.LoadFromFile(configFilePath); err != nil {
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
			{
				Name:  "remove",
				Usage: "Removes an installed package",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "yes",
						Aliases: []string{"y"},
						Value:   false,
						Usage:   "Do not prompt to confirm before before removal",
					},
				},
				Action: func(ctx *cli.Context) error {
					if ensureInitialized() {
						if ctx.NArg() == 0 {
							fmt.Println("Please specify the package to remove!")
							return nil
						}
						targetPkg := ctx.Args().First()
						pm := PackageManager{}
						if pm.initialize(GetConfig()) {
							pm.removePackage(targetPkg, ctx.Bool("yes"))
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
