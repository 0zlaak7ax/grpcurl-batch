// Package checkpoint provides persistent progress tracking across batch runs.
//
// A Store is backed by a JSON file on disk. After each request completes the
// outcome is recorded so that a subsequent run can skip requests that already
// succeeded, making large batches safe to interrupt and resume.
//
// Basic usage:
//
//	s, err := checkpoint.New("/tmp/run.checkpoint.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, req := range requests {
//		if s.Done(req.Key) {
//			continue // already succeeded in a previous run
//		}
//		err := execute(req)
//		_ = s.Record(req.Key, err == nil)
//	}
package checkpoint
