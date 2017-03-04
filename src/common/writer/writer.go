/*
* @Author: Ujjwal Kanth
* @Email: ujjwal.kanth@unbxd.com
* @File: log.go
 */
package writer

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"archive/zip"
	"io"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var errCreatingFile = errors.New("Error creating new file")
var errClosingFile = errors.New("Error closing File")
var errRenamingFile = errors.New("Error renaming File")

const (
	s3BucketName string = "unbxd-wingman-logs"
	awsRegion    string = "ap-southeast-1"
)

type RotateWriter struct {
	lock     sync.Mutex
	filename string // should be set to the actual filename
	duration int
	fp       *os.File
	maxDays  int
}

// Make a new RotateWriter. Return nil if error occurs during setup.
func NewWriter(filename string) (*RotateWriter, error) {
	// Check file before we initialize.
	// no logs will be cleaned up or archived for 0 maxDays
	return new(filename, 3600, 2)
}

func new(filename string, duration int, maxDays int) (*RotateWriter, error) {
	w := &RotateWriter{filename: filename, duration: duration, maxDays: maxDays}
	err := w.Rotate()
	if err != nil {
		return nil, err
	}
	// Trigger a rotation after every x interval.
	ticker := time.NewTicker(time.Duration(int64(w.duration)) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.Rotate()
			}
		}
	}()
	return w, nil
}

// Write satisfies the io.Writer interface.
func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.fp.Write(output)
}

func (w *RotateWriter) Close() (err error) {
	return w.fp.Close()
}

// Perform the actual act of rotating and reopening file.
func (w *RotateWriter) Rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// Close existing file if open
	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return errClosingFile
		}
	}
	// Rename dest file if it already exists
	_, err = os.Stat(w.filename)
	if err == nil {
		err = os.Rename(w.filename, w.filename+"."+time.Now().Format(time.RFC3339))
		if err != nil {
			return errRenamingFile
		}
	}

	// Create a file.
	w.fp, err = os.Create(w.filename)
	if err != nil {
		return errCreatingFile
	}

	err = w.cleanup()
	if err != nil {
		return err
	}
	return nil
}

func (w *RotateWriter) read() ([]string, error) {
	var buffer []string

	fp, err := os.Open(w.filename)
	if err != nil {
		return nil, err
	}
	// File is already open we use the scanner on it.
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		text := scanner.Text()
		buffer = append(buffer, text)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return buffer, nil
}

func (w *RotateWriter) cleanup() error {
	if w.maxDays == 0 {
		return nil
	}
	//fetch old files
	oldLogs, err := w.oldLogs()
	if err != nil {
		return err
	}
	w.cleanupWithFiles(oldLogs)
	return nil
}

func (w *RotateWriter) cleanupWithFiles(files []os.FileInfo) {
	go w.pushToS3AndDelete(files)
}

func (w *RotateWriter) pushToS3AndDelete(files []os.FileInfo) {
	w.pushToS3(files)
	w.deleteOldLogs(files)
}

func (w *RotateWriter) pushToS3(files []os.FileInfo) {
	svc := s3.New(session.New(), &aws.Config{Region: aws.String(awsRegion)})
	for _, file := range files {
		//get dir pattern
		timeStr := strings.TrimPrefix(file.Name(), w.filename+".")
		t, err := time.Parse(time.RFC3339, timeStr)
		dir := strconv.Itoa(t.Year()) + "/" + t.Month().String() + "/" + strconv.Itoa(t.Day()) + "/"
		//create zip file
		zipper(file.Name(), file.Name()+".zip")
		zipF, err := os.Open(file.Name() + ".zip")
		defer zipF.Close()
		if err != nil {
			log.Info("Error opening zipF "+zipF.Name(), err)
			return
		}
		fileInfo, _ := zipF.Stat()
		// todo: push to day directory
		var size int64 = fileInfo.Size()
		buffer := make([]byte, size)
		zipF.Read(buffer)
		fileBytes := bytes.NewReader(buffer)
		params := &s3.PutObjectInput{
			Bucket:        aws.String(s3BucketName),
			Key:           aws.String(dir + zipF.Name()),
			Body:          fileBytes,
			ContentLength: aws.Int64(size),
		}
		_, err = svc.PutObject(params)
		if err != nil {
			log.Info("Error on uploading ", err)
			return
		}
		//remove zip file post upload
		os.Remove(zipF.Name())
	}
}

func (w *RotateWriter) deleteOldLogs(files []os.FileInfo) {
	for _, f := range files {
		os.Remove(f.Name())
	}
}

func (w *RotateWriter) oldLogs() ([]os.FileInfo, error) {
	var logFiles []os.FileInfo
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return nil, err
	}
	maxDays := time.Duration(int64(24*time.Hour) * int64(w.maxDays))
	lastTime := time.Now().Add(-1 * maxDays)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !strings.HasPrefix(f.Name(), w.filename) {
			continue
		}
		timeStr := strings.TrimPrefix(f.Name(), w.filename+".")
		t, err := time.Parse(time.RFC3339, timeStr)
		if err == nil {
			if t.Before(lastTime) {
				logFiles = append(logFiles, f)
			}
		}
	}
	return logFiles, nil
}

func zipper(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
	return err
}
