package options

import "database/sql"

type Option interface {
	Name() string
	Description() string
	ID() int
	CourseID() int
}

type option struct {
	id          int
	name        string
	courseid    int
	description string
}

const GetOptionsSQL = "SELECT `optionid`,`courseid`,`name`,`description` FROM `options` WHERE `courseid` = ? ORDER BY name ASC"

func MakeOption(rows *sql.Rows) (newOption Option, err error) {
	retOption := &option{}
	newOption = retOption
	err = rows.Scan(&retOption.id, &retOption.courseid, &retOption.name, &retOption.description)
	return
}

func (this *option) Name() string {
	return this.name
}
func (this *option) Description() string {
	return this.description
}
func (this *option) ID() int {
	return this.id
}
func (this *option) CourseID() int {
	return this.courseid
}
