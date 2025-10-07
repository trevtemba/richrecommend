import grpc
from concurrent import futures
from agent.v1 import agent_pb2, agent_pb2_grpc 

class RecommendationService(agent_pb2_grpc.ProductAgentServicer):
    def GetRecommendation(self, request, context):
        print(f"Received request: {request.user_query}")
        # Example dummy recommendation
        return agent_pb2.RecommendationResponse(
            products=["SheaMoisture Curl Enhancing Smoothie", "Cantu Leave-In Conditioner"]
        )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    agent_pb2_grpc.add_ProductAgentServicer_to_server(RecommendationService(), server)
    server.add_insecure_port('[::]:50051')  # gRPC port
    print("Python gRPC server running on port 50051...")
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()