package internal

import "embed"

//go:embed database/migrations/*
var MigrationsFS embed.FS

//go:embed templates/*
var TemplateFS embed.FS

//go:embed static/*
var StaticFS embed.FS
