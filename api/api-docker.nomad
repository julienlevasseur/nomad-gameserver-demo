job "api" {
	datacenters = ["dc1"]

	type = "service"

	group "api" {

		count = 1

		task "api" {
		  driver = "docker"
		
		    config {
				image = "julienlevasseur/nomad-gameserver-demo:latest"

//				logging {
//					type = "gelf"
//
//					config {
//						gelf-address = "tcp://host.docker.internal:12201"
//						tag = "api"
//					}
//				}
			}

			env {
				NOMAD_ADDR       = "http://host.docker.internal:4646"
				CONSUL_HTTP_ADDR = "http://host.docker.internal:8500"
			}

		    resources {
				cpu = 50
				memory = 10 // minimum value accepted, but use less
                network {
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