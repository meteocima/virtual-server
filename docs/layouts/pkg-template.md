{{`{{ useLayout(".layout.njk") }}`}}
{{`{{ title("CIMA virtual-server") }}`}}
{{`{{ subtitle("`}}{{ .Name }}{{` package") }}`}}

# [virtual-server](./index) ⟶ {{`{{ meta.subtitle }}`}}

{{ if .IsCommand  }} 
# THIS IS A COMMAND
{{ end }}

{{ .EmitSynopsis }}

{{ .EmitUsage }}

