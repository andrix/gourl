package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "crypto/tls"
    "strings"
)

func main() {
    var headers = flag.Bool("I", false, "Print headers information")
    var insecure = flag.Bool("k", false, "Allow to perform insecure SSL connections")
    var user_agent = flag.String("A", "gourl/1.0", "Set the User-Agent")

    flag.Parse()
    if len(flag.Args()) <= 0 {
        fmt.Printf("Usage: gourl [options] url\n")
    } else {
        url := flag.Arg(0)
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            fmt.Printf("Error creating request: %s\n", err)
            return
        }
        req.Header.Add("User-Agent", *user_agent)

        // Create TLS config
        tlsConfig := tls.Config{RootCAs: nil}

        // If insecure, skip CA verfication
        if *insecure {
            tlsConfig.InsecureSkipVerify = true
        }

        tr := &http.Transport{
            TLSClientConfig:    &tlsConfig,
            DisableCompression: true,
        }
        client := &http.Client{Transport: tr}
        resp, err := client.Do(req)

        if err != nil {
            fmt.Printf("An error ocurred: %s\n", err)
            return
        }
        if *headers {
            fmt.Printf("%s %s\n", resp.Proto, resp.Status)
            for k, vals := range resp.Header {
                fmt.Printf("%s: %s\n", k, strings.Join(vals, ";"))
            }
            fmt.Println()
        } else {
            defer resp.Body.Close()
            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                fmt.Printf("An error ocurred: %s\n", err)
                return
            }
            fmt.Printf("%s", body)
        }
    }
}
