import grpc
from concurrent import futures
from agent.v1 import agent_pb2, agent_pb2_grpc 

class ProductAgentServicer(agent_pb2_grpc.ProductAgentServicer):
    def ParseProducts(self, request, context):
        print(f"Received request: {request.json_input}")
        # Example dummy recommendation
        return agent_pb2.ProductResponse(
            products=[
                agent_pb2.ParsedProduct(
                    name="SheaMoisture Curl Enhancing Smoothie",
                    description="Moisturizing styling cream for curly hair",
                    thumbnail="https://example.com/shea.jpg",
                    ingredients=["Shea Butter", "Coconut Oil"],
                    retailers=[
                        agent_pb2.Retailer(
                            name="Target",
                            link="https://target.com/product123",
                            rating=4.7,
                            price="$9.99",
                            in_stock=True
                        )
                    ]
                ),
                agent_pb2.ParsedProduct(
                    name="Cantu Leave-In Conditioner",
                    description="Repair cream for natural hair",
                    thumbnail="https://example.com/cantu.jpg",
                    ingredients=["Shea Butter", "Jojoba Oil"],
                    retailers=[
                        agent_pb2.Retailer(
                            name="Walmart",
                            link="https://walmart.com/product456",
                            rating=4.6,
                            price="$5.49",
                            in_stock=True
                        )
                    ]
                )
            ]
        )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    agent_pb2_grpc.add_ProductAgentServicer_to_server(ProductAgentServicer(), server)
    server.add_insecure_port('[::]:50051')  # gRPC port
    print("Python gRPC server running on port 50051...")
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()