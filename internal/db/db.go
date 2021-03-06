package db

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/mthuberty/movie-spots-api/internal/errs"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type DBIface interface {
	GetSeatReservations(url.Values) ([]SeatReservation, error)
	SaveSeatReservations([]SeatReservation) ([]SeatReservation, error)
	CheckSpotIsFilled(showingID string, seatRow string, seatNumber int) (bool, error)
	Close()
}

type DatabaseClient struct {
	DB *sql.DB
}

type SeatReservation struct {
	UUID       string `json:"uuid"`
	TheaterID  string `json:"theater_id"`
	LocationID string `json:"location_id"`
	ShowingID  string `json:"showing_id"`
	SeatRow    string `json:"seat_row"`
	SeatNumber int    `json:"seat_number"`
}

func NewDB(portstring, host, user, password, database string) (DBIface, error) {
	port, _ := strconv.Atoi(portstring)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)

	DB, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, errs.WrapTrace("db", "NewDB", fmt.Errorf("error attempting to connect to our database - %s", err))
	}

	if err = DB.Ping(); err != nil {
		return nil, errs.WrapTrace("db", "NewDB", fmt.Errorf("error attempting to verify database connection - %s", err))
	} else {
		log.Infoln("Database ping returned successfully - db is connected and responding.")
	}

	return &DatabaseClient{DB}, nil
}

func (dbc *DatabaseClient) Close() {
	dbc.DB.Close()
}

func (dbc *DatabaseClient) GetSeatReservations(params url.Values) ([]SeatReservation, error) {
	var seatReservations []SeatReservation

	var uuid string
	var theaterID string
	var locationID string
	var showingID string
	var seatRow string
	var seatNumber int

	sqlStatement := "SELECT * FROM seat_reservations"

	i := 0
	argCount := 1
	var args = make([]interface{}, 0)
	for key, valArr := range params {
		if i == 0 {
			sqlStatement = sqlStatement + " WHERE "
		}
		if key != "" {
			args = append(args, valArr[0])
			sqlStatement = sqlStatement + fmt.Sprintf("%s=$%d", key, argCount)
			argCount++
		}
		if i < len(params)-1 {
			sqlStatement = sqlStatement + " AND "
		}
		i++
	}

	rows, err := dbc.DB.Query(sqlStatement, args...)
	if err != nil {
		return nil, errs.WrapTrace("db", "GetSeatReservations", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&uuid, &theaterID, &locationID, &showingID, &seatRow, &seatNumber)
		if err != nil {
			return nil, errs.WrapTrace("db", "GetSeatReservations", err)
		}

		var seatReservation SeatReservation

		seatReservation.UUID = uuid
		seatReservation.TheaterID = theaterID
		seatReservation.LocationID = locationID
		seatReservation.ShowingID = showingID
		seatReservation.SeatRow = seatRow
		seatReservation.SeatNumber = seatNumber

		seatReservations = append(seatReservations, seatReservation)
	}

	return seatReservations, nil
}

func (dbc *DatabaseClient) SaveSeatReservations(seatReservations []SeatReservation) ([]SeatReservation, error) {

	for _, r := range seatReservations {
		if r.LocationID == "" {
			return nil, errs.WrapTrace("db", "SaveSeatReservations", errors.New("location_id must be provided"))
		}
		if r.TheaterID == "" {
			return nil, errs.WrapTrace("db", "SaveSeatReservations", errors.New("theater_id must be provided"))
		}
		if r.ShowingID == "" {
			return nil, errs.WrapTrace("db", "SaveSeatReservations", errors.New("showing_id must be provided"))
		}
		if r.SeatRow == "" {
			return nil, errs.WrapTrace("db", "SaveSeatReservations", errors.New("seat_row must be provided"))
		}
		if r.SeatNumber == 0 {
			return nil, errs.WrapTrace("db", "SaveSeatReservations", errors.New("seat_number must be provided"))
		}
	}

	sqlStr := `INSERT INTO seat_reservations(location_id, theater_id, showing_id, seat_row, seat_number) VALUES `

	var vals []interface{}

	count := 1

	for _, row := range seatReservations {
		sqlStr += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", count, count+1, count+2, count+3, count+4)
		vals = append(vals, row.TheaterID, row.LocationID, row.ShowingID, row.SeatRow, row.SeatNumber)
		count += 5
	}
	//trim the last , and return added rows
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	sqlStr += " RETURNING *;"

	prepSQL, err := dbc.DB.Prepare(sqlStr)
	if err != nil {
		return nil, errs.WrapTrace("db", "SaveSeatReservations", err)
	}

	rows, err := prepSQL.Query(vals...)
	if err != nil {
		return nil, errs.WrapTrace("db", "SaveSeatReservations", err)
	}
	defer rows.Close()

	var uuid string
	var theaterID string
	var locationID string
	var showingID string
	var seatRow string
	var seatNumber int

	var res []SeatReservation

	for rows.Next() {
		err := rows.Scan(&uuid, &theaterID, &locationID, &showingID, &seatRow, &seatNumber)
		if err != nil {
			return nil, errs.WrapTrace("db", "SaveSeatReservations", err)
		}

		var seatReservation SeatReservation

		seatReservation.UUID = uuid
		seatReservation.TheaterID = theaterID
		seatReservation.LocationID = locationID
		seatReservation.ShowingID = showingID
		seatReservation.SeatRow = seatRow
		seatReservation.SeatNumber = seatNumber

		res = append(res, seatReservation)
	}

	return res, nil
}

func (dbc *DatabaseClient) CheckSpotIsFilled(showingID string, seatRow string, seatNumber int) (bool, error) {

	sqlStatement := fmt.Sprintf("SELECT * FROM seat_reservations WHERE showing_id = $1 AND seat_row = $2 AND seat_number = $3")

	rows, err := dbc.DB.Query(sqlStatement, showingID, seatRow, seatNumber)
	if err != nil {
		return false, errs.WrapTrace("db", "CheckSpotIsFilled", err)
	}
	defer rows.Close()

	for rows.Next() {
		// there is already a reservation for this seat in the specified showing
		return true, nil
	}

	return false, nil
}
