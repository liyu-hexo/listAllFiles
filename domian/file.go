package domian

type File struct {
	Id   string `xorm:"ID pk"`
	Name string `xorm:"Name"`
	Path string `xorm:"Path"`
	Time int64  `xorm:"Time"`
}
