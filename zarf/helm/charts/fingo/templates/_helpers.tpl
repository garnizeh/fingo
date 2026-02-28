{{/*
Common labels
*/}}
{{- define "fingo.labels" -}}
app: fingo
app.kubernetes.io/name: fingo
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "fingo.selectorLabels" -}}
app: fingo
app.kubernetes.io/name: fingo
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Database environment variables for FinGo service
*/}}
{{- define "fingo.dbEnvVars" -}}
- name: FINGO_DB_USER
  valueFrom:
    configMapKeyRef:
      name: fingo-config
      key: db_user
- name: FINGO_DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: fingo-secret
      key: db_password
- name: FINGO_DB_HOST_PORT
  valueFrom:
    configMapKeyRef:
      name: fingo-config
      key: db_hostport
- name: FINGO_DB_DISABLE_TLS
  valueFrom:
    configMapKeyRef:
      name: fingo-config
      key: db_disabletls
{{- end -}}

{{/*
Kubernetes metadata environment variables
*/}}
{{- define "fingo.k8sEnvVars" -}}
- name: KUBERNETES_NAMESPACE
  valueFrom:
    fieldRef:
      fieldPath: metadata.namespace
- name: KUBERNETES_NAME
  valueFrom:
    fieldRef:
      fieldPath: metadata.name
- name: KUBERNETES_POD_IP
  valueFrom:
    fieldRef:
      fieldPath: status.podIP
- name: KUBERNETES_NODE_NAME
  valueFrom:
    fieldRef:
      fieldPath: spec.nodeName
{{- end -}}
