package getter

import (
	"github.com/cwchiu/ProxyPool/models"
	"github.com/parnurzeal/gorequest"
	"log"
	"regexp"
)

func USProxyOrg() (result []*models.IP) {
	pollURL := "https://www.us-proxy.org/"
	_, body, errs := gorequest.New().Get(pollURL).End()
	if errs != nil {
		log.Println(errs)
		return
	}
	// ip,port,code,country,Anonymity,google,Https
	re := regexp.MustCompile(`<tr><td>(.*?)</td><td>(.*?)</td>.*?<td class='\w+'>(no|yes)</td>.*?</tr>`)
	ret := re.FindAllStringSubmatch(body, -1)
	// resp.Body.Close()

	for _, match := range ret {
		ip := models.NewIP()
		ip.Data = match[1] + ":" + match[2]
		proxy_type := "HTTP"
		if match[3] == "yes" {
			proxy_type = "HTTPS"
		}
		ip.Type = proxy_type
		result = append(result, ip)
	}

	log.Println("us-proxy done.")
	return
}
