package setting

import "time"

type DatabaseSettingS struct {
	UserName     string
	Password     string
	Host         string
	DBName       string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
}

type RedisSettingS struct {
	Network string
	Address string
	Port    string
}

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type AppSettingS struct {
	DefaultPageSize       int
	MaxPageSize           int
	DefaultContextTimeout time.Duration
	LogSavePath           string
	LogFileName           string
	LogFileExt            string
	CsvSavePath           string
	DoneCsvSavePath       string
	LogV2SavePath         string
	LogV2FileName         string
	UploadSavePath        string
	UploadServerUrl       string
	UploadImageMaxSize    int
	UploadImageAllowExts  []string
}

var sections = make(map[string]interface{})

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}
