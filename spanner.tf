resource "google_spanner_instance" "global-game-spanner" {
  config       = var.spanner_config.location
  display_name = var.spanner_config.instance_name
  num_nodes    = var.spanner_config.num_nodes
}

resource "google_spanner_database" "spanner-database" {
  instance                 = google_spanner_instance.global-game-spanner.name
  name                     = var.spanner_config.db_name
  version_retention_period = "3d"
  deletion_protection = false
}
