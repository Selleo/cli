# cli

Install via homebrew:
```
brew tap Selleo/cli
brew install selleo

# upgrade
brew upgrade selleo
```

or download binary for your system.

## AWS

### Dev

You can fetch secrets from AWS store parameters and run the command:

```
selleo aws dev --region eu-central-1 --path /office/dev/api npm run start
```

### Generators

Generators are used to pre-generate templates that you can furhter adjust.

#### Terraform

Generate app environment:
```
selleo gen terraform app \
  --tf-cloud-org        selleo             \
  --tf-cloud-workspace  ict-til-production \
  --region              eu-west-1          \
  --stage               production         \
  --namespace           til                \
  --name                api                \
  --domain              til.selleo.com     \
  --subdomain           api
```

#### GitHub workflows

Generate staging frontend workflow:
```
selleo gen github frontend      \
  --workdir .                   \
  --domain  staging.example.com \
  --region  eu-west-1           \
  --app_id  office              \
  --stage   staging             
```

Generate production frontend workflow from subfolder and trigger build only on git tag push:
```
selleo gen github frontend     \
  --workdir packages/office    \
  --domain  example.com        \
  --region  eu-west-1          \
  --app_id  office             \
  --stage   production         \
  --tag-release
```

Generate staging backend workflow with no extra task running at the end:
```
selleo gen github backend           \
  --workdir     .                   \
  --domain      beta.example.com    \
  --subdomain   api                 \
  --region      eu-west-1           \
  --ecs-cluster rails-1234          \
  --ecs-service api                 \
  --stage       staging             
```

Generate production backend workflow with extra task run at the end:
```
selleo gen github backend           \
  --workdir     .                   \
  --domain      example.com         \
  --subdomain   api                 \
  --region      eu-west-1           \
  --ecs-cluster rails-1234          \
  --ecs-service api                 \
  --stage       production          \
  --one-off     migrate             \
  --tag-release
```

## About Selleo

![selleo](https://raw.githubusercontent.com/Selleo/selleo-resources/master/public/github_footer.png)

Software development teams with an entrepreneurial sense of ownership at their core delivering great digital products and building culture people want to belong to. We are a community of engaged co-workers passionate about crafting impactful web solutions which transform the way our clients do business.

All names and logos for [Selleo](https://selleo.com/about) are trademark of Selleo Labs Sp. z o.o. (formerly Selleo Sp. z o.o. Sp.k.)

