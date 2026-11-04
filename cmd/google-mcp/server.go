package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

// server is google-mcp-server run as Claude Code runs it (fnox exec -- google-mcp-server), spoken to over stdio:
// JSON-RPC lines in, its stdout and stderr kept line by line.
type server struct {
	cmd                *exec.Cmd
	in                 io.WriteCloser
	mu                 sync.Mutex
	outLines, errLines *[]string
	once               sync.Once
}

func startServer() (*server, error) {
	s := &server{cmd: exec.Command("fnox", "exec", "--", "google-mcp-server"), outLines: &[]string{}, errLines: &[]string{}}
	var err error
	if s.in, err = s.cmd.StdinPipe(); err != nil {
		return nil, err
	}
	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := s.cmd.Start(); err != nil {
		return nil, err
	}
	go s.collect(stdout, s.outLines)
	go s.collect(stderr, s.errLines)
	return s, nil
}

func (s *server) collect(r io.Reader, into *[]string) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	for sc.Scan() {
		s.mu.Lock()
		*into = append(*into, sc.Text())
		s.mu.Unlock()
	}
}

func (s *server) send(msg string) { fmt.Fprintln(s.in, msg) }

// initialize is the MCP handshake.
func (s *server) initialize() {
	s.send(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"mise","version":"0"}}}`)
	s.send(`{"jsonrpc":"2.0","method":"notifications/initialized"}`)
}

// wait is the first line of lines matching re, waiting up to d for it.
func (s *server) wait(lines *[]string, re *regexp.Regexp, d time.Duration) (string, bool) {
	for deadline := time.Now().Add(d); time.Now().Before(deadline); time.Sleep(200 * time.Millisecond) {
		s.mu.Lock()
		for _, l := range *lines {
			if re.MatchString(l) {
				s.mu.Unlock()
				return l, true
			}
		}
		s.mu.Unlock()
	}
	return "", false
}

// tail is the last n lines of lines.
func (s *server) tail(lines *[]string, n int) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	l := *lines
	if len(l) > n {
		l = l[len(l)-n:]
	}
	return strings.Join(l, "\n")
}

func (s *server) stop() {
	s.once.Do(func() {
		s.in.Close()
		if s.cmd.Process != nil {
			_ = s.cmd.Process.Kill()
			_ = s.cmd.Wait()
		}
	})
}
