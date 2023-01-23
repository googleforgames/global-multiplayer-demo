resource "google_spanner_instance" "main" {
  config       = var.spanner_config.location
  display_name = var.spanner_config.instacne_name
  num_nodes    = var.spanner_config.num_nodes
}

resource "google_spanner_database" "spanner-database" {
  instance                 = google_spanner_instance.main.name
  name                     = var.spanner_config.db_name
  version_retention_period = "3d"
  ddl = [
    "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
    "CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
  ]
  deletion_protection = false
}
