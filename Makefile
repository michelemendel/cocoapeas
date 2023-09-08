update_packages:
	@go get -u ./...

clean:
	@rm -rf bin


# -------------------------------------------------------------------------------
# Golang's Kafka producer and consumer

build_producer:
	@go build -o bin/producer ./producer/...
	
producer: build_producer
	@./bin/producer 2>&1

build_consumer:
	@go build -o bin/consumer ./consumer/...
	
consumer: build_consumer
	@./bin/consumer 2>&1


# -------------------------------------------------------------------------------
# Kafka

kafka_up:
	@docker-compose -f docker-compose-kafka.yml up -d

kafka_down:
	@docker-compose -f docker-compose-kafka.yml down