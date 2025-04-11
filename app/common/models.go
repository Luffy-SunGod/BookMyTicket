package common

import "time"

type Show struct {
	ShowName  string    `json:"show_name"`
	VenueID   int       `json:"venue_id"`
	HallID    int       `json:"hall_id"`
	Starttime time.Time `json:"show_start_time"`
	Endtime   time.Time `json:"show_end_time"`
	HallNo    int       `json:"hall_no"`
}

type ClaimSeatForm struct {
	SeatIDs    []string `json:"seat_ids"`
	ShowID     int      `json:"show_id"`
	BookedbyID int      `json:"user_id"` //user who is claiming
}

type UserSignin struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignup struct {
	UserName string `json:"name"`
	Password string `json:"password"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
	Number   string `json:"number"`
}
