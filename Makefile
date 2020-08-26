service_name := nalupi
version := 0.0.1
project_id := jjkoh95

go-run:
	go run cmd/main.go
go-build:
	CGO_ENABLED=0 GOOS=linux go build -v -o server cmd/server/main.go
docker-build-tag:
	docker build -t asia.gcr.io/${project_id}/${service_name}:${version} .
docker-push-google-registry:
	docker push asia.gcr.io/${project_id}/${service_name}:${version}
gcp-deploy-cloud-run:
	gcloud run deploy ${service_name} \
	--image asia.gcr.io/${project_id}/${service_name}:${version} \
	--region=asia-east1 --platform=managed \
	--concurrency=40 --memory 256Mi --timeout=600s --max-instances 1 --cpu 1 \
	--allow-unauthenticated
build-deploy:
	make go-build
	make docker-build-tag
	make docker-push-google-registry
	make gcp-deploy-cloud-run