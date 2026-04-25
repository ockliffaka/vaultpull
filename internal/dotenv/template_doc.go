// Package dotenv provides utilities for managing .env files, including
// template rendering.
//
// # Template Rendering
//
// The template sub-feature allows .env.tmpl files to reference Vault secrets
// using ${KEY} or $KEY syntax. This is useful for composing values from
// multiple secrets without duplicating them.
//
// Basic usage:
//
//	result, err := dotenv.Render("host=${DB_HOST}", secrets, dotenv.DefaultTemplateOptions())
//
// File-based rendering:
//
//	err := dotenv.RenderFile("app.env.tmpl", secrets, opts, "app.env")
//
// Use Strict mode to fail fast when a placeholder cannot be resolved:
//
//	opts := dotenv.DefaultTemplateOptions()
//	opts.Strict = true
package dotenv
