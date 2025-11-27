package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
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

	var flatData [][]string
	for _, chunk := range processed {
		flatData = append(flatData, chunk...)
	}

	for _, p := range s.Peers {
		if p != "" && p != s.Addr {
			fmt.Fprintf(os.Stderr, "Sending to peer: %s\n", p)
			if err := SendChunks(p, flatData); err != nil {
				fmt.Fprintf(os.Stderr, "Error sending to %s: %v\n", p, err)
			}
		}
	}

	for _, chunk := range flatData {
		for _, line := range chunk {
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
		fmt.Fprintf(os.Stderr, "Received %d lines\n", len(data))

		for _, arr := range data {
			for _, line := range arr {
				collected = append(collected, line)
			}
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
			out[i] = s.CutLines(chunk)
			wg.Done()
		}()
	}

	wg.Wait()
	return out
}

func SendChunks(addr string, chunks [][]string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to dial %s: %v", addr, err)
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(chunks); err != nil {
		return fmt.Errorf("failed to encode chunks: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Sent %d lines to %s\n", len(chunks), addr)
	return nil
}

func MergeAll(lines []string) []string {
	sort.Strings(lines)
	return lines
}

func SplitInput(r io.Reader) [][]string {
	scanner := bufio.NewScanner(r)
	chunk := make([]string, 0, 1000)
	var chunks [][]string

	for scanner.Scan() {
		chunk = append(chunk, scanner.Text())
		if len(chunk) == 1000 {
			chunks = append(chunks, chunk)
			chunk = make([]string, 0, 1000)
		}
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks
}
