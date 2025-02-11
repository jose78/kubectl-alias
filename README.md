# kubectl-alias

## Description

kubectl-ckdsasda is a kubectl plugin that allows customizing command output using SQL, enabling JOIN operations between different types of Kubernetes objects. This plugin helps users perform advanced queries to obtain detailed and structured information about their cluster.

## Installation

To install kubectl-ckdsasda, simply download it and add it to your $PATH. Ensure the binary has execution permissions:

```shell
chmod +x kubectl-ckdsasda
mv kubectl-ckdsasda /usr/local/bin/
```

## Configuration

The plugin requires the KUBEALIAS environment variable to point to a configuration file with the following YAML structure:

```yaml
---
version: v1
aliases:
  events:
    short: Display the description of each pod and linked service.
    long: This query retrieves all information from the event table, providing a comprehensive view of system events. It can be used to monitor and troubleshoot the cluster by showing detailed information about pod and service-related events.
    sql: |
      SELECT e.regarding.kind AS RESOURCE_TYPE, 
             e.type AS EVENT_TYPE_STATUS, 
             e.regarding.name AS event_resource_name, 
             p.metadata.name AS pod_name, 
             p.status.containerStatuses[0].state.waiting.message AS pod_Status_msg
      FROM event e, pod p
      WHERE e.regarding.name = p.metadata.name
  pods:
    short: Show the description of each pod and linked service.
    long: |
      This complex query joins multiple Kubernetes resources to provide a detailed view of the pod ecosystem. It returns:
      - The pod's namespace, prefixed with 'PEPE -> '
      - The name of the node where the pod is running
      - The name and UID of the ReplicaSet associated with the pod
      - The pod's name
      This information is gathered by joining the pod, service, deployment, node, and ReplicaSet resources. It filters the results to show only pods that are associated with services and deployments based on label matching, and links pods to their respective nodes and ReplicaSets. This query is particularly useful for understanding the relationships between different Kubernetes resources and troubleshooting deployment issues.
    sql: |
      SELECT 'PEPE   -> ' + p.metadata.namespace AS NAMESPACE, 
             n.metadata.name AS NODE_NAME, 
             rs.metadata.name AS RS_NAME, 
             rs.metadata.uid AS RS_UID,  
             p.metadata.name AS POD_NAME
      FROM pod p, svc s, deploy d, node n, rs
      WHERE p.metadata.labels.app = s.spec.selector.app 
        AND d.metadata.labels.app = s.spec.selector.app 
        AND n.metadata.name = p.spec.nodeName 
        AND p.metadata.name LIKE '%' + rs.metadata.name + '%'
```

Save this file in an accessible location and set the environment variable:

```shell
export KUBEALIAS=/path/to/file.yml
```

## Usage

Once installed and configured, you can run the plugin like any other kubectl command:

```shell
kubectl alias events
kubectl alias pods
```

## Requirements

kubectl installed and configured to access a Kubernetes cluster.

Go installed if you want to compile from source.

Build from Source

If you want to build the binary from the source code, run:

git clone https://github.com/jose78/kubectl-ckdsasda.git
cd kubectl-ckdsasda
go build -o kubectl-alias main.go

Then, move the binary to a directory in your $PATH.

## Contributing

Contributions are welcome. To report issues or suggest improvements, open an issue in the project repository.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.