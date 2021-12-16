package db

import (
	"flag"
	"fmt"
)

var (
	user string
	pass string
	host string
	port int
	name string
)

func LoadConfigsFromFlags(flagSet *flag.FlagSet) {
	flagSet.StringVar(&user, "db_user", "guest", "set database username")
	flagSet.StringVar(&pass, "db_pass", "guest", "set database password")
	flagSet.StringVar(&host, "db_host", "localhost", "set database host")
	flagSet.IntVar(&port, "db_port", 3306, "set database port")
	flagSet.StringVar(&name, "db_name", "sandbox", "set database name")
}

type Config interface {
	Source() string
}

type config struct {
	user string
	pass string
	host string
	port int
	name string
}

func NewConfig() Config {
	return &config{
		user: user,
		pass: pass,
		host: host,
		port: port,
		name: name,
	}
}

func (c *config) Source() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", c.user, c.pass, c.host, c.port, c.name)
}
