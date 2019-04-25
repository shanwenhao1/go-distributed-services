package spider

type Spider interface {
	GetBody() string
}

func GetGoVersion(s Spider) string {
	body := s.GetBody()
	return body
}
