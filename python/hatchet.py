from hatchet_sdk import Hatchet, Context

hatchet = Hatchet(debug=True)
 
@hatchet.workflow(name="first-workflow")
class QuickstartWorkflow:
    @hatchet.step(name="step1")
    def step1(self, context: Context):
        context.log("Called step1")

        return {
            "result": "This is a basic step in a DAG workflow."
        }
    
    @hatchet.step(name="step2", parents=["step1"])
    async def step2(self, context: Context):
        context.log("Called step2")

        res = await context.aio.spawn_workflow("quickstart-child-python", {})
        await res.result()

        return {
            "result": "This is a step which spawned a child workflow."
        }
    
@hatchet.workflow(name="quickstart-child-python")
class QuickstartChildWorkflow:
    @hatchet.step(name="child-step1")
    def step1(self, context):
        context.log("Called step1")

        return {
            "result": "This is a basic step in a DAG workflow."
        }

 
if __name__ == "__main__":
    worker = hatchet.worker('first-worker')
    worker.register_workflow(QuickstartWorkflow())
    worker.register_workflow(QuickstartChildWorkflow())
 
    worker.start()