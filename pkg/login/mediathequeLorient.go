package login

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

var cookieJar, _ = cookiejar.New(nil)

var client *http.Client = &http.Client{
	Timeout: time.Second * 10,
	Jar:     cookieJar,
}

// Logging with a SSO from Lorient Mediatheque

const LorientAccessHash string = "19a5ba6c15e9a9d4b4756ccbeea7f1e7d69e2ac6"
const LorientAccessID string = "410"

func LoginToKiosk(username string) (*Account, error) {
	log.Print("Login to LeKiosk with mediatheque Lorient SSO")

	data := url.Values{}
	data.Set("id", LorientAccessID)
	data.Set("AccessHash", LorientAccessHash)
	data.Set("email", username+"@medialorient.fr")

	req, err := http.NewRequest(http.MethodPost, KioskSSOLoginURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !(res.StatusCode == http.StatusNotFound) {
		return nil, fmt.Errorf("Login failed")
	}

	jwt, err := extractJWT(cookieJar)
	if err != nil {
		return nil, err
	}
	a := &Account{
		Jar: cookieJar,
		JWT: jwt,
	}
	return a, nil
}
