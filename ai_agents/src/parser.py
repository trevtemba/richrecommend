from pydantic import BaseModel
from agents import Agent, ModelSettings, TResponseInputItem, Runner, RunConfig
from openai.types.shared.reasoning import Reasoning

class ParserSchema__RetailersItem(BaseModel):
  name: str
  link: str
  rating: float
  price: float
  in_stock: bool


class ParserSchema__ProductData(BaseModel):
  name: str
  description: str
  thumbnail: str
  ingredients: list[str]
  retailers: list[ParserSchema__RetailersItem]


class ParserSchema(BaseModel):
  productData: ParserSchema__ProductData


parser = Agent(
  name="parser",
  instructions="""You are a product parsing assistant.
You will be given a product with product data containing unstructured text or mixed attributes.
Your task is to:

Extract and return only the fields defined in the JSON schema provided below.

Do not include any additional commentary, text, or explanation — only output valid JSON.

If a field cannot be determined, set its value to null or an empty list ([]) as appropriate.

If \"buying_options\" contains \"In stock\" then set the in_stock bool in the response schema to true, otherwise, make it false.

Only include retailers from the following list: [\"Target\", \"Walmart\", \"Amazon\", \"Ulta\", \"Sephora\"], anything else should be ignored.

Ensure that all field names exactly match the schema keys.""",
  model="gpt-5-mini-2025-08-07",
  output_type=ParserSchema,
  model_settings=ModelSettings(
    store=True,
    reasoning=Reasoning(
      effort="low",
      summary="auto"
    )
  )
)


class WorkflowInput(BaseModel):
  input_as_text: str


# Main code entrypoint
async def run_workflow(workflow_input: WorkflowInput):
  state = {

  }
  workflow = workflow_input.model_dump()
  conversation_history: list[TResponseInputItem] = [
    {
      "role": "user",
      "content": [
        {
          "type": "input_text",
          "text": workflow["input_as_text"]
        }
      ]
    }
  ]
  parser_result_temp = await Runner.run(
    parser,
    input=[
      *conversation_history
    ],
    run_config=RunConfig(trace_metadata={
      "__trace_source__": "agent-builder",
      "workflow_id": "wf_68e9bd5edef48190971c7dd14fa3f78d0fd86951122e18aa"
    })
  )

  conversation_history.extend([item.to_input_item() for item in parser_result_temp.new_items])

  parser_result = {
    "output_text": parser_result_temp.final_output.json(),
    "output_parsed": parser_result_temp.final_output.model_dump()
  }
