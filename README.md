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

| Name      | Version                 |
|-----------|-------------------------|
| Apktool   | 2.6.1                   |
| Bettercap | e224eea (June 13, 2022) |
| Cutter    | 2.1.0                   |
| Frida     | 15.2.2                  |
| Ghidra    | 10.1.5                  |
| JADX      | 1.4.3                   |