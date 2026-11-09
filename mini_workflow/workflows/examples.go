package workflows

import (
	core "mini_workflow/core"
	tasks "mini_workflow/tasks"
	"time"
)


func Demo() core.Workflow {
	var task1 = tasks.CreateSleepTask()
	task1.SetDuration(5 * time.Second)

	
	var task2 = tasks.CreateSleepTask()
	task2.SetDuration(10 * time.Second)

	var task3 = tasks.CreateSleepTask()
	task3.SetDuration(7 * time.Second)
	task3.AddNamedDependency("tasks2", task2)

	var tasks = []core.Task{task1, task2, task3}

	return core.Workflow{
		Tasks: tasks,
		CriticalTasks: []core.Task{task3},
	}

}
