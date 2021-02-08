package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"gotest.tools/v3/assert"
)

type ItemResponseTest struct {
	Items      []Item `json:"items"`
	PageNumber int    `json:"page_number"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
	TotalItems int    `json:"total_items"`
}

func TestAPIConfig_IntegrationTests(t *testing.T) {
	apiConfig := NewAPIConfig(GenerateSampleItems())
	s := httptest.NewServer(apiConfig.Init())
	defer s.Close()

	t.Run("returns 404 for invalid page", func(t *testing.T) {
		r, err := http.DefaultClient.Get(s.URL + "/items?pageSize=30&page=30")
		assert.NilError(t, err)
		assert.Equal(t, r.StatusCode, 404)
	})

	t.Run("returns 404 for invalid category", func(t *testing.T) {
		r, err := http.DefaultClient.Get(s.URL + "/items?pageSize=30&page=1&category=food")
		assert.NilError(t, err)
		assert.Equal(t, r.StatusCode, 404)
	})

	t.Run("returns all items when no category provided", func(t *testing.T) {
		r, err := http.DefaultClient.Get(s.URL + "/items?pageSize=15&page=1")
		assert.NilError(t, err)
		assert.Equal(t, r.StatusCode, 200)
		response := ItemResponseTest{}
		json.NewDecoder(r.Body).Decode(&response)
		assert.Equal(t, response.TotalItems, 15)
		assert.Equal(t, response.PageNumber, 1)
		assert.Equal(t, response.PageSize, 15)
		assert.Equal(t, response.TotalPages, 1)
	})

	t.Run("filters by category", func(t *testing.T) {
		r, err := http.DefaultClient.Get(s.URL + "/items?pageSize=15&page=1&category=Fruit")
		assert.NilError(t, err)
		assert.Equal(t, r.StatusCode, 200)
		response := ItemResponseTest{}
		json.NewDecoder(r.Body).Decode(&response)
		assert.Equal(t, response.TotalItems, 4)
		assert.Equal(t, response.PageNumber, 1)
		assert.Equal(t, response.PageSize, 15)
		assert.Equal(t, response.TotalPages, 1)
		for _, item := range response.Items {
			assert.Equal(t, item.Category, "Fruit")
		}
	})

	t.Run("multiple pages", func(t *testing.T) {
		for i := 1; i <= 4; i++ {

			v := url.Values{}
			v.Add("pageSize", "4")
			v.Add("page", strconv.Itoa(i))
			u, _ := url.Parse(s.URL + "/items")
			u.RawQuery = v.Encode()

			r, err := http.DefaultClient.Get(u.String())
			assert.NilError(t, err)
			assert.Equal(t, r.StatusCode, 200)
			response := ItemResponseTest{}
			json.NewDecoder(r.Body).Decode(&response)
			assert.Equal(t, response.TotalItems, 15)
			assert.Equal(t, response.PageNumber, i)
			assert.Equal(t, response.PageSize, 4)
			assert.Equal(t, response.TotalPages, 4)
		}
	})
}
