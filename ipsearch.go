package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
	"github.com/urfave/cli/v2"
)

//go:embed Country.mmdb
var ipDB []byte

func main() {
	app := &cli.App{
		Name:  "ipsearch",
		Usage: "resolve ip addresses to contries",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "file name for ips",
			},
		},
		Action: func(c *cli.Context) error {
			ip := c.Args().Get(0)
			if ip != "" {
				country, err := singleIP(ip)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(country)
			} else if c.String("f") != "" {
				batchIP(c.String("f"))
			} else {
				fmt.Println("Please read the usage.")
			}

			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func singleIP(ipString string) (string, error) {
	db, err := geoip2.FromBytes(ipDB)
	if err != nil {
		return "", err
	}

	ip := net.ParseIP(ipString)
	record, err := db.Country(ip)
	if err != nil {
		return "", err
	}
	return record.Country.Names["zh-CN"], nil
}

func batchIP(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	counts := make(map[string]int)
	db, err := geoip2.FromBytes(ipDB)
	if err != nil {
		return err
	}
	defer db.Close()

	for scanner.Scan() {
		ip := scanner.Text()
		name, err := printCountry(db, ip)
		if err != nil {
			fmt.Println(err)
		}
		counts[name]++
	}
	fmt.Println(counts)
	return nil
}

func printCountry(db *geoip2.Reader, ipAddr string) (string, error) {
	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(ipAddr)
	record, err := db.Country(ip)
	if err != nil {
		return "", err
	}
	return record.Country.Names["zh-CN"], nil
}
