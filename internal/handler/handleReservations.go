package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/mthuberty/movie-spots-api/internal/config"
	"github.com/mthuberty/movie-spots-api/internal/db"
	"github.com/mthuberty/movie-spots-api/internal/errs"
)

func HandleGetReservations(ac *config.AppConfiguration) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		params := r.URL.Query()

		var reservations []db.SeatReservation
		var err error
		reservations, err = ac.DBClient.GetSeatReservations(params)
		if err != nil {
			return errs.WrapTrace("handler", "HandleGetReservations", err)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(reservations)
		if err != nil {
			return errs.WrapTrace("handler", "HandleGetReservations", err)
		}
		return nil
	}
}

func HandlePostReservations(ac *config.AppConfiguration) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			return errs.WrapTrace("handler", "HandlePostReservations", err)
		}

		var reservations []db.SeatReservation

		err = json.Unmarshal(body, &reservations)
		if err != nil {
			return errs.WrapTrace("handler", "HandlePostReservations", err)
		}

		for _, r := range reservations {
			filled, err := ac.DBClient.CheckSpotIsFilled(r.ShowingID, r.SeatRow, r.SeatNumber)
			if err != nil {
				return errs.WrapTrace("handler", "HandlePostReservations", err)
			}
			if filled {
				return errs.WrapTrace("handler", "HandlePostReservations", errors.New("this spot already has a reservation"))
			}
		}

		savedReservations, err := ac.DBClient.SaveSeatReservations(reservations)
		if err != nil {
			return errs.WrapTrace("handler", "HandlePostReservations", err)
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(savedReservations)
		if err != nil {
			return errs.WrapTrace("handler", "HandlePostReservations", err)
		}
		return nil
	}
}
