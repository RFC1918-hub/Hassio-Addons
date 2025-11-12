package ultimateguitar

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetTabByID - Fetches the corresponding tab on UG
func (s *Scraper) GetTabByID(tabID int64) (TabResult, error) {
	tabResult := TabResult{}

	urlString := fmt.Sprintf("%s%s?tab_id=%d&tab_access_type=private", ugAPIEndpoint, AppPaths.TAB_INFO, tabID)
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return tabResult, err
	}

	s.ConfigureHeaders(req)

	res, err := s.Client.Do(req)
	if err != nil {
		return tabResult, err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&tabResult)
	if err != nil {
		return tabResult, err
	}

	return tabResult, nil
}
