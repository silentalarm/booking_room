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
	S3file       string
	Member       bool
	Owner        bool
	Color        string
}

type ClubMember struct {
	ID       int
	NickName string
	ClubName string
	Access   int
	JoinDate string
	IDIntra  string
}

func InsertNewClub(db *sql.DB, clubAbout, clubName, nickCreator, creationDate, slack string) {
	_, err := db.Exec(
		"INSERT INTO clubs (clubabout, clubname, nickcreator, creationdate, slack) values ($1, $2, $3, $4, $5)",
		clubAbout, clubName, nickCreator, creationDate, slack)
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
		row.NickOwner, row.IDOwner, _ = GetMemberByAccess(db, row.ClubName, 3)
		row.Member = IsUserInClub(db, nickName, idIntra, row.ClubName)
		row.Owner = IsUserClubOwner(db, nickName, idIntra, row.ClubName)
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
		row.NickOwner, row.IDOwner, _ = GetMemberByAccess(db, row.ClubName, 3)
		row.Size = clubsSize
		if row.Approved == approved {
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

	club.NickOwner, club.IDOwner, _ = GetMemberByAccess(db, club.ClubName, 3)

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
		row.NickOwner, row.IDOwner, _ = GetMemberByAccess(db, row.ClubName, 3)
		row.Member = IsUserInClub(db, nickName, idIntra, row.ClubName)
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

func DeleteAllClubReserves(db *sql.DB, clubName string) {
	_, err := db.Exec("DELETE FROM floor_2 WHERE clubname=$1", clubName)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DELETE FROM floor_3 WHERE clubname=$1", clubName)
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

func IsUserClubOwnerOrModer(db *sql.DB, nickName, idIntra, clubName string) bool {
	var access int
	err := db.QueryRow("SELECT memberaccess FROM clubmembers WHERE nickname=$1 and idintra=$2 and clubname=$3",
		nickName, idIntra, clubName).Scan(&access)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return false
	}
	if access == 3 || access == 2 {
		return true
	}
	return false
}

func ClubIsExist(db *sql.DB, clubName string) bool {
	err := db.QueryRow("SELECT clubname FROM clubs WHERE clubname=$1",
		clubName).Scan(&clubName)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return false
	}
	return true
}

func GetNameFile(db *sql.DB, clubName string) (string, error) {
	var nameClub string
	err := db.QueryRow("SELECT s3file FROM clubs WHERE clubname=$1",
		clubName).Scan(&nameClub)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return nameClub, err
	}
	return nameClub, nil
}

func GetClubColor(db *sql.DB, clubName string) string {
	color := "#039BE5"
	err := db.QueryRow("SELECT color FROM clubs WHERE clubname=$1",
		clubName).Scan(&color)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return color
	}
	return color
}

func GetMemberByAccess(db *sql.DB, clubName string, access int) (string, string, error) {
	var nameClub string
	var idIntra string
	err := db.QueryRow("SELECT nickname, idintra FROM clubmembers WHERE clubname=$1 and memberaccess=$2",
		clubName, access).Scan(&nameClub, &idIntra)
	if err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		}
		return nameClub, idIntra, err
	}
	return nameClub, idIntra, nil
}

func DeleteСlub(db *sql.DB, clubName string) {
	_, err := db.Exec("DELETE FROM clubs WHERE clubname=$1",
		clubName)
	if err != nil {
		panic(err)
	}
}

func DeleteClub(db *sql.DB, clubName string) {
	_, err := db.Exec("DELETE FROM clubs WHERE clubname=$1",
		clubName)
	if err != nil {
		panic(err)
	}
}

func DeleteUsersFromClub(db *sql.DB, clubName string) {
	_, err := db.Exec("DELETE FROM clubmembers WHERE clubname=$1",
		clubName)
	if err != nil {
		panic(err)
	}
}

func AppproveClub(db *sql.DB, clubName string) {
	_, err := db.Exec("UPDATE clubs SET approved=true WHERE clubname=$1",
		clubName)
	if err != nil {
		panic(err)
	}
}

func SetAboutClub(db *sql.DB, newAbout, clubName string) {
	_, err := db.Exec("UPDATE clubs SET clubabout=$1 WHERE clubname=$2",
		newAbout, clubName)
	if err != nil {
		panic(err)
	}
}

func SetSlackClub(db *sql.DB, slack, clubName string) {
	_, err := db.Exec("UPDATE clubs SET slack=$1 WHERE clubname=$2",
		slack, clubName)
	if err != nil {
		panic(err)
	}
}

func SetColorClub(db *sql.DB, color, clubName string) {
	_, err := db.Exec("UPDATE clubs SET color=$1 WHERE clubname=$2",
		color, clubName)
	if err != nil {
		panic(err)
	}
}

func SetAccess(db *sql.DB, nickName, clubName string, access int) error {
	_, err := db.Exec("UPDATE clubmembers SET memberaccess=$1 WHERE nickname=$2 and clubname=$3",
		access, nickName, clubName)

	return err
}

func SetImageName(db *sql.DB, clubName, fileName string) error {
	_, err := db.Exec("UPDATE clubs SET s3file=$1 WHERE clubname=$2",
		fileName, clubName)

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
		row := ClubMember{}
		err := acceptRowMember(rows, &row)

		if err != nil {
			continue
		}

		members = append(members, row)
	}
	return members, nil
}

func GetMemberClubsByAccess(db *sql.DB, nickName, intraID string, access int) ([]string, error) {
	rows, err := db.Query("SELECT * FROM clubmembers WHERE nickname=$1 and idintra=$2 and memberaccess=$3",
		nickName, intraID, access)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memberClubs := []string{}

	for rows.Next() {
		row := ClubMember{}
		err := acceptRowMember(rows, &row)
		if err != nil {
			continue
		}

		memberClubs = append(memberClubs, row.ClubName)
	}
	return memberClubs, nil
}

func acceptRow(rows *sql.Rows, Row *Club) error {
	err := rows.Scan(
		&Row.ID,
		&Row.About,
		&Row.ClubName,
		&Row.NickCreator,
		&Row.CreationDate,
		&Row.Approved,
		&Row.Slack,
		&Row.S3file,
		&Row.Color)

	return err
}

func acceptRowMember(rows *sql.Rows, Row *ClubMember) error {
	err := rows.Scan(
		&Row.ID,
		&Row.NickName,
		&Row.ClubName,
		&Row.Access,
		&Row.JoinDate,
		&Row.IDIntra)

	return err
}
