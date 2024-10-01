data "wxone_project" "default" {}

resource "wxone_key" "example" {
  name         = "example"
  public_key   = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCuDXAz3J2RzkQ3/DApu5HWjbuX3tmCpl90h9rCTvFhyotoPR29ggRtsXEmQXFbE8S2Yw05JdSu8sI8Y61bOyKZG/mgY6YV7+SXvMfk0+3ooLAH0dwRUcH8AkQwwQJfFADBC54z1C+F8r7e+gAjuoIJgfJQ/JtAjvVx8IhfSLgjtBEXDKqDGhwu3ukz4Gm+hr90yZN7RV+JGlFCKxXQrN/X2ldBAGGypwe0oQeENh6gGB/49U4RYVrSE6OWanywEWWMWkLtempF3D2gqKqbB/U3pzJYI9gHuNwtKxfR4d2yfSuBqB0SuwDgncHc6YBQYn48K/yKCFjdXpwTzUbVjiPBlTUQcjwT2VQIAxIxzovOrh/8RyVdVqKRWg+JLCJJtjrltIKnF6CE0rzs9D31PB9V1pitZaJFjGIokQKJ6tVqH/PmARmRlwyoRhim1hRSTDhF9YxaPmrGUHuwi0agpjTvvJi1HIJiJV68kMB1o+l120qAZ7p56cEqQT/OGiIUATU= ubuntu@cw-development"
  project_wide = true
  project_id   = data.wxone_project.default.id
}