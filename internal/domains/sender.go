package domains

type Sender interface {
	StartSenderWorker(input <-chan *Notification, done chan<- struct{})
}
