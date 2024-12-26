package pocketbase_ts_generator

import (
	"github.com/pocketbase/pocketbase"
	pb_core "github.com/pocketbase/pocketbase/core"
	"pocketbase-ts-generator/internal/cmd"
)

type GeneratorOptions struct {
	AllCollections     bool
	CollectionsInclude []string
	CollectionsExclude []string

	Output string
}

func RegisterHook(app *pocketbase.PocketBase, options *GeneratorOptions) {
	generatorFlags := &cmd.GeneratorFlags{
		AllCollections:     options.AllCollections,
		CollectionsInclude: options.CollectionsInclude,
		CollectionsExclude: options.CollectionsExclude,

		Output: options.Output,
	}

	app.OnCollectionAfterCreateSuccess().BindFunc(func(e *pb_core.CollectionEvent) error {
		_ = processFileGeneration(app, generatorFlags)

		return e.Next()
	})

	app.OnCollectionAfterUpdateSuccess().BindFunc(func(e *pb_core.CollectionEvent) error {
		_ = processFileGeneration(app, generatorFlags)

		return e.Next()
	})

	app.OnCollectionAfterDeleteSuccess().BindFunc(func(e *pb_core.CollectionEvent) error {
		_ = processFileGeneration(app, generatorFlags)

		return e.Next()
	})
}
