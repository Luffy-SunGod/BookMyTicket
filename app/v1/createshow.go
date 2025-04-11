package v1

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Luffy-SunGod/MyShowSeat/cmd/app/common"
	"github.com/Luffy-SunGod/MyShowSeat/cmd/app/database"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"strconv"
)

func CreateShowController(ctx *fiber.Ctx) error {
	var show common.Show
	fmt.Println("createshowcontroller")

	// we can do like this as well
	//err:=ctx.App().Config().JSONDecoder(ctx.Body(),&show)
	// this  is way  better
	err := ctx.BodyParser(&show)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(fmt.Sprintf("Failed to parse show form: %v", err))
	}

	db, err := database.ConnecttoDB()
	if err != nil {
		log.Fatalf("Unable to connect to database getting this %v error", err)
		return err
	}
	fmt.Println(show)
	err = database.CheckValidValues(db, show)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(fmt.Sprintf("Check failed %v", err))
	}

	hallCapacityStr, err := database.GetHallCapacity(db, show)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(fmt.Sprintf("Capacity failed %v", err))
	}
	fmt.Println("hallCapacityStr", hallCapacityStr)
	hallCapacity, err := strconv.Atoi(hallCapacityStr)

	if err != nil {
		log.Println("cannot convert hallcapacity to integer")
	}

	// Create show in show table, with capacity=hallcapacity, usage=0

	showid, err := database.InsertShow(db, show, hallCapacity)
	fmt.Println("showid", showid)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.Status(http.StatusBadRequest).SendString(fmt.Sprintf("show already exist with show id %d", showid))
		} else {
			// Handle other errors
			log.Println("Insert error")
			return ctx.Status(http.StatusBadRequest).SendString(fmt.Sprintf("Database error %v", err))
		}
	}

	//map<hallid,capacity> so that we don't need to query again and again in order to get capacity for a particular hall
	err = common.CreateRedisEntry(showid, hallCapacity)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).SendString(fmt.Sprintf("failed to updated Redis %d", http.StatusInternalServerError))
	}

	var seatIds []string
	err = database.GetSeatIds(db, show, &seatIds)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("Unable to get seatIds!! %v,%v", http.StatusInternalServerError, err))
	}

	err = database.InsertReservations(db, showid, seatIds)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(fmt.Sprintf("can't reserve seats getting this error -%v", err))
	}
	response := map[string]int{
		"showId": showid,
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}
