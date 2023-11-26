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
	"github.com/fatih/color"
	"golang.org/x/term"
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
	if key == "" || email == "" {
		fmt.Print("Enter your Cloudflare email: ")
		fmt.Scanln(&email)
		color.Green("Enter your Cloudflare API key: (input will be hidden)")

		byteKey, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			color.Red(err.Error())
		}
		key = string(byteKey)
		fmt.Println()
	}
	log.Printf("the email is: %s - api key is: %s", email, key)
	api, err = cloudflare.New(key, email)
	if err != nil {
		log.Printf("error in cf init: %s", err)
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
			color.Red("invalid domain")
			os.Exit(1)
		}
		zoneTld = tld[1] + "." + tld[2]
		log.Printf("the tld is: %s", zoneTld)
		*zoneID, err = getZoneID(zoneTld)
		if err != nil {
			color.Red(err.Error())
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

	color.Green("DNS record created successfully.")
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
