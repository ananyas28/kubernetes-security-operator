# Kubernetes Security & Observability Operator

A Kubernetes Operator built with **Go, Kubebuilder, and controller-runtime** for managing security-monitoring configuration and exposing operator observability metrics through **Prometheus and Grafana**.

The project demonstrates Kubernetes-native controller development, custom resources, reconciliation, RBAC, testing, metrics, monitoring, and deployment.

## Architecture

```text
                    ┌──────────────────────────┐
                    │    SecurityMonitor CR    │
                    │                          │
                    │  enabled                 │
                    │  severity                │
                    │  alertOn                 │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │ SecurityMonitor           │
                    │ Controller                │
                    │                           │
                    │ Go + controller-runtime   │
                    └────────────┬─────────────┘
                                 │
                         Reconciliation
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │ Kubernetes ConfigMap      │
                    │ security-monitor-config   │
                    └──────────────────────────┘

                    ┌──────────────────────────┐
                    │ Controller Metrics        │
                    │                           │
                    │ securitymonitor_*         │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │ Prometheus ServiceMonitor │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │ Prometheus                │
                    └────────────┬─────────────┘
                                 │
                                 ▼
                    ┌──────────────────────────┐
                    │ Grafana                   │
                    └──────────────────────────┘
```

## What the Operator Does

The operator introduces a custom Kubernetes resource called `SecurityMonitor`.

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

The controller watches `SecurityMonitor` resources and reconciles their desired configuration into Kubernetes resources.

Currently, the controller creates and synchronizes a ConfigMap containing the monitoring configuration.

## Implemented Features

* Custom `SecurityMonitor` Kubernetes CRD
* Go-based Kubernetes controller
* Kubebuilder and controller-runtime
* Kubernetes reconciliation loop
* Configurable monitoring severity
* Configurable security-event types through `alertOn`
* ConfigMap reconciliation
* Kubernetes RBAC configuration
* Controller and reconciliation tests
* Prometheus custom metrics
* Prometheus `ServiceMonitor` integration
* Grafana monitoring
* Docker Desktop Kubernetes deployment
* Kubernetes-native observability

## Custom Prometheus Metrics

The operator exposes custom metrics through the controller-runtime metrics endpoint.

### Reconciliation Count

```text
securitymonitor_reconcile_total
```

Tracks the total number of `SecurityMonitor` reconciliations.

### Reconciliation Errors

```text
securitymonitor_reconcile_errors_total
```

Tracks the total number of reconciliation errors.

The metrics were verified through Prometheus and Grafana.

Example:

```promql
securitymonitor_reconcile_total
```

```promql
securitymonitor_reconcile_errors_total
```

## Monitoring Architecture

The operator exposes its metrics through a Kubernetes Service.

Prometheus discovers the operator using a `ServiceMonitor`.

```text
SecurityMonitor Controller
          │
          ▼
   Metrics Endpoint
          │
          ▼
     ServiceMonitor
          │
          ▼
      Prometheus
          │
          ▼
       Grafana
```

The monitoring stack is deployed using the Prometheus Community `kube-prometheus-stack` Helm chart.

## Technology Stack

| Technology         | Purpose                                 |
| ------------------ | --------------------------------------- |
| Go                 | Operator implementation                 |
| Kubernetes         | Container orchestration platform        |
| Kubebuilder        | Operator scaffolding and CRD generation |
| controller-runtime | Controller and reconciliation framework |
| Docker             | Container runtime                       |
| Prometheus         | Metrics collection                      |
| ServiceMonitor     | Prometheus service discovery            |
| Grafana            | Metrics visualization                   |
| Helm               | Monitoring stack deployment             |
| EnvTest            | Controller testing                      |

## Project Structure

```text
kubernetes-security-operator/
│
├── api/
│   └── v1/
│       └── securitymonitor_types.go
│
├── cmd/
│   └── main.go
│
├── internal/
│   └── controller/
│       ├── securitymonitor_controller.go
│       └── securitymonitor_controller_test.go
│
├── config/
│   ├── crd/
│   ├── rbac/
│   ├── manager/
│   ├── default/
│   ├── samples/
│   └── monitoring-servicemonitor.yaml
│
├── test/
│   └── utils/
│
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## Running the Operator

### Prerequisites

* Go
* Docker
* Kubernetes
* kubectl
* Kubebuilder

### Run Tests

```bash
go test ./...
```

### Generate Kubernetes Manifests

```bash
make manifests
```

### Generate Code

```bash
make generate
```

### Build the Container

```bash
make docker-build IMG=security-operator:dev
```

### Deploy

The generated Kubernetes manifests under `config/` can be used to deploy the operator to a Kubernetes cluster.

## Local Kubernetes Environment

The project was developed and tested using:

* Docker Desktop Kubernetes
* Kubernetes 1.34.x
* Ubuntu/WSL2 development environment

The operator was deployed to a local Kubernetes cluster and verified through:

```bash
kubectl get pods -n projects-system
```

```bash
kubectl get securitymonitors
```

```bash
kubectl get configmaps
```

## Monitoring Verification

The operator metrics endpoint was verified directly from the Kubernetes cluster.

Prometheus successfully discovered the operator through the `ServiceMonitor`.

The custom metrics were then queried successfully from Grafana.

Verified metrics include:

```text
securitymonitor_reconcile_total
securitymonitor_reconcile_errors_total
```

At the time of verification, the reconciliation error metric returned:

```text
0
```

## Testing

Controller tests are implemented using controller-runtime's testing framework and Kubernetes EnvTest.

Run:

```bash
go test ./...
```

The test suite covers controller reconciliation behavior and Kubernetes resource handling.

## Future Improvements

Planned improvements include:

* Kubernetes audit-event processing
* Security-event detection
* Alert generation for configured security events
* Prometheus alerting rules
* Grafana security-focused dashboard
* Additional controller metrics
* Integration with Kubernetes security tools
* Automated CI/CD for container builds and deployments

## Author

Ananya S.

Built as a hands-on Kubernetes, DevOps, observability, and security engineering project.
