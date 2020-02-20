job "traefik" {
	datacenters = ["dc1"]

	type = "service"

	group "edge-routers" {

		count = 1

		task "docker" {
		  driver = "docker"
		
		    config {
			    image = "traefik"

			    cpu_hard_limit = true

                volumes = [
                    "/Users/julien.levasseur/go/src/github.com/julienlevasseur/nomad-gameserver-demo/edge-routing/:/etc/traefik/",
                    "/var/run/docker.sock:/var/run/docker.sock"
                ]
		    }

		    resources {
		    	cpu = 1000
		    	memory = 64
            	network {
            		mbits = 10
		    		port "traefik_http" {
		    			static = 8080
		    		}
            	}
            }

		    service {
		    	name = "${NOMAD_JOB_NAME}"
            	tags = ["traefik", "docker"]
            	port = "traefik_http"

		    	check {
               		name     = "${NOMAD_JOB_NAME}-tcp"
                	type     = "tcp"
                	interval = "10s"
                	timeout  = "2s"
               	}

                check {
                    name     = "${NOMAD_JOB_NAME}-http"
                    type     = "http"
                    protocol = "http"
                    port     = "traefik_http"
                    path     = "/dashboard"
                    interval = "5s"
                    timeout  = "2s"
                    method   = "GET"
                }


		    }
        }
	}
}
