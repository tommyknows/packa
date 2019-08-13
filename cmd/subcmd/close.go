package subcmd

import (
	"io"

	"github.com/tommyknows/packa/pkg/output"
)

// always try to close the controller, this will save the config
// file. if an error occurred on closing but we are already returning
// a non-nil error, we just log it.
// if no err would be returned normally, we overwrite it
// use it with:
//  defer func() {
//    err = close(ctl, err)
//  }()
func close(ctl io.Closer, inerr error) (outerr error) {
	outerr = inerr
	if err := ctl.Close(); err != nil {
		if inerr != nil {
			output.Warn(err.Error())
			return nil
		}
		outerr = err
	}
	return outerr
}
