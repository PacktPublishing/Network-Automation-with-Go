

## Templates

- [IOS XE](https://github.com/networktocode/ntc-templates/blob/master/ntc_templates/templates/cisco_ios_show_version.textfsm)
- [NX-OS](https://github.com/networktocode/ntc-templates/blob/master/ntc_templates/templates/cisco_nxos_show_version.textfsm)

## Outputs

### Devices

```bash
# show version | i time
csr1000v-1 uptime is 1 hour, 25 minutes
Uptime for this control processor is 1 hour, 27 minutes
```

```bash
# show version | i time
  BIOS compile time:  
  NXOS compile time:  12/22/2019 2:00:00 [12/22/2019 14:00:37]
Kernel uptime is 6 day(s), 6 hour(s), 0 minute(s), 38 second(s)
```

### Program


```bash
Hostname: sandbox-iosxe-latest-1.cisco.com
Hardware: [CSR1000V]
SW Version: 17.3.1a
Uptime: 1 day, 54 minutes

Hostname: sandbox-nxos-1.cisco.com
Hardware: C9300v
SW Version: 9.3(3)
Uptime: 4 day(s), 23 hour(s), 17 minute(s), 28 second(s)
```

or

```bash
Hostname: sandbox-iosxe-latest-1.cisco.com
Hardware: [CSR1000V]
SW Version: 17.3.1a
Uptime: 1 day, 1 hour, 1 minute

Hostname: sandbox-nxos-1.cisco.com
Hardware: C9300v
SW Version: 9.3(3)
Uptime: 4 day(s), 23 hour(s), 24 minute(s), 7 second(s)
```