package go_utils

import (
	"adnpa/go-utils/pkg/network"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	i := network.IpStr2Int("192.168.0.1")
	fmt.Println(i)
	fmt.Println(network.IpInt2Str(i))
}
