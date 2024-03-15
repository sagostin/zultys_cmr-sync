package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	ftpserverlib "github.com/fclairamb/ftpserverlib"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"os"
	"strings"
)

// MemFs is an in-memory filesystem for storing uploaded files temporarily.
var MemFs = afero.NewMemMapFs()

// CustomDriver implements the ftpserverlib.MainDriver interface.
type CustomDriver struct {
	// No additional fields needed for this simple example.
	Username    string
	Password    string
	ListenAddr  string
	CsvDataChan chan string
}

func (d *CustomDriver) GetTLSConfig() (*tls.Config, error) {
	//TODO implement me
	panic("implement me")
}

// ClientDriverExt extends the ftpserverlib.ClientDriver interface with custom methods.
type ClientDriverExt interface {
	ftpserverlib.ClientDriver
	SaveFileToMemory(path string, data []byte) error
}

// CustomClientDriver implements the ClientDriverExt interface.
// It holds the connection-specific context.
type CustomClientDriver struct {
	afero.Fs
	*CustomDriver
}

// AuthUser checks the username and password.
func (d *CustomDriver) AuthUser(cc ftpserverlib.ClientContext, user, pass string) (ftpserverlib.ClientDriver, error) {
	if user == d.Username && pass == d.Password {
		return &CustomClientDriver{Fs: MemFs, CustomDriver: d}, nil
	}
	return nil, errors.New("authentication failed")
}

// ClientConnected is called when a new client is connected.
func (d *CustomDriver) ClientConnected(cc ftpserverlib.ClientContext) (string, error) {
	return "Welcome to the FTP server", nil
}

// ClientDisconnected is called when the client disconnects.
func (d *CustomDriver) ClientDisconnected(cc ftpserverlib.ClientContext) {
	// Cleanup or logging can be done here.
}

// GetSettings returns the server settings.
func (d *CustomDriver) GetSettings() (*ftpserverlib.Settings, error) {
	return &ftpserverlib.Settings{
		// Customize your server settings here.
		// todo allow to specify the ftp server binding address??
		ListenAddr: d.ListenAddr,
	}, nil
}

// OpenFile opens a file for reading or writing.
func (d *CustomClientDriver) OpenFile(path string, flag int, perm os.FileMode) (afero.File, error) {
	file, err := d.Fs.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}

	// Wrap the file in a custom type to intercept the Close call
	// so we can process the file's contents before it's closed.
	return &CustomFile{File: file, path: path, CustomDriver: d.CustomDriver}, nil
}

// CustomFile wraps an afero.File to process the file's content on Close.
type CustomFile struct {
	afero.File
	path string
	*CustomDriver
}

// Close processes the file's content if it's a CSV before closing it.
func (f *CustomFile) Close() error {
	// Check if the file is indeed a CSV based on its content.
	var buf bytes.Buffer
	if _, err := f.Seek(0, os.SEEK_SET); err != nil {
		return err
	}
	if _, err := buf.ReadFrom(f.File); err != nil {
		return err
	}

	content := buf.String()
	if strings.Contains(content, ",") && strings.HasSuffix(f.path, ".csv") {
		// Process the CSV content here.
		fmt.Printf("Processing CSV content from path %s: \n%s\n\n", f.path, content)
		f.CsvDataChan <- content
		log.Info("sending content to chan")
		// todo output file to channel for processing?????
	} else {
		fmt.Println("Uploaded file is not a CSV.")
	}

	return f.File.Close()
}
