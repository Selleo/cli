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

#### GitHub workflows

Generate staging and production workflows:
```
selleo gen github frontend --workdir packages/client --domain selleo.com --region eu-central-1 --app_id website
```

## About Selleo

![selleo](https://raw.githubusercontent.com/Selleo/selleo-resources/master/public/github_footer.png)

Software development teams with an entrepreneurial sense of ownership at their core delivering great digital products and building culture people want to belong to. We are a community of engaged co-workers passionate about crafting impactful web solutions which transform the way our clients do business.

All names and logos for [Selleo](https://selleo.com/about) are trademark of Selleo Labs Sp. z o.o. (formerly Selleo Sp. z o.o. Sp.k.)

