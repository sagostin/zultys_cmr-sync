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

// CustomFtpDriver implements the ftpserverlib.MainDriver interface.
type CustomFtpDriver struct {
	// No additional fields needed for this simple example.
	Username    string
	Password    string
	ListenAddr  string
	CsvDataChan chan string
}

func (d *CustomFtpDriver) GetTLSConfig() (*tls.Config, error) {
	//TODO implement me
	panic("implement me")
}

// CustomFtpClientDriver implements the ClientDriverExt interface.
// It holds the connection-specific context.
type CustomFtpClientDriver struct {
	afero.Fs
	*CustomFtpDriver
}

// AuthUser checks the username and password.
func (d *CustomFtpDriver) AuthUser(cc ftpserverlib.ClientContext, user, pass string) (ftpserverlib.ClientDriver, error) {
	if user == d.Username && pass == d.Password {
		return &CustomFtpClientDriver{Fs: MemFs, CustomFtpDriver: d}, nil
	}
	return nil, errors.New("authentication failed")
}

// ClientConnected is called when a new client is connected.
func (d *CustomFtpDriver) ClientConnected(cc ftpserverlib.ClientContext) (string, error) {
	return "Welcome to the FTP server", nil
}

// ClientDisconnected is called when the client disconnects.
func (d *CustomFtpDriver) ClientDisconnected(cc ftpserverlib.ClientContext) {
	// Cleanup or logging can be done here.
}

// GetSettings returns the server settings.
func (d *CustomFtpDriver) GetSettings() (*ftpserverlib.Settings, error) {
	return &ftpserverlib.Settings{
		// Customize your server settings here.
		// todo allow to specify the ftp server binding address??
		ListenAddr: d.ListenAddr,
	}, nil
}

// OpenFile opens a file for reading or writing.
func (d *CustomFtpClientDriver) OpenFile(path string, flag int, perm os.FileMode) (afero.File, error) {
	file, err := d.Fs.OpenFile(path, flag, perm)
	if err != nil {
		return nil, err
	}

	// Wrap the file in a custom type to intercept the Close call
	// so we can process the file's contents before it's closed.
	return &CustomFtpFile{File: file, path: path, CustomFtpDriver: d.CustomFtpDriver}, nil
}

// CustomFtpFile wraps an afero.File to process the file's content on Close.
type CustomFtpFile struct {
	afero.File
	path string
	*CustomFtpDriver
}

// Close processes the file's content if it's a CSV before closing it.
func (f *CustomFtpFile) Close() error {
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
		fmt.Printf("Processing CSV content from path %s", f.path)
		f.CsvDataChan <- content
		log.Info("sending content to chan")
		// todo output file to channel for processing?????
	} else {
		fmt.Println("Uploaded file is not a CSV.")
	}

	return f.File.Close()
}
