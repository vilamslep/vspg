package notice

type Sender interface {
	Send(Messager) error
}