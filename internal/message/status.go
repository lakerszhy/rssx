package message

const (
	statusInitial status = iota
	statusInProgress
	statusSuccessful
	statusFailed
)

type status int

func (s status) IsInitial() bool {
	return s == statusInitial
}

func (s status) IsInProgress() bool {
	return s == statusInProgress
}

func (s status) IsSuccessful() bool {
	return s == statusSuccessful
}

func (s status) IsFailed() bool {
	return s == statusFailed
}
