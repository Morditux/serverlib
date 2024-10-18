package sessions

type Session interface {
	Id() string
	Get(key string) (any, bool)
	Set(key string, value any)
}

type Sessions interface {
	Get(id string) (Session, bool)
	Set(id string, session Session)
	Delete(id string)
}
