package service

import (
	"text/template"
)

const NotificationTempl = `{{- $status := .Status -}}
{{- range .OrderEvents -}}
{{- if  eq .Type "EXECUTION" -}}
*The order has been EXECUTED*:
{{.ExecOrder.Side}} {{.Amount}} *{{.ExecOrder.Symbol}}* at {{.Price}}
{{- else -}}
*The order has been CANCELED:*
{{.Order.Side}} {{.Order.Quantity}} *{{.Order.Symbol}}* at {{.Order.LimitPrice}}
Status: {{$status}}
{{- end -}}
{{- else -}}
{{- if ne $status "placed" -}}
*Placing the order FAILED.*
Status: {{$status}}
{{- end -}}
{{- end -}}`

var NotificationTemplate = template.Must(template.New("notification").Parse(NotificationTempl))
