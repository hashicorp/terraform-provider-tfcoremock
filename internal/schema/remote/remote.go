package main

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/terraform-provider-mock/internal/schema"
	proto5 "github.com/hashicorp/terraform-provider-mock/internal/schema/remote/plugin/tfplugin5"
	proto6 "github.com/hashicorp/terraform-provider-mock/internal/schema/remote/plugin/tfplugin6"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

// PluginDir is the directory we will install our plugins to. We install the
// plugins we need into our own directory (potentially duplicating if Terraform
// is installing the same plugins) because (a) we don't always know where
// Terraform is installing the plugins, and (b) we might cause strange behaviour
// if both Terraform itself and this plugin (the mock provider) are both reading
// and writing to the same directory potentially at the same time.
const PluginDir = ".terraform.plugin"

type Cache struct {
	Directory string

	Providers map[string]string
}

func Open(currentDirectory string) (Cache, error) {
	cache := Cache{
		Directory: path.Join(currentDirectory, PluginDir),
		Providers: make(map[string]string),
	}

	targetDirectory := path.Join(currentDirectory, PluginDir)
	if err := os.MkdirAll(targetDirectory, os.ModePerm); err != nil {
		return cache, err
	}

	if err := cache.getProviders(); err != nil {
		return cache, err
	}

	return cache, nil
}

func (cache *Cache) GetSchema(ctx context.Context, key string) (map[string]schema.Schema, map[string]schema.Schema, error) {
	if _, ok := cache.Providers[key]; !ok {
		return nil, nil, errors.New("missing provider: " + key)
	}

	config := &plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  4,
			MagicCookieKey:   "TF_PLUGIN_MAGIC_COOKIE",
			MagicCookieValue: "d602bf8f470bc67ca7faa0386276bbdd4330efaf76d1a219cb4d6991ca9872b2",
		},
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Managed:          true,
		Cmd:              exec.Command(cache.Providers[key]),
		AutoMTLS:         true,
		VersionedPlugins: map[int]plugin.PluginSet{
			5: {
				"provider": &proto5.ProviderPlugin{},
			},
			6: {
				"provider": &proto6.ProviderPlugin{},
			},
		},
	}

	client := plugin.NewClient(config)
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, nil, err
	}
	defer rpcClient.Close()
	defer client.Kill()

	raw, err := rpcClient.Dispense("provider")
	if err != nil {
		return nil, nil, err
	}

	switch ver := client.NegotiatedVersion(); ver {
	case 5:
		p := raw.(*proto5.Provider)
		return p.GetSchema(ctx)
	case 6:
		p := raw.(*proto6.Provider)
		return p.GetSchema(ctx)
	default:
		return nil, nil, errors.New(fmt.Sprintf("unrecognized version: %d", ver))
	}
}

func (cache *Cache) InstallProvider(key, target string) error {
	if _, ok := cache.Providers[key]; ok {
		return nil
	}

	resp, err := http.Get(target)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Unsuccesful HTTP request. %d: %s", resp.StatusCode, resp.Status))
	}

	archive, err := ioutil.TempFile(cache.Directory, "terraform-provider")
	if err != nil {
		return err
	}
	defer archive.Close()
	defer os.Remove(archive.Name())

	if _, err := io.Copy(archive, resp.Body); err != nil {
		return err
	}

	reader, err := zip.OpenReader(archive.Name())
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		wantPrefix := fmt.Sprintf("terraform-provider-%s", key)
		if !strings.HasPrefix(f.Name, wantPrefix) {
			continue
		}

		dst := path.Join(cache.Directory, key)
		if err := copyFromArchive(f, dst); err != nil {
			return err
		}

		if err := os.Chmod(dst, os.ModePerm); err != nil {
			return err
		}
		cache.Providers[key] = dst
	}

	return nil
}

func copyFromArchive(src *zip.File, dst string) error {
	srcF, err := src.Open()
	if err != nil {
		return err
	}
	defer srcF.Close()

	dstF, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstF.Close()

	_, err = io.Copy(dstF, srcF)
	return err
}

func (cache *Cache) getProviders() error {
	files, err := ioutil.ReadDir(cache.Directory)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		cache.Providers[f.Name()] = path.Join(cache.Directory, f.Name())
	}

	return nil
}
