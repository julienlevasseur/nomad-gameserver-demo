job "minecraft" {
    datacenters = ["dc1"]

    type = "batch"

    parameterized {
        payload       = "forbidden"
        meta_optional = ["SERVER_PROPERTIES"]
    }

    task "minecraft-server" {
        driver = "docker"
	
        config {
            image = "julienlevasseur/minecraft-server:1.0-alpine"

            cpu_hard_limit = true

            logging {
                type = "gelf"
    			
                config {
                    gelf-address = "tcp://host.docker.internal:12201"
                    tag = "${NOMAD_JOB_NAME}"
                }
            }
        }

        env {
            EULA = "TRUE"
        }

        resources {
            cpu = 8000
            memory = 1280
            network {
                mbits = 200
                port "minecraft" {}
            }
        }

        service {
            name = "${JOB}-${NOMAD_ALLOC_ID}"
            tags = [
				"minecraft",
				"docker",
				"${NOMAD_ALLOC_NAME}",
				"traefik.enable=true",
				"traefik.tcp.routers.TCPRouter.entrypoints=TCP",
				"traefik.tcp.routers.TCPRouter.service=${JOB}"
			]
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
}
