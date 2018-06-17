package main 

import (
	"net/http"
	"io/ioutil"
	"time"
  "strings"
  "bytes"
  "crypto/tls"
//	"log"

)

type HttpRequest struct {
  Url string `json:"url"`
  Method string `json:"method"`
  Headers []string `json:"headers"`
  Body  string `json:"body"`
}

type HttpResponse struct {
  Status string `json:"status"`
  Headers []string `json:"headers"`
  Body string `json:"body"`
}

func doProxyRequest (reqData HttpRequest) HttpResponse {
  req, _ := http.NewRequest (reqData.Method, reqData.Url, nil)
  for _, h := range reqData.Headers {
    div := strings.Index(h, ":")
    key := h[0:div]
    val := h[div + 2:]
    req.Header.Add(key, val)
  }
  var respData HttpResponse
	netClient := &http.Client{
		Timeout: 10 * time.Second,
    Transport: &http.Transport{// disable certifcate checking
      TLSClientConfig: &tls.Config{InsecureSkipVerify: true,},
    },    
	}
	resp, _ := netClient.Do(req)
  if resp != nil {
    respData.Status = resp.Status
    buf := new(bytes.Buffer)
    buf.ReadFrom(resp.Body)
    respData.Body = buf.String()
    var val string ;
    for k, v := range resp.Header{
      val = ""
      for _, vi := range v {
        val += vi
      }
      respData.Headers = append(respData.Headers, k + ": " + val)
    } 
  } else {
    lo("null response")
    respData.Headers = append(respData.Headers, "one: two")
  }
  return respData
}

func checkIpAddress (ip string) bool {
	url := "http://check.getipintel.net/check.php?ip=" + ip + "&contact=GeneKnight4@GMail.com&flags=m"
//	log.Println(url)
	var netClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, _ := netClient.Get(url)
	buf, _ := ioutil.ReadAll(resp.Body) // returns array of bytes
	return buf[0] == 49 // == '1'
}
