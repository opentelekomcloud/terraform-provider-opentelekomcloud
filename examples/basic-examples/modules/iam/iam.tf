
# create random number
resource "random_id" "project" {
 byte_length = 4
}

resource "opentelekomcloud_identity_project_v3" "project_1" {
  provider = "opentelekomcloud.iam"
  name = "${var.project_name}_${random_id.project.id}" 
}

data "opentelekomcloud_identity_project_v3" "project_3" {
 provider = "opentelekomcloud.iam"
  name = "eu-de"
}

resource "opentelekomcloud_identity_project_v3" "project_2" {
 provider = "opentelekomcloud.iam"
  name = "${var.project_name}_${random_id.project.id}_2" 
  description = "${var.project_desc}"
  region = "${var.region}"
  domain_id = "${var.domain_id}"
  parent_id = "${var.parent_id}"
}

resource "opentelekomcloud_identity_user_v3" "user_1" {
 provider = "opentelekomcloud.iam"
  name = "${var.user_name}_1"
}

resource "opentelekomcloud_identity_user_v3" "user_2" {
 provider = "opentelekomcloud.iam"
  name = "${var.user_name}_2"
  #description  = "${var.user_desc}"
  password = "${var.user_passd}"
  default_project_id  = "${opentelekomcloud_identity_project_v3.project_1.id}"
  domain_id = "${var.domain_id}"
  region = "${var.region}"
  enabled = "${var.user_status}"
}
data "opentelekomcloud_identity_user_v3" "user_3" {
 provider = "opentelekomcloud.iam"
  name = "h00454348"
}

resource "opentelekomcloud_identity_group_v3" "group_1" {
 provider = "opentelekomcloud.iam"
  name = "group_1"
}

resource "opentelekomcloud_identity_group_v3" "group_2" {
 provider = "opentelekomcloud.iam"
  name = "${var.group_name}"
  description = "${var.group_desc}"
  domain_id = "${var.domain_id}"
  region = "${var.region}"
}
data "opentelekomcloud_identity_group_v3" "admins" {
 provider = "opentelekomcloud.iam"
  name = "admin"
}

resource "opentelekomcloud_identity_group_membership_v3" "membership_1" {
    provider = "opentelekomcloud.iam"
        group = "${opentelekomcloud_identity_group_v3.group_1.id}"
        users = ["${opentelekomcloud_identity_user_v3.user_1.id}"   ,
                "${opentelekomcloud_identity_user_v3.user_2.id}"
                ]
}

data "opentelekomcloud_identity_role_v3" "role_1" {
 provider = "opentelekomcloud.iam"
 name = "${var.role_name}" #security admin
}

resource "opentelekomcloud_identity_role_assignment_v3" "role_assignment_1" {
 provider = "opentelekomcloud.iam"
  group_id = "${opentelekomcloud_identity_group_v3.group_1.id}"
  domain_id = "${var.domain_id}"
  role_id = "${data.opentelekomcloud_identity_role_v3.role_1.id}"
} 