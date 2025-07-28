# GBFRSaveParser

This is a tool that imports **Granblue Fantasy: Relink** save files into **MySQL**.

**Environment:**

- Recommended: Go 1.23.10, MySQL 5.7.37

**Prerequisites:**

- Go and MySQL installed

**Installation**

```
git clone https://github.com/cuiyan142857/GBFRSaveParser.git
cd GBFRSaveParser
go mod download

Edit line 42 of GBFRSaveParser.go to set your MySQL username, password, and database name (default: relink).

go run GBFRSaveParser.go import SaveData1.dat

On success, you’ll see:
Import completed, 535771 rows written
```

**Acknowledgments**

The data structures and parsing methods are derived from the project https://github.com/Nenkai/GBFRDataTools. Many thanks to nenkai for the project.