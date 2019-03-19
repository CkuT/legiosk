package magazine

import (
	"fmt"
)

// Magazine holds the informations about a Magazine
type Magazine struct {
	KioskInfos     KioskMagazineInfos
	Title          string      `json:"title"`
	ReleaseDate    string      `json:"releaseDate"`
	TotalPages     int         `json:"pageCount"`
	CoverURL       string      `json:"coverUrl"`
	RawPages       interface{} `json:"signedUrls"`
	Pages          []*Page
	URL            string `json:"url"`
	HasHighResURLs bool   `json:"hasHighResUrls"`
	IsPurchased    bool   `json:"isPurchased"`
	Publication    string
	Issue          string
	RawURL         string
}

// Page represents a magazine page
type Page struct {
	PageNumber int
	PdfURL     string `json:"pdfUrl"`
	BigURL     string `json:"bigUrl"`
	HiresURL   string `json:"hiresUrl"`
	TmpFile    string
}

// KioskMagazineInfos represents infos from LeKiosk website
type KioskMagazineInfos struct {
	Issue       string
	Publication string
}

func (k *KioskMagazineInfos) GetReferer() string {
	return fmt.Sprintf("https://pros.lekiosk.com/fr/reader/%s/%s", k.Publication, k.Issue)
}

func (k *KioskMagazineInfos) GetMagazineURL() string {
	return fmt.Sprintf("https://apipros.lekiosk.com/publications/%s/issues/%s/signedurls", k.Publication, k.Issue)
}

func (m *Magazine) formatMagazineTitle() string {
	return m.ReleaseDate + m.Title + ".pdf"
}
