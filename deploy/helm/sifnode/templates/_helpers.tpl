{{- define "sifnode.name" -}}
{{- printf "%s-sifnode" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "sifnode.labels" -}}
app.kubernetes.io/name: sifnode
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "sifnode.storageClass" -}}
{{- if .Values.persistence.storageClassName -}}
{{- .Values.persistence.storageClassName -}}
{{- else if eq .Values.cloudProvider "azure" -}}
managed-csi-premium
{{- else -}}
{{- printf "%s-%s-gp3" .Release.Namespace (include "sifnode.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
