package b_rawEvent

import (
	"fmt"
	"github.com/Sapomie/wayne-data/global"
	"github.com/gocarina/gocsv"
	"io/ioutil"
	"os"
	"strings"
)

func getRawEventFromCsvFile() ([]*RawEvent, string, error) {
	files, filesRemove, err := getFileAndRemove()
	if err != nil {
		return nil, "", err
	}
	var rawEvents []*RawEvent
	var count int
	for i, dir := range files {
		in, err := os.Open(dir)
		if err != nil {
			return nil, "", err
		}
		rEvents := make([]*RawEvent, 0)
		if err = gocsv.UnmarshalFile(in, &rEvents); err != nil {
			return nil, "", err
		}
		for _, rEvent := range rEvents {
			rawEvents = append(rawEvents, rEvent)
			count++
		}
		err = in.Close()
		if err != nil {
			return nil, "", err
		}
		err = os.Rename(files[i], filesRemove[i])
		if err != nil {
			return nil, "", err
		}
	}

	info := fmt.Sprintf("getting %d raw Events from Csv files", count)

	return rawEvents, info, nil
}

func getFileAndRemove() (csvFiles []string, csvFilesRename []string, err error) {
	rd, err := ioutil.ReadDir(global.AppSetting.CsvSavePath)
	for _, fi := range rd {
		if !strings.HasSuffix(fi.Name(), ".csv") {
			continue
		} else {
			filename := fmt.Sprint(global.AppSetting.CsvSavePath + "/" + fi.Name())
			csvFiles = append(csvFiles, filename)
			fileNameRemoved := fmt.Sprint(global.AppSetting.DoneCsvSavePath + "/" + fi.Name())
			csvFilesRename = append(csvFilesRename, fileNameRemoved)
		}
	}
	return
}
