// simple tool to parse ip to geo locations
package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
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
		Name:  "iplookup",
		Usage: "resolve ip addresses to countries",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "file name for IP list",
			},
			&cli.StringFlag{
				Name:    "language",
				Value:   "en",
				Aliases: []string{"l"},
				Usage:   "language for country name",
			},
			&cli.StringFlag{
				Name:    "domain",
				Value:   "",
				Aliases: []string{"d"},
				Usage:   "domain to resolve for ip",
			},
		},
		Action: commandHandler,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func commandHandler(c *cli.Context) error {
	ip := c.Args().Get(0)
	lang := c.String("l")
	db, err := geoip2.FromBytes(ipDB)
	if err != nil {
		return err
	}

	if domain := c.String("d"); domain != "" {
		ps, err := net.LookupIP(domain)
		if err != nil {
			return err
		}
		for _, ip := range ps {
			country, err := singleIP(db, ip.String(), lang)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(country.country, ":", ip)
		}
	} else if ip != "" {
		country, err := singleIP(db, ip, lang)
		if err != nil {
			return err
		}
		fmt.Println(country.country)
	} else if isInputFromPipe() {
		countries, err := batchIP(db, os.Stdin, lang)
		if err != nil {
			return err
		}
		printCountries(countries)
	} else if c.String("f") != "" {
		file, err := os.Open(c.String("file"))
		if err != nil {
			return err
		}
		defer file.Close()
		countries, err := batchIP(db, file, lang)
		if err != nil {
			return err
		}
		printCountries(countries)
	} else {
		fmt.Println("Please read the usage.")
	}

	return nil
}

func printCountries(countries []ipCountry) {
	counts := make(map[string]int)
	for _, country := range countries {
		counts[country.country]++
	}
	for country, count := range counts {
		fmt.Printf("%s: %d\n", country, count)
	}
}

func singleIP(db *geoip2.Reader, ipString string, language string) (ipCountry, error) {
	ip := net.ParseIP(ipString)
	record, err := db.Country(ip)
	if err != nil {
		return ipCountry{}, err
	}
	return ipCountry{ipString, record.Country.Names[language]}, nil
}

type ipCountry struct {
	ip      string
	country string
}

func batchIP(db *geoip2.Reader, r io.Reader, language string) ([]ipCountry, error) {
	scanner := bufio.NewScanner(r)
	var countries []ipCountry

	for scanner.Scan() {
		ip := scanner.Text()
		name, err := countryName(db, ip, language)
		if err != nil {
			return nil, err
		}
		countries = append(countries, ipCountry{ip, name})
	}
	return countries, nil
}

func countryName(db *geoip2.Reader, ipAddr string, language string) (string, error) {
	ip := net.ParseIP(ipAddr)
	record, err := db.Country(ip)
	if err != nil {
		return "", err
	}
	return record.Country.Names[language], nil
}

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}
