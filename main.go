package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	//auth "bitbucket.org/valopaaltd/vp-api-auth"
	auth "github.com/Shemeikka/gauth"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Return URL-encoded string
func urlEncoded(str string) string {
	u, err := url.Parse(str)
	if err != nil {
		return ""
	}
	return u.String()
}

// Print HTTP request response
func printResponse(resp *http.Response) {
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response Body:", string(body))
	defer resp.Body.Close()
}

// Create HTTP request
func request(u string, m string, msg []byte) (*http.Request, error) {
	if strings.ToUpper(m) == "GET" {
		msg = []byte("")
	}

	req, err := http.NewRequest(m, u, bytes.NewBuffer(msg))
	if err != nil {
		return req, err
	}
	req.Header.Set("content-type", "application/json")
	return req, nil
}

// Send HTTP request
func sendRequest(req *http.Request, route string) {
	tr := &http.Transport{
		// Don't verify TLS certificate
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	printResponse(resp)
}

type ConfigRequest struct {
	Route  string
	Method string
	Body   string
}

type ConfigApi struct {
	ID  string
	Key string
}

func main() {
	configFile := kingpin.Flag("config", "Set configuration filename without extension.").Short('c').String()
	// Parse cmd arguments
	kingpin.Parse()

	// If config-file is given as a command line argument, use it, otherwise use default
	if *configFile != "" {
		fmt.Printf("Config file: %s\n", *configFile)
		viper.SetConfigName(*configFile)
	} else {
		fmt.Println("Config file: config")
		viper.SetConfigName("config")
	}

	// Read config file from current dir
	viper.AddConfigPath("")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalln("No configuration file loaded")
	}

	// Read server information from config
	serverInfo := viper.GetStringMap("server")

	var clientRequest ConfigRequest
	err = viper.MarshalKey("request", &clientRequest)
	if err != nil {
		log.Fatalln("Couldn't load config request: " + err.Error())
	}

	// Get API information from config
	var clientAPI map[string]ConfigApi

	err = viper.MarshalKey("api", &clientAPI)
	if err != nil {
		log.Fatalln("Couldn't load config request: " + err.Error())
	}

	// Init authentication object
	a := auth.Config{
		RootURL:     serverInfo["root_url"].(string),
		PrivateKey:  clientAPI["master"].Key,
		PublicID:    clientAPI["master"].ID,
		SignHeaders: []string{"date", "content-type", "content-md5"},
	}

	// Destination url
	url := serverInfo["url"].(string) + serverInfo["root_url"].(string) + "/" + clientRequest.Route

	fmt.Printf("Sending\n%s\n[%s]\nto %s\n\n", clientRequest.Method, clientRequest.Body, url)
	fmt.Printf("URL: %s\n", urlEncoded(url))

	req, err := request(urlEncoded(url), clientRequest.Method, []byte(clientRequest.Body))
	if err != nil {
		log.Fatal(err)
	}

	a.AddAuth(req)
	sendRequest(req, clientRequest.Route)
}
