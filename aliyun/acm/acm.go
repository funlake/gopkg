package acm

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/funlake/gopkg/utils/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	url2 "net/url"
	"strings"
	"time"
)

type configration struct {
	ack string //access key
	sek string //secure key
	ep  string // endpoint
	gp  string // group
	ns  string // namespace
	pt  string // port
}

func int() {
	rand.Seed(time.Now().UnixNano())
}
func (c *configration) GetServers() []string {
	url := "http://" + c.ep + ":" + c.pt + "/diamond-server/diamond"
	resp, err := http.Get(url)
	servers := []string{}
	if err != nil {
		log.Error(err.Error())
	} else {
		c, _ := ioutil.ReadAll(resp.Body)
		servers = strings.Split(string(c), "\n")
		//strip empty line
		servers = servers[:len(servers)-1]
		defer resp.Body.Close()
	}
	return servers
}

func (c *configration) GetServer() string {
	servers := c.GetServers()
	if len(servers) < 1 {
		panic("No available servers")
	} else {
		return servers[rand.Intn(len(servers))]
	}
	return ""
}
func (c *configration) GetSign(params []string) string {
	mac := hmac.New(sha1.New, []byte(c.sek))
	mac.Write([]byte(strings.Join(params, "+")))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
func (c *configration) GetConfig(dataId string, group string) {
	server := c.GetServer()
	//dataid := strings.Replace(dataId,":","%3A",1)
	dataid := url2.QueryEscape(dataId)
	group = url2.QueryEscape(group)
	ns := url2.QueryEscape(c.ns)
	url := fmt.Sprintf("http://%s:%s/diamond-server/config.co?dataId=%s&group=%s&tenant=%s", server, c.pt, dataid, group, ns)
	req, _ := http.NewRequest("GET", url, nil)
	//req.Header.Set("Diamond-Client-AppName","ACM-SDK-GO")
	//req.Header.Set("Client-Version","0.0.1")
	//req.Header.Set("Content-Type","application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("ExConfigInfo", "true")
	req.Header.Set("Spas-AccessKey", c.ack)
	ts := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	req.Header.Set("TimeStamp", ts)
	//s := c.ns + "+" + c.gp + "+" + ts
	//mac := hmac.New(sha1.New,[]byte(c.sek))
	//mac.Write([]byte(s))
	//sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	sign := c.GetSign([]string{c.ns, c.gp, ts})
	req.Header.Set("Spas-Signature", sign)
	//log.Success("%s",sign)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	cc, _ := ioutil.ReadAll(resp.Body)
	log.Success("%s,%s , %s", url, req.Header, cc)
}
