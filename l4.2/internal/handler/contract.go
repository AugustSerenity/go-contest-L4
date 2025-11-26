package handler

type Service interface {
	HandleIncomingChunk([][]string)
}
