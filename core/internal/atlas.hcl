data "external_schema" "sqlite" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./migrator",
    "sqlite"
  ]
}

env "sqlite" {
  src = data.external_schema.sqlite.url
  dev = "sqlite://file?mode=memory&_fk=1"
  migration {
    dir = "file://database/generated/migrations/sqlite?format=goose"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

data "external_schema" "postgres" {
  program = [
    "go", "run", "-mod=mod", "./migrator", "postgres",
  ]
}

env "postgres" {
  src = data.external_schema.postgres.url
  dev = "docker://postgres/16/dev"
  migration {
    dir = "file://database/generated/migrations/postgres?format=goose"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}