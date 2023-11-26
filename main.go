package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	_ "embed"

	"github.com/cloudflare/cloudflare-go"
)

//go:embed .secrets.json
var secretFile []byte

var (
	api *cloudflare.API
	err error
)

func init() {
	var data secrets
	var email, key string
	json.Unmarshal(secretFile, &data)
	key = os.Getenv("CF_API_KEY")
	email = os.Getenv("CF_EMAIL")
	if key == "" || email == "" {
		key = data.Token
		email = data.Email
	}
	api, err = cloudflare.New(key, email)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func main() {
	dnsType := flag.String("type", "", "The type of the DNS record (A, CNAME, etc.)")
	dnsName := flag.String("name", "", "The name of the DNS record")
	dnsContent := flag.String("content", "", "The content of the DNS record")
	dnsTTL := flag.Int("ttl", 120, "The TTL of the DNS record")
	zoneID := flag.String("zoneid", "", "The ID of the zone to add the DNS record to")

	flag.Parse()

	var zoneTld string

	if zoneID == nil || *zoneID == "" {
		tld := strings.Split(*dnsName, ".")
		if len(tld) < 2 {
			os.Exit(1)
		}
		zoneTld = tld[1] + "." + tld[2]
		log.Printf("the tld is: %s", zoneTld)
		*zoneID, err = getZoneID(zoneTld)
		if err != nil {
			os.Exit(1)
		}
		log.Printf("the zone is: %s", zoneID)
	}

	params := cloudflare.CreateDNSRecordParams{
		Type:    *dnsType,
		Name:    *dnsName,
		Content: *dnsContent,
		TTL:     *dnsTTL,
	}

	rc := &cloudflare.ResourceContainer{Identifier: *zoneID}
	_, err = api.CreateDNSRecord(context.TODO(), rc, params)
	if err != nil {
		fmt.Printf("error creating DNS record: %s", err)
		os.Exit(1)
	}

	fmt.Println("DNS record created successfully.")
}

func getZoneID(domain string) (string, error) {
	zones, err := api.ListZones(context.Background())
	if err != nil {
		return "", err
	}

	for _, zone := range zones {
		if zone.Name == domain {
			return zone.ID, nil
		}
	}

	return "", fmt.Errorf("no zone found for domain: %s", domain)
}

type secrets struct {
	Email string `json:"email"`
	Token string `json:"api_key"`
}