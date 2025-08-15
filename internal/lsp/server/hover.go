package server

import (
	"context"
	"fmt"

	"cuelang.org/go/internal/golangorgx/gopls/protocol"
	"cuelang.org/go/internal/lsp/cache"
)

func (s *server) Hover(
	ctx context.Context,
	params *protocol.HoverParams,
) (_ *protocol.Hover, rerr error) {
	uri := params.TextDocument.URI
	mod, err := s.workspace.FindModuleForFile(uri)
	if err != nil {
		return nil, err
	} else if mod == nil {
		return nil, fmt.Errorf("no module found for %v", uri)
	}
	pkgs, err := mod.FindPackagesOrModulesForFile(uri)
	if err != nil {
		return nil, err
	} else if len(pkgs) == 0 {
		return nil, fmt.Errorf("no pkgs found for %v", uri)
	}
	// The first package will be the "most specific". I.e. the package
	// with root at the same directory as the file itself. There's
	// definitely an argument that we should be calling Definition for
	// all packages, and merging the results. This would find
	// definitions that exist due to ancestor imports. TODO
	pkg, ok := pkgs[0].(*cache.Package)
	if !ok {
		return nil, fmt.Errorf("no pkgs found for %v", uri)
	}
	return pkg.Hover(uri, params.Position), nil
}
