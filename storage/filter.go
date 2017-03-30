package storage

import (
	"ProxyPool/models"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
)

// CheckProxy .
func CheckProxy(ip string) {
	if CheckIP(ip) {
		ProxyAdd(ip)
	}
}

// CheckIP is to check the ip work or not
func CheckIP(ip string) bool {
	pollURL := "http://www.baidu.com"
	resp, _, errs := gorequest.New().Proxy("http://" + ip).Get(pollURL).End()
	if errs != nil {
		return false
	}
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

// CheckProxyDB to check the ip in DB
func CheckProxyDB() {
	conn := NewStorage()
	x := conn.Count()
	log.Println("Before check, DB has:", x, "records.")
	ips, err := conn.GetAll()
	if err != nil {
		log.Println(err.Error())
		return
	}
	var wg sync.WaitGroup
	for _, v := range ips {
		wg.Add(1)
		go func(v *models.IP) {
			if !CheckIP(v.Data) {
				ProxyDel(v)
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
	x = conn.Count()
	log.Println("After check, DB has:", x, "records.")
}

// ProxyRandom .
func ProxyRandom() (ip *models.IP) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	conn := NewStorage()
	ips, _ := conn.GetAll()
	x := conn.Count()

	return ips[r.Intn(x)]
}

// ProxyAdd .
func ProxyAdd(ip string) {
	conn := NewStorage()
	_, err := conn.GetOne(ip)
	if err != nil {
		things := models.NewIP()
		things.Data = ip
		conn.Create(things)
	}
}

// ProxyDel .
func ProxyDel(ip *models.IP) {
	conn := NewStorage()
	if err := conn.Delete(ip); err != nil {
		log.Println(err.Error())
	}
}
