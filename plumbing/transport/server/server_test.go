package server_test

import (
	"testing"

	"github.com/nulab/go-git/plumbing/cache"
	"github.com/nulab/go-git/plumbing/transport"
	"github.com/nulab/go-git/plumbing/transport/client"
	"github.com/nulab/go-git/plumbing/transport/server"
	"github.com/nulab/go-git/plumbing/transport/test"
	"github.com/nulab/go-git/storage/filesystem"
	"github.com/nulab/go-git/storage/memory"

	. "gopkg.in/check.v1"
	"github.com/nulab/go-git-fixtures"
)

func Test(t *testing.T) { TestingT(t) }

type BaseSuite struct {
	fixtures.Suite
	test.ReceivePackSuite

	loader       server.MapLoader
	client       transport.Transport
	clientBackup transport.Transport
	asClient     bool
}

func (s *BaseSuite) SetUpSuite(c *C) {
	s.Suite.SetUpSuite(c)
	s.loader = server.MapLoader{}
	if s.asClient {
		s.client = server.NewClient(s.loader)
	} else {
		s.client = server.NewServer(s.loader)
	}

	s.clientBackup = client.Protocols["file"]
	client.Protocols["file"] = s.client
}

func (s *BaseSuite) TearDownSuite(c *C) {
	if s.clientBackup == nil {
		delete(client.Protocols, "file")
	} else {
		client.Protocols["file"] = s.clientBackup
	}
}

func (s *BaseSuite) prepareRepositories(c *C) {
	var err error

	fs := fixtures.Basic().One().DotGit()
	s.Endpoint, err = transport.NewEndpoint(fs.Root())
	c.Assert(err, IsNil)
	s.loader[s.Endpoint.String()] = filesystem.NewStorage(fs, cache.NewObjectLRUDefault())

	s.EmptyEndpoint, err = transport.NewEndpoint("/empty.git")
	c.Assert(err, IsNil)
	s.loader[s.EmptyEndpoint.String()] = memory.NewStorage()

	s.NonExistentEndpoint, err = transport.NewEndpoint("/non-existent.git")
	c.Assert(err, IsNil)
}
