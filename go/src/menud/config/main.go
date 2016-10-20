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

const bindSection = "bind"
const bindAddr = "addr"
const bindPort = "port"

const httpSection = "http"
const externalHost = "host"
const prefixPath = "path"
const cookieName = "cookie"
const secure = "secure"

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

func getBoolWithDefault(section, key string, defValue bool) (value bool) {
	value, ok := iniFile.GetBool(section, key)
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


func BindString() string {
	return fmt.Sprintf("%s:%d", getWithDefault(bindSection, bindAddr, "127.0.0.1"), getIntWithDefault(bindSection, bindAddr, 8000))
}

func CookieName() string {
	return getWithDefault(httpSection, cookieName, "sessid")
}

func ExternalHost() string {
	return getWithDefault(httpSection, externalHost, "")
}

func PathPrefix() string {
	return getWithDefault(httpSection, prefixPath, "/")
}

func UseHTTPS() bool {
	return getBoolWithDefault(httpSection, secure, false)
}
