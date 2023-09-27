# Setup an email alias

resource "glesys_emailalias" "alice" {
  emailalias  = "info@example.com"
  goto      = "alice@example.com,bob@example.com"
}
