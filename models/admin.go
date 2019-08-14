package models

//ExistAdmin check existence with pid in table `admin`
func ExistAdmin(pid string,token string) (exist bool, err error) {
	var count int
	sql := "select count(*) from `admin` where pid=? and token = ?"
	rows, err := db.Query(sql, pid,token)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&count)
	}
	if count > 0 {
		exist = true
	}
	return
}
