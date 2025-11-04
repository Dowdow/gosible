package runner

import (
	"bufio"
	"bytes"

	tea "github.com/charmbracelet/bubbletea"
)

type teaWriterStdout struct {
	ch  chan tea.Msg
	buf bytes.Buffer
}

func newTeaWriterStdout(ch chan tea.Msg) *teaWriterStdout {
	return &teaWriterStdout{ch: ch}
}

func (w *teaWriterStdout) Write(p []byte) (n int, err error) {
	n, err = w.buf.Write(p)
	scanner := bufio.NewScanner(&w.buf)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			w.ch <- StdoutMsg{Msg: line}
		}
	}
	w.buf.Reset()
	return n, err
}

type teaWriterSterr struct {
	ch  chan tea.Msg
	buf bytes.Buffer
}

func newTeaWriterStderr(ch chan tea.Msg) *teaWriterSterr {
	return &teaWriterSterr{ch: ch}
}

func (w *teaWriterSterr) Write(p []byte) (n int, err error) {
	n, err = w.buf.Write(p)
	scanner := bufio.NewScanner(&w.buf)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			w.ch <- StderrMsg{Msg: line}
		}
	}
	w.buf.Reset()
	return n, err
}
