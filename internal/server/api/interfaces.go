package api

type Config interface {
	ServAPIAddr() string
	DirWithFiles() string
	LogLevel() string
	String() string
}
