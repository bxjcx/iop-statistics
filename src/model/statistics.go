package model

import (
	"config"
	"fmt"
	"log"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

const (
	RecordTypeTDoll = iota
	RecordTypeEquip
	RecordTypeFairy
)

// Record 单条记录
type Record struct {
	Formula   Formula `json:"formula" bson:"formula"`       // 公式
	BuildSlot int     `json:"build_slot" json:"build_slot"` // 建造栏位
	DevLevel  int     `json:"dev_lv"`                       // 建造等级(玩家)
	ID        int     `json:"id" bson:"id"`                 // 建造结果ID
	DevType   int     `json:"dev_type" bson:"dev_type"`     // 建造类型
	DevTime   int     `json:"dev_time" bson:"dev_time"`     // 记录时间
}

// StatsRecord 统计记录，在_stats表中
type StatsRecord struct {
	Formula Formula `json:"formula" bson:"formula"`
	ID      int     `json:"id" bson:"id"`
	Type    int     `json:"type" bson:"type"`
	Count   int     `json:"count" bson:"count"`
	Date    string     `json:"date" bson:"date"` // TODO: 为什么会是一个string
}

// 公式
type Formula struct {
	Mp         int `bson:"mp" json:"mp"`
	Ammo       int `bson:"ammo" json:"ammo"`
	Mre        int `bson:"mre" json:"mre"`
	Part       int `bson:"part" json:"part"`
	InputLevel int `bson:"input_level" json:"input_level"`
}

// 信息
type Info struct {
	LastUpdate int    `json:"last_update"`
	GunCount   int    `json:"gun_count"`
	EquipCount int    `json:"equip_count"`
	Info       string `json:"info"`
}

// Stats 统计结果
type Stats struct {
	Formula *Formula `json:"formula,omitempty" bson:"formula"`
	ID      int      `json:"id,omitempty" bson:"formula"`
	Type    int      `json:"type" bson:"type"`
	Count   int      `json:"count" bson:"count"`
	Total   int      `json:"total,omitempty" bson:"total"`
}

// TimeRecord 新加枪的加入时间
type TimeRecord struct {
	Type int `json:"time" bson:"time"`
	ID   int `json:"id" bson:"id"`
	Date int `json:"date" bson:"date"`
}

// KVPair 键值对
type KVPair struct {
	Key   string      `bson:"key"`
	Value interface{} `bson:"value"`
}

func GetInfo() (Info, error) {
	ret := Info{}

	var err error

	ret.LastUpdate, err = getInfoLastUpdate()
	if err != nil {
		return ret, err
	}

	ret.Info, err = getInfoInfo()
	if err != nil {
		return ret, err
	}

	ret.GunCount, err = getCollectionCount(ColTDollRecord)
	if err != nil {
		return ret, err
	}

	ret.EquipCount, err = getCollectionCount(ColEquipRecord)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

// GetIDStatsData 获得由ID为查询条件的记录，返回值为记录
func GetIDStatsData(t int, id int, from int, to int) ([]*Stats, error) {
	s := mongoSession.Copy()
	defer s.Close()
	var col string

	// 设置对应表名
	switch t {
	case RecordTypeTDoll:
		col = ColTDollStats
	case RecordTypeEquip:
		fallthrough
	case RecordTypeFairy:
		col = ColEquipStats
	default:
		panic("Invalid type at GetIDStatsData")
	}

	c := s.DB(DBName).C(col)


	query := bson.D{ // 查询条件为ID
		//{"date", to},
		{"id", id},
	}

	if t != RecordTypeTDoll {
		query = append(query, bson.DocElem{Name:"type", Value:t}) // 如果建造类型不是人形，需要加入type
	}

	boundary, err := getBoundaryDate(col, query, to, 2) // 获得结束时间的boundary
	if err != nil {
		return nil, err
	}

	query = append(query, bson.DocElem{Name:"date", Value:strconv.Itoa(boundary)}) // TODO: 暂时换成string，因为数据库错了

	cur := c.Find(query).Sort("-count")
	it := cur.Iter()

	log.Printf("Query 1: %+v, Boundary: %d\n", query, boundary)

	var i StatsRecord
	var ret []*Stats
	counter := 0
	for it.Next(&i) { // 遍历每一条记录
		if i.Count <= config.GlobalConfig.Statistics.RecordThreshold &&
			counter > config.GlobalConfig.Statistics.MinRecordCount {
			break
		}


		minus := getSingleStatsCount(col, query, from, 1) // 获得开始时间的建造数量，后面要减去
		log.Printf("Query 2: %+v, From: %d, Minus: %d\n", i, from, minus)

		query = bson.D{
			{"formula", i.Formula},
		}

		total := getStatsCount(col, query, from, to) // 获得from到to内的单公式出货数
		log.Printf("Query 3: Formula:%+v, Total:%d", i.Formula, total)
		// 获得from_time的值

		s := new(Stats)
		s.Formula = &i.Formula
		// s.ID = i.ID
		s.Type = i.Type
		s.Count = i.Count - minus
		s.Total = total

		//s = Stats{
		//	Formula: &i.Formula,
		//	//ID:      i.ID,
		//	Type:  i.Type,
		//	Count: i.Count - minus,
		//	Total: total,
		//}

		log.Printf("Result:%+v\n", s)
		ret = append(ret, s)
		counter++
	}

	return ret, nil

}

// TODO: 是不是可以加个缓存？
// GetFormulaStatsData 获得由公式为查询条件的记录，返回值为 记录
func GetFormulaStatsData(t int, formula Formula, from int, to int) ([]*Stats, error) {
	s := mongoSession.Copy()
	defer s.Close()
	var col string

	switch t {
	case RecordTypeTDoll:
		col = ColTDollStats
	case RecordTypeEquip:
		fallthrough
	case RecordTypeFairy:
		col = ColEquipStats
	default:
		panic("Invalid type at GetFormulaStatsData")
	}

	c := s.DB(DBName).C(col)

	var ret []*Stats
	var recs []StatsRecord

	boundary, err := getBoundaryDate(col, bson.D{{"formula", formula},}, to, 2)
	if err != nil {
		return ret, err
	}

	query := bson.D{
		{"date", boundary}, // TODO: 边界判断
		{"formula", formula},
	}

	err = c.Find(query).All(&recs)
	if err != nil {
		return ret, err
	}

	m := make(map[int]*Stats)
	for _, v := range recs {
		s := &Stats{
			ID:    v.ID,
			Count: v.Count,
		}

		m[s.ID] = s
		ret = append(ret, s)

	}

	query[0].Value = from // 第一个是 date

	recs = []StatsRecord{} // 防止出事儿，清空掉
	err = c.Find(query).All(&recs)
	if err != nil {
		return ret, err
	}

	for _, v := range recs {
		s, ok := m[v.ID]
		if ok {
			s.Count -= v.Count
		}
	}

	return ret, nil
}

// getStatsCount 前减后的方式获得统计数目
func getStatsCount(col string, query bson.D, from int, to int) int {

	time1 := getSingleStatsCount(col, query, from, 1)
	time2 := getSingleStatsCount(col, query, to, 2)

	return time2 - time1
}

func getSingleStatsCount(col string, query bson.D, date int, lr int) int {
	if date == 0 {
		return 0
	}

	// Redis 缓存！
	key := fmt.Sprintf("%s:%v:%d:%d", col, query, date, lr)
	r := NewRedisDBCntlr()
	defer r.Close()

	if ok, _ := r.EXISTS(key); ok == 1 {
		n, err := r.GETINT(key)
		if err != nil {
			log.Printf("Get redis cache error at getSingleStatsCount: %s\n", err.Error())
		} else {
			return n
		}
	}

	s := mongoSession.Copy()
	defer s.Close()
	c := s.DB(DBName).C(col)

	var rec StatsRecord

	//dateQuery := bson.M{}
	//
	//if lr == 1 { // l
	//	dateQuery["$gte"] = strconv.Itoa(date) //"$gte"
	//} else {
	//	dateQuery["$lte"] = strconv.Itoa(date) //"$lte"
	//}

	boundary, err := getBoundaryDate(col, query, date, lr)
	if err != nil {
		return 0 // 是不是不太好？
	}

	query = append(query, bson.DocElem{Name:"date", Value:boundary})

	log.Printf("Query:%+v\n", query)

	err = c.Find(query).One(&rec)

	log.Printf("Record:%+v\n", rec)

	if err != nil {
		return 0
	}

	_, err = r.SET(key, rec.Count)
	if err != nil {
		log.Printf("Set redis cache error at getSingleStatsCount: %s\n", err.Error())
	}

	return rec.Count

}

func getCollectionCount(col string) (int, error) {
	s := mongoSession.Copy()
	defer s.Close()
	c := s.DB(DBName).C(col)

	return c.Count()
}

func getInfoValueByKey(key string) (interface{}, error) {
	s := mongoSession.Copy()
	defer s.Close()
	c := s.DB(DBName).C(ColInfo)

	query := bson.M{
		"key": key,
	}

	result := KVPair{}

	err := c.Find(query).One(&result)

	if err != nil {
		return nil, err
	}

	return result.Value, nil
}

func getInfoLastUpdate() (int, error) {
	ret, err := getInfoValueByKey("last_update")
	if err != nil {
		return 0, err
	}

	return ret.(int), err
}

func getInfoInfo() (string, error) {
	ret, err := getInfoValueByKey("info")
	if err != nil {
		return "", err
	}

	return ret.(string), err
}

func getFormulaTotalCount(f Formula, t int) {
	s := mongoSession.Copy()
	defer s.Close()
	defer s.Close()
	var col string

	switch t {
	case RecordTypeTDoll:
		col = ColTDollStats
	case RecordTypeEquip:
		fallthrough
	case RecordTypeFairy:
		col = ColEquipStats
	default:
		panic("Invalid type at GetFormulaTotalCount")
	}

	c := s.DB(DBName).C(col)


}

// GetAddedDate 获得加入建造池的时间 没有返回0
func GetAddedDate(t int, id int) int {
	s := mongoSession.Copy()
	defer s.Close()
	c := s.DB(DBName).C(ColDate)

	query := bson.D{
		{"id", id},
		{"type", t},
	}

	var ret TimeRecord
	err := c.Find(query).One(&ret)
	if err != nil {
		return 0
	}

	return ret.Date
}

// getBoundaryDate 获得有记录存在的边界时间，L为1 R为2
func getBoundaryDate(col string, query bson.D, date int, lr int) (int, error) { // TODO: 可能的redis缓存
	s := mongoSession.Copy()
	defer s.Close()
	c := s.DB(DBName).C(col)

	//TODO: wtf????
	dateQuery := bson.M{} // 心态炸了

	if lr == 1 { // l
		dateQuery["$gte"] = strconv.Itoa(date) //"$gte"
	} else {
		dateQuery["$lte"] = strconv.Itoa(date) //"$lte"
	}

	query = append(query, bson.DocElem{Name: "date", Value: dateQuery})

	var r StatsRecord

	err := c.Find(query).Select(bson.M{"date": 1}).Sort("-date").Limit(1).One(&r)
	if err != nil {
		return 0, err
	}

	log.Printf("%+v\n", r)

	retn, _ :=strconv.Atoi(r.Date)

	return retn, nil
}