package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"reflect"
	"strings"

	"xorm.io/builder"
	"xorm.io/xorm"
)

var MasterDb *xorm.Engine

func InitDatabase(db_url string) error {
	var err error
	if MasterDb, err = xorm.NewEngine("mysql", db_url); err != nil {
		return err
	}
	//MasterDb.SetConnMaxLifetime(config.GetOpts().Mysql.ConnMaxLifetime)
	MasterDb.SetMaxOpenConns(1000)
	MasterDb.SetMaxIdleConns(10)

	return nil
}

type Joblist struct {
	Id   int64  `orm:"column(id);auto"`
	Name string `orm:"column(name);size(64)"`
	Pass string `orm:"column(pass);size(512)"`
}

func (t *Joblist) TableName() string {
	return "job_list"
}

//func init() {
//	orm.RegisterModel(new(Joblist))
//}

// func GetUserByNamePass(user *Joblist) error {
// 	return orm.NewOrm().Read(user, "name", "pass")
// }

func AddJoblist(job *Joblist) error {
	_, err := MasterDb.InsertOne(job)
	//_, err := orm.NewOrm().Insert(job)
	return err
}

func DeleteJob(job *Joblist) error {
	//语句示例，未执行
	job.Id = 2
	_, err := MasterDb.Delete(job)
	//_, err := orm.NewOrm().Delete(job)
	return err
}

func UpdateJobById(job *Joblist) error {
	//语句示例，未执行
	_, err := MasterDb.Delete(job)
	//_, err := orm.NewOrm().Update(job)
	return err
}

func SearchJob() error {
	//语句示例，未执行
	sql := builder.Dialect(builder.MYSQL).
		Select("*").
		From("joblist").
		Where(builder.Eq{
			".is_del": 1,
		}).And(builder.In("id", 1)).OrderBy("aaa ASC")

	var ps []Joblist
	if err := MasterDb.SQL(sql).Find(&ps); err != nil {

	}
	return nil
}

//根据条件筛选数据
//query 根据字段值筛选，类似于：where name=“zhangsan”
//fields 查询字段，类似于： select name，pass from 。。。
//sortby 排序字段，例如根据Id排序: order by Id ...
//order 结合sortby字段排序方式，asc ：正序，desc : order by Id asc ...
//offset 游标 ，表示从搜索数据中的第几条开始，常用来分页
//limit  筛选数量
func GetJobs(keys string, upSala, lowSala int64, query map[string]string, fields []string, sortby []string, order []string, offset int64, limit int64) (ml []interface{}, err error) {
	//cond := orm.NewCondition()
	//cond.Or()
	o := orm.NewOrm()
	qs := o.QueryTable(new(Joblist))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Joblist
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}
