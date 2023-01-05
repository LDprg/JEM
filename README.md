# JEM
[![Publish Dev Release](https://github.com/LDprg/JEM/actions/workflows/build.yml/badge.svg)](https://github.com/LDprg/JEM/actions/workflows/build.yml)

JEM is a Java Environment Manager written in go 

## Installation
Download the latest release from the [releases page](https://github.com/LDprg/JEM/releases) and extract it to a folder of your choice.

> Uninstall all other Java versions before installing JEM.

You will need to restart your terminal after installation in order to update your PATH.

### Linux & Mac
Run:
```bash
./jem_install
```
You may need to reload the .bashrc file after installation:
```bash
source ~/.bashrc
```

### Windows
Run:
```bash
jem_install.exe
```

## Usage
```bash
jem [command] [options]
```

### Commands
#### list
```bash
jem list
```
Lists the current, installed and installable versions of java

#### install
```bash
jem install [version]
```
Installs the specified version of java

#### use
```bash
jem use [version]
```
Sets the current version of java

#### uninstall
```bash
jem uninstall [version]
```
Uninstalls the specified version of java

## Uninstallation

> Notice: This will remove all installed versions of java

### Linux & Mac
Run:
```bash
./jem_uninstall
```

### Windows
Run:
```bash
jem_uninstall.exe
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
This projets use the [MIT](https://choosealicense.com/licenses/mit/) License.

## Where does the java versions come from?
The java versions are downloaded with [Adoptium Api](https://adoptium.net/).