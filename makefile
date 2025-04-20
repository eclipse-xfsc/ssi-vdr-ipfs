include $(PWD)/.env
include $(PWD)/.env.credentials

docker-run:
	docker rm -f vdr-ipfs
	docker build -t vdr-ipfs-image --build-arg="GITLAB_TOKEN=$(GITLAB_TOKEN)" -f deployment/docker/Dockerfile .
	docker run --name vdr-ipfs --env-file=".env" --env-file=".env.credentials" -d vdr-ipfs-image