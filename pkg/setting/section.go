package setting

import "time"

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
type AppSettingS struct {
	DefaultPageSize int
	MaxPageSize     int
	LogSavePath     string
	LogFileName     string
	LogFileExt      string
}
type DatabaseSettingS struct {
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

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
