data "cloudfoundry_remote_mtar_hash" "hash" {
  url = "https://github.com/Dray56/mtar-archive/releases/download/v1.0.0/a.cf.app.mtar"
}

output "sha_sum" {
  value = data.cloudfoundry_remote_mtar_hash.hash.id
}