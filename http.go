package main 

import (
	"net/http"
	"io/ioutil"
	"time"
//	"log"

)

func checkIpAddress (ip string) bool {
	url := "http://check.getipintel.net/check.php?ip=" + ip + "&contact=GeneKnight4@GMail.com&flags=m"
//	log.Println(url)
	var netClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, _ := netClient.Get(url)
	buf, _ := ioutil.ReadAll(resp.Body) // returns array of bytes
	return buf[0] == 49 // == '1'
//	log.Println(string(buf))
//	log.Println(err)
}
