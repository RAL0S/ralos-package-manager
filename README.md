# AttifyOS Package Manager

The source code of AttifyOS package manager.

## Compile

```
go build -o apm
```

## Usage

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

From now on, you can install packages by running
```
./apm install <package_name>
```

```
./apm install ghidra
./apm install cutter
```
