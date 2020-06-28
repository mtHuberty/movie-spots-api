package restclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mthuberty/movie-spots-api/internal/errs"
)

type RestClient struct {
	RestClientIface
	OAuthToken string
}

type RestClientIface interface {
	Get(url string, qp map[string]string, v interface{}) error
}

func (rc RestClient) Get(url string, qp map[string]string, v interface{}) error {
	token := rc.OAuthToken
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errs.WrapTrace("restclient", "Get", err)
	}

	if qp != nil {
		q := request.URL.Query()
		for k, v := range qp {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	client := http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return errs.WrapTrace("restclient", "Get", fmt.Errorf("Failed to GET - %v", err))
	}
	defer resp.Body.Close()

	rBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errs.WrapTrace("restclient", "Get", fmt.Errorf("Failed to read the response body - %v", err))
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		err = json.Unmarshal([]byte(rBytes), v)
		if err != nil {
			return errs.WrapTrace("restclient", "Get", fmt.Errorf("Failed to unmarshall json data - %v", err))
		}

		return nil
	}

	// We should return an error instead of nil if statusCode isn't a 200 or 300 level code
	return nil
}
