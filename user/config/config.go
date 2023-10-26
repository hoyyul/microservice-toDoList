package config

type Config struct {
	Server Server `yaml:server`
	Mysql  Mysql  `yaml:mysql`
	Etcd   Etcd   `yaml:etcd`
	Redis  Redis  `yaml:redis`
}

type Server struct {
	Port    string `yaml:port`
	ENV     string `yaml:env`
	Version string `yaml:version`
}
type Mysql struct {
	Host     string `yaml:host`
	Port     string `yaml:port`
	Database string `yaml:database`
	Username string `yaml:username `
	Password string `yaml:password`
	Charset  string `yaml:charset`
}

type Etcd struct {
	Address string `yaml:address`
}

type Redis struct {
	Address  string `yaml:redis`
	Password string `yaml:password`
}
