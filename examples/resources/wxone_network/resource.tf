data "wxone_project" "default" {}

resource "wxone_network" "example" {
  name              = "example"
  availability_zone = "wx_dus_1"
  project_id        = data.wxone_project.default.id
  subnets = [
    {
      name       = "example"
      ip_version = "IPv4"
      cidr       = "10.0.0.0/24"
    }
  ]
}
