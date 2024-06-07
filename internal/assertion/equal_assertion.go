package assertion

type EqualAssertion struct {
}

func (a *EqualAssertion) Name() string {
	return "eq"
}

func (a *EqualAssertion) Assert(actual, expected interface{}) (bool, error) {
	return false, nil
}
