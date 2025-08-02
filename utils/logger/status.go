package logger

type Status int

const (
	Pending Status = iota
	InProgress
	Completed
	Failed
)

func (s Status) String() string {
	switch s {
	case Pending:
		return "Pending"
	case InProgress:
		return "InProgress"
	case Completed:
		return "Completed"
	case Failed:
		return "Failed"
	default:
		return "Unknown"
	}
}
