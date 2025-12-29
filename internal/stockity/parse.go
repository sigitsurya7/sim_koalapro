package stockity

import "strconv"

func ParseUserID(raw string) (int64, error) {
	return strconv.ParseInt(raw, 10, 64)
}
