package domain

import "io"

type PutFileOptions struct {
	ObjectName  string
	Body        io.Reader
	Size        int64
	ContentType string
}
