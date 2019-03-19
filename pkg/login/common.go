package login

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
)

const KioskSSOLoginURL string = "https://apipros.lekiosk.com/login/ssoform"

// Account holds a Kiosk account
type Account struct {
	JWT string
	Jar *cookiejar.Jar
}

func extractJWT(jar *cookiejar.Jar) (string, error) {
	req, err := http.NewRequest("GET", KioskSSOLoginURL, nil)
	if err != nil {
		return "", err
	}

	cookie := client.Jar.Cookies(req.URL)
	for _, c := range cookie {
		if c.Name == "authToken" {
			return c.Value, nil
		}
	}
	return "", fmt.Errorf("Cannot find authToken header")
}
