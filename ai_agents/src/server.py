import grpc
import parser.parser_single as parser_single
import parser.parser_batch as parser_batch
import json
import asyncio
from concurrent import futures
from grpc import aio
from agent.v1 import agent_pb2, agent_pb2_grpc 

class ProductAgentServicer(agent_pb2_grpc.ProductAgentServicer):
    async def ParseProducts(self, request, context):
        try:
            print(f"Received request: {request.json_input}")
            # # Parse input JSON
            # data = json.loads(request.json_input)
            # text_input = data.get("text") or json.dumps(data)  # fallback if raw text not present

            # Run the parsing agent
            workflow_input = parser_batch.WorkflowInput(input_as_text=request.json_input)
            parser_output = await parser_batch.run_workflow(workflow_input)

            print(f"OUTPUT FROM PARSER: {parser_output}")
            # Safely extract parsed products
            products = parser_output.get("products", [])

            grpc_products = []

            for product in products:
                # Skip incomplete or malformed entries
                if not product.get("name"):
                    continue

                retailers = []
                for r in product.get("retailers", []):
                    retailers.append(
                        agent_pb2.Retailer(
                            name=r.get("name", ""),
                            link=r.get("link", ""),
                            rating=float(r.get("rating", 0.0)),
                            price=str(r.get("price", 0.0)),
                            in_stock=bool(r.get("in_stock", False))
                        )
                    )

                grpc_products.append(
                    agent_pb2.ParsedProduct(
                        name=product.get("name", ""),
                        description=product.get("description", ""),
                        thumbnail=product.get("thumbnail", ""),
                        ingredients=product.get("ingredients", []),
                        retailers=retailers
                    )
                )
            print("Product response sent!")
            return agent_pb2.ProductResponse(products=grpc_products)
        except Exception as e:
            await context.abort(grpc.StatusCode.INTERNAL, f"Failed to parse products: {e}")

async def serve():
    server = aio.server()
    agent_pb2_grpc.add_ProductAgentServicer_to_server(ProductAgentServicer(), server)
    server.add_insecure_port('[::]:50051')  # gRPC port
    print("Python gRPC server running on port 50051...")
    await server.start()
    await server.wait_for_termination()

if __name__ == "__main__":
    asyncio.run(serve())