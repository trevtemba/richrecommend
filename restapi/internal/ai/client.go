package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	pb "github.com/trevtemba/richrecommend/agent/v1" // path to generated Go stubs
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.ProductAgentClient
}

func NewClient(addr string) (*Client, error) {
	var opts []grpc.DialOption
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: pb.NewProductAgentClient(conn),
	}, nil
}

func (c *Client) GetProductData(userQuery map[string]any) ([]*pb.ParsedProduct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	productDataJson, err := json.Marshal(userQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal product query: %w", err)
	}
	resp, err := c.client.ParseProducts(ctx, &pb.ProductRequest{
		JsonInput: string(productDataJson),
	})
	if err != nil {
		return nil, err
	}
	return resp.Products, nil
}
