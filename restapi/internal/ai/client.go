package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	pb "github.com/trevtemba/richrecommend/agent/v1" // path to generated Go stubs
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.ProductAgentClient
}

func NewClient(addr string) (*Client, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
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
	fmt.Println(string(productDataJson))
	resp, err := c.client.ParseProducts(ctx, &pb.ProductRequest{
		JsonInput: string(productDataJson),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			log.Printf("gRPC error: code=%v message=%v", st.Code(), st.Message())
		} else {
			log.Printf("RPC failed: %v", err)
		}
		return nil, err
	}
	log.Printf("RPC success, products received: %d", len(resp.Products))
	return resp.Products, nil
}
