package aggregator

import (
	"fmt"
	"log"
	"os"
	"time"

	"../common"
)

type FileWriter struct {
	Files      map[string]*os.File
	DateString string
}

func NewFileWriter() *FileWriter {
	f := &FileWriter{
		Files:      make(map[string]*os.File),
		DateString: dateString(),
	}
	return f
}

func (T *FileWriter) WriteData(datum *common.Datum) {
	path := fmt.Sprintf("./data/ticker/%s/", datum.Symbol)
	if _, ok := T.Files[datum.Symbol]; !ok || T.DateString != dateString() {
		T.DateString = dateString()
		T.Files[datum.Symbol] = openFile(path, T.DateString)
	}
	f := T.Files[datum.Symbol]
	serialized := serialize(datum)
	_, err := f.WriteString(serialized)
	if err != nil {
		log.Printf("error: could not write to file (%+v)\n", err)
		return
	}
}

func serialize(datum *common.Datum) string {
	serialized := fmt.Sprintf("%+v,", datum)
	return serialized
}

func dateString() string {
	y, m, d := time.Now().Date()
	dateString := fmt.Sprintf("%d_%d_%d.data", y, m, d)
	return dateString
}

func openFile(dir, filename string) *os.File {
	var file *os.File
	var err error
	fullPath := dir + filename
	if _, err = os.Stat(fullPath); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			log.Fatalf("error: could not create path %s (%+v)\n", fullPath, err)
		}
		file, err = os.Create(fullPath)
		if err != nil {
			log.Fatalf("error: could not create new file (%+v)\n", err)
		}
	} else {
		file, err = os.OpenFile(fullPath, os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error: could not open data file (%+v)\n", err)
		}
	}
	return file
}
