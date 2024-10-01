data "wxone_project" "default" {}

data "wxone_image" "ubuntu" {
  project_id = data.wxone_project.default.id
  name       = "ubuntu"
}
