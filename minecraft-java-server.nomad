job "minecraft-java" {
	datacenters = ["dc1"]

	type = "service"

	group "minecraft-java-servers" {

		count = 1

		volume "minecraft-java" {
          type      = "host"
          source    = "minecraft-java"
        }
		task "server" {
		  driver = "java"
		
		  config {
		    jar_path    = "local/server.jar"
		    jvm_options = ["-Xmx1024m", "-Xms1024m"]
			args        = ["nogui"]
		  }

		  artifact {
      	    mode        = "file"
      	    source      = "https://launcher.mojang.com/v1/objects/bb2b6b1aefcd70dfd1892149ac3a215f6c636b07/server.jar"
      	  }


		  template {
            data        = "eula=true"
            destination = "eula.txt"
          }

		  template {
			  data = <<EOH
#Minecraft server properties
spawn-protection=16
max-tick-time=60000
query.port=25565
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
rcon.port=25675
server-port=25665
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
			  destination = "server.properties"
		  }

//		  artifact {
//        	source      = "https://jlevasseur-artifact-storage.s3.amazonaws.com/julienlevasseur/server.properties"
//			destination = "server.properties"
//          }

		  volume_mount {
            volume      = "minecraft-java"
            destination = "/local"
          }

		  resources {
			cpu = 200
			memory = 1024
        	network {
        	  mbits = 200
        	  port "minecraft_java" {}
        	}
          }

		  service {
        	tags = ["minecraft-java"]

        	port = "minecraft_java"

			check {
              type     = "script"
              command  = "/bin/bash"
              args     = ["-c",
			    "telnet",
				"${NOMAD_IP_minecraft_java}",
				"${NOMAD_PORT_minecraft_java}"
			  ]
			  interval = "10s"
			  timeout  = "2s"
            }
		  }
		}
	}
}
