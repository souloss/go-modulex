//go:generate go run .
package main

import (
	"path/filepath"

	"github.com/souloss/go-clean-arch/entities"
	"github.com/souloss/go-clean-arch/pkg/pkgscan"
	"gorm.io/gen"
)

const (
	packageName   = "entities"
	outputPackage = "internal/query"
)

var tags = []string{"gorm"}

// genGormCode 生成 GORM 相关的代码
func genGormCode() {
	output := filepath.Join(pkgscan.FindProjectRoot(), outputPackage)
	g := gen.NewGenerator(gen.Config{
		OutPath: output,
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	// Generate basic type-safe DAO API
	for _, st := range entities.GetAllEntities() {
		g.ApplyBasic(st)
	}

	// Generate the code
	g.Execute()
}

func main() {
	genGormCode()
}
