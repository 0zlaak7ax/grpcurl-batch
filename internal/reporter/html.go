package reporter

import (
	"html/template"
	"os"
	"time"

	"grpcurl-batch/internal/formatter"
)

const htmlTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>grpcurl-batch Report</title>
<style>
  body { font-family: sans-serif; margin: 2rem; }
  h1 { color: #333; }
  .summary { margin-bottom: 1.5rem; }
  table { border-collapse: collapse; width: 100%; }
  th, td { border: 1px solid #ccc; padding: 0.5rem 1rem; text-align: left; }
  th { background: #f0f0f0; }
  .pass { color: green; font-weight: bold; }
  .fail { color: red; font-weight: bold; }
</style>
</head>
<body>
<h1>grpcurl-batch Report</h1>
<div class="summary">
  <strong>Generated:</strong> {{.Generated}}<br>
  <strong>Total:</strong> {{.Total}} &nbsp;
  <strong>Passed:</strong> {{.Passed}} &nbsp;
  <strong>Failed:</strong> {{.Failed}}
</div>
<table>
  <thead><tr><th>#</th><th>Method</th><th>Status</th><th>Attempts</th><th>Output / Error</th></tr></thead>
  <tbody>
  {{range $i, $r := .Results}}
  <tr>
    <td>{{inc $i}}</td>
    <td>{{$r.Method}}</td>
    <td class="{{if $r.Err}}fail{{else}}pass"}}>{{if $r.Err}}FAIL{{else}}PASS{{end}}</td>
    <td>{{$r.Attempts}}</td>
    <td><pre>{{if $r.Err}}{{$r.Err}}{{else}}{{$r.Output}}{{end}}</pre></td>
  </tr>
  {{end}}
  </tbody>
</table>
</body>
</html>`

type htmlData struct {
	Generated string
	Total     int
	Passed    int
	Failed    int
	Results   []formatter.Result
}

// WriteHTML renders results as an HTML report to the given file path.
func WriteHTML(path string, results []formatter.Result) error {
	passed, failed := 0, 0
	for _, r := range results {
		if r.Err != nil {
			failed++
		} else {
			passed++
		}
	}

	funcs := template.FuncMap{
		"inc": func(i int) int { return i + 1 },
	}

	tmpl, err := template.New("report").Funcs(funcs).Parse(htmlTmpl)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, htmlData{
		Generated: time.Now().Format(time.RFC1123),
		Total:     len(results),
		Passed:    passed,
		Failed:    failed,
		Results:   results,
	})
}
