package shell

import "context"

type printCommands struct{}

type silent struct{}

type envVars struct{}

func WithPrintingCommands(ctx context.Context) context.Context {
	return context.WithValue(ctx, printCommands{}, true)
}

func WithSilence(ctx context.Context) context.Context {
	return context.WithValue(ctx, silent{}, true)
}

func WithEnv(ctx context.Context, env map[string]string) context.Context {
	return context.WithValue(ctx, envVars{}, env)
}

func shouldPrintCommands(ctx context.Context) bool {
	p, ok := ctx.Value(printCommands{}).(bool)
	return ok && p
}

func isSilent(ctx context.Context) bool {
	s, ok := ctx.Value(silent{}).(bool)
	return ok && s
}

func getEnvVars(ctx context.Context) map[string]string {
	env, ok := ctx.Value(envVars{}).(map[string]string)
	if !ok {
		return nil
	}
	return env
}
