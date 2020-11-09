openrc:
	docker-compose -f deployments/svc/docker-compose.yaml --project-directory . up
systemd:
	docker-compose -f deployments/svc-systemd/docker-compose.yaml --project-directory . up
