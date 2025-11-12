package ultimateguitar

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"
)

// API constants
const ugAPIEndpoint = "https://api.ultimate-guitar.com/api/v1" // API Base URL
const ugUserAgent = "UGT_ANDROID/4.11.1 (Pixel; 8.1.0)"        // Default useragent for a pixel device
const UGTimeFormat = "2006-01-02"                              // API Key datetime format

// Default headers
var ugHeaders = map[string]string{
	"Accept-Charset": "utf-8",
	"Accept":         "application/json",
	"User-Agent":     ugUserAgent,
	"Connection":     "close",
}

// Scraper struct
type Scraper struct {
	Client   *http.Client
	DeviceID string
	APIKey   string
	Token    string
}

// Generates a new device id for the scraper instances. This value is used in the request headers and to generate X-UG-API-KEY.
func (s *Scraper) generateDeviceID() {
	raw := make([]byte, 16)
	_, err := rand.Read(raw)
	if err != nil {
		log.Fatal(err)
	}
	s.DeviceID = fmt.Sprintf("%x", raw)[:16]
}

// Generate the X-UG-API-KEY for this request. The payload is the MD5 result of the concatenated value of device id + "2006-01-02:15" (utc) + "createLog()"
func (s *Scraper) generateAPIKey() string {
	hour := time.Now().UTC().Hour()
	formattedDate := fmt.Sprintf("%s:%d", time.Now().UTC().Format(UGTimeFormat), hour)

	hashed := md5.Sum([]byte(fmt.Sprintf("%s%s%s", s.DeviceID, formattedDate, "createLog()")))
	return fmt.Sprintf("%x", hashed)
}

func (s *Scraper) ConfigureHeaders(req *http.Request) {
	for key := range ugHeaders {
		req.Header[key] = []string{ugHeaders[key]}
	}
	req.Header["X-UG-CLIENT-ID"] = []string{s.DeviceID}
	req.Header["X-UG-API-KEY"] = []string{s.generateAPIKey()}

	// This header isn't sent in the app, so we remove it.
	req.Header.Del("Accept-Encoding")
}

// New creates a new Scraper instance
func New() Scraper {
	s := Scraper{
		Client: &http.Client{},
	}
	s.generateDeviceID()
	s.generateAPIKey()
	return s
}
