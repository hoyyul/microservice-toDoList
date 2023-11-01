package eTcd

import (
	"encoding/json"
	"fmt"
)

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

func ParseValue(data []byte) (Server, error) {
	var server Server
	err := json.Unmarshal(data, &server)
	if err != nil {
		return server, err
	}
	return server, nil
}
