terraform {
  required_providers {
    filedemo = {
      source = "github.com/goneup/terraform-provider-filedemo"
    }
  }
}

resource "filedemo_file" "example" {
  content  = "Some sdfsdfsd contentsdd222"
  filename = "./test2.txt"
}
