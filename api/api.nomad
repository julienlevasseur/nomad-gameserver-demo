job "api" {
	datacenters = ["dc1"]

	type = "service"

	group "api" {

		count = 1

		task "api" {
		  driver = "docker"
		
		    config {
				image = "nomad-gameserver-demo:develop"
			}

			env {
				NOMAD_ADDR       = "http://host.docker.internal:4646"
				CONSUL_HTTP_ADDR = "http://host.docker.internal:8500"
			}

		    resources {
				cpu = 50
				memory = 10 // minimum value accepted, but use less
        		network {
        			mbits = 1
					port "api" {
						static = 7070
					}
        		}
        	}

			service {
				name = "${JOB}-${NOMAD_ALLOC_ID}"
        		tags = ["api", "docker", "${NOMAD_ALLOC_NAME}"]
        		port = "api"

				check {
              		name     = "${NOMAD_JOB_NAME}:${NOMAD_ALLOC_INDEX}-tcp"
            		type     = "tcp"
            		interval = "10s"
            		timeout  = "2s"
            	}
			}
		}
	}
}
