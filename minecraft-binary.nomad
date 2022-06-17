job "minecraft-binary" {
	datacenters = ["dc1"]

	type = "service"

	group "minecraft-binary-servers" {

		count = 1

		volume "minecraft-binary" {
          type      = "host"
          source    = "minecraft-binary"
        }
		task "server" {
		  driver = "raw_exec"
		
		  config {
			command = "/usr/bin/java"
    		args    = [
				"-Xmx1024M",
				"-Xms1024M",
				"-jar",
				"local/server.jar",
				"nogui",
			]
		  }

		  artifact {
      	    mode        = "file"
      	    source      = "https://launcher.mojang.com/v1/objects/bb2b6b1aefcd70dfd1892149ac3a215f6c636b07/server.jar"
      	  }


		  template {
            data        = "eula=true"
            destination = "eula.txt"
          }

		  artifact {
        	source      = "https://jlevasseur-artifact-storage.s3.amazonaws.com/julienlevasseur/server.properties"
			destination = "server.properties"
          }

		  volume_mount {
            volume      = "minecraft-binary"
            destination = "/local"
          }

		  resources {
			cpu = 200
			memory = 1024
        	network {
        	  mbits = 200
        	  port "minecraft" {}
        	}
          }

		  service {
        	tags = ["minecraft-binary"]

        	port = "minecraft"

			check {
              type     = "script"
              command  = "/bin/bash"
              args     = ["-c",
			    "telnet",
				"${NOMAD_IP_minecraft}",
				"${NOMAD_PORT_minecraft}"
			  ]
			  interval = "10s"
			  timeout  = "2s"
            }
		  }
		}
	}
}
