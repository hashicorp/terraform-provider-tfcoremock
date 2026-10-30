// Copyright IBM Corp. 2022, 2026
// SPDX-License-Identifier: MPL-2.0

package state_store

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/statestore"
	schema "github.com/hashicorp/terraform-plugin-framework/statestore/schema"
)

const (
	defaultStateName       = "default"
	defaultWorkspaceDir    = "terraform.tfstate.d"
	defaultWorkspaceFile   = "environment"
	defaultStateFilename   = "terraform.tfstate"
	defaultBackupExtension = ".backup"
)

type Fs struct {
	Name string

	// WorkspaceDir is the path to the folder containing data for ALL workspaces.
	// Workspaces will be created as subdirectories within this folder.
	// Defaults to `defaultWorkspaceDir` value if not set.
	WorkspaceDir string
}

type fsModel struct {
	WorkspaceDir string `tfsdk:"workspace_dir"`
}

var (
	_ statestore.StateStore              = (*Fs)(nil)
	_ statestore.StateStoreWithConfigure = (*Fs)(nil)
)

func (f *Fs) Metadata(ctx context.Context, req statestore.MetadataRequest, resp *statestore.MetadataResponse) {
	resp.TypeName = f.Name
}

func (f *Fs) Schema(ctx context.Context, req statestore.SchemaRequest, resp *statestore.SchemaResponse) {
	resp.Schema = fsStateStoreSchema()
}

func (f *Fs) Initialize(ctx context.Context, req statestore.InitializeRequest, resp *statestore.InitializeResponse) {
	data := fsModel{}

	// We could pull data from the provider's configuration here if needed: req.ProviderData.

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	resp.StateStoreData = &data // Passes the data off to the Configure function
}

func (f *Fs) Configure(ctx context.Context, req statestore.ConfigureRequest, resp *statestore.ConfigureResponse) {
	if req.StateStoreData != nil {
		data, ok := req.StateStoreData.(*fsModel)
		if !ok {
			panic(
				fmt.Sprintf(
					"The %s store's Configure method encountered an unexpected type for StateStoreData. Expected *fsModel, got %T",
					f.Name,
					data,
				),
			)
		}

		// workspace_dir attribute
		switch {
		case data.WorkspaceDir != "":
			f.WorkspaceDir = data.WorkspaceDir
		default:
			f.WorkspaceDir = defaultWorkspaceDir
		}
	}
}

func (f *Fs) GetStates(ctx context.Context, req statestore.GetStatesRequest, resp *statestore.GetStatesResponse) {
	resp.StateIDs = []string{"default"} // TODO: stop always reporting default workspace exists

	entries, err := os.ReadDir(f.WorkspaceDir)
	// no error if there's no envs configured
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		resp.Diagnostics = append(
			resp.Diagnostics, diag.NewErrorDiagnostic(
				"Failed to read custom workspace directory", fmt.Sprintf("Error: %v", err),
			),
		)
		return
	}

	var listed []string
	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == defaultStateName {
				// Avoid listing twice
				continue
			}
			listed = append(listed, filepath.Base(entry.Name()))
		}
	}

	sort.Strings(listed)
	resp.StateIDs = append(resp.StateIDs, listed...)
}

func (f *Fs) DeleteState(ctx context.Context, req statestore.DeleteStateRequest, resp *statestore.DeleteStateResponse) {
	err := os.RemoveAll(filepath.Join(f.WorkspaceDir, req.StateID))
	if err != nil {
		resp.Diagnostics = append(
			resp.Diagnostics,
			diag.NewErrorDiagnostic("Failed to delete workspace", fmt.Sprintf("Error: %v", err)),
		)
		return
	}
}

func (f *Fs) Read(ctx context.Context, req statestore.ReadRequest, resp *statestore.ReadResponse) {
	file, err := f.openWorkspaceFile(req.StateID)
	_, ok := err.(*os.PathError)
	if err != nil && !ok {
		resp.Diagnostics = append(
			resp.Diagnostics,
			diag.NewErrorDiagnostic("Failed open workspace file", fmt.Sprintf("Error: %v", err)),
		)
		return
	}
	if ok {
		// it's ok for a state file to not exist yet; return empty state
		return
	}

	defer file.Close()

	// Read the contents of the file into resp.StateBytes
	resp.StateBytes, err = io.ReadAll(file)
	if err != nil {
		resp.Diagnostics = append(
			resp.Diagnostics,
			diag.NewErrorDiagnostic("Failed to read workspace file", fmt.Sprintf("Error: %v", err)),
		)
		return
	}
}

func (f *Fs) Write(ctx context.Context, req statestore.WriteRequest, resp *statestore.WriteResponse) {
	file, err := f.createOrOpenWorkspaceFile(req.StateID)
	if err != nil {
		resp.Diagnostics = append(
			resp.Diagnostics,
			diag.NewErrorDiagnostic("Failed to create or open workspace file", fmt.Sprintf("Error: %v", err)),
		)
		return
	}
	defer file.Close()

	_, err = file.Write(req.StateBytes)
	if err != nil {
		resp.Diagnostics = append(
			resp.Diagnostics,
			diag.NewErrorDiagnostic("Failed to write state to workspace file", fmt.Sprintf("Error: %v", err)),
		)
		return
	}
}

func (f *Fs) Lock(ctx context.Context, req statestore.LockRequest, resp *statestore.LockResponse) {
	stateFilePath := f.stateFilePath(req.StateID)

	file, err := os.OpenFile(stateFilePath, os.O_RDWR, 0666)
	_, ok := err.(*os.PathError)
	if err != nil && !ok {
		resp.Diagnostics = append(
			resp.Diagnostics,
			diag.NewErrorDiagnostic("Failed to lock workspace file", fmt.Sprintf("Error when opening file: %v", err)),
		)
		return
	}
	if ok {
		// handle state file not existing yet by making an empty file
		file, err = f.createOrOpenWorkspaceFile(req.StateID)
		if err != nil {
			resp.Diagnostics = append(
				resp.Diagnostics,
				diag.NewErrorDiagnostic("Failed to create missing workspace file during locking", fmt.Sprintf("Error: %v", err)),
			)
			return
		}
	}

	defer file.Close()

	err = f.lock(file)
	if err != nil {
		resp.Diagnostics = append(
			resp.Diagnostics,
			diag.NewErrorDiagnostic("Failed to lock workspace file", fmt.Sprintf("Error when locking file: %v", err)),
		)
		return
	}
}

func (f *Fs) Unlock(ctx context.Context, req statestore.UnlockRequest, resp *statestore.UnlockResponse) {
	stateFilePath := f.stateFilePath(req.StateID)

	file, err := os.OpenFile(stateFilePath, os.O_RDWR, 0666)
	_, ok := err.(*os.PathError)
	if err != nil && !ok {
		resp.Diagnostics = append(
			resp.Diagnostics,
			diag.NewErrorDiagnostic("Failed to unlock workspace file", fmt.Sprintf("Error when opening file: %v", err)),
		)
		return
	}
	if ok {
		// handle state file not existing yet by making an empty file
		file, err = f.createOrOpenWorkspaceFile(req.StateID)
		if err != nil {
			resp.Diagnostics = append(
				resp.Diagnostics,
				diag.NewErrorDiagnostic("Failed to create missing workspace file during unlocking", fmt.Sprintf("Error: %v", err)),
			)
			return
		}
	}

	defer file.Close()

	err = f.unlock(file)
	if err != nil {
		resp.Diagnostics = append(
			resp.Diagnostics,
			diag.NewErrorDiagnostic("Failed to unlock workspace file", fmt.Sprintf("Error when unlocking file: %v", err)),
		)
		return
	}
}

func fsStateStoreSchema() schema.Schema {
	return schema.Schema{
		Description: "Filesystem-based state store for tfcoremock provider",
		Attributes: map[string]schema.Attribute{
			"workspace_dir": schema.StringAttribute{
				Description: "The path to the directory where custom workspace state is stored.",
				Optional:    true,
			},
		},
	}
}

// stateFilePath returns the path to the active state file based on the workspace name provided
func (f *Fs) stateFilePath(name string) string {
	return filepath.Join(f.WorkspaceDir, name, defaultStateFilename)
}

// openWorkspaceFile attempts to open a state file
// Calling code is responsible for handling the case where the file does not exist.
func (f *Fs) openWorkspaceFile(workspaceName string) (*os.File, error) {
	stateFilePath := f.stateFilePath(workspaceName)
	return os.OpenFile(stateFilePath, os.O_RDWR, 0666)
}

// createOrOpenWorkspaceFile will open an existing state file or create it if it does not exist.
// This includes creating any necessary directories leading up to the state file location.
func (f *Fs) createOrOpenWorkspaceFile(workspaceName string) (*os.File, error) {
	stateFilePath := f.stateFilePath(workspaceName)
	_, err := os.Stat(stateFilePath)
	if _, ok := err.(*os.PathError); ok {
		// We need to create the file, after any directories leading up to it
		if err := os.MkdirAll(filepath.Dir(stateFilePath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create directories for workspace %q: %w", workspaceName, err)
		}
		file, err := os.OpenFile(stateFilePath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return file, fmt.Errorf("failed to create state file for workspace %q: %w", workspaceName, err)
		}
		return file, nil
	}
	// File exists
	return os.OpenFile(stateFilePath, os.O_RDWR, 0666)
}
