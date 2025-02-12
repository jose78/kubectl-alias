# kubectl-alias

## Description

kubectl-alias is a kubectl plugin that allows customizing command output using SQL, enabling JOIN operations between different types of Kubernetes objects. This plugin helps users perform advanced queries to obtain detailed and structured information about their cluster.

## Installation

There a lot of ways to install the SW:

* Download directly from the Relerase
To install kubectl-alias, simply download it and add it to your $PATH. Ensure the binary has execution permissions:

```shell
chmod +x kubectl-alias
mv kubectl-alias /usr/local/bin/
```

* Clone the repo and build your own executable:
```shell 
go mod tidy
go build -o kubectl-alias
mv kubectl-alias /usr/local/bin/
```

## Configuration

The plugin requires the KUBEALIAS environment variable to point to a configuration file with the following YAML structure:

```yaml
---
version: v1 $ The current versión of the api.
aliases:
  PUT_HERE_YOUR_SUBCOMMAND:
    short: Display the description of each pod and linked service.
    long: This query retrieves all information from the event table, providing a comprehensive view of system events. It can be used to monitor and troubleshoot the cluster by showing detailed information about pod and service-related events.
    sql: Select to be executed

  PUT_HERE_YOUR_ANOTHER_SUBCOMMAND:
    short: Display the description of each pod and linked service.
    long: This query retrieves all information from the event table, providing a comprehensive view of system events. It can be used to monitor and troubleshoot the cluster by showing detailed information about pod and service-related events.
    sql: Select to be executed
    
```

Save this file in an accessible location and set the environment variable:

```shell
export KUBEALIAS=/path/to/file.yml

```

## Usage

Once installed and configured, you can run the plugin like any other kubectl command:

```shell

kubectl alias PUT_HERE_YOUR_SUBCOMMANDS

```

## Example

In my case I have now only one alias called simple_pod.

```yaml
---
version: v1
aliases:
  simple_pod:
    short: Show the description of each pod and linked service.
    args:
         - POD_NAME
    long: |
        This complex query joins multiple Kubernetes resources to provide a detailed view of the pod ecosystem. It returns:
            * The pod's namespace
            * The name of the node where the pod is running
            * The name of the pod 
            * Service info
    sql: |
          SELECT  p.metadata.namespace as NAMESPACE , n.metadata.name as NODE_NAME , p.metadata.name POD_NAME, 'service ' || s.metadata.name || ':'  ||  s.spec.ports[0].port || ' -> ' || s.spec.ports[0].targetPort as SVC
          FROM 
              pod p, svc s, deploy d , node n
          where p.metadata.labels.app = s.spec.selector.app 
          and d.metadata.labels.app = s.spec.selector.app  
          and  p.metadata.name like '%POD_NAME%'
```

Executing the next command: 
```yaml
kubectl alias simple_pod test
```
I have the next result

![alt text](resources/screenshot.png "Execution of alias simple_pod")


## Road Map

  * Add another outputs formats (json, yaml, csv) or custom using a go template. 
  * Add Openshift functionality.
  * Allow to use more formats to configure the KUBEALIAS like TOML, json, etc.
  * Enamble to use functions of go Tamplate.
  * Manager to add online subCommans.

## Contributing

Contributions are welcome. To report issues or suggest improvements, open an issue in the project repository.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.