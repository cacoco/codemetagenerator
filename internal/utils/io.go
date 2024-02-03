package utils

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// WRITER

type Writer interface {
	Stdout() io.WriteCloser
	Println(a ...any) (n int, err error)
	Print(a ...any) (n int, err error)
	Errorf(format string, a ...any) error
}

// TEST WRITER

type nopWriteCloser struct {
	io.WriteCloser
}

func (*nopWriteCloser) Write(buf []byte) (int, error) {
	return len(buf), nil
}

func (*nopWriteCloser) Close() error {
	return nil
}

type (
	TestWriter struct {
		Writer
	}
)

func (r *TestWriter) Stdout() io.WriteCloser {
	return &nopWriteCloser{}
}

func (r *TestWriter) Print(a ...any) (n int, err error) {
	return 0, nil
}
func (r *TestWriter) Println(a ...any) (n int, err error) {
	return 0, nil
}
func (r *TestWriter) Errorf(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

// STDOUT WRITER

type StdoutWriter struct {
	Writer
}

func (*StdoutWriter) Stdout() io.WriteCloser {
	return os.Stdout
}
func (s *StdoutWriter) Print(a ...any) (n int, err error) {
	return fmt.Fprint(s.Stdout(), a...)
}
func (s *StdoutWriter) Println(a ...any) (n int, err error) {
	return fmt.Fprintln(s.Stdout(), a...)
}
func (s *StdoutWriter) Errorf(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

// READER

type Reader interface {
	Stdin() io.ReadCloser
}

// TEST READER

type (
	TestReader struct {
		Reader
		In TestStdin
	}
	TestStdin struct {
		m sync.Mutex
		io.ReadCloser
		Data Stack[string]
		done bool
	}
)

func (r *TestReader) Stdin() io.ReadCloser {
	return &r.In
}
func (r *TestStdin) Read(buf []byte) (int, error) {
	r.m.Lock()
	defer r.m.Unlock()
	if r.Data.Len() == 0 {
		return 0, io.EOF
	}
	d := r.Data.Pop()
	for i, b := range []byte(*d) {
		buf[i] = b
	}
	r.done = true
	return len(*d), nil
}
func (*TestStdin) Close() error { return nil }

// STDIN READER

type StdinReader struct {
	Reader
}

func (*StdinReader) Stdin() io.ReadCloser {
	return os.Stdin
}
