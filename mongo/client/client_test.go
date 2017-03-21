package client

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	addresses, err := splitHosts("127.0.0.1:1234, 172.16.0.1,192.168.1.1", "9999")

	fmt.Println("==============")
	fmt.Println(err)
	fmt.Println(addresses)
	fmt.Println("--------------")
}
