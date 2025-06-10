package http

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/go-git/go-git/v6/internal/transport/test"
	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/go-git/go-git/v6/storage/filesystem"
	"github.com/go-git/go-git/v6/storage/memory"
	"github.com/stretchr/testify/suite"

	fixtures "github.com/go-git/go-git-fixtures/v5"
)

func TestUploadPackSuite(t *testing.T) {
	suite.Run(t, new(UploadPackSuite))
}

type UploadPackSuite struct {
	test.UploadPackSuite
	server *http.Server
}

func (s *UploadPackSuite) SetupTest() {
	server, base, port := setupServer(s.T(), true)
	s.server = server

	s.Client = DefaultTransport

	basic := test.PrepareRepository(s.T(), fixtures.Basic().One(), base, "basic.git")
	empty := test.PrepareRepository(s.T(), fixtures.ByTag("empty").One(), base, "empty.git")

	s.Endpoint = newEndpoint(s.T(), port, "basic.git")
	s.Storer = filesystem.NewStorage(basic, nil)

	s.EmptyEndpoint = newEndpoint(s.T(), port, "empty.git")
	s.EmptyStorer = filesystem.NewStorage(empty, nil)

	s.NonExistentEndpoint = newEndpoint(s.T(), port, "non-existent.git")
	s.NonExistentStorer = memory.NewStorage()
}

func (s *UploadPackSuite) TearDownTest() {
	s.Require().NoError(s.server.Close())
}

// Overwritten, different behaviour for HTTP.
func (s *UploadPackSuite) TestAdvertisedReferencesNotExists() {
	r, err := s.Client.NewSession(s.Storer, s.NonExistentEndpoint, s.EmptyAuth)
	s.Nil(err)
	conn, err := r.Handshake(context.TODO(), transport.UploadPackService)
	s.Error(err)
	s.Nil(conn)
}

// func (s *UploadPackSuite) TestuploadPackRequestToReader() {
// 	r := &transport.FetchRequest{}
// 	r.Wants = append(r.Wants, plumbing.NewHash("d82f291cde9987322c8a0c81a325e1ba6159684c"))
// 	r.Wants = append(r.Wants, plumbing.NewHash("2b41ef280fdb67a9b250678686a0c3e03b0a9989"))
// 	r.Haves = append(r.Haves, plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5"))
//
// 	sr, err := uploadPackRequestToReader(r)
// 	s.Nil(err)
// 	b, _ := io.ReadAll(sr)
// 	s.Equal(string(b),
// 		"0032want 2b41ef280fdb67a9b250678686a0c3e03b0a9989\n"+
// 			"0032want d82f291cde9987322c8a0c81a325e1ba6159684c\n0000"+
// 			"0032have 6ecf0ef2c2dffb796033e5a02219af86ec6584e5\n"+
// 			"0009done\n",
// 	)
// }

func (s *UploadPackSuite) TestAdvertisedReferencesRedirectPath() {
	endpoint, _ := transport.NewEndpoint("https://gitlab.com/gitlab-org/gitter/webapp")

	session, err := s.Client.NewSession(s.Storer, endpoint, s.EmptyAuth)
	s.NoError(err)
	conn, err := session.Handshake(context.TODO(), transport.UploadPackService)
	s.NoError(err)
	defer func() { s.Nil(conn.Close()) }()

	info, err := conn.GetRemoteRefs(context.TODO())
	s.NoError(err)
	s.NotNil(info)

	url := conn.(*HTTPSession).ep.String()
	s.Equal("https://gitlab.com/gitlab-org/gitter/webapp.git", url)
}

func (s *UploadPackSuite) TestAdvertisedReferencesRedirectSchema() {
	endpoint, _ := transport.NewEndpoint("http://github.com/git-fixtures/basic")

	session, err := s.Client.NewSession(s.Storer, endpoint, s.EmptyAuth)
	s.NoError(err)
	conn, err := session.Handshake(context.TODO(), transport.UploadPackService)
	s.NoError(err)
	defer func() { s.Nil(conn.Close()) }()

	info, err := conn.GetRemoteRefs(context.TODO())
	s.NoError(err)
	s.NotNil(info)

	url := conn.(*HTTPSession).ep.String()
	s.Equal("https://github.com/git-fixtures/basic", url)
}

func (s *UploadPackSuite) TestAdvertisedReferencesContext() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	endpoint, _ := transport.NewEndpoint("http://github.com/git-fixtures/basic")

	session, err := s.Client.NewSession(s.Storer, endpoint, s.EmptyAuth)
	s.NoError(err)
	conn, err := session.Handshake(ctx, transport.UploadPackService)
	s.NoError(err)
	defer func() { s.Nil(conn.Close()) }()

	info, err := conn.GetRemoteRefs(ctx)
	s.NoError(err)
	s.NotNil(info)

	url := conn.(*HTTPSession).ep.String()
	s.Equal("https://github.com/git-fixtures/basic", url)
}

func (s *UploadPackSuite) TestAdvertisedReferencesContextCanceled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	endpoint, _ := transport.NewEndpoint("http://github.com/git-fixtures/basic")

	session, err := s.Client.NewSession(s.Storer, endpoint, s.EmptyAuth)
	s.NoError(err)
	conn, err := session.Handshake(ctx, transport.UploadPackService)
	s.Error(err)
	s.Nil(conn)
	s.Equal(&url.Error{Op: "Get", URL: "http://github.com/git-fixtures/basic/info/refs?service=git-upload-pack", Err: context.Canceled}, err)
}
