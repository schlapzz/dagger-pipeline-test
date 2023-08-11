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
	project := client.Host().Directory(".")

	goMod := client.CacheVolume("go")
	goChache := client.CacheVolume("go-cache")

	builder := client.Container(dagger.ContainerOpts{Platform: "linux/amd64"}).
		From("golang:1.20").
		WithMountedCache("/go/src", goMod).
		WithMountedCache("/root/.cache/go-build", goChache).
		WithDirectory("/src", project).
		WithWorkdir("/src").
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
