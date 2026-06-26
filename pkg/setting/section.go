package setting

import (
	"time"
)

type ServerSetting struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
type AppSetting struct {
	DefaultPageSize       int
	MaxPageSize           int
	DefaultContextTimeout time.Duration
	LogSavePath           string
	LogFileName           string
	LogFileExt            string
	UploadSavePath        string
	UploadServerUrl       string
	UploadImageMaxSize    int
	UploadImageAllowExts  []string
}
type EmailSetting struct {
	Host     string
	Port     int
	UserName string
	Password string
	IsSSL    bool
	From     string
	To       []string
}
type JWTSetting struct {
	Secret string
	Issuer string
	Expire time.Duration
}
type DatabaseSetting struct {
	DBType       string `mapstructure:"DBType"`
	Username     string `mapstructure:"Username"`
	Password     string `mapstructure:"Password"`
	Host         string `mapstructure:"Host"`
	Port         int    `mapstructure:"Port"`
	DBName       string `mapstructure:"DBName"`
	TablePrefix  string `mapstructure:"TablePrefix"`
	Charset      string `mapstructure:"Charset"`
	ParseTime    bool   `mapstructure:"ParseTime"`
	MaxIdleConns int    `mapstructure:"MaxIdleConns"`
	MaxOpenConns int    `mapstructure:"MaxOpenConns"`
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
func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
