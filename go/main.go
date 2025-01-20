package main

import (
	"fmt"
	"time"

	"github.com/hatchet-dev/hatchet/pkg/client"
	"github.com/hatchet-dev/hatchet/pkg/cmdutils"
	"github.com/hatchet-dev/hatchet/pkg/worker"
)

type stepOutput struct {
	Result string `json:"result"`
}

func main() {
	c, err := client.New()

	if err != nil {
		panic(fmt.Sprintf("error creating client: %v", err))
	}

	w, err := worker.NewWorker(
		worker.WithClient(
			c,
		),
		worker.WithMaxRuns(1),
	)

	if err != nil {
		panic(fmt.Sprintf("error creating worker: %v", err))
	}

	err = w.RegisterWorkflow(
		&worker.WorkflowJob{
			Name:        "quickstart-go",
			Description: "This is an example Go workflow.",
			On:          worker.NoTrigger(),
			Steps: []*worker.WorkflowStep{
				worker.Fn(func(ctx worker.HatchetContext) (result *stepOutput, err error) {
					ctx.Log(fmt.Sprintf("This step was called at %s", time.Now().String()))

					return &stepOutput{
						Result: "This is a basic step in a DAG workflow.",
					}, nil
				},
				).SetName("step1"),
				worker.Fn(func(ctx worker.HatchetContext) (result *stepOutput, err error) {
					ctx.Log(fmt.Sprintf("This step was called at %s", time.Now().String()))

					child, err := ctx.SpawnWorkflow("quickstart-go-child", nil, nil)

					if err != nil {
						return nil, err
					}

					_, err = child.Result()

					if err != nil {
						return nil, err
					}

					return &stepOutput{
						Result: "This is a step which spawned a child workflow.",
					}, nil
				},
				).SetName("step2").AddParents("step1"),
			},
		},
	)

	if err != nil {
		panic(fmt.Sprintf("error registering workflow: %v", err))
	}

	err = w.RegisterWorkflow(
		&worker.WorkflowJob{
			Name:        "quickstart-child-go",
			Description: "This is an example Go child workflow. This gets spawned by the parent workflow.",
			On:          worker.NoTrigger(),
			Steps: []*worker.WorkflowStep{
				worker.Fn(func(ctx worker.HatchetContext) (result *stepOutput, err error) {
					ctx.Log(fmt.Sprintf("This step was called at %s", time.Now().String()))

					return &stepOutput{
						Result: "This is a basic step in the child workflow.",
					}, nil
				},
				).SetName("child-step1"),
			},
		},
	)

	if err != nil {
		panic(fmt.Sprintf("error registering workflow: %v", err))
	}

	interruptCtx, cancel := cmdutils.InterruptContextFromChan(cmdutils.InterruptChan())
	defer cancel()

	cleanup, err := w.Start()
	if err != nil {
		panic(fmt.Sprintf("error starting worker: %v", err))
	}

	<-interruptCtx.Done()
	if err := cleanup(); err != nil {
		panic(err)
	}
}
