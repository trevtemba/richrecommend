import grpc
import parser
import json
import asyncio
from concurrent import futures
from agent.v1 import agent_pb2, agent_pb2_grpc 

class ProductAgentServicer(agent_pb2_grpc.ProductAgentServicer):
    async def ParseProducts(self, request, context):
        try:
            print(f"Received request: {request.json_input}")
            # Parse input JSON
            data = json.loads(request.json_input)
            text_input = data.get("text") or json.dumps(data)  # fallback if raw text not present

            # Run the parsing agent
            workflow_input = parser.WorkflowInput(input_as_text=text_input)
            parser_output = await parser.run_workflow(workflow_input)

            # Safely extract parsed products
            parsed = parser_output["output_parsed"]
            products = parsed.get("products", [])

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
                            rating=r.get("rating", 0.0),
                            price=r.get("price", 0.0),
                            in_stock=r.get("in_stock", False)
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
            return agent_pb2.ProductResponse(products=grpc_products)
        except Exception as e:
            await context.abort(grpc.StatusCode.INTERNAL, f"Failed to parse products: {e}")

        # return agent_pb2.ProductResponse(
        #     products=[
        #         agent_pb2.ParsedProduct(
        #             name="SheaMoisture Curl Enhancing Smoothie",
        #             description="Moisturizing styling cream for curly hair",
        #             thumbnail="https://example.com/shea.jpg",
        #             ingredients=["Shea Butter", "Coconut Oil"],
        #             retailers=[
        #                 agent_pb2.Retailer(
        #                     name="Target",
        #                     link="https://target.com/product123",
        #                     rating=4.7,
        #                     price="$9.99",
        #                     in_stock=True
        #                 )
        #             ]
        #         ),
        #         agent_pb2.ParsedProduct(
        #             name="Cantu Leave-In Conditioner",
        #             description="Repair cream for natural hair",
        #             thumbnail="https://example.com/cantu.jpg",
        #             ingredients=["Shea Butter", "Jojoba Oil"],
        #             retailers=[
        #                 agent_pb2.Retailer(
        #                     name="Walmart",
        #                     link="https://walmart.com/product456",
        #                     rating=4.6,
        #                     price="$5.49",
        #                     in_stock=True
        #                 )
        #             ]
        #         )
        #     ]
        # )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    agent_pb2_grpc.add_ProductAgentServicer_to_server(ProductAgentServicer(), server)
    server.add_insecure_port('[::]:50051')  # gRPC port
    print("Python gRPC server running on port 50051...")
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    serve()