package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Getter gets a string from a remote location
type Getter interface {
	Get(context.Context) (string, error)
}

// IPGetter retrieves the external IP address for a host as a string
type IPGetter struct {
	Getter
}

// Get will use canihazip to retrieve teh current external IP address
func (i *IPGetter) Get(ctx context.Context) (string, error) {
	r, err := http.Get("https://www.canihazip.com/s")
	if err != nil {
		return "", fmt.Errorf("Couldn't retrieve external IP address")
	}

	ip, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("Could not read IP address from response")
	}

	return fmt.Sprintf("%s", ip), nil
}
