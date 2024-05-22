# Red Alert Labs OS Package Manager

The source code of Red Alert Labs OS package manager.

## Compile

```
go build -o ralpm
```

## Usage

```
$ ralpm help
NAME:
   Red Alert Labs OS Package Manager (ralpm) - A package management for Red Alert Labs OS

USAGE:
   Red Alert Labs OS Package Manager (ralpm) [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   init, initialize, initialise  Initialize Red Alert Labs OS package manager
   list                          List installed packages
   install                       Install a package
   remove                        Removes an installed package
   help, h                       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## Prerequisites

The following packages must be installed before using ralpm.

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

Copy the ralpm binary to an empty folder, preferably within the home directory.

```
/home/kali/ralos/ralpm
```

Henceforth packages will be installed within the `/home/kali/ralos/` directory.

For the first time, inititalize the package manager by running

```
./ralpm init
```

This will create the config file `ralpm.toml` within the same directory.

For quicker access to the installed tools, add the bin directory (`/home/kali/ralos/bin`) to the system path. This can be done by adding the line `export PATH=$PATH:/home/kali/ralos/bin/` at the end of `.bashrc`.

## Package installation 

Packages can be installed by running
```
./ralpm install <package_name>
```

Example:
```
./ralpm install ghidra
./ralpm install cutter
```

## Package removal

To uninstall a package run
```
./ralpm remove <package_name>
```

This will prompt for confirmation before removing the package.

To uninstall without prompt run,
```
./ralpm remove --yes <package_name>
```

## List of available packages

The list of available packages can be found on the [package-index](https://github.com/RAL0S/package-index)