package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/olafstar/salejobs-api/internal/types"
)

func (s *APIServer) handleSpecificAdvertisement(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.getSpecificAdvertisement(w, r)
	}

	return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Method not allowed"}
}

func (s *APIServer) handleAdvertisements(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.getAdvertisements(w, r)
	}
	if r.Method == "POST"{
		return s.createAdvertisements(w, r)
	}

	return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Method not allowed"}
}

func (s *APIServer) handleAdvertisementsCounter(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.countAdvertisements(w)
	}

	return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Method not allowed"}
}

func (s *APIServer) countAdvertisements(w http.ResponseWriter) error {
	advCounterBytes, okCounter := s.cache.read(CacheIDAdvc)

	var totalAds int64
	var err error
	if !okCounter {
		totalAds, err = s.store.CountAdvertisements()
		if err != nil {
			return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Failed to fetch total advertisement count"}
		}
		s.cache.update(CacheIDAdvc, totalAds)
	} else {
		err := json.Unmarshal(advCounterBytes, &totalAds)

		if err != nil {
			return err
		}
	}

	return WriteJSON(w, http.StatusOK, &types.AdvertismentCounterResponse{
		Total: totalAds,
	})
}

func (s *APIServer) getSpecificAdvertisement(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	cacheID := CacheID(fmt.Sprintf(string(CacheIDAdv), id))
	advCache, ok := s.cache.read(cacheID)

	var adv *types.CreateAdvertisementResponse
	var err error
	if !ok {
		advertisement, err := s.store.QueryAdvertisement(id)
		if err != nil {
			return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Failed to fetch advertisements"}
		}
		s.cache.update(cacheID, advertisement)
		adv = advertisement
	} else {
		err = json.Unmarshal(advCache, &adv)
		if err != nil {
			return err
		}
	}


	return WriteJSON(w, http.StatusOK, adv) 
}

func (s *APIServer) getAdvertisements(w http.ResponseWriter, r *http.Request) error {
	defaultParams := types.GetAdvertismentsBody{
		Page:  1,
		Limit: 10,
	}

	var params types.GetAdvertismentsBody = defaultParams
	var err error

	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	if page == "" && limit == "" {
		params = defaultParams
	}

	if page != "" {
		pageInt, err := strconv.ParseInt(page, 10, 64)
	
		if err != nil {
			return err
		}

		params = types.GetAdvertismentsBody{
			Page: pageInt,
			Limit: params.Limit,
		}
	}

	if limit != "" {
		limitInt, err := strconv.ParseInt(limit, 10, 64)
	
		if err != nil {
			return err
		}

		params = types.GetAdvertismentsBody{
			Page: params.Page,
			Limit: limitInt,
		}
	}


	if params.Page < 1 {
		params.Page = defaultParams.Page
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = defaultParams.Limit
	}

	cacheID := CacheID(fmt.Sprintf(string(CacheIDAdvBase), params.Page, params.Limit))

	advCache, ok := s.cache.read(cacheID)
	advCounterBytes, okCounter := s.cache.read(CacheIDAdvc)
	var totalAds int64
	if !okCounter {
		totalAds, err = s.store.CountAdvertisements()
		if err != nil {
			return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Failed to fetch total advertisement count"}
		}
		s.cache.update(CacheIDAdvc, totalAds)
	} else {
		err = json.Unmarshal(advCounterBytes, &totalAds)
		if err != nil {
			return err
		}
	}

	var adv []types.AdvertisementsCard
	if !ok {
		advertisements, err := s.store.QueryAdvertisementsCards(params.Page, params.Limit)
		if err != nil {
			return &HTTPError{StatusCode: http.StatusInternalServerError, Message: "Failed to fetch advertisements"}
		}
		s.cache.update(cacheID, advertisements)
		adv = advertisements
	} else {
		err = json.Unmarshal(advCache, &adv)
		if err != nil {
			return err
		}
	}

	lastPage := (totalAds + int64(params.Limit) - 1) / int64(params.Limit)

	if adv == nil {
		adv = []types.AdvertisementsCard{}
	}

	response := types.GetAdvertismentsResponse{
		CurrentPage:    int64(params.Page),
		Total:          totalAds,
		Last:           lastPage,
		Advertisements: adv,
	}

	return WriteJSON(w, http.StatusOK, response)
}

func (s *APIServer) createAdvertisements(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var body types.CreateAdvertisementBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return &HTTPError{StatusCode: http.StatusBadRequest, Message: "Invalid request body"}
	}
	defer r.Body.Close()

	if err := validateAdvertismentBody(body); err != nil {
		return err 
	}

	err := s.store.CreateAdvertisement(body)

	if err != nil {
		return err
	}

	s.cache.clearByPattern("adv_page_")
	s.cache.clear(CacheIDAdvc)

	return WriteJSON(w, http.StatusOK, body)
}