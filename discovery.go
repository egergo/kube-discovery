package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/pflag"
)

type Address struct {
	IP string `json:"ip"`
}

type Port struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

type Subset struct {
	Addresses []Address `json:"addresses"`
	Ports     []Port    `json:"ports"`
}

type Resp struct {
	Subsets []Subset `json:"subsets"`
}

type Options struct {
	Seed      []string
	SeedURL   []*url.URL
	Namespace string
	Service   string
	Timeout   int64
	Help      bool
}

func NewOptions() *Options {
	return &Options{
		Namespace: "default",
		Service:   "kubernetes",
		Timeout:   10000,
	}
}

func checkSeed(u *url.URL, opt *Options) string {
	if len(u.Path) == 0 || u.Path[len(u.Path)-1] != '/' {
		u.Path = u.Path + "/"
	}
	u.Path = fmt.Sprintf("%sapi/v1/namespaces/%s/endpoints/%s", u.Path, url.QueryEscape(opt.Namespace), url.QueryEscape(opt.Service))
	addr := u.String()

	log.Print(addr, " Probing URL")

	timeout := time.Duration(opt.Timeout) * time.Millisecond
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(addr)
	if err != nil {
		log.Print(addr, " Could not load URL: ", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Print(addr, " Invalid response code: ", resp.StatusCode)
		return ""
	}

	decoder := json.NewDecoder(resp.Body)
	var r Resp
	err = decoder.Decode(&r)
	if err != nil {
		log.Print(addr, " Cannot decode JSON: ", err)
		return ""
	}

	valid := make([]Subset, 0, len(r.Subsets))
	for _, subset := range r.Subsets {
		if len(subset.Addresses) > 0 && len(subset.Ports) > 0 {
			valid = append(valid, subset)
		}
	}

	if len(valid) == 0 {
		log.Print(addr, " No valid subset found")
		return ""
	}

	chosenSubset := valid[rand.Intn(len(valid))]
	chosenAddress := chosenSubset.Addresses[rand.Intn(len(chosenSubset.Addresses))]
	chosenPort := chosenSubset.Ports[0]

	return fmt.Sprintf("%s:%d\n", chosenAddress.IP, chosenPort.Port)
}

func checkSeedAsync(u *url.URL, opt *Options, ch chan<- string) {
	ch <- checkSeed(u, opt)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetOutput(os.Stderr)

	opt := NewOptions()

	var fs = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	fs.StringSliceVar(&opt.Seed, "seed", opt.Seed, "Server seeds")
	fs.StringVar(&opt.Namespace, "namespace", opt.Namespace, "Namespace of the service")
	fs.StringVar(&opt.Service, "service", opt.Service, "Name of the service")
	fs.Int64Var(&opt.Timeout, "timeout", opt.Timeout, "Timeout to reach each seed server")
	fs.BoolVar(&opt.Help, "help", opt.Help, "Display usage information")

	fs.Parse(os.Args)

	if opt.Help {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fs.PrintDefaults()
		return
	}

	opt.SeedURL = make([]*url.URL, len(opt.Seed))

	for index, server := range opt.Seed {
		u, err := url.Parse(server)
		if err != nil {
			log.Fatal(err)
		}
		if len(u.Scheme) == 0 || len(u.Host) == 0 {
			log.Fatal("Invalid URL: ", server)
		}
		opt.SeedURL[index] = u
	}

	ch := make(chan string)

	for _, u := range opt.SeedURL {
		go checkSeedAsync(u, opt, ch)
	}

	for _, _ = range opt.SeedURL {
		result := <-ch
		if len(result) > 0 {
			fmt.Printf("%v\n", result)
			os.Exit(0)
		}
	}

	log.Fatal("No suitable API server found")
}
