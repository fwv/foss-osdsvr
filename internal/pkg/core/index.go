package core

type Index interface {
	Put(oid int64, entry *Entry) error
	Get(oid int64) (entry *Entry)
}
