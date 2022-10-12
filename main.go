package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kong"
	"github.com/google/pprof/profile"
	"github.com/pkg/browser"
	sharev1alpha1 "go.buf.build/grpc/go/parca-dev/parca/parca/share/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Context struct {
	Debug  bool
	Server string
}

type UploadCmd struct {
	Description string `help:"Description of the profile." short:"d"`
	Path        string `arg:"" name:"path" help:"Path to pprof profile." type:"path"`
}

func (u *UploadCmd) Run(ctx *Context) error {
	return runProfileUpload(ctx.Server, u.Path, u.Description)
}

var cli struct {
	Debug  bool   `help:"Enable debug mode."`
	Server string `help:"Server to upload to." default:"api.pprof.me:443"`

	Upload UploadCmd `cmd:"upload" description:"Upload a profile"`
}

func main() {
	ctx := kong.Parse(&cli)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run(&Context{
		Server: cli.Server,
		Debug:  cli.Debug,
	})
	ctx.FatalIfErrorf(err)
}

func runProfileUpload(server, target, description string) error {
	i, err := os.Stat(target)
	if err != nil {
		// TODO(brancz): pprof allows this to be a URL, maybe we should too?
		return err
	}

	// The size is also validated on the server side, but the error is not as clear.
	if i.Size() > 2*1024*1024 {
		return errors.New("profile is too large, max 2MB allowed")
	}

	if i.IsDir() {
		return errors.New("target is a directory not a pprof profile")
	}

	if description == "" {
		description = filepath.Base(target)
	}

	f, err := os.Open(target)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("close file: %w", err)
	}

	// Parsed server-side as well but doesn't hurt to do before we upload.
	p, err := profile.ParseData(data)
	if err != nil {
		return fmt.Errorf("parse profile: %w", err)
	}

	// Validated server-side as well but doesn't hurt to do before we upload.
	if err := p.CheckValid(); err != nil {
		return fmt.Errorf("validate profile: %w", err)
	}

	conn, err := grpc.Dial(
		server,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
	)
	if err != nil {
		return fmt.Errorf("dial server: %w", err)
	}

	client := sharev1alpha1.NewShareServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Upload(ctx, &sharev1alpha1.UploadRequest{
		Description: description,
		Profile:     data,
	})
	if err != nil {
		return fmt.Errorf("upload profile: %w", err)
	}

	fmt.Println(resp.Link)
	return browser.OpenURL(resp.Link)
}
