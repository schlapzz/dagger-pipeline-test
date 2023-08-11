package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	// create dagger client
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// get host directory
	goMod := client.CacheVolume("go")
	goChache := client.CacheVolume("go-cache")

	base := client.Container(dagger.ContainerOpts{Platform: "linux/amd64"}).
		From("golang:1.20").
		WithMountedCache("/go/src", goMod).
		WithMountedCache("/root/.cache/go-build", goChache).
		WithWorkdir("/src")

	// get working directory on host
	modDir := client.Host().Directory(".", dagger.HostDirectoryOpts{
		Include: []string{"go.mod", "go.sum"},
	})

	//download go modules
	builder := base.WithDirectory("/src", modDir).
		WithExec([]string{"go", "mod", "download"})

		//get working directory on host, load files, exclude build and ci directory
	source := client.Host().Directory(".", dagger.HostDirectoryOpts{
		Exclude: []string{"ci/", "build/"},
	})

	builder = builder.WithDirectory("/src", source).
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "build", "-o", "myapp1"})
	// publish binary on alpine base
	prodImage := client.Container().
		From("alpine").
		WithFile("/bin/myapp", builder.File("/src/myapp1")).
		WithEntrypoint([]string{"/bin/myapp"})

	addr, err := prodImage.Publish(ctx, "ttl.sh/soaisbvj:5m")
	if err != nil {
		panic(err)
	}

	fmt.Println(addr)
}
