# {{.Title}}

> {{oneline .Description}}

{{if eq .Lang "fr" -}}
Publié le {{.DisplayDate}}{{if .Updated}}, mis à jour le {{.DisplayUpdated}}{{end}} par Benoît Hervier, lead developer et architecte logiciel freelance (https://rvier.fr/).
Version HTML : {{.URL}}
{{- else -}}
Published {{.DisplayDate}}{{if .Updated}}, updated {{.DisplayUpdated}}{{end}} by Benoît Hervier, freelance lead developer and software architect (https://rvier.fr/).
HTML version: {{.URL}}
{{- end}}

{{.Source}}
