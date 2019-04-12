package dbapi

type UserAPI interface {
	Auth(id, authstr string) bool
	Detail(id, authstr string) (map[string]string, error)
}

type StudentAPI interface {
	ListTutorial() ([]map[string]string, error)
}
