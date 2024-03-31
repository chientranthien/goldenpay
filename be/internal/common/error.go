package common

func ErrToStr(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}

func FatalIfErr(err error, msg ...any) {
	if err == nil {
		return
	}
	L().Fatalw("fatalDueToErr", "err", err, msg)
}
