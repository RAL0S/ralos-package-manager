package main

import (
	"log"
	"os"
	"path/filepath"
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

func (pm *PackageManager) isInstalled(localName string) bool {
	cfg := GetConfig()
	for pkgName := range cfg.Packages {
		if pkgName == localName {
			return true
		}
	}
	return false
}

func (pm *PackageManager) getInstalledVersion(localName string) string {
	cfg := GetConfig()
	return cfg.Packages[localName].Version
}

func (pm *PackageManager) removePackage(pkgToRemove string) {
	cfg := GetConfig()
	pkgDetails := cfg.Packages[pkgToRemove]
	pkgPath := filepath.Join(cfg.InstallPath, "packages", pkgDetails.Name+"-"+pkgDetails.Version)
	os.RemoveAll(pkgPath)
}

func (pm *PackageManager) installPackage(pkgToInstall string) {
	pkgToInstall = strings.ToLower(pkgToInstall)
	pkg, exists := pm.pkgIndex.Packages[pkgToInstall]
	if !exists {
		log.Printf("No package with name %s exists\n", pkgToInstall)
		return
	}
	isUpgrade := false
	if pm.isInstalled(pkgToInstall) {
		installedVersion := pm.getInstalledVersion(pkgToInstall)
		if installedVersion == pkg.Version {
			log.Printf("Package %s, version %s is already installed\n", pkgToInstall, pkg.Version)
			return
		}
		log.Printf("Found updated package, installed=%s, available=%s, upgrading\n", installedVersion, pkg.Version)
		isUpgrade = true		
	}

	log.Printf("Installing %s, version %s\n", pkgToInstall, pkg.Version)
	pkgInstaller := PackageInstaller{}
	pkgInstaller.New(pm.cfg, pkg)
	status := pkgInstaller.Install()
	if status {
		//No previously installed package
		if len(pm.cfg.Packages) == 0 {
			pm.cfg.Packages = make(map[string]PackageInfo)
		}

		if isUpgrade {
			log.Println("Removing old package version")
			pm.removePackage(pkgToInstall)
		}

		log.Println("Updating local package index")
		pm.cfg.Packages[pkgToInstall] = pkg
		if err := pm.cfg.WriteToFile(filepath.Join(pm.cfg.InstallPath, CONFIG_FILE_NAME)); err != nil {
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
