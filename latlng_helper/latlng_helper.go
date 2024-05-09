package latlng_helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chassis/foundation/httputil"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/openlog"
)

type AWSLocationOutput struct {
	Results []struct {
		Place struct {
			Categories []string
			Country    string
			Geometry   struct {
				Point []float64
			}
			Interpolated bool
			Label        string
			Municipality string
			Region       string
			SubRegion    string
		}
		Relevance float64
	}
	Summary struct {
		DataSource      string
		FilterCountries []string
		MaxResults      int
		ResultBBox      []float64
		Text            string
	}
}

func GetLocationsFromText(address string) (float64, float64, error) {
	endpointURL := "https://places.geo.ap-southeast-1.amazonaws.com/places/v0/indexes/LocationESRI/search/text"
	// aws location places can only take in text of max length = 200
	if len([]rune(address)) > 200 {
		return -1, -1, errors.New("length of address > 200")
	}
	// payload
	payload := make(map[string]interface{})
	payload["Text"] = address
	payload["MaxResults"] = 1
	payload["FilterCountries"] = []string{"THA"}

	p_string, err := json.Marshal(payload)
	if err != nil {
		openlog.Error("Error while marshalling payload")
		return -1, -1, err
	}
	fmt.Println("Request payload: " + string(p_string))
	req, err := http.NewRequest(http.MethodPost, endpointURL, bytes.NewBuffer([]byte(p_string)))
	if err != nil {
		openlog.Error("Error creating request")
		return -1, -1, err
	}
	openlog.Debug("Request URL: " + req.URL.String())
	q := req.URL.Query()
	apiKey := archaius.GetString("aws.apiKey", "")
	q.Add("key", apiKey)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		openlog.Error("Error sending request")
		return -1, -1, err
	}
	defer resp.Body.Close()

	results := AWSLocationOutput{}
	err = json.Unmarshal(httputil.ReadBody(resp), &results)
	if err != nil {
		openlog.Error("Error while decoding response: " + err.Error())
	}

	if len(results.Results) != 0 {
		lat := results.Results[0].Place.Geometry.Point[1]
		lng := results.Results[0].Place.Geometry.Point[0]
		return lat, lng, nil
	}
	return -1, -1, err
}
