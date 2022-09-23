# Setup an email account

resource "glesys_emailaccount" "bob" {
  emailaccount       = "bob@example.com"
  password           = "SecretPassword123"
  autorespond        = "yes"
  autorespondmessage = "I'm away."
  quotaingib         = 2
}
resource "glesys_emailaccount" "alice" {
  emailaccount  = "alice@example.com"
  password      = "PasswordSecret321"
  antispamlevel = 5
}
