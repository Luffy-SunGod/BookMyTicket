package database

import (
	"fmt"
	"github.com/Luffy-SunGod/MyShowSeat/cmd/app/common"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func CheckValidValues(db *sqlx.DB, show common.Show) error {
	//validate the venueid and hallid

	//check if venue exist
	var venueExist bool
	err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM venue WHERE venueid=$1)`, show.VenueID).Scan(&venueExist)
	if err != nil {
		return fmt.Errorf("venueExsits Error: %v, sommething wrong went while executing query", err)
	}

	if !venueExist {
		return fmt.Errorf("no venue exist with %d", show.VenueID)
	}

	var hallExist bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM hall WHERE venueid=$1 and hallno=$2)`, show.VenueID, show.HallNo).Scan(&hallExist)

	if err != nil {
		return fmt.Errorf("hallexist Error: %v, sommething wrong went while executing query", err)
	}
	if !hallExist {
		return fmt.Errorf("no hall exist with %d", show.VenueID)
	}
	return nil
}

func GetHallCapacity(db *sqlx.DB, show common.Show) (string, error) {
	var capacity int
	err := db.QueryRow("SELECT capacity from hall WHERE venueid=$1 and hallno=$2", show.VenueID, show.HallNo).Scan(&capacity)
	if err != nil {
		return "", fmt.Errorf("something went wrong while executing query %v", err)
	}
	fmt.Println(capacity)
	return strconv.Itoa(capacity), nil

}

func GetSeatIds(db *sqlx.DB, show common.Show, seatIds *[]string) error {
	fmt.Println("Inside GetSeatIds")
	rows, err := db.Query(`SELECT seatId from seat where hallid=$1 and venueid=$2`, show.HallID, show.VenueID)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	fmt.Println("Trying to print the rows!!")

	if err != nil {
		return err
	}
	if !rows.Next() {
		fmt.Println("No rows found")
	}
	for rows.Next() {
		var seatid string
		if err := rows.Scan(&seatid); err != nil {
			return err
		}
		fmt.Println("seatid", seatid)
		*seatIds = append(*seatIds, seatid)
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	fmt.Println("exiting GetSeatIds")
	return nil
}

func InsertReservations(db *sqlx.DB, showid int, seatIDs []string) error {
	//begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // rollback transaction if not commited

	for _, seatid := range seatIDs {
		seatReservationId := "SH_" + strconv.Itoa(showid) + "_ST_" + seatid
		_, err := tx.Exec(`INSERT INTO reservation (seatReservationid,showid) VALUES ($1,$2)`, seatReservationId, showid)
		if err != nil {
			return err
		}
	}
	//commit the transaction if all inserts were succesfull
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func InsertShow(db *sqlx.DB, show common.Show, hallCapacity int) (int, error) {
	var showid int
	err := db.QueryRow(`INSERT INTO Show (ShowName, VenueID, HallID, Time_start, Time_end, totalcapacity, currentusage)
                    VALUES ($1, $2, $3, $4, $5, $6, $7)
                    ON CONFLICT (ShowName, VenueID, HallID, Time_start) 
                    DO UPDATE SET Time_end = EXCLUDED.Time_end
                    RETURNING ShowID`,
		show.ShowName, show.VenueID, show.HallID, show.Starttime, show.Endtime, hallCapacity, 0).Scan(&showid)

	if err != nil {
		fmt.Println("Database error:", err)
	} else {
		fmt.Println("Inserted/Updated Show ID:", showid)
	}
	return showid, err
}
