package logging

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/logger.go -package=mocks -mock_names=Interface=MockLogger

type Interface interface {
	Success(msg string)
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
	SetDebugLevel()
}
