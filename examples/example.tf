terraform {
  required_version = ">= 1.0"
  required_providers {
    solus = {
      version = "0.1.0"
      source  = "solusio/solus"
    }
  }
}

data "solus_location" "default" {
  id = 1
}

resource "solus_project" "new" {
  name = "New Project"
}

resource "solus_os_image" "alpine" {
  name = "Alpine"
}

resource "solus_os_image_version" "alpine_3_15" {
  os_image_id         = solus_os_image.alpine.id
  version             = "3.15"
  url                 = "https://images.prod.solus.io/solus-alpine-3.15.qcow2"
  cloud_init_version  = "v2"
  virtualization_type = "kvm"
}

resource "solus_plan" "kvm_fb_qcow2_10_1024_1" {
  name                = "KVM FB QCOW2 10-1024-1"
  virtualization_type = "kvm"
  storage_type        = "fb"
  image_format        = "qcow2"
  params {
    disk   = 10
    ram_mb = 1024
    vcpu   = 1
  }
  available_locations = [
    data.solus_location.default.id
  ]
  available_os_image_versions = [
    solus_os_image_version.alpine_3_15.id
  ]
}

resource "solus_ssh_key" "ssh_key" {
  name = "My SSH Key"
  body = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFbNUAbScm4GGCTjKwgC4T/zitU9kdHFKvOp3U//bVFQ"
}

resource "solus_virtual_server" "vs_1" {
  hostname    = "vs.example.com"
  description = "VS created through terraform"

  location_id         = data.solus_location.default.id
  os_image_version_id = solus_os_image_version.alpine_3_15.id
  plan_id             = solus_plan.kvm_fb_qcow2_10_1024_1.id
  project_id          = solus_project.new.id
  ssh_keys = [
    solus_ssh_key.ssh_key.id,
  ]

  user_data = <<EOT
#cloud-config
runcmd:
  - echo "echo "Hello from User Data"" > /root/hello-world-user-data.sh
  - [ chmod, +x, "/root/hello-world-user-data.sh" ]
EOT
}
