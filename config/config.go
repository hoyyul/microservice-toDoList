package config

type Config struct {
	Server   *Server             `yaml:"server"`
	Mysql    *Mysql              `yaml:"mysql"`
	Etcd     *Etcd               `yaml:"etcd"`
	Redis    *Redis              `yaml:"redis"`
	Services map[string]*Service `yaml:"services"`
}

type Server struct {
	Address string `yaml:"address"`
	Env     string `yaml:"env"`
	Jwt     string `yaml:"jwt"`
	Version string `yaml:"version"`
}
type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
}

type Etcd struct {
	Address string `yaml:"address"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	PoolSize int    `yaml:"pool_size"`
}

type Service struct {
	Name        string `yaml:"name"`
	LoadBalance bool   `yaml:"loadBalance"`
	Address     string `yaml:"address"`
}
