package config

type AppConfig struct {
	LogLevel    string `json:"log_level"`
	Port        int    `json:"port"`
	ConfigPath  string `json:"config_path"`
	StoragePath string `json:"storage_path"`
	TmpPath     string `json:"tmp_path"`
	DataPath    string `json:"data_path"`
}
