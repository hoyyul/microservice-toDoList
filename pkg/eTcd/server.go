package eTcd

import "fmt"

type Server struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

func BuildPrefix(server Server) string {
	return fmt.Sprintf("/%s/", server.Name)
}

func BuildRegisterPath(server Server) string {
	return fmt.Sprintf("%s%s", BuildPrefix(server), server.Addr)
}
