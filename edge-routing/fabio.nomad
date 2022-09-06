job "fabio" {

    datacenters = ["dc1"]

    type = "system"

    group "edge-routers" {

        count = 1

        network {
            
            port "lb" {
              static = 9999
            }
            
            port "ui" {
              static = 9998
            }
        }

        task "fabio" {

            driver = "docker"

            config {

                image        = "fabiolb/fabio"
                ports        = ["lb","ui"]

                cpu_hard_limit = true
            }

            env {
                REGISTRY_CONSUL_ADDR = "${NOMAD_IP_lb}:8500"
            }

            resources {
                cpu    = 200
                memory = 128
            }

            service {
                name = "${NOMAD_JOB_NAME}"
                tags = ["fabio", "docker"]
                port = "ui"
            }
        }
    }
}