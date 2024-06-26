{{- define "webhook.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "webhook.determineFullname" -}}
{{- if contains .ChartName .ReleaseName }}
{{- .ReleaseName | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .ReleaseName .ChartName | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{- define "webhook.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- include "webhook.determineFullname" (dict "ChartName" $name "ReleaseName" .Release.Name) }}
{{- end }}
{{- end }}

{{- define "webhook.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "webhook.labels" -}}
helm.sh/chart: {{ include "webhook.chart" . }}
{{ include "webhook.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "webhook.selectorLabels" -}}
app.kubernetes.io/name: {{ include "webhook.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "webhook.selfSignedIssuer" -}}
{{ printf "%s-selfsign" (include "webhook.fullname" .) }}
{{- end -}}

{{- define "webhook.rootCAIssuer" -}}
{{ printf "%s-ca" (include "webhook.fullname" .) }}
{{- end -}}

{{- define "webhook.rootCACertificate" -}}
{{ printf "%s-ca" (include "webhook.fullname" .) }}
{{- end -}}

{{- define "webhook.servingCertificate" -}}
{{ printf "%s-webhook-tls" (include "webhook.fullname" .) }}
{{- end -}}
