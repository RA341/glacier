data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/migrator",
  ]
}
env "gorm" {
  src = data.external_schema.gorm.url
  dev = "sqlite://file?mode=memory&_fk=1"
  migration {
    dir = "file://internal/database/generated/migrations?format=goose"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}