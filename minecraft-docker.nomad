job "minecraft" {
	datacenters = ["dc1"]

	type = "batch"

	parameterized {
  		payload       = "forbidden"
		meta_optional = ["SERVER_PROPERTIES"]
  	}

	group "servers" {

		count = 1

		task "docker" {
		  driver = "docker"
		
		  config {
				image = "openhack/minecraft-server"

				cpu_hard_limit = true
			}

		  env {
				EULA = "TRUE"
			}

			artifact {
        		source      = "${NOMAD_META_SERVER_PROPERTIES}"
				destination = "server.properties"
        	}

		  resources {
				cpu = 400
				memory = 1280
        		network {
        			mbits = 200
					port "minecraft" {}
        		}
        	}

			service {
				name = "${JOB}-${NOMAD_ALLOC_ID}"
        		tags = ["minecraft", "docker", "${NOMAD_ALLOC_NAME}"]
        		port = "minecraft"

				check {
              		name     = "${NOMAD_JOB_NAME}:${NOMAD_ALLOC_INDEX}-tcp"
            		type     = "tcp"
            		interval = "10s"
            		timeout  = "2s"
            	}
			}
		}

		meta {
			SERVER_PROPERTIES = "https://jlevasseur-artifact-storage.s3.amazonaws.com/julienlevasseur/server.properties"
		}
	}
}
