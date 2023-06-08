package migrations

type Schema_migration struct {
	Version int  `orm:"version" json:"version"`
	Dirty   bool `orm:"dirty" json:"dirty"`
}
