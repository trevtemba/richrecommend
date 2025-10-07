package ai

import (
	"context"
	"time"

	pb "github.com/trevtemba/richrecommend/agent/v1" // path to generated Go stubs
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.ProductAgentClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: pb.NewProductAgentClient(conn),
	}, nil
}

func (c *Client) GetRecommendation(userQuery string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := c.client.GetRecommendation(ctx, &pb.RecommendationRequest{
		UserQuery: userQuery,
	})
	if err != nil {
		return nil, err
	}
	return resp.Products, nil
}
