package network

import (
	"encoding/binary"
	"net"
)

func IpStr2Int(ipStr string) uint32 {
	ip := net.ParseIP(ipStr)
	return binary.BigEndian.Uint32(ip.To4())
}

func IpInt2Str(i uint32) string {
	return net.IPv4(byte(i>>24), byte(i>>16&0xFF), byte(i>>8&0xFF), byte(i&0xFF)).String()
}
