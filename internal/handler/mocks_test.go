package handler

type mockURLService struct {
	resolveURL func(id string) (string, error)
	shortenURL func(url string) (string, error)
}

func (m *mockURLService) ShortenURL(s string) (string, error) {
	return m.shortenURL(s)
}

func (m *mockURLService) ResolveURL(id string) (string, error) {
	return m.resolveURL(id)
}
