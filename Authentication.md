# Authentication Mechanisms for Cloudfoundry Terraform Provider

The cloudfoundry terraform provider supports any of the following authentication mechanism currently:

## USERNAME-PASSWORD

Use the env variables `CF_API_URL`, `CF_USER` and `CF_PASSWORD`.

Alternatively,

```hcl
provider "cloudfoundry" {
    api_url = "<CF-API-URL>"
    user = "<USER-ID>"
    password = "<PASSWORD>"
}
```

## CLIENT ID-CLIENT SECRET

Use the env variables `CF_API_URL`, `CF_CF_CLIENT_ID` and `CF_CF_CLIENT_SECRET`.

Alternatively, 

```hcl
provider "cloudfoundry" {
    api_url = "<CF-API-URL>"
    cf_client_id = "<CF-CLIENT-ID>"
    cf_client_secret = "<CF-CLIENT-SECRET>"
}
```

## Using cf-cli configuration.

If you have installed the [cf-cli](https://docs.cloudfoundry.org/cf-cli/) and have [logged into the environment](https://docs.cloudfoundry.org/cf-cli/getting-started.html#login), the the cloudfoundry terraform provider can use the default configuration of the cf-cli (present in `~/.cf` folder) to connect to the environment.

```hcl
provider cloudfoundry {}
```

If the provider is initialized without any parameters and no environment variables are set, then the provider will try to connect this way.


