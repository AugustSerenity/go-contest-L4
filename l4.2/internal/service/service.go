package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/AugustSerenity/go-contest-L4/l4.2_CLI-sort/internal/handler"
)

type Service struct {
	Addr     string
	Peers    []string
	Opts     Options
	Incoming chan [][]string
}

func New(addr string, peers []string, opts Options) *Service {
	return &Service{
		Addr:     addr,
		Peers:    peers,
		Opts:     opts,
		Incoming: make(chan [][]string, 10),
	}
}

func (s *Service) RunServer(in io.Reader, out io.Writer) {
	fmt.Fprintf(os.Stderr, "Starting server on %s\n", s.Addr)

	go handler.StartListener(s.Addr, s)

	chunks := SplitInput(in)
	processed := s.ProcessChunks(chunks)

	time.Sleep(100 * time.Millisecond)

	for _, p := range s.Peers {
		if p != "" && p != s.Addr {
			fmt.Fprintf(os.Stderr, "Sending to peer: %s\n", p)
			if err := SendChunks(p, processed); err != nil {
				fmt.Fprintf(os.Stderr, "Error sending to %s: %v\n", p, err)
			}
		}
	}

	for _, part := range processed {
		for _, line := range part {
			fmt.Fprintln(out, line)
		}
	}

	fmt.Fprintf(os.Stderr, "Server %s completed\n", s.Addr)
}

func (s *Service) RunCoordinator(out io.Writer) {
	fmt.Fprintln(os.Stderr, "Running coordinator on", s.Addr)

	go handler.StartListener(s.Addr, s)

	N := len(s.Peers)
	collected := []string{}

	for i := 0; i < N; i++ {
		fmt.Fprintf(os.Stderr, "Waiting for data from peers (%d/%d)\n", i+1, N)
		data := <-s.Incoming
		fmt.Fprintf(os.Stderr, "Received chunk with %d lines\n", len(data))
		for _, chunk := range data {
			collected = append(collected, chunk...)
		}
	}

	fmt.Fprintf(os.Stderr, "Total collected lines: %d\n", len(collected))
	merged := MergeAll(collected)

	for _, line := range merged {
		fmt.Fprintln(out, line)
	}

	fmt.Fprintf(os.Stderr, "Coordinator completed\n")
}

func (s *Service) HandleIncomingChunk(data [][]string) {
	s.Incoming <- data
}

func (s *Service) ProcessChunks(chunks [][]string) [][][]string {
	out := make([][][]string, len(chunks))
	wg := sync.WaitGroup{}

	for i, chunk := range chunks {
		i, chunk := i, chunk
		wg.Add(1)
		go func() {
			out[i] = [][]string{s.CutLines(chunk)}
			wg.Done()
		}()
	}

	wg.Wait()
	return out
}

func SendChunks(addr string, chunks [][][]string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to dial %s: %v", addr, err)
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(chunks); err != nil {
		return fmt.Errorf("failed to encode chunks: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Sent %d chunks to %s\n", len(chunks), addr)
	return nil
}
