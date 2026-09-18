# Kubernetes Security & Observability Operator

A Kubernetes Operator built with **Go, Kubebuilder, and controller-runtime** to manage security-monitoring configuration through a custom Kubernetes resource.

The operator introduces a `SecurityMonitor` Custom Resource that allows security-monitoring settings to be defined declaratively. The controller continuously observes the desired state and creates or updates a Kubernetes `ConfigMap` containing the configured monitoring settings.

## Architecture

```
                    Kubernetes Cluster
                           │
                           ▼
                 ┌───────────────────┐
                 │  SecurityMonitor   │
                 │    Custom Resource │
                 └─────────┬─────────┘
                           │
                           ▼
                 ┌───────────────────┐
                 │ SecurityMonitor   │
                 │    Controller     │
                 └─────────┬─────────┘
                           │
                    Reconciliation
                           │
                           ▼
                 ┌───────────────────┐
                 │  ConfigMap        │
                 │ security-monitor- │
                 │      config       │
                 └───────────────────┘
```

The controller follows the Kubernetes reconciliation pattern: it watches the `SecurityMonitor` resource and works to make the cluster state match the desired configuration.

## Features

* Custom `SecurityMonitor` Kubernetes resource
* Declarative security-monitoring configuration
* Configurable monitoring enable/disable state
* Configurable minimum severity:

  * `low`
  * `medium`
  * `high`
  * `critical`
* Configurable security events through `alertOn`
* Automatic ConfigMap creation
* Automatic ConfigMap updates when the `SecurityMonitor` changes
* Kubernetes status conditions
* RBAC permissions generated through Kubebuilder markers
* Controller logging
* Prometheus-compatible controller metrics configuration
* Unit/controller tests using Kubernetes EnvTest
* Docker-based operator image
* Kubernetes deployment manifests generated with Kustomize
* GitHub Actions workflows for testing and linting

## Custom Resource

Example:

```yaml
apiVersion: security.example.com/v1
kind: SecurityMonitor
metadata:
  name: security-monitor
spec:
  enabled: true
  severity: high
  alertOn:
    - authentication-failure
    - forbidden-request
    - container-restart
```

Apply it with:

```bash
kubectl apply -f config/samples/security_v1_securitymonitor.yaml
```

Check the resource:

```bash
kubectl get securitymonitors
```

Check the generated ConfigMap:

```bash
kubectl get configmap security-monitor-config -o yaml
```

The resulting configuration contains values such as:

```yaml
data:
  enabled: "true"
  severity: high
  alertOn: authentication-failure,forbidden-request,container-restart
```

## How Reconciliation Works

When a `SecurityMonitor` is created or modified, the controller:

1. Retrieves the `SecurityMonitor` resource.
2. Reads its desired configuration.
3. Checks whether the corresponding ConfigMap exists.
4. Creates the ConfigMap if it does not exist.
5. Updates the ConfigMap when the desired configuration changes.
6. Updates the `SecurityMonitor` status condition.
7. Continues observing the resource for future changes.

For example, changing:

```yaml
severity: high
```

to:

```yaml
severity: critical
```

causes the controller to reconcile the resource and update the ConfigMap accordingly.

## Technology Stack

| Technology         | Purpose                                  |
| ------------------ | ---------------------------------------- |
| Go                 | Operator implementation                  |
| Kubebuilder        | Kubernetes Operator scaffolding and APIs |
| controller-runtime | Controller and reconciliation framework  |
| Kubernetes         | Runtime platform                         |
| Docker             | Containerization                         |
| Kustomize          | Kubernetes manifest management           |
| EnvTest            | Controller testing                       |
| Prometheus         | Controller metrics integration           |
| GitHub Actions     | CI workflows                             |

## Project Structure

```text
kubernetes-security-operator/
├── api/
│   └── v1/
│       └── securitymonitor_types.go
├── cmd/
│   └── main.go
├── config/
│   ├── crd/
│   ├── default/
│   ├── manager/
│   ├── prometheus/
│   ├── rbac/
│   └── samples/
├── internal/
│   └── controller/
│       ├── securitymonitor_controller.go
│       └── securitymonitor_controller_test.go
├── test/
│   └── e2e/
├── Dockerfile
├── Makefile
├── PROJECT
├── go.mod
└── README.md
```

## Prerequisites

* Go 1.27+
* Docker
* kubectl
* A Kubernetes cluster

The project can be tested locally using Docker Desktop's Kubernetes cluster or another Kubernetes environment.

## Run Tests

Run the complete Go test suite:

```bash
go test ./...
```

The controller tests use Kubernetes EnvTest to test reconciliation behavior against a Kubernetes API environment.

## Run Locally

To run the controller locally against your current Kubernetes context:

```bash
make run
```

Make sure your `kubectl` context points to the intended cluster before starting the controller.

Check the current context:

```bash
kubectl config current-context
```

## Build the Operator

Build the manager binary:

```bash
make build
```

Build the Docker image:

```bash
make docker-build IMG=<registry>/kubernetes-security-operator:tag
```

Push the image to a container registry:

```bash
make docker-push IMG=<registry>/kubernetes-security-operator:tag
```

## Deploy to Kubernetes

Install the CRDs:

```bash
make install
```

Deploy the controller:

```bash
make deploy IMG=<registry>/kubernetes-security-operator:tag
```

Verify the deployment:

```bash
kubectl get pods -n projects-system
```

Create the sample `SecurityMonitor`:

```bash
kubectl apply -f config/samples/security_v1_securitymonitor.yaml
```

Verify:

```bash
kubectl get securitymonitors
kubectl get configmaps
```

## Remove the Deployment

Delete the sample resource:

```bash
kubectl delete -f config/samples/security_v1_securitymonitor.yaml
```

Uninstall the CRDs:

```bash
make uninstall
```

Remove the controller deployment:

```bash
make undeploy
```

## Observability

The operator exposes controller metrics through the Kubernetes monitoring configuration generated by Kubebuilder.

The project also includes Prometheus-related Kubernetes manifests under:

```text
config/prometheus/
```

These provide the foundation for monitoring controller-runtime metrics and can be connected to a Prometheus/Grafana observability stack in a Kubernetes environment.

## Security Design

The operator uses Kubernetes RBAC to follow the principle of granting the controller only the permissions required for its resources.

The controller has permissions for:

* `SecurityMonitor` resources
* `SecurityMonitor` status
* `SecurityMonitor` finalizers
* ConfigMaps

RBAC rules are defined using Kubebuilder annotations in the controller source and generated into Kubernetes manifests.

## CI

GitHub Actions workflows are included under:

```text
.github/workflows/
```

They provide automated project checks including:

* Go tests
* Linting
* End-to-end test workflow

## Future Enhancements

Potential extensions include:

* Integration with actual security-event sources
* Alert delivery through external notification systems
* Prometheus alert rules for security events
* Grafana security dashboards
* Loki-based security-event log collection
* Additional Kubernetes security policies
* Automated container-image publishing
* Expanded end-to-end security scenarios

## License

Copyright 2026.

Licensed under the Apache License, Version 2.0.
