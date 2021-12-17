package b_raw_event

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"io/ioutil"
	"os"
	"strings"
)

func (svc RawEventService) getRawEventFromCsvFile() ([]*RawEvent, string, error) {
	files, filesRemove, err := svc.getFileAndRemove()
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

func (svc RawEventService) getFileAndRemove() (csvFiles []string, csvFilesRename []string, err error) {
	rd, err := ioutil.ReadDir(svc.appSetting.CsvSavePath)
	for _, fi := range rd {
		if !strings.HasSuffix(fi.Name(), ".csv") {
			continue
		} else {
			filename := fmt.Sprint(svc.appSetting.CsvSavePath + "/" + fi.Name())
			csvFiles = append(csvFiles, filename)
			fileNameRemoved := fmt.Sprint(svc.appSetting.DoneCsvSavePath + "/" + fi.Name())
			csvFilesRename = append(csvFilesRename, fileNameRemoved)
		}
	}
	return
}
