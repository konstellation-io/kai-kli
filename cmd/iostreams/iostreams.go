package iostreams

import (
	"io"
	"os"
	"time"

	"github.com/mattn/go-isatty"

	"github.com/mattn/go-colorable"

	"github.com/briandowns/spinner"
)

const (
	spinnerDuration = 400 * time.Millisecond
)

// IOStreams contains information for input/output used in the CLI.
type IOStreams struct {
	In     io.ReadCloser
	Out    io.Writer
	ErrOut io.Writer

	progressIndicator *spinner.Spinner

	progressIndicatorEnabled bool

	stdinTTYOverride  bool
	stdinIsTTY        bool
	stdoutTTYOverride bool
	stdoutIsTTY       bool
	stderrTTYOverride bool
	stderrIsTTY       bool
}

// System returns an IOStreams object with prepared input/output.
func System() *IOStreams {
	stdoutIsTTY := isTerminal(os.Stdout)
	stderrIsTTY := isTerminal(os.Stderr)

	st := &IOStreams{
		In:     os.Stdin,
		Out:    colorable.NewColorable(os.Stdout),
		ErrOut: colorable.NewColorable(os.Stderr),
	}

	if stdoutIsTTY && stderrIsTTY {
		st.progressIndicatorEnabled = true
	}

	// prevent duplicate isTerminal queries now that we know the answer
	st.SetStdoutTTY(stdoutIsTTY)
	st.SetStderrTTY(stderrIsTTY)

	return st
}

// SetStdinTTY marks stdin as TTY.
func (s *IOStreams) SetStdinTTY(isTTY bool) {
	s.stdinTTYOverride = true
	s.stdinIsTTY = isTTY
}

// IsStdinTTY checks if stdin is a TTY.
func (s *IOStreams) IsStdinTTY() bool {
	if s.stdinTTYOverride {
		return s.stdinIsTTY
	}

	if stdin, ok := s.In.(*os.File); ok {
		return isTerminal(stdin)
	}

	return false
}

// SetStdoutTTY marks stdout as TTY.
func (s *IOStreams) SetStdoutTTY(isTTY bool) {
	s.stdoutTTYOverride = true
	s.stdoutIsTTY = isTTY
}

// IsStdoutTTY checks if stdout is a TTY.
func (s *IOStreams) IsStdoutTTY() bool {
	if s.stdoutTTYOverride {
		return s.stdoutIsTTY
	}

	if stdout, ok := s.Out.(*os.File); ok {
		return isTerminal(stdout)
	}

	return false
}

// SetStderrTTY marks stderr as TTY.
func (s *IOStreams) SetStderrTTY(isTTY bool) {
	s.stderrTTYOverride = true
	s.stderrIsTTY = isTTY
}

// IsStderrTTY checks if stderr is a TTY.
func (s *IOStreams) IsStderrTTY() bool {
	if s.stderrTTYOverride {
		return s.stderrIsTTY
	}

	if stderr, ok := s.ErrOut.(*os.File); ok {
		return isTerminal(stderr)
	}

	return false
}

// StartProgressIndicator creates a new spinner indicator.
func (s *IOStreams) StartProgressIndicator(str string) {
	if !s.progressIndicatorEnabled {
		return
	}

	sp := spinner.New(spinner.CharSets[11], spinnerDuration, spinner.WithWriter(s.ErrOut))
	sp.Start()

	if str != "" {
		sp.Suffix = str
	}

	s.progressIndicator = sp
}

// StopProgressIndicator stops the running spinner indicator.
func (s *IOStreams) StopProgressIndicator() {
	if s.progressIndicator == nil {
		return
	}

	s.progressIndicator.Stop()
	s.progressIndicator = nil
}

func isTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}
