//  GBFRSaveParser.go
//
//  go run GBFRSaveParser.go import SaveData1.dat
//  go run GBFRSaveParser.go export dump.csv
//
// Dependencies：
//   go get github.com/google/flatbuffers/go
//   go get github.com/go-sql-driver/mysql
//
// save_units
//   id        BIGINT AUTO_INCREMENT PRIMARY KEY
//   idtype    INT
//   unitid    INT
//   val_index INT
//   value     BIGINT
//   table_ty  VARCHAR(8)

package main

import (
	"database/sql"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"

	fb "GBFRDataTools/FlatBuffers"

	_ "github.com/go-sql-driver/mysql"
)

type header struct {
	MainVer, Reserved, SubVer         uint32
	SteamID                           uint64
	Offset1, SlotOff, Size1, SlotSize uint64
}

const dsn = "user:password@tcp(127.0.0.1:3306)/relink?charset=utf8mb4&parseTime=true"

func ensureTable(db *sql.DB) {
	sqlStmt := `CREATE TABLE IF NOT EXISTS save_units(
		id        BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
		idtype    INT UNSIGNED,
		unitid    INT UNSIGNED,
		val_index INT UNSIGNED,
		value     BIGINT,
		table_ty  VARCHAR(8),
		PRIMARY KEY(id),
		UNIQUE KEY uk(idtype,unitid,val_index,table_ty)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	if _, err := db.Exec(sqlStmt); err != nil {
		log.Fatalf("create table: %v", err)
	}
}

func importBin(db *sql.DB, path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var h header
	if err := binary.Read(f, binary.LittleEndian, &h); err != nil {
		log.Fatal(err)
	}

	readSeg := func(off, size uint64) []byte {
		buf := make([]byte, size)
		if _, err := f.ReadAt(buf, int64(off)); err != nil && err != io.EOF {
			log.Fatalf("read segment: %v", err)
		}
		return buf
	}

	seg1 := readSeg(h.Offset1, h.Size1)
	slot := readSeg(h.SlotOff, h.SlotSize)

	process := func(buf []byte, stmt *sql.Stmt) int {
		root := fb.GetRootAsSaveDataBinary(buf, 0)
		rows := 0

		/* --- 1. BoolTable --- */
		bu := new(fb.BoolSaveDataUnit)
		for i := 0; i < root.BoolTableLength(); i++ {
			if !root.BoolTable(bu, i) {
				continue
			}
			for j := 0; j < bu.ValueDataLength(); j++ {
				val := 0
				if bu.ValueData(j) {
					val = 1
				}
				stmt.Exec(bu.Idtype(), bu.UnitId(), uint32(j), uint64(val), "Bool")
				rows++
			}
		}

		/* --- 2. ByteTable --- */
		bb := new(fb.ByteSaveDataUnit)
		for i := 0; i < root.ByteTableLength(); i++ {
			if !root.ByteTable(bb, i) {
				continue
			}
			for j := 0; j < bb.ValueDataLength(); j++ {
				v := int8(bb.ValueData(j))
				stmt.Exec(bb.Idtype(), bb.UnitId(), uint32(j), uint64(int64(v)), "Byte")
				rows++
			}
		}

		/* --- 3. UByteTable --- */
		ub := new(fb.UByteSaveDataUnit)
		for i := 0; i < root.UbyteTableLength(); i++ {
			if !root.UbyteTable(ub, i) {
				continue
			}
			for j := 0; j < ub.ValueDataLength(); j++ {
				stmt.Exec(ub.Idtype(), ub.UnitId(), uint32(j), uint64(ub.ValueData(j)), "UByte")
				rows++
			}
		}

		/* --- 4. ShortTable --- */
		sh := new(fb.ShortSaveDataUnit)
		for i := 0; i < root.ShortTableLength(); i++ {
			if !root.ShortTable(sh, i) {
				continue
			}
			for j := 0; j < sh.ValueDataLength(); j++ {
				v := int16(sh.ValueData(j))
				stmt.Exec(sh.Idtype(), sh.UnitId(), uint32(j), uint64(int64(v)), "Short")
				rows++
			}
		}

		/* --- 5. UShortTable --- */
		us := new(fb.UShortSaveDataUnit)
		for i := 0; i < root.UshortTableLength(); i++ {
			if !root.UshortTable(us, i) {
				continue
			}
			for j := 0; j < us.ValueDataLength(); j++ {
				stmt.Exec(us.Idtype(), us.UnitId(), uint32(j), uint64(us.ValueData(j)), "UShort")
				rows++
			}
		}

		/* --- 6. IntTable --- */
		it := new(fb.IntSaveDataUnit)
		for i := 0; i < root.IntTableLength(); i++ {
			if !root.IntTable(it, i) {
				continue
			}
			for j := 0; j < it.ValueDataLength(); j++ {
				v := int32(it.ValueData(j))
				stmt.Exec(it.Idtype(), it.UnitId(), uint32(j), uint64(int64(v)), "Int")
				rows++
			}
		}

		/* --- 7. UIntTable --- */
		ui := new(fb.UIntSaveDataUnit)
		for i := 0; i < root.UintTableLength(); i++ {
			if !root.UintTable(ui, i) {
				continue
			}
			for j := 0; j < ui.ValueDataLength(); j++ {
				stmt.Exec(ui.Idtype(), ui.UnitId(), uint32(j), uint64(ui.ValueData(j)), "UInt")
				rows++
			}
		}

		/* --- 8. LongTable --- */
		lo := new(fb.LongSaveDataUnit)
		for i := 0; i < root.LongTableLength(); i++ {
			if !root.LongTable(lo, i) {
				continue
			}
			for j := 0; j < lo.ValueDataLength(); j++ {
				v := lo.ValueData(j)
				stmt.Exec(lo.Idtype(), lo.UnitId(), uint32(j), uint64(v), "Long")
				rows++
			}
		}

		/* --- 9. ULongTable --- */
		ul := new(fb.ULongSaveDataUnit)
		for i := 0; i < root.UlongTableLength(); i++ {
			if !root.UlongTable(ul, i) {
				continue
			}
			for j := 0; j < ul.ValueDataLength(); j++ {
				stmt.Exec(ul.Idtype(), ul.UnitId(), uint32(j), ul.ValueData(j), "ULong")
				rows++
			}
		}

		/* ---10. FloatTable --- */
		ft := new(fb.FloatSaveDataUnit)
		for i := 0; i < root.FloatTableLength(); i++ {
			if !root.FloatTable(ft, i) {
				continue
			}
			for j := 0; j < ft.ValueDataLength(); j++ {
				bits64 := ft.ValueData(j)
				bits32 := uint32(bits64)
				fv := math.Float32frombits(bits32)
				_ = fv
				if _, err := stmt.Exec(
					ft.Idtype(), ft.UnitId(), uint32(j),
					uint64(bits32),
					"Float",
				); err != nil {
					log.Fatalf("insert float row: %v", err)
				}
				rows++
			}
		}
		return rows
	}

	stmt, err := db.Prepare(`INSERT INTO save_units(idtype,unitid,val_index,value,table_ty)
		VALUES(?,?,?,?,?)
		ON DUPLICATE KEY UPDATE value=VALUES(value)`)
	if err != nil {
		log.Fatalf("prepare: %v", err)
	}
	defer stmt.Close()

	total := process(seg1, stmt) + process(slot, stmt)
	fmt.Printf("Import completed, %d rows written\n", total)
}

func exportCSV(db *sql.DB, outfile string) {
	f, err := os.Create(outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.Write([]string{"id", "idtype", "unitid", "val_index", "value", "table_ty"})

	rows, err := db.Query(`SELECT id, idtype, unitid, val_index, value, table_ty
	                       FROM save_units ORDER BY id`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, idt, uid, vidx int64
		var val uint64
		var ty string
		if err := rows.Scan(&id, &idt, &uid, &vidx, &val, &ty); err != nil {
			log.Fatal(err)
		}
		rec := []string{
			strconv.FormatInt(id, 10),
			strconv.FormatInt(idt, 10),
			strconv.FormatInt(uid, 10),
			strconv.FormatInt(vidx, 10),
			strconv.FormatUint(val, 10),
			ty,
		}
		w.Write(rec)
	}
	w.Flush()
	fmt.Println("CSV export completed →", outfile)
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s import <SaveData.dat> | export <out.csv>", os.Args[0])
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	ensureTable(db)

	switch os.Args[1] {
	case "import":
		importBin(db, os.Args[2])
	case "export":
		exportCSV(db, os.Args[2])
	default:
		log.Fatalf("Unknown args %s", os.Args[1])
	}
}
