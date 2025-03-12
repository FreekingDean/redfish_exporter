{{/*
Expand the name of the chart.
*/}}
{{- define "redfish_exporter.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "redfish_exporter.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "redfish_exporter.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "redfish_exporter.labels" -}}
helm.sh/chart: {{ include "redfish_exporter.chart" . }}
{{ include "redfish_exporter.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "redfish_exporter.selectorLabels" -}}
app.kubernetes.io/name: {{ include "redfish_exporter.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "redfish_exporter.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "redfish_exporter.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "redfish_exporter.volumes" -}}
{{- if .Values.redfish_exporter_config.enabled -}}
- name: config
  configMap:
    name: {{ include "redfish_exporter.fullname" . }}
{{- end -}}
{{- with .Values.volumes -}}
{{- toYaml . | nindent 0 }}
{{- end -}}
{{- end -}}

{{- define "redfish_exporter.volume_mounts" -}}
{{- if .Values.redfish_exporter_config.enabled -}}
- name: config
  mountPath: "/etc/redfish_exporter"
{{- end -}}
{{- with .Values.volumeMounts -}}
{{- toYaml . | nindent 0 }}
{{- end -}}
{{- end -}}
