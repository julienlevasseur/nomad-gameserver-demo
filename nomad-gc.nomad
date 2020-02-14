job "nomad-gc" {
  datacenters = ["dc1"]

  type = "batch"

  periodic {
    # Every 5 minutes
    cron             = "*/5 * * * * * *"
    prohibit_overlap = true
  }

  group "nomad-gc" {
    count = 1

    task "nomad-gc" {
      driver = "raw_exec"

      resources {
        cpu    = 20
        memory = 10

        network {
          mbits = 1
        }
      }

      config {
        command = "local/rungc.sh"
      }

      env {
        CLUSTER_NAME = "${meta.cluster_name}"
      }

      template {
        destination = "local/rungc.sh"
        perms = "644"
        data = <<EOB
#!/bin/bash
set -eu

NOMAD_ADDR=http://127.0.0.1:4646

curl  -ik --request PUT $NOMAD_ADDR/v1/system/gc

EOB
      }
    }
  }
}
