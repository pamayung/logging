package wlog

type m map[string]interface{}

func InternalError() m {
	return m{"status": 500, "message": "Internal Server Error"}
}
