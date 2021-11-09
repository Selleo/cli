# cli

## AWS

### ECS

Deploy new docker image to ECS cluster and given service. This will create a new task revision and create a deployment for given service.

```
selleo aws ecs deploy --cluster CLUSTER_ID --service SERVICE_NAME --docker-image DOCKER_IMAGE --region AWS_REGION
```

