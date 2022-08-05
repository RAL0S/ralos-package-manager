package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/pelletier/go-toml/v2"
)

type PackageIndex struct {
	LastUpdated time.Time              `toml:"last_updated"`
	Packages    map[string]PackageInfo `toml:"packages"`
}

type PackageManager struct {
	cfg      *APMConfig
	pkgIndex PackageIndex
}

func (pm *PackageManager) initialize(cfg *APMConfig) bool {
	pm.cfg = cfg

	log.Println("Fetching package index")
	pkgIndexBytes, err := fetchPackageIndex()
	if err != nil {
		log.Println("Failed to fetch package index")
		return false
	}

	if err := toml.Unmarshal(pkgIndexBytes, &pm.pkgIndex); err != nil {
		log.Println("Failed to parse package index")
		return false
	}
	return true
}

func (pm *PackageManager) isInstalled(localName string, isTesting bool) bool {	
	if isTesting {return false}
	cfg := GetConfig()
	for pkgName, pkgInfo := range cfg.Packages {
		if pkgName == localName && pkgInfo.Testing == isTesting {
			return true
		}
	}
	return false
}

func (pm *PackageManager) getInstalledVersion(localName string) string {
	cfg := GetConfig()
	return cfg.Packages[localName].Version
}

func (pm *PackageManager) removePackage(pkgToRemove string, noPrompt bool) {
	pkgToRemove = strings.ToLower(pkgToRemove)
	cfg := GetConfig()
	pkg, exists := cfg.Packages[pkgToRemove]
	if !exists {
		log.Printf("No such package '%s' is installed\n", pkgToRemove)
		return
	}
	if !noPrompt {
		if pkg.Testing {
			fmt.Printf("This will remove %s (testing)\n", pkg.Name)
		} else {
			fmt.Printf("This will remove %s version %s\n", pkg.Name, pkg.Version)
		}
		fmt.Print("Are you sure (y/n)? ")

		var choice string
		fmt.Scanln(&choice)
		if strings.ToLower(choice) != "y" {
			fmt.Println("Aborting.")
			return
		}
	}
	log.Printf("Removing package %s version %s\n", pkgToRemove, pkg.Version)

	pkgUninstaller := PackageUninstaller{}
	pkgUninstaller.New(pm.cfg, pkg)
	status := pkgUninstaller.Uninstall()
	if status {
		delete(cfg.Packages, pkgToRemove)
		cfg.Save()
		log.Println("Uninstallation successful")
	} else {
		log.Println("Uninstallation failed")
	}
}

func (pm *PackageManager) installPackage(pkgToInstall string, installTesting bool) {
	pkgToInstall = strings.ToLower(pkgToInstall)
	pkg, exists := pm.pkgIndex.Packages[pkgToInstall]
	if !exists {
		log.Printf("No package with name %s exists\n", pkgToInstall)
		return
	}
	isUpgrade := false
	if pm.isInstalled(pkgToInstall, false) {
		// Stable version installed
	
		// Don't install testing version if stable already installed
		if installTesting {
			log.Println("Can't install testing version as stable version is already installed")
			return
		}

		installedVersion := pm.getInstalledVersion(pkgToInstall)
		if installedVersion == pkg.Version {
			log.Printf("Package %s, version %s is already installed\n", pkgToInstall, pkg.Version)
			return
		}
		log.Printf("Found updated package, installed=%s, available=%s, upgrading\n", installedVersion, pkg.Version)
		isUpgrade = true
	}

	if installTesting {
		log.Printf("Installing %s (*testing)\n", pkgToInstall)
	} else {
		log.Printf("Installing %s, version %s\n", pkgToInstall, pkg.Version)
	}
	pkgInstaller := PackageInstaller{}
	pkgInstaller.New(pm.cfg, pkg, installTesting)
	status := pkgInstaller.Install()
	if status {
		//No previously installed package
		if len(pm.cfg.Packages) == 0 {
			pm.cfg.Packages = make(map[string]PackageInfo)
		}

		if isUpgrade {
			log.Println("Removing old package version")
			pm.removePackage(pkgToInstall, true)
		}

		log.Println("Updating local package index")
		if installTesting {	pkg.Testing = true }
		pm.cfg.Packages[pkgToInstall] = pkg
		
		if err := pm.cfg.Save(); err != nil {
			log.Println("Successfully installed but failed to update  the local package index")
		} else {
			log.Printf("Successfully installed %s, version %s\n", pkgToInstall, pkg.Version)
		}
		return
	}
	log.Println("Installation failed")
}

func fetchPackageIndex() ([]byte, error) {
	req, err := grab.NewRequest("", PACKAGE_INDEX_URL)
	if err != nil {
		return nil, err
	}
	req.NoStore = true

	client := grab.NewClient()
	resp := client.Do(req)
	data, err := resp.Bytes()
	if err != nil {
		return nil, err
	}
	return data, nil
}
