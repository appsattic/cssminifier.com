[program:cssminifier]
directory = /home/chilts/src/appsattic-cssminifier.com
command = /home/chilts/src/appsattic-cssminifier.com/bin/cssminifier
user = chilts
autostart = true
autorestart = true
stdout_logfile = /var/log/chilts/cssminifier-stdout.log
stderr_logfile = /var/log/chilts/cssminifier-stderr.log
environment =
    CSSMINIFIER_PORT="__CSSMINIFIER_PORT__",
    CSSMINIFIER_APEX="__CSSMINIFIER_APEX__",
    CSSMINIFIER_BASE_URL="__CSSMINIFIER_BASE_URL__",
    CSSMINIFIER_DIR="__CSSMINIFIER_DIR__",
    CSSMINIFIER_GOOGLE_ANALYTICS="__CSSMINIFIER_GOOGLE_ANALYTICS__"
