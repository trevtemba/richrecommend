package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	pb "github.com/trevtemba/richrecommend/agent/v2" // path to generated Go stubs
	"github.com/trevtemba/richrecommend/internal/models"
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

func (c *Client) GetProductData(normalizerCtx context.Context, userQueryBatch []map[string]map[string]any) ([]models.ProductData, error) {
	ctx, cancel := context.WithTimeout(normalizerCtx, time.Second*120)
	defer cancel()

	productDataJson, err := json.Marshal(userQueryBatch)
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
	var parserResponses []models.ProductData
	for _, parsedProduct := range resp.Products {
		var retailerSlice []models.Retailer

		var productData models.ProductData
		productData.Name = parsedProduct.Name
		productData.Description = parsedProduct.Description
		productData.Thumbnail = parsedProduct.Thumbnail
		productData.Ingredients = parsedProduct.Ingredients

		for _, retailer := range parsedProduct.Retailers {
			retailerSlice = append(retailerSlice, models.Retailer{
				Name:    retailer.Name,
				Link:    retailer.Link,
				Rating:  retailer.Rating,
				Price:   retailer.Price,
				InStock: retailer.InStock,
			})
		}
		productData.Retailers = retailerSlice
		parserResponses = append(parserResponses, productData)
	}
	return parserResponses, nil
}
