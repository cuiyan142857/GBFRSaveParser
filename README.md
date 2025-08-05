# GBFRSaveParser

This is a tool that imports **Granblue Fantasy: Relink** save files into **MySQL**.

**Usage**

You can run the packaged binary directly. Choose the binary for your platform (Linux / Windows / macOS).

```bash
Usage:
  ./GBFRSaveParser_for_mac import <SaveData1.dat> [db flags]
  ./GBFRSaveParser_for_mac export <out.csv> [db flags]

DB flags:
  -u, --user <user>         (default: root)
  -p, --password <pass>
  -H, --host <host:port>    (default: 127.0.0.1:3306)
  -d, --db, --database <db> (default: relink)

Example:
  ./GBFRSaveParser_for_mac import SaveData1.dat -u root -p mypassword -H 127.0.0.1:3306 -d relink
  ./GBFRSaveParser_for_mac export output.csv -u root -p mypassword -H 127.0.0.1:3306 -d relink
```

On success, you’ll see:

Import completed, 535771 rows written



Below is the development environment.

**Environment**

- Recommended: Go 1.23.10, MySQL 5.7.37

**Prerequisites**

- Go and MySQL installed

**Installation**

```bash
git clone https://github.com/cuiyan142857/GBFRSaveParser.git
cd GBFRSaveParser
go mod download

go run GBFRSaveParser.go import SaveData1.dat -u root -p mypassword -H 127.0.0.1:3306 -d relink

# On success, you’ll see:
# Import completed, 535771 rows written
```

**Acknowledgments:**

This project is licensed under the MIT License and includes code from [Nenkai](https://github.com/Nenkai/GBFRDataTools)’s repository.

Many thanks to nenkai for the project.