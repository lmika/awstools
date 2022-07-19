package workspaces

import (
	"github.com/asdine/storm"
	"github.com/pkg/errors"
	"os"
)

type MetaInfo struct {
	Command string
}

type Manager struct {
	metainfo MetaInfo
}

func New(metaInfo MetaInfo) *Manager {
	return &Manager{metainfo: metaInfo}
}

func (m *Manager) Open(filename string) (*Workspace, error) {
	db, err := storm.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open workspace at %v", filename)
	}
	return &Workspace{db: db}, nil
}

func (m *Manager) CreateTemp() (*Workspace, error) {
	workspaceFile, err := os.CreateTemp("", m.metainfo.Command+"*.workspace")
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create workspace file")
	}
	workspaceFile.Close() // We just need the filename

	return m.Open(workspaceFile.Name())
}
