# API

[Table of Contents](#table-of-contents)

[Jobs](#jobs)

* [List](#list)

* [Register](#register)

* [Deregister](#deregister)

* [Info](#info)

* [Dispatch](#dispatch)

[Services](#services)

* [List](#list)

* [Info](#info)


## Jobs

### List

Returns the list of registered jobs.

* **URL**

    /v1/jobs/

* **Method**

    `GET`

* **URL Params**

    None

* **Data Params**

    None

* **Sample Request**

    ```bash
    curl -X GET http://127.0.0.1:7070/v1/jobs/
    ```

* **Sample Response**

    ```json
    {
      "Jobs": [
        {
          "ID": "JobID",
          "ParentID": "",
          "Name": "JobName",
          "Datacenters": [
            "DC"
          ],
          "Type": "JobType",
          "Priority": 50,
          "Periodic": false,
          "ParameterizedJob": false,
          "Stop": false,
          "Status": "running",
          "StatusDescription": "",
          "JobSummary": {
            "JobID": "JobID",
            "Namespace": "default",
            "Summary": {
              "api": {
                "Queued": 0,
                "Complete": 0,
                "Failed": 0,
                "Running": 1,
                "Starting": 0,
                "Lost": 0
              }
            },
            "Children": {
              "Pending": 0,
              "Running": 0,
              "Dead": 0
            },
            "CreateIndex": 4861,
            "ModifyIndex": 4866
          },
          "CreateIndex": 4861,
          "ModifyIndex": 4900,
          "JobModifyIndex": 4861,
          "SubmitTime": 1581970238593228000
        }
      ]
    }
    ```

---

### Register

Register a job.

* **URL**

    /v1/jobs/

* **Method**

    `POST`

* **URL Params**

    None

* **Data Params**

    **Required:**

    HCL or JSON payload.

* **Sample Request**

    ```bash
    curl -X POST http://127.0.0.1:7070/v1/jobs/ --data "@job.nomad"
    ```

* **Sample Response**

    ```json
    {
      "JobRegisterResponse": {
        "EvalID": "",
        "EvalCreateIndex": 0,
        "JobModifyIndex": 14376,
        "Warnings": "",
        "LastIndex": 0,
        "LastContact": 0,
        "KnownLeader": false,
        "RequestTime": 0
      }
    }
    ```

---

### Deregister

Deregister a job.

* **URL**

    /v1/jobs/

* **Method**

    `DELETE`

* **URL Params**

    None

* **Data Params**

    **Required:**

    `JobID=[string]`

* **Sample Request**

    ```bash
    curl -X DELETE http://127.0.0.1:7070/v1/jobs/ --data '{"JobID": "job_id"}'
    ```

* **Sample Response**

    ```json
    {}
    ```

---

### Info

Get job details.

* **URL**

    /v1/jobs/{id}

* **Method**

    `GET`

* **URL Params**

    **Required**

    `ID=[string]` The ID of the of the job.

* **Data Params**

    None

* **Sample Request**

    ```bash
    curl -X GET http://127.0.0.1:7070/v1/jobs/{job_id}
    ```

* **Sample Response**

    ```json
    {
      "Job": {
        "Stop": false,
        "Region": "global",
        "Namespace": "default",
        "ID": "JobID",
        "ParentID": "",
        "Name": "JobID",
        "Type": "batch",
        "Priority": 50,
        "AllAtOnce": false,
        "Datacenters": [
          "DC"
        ],
        "Constraints": null,
        "Affinities": null,
        "TaskGroups": [
          {
            "Name": "servers",
            "Count": 1,
            "Constraints": null,
            "Affinities": null,
            "Tasks": [
              {
                "Name": "docker",
                "Driver": "docker",
                "User": "",
                "Config": {
                  "cpu_hard_limit": true,
                  "image": "ImageName"
                },
                "Constraints": null,
                "Affinities": null,
                "Env": {
                  "EULA": "TRUE"
                },
                "Services": [
                  {
                    "Id": "",
                    "Name": "service-${NOMAD_ALLOC_ID}",
                    "Tags": [
                      "service",
                      "docker",
                      "${NOMAD_ALLOC_NAME}"
                    ],
                    "CanaryTags": null,
                    "EnableTagOverride": false,
                    "PortLabel": "PortName",
                    "AddressMode": "auto",
                    "Checks": [
                      {
                        "Id": "",
                        "Name": "${NOMAD_JOB_NAME}:${NOMAD_ALLOC_INDEX}-tcp",
                        "Type": "tcp",
                        "Command": "",
                        "Args": null,
                        "Path": "",
                        "Protocol": "",
                        "PortLabel": "",
                        "AddressMode": "",
                        "Interval": 10000000000,
                        "Timeout": 2000000000,
                        "InitialStatus": "",
                        "TLSSkipVerify": false,
                        "Header": null,
                        "Method": "",
                        "CheckRestart": null,
                        "GRPCService": "",
                        "GRPCUseTLS": false,
                        "TaskName": ""
                      }
                    ],
                    "CheckRestart": null,
                    "Connect": null,
                    "Meta": null,
                    "CanaryMeta": null
                  }
                ],
                "Resources": {
                  "CPU": 400,
                  "MemoryMB": 1280,
                  "DiskMB": 0,
                  "Networks": [
                    {
                      "Mode": "",
                      "Device": "",
                      "CIDR": "",
                      "IP": "",
                      "MBits": 200,
                      "ReservedPorts": null,
                      "DynamicPorts": [
                        {
                          "Label": "PortName",
                          "Value": 0,
                          "To": 0
                        }
                      ]
                    }
                  ],
                  "Devices": null,
                  "IOPS": 0
                },
                "Meta": null,
                "KillTimeout": 5000000000,
                "LogConfig": {
                  "MaxFiles": 10,
                  "MaxFileSizeMB": 10
                },
                "Artifacts": [
                  {
                    "GetterSource": "${NOMAD_META_SERVER_PROPERTIES}",
                    "GetterOptions": null,
                    "GetterMode": "any",
                    "RelativeDest": "server.properties"
                  }
                ],
                "Vault": null,
                "Templates": null,
                "DispatchPayload": null,
                "VolumeMounts": null,
                "Leader": false,
                "ShutdownDelay": 0,
                "KillSignal": "",
                "Kind": ""
              }
            ],
            "Spreads": null,
            "Volumes": null,
            "RestartPolicy": {
              "Interval": 86400000000000,
              "Attempts": 3,
              "Delay": 15000000000,
              "Mode": "fail"
            },
            "ReschedulePolicy": {
              "Attempts": 1,
              "Interval": 86400000000000,
              "Delay": 5000000000,
              "DelayFunction": "constant",
              "MaxDelay": 0,
              "Unlimited": false
            },
            "EphemeralDisk": {
              "Sticky": false,
              "Migrate": false,
              "SizeMB": 300
            },
            "Update": null,
            "Migrate": null,
            "Networks": null,
            "Meta": {
              "SERVER_PROPERTIES": "Artifact_Address"
            },
            "Services": null,
            "ShutdownDelay": null
          }
        ],
        "Update": {
          "Stagger": 0,
          "MaxParallel": 0,
          "HealthCheck": "",
          "MinHealthyTime": 0,
          "HealthyDeadline": 0,
          "ProgressDeadline": 0,
          "Canary": 0,
          "AutoRevert": false,
          "AutoPromote": false
        },
        "Spreads": null,
        "Periodic": null,
        "ParameterizedJob": {
          "Payload": "forbidden",
          "MetaRequired": null,
          "MetaOptional": [
            "SERVER_PROPERTIES"
          ]
        },
        "Dispatched": false,
        "Payload": null,
        "Reschedule": null,
        "Migrate": null,
        "Meta": null,
        "ConsulToken": null,
        "VaultToken": "",
        "Status": "running",
        "StatusDescription": "",
        "Stable": false,
        "Version": 0,
        "SubmitTime": 1581975403338654000,
        "CreateIndex": 15336,
        "ModifyIndex": 15336,
        "JobModifyIndex": 15336
      }
    }
    ```

---

### Dispatch

Dispatch a parameterized job.

* **URL**

    /v1/jobs/{id}/dispatch

* **Method**

    `POST`

* **URL Params**

    **Required**

    `ID=[string]` The ID of the job.

* **Data Params**

    **Optional**

    `Meta=[meta<string|string>: nil]` Specifies arbitrary metadata to pass to the job.

    `Payload=[string]` Specifies a base64 encoded string containing the payload. This is limited to 15 KB.

* **Sample Request**

    ```bash
    curl -X POST http://127.0.0.1:7070/v1/jobs/{job_id}/dispatch -d '{}'
    ```

* **Sample Response**

    ```json
    {
      "JobDispatchResponse": {
        "DispatchedJobID": "JobID/dispatch-1581975745-d78a8eee",
        "EvalID": "13c4a9be-0c6f-1302-02a7-c2a92daffbd6",
        "EvalCreateIndex": 16033,
        "JobCreateIndex": 16032,
        "LastIndex": 0,
        "RequestTime": 0
      }
    }
    ```

## Services

### List

Returns the list of registered services.

* **URL**

    /v1/services/

* **Method**

    `GET`

* **URL Params**

    None

* **Data Params**

    None

* **Sample Request**

    ```bash
    curl -X GET http://127.0.0.1:7070/v1/services/
    ```

* **Sample Response**

    ```json
    {
      "Services": {
        "api-12e7b412-c454-8307-2b99-9989a75f19eb": [
          "api",
          "docker",
          "api.api[0]"
        ],
        "consul": [],
        "gameserver-9d879409-0807-d560-c2f0-7f87e22c505d": [
          "docker",
          "gameserver/dispatch-1581977026-667046a1.servers[0]",
          "gameserver"
        ],
        "nomad": [
          "rpc",
          "http",
          "serf"
        ],
        "nomad-client": [
          "http"
        ]
      }
    }
    ```

---

### Info

Get service details.

* **URL**

    /v1/services/{ServiceID}

* **Method**

    `GET`

* **URL Params**

    **Required:**

    `ServiceID=[string]`

* **Data Params**

    None

* **Sample Request**

    ```bash
    curl -X GET http://127.0.0.1:7070/v1/services/{service_id}
    ```

* **Sample Response**

    ```json
    {
      "Services": [
        {
          "ID": "f0dc6890-ebf3-ccdf-1a93-e8c62ae1aefd",
          "Node": "hostname.local",
          "Address": "127.0.0.1",
          "Datacenter": "dc1",
          "TaggedAddresses": {
            "lan": "127.0.0.1",
            "wan": "127.0.0.1"
          },
          "NodeMeta": {
            "consul-network-segment": ""
          },
          "ServiceID": "_nomad-task-12e7b412-c454-8307-2b99-9989a75f19eb-api-api-12e7b412-c454-8307-2b99-9989a75f19eb-api",
          "ServiceName": "api-12e7b412-c454-8307-2b99-9989a75f19eb",
          "ServiceAddress": "127.0.0.1",
          "ServiceTaggedAddresses": null,
          "ServiceTags": [
            "api",
            "docker",
            "api.api[0]"
          ],
          "ServiceMeta": {
            "external-source": "nomad"
          },
          "ServicePort": 7070,
          "ServiceWeights": {
            "Passing": 1,
            "Warning": 1
          },
          "ServiceEnableTagOverride": false,
          "ServiceProxy": {
            "MeshGateway": {},
            "Expose": {}
          },
          "CreateIndex": 803,
          "Checks": null,
          "ModifyIndex": 803
        }
      ]
    }
    ```