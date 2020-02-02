job "home-updater" {
  datacenters = ["pi"]
  type        = "service"

  task "home-updater" {
    driver = "exec"

    config {
      command = "/usr/local/bin/home-updater"
    }

    template {
      source      = "/usr/local/bin/home-updater.env"
      destination = "secrets/file.env"
      env         = true
    }
  }
}
