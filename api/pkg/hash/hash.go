package hash

type Hash interface {
	New(str string, opt *Options) string
	Compare(str, hash string) bool
}

func New(str string, opt *Options) string {
	return Backend.New(str, opt)
}

func Compare(str, hash string) bool {
	return Backend.Compare(str, hash)
}
