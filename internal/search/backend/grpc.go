package backend

import (
	"context"
	"io"

	"github.com/sourcegraph/sourcegraph/internal/grpc"
	"github.com/sourcegraph/zoekt"
	v1 "github.com/sourcegraph/zoekt/grpc/v1"
	"github.com/sourcegraph/zoekt/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// switchableZoektGRPCClient is a zoekt.Streamer that can switch between
// gRPC and HTTP backends.
type switchableZoektGRPCClient struct {
	httpClient zoekt.Streamer
	grpcClient zoekt.Streamer
}

func (c *switchableZoektGRPCClient) StreamSearch(ctx context.Context, q query.Q, opts *zoekt.SearchOptions, sender zoekt.Sender) error {
	if grpc.IsGRPCEnabled(ctx) {
		return c.grpcClient.StreamSearch(ctx, q, opts, sender)
	} else {
		return c.httpClient.StreamSearch(ctx, q, opts, sender)
	}
}

func (c *switchableZoektGRPCClient) Search(ctx context.Context, q query.Q, opts *zoekt.SearchOptions) (*zoekt.SearchResult, error) {
	if grpc.IsGRPCEnabled(ctx) {
		return c.grpcClient.Search(ctx, q, opts)
	} else {
		return c.httpClient.Search(ctx, q, opts)
	}
}

func (c *switchableZoektGRPCClient) List(ctx context.Context, q query.Q, opts *zoekt.ListOptions) (*zoekt.RepoList, error) {
	if grpc.IsGRPCEnabled(ctx) {
		return c.grpcClient.List(ctx, q, opts)
	} else {
		return c.httpClient.List(ctx, q, opts)
	}
}

func (c *switchableZoektGRPCClient) Close() {
	c.httpClient.Close()
}

func (c *switchableZoektGRPCClient) String() string {
	return c.httpClient.String()
}

type zoektGRPCClient struct {
	endpoint string
	client   v1.WebserverServiceClient
}

var _ zoekt.Streamer = (*zoektGRPCClient)(nil)

func (z *zoektGRPCClient) StreamSearch(ctx context.Context, q query.Q, opts *zoekt.SearchOptions, sender zoekt.Sender) error {
	req := &v1.SearchRequest{
		Query: query.QToProto(q),
		Opts:  opts.ToProto(),
	}

	ss, err := z.client.StreamSearch(ctx, req)
	if err != nil {
		return convertError(err)
	}

	for {
		msg, err := ss.Recv()
		if err != nil {
			return convertError(err)
		}

		sender.Send(zoekt.SearchResultFromProto(msg))
	}
}

func (z *zoektGRPCClient) Search(ctx context.Context, q query.Q, opts *zoekt.SearchOptions) (*zoekt.SearchResult, error) {
	req := &v1.SearchRequest{
		Query: query.QToProto(q),
		Opts:  opts.ToProto(),
	}

	resp, err := z.client.Search(ctx, req)
	if err != nil {
		return nil, convertError(err)
	}

	return zoekt.SearchResultFromProto(resp), nil
}

// List lists repositories. The query `q` can only contain
// query.Repo atoms.
func (z *zoektGRPCClient) List(ctx context.Context, q query.Q, opts *zoekt.ListOptions) (*zoekt.RepoList, error) {
	req := &v1.ListRequest{
		Query: query.QToProto(q),
		Opts:  opts.ToProto(),
	}

	resp, err := z.client.List(ctx, req)
	if err != nil {
		return nil, convertError(err)
	}

	return zoekt.RepoListFromProto(resp), nil
}

func (z *zoektGRPCClient) Close()         {}
func (z *zoektGRPCClient) String() string { return z.endpoint }

func convertError(err error) error {
	if err == nil || err == io.EOF {
		return nil
	}

	if status.Code(err) == codes.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	if status.Code(err) == codes.Canceled {
		return context.Canceled
	}

	return err
}
