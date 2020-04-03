package mongodb

import (
	"errors"
	"strings"
	"sync"

	"gopkg.in/mgo.v2/bson"
)

var (
	Mongo    *MongoDb
	initonce sync.Once
)

type MongoDb struct {
	Url string
	DB  string
	c   *dialContext
}

func DitalMongo(url string) {
	initonce.Do(func() {
		initMongDb(url)
	})
}

func CloseMongo() {
	Mongo.CloseMongDb()
}

func initMongDb(url string) {
	con, err := dial(url, 10)
	if err != nil {
		panic(err)
	}
	db := url[strings.LastIndex(url, "/")+1:]
	if len(db) <= 0 {
		panic(errors.New("url format error"))
	}
	Mongo = &MongoDb{
		Url: url,
		DB:  db,
		c:   con,
	}
}

func (m *MongoDb) InsertData(collection string, docs interface{}) (err error) {
	s := m.c.Ref()
	defer m.c.UnRef(s)
	return s.DB(m.DB).C(collection).Insert(docs)
}

func (m *MongoDb) IsExistKey(collection string, bs bson.M) (bool, error) {
	s := m.c.Ref()
	defer m.c.UnRef(s)
	var (
		value interface{}
	)
	if err := s.DB(m.DB).C(collection).Find(bs).One(value); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *MongoDb) QueryData(collection string, bs bson.M, value interface{}) (err error) {
	s := m.c.Ref()
	defer m.c.UnRef(s)
	err = s.DB(m.DB).C(collection).Find(bs).One(value)
	return
}

func (m *MongoDb) GetCount(collection string, bs bson.M) (count int, err error) {
	s := m.c.Ref()
	defer m.c.UnRef(s)
	return s.DB(m.DB).C(collection).Find(bs).Count()
}

func (m *MongoDb) QueryAllData(collection string, bs bson.M, value interface{}) (err error) {
	s := m.c.Ref()
	defer m.c.UnRef(s)
	return s.DB(m.DB).C(collection).Find(bs).All(value)
}

func (m *MongoDb) UpdateData(collection string, bs bson.M, value interface{}) (err error) {
	s := m.c.Ref()
	defer m.c.UnRef(s)
	return s.DB(m.DB).C(collection).Update(bs, value)
}

func (m *MongoDb) UpdateDataById(collection string, id, value interface{}) (err error) {
	s := m.c.Ref()
	defer m.c.UnRef(s)
	return s.DB(m.DB).C(collection).UpdateId(id, value)
}

func (m *MongoDb) DelDataById(collection string, id interface{}) (err error) {
	s := m.c.Ref()
	defer m.c.UnRef(s)
	return s.DB(m.DB).C(collection).RemoveId(id)

}
func (m *MongoDb) DelData(collection string, bs bson.M) error {
	s := m.c.Ref()
	defer m.c.UnRef(s)
	return s.DB(m.DB).C(collection).Remove(bs)
}
func (m *MongoDb) CloseMongDb() {
	m.c.Close()
}
