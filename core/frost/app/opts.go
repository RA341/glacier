package app

import (
	"context"
	"fmt"
	"io/fs"
	"os"
)

type ServerOpt func(o *Server)

func ParseOpts(srv *Server, opts ...ServerOpt) {
	for _, opt := range opts {
		opt(srv)
	}

	if srv.Ctx == nil {
		srv.Ctx = context.Background()
	}
}

func WithCtx(ctx context.Context) ServerOpt {
	return func(o *Server) {
		o.Ctx = ctx
	}
}

func LoadUIFromDir(path string) (fs.FS, error) {
	root, err := os.OpenRoot(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open UI path: %s : %w", path, err)
	}
	return root.FS(), nil
}

func WithUIFS(uiFs fs.FS) ServerOpt {
	return func(o *Server) {
		o.UIFS = uiFs
	}
}
