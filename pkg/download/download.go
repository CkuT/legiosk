package download

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"

	"time"

	"github.com/ckut/legiosk/pkg/login"
	"github.com/ckut/legiosk/pkg/magazine"
	"github.com/h2non/filetype"
)

var ConnectionError error

type MagazineResult struct {
	Success bool              `json:"success"`
	Result  magazine.Magazine `json:"result"`
}

const defaultUA string = "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:65.0) Gecko/20100101 Firefox/65.0"

var kioskAPIURLRegex = regexp.MustCompile(`https://pros.lekiosk.com/fr/reader/(?P<publication>\d+)/(?P<issue>\d+)`)
var kioskAPIURLActivateRegex = regexp.MustCompile(`https://apipros.lekiosk.com/publications/(?P<publication>\d+)/issue/(?P<issue>\d+)/activate`) // GET Request
var kioskAPIURLDeactivateRegex = regexp.MustCompile(`https://apipros.lekiosk.com/users/library/action/delete`)                                   // PUT Request. Body Param: issue (--data '[21448847]')
var kioskPageProductURL = regexp.MustCompile(`https://pros.lekiosk.com/fr/pageproduct/(?P<publication>\d+)/(?P<issue>\d+)`)

var account *login.Account

var client *http.Client = &http.Client{
	Timeout: time.Second * 10,
}

func Download(url string, username string) error {
	// Decide what to do
	matchedAPI, _ := regexp.MatchString(kioskAPIURLRegex.String(), url)
	switch {
	case matchedAPI:
		err := download(url, username)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("Unknown URL type")
	}
	return nil
}

// ExtractIssuePublication tries to get the issue and publication for an url
func ExtractIssuePublication(url string) (*magazine.KioskMagazineInfos, error) {
	result := make(map[string]string)

	if t := kioskPageProductURL.FindStringSubmatch(url); len(t) > 0 {
		for i, name := range kioskPageProductURL.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = t[i]
			}
		}
	}
	if t := kioskAPIURLRegex.FindStringSubmatch(url); len(t) > 0 {
		for i, name := range kioskAPIURLRegex.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = t[i]
			}
		}
	}

	if len(result) != 2 {
		return nil, fmt.Errorf("Cannot extract publication and issue")
	}

	return &magazine.KioskMagazineInfos{
		Issue:       result["issue"],
		Publication: result["publication"],
	}, nil
}

func GetMag(url string, mi *magazine.KioskMagazineInfos) (*magazine.Magazine, error) {
	magURL := mi.GetMagazineURL()
	referer := mi.GetReferer()

	req, err := http.NewRequest("GET", magURL, nil)

	req.Header.Set("User-Agent", defaultUA)
	req.Header.Set("Referer", referer)
	req.Header.Set("Authorization", "JWT "+account.JWT)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var mr MagazineResult
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(content, &mr)
	if !mr.Success {
		return nil, fmt.Errorf("Cannot unmarshal magazine")
	}

	pages := []*magazine.Page{}
	for pageNumber, page := range mr.Result.RawPages.(map[string]interface{}) {
		tmp := page.(map[string]interface{})
		p, _ := strconv.Atoi(pageNumber)
		page := &magazine.Page{
			PageNumber: p,
			PdfURL:     tmp["pdfUrl"].(string),
			HiresURL:   tmp["hiresUrl"].(string),
			BigURL:     tmp["bigUrl"].(string),
		}
		pages = append(pages, page)
	}
	mr.Result.Pages = pages
	mr.Result.URL = url
	return &mr.Result, nil
}

func download(url string, username string) error {
	ac, err := login.LoginToKiosk(username)
	if err != nil {
		return err
	}
	account = ac
	client.Jar = account.Jar

	mi, err := ExtractIssuePublication(url)
	if err != nil {
		return err
	}
	mag, err := GetMag(url, mi)
	if err != nil {
		return err
	}
	downloadMag(mag)
	return nil
}

func downloadPage(page *magazine.Page, outputDir, referer string, wg *sync.WaitGroup) error {
	defer wg.Done()
	req, err := http.NewRequest("GET", page.HiresURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Referer", referer)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	bodyContent, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Guessing mimetype
	kind, _ := filetype.Match(bodyContent)
	if kind == filetype.Unknown {
		return fmt.Errorf("Unknown file type")
	}

	tmpfile, err := ioutil.TempFile(outputDir, "*."+kind.Extension)
	if err != nil {
		log.Fatal(err)
	}
	tmpfile.Write(bodyContent)

	page.TmpFile = tmpfile.Name()
	return nil
}

func downloadMag(mag *magazine.Magazine) error {
	log.Print("Downloading magazine")

	ch := make(chan *magazine.Page, 4)

	dir, err := ioutil.TempDir(os.TempDir(), "giosk")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	referer := mag.KioskInfos.GetReferer()
	var wg sync.WaitGroup
	for _, page := range mag.Pages {
		ch <- page
		go func(page *magazine.Page) {
			wg.Add(1)
			downloadPage(page, dir, referer, &wg)
			<-ch
		}(page)
	}
	wg.Wait()
	magazine.CreatePDF(mag)
	close(ch)

	return nil
}
