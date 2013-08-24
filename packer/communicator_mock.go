package packer

import (
	"bytes"
	"io"
	"sync"
)

// MockCommunicator is a valid Communicator implementation that can be
// used for tests.
type MockCommunicator struct {
	StartCalled     bool
	StartCmd        *RemoteCmd
	StartStderr     string
	StartStdout     string
	StartStdin      string
	StartExitStatus int

	UploadCalled bool
	UploadPath   string
	UploadData   string

	DownloadCalled bool
	DownloadPath   string
	DownloadData   string
}

func (c *MockCommunicator) Start(rc *RemoteCmd) error {
	c.StartCalled = true
	c.StartCmd = rc

	go func() {
		var wg sync.WaitGroup
		if rc.Stdout != nil && c.StartStdout != "" {
			wg.Add(1)
			go func() {
				rc.Stdout.Write([]byte(c.StartStdout))
				wg.Done()
			}()
		}

		if rc.Stderr != nil && c.StartStderr != "" {
			wg.Add(1)
			go func() {
				rc.Stderr.Write([]byte(c.StartStderr))
				wg.Done()
			}()
		}

		if rc.Stdin != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var data bytes.Buffer
				io.Copy(&data, rc.Stdin)
				c.StartStdin = data.String()
			}()
		}

		wg.Wait()
		rc.SetExited(c.StartExitStatus)
	}()

	return nil
}

func (c *MockCommunicator) Upload(path string, r io.Reader) error {
	c.UploadCalled = true
	c.UploadPath = path

	var data bytes.Buffer
	if _, err := io.Copy(&data, r); err != nil {
		panic(err)
	}

	c.UploadData = data.String()

	return nil
}

func (c *MockCommunicator) UploadDir(string, string, []string) error {
	return nil
}

func (c *MockCommunicator) Download(path string, w io.Writer) error {
	c.DownloadCalled = true
	c.DownloadPath = path
	w.Write([]byte(c.DownloadData))

	return nil
}