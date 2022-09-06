# AttifyOS Package Manager

The source code of AttifyOS package manager.

## Compile

```
go build -o apm
```

## Usage

```
$ apm help
NAME:
   AttifyOS Package Manager (apm) - A package management for AttifyOS

USAGE:
   AttifyOS Package Manager (apm) [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   init, initialize, initialise  Initialize AttifyOS package manager
   list                          List installed packages
   install                       Install a package
   remove                        Removes an installed package
   help, h                       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## Prerequisites

The following packages must be installed before using apm.

```
fuse   # For running appimages
unzip
tar
bzip2
xz-utils
wget
snapd
git
```

You may skip installing the packages if they are already installed.
(Recent versions of Ubuntu already have snapd pre-installed)

## Setup

Copy the apm binary to an empty folder, preferably within the home directory.

```
/home/ubuntu/attifyos/apm
```

Henceforth packages will be installed within the `/home/ubuntu/attifyos/` directory.

For the first time, inititalize the package manager by running

```
./apm init
```

This will create the config file `apm.toml` within the same directory.

For quicker access to the installed tools, add the bin directory (`/home/ubuntu/attifyos/bin`) to the system path. This can be done by adding the line `export PATH=$PATH:/home/ubuntu/attifyos/bin/` at the end of `.bashrc`.

## Package installation 

Packages can be installed by running
```
./apm install <package_name>
```

Example:
```
./apm install ghidra
./apm install cutter
```

## Package removal

To uninstall a package run
```
./apm remove <package_name>
```

This will prompt for confirmation before removing the package.

To uninstall without prompt run,
```
./apm remove --yes <package_name>
```

## List of available packages

| Name              | Version                   | Comments                              | Repository URL                                | Build date   |
|-------------------|---------------------------|---------------------------------------|-----------------------------------------------|--------------|
| Apktool           | 2.6.1                     |                                       | https://github.com/AttifyOS/apktool           | Aug 8, 2022  |
| apk-mitm          | 1.2.1                     |                                       | https://github.com/AttifyOS/apk-mitm          | Sep 6, 2022  |
| Arduino           | 2.0.0-rc9.2               |                                       | https://github.com/AttifyOS/ArduinoIDE        | Aug 12, 2022 |
| attify-badge-tool | 0.0.1                     |                                       | https://github.com/AttifyOS/attify-badge-tool | Aug 16, 2022 |
| Baudrate          | 1.0                       |                                       | https://github.com/AttifyOS/baudrate          | Aug 17, 2022 |
| bdaddr            | 5.64                      |                                       | https://github.com/AttifyOS/bdaddr            | Aug 17, 2022 |
| Bettercap         | e224eea (June 13, 2022)   |                                       | https://github.com/AttifyOS/bettercap         | Aug 8, 2022  |
| binwalk           | 2.3.3                     | snap package                          | https://github.com/AttifyOS/binwalk           | Aug 23, 2022 |
| bluez             | latest                    | official snap package (auto-updating) | https://github.com/AttifyOS/bluez             | Aug 17, 2022 |
| cantoolz          | 3.7.0                     |                                       | https://github.com/AttifyOS/CANToolz          | Aug 29, 2022 |
| colordiff         | 1.0.20                    |                                       | https://github.com/AttifyOS/colordiff         | Aug 24, 2022 |
| cotopaxi          | 1.6.0                     |                                       | https://github.com/AttifyOS/cotopaxi          | Aug 29, 2022 |
| Cutter            | 2.1.0                     |                                       | https://github.com/AttifyOS/cutter            | Aug 5, 2022  |
| drozer            | a80c5f1 (Aug 11, 2022)    |                                       | https://github.com/AttifyOS/drozer            | Sep 05, 2022 |
| dump1090          | de61bd5 (Feb 4, 2020)     |                                       | https://github.com/AttifyOS/dump1090          | Aug 29, 2022 |
| dumpflash         | 0.0.1                     |                                       | https://github.com/AttifyOS/dumpflash         | Aug 29, 2022 |
| dwarf             | 1.1.3                     |                                       | https://github.com/AttifyOS/Dwarf             | Aug 29, 2022 |
| Flashrom          | 1.2                       | snap package                          | https://github.com/AttifyOS/flashrom          | Aug 12, 2022 |
| fmk               | 0.99--8403a17             | snap package                          | https://github.com/AttifyOS/firmware-mod-kit  | Aug 18, 2022 |
| Frida             | 15.2.2                    |                                       | https://github.com/AttifyOS/frida             | Aug 8, 2022  |
| fwanalyzer        | 1.4.3                     |                                       | https://github.com/AttifyOS/fwanalyzer        | Sep 2, 2022  |
| gdb-multiarch     | 12.1                      |                                       | https://github.com/AttifyOS/gdb-multiarch     | Sep 4, 2022  |
| Ghidra            | 10.1.5                    |                                       | https://github.com/AttifyOS/ghidra            | Aug 5, 2022  |
| gnuradio          | 3.10.3.0                  |                                       | https://github.com/AttifyOS/gnuradio          | Aug 18, 2022 |
| Gqrx              | v2.15.9                   |                                       | https://github.com/AttifyOS/gqrx              | Aug 16, 2022 |
| inspectrum        | 98b998ff (March 27, 2022) |                                       | https://github.com/AttifyOS/inspectrum        | Aug 25, 2022 |
| iotsecfuzz        | 1.0.0                     |                                       | https://github.com/AttifyOS/iotsecfuzz        | Sep 5, 2022  |
| JADX              | 1.4.3                     |                                       | https://github.com/AttifyOS/jadx              | Aug 5, 2022  |
| kalibrate_rtl     | 0.4.1                     |                                       | https://github.com/AttifyOS/kalibrate-rtl     | Aug 30, 2022 |
| killerbee         | 3.0.0-beta.2              |                                       | https://github.com/AttifyOS/killerbee         | Aug 13, 2022 |
| linuxrouter       | 0.6.6                     |                                       | https://github.com/AttifyOS/linux-router      | Aug 19, 2022 |
| LTE_Cell_Scanner  | bef6ef4 (April 25, 2022)  |                                       | https://github.com/AttifyOS/LTE-Cell-Scanner  | Aug 31, 2022 |
| Nmap              | 7.92                      | snap package                          | https://github.com/AttifyOS/nmap              | Aug 11, 2022 |
| Objection         | 1.11.0                    | snap package                          | https://github.com/AttifyOS/objection         | Aug 10, 2022 |
| ook-decoder       | cc108f9 (July 19, 2021)   |                                       | https://github.com/AttifyOS/ook-decoder       | Aug 26, 2022 |
| OpenOCD           | 0.11.0-4                  |                                       | https://github.com/AttifyOS/OpenOCD           | Aug 26, 2022 |
| radare2           | 5.7.6                     |                                       | https://github.com/AttifyOS/radare2           | Aug 23, 2022 |
| routersploit      | 3.4.0                     |                                       | https://github.com/AttifyOS/routersploit      | Aug 29, 2022 |
| rtl_433           | 21.12-146                 |                                       | https://github.com/AttifyOS/rtl_433           | Sep 5, 2022  |
| spectrum_painter  | 0.1                       |                                       | https://github.com/AttifyOS/spectrum_painter  | Aug 31, 2022 |
| SRecord           | 1.64                      |                                       | https://github.com/AttifyOS/SRecord           | Aug 13, 2022 |
| urh               | 2.9.3                     |                                       | https://github.com/AttifyOS/urh               | Aug 29, 2022 |