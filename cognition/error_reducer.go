//cognition/error_reducer.go
package cognition
type ErrorReducer interface {
    Refine(Intent) Intent
    ConfidenceAdjust(Intent) Intent
}

type AgentRuntime struct {
    interpreter IntentInterpreter
    reducer     ErrorReducer
    planner     TaskPlanner
    router      ExecutionRouter
}

intent, err := a.interpreter.Parse(input)

intent = a.reducer.Refine(intent)

if intent.Confidence < ctx.Policy.MinConfidence {
    return fmt.Errorf("intent rejected")
}