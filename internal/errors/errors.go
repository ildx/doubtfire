package errors

// Error messages
const (
	ErrCopyDir           = "error copying existing destination directory"
	ErrCopyFile          = "error copying file"
	ErrCreateDir         = "error creating destination directory"
	ErrDeleteOldDir      = "error deleting old destination directory"
	ErrEmptyDir          = "the destination directory cannot be empty"
	ErrHomeDir           = "the destination directory cannot be the home directory"
	ErrMoveFile          = "error moving file"
	ErrReadDir           = "error reading directory"
	ErrResolveConflict   = "conflict detected, new destination path"
	ErrRunTUI            = "error running TUI"
	ErrSaveConfig        = "error saving configuration"
	ErrUpdateCleanupDate = "error updating last cleanup date"
)
