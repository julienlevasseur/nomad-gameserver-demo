job "minecraft" {
	datacenters = ["dc1"]

	type = "service"

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

			//artifact {
        	//	source      = "${NOMAD_META_SERVER_PROPERTIES}"
			//	//destination = "server.properties"
			//	destination = "tmp"
        	//}

		  resources {
				cpu = 8000
				memory = 1280
        		network {
        			mbits = 200
					port "minecraft" {
						//static = 25565
					}

					port "rcon" {
						//static = 25575
					}
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

			template {
				data = <<EOH
spawn-protection=16
max-tick-time=60000
query.port={{ env "NOMAD_PORT_minecraft" }}
generator-settings=
force-gamemode=false
allow-nether=true
enforce-whitelist=false
gamemode=survival
broadcast-console-to-ops=true
enable-query=false
player-idle-timeout=0
difficulty=easy
broadcast-rcon-to-ops=true
spawn-monsters=true
op-permission-level=4
pvp=true
snooper-enabled=true
level-type=default
hardcore=false
enable-command-block=false
network-compression-threshold=256
max-players=20
max-world-size=29999984
resource-pack-sha1=
function-permission-level=2
rcon.port={{ env "NOMAD_PORT_rcon" }}
server-port={{ env "NOMAD_PORT_minecraft" }}
server-ip=
spawn-npcs=true
allow-flight=false
level-name=world
view-distance=10
resource-pack=
spawn-animals=true
white-list=false
rcon.password=
generate-structures=true
online-mode=true
max-build-height=256
level-seed=
prevent-proxy-connections=false
use-native-transport=true
motd=A Minecraft Server
enable-rcon=false
				EOH
				destination = "local/server.properties"
			}
		}

		//meta {
		//	SERVER_PROPERTIES = "https://jlevasseur-artifact-storage.s3.amazonaws.com/julienlevasseur/server.properties"
		//}
	}
}
