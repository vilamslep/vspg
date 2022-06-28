package notice

type Messager interface {
	BuildMessage() ([]byte, error)
}