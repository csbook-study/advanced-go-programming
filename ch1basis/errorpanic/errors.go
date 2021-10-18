package errorpanic

// github.com/chai2010/errors
type Error interface {
	Caller() []CallerInfo
	Wraped() []error
	Code() int
	error

	private()
}

type CallerInfo struct {
	FuncName string
	FileName string
	FileLine int
}

func New(msg string) error {
	return nil
}

func NewWithCode(code int, msg string) error {
	return nil
}

func Wrap(err error, msg string) error {
	return nil
}

func WrapWithCode(code int, err error, msg string) error {
	return nil
}

func FromJson(json string) (Error, error) {
	return nil, nil
}

func ToJson(err error) string {
	return ""
}
