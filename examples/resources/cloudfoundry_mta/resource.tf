resource "cloudfoundry_mta" "mtar" {
  space     = "02c0cc92-6ecc-44b1-b7b2-096ca19ee143"
  mtar_url  = "https://github.com/Dray56/mtar-archive/releases/download/v1.0.0/a.cf.app.mtar"
  namespace = "test"
}
