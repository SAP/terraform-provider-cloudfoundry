# Migration 

 Although the definitions look similar the [cloudfoundry-community](https://github.com/cloudfoundry-community/terraform-provider-cloudfoundry) and [SAP cloudfoundry](https://github.com/SAP/terraform-provider-cloudfoundry) providers are not the same. There would be differences in the definition as well as the state structure among the existing cloudfoundry providers. 
 
Therefore, providers cannot be simply interchanged and one will need to write new configuration to use the resource and datasource definitions available with this provider. Additionally one will need to import the existing resources into the state for the new provider.

 **Keep a backup of the existing configuration and state. Also, before one starts the migration process, please understand the differences in the structure of the resources and datasources as we support the latest V3 API's from CloudFoundry API. Refer our documentation for more details.**

The newer v3 APIs have brought in changes to certain parameters of the existing resources and datasources to enhance the functionality and optimise the way the APIs behave. There could be chances that some of the parameters might have been removed and newer parameters would have been added. Please refer to the [V3 API Specifications](https://v3-apidocs.cloudfoundry.org/version/3.166.0/index.html).

This migration guide helps one in transitioning from the existing [community CF provider](https://github.com/cloudfoundry-community/terraform-provider-cloudfoundry) to the [current CF provider](https://github.com/SAP/terraform-provider-cloudfoundry) and states the difference between the two.

## Maintaining both providers

A one shot rewrite of the entire configuration may not be practical especially when one has a large configuration and some of the resources/datasources are not available in the newer provider. In this case it makes sense to keep both the providers and step by step rewrite the configuration moving from the old to the new provider.

Steps involved here would be:

- Obtain details (required fields) of resource/datasource from state
- Remove the reource/datasource from the state
- Update  resource/datasource in the configuration to the new providers structure.
- Import the resource into the state file

### Example

Consider the following terraform configuration:

```
terraform {
  required_providers {
    cloudfoundry = {
     source = "cloudfoundry-community/cloudfoundry"
      version = "0.53.1"
    }
  }
}

provider "cloudfoundry" { 
  api_url = "https://api.example.com"   
  username = "your-username"            
  password = "your-password"                                 
}

resource "cloudfoundry_org" "org" {
  name = "tf-test"
}

resource "cloudfoundry_space" "space" {
  name = "my-space"                     
  org  =  cloudfoundry_org.org.id                 
}

resource "cloudfoundry_app" "my-app" {
  name  = "my-app"                      
  space = cloudfoundry_space.my_space.id
  path = "/path/to/your/app"            
}
```

Let's assume that the resource `cloudfoundry_app` - `my-app` is being migrated to the new provider.

Steps:

1) Get details of `my-app` from the state.

```

➜ tf state show cloudfoundry_app.my-app | grep -m 1 id
    id                              = "d0348ed0-6e89-4836-80db-16479526a748"
```

2) Remove `my-app` from the state.

```
➜ terraform state rm cloudfoundry_app.my-app
```

3) Replace the configuration of the `cloudfoundry_app` resource with the new provider and initialise if required with ` terraform init -upgrade`.

```
terraform {
  required_providers {
    cloudfoundry = {
     source = "cloudfoundry-community/cloudfoundry"
      version = "0.53.1"
    }
    
    #Add new provider
    cloudfoundry-v3-new = {
      source = "SAP/cloudfoundry"
      version = "0.2.0-beta"
    }
  }
}

#Old provider remains to handle non migrated resources
provider "cloudfoundry" { 
  api_url = "https://api.example.com"   
  username = "your-username"            
  password = "your-password"            
}

#Configure the new SAP Terraform Cloud foundry provider with alias
provider "cloudfoundry-v3-new" {  
  api_url = "https://api.example.com"   
  username = "your-username"            
  password = "your-password"      
}

resource "cloudfoundry_org" "org" {
  name = "tf-test"
}

resource "cloudfoundry_space" "space" {
  name = "my-space"                     
  org  =  cloudfoundry_org.org.id                 
}

#  --- Use provider meta-argument to differentiate the providers ---
 
resource "cloudfoundry_app" "my-app" {
  provider = cloudfoundry-v3-new
  name = "my-app"                     
  org_name  = cloudfoundry_org.org.name
  space_name = cloudfoundry_space.my_space.name               
}
```

4) Import the state using the updated resource definition:

```
terraform import cloudfoundry_app.my-app d0348ed0-6e89-4836-80db-16479526a748
```

5) The state should be imported. After successfull import run `terraform plan` to verify. It might prompt to reapply the configuration to populate some attribute values after which it is ready for use.

## Resource Definition Changes

> [!NOTE]  
> For every resource and datasource, the computed attributes `created_at` and `updated_at` have been added in the new provider in line with V3.

While most resources have  maintained  the same structure, some resources needed changes to follow the V3 API structure. Following is a list of resources that have changed

- [Application](./resources/app.md)
- [Domain](./resources/domain.md)
- [Org Quota](./resources/org_quota.md)
- [Org Role](./resources/org_role.md)
- [Organisation](./resources/org.md)
- [Route](./resources/route.md)
- [Security Group](./resources/security_group.md)
- [Service credential Binding](./resources/service_credential_binding.md)
- [Service Instance](./resources/service_instance.md)
- [Space Quota](./resources/space_quota.md)
- [Space Role](./resources/space_role.md)
- [Space](./resources/space.md)
- [User](./resources/user.md)


## DataSource Definition Changes

While most dataSources have  maintaied the same structure, some dataSources needed changes to follow the V3 API structure. Following is a list of datasources that have changed

- [Application](./data-sources/app.md)
- [Domain](./data-sources/domain.md)
- [Org Quota](./data-sources/org_quota.md)
- [Organisation](./data-sources/org.md)
- [Route](./data-sources/route.md)
- [Security Group](./data-sources/security_group.md)
- [Service credential Binding](./data-sources/service_credential_binding.md)
- [Service Instance](./data-sources/service_instance.md)
- [Service](./data-sources/service.md)
- [Space Quota](./data-sources/space_quota.md)
- [Space](./data-sources/space.md)
- [User](./data-sources/user.md)



