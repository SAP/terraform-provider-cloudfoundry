# terraform import cloudfoundry_mtar.<resource_name> <space_guid/mta_id> OR 
# terraform import cloudfoundry_mtar.<resource_name> <space_guid/mta_id/namespace> if MTA in custom namespace

terraform import cloudfoundry_mtar.my_mtar 02c0cc92-6ecc-44b1-b7b2-096ca19ee143/a.cf.app/hello

