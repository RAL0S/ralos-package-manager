package main

import (
	"log"
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

func (pm *PackageManager) installPackage(pkgToInstall string) {
	pkgToInstall = strings.ToLower(pkgToInstall)
	pkg, exists := pm.pkgIndex.Packages[pkgToInstall]
	if !exists {
		log.Printf("No package with name %s exists\n", pkgToInstall)
		return
	}
	log.Printf("Installing %s, version %s\n", pkgToInstall, pkg.Version)
	pkgInstaller := PackageInstaller{}
	pkgInstaller.New(pm.cfg, pkg)
	status := pkgInstaller.Install()
	if status {
		log.Println("Updating local package index")
		
		//No previously installed package
		if len(pm.cfg.Packages) == 0 {
			pm.cfg.Packages = make(map[string]PackageInfo)
		}
		pm.cfg.Packages[pkgToInstall] = pkg
		if err := pm.cfg.WriteToFile(filepath.Join(pm.cfg.InstallPath, CONFIG_FILE_NAME)); err != nil {
			log.Println("Successfully installed but failed to update  the local package index")
		} else {
			log.Printf("Successfully installed %s\n", pkgToInstall)			
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
