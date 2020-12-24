package database

import (
	"database/sql"
)

type Club struct {
	ID           int
	About        string
	NickOwner    string
	IDOwner      string
	ClubName     string
	NickCreator  string
	CreationDate string
	Approved     bool
	Size         int
	Slack        string
	Member       bool
	Owner        bool
}

type ClubMember struct {
	ID       int
	NickName string
	ClubName string
	Access   int
	JoinDate string
	IDIntra  string
}

func InsertNewClub(db *sql.DB, clubAbout, nickOwner, idOwner, clubName, nickCreator, creationDate, slack string) {
	_, err := db.Exec(
		"INSERT INTO clubs (clubabout, nickowner, idowner, clubname, nickcreator, creationdate, slack) values ($1, $2, $3, $4, $5, $6, $7)",
		clubAbout, nickOwner, idOwner, clubName, nickCreator, creationDate, slack)
	if err != nil {
		panic(err)
	}
}

func GetClubs(db *sql.DB, approved bool, nickName, idIntra string) ([]Club, error) {
	rows, err := getDBRows(db, "clubs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	clubs := []Club{}

	for rows.Next() {
		row := Club{}
		err := acceptRow(rows, &row)
		if err != nil {
			continue
		}
		clubsSize, _ := getClubSize(db, row.ClubName)
		row.Member = IsUserInClub(db, nickName, idIntra, row.ClubName)
		row.Size = clubsSize
		if row.Approved == approved {
			clubs = append(clubs, row)
		}
	}
	return clubs, nil
}

func GetClubsToApprove(db *sql.DB, approved bool) ([]Club, error) {
	rows, err := getDBRows(db, "clubs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	clubs := []Club{}

	for rows.Next() {
		row := Club{}
		err := acceptRow(rows, &row)
		if err != nil {
			continue
		}
		clubsSize, _ := getClubSize(db, row.ClubName)
		row.Size = clubsSize
		if row.Approved == approved {
			clubs = append(clubs, row)
		}
	}
	return clubs, nil
}

func GetOwnerClubs(db *sql.DB, ownerName string, approved bool) ([]Club, error) {
	rows, err := getDBRows(db, "clubs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	clubs := []Club{}

	for rows.Next() {
		row := Club{}
		err := acceptRow(rows, &row)
		if err != nil {
			continue
		}
		clubsSize, _ := getClubSize(db, row.ClubName)
		row.Size = clubsSize
		if row.Approved == approved && row.NickOwner == ownerName {
			clubs = append(clubs, row)
		}
	}
	return clubs, nil
}

func GetClub(db *sql.DB, clubName string, approved bool) (*Club, error) {
	row, err := db.Query("SELECT * FROM clubs WHERE clubname=$1 and approved=$2", clubName, approved)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	club := Club{}

	for row.Next() {
		err = acceptRow(row, &club)
	}
	clubsSize, _ := getClubSize(db, club.ClubName)
	club.Size = clubsSize

	return &club, nil
}

func GetUserClubs(db *sql.DB, approved bool, nickName, idIntra string) ([]Club, error) {
	rows, err := db.Query("SELECT * FROM clubs WHERE approved=$1", approved)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	clubs := []Club{}

	for rows.Next() {
		row := Club{}
		err := acceptRow(rows, &row)
		if err != nil {
			continue
		}

		clubsSize, _ := getClubSize(db, row.ClubName)
		userJoined := IsUserInClub(db, nickName, idIntra, row.ClubName)
		row.Owner = IsUserClubOwner(db, nickName, idIntra, row.ClubName)
		row.Size = clubsSize
		if row.Approved == approved && userJoined == true {
			clubs = append(clubs, row)
		}
	}
	return clubs, nil
}

func UserJoinСlub(db *sql.DB, nickName, clubName string, memberAccess int, joinDate, idIntra string) {
	clubMember := IsUserInClub(db, nickName, idIntra, clubName)

	if clubMember == true {
		return
	}

	_, err := db.Exec(
		"INSERT INTO clubmembers (nickname, clubname, memberaccess, joindate, idintra) values ($1, $2, $3, $4, $5)",
		nickName, clubName, memberAccess, joinDate, idIntra)
	if err != nil {
		panic(err)
	}
}

func UserLeaveСlub(db *sql.DB, nickName, idIntra, clubName string) {
	clubMember := IsUserInClub(db, nickName, idIntra, clubName)

	if clubMember == false {
		return
	}

	_, err := db.Exec(
		"DELETE FROM clubmembers WHERE nickname=$1 and idintra=$2 and clubname=$3",
		nickName, idIntra, clubName)
	if err != nil {
		panic(err)
	}
}

func getClubSize(db *sql.DB, clubName string) (int, error) {
	rows, err := db.Query("SELECT * FROM clubmembers WHERE clubname=$1", clubName)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	counter := 0
	for rows.Next() {
		counter++
	}

	return counter, nil
}

func IsUserClubMember(db *sql.DB, nickName, idIntra string) bool {
	err := db.QueryRow("SELECT nickname FROM clubmembers WHERE nickname=$1 and idintra=$2",
		nickName, idIntra).Scan(&nickName)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return false
	}
	return true
}

func IsUserInClub(db *sql.DB, nickName, idIntra, clubName string) bool {
	err := db.QueryRow("SELECT nickname FROM clubmembers WHERE nickname=$1 and idintra=$2 and clubname=$3",
		nickName, idIntra, clubName).Scan(&nickName)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return false
	}
	return true
}

func IsUserClubOwner(db *sql.DB, nickName, idIntra, clubName string) bool {
	var access int
	err := db.QueryRow("SELECT memberaccess FROM clubmembers WHERE nickname=$1 and idintra=$2 and clubname=$3",
		nickName, idIntra, clubName).Scan(&access)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return false
	}
	if access != 3 {
		return false
	}
	return true
}

func DeleteСlub(db *sql.DB, clubName string) {
	_, err := db.Exec("DELETE FROM clubs WHERE clubname=$1", clubName)
	if err != nil {
		panic(err)
	}
}

func DeleteClubByOwner(db *sql.DB, nickName, idIntra, clubName string) {
	_, err := db.Exec("DELETE FROM clubs WHERE nickowner=$1 and idowner=$2 and clubname=$3",
		nickName, idIntra, clubName)
	if err != nil {
		panic(err)
	}
}

func DeleteUsersFromClub(db *sql.DB, clubName string) {
	_, err := db.Exec("DELETE FROM clubmembers WHERE clubname=$1", clubName)
	if err != nil {
		panic(err)
	}
}

func AppproveClub(db *sql.DB, clubName string) {
	_, err := db.Exec("UPDATE clubs SET approved=true WHERE clubname=$1", clubName)
	if err != nil {
		panic(err)
	}
}

func SetAboutClub(db *sql.DB, newAbout, nickOwner, idOwner, clubName string) {
	_, err := db.Exec("UPDATE clubs SET clubabout=$1 WHERE nickowner=$2 and idowner=$3 and clubname=$4",
		newAbout, nickOwner, idOwner, clubName)
	if err != nil {
		panic(err)
	}
}

func SetAccess(db *sql.DB, nickName, clubName string, access int) error {
	_, err := db.Exec("UPDATE clubmembers SET memberaccess=$1 WHERE nickname=$2 and clubname=$3",
		access, nickName, clubName)

	return err
}

func SetClubOwner(db *sql.DB, newOwner, clubName string) error {
	_, err := db.Exec("UPDATE clubs SET nickowner=$1 WHERE clubname=$2",
		newOwner, clubName)

	return err
}

func GetClubMembers(db *sql.DB, clubName string) ([]ClubMember, error) {
	rows, err := db.Query("SELECT * FROM clubmembers WHERE clubname=$1", clubName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []ClubMember{}

	for rows.Next() {
		Row := ClubMember{}
		err := rows.Scan(
			&Row.ID,
			&Row.NickName,
			&Row.ClubName,
			&Row.Access,
			&Row.JoinDate,
			&Row.IDIntra)

		if err != nil {
			continue
		}

		members = append(members, Row)
	}
	return members, nil
}

func acceptRow(rows *sql.Rows, Row *Club) error {
	err := rows.Scan(
		&Row.ID,
		&Row.About,
		&Row.NickOwner,
		&Row.IDOwner,
		&Row.ClubName,
		&Row.NickCreator,
		&Row.CreationDate,
		&Row.Approved,
		&Row.Slack)

	return err
}