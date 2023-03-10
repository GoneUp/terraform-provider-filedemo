terraform {
  required_providers {
    filedemo = {
      source = "github.com/goneup/terraform-provider-filedemo"
    }
  }
}

data "filedemo_data_file" "test_file" {
  filename = "./file.txt"
}

output "info" {
  value = data.filedemo_data_file.test_file
}
