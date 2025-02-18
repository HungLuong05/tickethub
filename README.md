# Ticket Master Backend

## Problem Requirements

1. User can claim ticket synchronously, supporting 100+ concurrent request
2. User can view events and get realtime updates on the amount of tickets left

## Tech Stack

**Programming Language**: Go

**Database**: PostgreSQL, Redis

**Ingress controller**: Nginx

**Messaging**: Apache Kafka, HTTP

**Containerization**: Docker + Kubernetes

## Installation guide

Prerequisite: A Kubernetes cluster (minikube, ...)

First, install Strimzi Operator to manage the Kafka cluster:

```bash
kubectl apply -f 'https://strimzi.io/install/latest?namespace=default' -n default
```

Then, set up the infrastructure by deploying all yaml file inside ./k8s folder:

```bash
kubectl apply -f ./k8s/*.yaml
```

## Future development

- Deploy more kafka and zookeepers pods to achieve better fault tolerant
- Deploy more event microservice to serve thousands of people watching remaning tickets of events
- Write tests for correctness and performance