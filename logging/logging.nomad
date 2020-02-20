job "logging" {
    datacenters = ["dc1"]

    type = "service"

    task "mongo" {
        driver = "docker"

        config {
            image = "mongo:3"
        }

        resources {
			cpu = 100
			memory = 256
        	network {
        		mbits = 1
				port "mongo" {
					static = 27017
				}
        	}
        }

        service {
			name = "${JOB}-${NOMAD_ALLOC_ID}"
        	tags = ["mongo", "docker", "${NOMAD_ALLOC_NAME}"]
        	port = "mongo"

			check {
          		name     = "${NOMAD_JOB_NAME}:${NOMAD_ALLOC_INDEX}-tcp"
        		type     = "tcp"
        		interval = "10s"
        		timeout  = "2s"
        	}
		}
    }

    task "elasticsearch" {
        driver = "docker"

        config {
            image = "docker.elastic.co/elasticsearch/elasticsearch-oss:6.8.2"
        }

        env {
            http.host = "0.0.0.0"
            ES_JAVA_OPTS = "-Xms512m -Xmx512m"
        }

        resources {
			cpu = 100
			memory = 768
        	network {
        		mbits = 1
				port "elasticsearch" {
                    static = 9200
                }
        	}
        }

        service {
			name = "${JOB}-${NOMAD_ALLOC_ID}"
        	tags = ["elasticsearch", "docker", "${NOMAD_ALLOC_NAME}"]
        	port = "elasticsearch"

			check {
          		name     = "${NOMAD_JOB_NAME}:${NOMAD_ALLOC_INDEX}-tcp"
        		type     = "tcp"
        		interval = "10s"
        		timeout  = "2s"
        	}
		}
    }

    task "graylog" {
        driver = "docker"

        config {
            image = "graylog/graylog:3.1"

            //network_mode = "host"
        }

        env {
            GRAYLOG_HTTP_EXTERNAL_URI = "http://127.0.0.1:9000/"
            GRAYLOG_MONGODB_URI = "mongodb://${NOMAD_ADDR_mongo}"
            //GRAYLOG_ELASTICSEARCH_HOSTS = "http://${NOMAD_ADDR_elasticsearch}"
        }

        resources {
			cpu = 100
			memory = 768
        	network {
        		mbits = 1
				port "http" {
                    static = 9000
                }

                port "gelf" {
                    static = 12201
                }
        	}
        }

        service {
			name = "${JOB}-${NOMAD_ALLOC_ID}"
        	tags = ["graylog", "docker", "${NOMAD_ALLOC_NAME}"]
        	port = "http"

			check {
          		name     = "${NOMAD_JOB_NAME}:${NOMAD_ALLOC_INDEX}-tcp"
        		type     = "tcp"
        		interval = "10s"
        		timeout  = "2s"
        	}
		}
    }
}