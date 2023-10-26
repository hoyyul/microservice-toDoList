package errmsg

var Msgflags = map[int]string{
	SUCCESS:       "Success",
	FAILURE:       "Failure",
	InvalidParams: "Invalid parameter",
}

func GetMsg(code int) string {
	return Msgflags[code]
}
