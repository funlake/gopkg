package acm

import (
	"fmt"
	"testing"
)

func TestConfigration_GetServers(t *testing.T) {
	c := &configration{
		ack: "e07a5cfe85ad440e9ded33db25e7236f",
		sek: "96UVhMYbVhredzdLmr5UvQ6NjZI=",
		ep:  "acm.aliyun.com",
		pt:  "8080",
		ns:  "a03fb67c-0c95-4d8d-8f1e-671e80759f9a",
		gp:  "DEFAULT_GROUP",
	}
	//t.Log(c.GetServer())
	c.GetConfig("com.etcchebao.acm:base.properties", "DEFAULT_GROUP")
	fmt.Printf("%s", []string{"hello", "hi", "o", ""}[:len([]string{"hello", "hi", "o", ""})-1])
}
