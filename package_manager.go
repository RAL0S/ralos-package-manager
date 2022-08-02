package main

import (
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

func (pm *PackageManager) installPackage(pkgToInstall string) {
	pkgToInstall = strings.ToLower(pkgToInstall)
	pkg, exists := pm.pkgIndex.Packages[pkgToInstall]
	if !exists {
		log.Printf("No package with name %s exists\n", pkgToInstall)
	}
	log.Printf("Installing %s, version %s\n", pkgToInstall, pkg.Version)
	pkgInstaller := PackageInstaller{}
	pkgInstaller.New(pm.cfg, pkg)
	status := pkgInstaller.Install()
	if status {
		log.Printf("Successfully installed %s\n", pkgToInstall)
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
