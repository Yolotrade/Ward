package aggregator

import (
	"fmt"
	"log"
	"os"

	"../common"
)

type FileWriter struct {
	Files map[string]*os.File
}

func NewFileWriter() *FileWriter {
	f := &FileWriter{
		Files: make(map[string]*os.File),
	}
	return f
}

func (T *FileWriter) WriteData(datum *common.Datum) {
	path := fmt.Sprintf("data/%s.data", datum.Symbol)
	if _, ok := T.Files[datum.Symbol]; !ok {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			T.Files[datum.Symbol], err = os.Create(path)
			if err != nil {
				log.Printf("error: could not create new file (%+v)\n", err)
				return
			}
		} else {
			T.Files[datum.Symbol], err = os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				log.Printf("error: could not open data file (%+v)\n", err)
				return
			}
		}
	}
	f := T.Files[datum.Symbol]
	_, err := f.WriteString(fmt.Sprintf("%+v,", datum))
	if err != nil {
		log.Printf("error: could not write to file (%+v)\n", err)
		return
	}
}
