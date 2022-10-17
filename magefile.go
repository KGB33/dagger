//go:build mage
// +build mage

package main

import (
	"context"
	"errors"
	"fmt"

	"dagger.io/dagger"
	"github.com/dagger/dagger/codegen/generator"
	"github.com/google/go-cmp/cmp"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

type Lint mg.Namespace

// Run all lint targets
func (t Lint) All(ctx context.Context) error {
	mg.Deps(t.Codegen)
	return nil
}

// Lint SDK code generation
func (Lint) Codegen(ctx context.Context) error {
	c, err := dagger.Connect(ctx)
	if err != nil {
		return err

	}
	defer c.Close()

	generated, err := generator.IntrospectAndGenerate(ctx, c, generator.Config{
		Package: "dagger",
	})
	if err != nil {
		return err
	}

	// grab the file from the repo
	src, err := c.
		Host().
		Workdir().
		Read().
		File("sdk/go/api.gen.go").
		Contents(ctx)
	if err != nil {
		return err
	}

	// compare the two
	diff := cmp.Diff(string(generated), src)
	if diff != "" {
		return fmt.Errorf("generated api mismatch. please run `go generate ./...`:\n%s", diff)
	}

	return nil
}

type Build mg.Namespace

// Dagger will build the dagger binary
func (Build) Dagger(ctx context.Context) error {
	c, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	builder := c.Container().From("golang:1.18.2-alpine")

	workdir := c.Host().Workdir()
	goMod, err := workdir.Read().File("go.mod").ID(ctx)
	if err != nil {
		return err
	}

	goSum, err := workdir.Read().File("go.sum").ID(ctx)
	if err != nil {
		return err
	}

	builder = builder.WithMountedFile("/app/go.mod", goMod)
	builder = builder.WithMountedFile("/app/go.sum", goSum).WithWorkdir("/app")
	builder = builder.Exec(dagger.ContainerExecOpts{
		Args: []string{"go", "mod", "download"},
	})

	src, err := workdir.Read().ID(ctx)
	if err != nil {
		return err
	}

	builder = builder.WithMountedDirectory("/app", src).WithWorkdir("/app")

	builder = builder.Exec(dagger.ContainerExecOpts{
		Args: []string{"mkdir", "/app/build"},
	})

	builder = builder.Exec(dagger.ContainerExecOpts{
		Args: []string{"go", "build", "-o", "/app/build/dagger", "/app/cmd/dagger"},
	})

	daggerBuildDir, err := builder.Directory("./build").ID(ctx)
	if err != nil {
		return err
	}

	ok, err := c.Host().Workdir().Write(ctx, daggerBuildDir, dagger.HostDirectoryWriteOpts{Path: "."})
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("HostDirectoryWrite not ok")
	}
	return nil
}
