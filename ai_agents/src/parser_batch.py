from pydantic import BaseModel
from agents import Agent, ModelSettings, TResponseInputItem, Runner, RunConfig
from openai.types.shared.reasoning import Reasoning

class BatchParserSchema__RetailersItem(BaseModel):
  name: str
  link: str
  rating: float
  price: float
  in_stock: bool


class BatchParserSchema__ProductsItem(BaseModel):
  name: str
  description: str
  thumbnail: str
  ingredients: list[str]
  retailers: list[BatchParserSchema__RetailersItem]


class BatchParserSchema(BaseModel):
  products: list[BatchParserSchema__ProductsItem]


batch_parser = Agent(
  name="Batch Parser",
  instructions="""You are a product parsing assistant.
You will be a list of products with product data containing unstructured text or mixed attributes.
Your task is to:

Extract and return only the fields defined in the JSON schema provided below.

Do not include any additional commentary, text, or explanation — only output valid JSON.

If a field cannot be determined, set its value to null or an empty list ([]) as appropriate.

If \"buying_options\" contains \"In stock\" then set the in_stock bool in the response schema to true, otherwise, make it false.

Only include retailers from the following list: [\"Target\", \"Walmart\", \"Amazon\", \"Ulta\", \"Sephora\"], anything else should be ignored.

Ensure that all field names exactly match the schema keys.""",
  model="gpt-5",
  output_type=BatchParserSchema,
  model_settings=ModelSettings(
    store=True,
    reasoning=Reasoning(
      effort="low"
    )
  )
)


class WorkflowInput(BaseModel):
  input_as_text: str


# Main code entrypoint
async def run_workflow(workflow_input: WorkflowInput):
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
  batch_parser_result_temp = await Runner.run(
    batch_parser,
    input=[
      *conversation_history
    ],
    run_config=RunConfig(trace_metadata={
      "__trace_source__": "agent-builder",
      "workflow_id": "wf_68ec2a1b9d7c8190b4ede250c99b6d920eaf333ed1cefb00"
    })
  )

  conversation_history.extend([item.to_input_item() for item in batch_parser_result_temp.new_items])

  batch_parser_result = {
    "output_text": batch_parser_result_temp.final_output.json(),
    "output_parsed": batch_parser_result_temp.final_output.model_dump()
  }
