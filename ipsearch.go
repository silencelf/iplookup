package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "IpSearch",
		Usage: "resolve ip addresses to contries",
		Action: func(c *cli.Context) error {
			arg := c.Args().Get(0)
			fmt.Println(arg)
			file, err := os.Open("ip.txt")
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			counts := make(map[string]int)
			db, err := geoip2.Open("country.mmdb")
			if err != nil {
				return err
			}
			defer db.Close()

			for scanner.Scan() {
				ip := scanner.Text()
				name := printCountry(db, ip)
				counts[name]++
			}

			fmt.Println(counts)
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func printCountry(db *geoip2.Reader, ipAddr string) string {
	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(ipAddr)
	record, err := db.Country(ip)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Chinese country name: %v\n", record.Country.Names["zh-CN"])
	//fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
	//fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
	return record.Country.Names["zh-CN"]
}
