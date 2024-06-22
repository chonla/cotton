package httphelper

type RequestParser interface {
	Parse(requestString string) (Request, error)
}
