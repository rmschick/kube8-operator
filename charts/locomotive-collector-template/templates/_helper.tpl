{{/*
Takes values and transforms them into the configuration that the collector expects
*/}}
{{- define "standard-collector.configuration" -}}
{{- /* Configuration Tied To Chart Values */ -}}
{{- $metadataConfiguration := dict "environment" .Values.environment -}}
{{- $cursorConfiguration := dict "firestore" (dict "projectID" "cyderes-dev") -}}
{{- if eq .Values.environment "production" -}}
{{- $_ := set $cursorConfiguration.firestore "projectID" "cyderes-prod" -}}
{{- end -}}
{{- $lifecycleConfiguration := dict "drainTimeout" "5s" "gracefulTimeout" "5s" "shutdownTimeout" "5s" -}}
{{- $entrypointConfiguration := dict -}}
{{- if .Values.introspection.enabled -}}
{{- $endpoints := list "debug" "meta" "metrics" -}}
{{- $accessLogDefaults := dict "write" false "includeBody" false -}}
{{- $accessLogConfiguration := dict "requests" $accessLogDefaults "responses" $accessLogDefaults -}}
{{- $introspectionEntrypoint := dict "port" .Values.introspection.port "host" "0.0.0.0" "readTimeout" "3s" "readHeaderTimeout" "3s" "writeTimeout" "3s" "idleTimeout" "3s" "endpoints" $endpoints "middleware" (list) "metrics" (dict "enabled" false) "accesslog" $accessLogConfiguration -}}
{{- $_ := set $entrypointConfiguration "introspection" (dict "http" $introspectionEntrypoint) -}}
{{- end -}}
{{- $serverConfiguration := dict "lifecycle" $lifecycleConfiguration "entrypoints" $entrypointConfiguration -}}
{{- $transformedConfiguration := dict "metadata" $metadataConfiguration "cursor" $cursorConfiguration "server" $serverConfiguration -}}
{{- /* Merging of Configuration Together */ -}}
{{- $configuration := merge (.Values.cms | merge dict) $transformedConfiguration (.Values.configuration | merge dict) -}}
{{- toYaml $configuration | trim -}}
{{- end -}}
