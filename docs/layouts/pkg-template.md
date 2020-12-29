{{`{{ useLayout(".layout.njk") }}`}}
{{`{{ title("CIMA wrfda-runner") }}`}}
{{`{{ subtitle("`}}{{ .Name }}{{` package") }}`}}

# wrfda-runner ⟶ {{`{{ meta.subtitle }}`}}

{{ if .IsCommand  }} 
# THIS IS A COMMAND
{{ end }}

{{ .EmitSynopsis }}

{{ .EmitUsage }}

