{{/* Common labels for all Kubernetes objects */}}
{{- define "eks-spot-termination-monitor.commonLabels" -}}
app: eks-spot-termination-monitor
{{- range $key, $value := .Values.customlabels }}
{{ $key }}: {{ $value | quote }}
{{- end }}
{{- end }}