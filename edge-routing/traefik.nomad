job "traefik" {
	datacenters = ["dc1"]

	type = "system"

	group "edge-routers" {

		count = 1

		task "docker" {

			driver = "docker"
		
			config {
				image = "traefik"

				cpu_hard_limit = true
			}

			env {
				TRAEFIK_API_INSECURE = true
				TRAEFIK_PROVIDERS_CONSULCATALOG_REFRESHINTERVAL = "5s"
				TRAEFIK_PROVIDERS_CONSULCATALOG_ENDPOINT_ADDRESS = "http://host.docker.internal:8500"
				TRAEFIK_ENTRYPOINTS_HTTP = true
				TRAEFIK_ENTRYPOINTS_HTTP_ADDRESS = ":80"
				TRAEFIK_ENTRYPOINTS_TCP = true
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
