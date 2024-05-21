package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type PackageUninstaller struct {
	cfg     *RALPMConfig
	pkgInfo PackageInfo
}

func (pu *PackageUninstaller) New(cfg *RALPMConfig, pkgInfo PackageInfo) {
	pu.cfg = cfg
	pu.pkgInfo = pkgInfo
}

func (pu *PackageUninstaller) Uninstall() bool {
	pkgCloneDir := filepath.Join(pu.cfg.InstallPath, CLONE_DIR_NAME, pu.pkgInfo.Name+"-"+pu.pkgInfo.Version)
	pkgTomlPath := filepath.Join(pkgCloneDir, "package.toml")
	pkgTomlBytes, err := os.ReadFile(pkgTomlPath)
	if err != nil {
		log.Println("Failed to read package config file", pkgTomlPath)
		return false
	}

	var pkgDescriptor PackageDescriptor
	if err := toml.Unmarshal(pkgTomlBytes, &pkgDescriptor); err != nil {
		log.Println("Failed to parse config file", pkgTomlPath)
		return false
	}

	installScriptPath := filepath.Join(pkgCloneDir, pkgDescriptor.PkgDetail.InstallScript)
	if _, err := os.Stat(installScriptPath); err != nil {
		log.Println("Package installer script not found:", installScriptPath)
		return false
	}

	// Make install script executable
	os.Chmod(installScriptPath, 0744)

	pkgPath := filepath.Join(pu.cfg.InstallPath, PKG_DIR_NAME, pu.pkgInfo.Name+"-"+pu.pkgInfo.Version)
	tmpDir, _ := os.MkdirTemp("", "ralpm_tmp")
	defer os.RemoveAll(tmpDir)

	cmd := exec.Command(installScriptPath, "uninstall")
	cmd.Env = append(
		os.Environ(),
		"RALPM_TMP_DIR="+tmpDir,
		"RALPM_PKG_INSTALL_DIR="+pkgPath,
		"RALPM_PKG_BIN_DIR="+filepath.Join(pu.cfg.InstallPath, BIN_DIR_NAME),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("Running package uninstallation script, please wait...")
	err = cmd.Run()
	if err != nil {
		log.Println("Failed to execute uninstall script")
		return false
	}

	if cmd.ProcessState.ExitCode() != 0 {
		log.Println("Uninstall script returned non zero status code")
		return false
	}

	os.RemoveAll(pkgPath)
	os.RemoveAll(pkgCloneDir)
	return true
}
