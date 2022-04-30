package models

type StatusType int

const (
	NEW = iota + 1
	PROCESSING
	INVALID
	PROCESSED
	WITHDRAWN
)

func (s StatusType) String() string {
	return [...]string{"NEW", "PROCESSING", "INVALID", "PROCESSED", "WITHDRAWN"}[s-1]
}

func (s StatusType) EnumIndex() int {
	return int(s)
}
