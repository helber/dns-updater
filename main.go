package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

func main() {
	verb := flag.Bool("v", false, "Verbose")
	flag.Parse()
	host := os.Getenv("A_HOST")
	if host == "" {
		log.Fatal("ENV var A_HOST need suppliend")
	}
	// get external ip addr
	externalIP := os.Getenv("HOST_IP_GET")
	if externalIP == "" {
		externalIP = "http://ifconfig.io/ip;http://ifconfig.co/ip"
	}
	extIP := ""
	sites := strings.Split(externalIP, ";")
	for _, site := range sites {
		if *verb {
			log.Println("Call -> ", site)
		}
		resp, err := http.Get(site)
		if err != nil {
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		rextIP := strings.TrimSpace(string(body))
		trial := net.ParseIP(rextIP).To4()
		if trial == nil {
			if *verb {
				log.Println("invalid IP")
			}
			continue
		}
		extIP = rextIP
		if *verb {
			log.Println("My external IP", extIP)
		}
	}
	if extIP == "" {
		log.Fatal("Can not determine external IP address")
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Fatal(err)
	}
	if *verb {
		log.Println(ips)
	}
	dnsIP := strings.TrimSpace(ips[0].String())
	if dnsIP == extIP {
		if *verb {
			log.Println("IPs not cahnged")
		}
		os.Exit(0)
	}
	// CloudFlare API
	api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the zone ID
	id, err := api.ZoneIDByName(host)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch records
	records, err := api.DNSRecords(id, cloudflare.DNSRecord{Type: "A", Name: host})
	if err != nil {
		log.Fatal(err)
	}
	for i, rec := range records {
		log.Println(i, rec)
		if rec.Content == extIP {
			log.Fatal("record not propaged yet")
			os.Exit(0)
		}
		err := api.UpdateDNSRecord(id, rec.ID, cloudflare.DNSRecord{Content: extIP, Type: "A", Name: host})
		if err != nil {
			log.Fatal("change record", err)
		}
	}
}
