package tracer

import (
	"go.opentelemetry.io/otel/attribute"
	"net"
	"net/http"
	"strconv"
)

var (
	SystemDbPostgres = attribute.String("db.system", "postgres")
	SystemDbRedis    = attribute.String("db.system", "redis")
)

func NewHttpAttributes(network string, request *http.Request) []attribute.KeyValue {
	attributes := make([]attribute.KeyValue, 11)

	switch network {
	case "tcp", "tcp4", "tcp6":
		attributes = append(attributes, attribute.String("network.transport", "tcp"))
	case "udp", "udp4", "udp6":
		attributes = append(attributes, attribute.String("network.transport", "udp"))
	}

	peerIP, _, peerPort := hostIPNamePort(request.RemoteAddr)
	if peerIP != "" {
		attributes = append(attributes, attribute.String("network.peer.address", peerIP))
	}
	if peerPort != 0 {
		attributes = append(attributes, attribute.Int("network.peer.port", peerPort))
	}

	attributes = append(attributes, attribute.String("http.request.method", request.Method))
	attributes = append(attributes, attribute.String("user_agent.original", request.UserAgent()))
	attributes = append(attributes, attribute.Int64("http.request.header.content-length", request.ContentLength))
	attributes = append(attributes, attribute.String("url.scheme", request.URL.Scheme))
	attributes = append(attributes, attribute.String("url.path", request.URL.Path))
	attributes = append(attributes, attribute.String("url.query", request.URL.RawQuery))

	serverIP, serverName, serverPort := hostIPNamePort(request.Host)
	if serverIP != "" {
		attributes = append(attributes, attribute.String("server.address", serverIP))
	} else if serverName != "" {
		attributes = append(attributes, attribute.String("server.address", serverName))
	}
	if serverPort != 0 {
		attributes = append(attributes, attribute.Int("server.port", serverPort))
	}

	return attributes
}

func hostIPNamePort(hostWithPort string) (ip string, name string, port int) {
	var (
		hostPart, portPart string
		parsedPort         uint64
		err                error
	)
	if hostPart, portPart, err = net.SplitHostPort(hostWithPort); err != nil {
		hostPart, portPart = hostWithPort, ""
	}
	if parsedIP := net.ParseIP(hostPart); parsedIP != nil {
		ip = parsedIP.String()
	} else {
		name = hostPart
	}
	if parsedPort, err = strconv.ParseUint(portPart, 10, 16); err == nil {
		port = int(parsedPort)
	}
	return
}
