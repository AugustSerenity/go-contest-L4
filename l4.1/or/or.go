package or

func Or(channels ...<-chan interface{}) <-chan interface{} {
	done := make(chan interface{})
	signal := make(chan struct{})

	if len(channels) == 0 {
		close(done)
		return done
	}

	for _, ch := range channels {
		go func(inCh <-chan interface{}) {
			<-inCh

			select {
			case signal <- struct{}{}:
			default:
			}
		}(ch)
	}

	go func() {
		<-signal
		close(done)
	}()

	return done
}
