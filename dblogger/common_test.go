package dblogger

const (
	tableBasicName = "basic"
)

type Basic struct {
	Uid  float64 `json:"uid" gorm:"column:uid;autoIncrement;primaryKey;type:bigint;comment:自增id"`
	Name string  `json:"name" gorm:"column:name;type:varchar(256);comment:名称"`
	Desc string  `json:"desc" gorm:"column:desc;type:varchar(256);comment:描述"`
	Age  int     `json:"age" gorm:"column:age;type:int;comment:年龄"`
	Sex  int     `json:"sex" gorm:"column:sex;type:int;comment:性别"`
}

func (*Basic) TableName() string {
	return tableBasicName
}
