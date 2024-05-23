# Migration 

 Although the definitions look similar the providers are not the same. There would be differences in the definition as well as the state structure among different providers. 
 
 Therefore, provders cannot be simply interchanged and you will need to write new configuration to use the resource and datasource definitions available with this provider. Additionally you will need to import the existing resources into the state for the new provider.

 **Keep a backup of your existing configuration and state. Also, before you start the migration process, please understand the differences in the resources and datasources of the the new provider from the documentation.**

 The newer v3 APIs have brought in changes to certain parameters to enhance funcationality and optimise the way the APIs behave. As a result some of the parameters may have been removed and newer parameters may have been added. Please refer to the (v3 api secifications)[https://v3-apidocs.cloudfoundry.org/version/3.166.0/index.html]

## Maintaining both providers

A one shot rewrite of the entire configuration may not be practical especially when you have a large configuration and some of the resources/datasources are not available in the newer provider. In this case it makes sense to keep both the providers and step by step rewirte the configuration moving from the old to the new provider. Introduce the new provider into the configuration with `alias`.

Steps involved here would be:

- Obtain details (required fields) of resource/datasource from state
- Remove the reource/datasource from the state
- Update  resource/datasource in the configuration to the new providers structure.
- Import the resource into the state file

### Example

Consider the following terraform configuration:

```
provider "cloudfoundry" { #community terraform provider
  api_url = "https://api.example.com"   
  username = "your-username"            
  password = "your-password"            
  org = "your-org"                             
}

resource "cloudfoundry_space" "my_space" {
  name = "my-space"                     
  org  = "your-org"                     
}

resource "cloudfoundry_app" "my_app" {
  name  = "my-app"                      
  space = cloudfoundry_space.my_space.id
  memory = 1G                         
  # Path to your application directory
  path = "/path/to/your/app"            
}
```
Let's assume that the resource `cloudfoundry_app` - `my_app` is being migrated to the new provider.

Steps:

1) Get details of `my_app` from the state

```
echo cloudfoundry_app.my_app.id | terraform console

"12345b61-9551-4d43-bcd3-3f5c8fb775fd"

```
Obtain the required fields (for the new provider) : name, org_name, space_name

2) Remove the `cloudfoundry_space` from the state:

```
terraform state rm cloudfoundry_space.my_space
```

3) Replace the configuration of the `cloudfoundry_space` resource from old to the new provider:

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
provider "cloudfoundry" { #community terraform provider
  api_url = "https://api.example.com"   
  username = "your-username"            
  password = "your-password"            
  org = "your-org"                             
}

#Configure the new SAP Terraform Cloud foundry provider with alias

provider "cloudfoundry-v3-new" { 
   alias  = "new-provider"    # Add alias                         
}


#Old provider remains to handle non migrated resources
provider "cloudfoundry" { #community terraform provider
  api_url = "https://api.example.com"   
  username = "your-username"            
  password = "your-password"            
  org = "your-org"                             
}


#  ---- Use Alias to diffrentiate the providers ---
 
resource "cloudfoundry_space" "my_space" {
  provider = cloudfoundry-v3-new.new-provider 
  name = "my-space"                     
  org  = "your-org"                     
}
```

4) Import the state using the updated resource definition:

```
terraform import cloudfoundry_space.my_space 12345b61-9551-4d43-bcd3-3f5c8fb775fd

terraform import cloudfoundry_space.my_space <SPACE-ID>
```

5) The state should be imported. After successfull import run `terraform plan`. `No changes` should be present

## Resource Definition Changes

While most resources have be maintaied with the same structure, some resources needed changes to follow the V3 API structure. Following is a list of resources/datasources that have changed

- [Domain](domain.md)
- [Org Quota](org_quota.md)
- [Org Role](org_role.md)
- [Organisation](organisation.md)





