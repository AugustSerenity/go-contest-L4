package or

func or(channels ...<-chan interface{}) <-chan interface{} {
	done := make(chan interface{})
	go func() {
		defer close(done)

		for _, ch := range channels {
			go func(c <-chan interface{}) {
				<-c
				done <- struct{}{}
			}(ch)
		}
		<-done
	}()

	return done
}
