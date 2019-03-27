package model

import (
	"config"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

const (
	ColInfo = "info"
	ColDate = "date"
	ColTDollRecord = "tdoll_record"
	ColEquipRecord = "equip_record" // 包括 Fairy
	ColTDollStats = "tdoll_stats"
	ColEquipStats = "equip_stats"
)

var (
	DBName = config.GlobalConfig.Mongo.DB
)

var (
	mongoSession *mgo.Session
	redisConn    *redis.Conn
)

func init() {
	log.Println("Init database...")

	session, err := getMongoSession()
	if err != nil {
		log.Println("MongoDB init error!")
		log.Panic(err)
		return
	}
	mongoSession = session

	log.Println("Database init done!")
}

func getMongoSession() (*mgo.Session, error) {
	mgosession, err := mgo.DialWithTimeout(fmt.Sprintf("mongodb://%s:%s", config.GlobalConfig.Mongo.Host, config.GlobalConfig.Mongo.Port), time.Second*15)
	if err != nil {
		log.Println("Mongodb dial error!")
		log.Panic(err)
		return nil, err
	}

	mgosession.SetMode(mgo.Monotonic, true)
	mgosession.SetPoolLimit(300) // why 300 >_<

	return mgosession, nil
}
