# Automation Frameworks

`vi /lib/systemd/system/nvued.service`

```
[Service]
LimitNOFILE=1024
```

`systemctl daemon-reload`