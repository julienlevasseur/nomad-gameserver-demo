job "api" {
	datacenters = ["dc1"]

	type = "service"

	group "api" {

		count = 1

		task "api" {
		  driver = "raw_exec"
		
		    config {
				command = "/usr/local/bin/api"
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
				name = "${JOB}"
        		tags = [
					"api",
					"docker",
					"${NOMAD_JOB_NAME}",
					"traefik.enable=true"
				]
        		port = "api"

				check {
              		name     = "${NOMAD_JOB_NAME}-tcp"
            		type     = "tcp"
            		interval = "10s"
            		timeout  = "2s"
            	}
			}
		}
	}
}
