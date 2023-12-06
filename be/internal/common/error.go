package common

func ErrToStr(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}
