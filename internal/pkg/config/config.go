package config

type Configuration struct {
	Log struct {
		Level string
		Type  string
	}
	Database struct {
		Host string
		User string
		Pass string
		Name string
	}
}
