package handler

import (
	"carinfo/internal/config"
	"carinfo/internal/entity"
	"carinfo/pkg/logger"
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// getCarInfoFromExternalAPI обращается к стороннему API для обогащения результатов запроса
// ссылку для стороннего API необходимо указать в /internal/config/.env в строке FOREIGN_API_URL
func getCarInfoFromExternalAPI(regNums entity.RegNums) ([]entity.Car, error) {
	logger.Log.Debug("getCarInfoFromExternalAPI starting")
	cars := make([]entity.Car, 0, len(regNums.RegNum))
	c := http.Client{}
	var err error
	for _, regNum := range regNums.RegNum {
		car := entity.NewCar()
		link, err := url.Parse(config.All.ForeignApiUrl)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return nil, err
		}
		q := link.Query()
		q.Set("regNum", regNum)
		link.RawQuery = q.Encode()
		req, err := http.NewRequest(http.MethodGet, link.String(), nil)
		logger.Log.Debug("getCarInfoFromExternalAPI request starting for regNum", regNum)
		resp, err := c.Do(req)
		if err != nil {
			return cars, err
		}
		err = json.NewDecoder(resp.Body).Decode(&car)
		if err != nil {
			return cars, err
		}
		car.RegNum = regNum
		cars = append(cars, car)
	}
	logger.Log.Debug("getCarInfoFromExternalAPI ended")
	return cars, err
}

func (h *Handler) handleCreateNewCar(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("handleCreateNewCar starting")
	switch r.Method {
	case http.MethodPost:
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.All.TimeOut)*time.Second)
		defer cancel()
		w.Header().Set("Content-Type", "application/json")
		var newCars []entity.Car
		var nums entity.RegNums
		err := json.NewDecoder(r.Body).Decode(&nums)
		if err != nil {
			ms := "JSON decoding error: "
			logger.Log.Error(ms, err.Error(), "")
			http.Error(w, ms, http.StatusInternalServerError)
		}
		cars, err := getCarInfoFromExternalAPI(nums)
		newCars, err = h.CarService.Create(ctx, cars)
		if err != nil {
			logger.Log.Error("car creating error:", err.Error(), "!")
			http.Error(w, "car creating error", http.StatusInternalServerError)
		}
		jsonCars, err := json.Marshal(newCars)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.Error("marshaling error")
		}
		_, err = fmt.Fprintf(w, string(jsonCars))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.Error("marshaling error")
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	logger.Log.Debug("handleCreateNewCar ended")
}
func (h *Handler) handleListCars(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("handleListCars starting")
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.All.TimeOut)*time.Second)
		defer cancel()
		w.Header().Set("Content-Type", "application/json")
		options := getQueries(r)
		cars, err := h.CarService.ReadAll(ctx, options)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.Error("all cars getting error", err.Error(), "!")
		}
		jsonCars, err := json.Marshal(cars)
		fmt.Fprintf(w, string(jsonCars))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	logger.Log.Debug("handleListCars ending")
}
func getQueries(r *http.Request) map[string]string {
	logger.Log.Debug("getQueries starting")
	page := r.URL.Query().Get("page")
	perPage := r.URL.Query().Get("per_page")
	mark := r.URL.Query().Get("mark")
	model := r.URL.Query().Get("model")
	year := r.URL.Query().Get("year")
	regNum := r.URL.Query().Get("reg_num")
	ownerName := r.URL.Query().Get("owner_name")
	ownerSurname := r.URL.Query().Get("owner_surname")
	ownerPatronymic := r.URL.Query().Get("owner_patronymic")
	options := make(map[string]string)
	options["page"] = page
	options["per_page"] = perPage
	options["mark"] = mark
	options["model"] = model
	options["year"] = year
	options["reg_num"] = regNum
	options["owner_name"] = ownerName
	options["owner_surname"] = ownerSurname
	options["owner_patronymic"] = ownerPatronymic
	for k, v := range options {
		if v == "" {
			delete(options, k)
		}
	}
	logger.Log.Debug("getQueries ended")
	return options
}

func (h *Handler) handleGetCarById(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("handleGetCarById starting")
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.All.TimeOut)*time.Second)
		defer cancel()
		w.Header().Set("Content-Type", "application/json")
		id := r.PathValue("id")
		idInt, err := strconv.Atoi(id)
		if err != nil || idInt <= 0 {
			logger.Log.Error("car id error")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, fmt.Sprintf("{\"car bad id\":%s}", id))
			return
		}
		car, err := h.CarService.ReadById(ctx, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "one car getting error")
			logger.Log.Error("one car getting error:", err.Error(), "!")
		}
		if car.ID == "" {
			fmt.Fprintf(w, fmt.Sprintf("{\"car not found\":%s}", id))
			return
		}
		jsonCar, err := json.Marshal(car)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Log.Error("marshaling error")
		}
		fmt.Fprintf(w, string(jsonCar))
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	logger.Log.Debug("handleGetCarById ended")
}

func (h *Handler) handleEditCarById(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("handleEditCarById starting")
	switch r.Method {
	case http.MethodPut:
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.All.TimeOut)*time.Second)
		defer cancel()
		id := r.PathValue("id")
		idInt, err := strconv.Atoi(id)
		w.Header().Set("Content-Type", "application/json")
		if err != nil || idInt <= 0 {
			logger.Log.Error("car id error")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		car := entity.NewCar()
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&car); err != nil {
			ms := "JSON decoding error: "
			logger.Log.Error(ms, err.Error(), "!")
			http.Error(w, ms, http.StatusInternalServerError)
		}
		car, err = h.CarService.Update(ctx, id, car)
		if err != nil {
			http.Error(w, "one car getting error", http.StatusInternalServerError)
			logger.Log.Error("one car getting error", err.Error(), "!")
		}
		jsonCar, err := json.Marshal(car)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		fmt.Fprintf(w, string(jsonCar))
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	logger.Log.Debug("handleEditCarById ended")
}

func (h *Handler) handleDeleteCarById(w http.ResponseWriter, r *http.Request) {
	logger.Log.Debug("handleDeleteCarById starting")
	switch r.Method {
	case http.MethodDelete:
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.All.TimeOut)*time.Second)
		defer cancel()
		w.Header().Set("Content-Type", "application/json")
		id := r.PathValue("id")
		idInt, err := strconv.Atoi(id)
		if err != nil || idInt <= 0 {
			logger.Log.Error("car id error")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = h.CarService.Delete(ctx, id)
		if err != nil {
			http.Error(w, "car deleting error", http.StatusInternalServerError)
			logger.Log.Error("car deleting error", err.Error(), "!")
		}
		fmt.Fprintf(w, "car deleted")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	logger.Log.Debug("handleDeleteCarById ended")
}
