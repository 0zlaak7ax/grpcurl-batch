// Package notify provides pluggable notification backends used by
// grpcurl-batch to broadcast batch-run outcomes after every execution.
//
// # Backends
//
//   - [Log]     – writes a summary line to any [io.Writer].
//   - [Webhook] – HTTP POST a JSON payload to a configurable endpoint.
//   - [Multi]   – fans out to multiple backends simultaneously.
//
// # Usage
//
//	m := notify.NewMulti(
//		&notify.Log{Out: os.Stderr},
//		notify.NewWebhook("https://hooks.example.com/grpcurl"),
//	)
//	_ = m.Notify(ctx, notify.Summary{Total: 10, Passed: 9, Failed: 1})
package notify
