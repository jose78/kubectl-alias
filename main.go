package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"log"
	"os"
)

var (
	dataK8s string = `
{

"pod": [
        {
            "apiVersion": "v1",
            "kind": "Pod",
            "metadata": {
                "creationTimestamp": "2024-10-23T09:45:18Z",
                "generateName": "test-b4b6cf8f6-",
                "labels": {
                    "app": "test",
                    "pod-template-hash": "b4b6cf8f6"
                },
                "name": "test-b4b6cf8f6-nsrsw",
                "namespace": "default",
                "ownerReferences": [
                    {
                        "apiVersion": "apps/v1",
                        "blockOwnerDeletion": true,
                        "controller": true,
                        "kind": "ReplicaSet",
                        "name": "test-b4b6cf8f6",
                        "uid": "cdda94dc-5140-45f0-a2b5-755ef21802de"
                    }
                ],
                "resourceVersion": "902116",
                "uid": "fbea3800-2cf2-4216-b109-ad9fb3192bc9"
            },
            "spec": {
                "containers": [
                    {
                        "env": [
                            {
                                "name": "name",
                                "value": "test"
                            }
                        ],
                        "image": "nginx",
                        "imagePullPolicy": "Always",
                        "name": "nginx",
                        "resources": {},
                        "terminationMessagePath": "/dev/termination-log",
                        "terminationMessagePolicy": "File",
                        "volumeMounts": [
                            {
                                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                                "name": "kube-api-access-rphdr",
                                "readOnly": true
                            }
                        ]
                    }
                ],
                "dnsPolicy": "ClusterFirst",
                "enableServiceLinks": true,
                "nodeName": "kind-control-plane",
                "preemptionPolicy": "PreemptLowerPriority",
                "priority": 0,
                "restartPolicy": "Always",
                "schedulerName": "default-scheduler",
                "securityContext": {},
                "serviceAccount": "default",
                "serviceAccountName": "default",
                "terminationGracePeriodSeconds": 30,
                "tolerations": [
                    {
                        "effect": "NoExecute",
                        "key": "node.kubernetes.io/not-ready",
                        "operator": "Exists",
                        "tolerationSeconds": 300
                    },
                    {
                        "effect": "NoExecute",
                        "key": "node.kubernetes.io/unreachable",
                        "operator": "Exists",
                        "tolerationSeconds": 300
                    }
                ],
                "volumes": [
                    {
                        "name": "kube-api-access-rphdr",
                        "projected": {
                            "defaultMode": 420,
                            "sources": [
                                {
                                    "serviceAccountToken": {
                                        "expirationSeconds": 3607,
                                        "path": "token"
                                    }
                                },
                                {
                                    "configMap": {
                                        "items": [
                                            {
                                                "key": "ca.crt",
                                                "path": "ca.crt"
                                            }
                                        ],
                                        "name": "kube-root-ca.crt"
                                    }
                                },
                                {
                                    "downwardAPI": {
                                        "items": [
                                            {
                                                "fieldRef": {
                                                    "apiVersion": "v1",
                                                    "fieldPath": "metadata.namespace"
                                                },
                                                "path": "namespace"
                                            }
                                        ]
                                    }
                                }
                            ]
                        }
                    }
                ]
            },
            "status": {
                "conditions": [
                    {
                        "lastProbeTime": null,
                        "lastTransitionTime": "2024-10-23T09:45:18Z",
                        "status": "True",
                        "type": "Initialized"
                    },
                    {
                        "lastProbeTime": null,
                        "lastTransitionTime": "2024-10-30T21:17:51Z",
                        "status": "True",
                        "type": "Ready"
                    },
                    {
                        "lastProbeTime": null,
                        "lastTransitionTime": "2024-10-30T21:17:51Z",
                        "status": "True",
                        "type": "ContainersReady"
                    },
                    {
                        "lastProbeTime": null,
                        "lastTransitionTime": "2024-10-23T09:45:18Z",
                        "status": "True",
                        "type": "PodScheduled"
                    }
                ],
                "containerStatuses": [
                    {
                        "containerID": "containerd://f50e463379cbbfd75ed79ad5ad59d99794c7bb9437c3076c1ab257732b1022a5",
                        "image": "docker.io/library/nginx:latest",
                        "imageID": "docker.io/library/nginx@sha256:28402db69fec7c17e179ea87882667f1e054391138f77ffaf0c3eb388efc3ffb",
                        "lastState": {
                            "terminated": {
                                "containerID": "containerd://ec6d920c61f05618d85b202b80fd5817d5a9501a1ea465010c75c3a0dd60215d",
                                "exitCode": 255,
                                "finishedAt": "2024-10-30T21:17:42Z",
                                "reason": "Unknown",
                                "startedAt": "2024-10-23T21:49:15Z"
                            }
                        },
                        "name": "nginx",
                        "ready": true,
                        "restartCount": 2,
                        "started": true,
                        "state": {
                            "running": {
                                "startedAt": "2024-10-30T21:17:51Z"
                            }
                        }
                    }
                ],
                "hostIP": "172.18.0.2",
                "phase": "Running",
                "podIP": "10.244.0.4",
                "podIPs": [
                    {
                        "ip": "10.244.0.4"
                    }
                ],
                "qosClass": "BestEffort",
                "startTime": "2024-10-23T09:45:18Z"
            }
        }
    ],

"deploy":  [
        {
            "apiVersion": "apps/v1",
            "kind": "Deployment",
            "metadata": {
                "annotations": {
                    "deployment.kubernetes.io/revision": "2"
                },
                "creationTimestamp": "2024-10-23T06:50:36Z",
                "generation": 2,
                "labels": {
                    "app": "test",
                    "version": "1.0"
                },
                "name": "test",
                "namespace": "default",
                "resourceVersion": "65792",
                "uid": "f2154744-8908-4389-bae9-9a4bdd748f35"
            },
            "spec": {
                "progressDeadlineSeconds": 600,
                "replicas": 1,
                "revisionHistoryLimit": 10,
                "selector": {
                    "matchLabels": {
                        "app": "test"
                    }
                },
                "strategy": {
                    "rollingUpdate": {
                        "maxSurge": "25%",
                        "maxUnavailable": "25%"
                    },
                    "type": "RollingUpdate"
                },
                "template": {
                    "metadata": {
                        "creationTimestamp": null,
                        "labels": {
                            "app": "test"
                        }
                    },
                    "spec": {
                        "containers": [
                            {
                                "env": [
                                    {
                                        "name": "name",
                                        "value": "test"
                                    }
                                ],
                                "image": "nginx",
                                "imagePullPolicy": "Always",
                                "name": "nginx",
                                "resources": {},
                                "terminationMessagePath": "/dev/termination-log",
                                "terminationMessagePolicy": "File"
                            }
                        ],
                        "dnsPolicy": "ClusterFirst",
                        "restartPolicy": "Always",
                        "schedulerName": "default-scheduler",
                        "securityContext": {},
                        "terminationGracePeriodSeconds": 30
                    }
                }
            },
            "status": {
                "availableReplicas": 1,
                "conditions": [
                    {
                        "lastTransitionTime": "2024-10-23T06:51:13Z",
                        "lastUpdateTime": "2024-10-23T06:51:13Z",
                        "message": "Deployment has minimum availability.",
                        "reason": "MinimumReplicasAvailable",
                        "status": "True",
                        "type": "Available"
                    },
                    {
                        "lastTransitionTime": "2024-10-23T06:50:36Z",
                        "lastUpdateTime": "2024-10-23T09:45:20Z",
                        "message": "ReplicaSet \"test-b4b6cf8f6\" has successfully progressed.",
                        "reason": "NewReplicaSetAvailable",
                        "status": "True",
                        "type": "Progressing"
                    }
                ],
                "observedGeneration": 2,
                "readyReplicas": 1,
                "replicas": 1,
                "updatedReplicas": 1
            }
        }
    ],
    
    "replicationSet": [
        {
            "apiVersion": "apps/v1",
            "kind": "ReplicaSet",
            "metadata": {
                "annotations": {
                    "deployment.kubernetes.io/desired-replicas": "1",
                    "deployment.kubernetes.io/max-replicas": "2",
                    "deployment.kubernetes.io/revision": "1"
                },
                "creationTimestamp": "2024-10-23T06:50:36Z",
                "generation": 2,
                "labels": {
                    "app": "test",
                    "pod-template-hash": "5746d4c59f"
                },
                "name": "test-5746d4c59f",
                "namespace": "default",
                "ownerReferences": [
                    {
                        "apiVersion": "apps/v1",
                        "blockOwnerDeletion": true,
                        "controller": true,
                        "kind": "Deployment",
                        "name": "test",
                        "uid": "f2154744-8908-4389-bae9-9a4bdd748f35"
                    }
                ],
                "resourceVersion": "65791",
                "uid": "325397bc-5d4d-4057-a846-c0a00639c3a0"
            },
            "spec": {
                "replicas": 0,
                "selector": {
                    "matchLabels": {
                        "app": "test",
                        "pod-template-hash": "5746d4c59f"
                    }
                },
                "template": {
                    "metadata": {
                        "creationTimestamp": null,
                        "labels": {
                            "app": "test",
                            "pod-template-hash": "5746d4c59f"
                        }
                    },
                    "spec": {
                        "containers": [
                            {
                                "image": "nginx",
                                "imagePullPolicy": "Always",
                                "name": "nginx",
                                "resources": {},
                                "terminationMessagePath": "/dev/termination-log",
                                "terminationMessagePolicy": "File"
                            }
                        ],
                        "dnsPolicy": "ClusterFirst",
                        "restartPolicy": "Always",
                        "schedulerName": "default-scheduler",
                        "securityContext": {},
                        "terminationGracePeriodSeconds": 30
                    }
                }
            },
            "status": {
                "observedGeneration": 2,
                "replicas": 0
            }
        },
        {
            "apiVersion": "apps/v1",
            "kind": "ReplicaSet",
            "metadata": {
                "annotations": {
                    "deployment.kubernetes.io/desired-replicas": "1",
                    "deployment.kubernetes.io/max-replicas": "2",
                    "deployment.kubernetes.io/revision": "2"
                },
                "creationTimestamp": "2024-10-23T09:45:18Z",
                "generation": 1,
                "labels": {
                    "app": "test",
                    "pod-template-hash": "b4b6cf8f6"
                },
                "name": "test-b4b6cf8f6",
                "namespace": "default",
                "ownerReferences": [
                    {
                        "apiVersion": "apps/v1",
                        "blockOwnerDeletion": true,
                        "controller": true,
                        "kind": "Deployment",
                        "name": "test",
                        "uid": "f2154744-8908-4389-bae9-9a4bdd748f35"
                    }
                ],
                "resourceVersion": "65783",
                "uid": "cdda94dc-5140-45f0-a2b5-755ef21802de"
            },
            "spec": {
                "replicas": 1,
                "selector": {
                    "matchLabels": {
                        "app": "test",
                        "pod-template-hash": "b4b6cf8f6"
                    }
                },
                "template": {
                    "metadata": {
                        "creationTimestamp": null,
                        "labels": {
                            "app": "test",
                            "pod-template-hash": "b4b6cf8f6"
                        }
                    },
                    "spec": {
                        "containers": [
                            {
                                "env": [
                                    {
                                        "name": "name",
                                        "value": "test"
                                    }
                                ],
                                "image": "nginx",
                                "imagePullPolicy": "Always",
                                "name": "nginx",
                                "resources": {},
                                "terminationMessagePath": "/dev/termination-log",
                                "terminationMessagePolicy": "File"
                            }
                        ],
                        "dnsPolicy": "ClusterFirst",
                        "restartPolicy": "Always",
                        "schedulerName": "default-scheduler",
                        "securityContext": {},
                        "terminationGracePeriodSeconds": 30
                    }
                }
            },
            "status": {
                "availableReplicas": 1,
                "fullyLabeledReplicas": 1,
                "observedGeneration": 1,
                "readyReplicas": 1,
                "replicas": 1
            }
        }
    ]
   
}
    `
)

func main() {
	os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	defer sqliteDatabase.Close()                                     // Defer Closing the database
	createTable(sqliteDatabase)                                      // Create Database Tables
	insert(sqliteDatabase)
	sqlSelect := `SELECT 
    json_extract(p.data, '$.metadata.name') AS pod_name,
    json_extract(d.data, '$.metadata.name') AS deploy_name
FROM 
    pod AS p
JOIN 
    deploy AS d ON json_extract(p.data, '$.metadata.labels.app') = json_extract(d.data, '$.metadata.name')
;`

	dataSelect_1, _ := evaluateSelect(sqliteDatabase, sqlSelect)
    fmt.Println(dataSelect_1)

}

func evaluateSelect(db *sql.DB, sqlSelect string) ([][]interface{}, error) {

    rows, errSelect := db.Query(sqlSelect)
    if errSelect != nil {
        return nil, errSelect
    }
    defer rows.Close() // Ensure rows are closed even if errors occur

    columnNames, err := rows.Columns()
    if err != nil {
        return nil, fmt.Errorf("failed to get column names: %w", err)
    }

    var results [][]interface{} // Store results as a slice of slices (interface{})

    // Prepare destination variables for scanning
    dest := make([]interface{}, len(columnNames))
    for i, _ := range columnNames {
        dest[i] = new(interface{}) // Allocate memory for each interface{}
    }

    for rows.Next() {
        // Scan row values into destination variables
        if err := rows.Scan(dest...); err != nil {
            return nil, fmt.Errorf("failed to scan row: %w", err)
        }

        // Convert scanned values to desired types if necessary
        rowValues := make([]interface{}, len(columnNames))
        for _, val := range dest {
            dataVal := *val.(*interface{})
            rowValues = append(rowValues, dataVal)
        }
        results = append(results, rowValues) // Append the row values to the results slice
    }
    return results, nil
}

func createTable(db *sql.DB) {

	var k8sObj map[string][]interface{}

	elements := []string{"pod", "deploy", "replicationSet", "secret", "configMap"}
	errK8s := json.Unmarshal([]byte(dataK8s), &k8sObj)
	if errK8s != nil {
		log.Fatal(errK8s)
	}

	data := `CREATE TABLE %s(
        id INTEGER PRIMARY KEY,
        data TEXT NOT NULL
    );
`
	for _, object := range elements {
		table := fmt.Sprintf(data, object)
		statement, err := db.Prepare(table)
		if err != nil {
			log.Fatal(err.Error())
		}
		statement.Exec()
		log.Printf("%s table created\n", object)
	}
}

func insert(db *sql.DB) {

	var k8sObj map[string][]interface{}

	elements := []string{"pod", "deploy", "replicationSet", "secret", "configMap"}
	errK8s := json.Unmarshal([]byte(dataK8s), &k8sObj)
	if errK8s != nil {
		log.Fatal(errK8s)
	}

	for _, object := range elements {

		elementToIterate := k8sObj[object]
		for _, value := range elementToIterate {

			valueByte, errK8sObj := json.Marshal(value)
			if errK8sObj != nil {
				log.Fatal(errK8sObj)
			}

			valueStr := fmt.Sprintf("INSERT INTO %s(data) VALUES('%s');", object, string(valueByte))
			statement, err := db.Prepare(valueStr) // Prepare statement.
			// This is good to avoid SQL injections
			if err != nil {
				log.Fatalln(err.Error())
			}
			_, err = statement.Exec()
			if err != nil {
				log.Fatalln(err.Error())
			}
		}
	}
}
