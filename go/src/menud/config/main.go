package config

import (
	"fmt"
	"github.com/dombenson/go-ini"
)

const iniFilePath = "config.ini"

const poolSection = "pool"
const poolSize = "size"

const dbSection = "db"

const userKey = "user"
const passKey = "pass"
const hostKey = "host"
const dbKey = "database"

var iniFile ini.Getter

func init() {
	var err error
	iniFile, err = ini.LoadFile(iniFilePath)
	if err != nil {
		panic(err)
	}
}

func getWithDefault(section, key, defValue string) (value string) {
	value, ok := iniFile.Get(section, key)
	if !ok {
		value = defValue
	}
	return
}

func getIntWithDefault(section, key string, defValue int) (value int) {
	value, ok := iniFile.GetInt(section, key)
	if !ok {
		value = defValue
	}
	return
}

func DBConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", getWithDefault(dbSection, userKey, "menud"), getWithDefault(dbSection, passKey, ""), getWithDefault(dbSection, hostKey, "127.0.0.1"), getWithDefault(dbSection, dbKey, "menud"))
}

func ConnectionPoolSize() int {
	return getIntWithDefault(poolSection, poolSize, 4)
}
