__CSSMINIFIER_APEX__ {
  proxy / localhost:__CSSMINIFIER_PORT__ {
    transparent
  }
  tls chilts@appsattic.com
  log stdout
  errors stderr
}

www.__CSSMINIFIER_APEX__ {
  redir http://__CSSMINIFIER_APEX__{uri} 302
}
