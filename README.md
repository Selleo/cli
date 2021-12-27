# cli

## AWS

### ECS

Deploy new docker image to ECS cluster and service. This will create a new task revision and create a deployment.

```
selleo aws ecs deploy --cluster CLUSTER_ID --service SERVICE_NAME --docker-image DOCKER_IMAGE --region AWS_REGION
```

### Secrets

Get all the secrets and export them in shell:

```
$(selleo aws secrets export --region REGION --id SECRET_ID)
```

And new secret (KEY/VALUE are positional arguments):

```
aws secrets set --region REGION --id SECRET_ID  KEY VALUE
```


## About Selleo

![selleo](https://raw.githubusercontent.com/Selleo/selleo-resources/master/public/github_footer.png)

Software development teams with an entrepreneurial sense of ownership at their core delivering great digital products and building culture people want to belong to. We are a community of engaged co-workers passionate about crafting impactful web solutions which transform the way our clients do business.

All names and logos for [Selleo](https://selleo.com/about) are trademark of Selleo Labs Sp. z o.o. (formerly Selleo Sp. z o.o. Sp.k.)

