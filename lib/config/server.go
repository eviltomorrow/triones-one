package config

import jsoniter "github.com/json-iterator/go"

type Server struct {
	GrpcAddressList []string `mapstructure:"grpc-address-list" json:"grpc-address-list"`
	GrpcPort        int      `mapstructure:"grpc-port" json:"grpc-port"`
	TlsCaFile       string   `mapstructure:"tls-ca-file" json:"tls-ca-file"`
	TlsServerCert   string   `mapstructure:"tls-server-cert" json:"tls-server-cert"`
	TlsServerPem    string   `mapstructure:"tls-server-pem" json:"tls-server-pem"`
}

func (c *Server) String() string {
	buf, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(c)
	return string(buf)
}
