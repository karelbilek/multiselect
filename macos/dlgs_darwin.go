package macos

import (
	"github.com/karelbilek/multiselect/macos/cocoa"
)

func (b *FileBuilder) loadMultiple() ([]string, error) {
	return b.run(false, true)
}

func (b *FileBuilder) run(save bool, multiple bool) ([]string, error) {
	star := false
	var exts []string
	for _, filt := range b.Filters {
		for _, ext := range filt.Extensions {
			if ext == "*" {
				star = true
			} else {
				exts = append(exts, ext)
			}
		}
	}
	if star && save {
		/* OSX doesn't allow the user to switch visible file types/extensions. Also
		** NSSavePanel's allowsOtherFileTypes property has no effect for an open
		** dialog, so if "*" is a possible extension we must always show all files. */
		exts = nil
	}
	f, err := cocoa.FileDlg(save, b.Dlg.Title, exts, star, multiple, b.StartDir, b.StartFile)
	if f == nil && err == nil {
		return nil, ErrCancelled
	}
	return f, err
}
