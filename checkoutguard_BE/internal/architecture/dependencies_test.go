package architecture

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestDependencyDirection(t *testing.T) {
	err := filepath.WalkDir("..", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, spec := range file.Imports {
			dependency, _ := strconv.Unquote(spec.Path.Value)
			internal := strings.HasPrefix(dependency, "checkoutguard/internal/")
			thirdParty := strings.Contains(strings.Split(dependency, "/")[0], ".")
			switch {
			case strings.Contains(path, "/domain/"):
				if internal || thirdParty {
					t.Errorf("domain imports outward dependency: %s -> %s", path, dependency)
				}
			case strings.Contains(path, "/application/"):
				if thirdParty || (internal && dependency != "checkoutguard/internal/domain") {
					t.Errorf("application imports outward dependency: %s -> %s", path, dependency)
				}
			case strings.Contains(path, "/adapters/"):
				if strings.Contains(dependency, "/delivery/") {
					t.Errorf("adapter imports delivery: %s -> %s", path, dependency)
				}
			case strings.Contains(path, "/delivery/"):
				if strings.Contains(dependency, "/adapters/") {
					t.Errorf("delivery imports a concrete adapter: %s -> %s", path, dependency)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
