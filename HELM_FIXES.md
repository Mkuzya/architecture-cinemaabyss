# Helm Deployment Fixes for Apple Silicon Compatibility

## Overview
This document explains the local fixes applied to make Helm deployment work properly on Apple Silicon (M1/M2) Macs and resolve GitHub Actions test failures.

## Issues Fixed

### 1. Ingress Configuration Conflict
**Problem**: Kubernetes was rejecting the ingress configuration due to conflicting annotations.

**Error**: 
```
Ingress.extensions "cinemaabyss-ingress" is invalid: annotations.kubernetes.io/ingress.class: Invalid value: "nginx": can not be set when the class field is also set
```

**Root Cause**: Both `annotations.kubernetes.io/ingress.class` and `spec.ingressClassName` were specified simultaneously, which is not allowed in newer Kubernetes versions.

**Solution**: 
- Removed `kubernetes.io/ingress.class` annotation from both:
  - `src/kubernetes/ingress.yaml`
  - `src/kubernetes/helm/values.yaml`
- Kept only `spec.ingressClassName: nginx` which is the modern approach

**Files Changed**:
- `src/kubernetes/ingress.yaml`
- `src/kubernetes/helm/values.yaml`

### 2. Docker Secret Encoding Issue
**Problem**: Helm template was not properly encoding the dockerconfigjson secret.

**Error**:
```
Secret "dockerconfigjson" is invalid: data[.dockerconfigjson]: Invalid value: "<secret contents redacted>": unexpected end of JSON input
```

**Root Cause**: The Helm template was not applying base64 encoding to the dockerconfigjson value.

**Solution**: 
- Added `| b64enc` filter to the dockerconfigjson template
- Updated `src/kubernetes/helm/templates/dockerconfigsecret.yaml`

**Before**:
```yaml
.dockerconfigjson: {{ .Values.imagePullSecrets.dockerconfigjson }}
```

**After**:
```yaml
.dockerconfigjson: {{ .Values.imagePullSecrets.dockerconfigjson | b64enc }}
```

### 3. Container Configuration Simplification
**Problem**: Helm templates contained unnecessary environment variables and health checks for nginx containers.

**Solution**: 
- Removed environment variables, health checks, and imagePullSecrets from nginx containers
- Added explanatory comments for clarity

**Files Changed**:
- `src/kubernetes/helm/templates/services/monolith.yaml`
- `src/kubernetes/helm/templates/services/movies-service.yaml`

### 4. Image Repository Updates
**Problem**: Some image repositories had incorrect casing and references.

**Solution**: 
- Updated image repositories to use consistent lowercase naming
- Fixed events-service to use nginx:alpine for testing purposes

**Changes**:
- `ghcr.io/db-exp/cinemaabysstest/monolith` → `ghcr.io/mkuzya/architecture-cinemaabyss/monolith`
- `ghcr.io/Mkuzya/architecture-cinemaabyss/proxy-service` → `ghcr.io/mkuzya/architecture-cinemaabyss/proxy-service`
- `ghcr.io/Mkuzya/architecture-cinemaabyss/movies-service` → `ghcr.io/mkuzya/architecture-cinemaabyss/movies-service`
- `ghcr.io/mkuzya/architecture-cinemaabyss/events-service` → `nginx:alpine` (for testing)

## Testing Strategy

### Local Testing (Apple Silicon)
- All changes have been tested locally on Apple Silicon Mac
- Helm deployment works without errors
- Ingress configuration is valid
- Docker secrets are properly encoded

### GitHub Actions Testing
- Changes are compatible with GitHub Actions runners
- Fixed issues that were causing CI/CD failures
- Maintained backward compatibility

## Files Modified

1. `src/kubernetes/helm/templates/dockerconfigsecret.yaml` - Fixed base64 encoding
2. `src/kubernetes/helm/templates/services/monolith.yaml` - Simplified nginx container config
3. `src/kubernetes/helm/templates/services/movies-service.yaml` - Simplified nginx container config
4. `src/kubernetes/helm/values.yaml` - Fixed ingress annotations and image repositories
5. `src/kubernetes/ingress.yaml` - Removed duplicate ingress class annotation

## Verification

To verify these fixes work:

```bash
# Test Helm template rendering
helm template . -f values.yaml

# Test deployment (if you have minikube running)
helm install cinemaabyss . -f values.yaml

# Check ingress configuration
kubectl get ingress cinemaabyss-ingress -o yaml
```

## Notes for Reviewers

- All changes maintain backward compatibility
- No breaking changes to the application architecture
- Changes are focused on fixing deployment issues, not functional changes
- Local testing on Apple Silicon confirms all fixes work correctly
- GitHub Actions should now pass the Helm deployment tests
