package playwrightdebug_testhook

import (
	"encoding/json"
	"fmt"

	"github.com/xhd2015/skills/playwrightdebug"
)

type runOptionsInput struct {
	ScriptPath string   `json:"script_path"`
	ScriptArgs []string `json:"script_args"`
	CacheDir   string   `json:"cache_dir"`
	SkipEnsure bool     `json:"skip_ensure"`
}

// InspectRunPlan records the node argv and cache layout for RunOptions without executing node.
func InspectRunPlan(optsJSON []byte) ([]byte, error) {
	var input runOptionsInput
	if err := json.Unmarshal(optsJSON, &input); err != nil {
		return nil, fmt.Errorf("decode run options: %w", err)
	}

	opts := playwrightdebug.RunOptions{
		ScriptPath: input.ScriptPath,
		ScriptArgs: input.ScriptArgs,
		CacheDir:   input.CacheDir,
		SkipEnsure: input.SkipEnsure,
	}
	plan, err := playwrightdebug.BuildRunPlan(opts)
	if err != nil {
		return json.Marshal(plan)
	}
	return json.Marshal(plan)
}

// InspectCLIPlan records the node argv for CLI file-mode argv without executing node.
func InspectCLIPlan(argv []string) ([]byte, error) {
	scriptPath, scriptArgs, err := playwrightdebug.ParseCLIFileRun(argv)
	if err != nil {
		plan := playwrightdebug.RunPlanResult{Error: err.Error()}
		return json.Marshal(plan)
	}

	opts := playwrightdebug.RunOptions{
		ScriptPath: scriptPath,
		ScriptArgs: scriptArgs,
		SkipEnsure: true,
	}
	plan, err := playwrightdebug.BuildRunPlan(opts)
	if err != nil {
		return json.Marshal(plan)
	}
	return json.Marshal(plan)
}