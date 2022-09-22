package macos

import (
	"errors"
)

// ErrCancelled is an error returned when a user cancels/closes a dialog.
var ErrCancelled = errors.New("Cancelled")

// Cancelled refers to ErrCancelled.
// Deprecated: Use ErrCancelled instead.
var Cancelled = ErrCancelled

// Dlg is the common type for dialogs.
type Dlg struct {
	Title string
}

// FileFilter represents a category of files (eg. audio files, spreadsheets).
type FileFilter struct {
	Desc       string
	Extensions []string
}

// FileBuilder is used for creating file browsing dialogs.
type FileBuilder struct {
	Dlg
	StartDir  string
	StartFile string
	Filters   []FileFilter
}

// File initialises a FileBuilder using the default configuration.
func File() *FileBuilder {
	return &FileBuilder{}
}

// Title specifies the title to be used for the dialog.
func (b *FileBuilder) Title(title string) *FileBuilder {
	b.Dlg.Title = title
	return b
}

// Filter adds a category of files to the types allowed by the dialog. Multiple
// calls to Filter are cumulative - any of the provided categories will be allowed.
// By default all files can be selected.
//
// The special extension '*' allows all files to be selected when the Filter is active.
func (b *FileBuilder) Filter(desc string, extensions ...string) *FileBuilder {
	filt := FileFilter{desc, extensions}
	if len(filt.Extensions) == 0 {
		filt.Extensions = append(filt.Extensions, "*")
	}
	b.Filters = append(b.Filters, filt)
	return b
}

// SetStartDir specifies the initial directory of the dialog.
func (b *FileBuilder) SetStartDir(startDir string) *FileBuilder {
	b.StartDir = startDir
	return b
}

// SetStartFile specifies the initial file name of the dialog.
func (b *FileBuilder) SetStartFile(startFile string) *FileBuilder {
	b.StartFile = startFile
	return b
}

func (b *FileBuilder) LoadMultiple() ([]string, error) {
	return b.loadMultiple()
}
