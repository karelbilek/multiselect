package multiselect

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/harry1453/go-common-file-dialog/cfd"
	"github.com/harry1453/go-common-file-dialog/cfdutil"
	"github.com/karelbilek/multiselect/macos"
)

var ErrCancelled = fmt.Errorf("cancelled by user")

func Fileselect(title string, ext string, extDesc string) ([]string, error) {
	if runtime.GOOS == "windows" {
		var ffs []cfd.FileFilter
		if ext != "" {
			ffs = []cfd.FileFilter{{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			}}
		} else {
			ffs = []cfd.FileFilter{
				{
					DisplayName: "All Files (*.*)",
					Pattern:     "*.*",
				},
				{
					DisplayName: extDesc,
					Pattern:     "*." + ext,
				},
			}
		}
		results, err := cfdutil.ShowOpenMultipleFilesDialog(cfd.DialogConfig{
			Title:            title,
			FileFilters:      ffs,
			DefaultExtension: ext,
		})
		if errors.Is(err, cfd.ErrorCancelled) {
			return nil, ErrCancelled
		}
		return results, err
	}
	if runtime.GOOS == "darwin" {
		df := macos.File().Title(title)
		if ext != "" {
			df.Filter(extDesc, ext)
			//df.Filter("All files", "*")
		}
		results, err := df.LoadMultiple()
		if errors.Is(err, macos.ErrCancelled) {
			return nil, ErrCancelled
		}

		return results, err
	}
	return nil, fmt.Errorf("your os does not worl, sorry")
}
