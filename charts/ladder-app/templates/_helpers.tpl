{{- define "ladder-app.name" -}}{{ .Chart.Name }}{{- end }}
{{- define "ladder-app.fullname" -}}{{ .Release.Name }}{{- end }}
{{- define "ladder-app.labels" -}}
app.kubernetes.io/name: {{ include "ladder-app.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Values.image.tag | default .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
{{- end }}
{{- define "ladder-app.selectorLabels" -}}
app.kubernetes.io/name: {{ include "ladder-app.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
