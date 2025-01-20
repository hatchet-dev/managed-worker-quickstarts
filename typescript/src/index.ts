import Hatchet, { Workflow, Context } from "@hatchet-dev/typescript-sdk";

const hatchet = Hatchet.init();

const parentWorkflow: Workflow = {
  id: "quickstart-typescript",
  description: "This is an example Typescript workflow.",
  steps: [
    {
      name: "step1",
      run: async (ctx) => {
        ctx.log(`This step was called at ${new Date().toISOString()}`);

        return {
          result: "This is a basic step in a DAG workflow.",
        };
      },
    },
    {
      name: "step2",
      parents: ["step1"],
      run: async (ctx) => {
        ctx.log(`This step was called at ${new Date().toISOString()}`);

        await ctx.spawnWorkflow("quickstart-child-workflow", {}).result();

        return {
          result: "This is a step which spawned a child workflow.",
        };
      },
    },
  ],
};

const childWorkflow: Workflow = {
  id: "quickstart-child-typescript",
  description:
    "This is an example Typescript child workflow. This gets spawned by the parent workflow.",
  steps: [
    {
      name: "child-step1",
      run: async (ctx) => {
        ctx.log(`This step was called at ${new Date().toISOString()}`);

        return {
          result: "This is a basic step in the child workflow.",
        };
      },
    },
  ],
};

async function main() {
  const worker = await hatchet.worker("my-worker");
  await worker.registerWorkflow(parentWorkflow);
  await worker.registerWorkflow(childWorkflow);
  worker.start();
}

main();
