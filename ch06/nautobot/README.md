## Nautobot

### Dependencies

```bash
cd ch06/nautobot/client
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
wget https://demo.nautobot.com/api/swagger.yaml\?api_version\=1.3 -O swagger.yaml
oapi-codegen --config config.yaml swagger.yaml
go mod init github.com/nautobot/go-nautobot
```

```bash
$ go mod edit -replace github.com/nautobot/go-nautobot=./client
```

### Running the program

```bash
$ go run nautobot
Manufacturer ID: 1df6255f-96d0-4264-bbf7-557cb04adc3c
Manufacturer ID: 14a630d1-4602-4bb9-b88b-813cf83ba950
Manufacturer ID: 0266f513-dae8-4183-9d45-b805b2e7df1c
Manufacturer ID: 8e570f17-fbfb-4d46-89fb-c9253f206784
Device-Type ID: ac0763ca-64e4-466b-87c8-6e9233df70ed
Device-Type ID: 06f51214-caaa-4276-80c7-6c75b46c5919
Device-Type ID: 1d6cfd62-c09f-402a-a732-865d043d2498
Device-Role ID: 8fd80869-d2b0-4cd8-ab86-e3bf13410fd7
ID: 06fee8e5-2413-4ba3-a830-79b0c471f520
DeviceRole: 8fd80869-d2b0-4cd8-ab86-e3bf13410fd7
DeviceType: 06f51214-caaa-4276-80c7-6c75b46c5919
Site: cc966ae0-c17e-4901-953f-51b85c283948
Device already present: {"count":1,"next":null,"previous":null,"results":[{"id":"06fee8e5-2413-4ba3-a830-79b0c471f520","url":"https://demo.nautobot.com/api/dcim/devices/06fee8e5-2413-4ba3-a830-79b0c471f520/","name":"ams01-ceos-02","device_type":{"id":"06f51214-caaa-4276-80c7-6c75b46c5919","url":"https://demo.nautobot.com/api/dcim/device-types/06f51214-caaa-4276-80c7-6c75b46c5919/","manufacturer":{"id":"1df6255f-96d0-4264-bbf7-557cb04adc3c","url":"https://demo.nautobot.com/api/dcim/manufacturers/1df6255f-96d0-4264-bbf7-557cb04adc3c/","name":"Arista","slug":"arista","display":"Arista"},"model":"cEOS","slug":"ceos","display":"Arista cEOS"},"device_role":{"id":"8fd80869-d2b0-4cd8-ab86-e3bf13410fd7","url":"https://demo.nautobot.com/api/dcim/device-roles/8fd80869-d2b0-4cd8-ab86-e3bf13410fd7/","name":"Router","slug":"router","display":"Router"},"tenant":null,"platform":null,"serial":"","asset_tag":null,"site":{"id":"cc966ae0-c17e-4901-953f-51b85c283948","url":"https://demo.nautobot.com/api/dcim/sites/cc966ae0-c17e-4901-953f-51b85c283948/","name":"AMS01","slug":"ams01","display":"AMS01"},"rack":null,"position":null,"face":null,"parent_device":null,"status":{"value":"active","label":"Active"},"primary_ip":null,"primary_ip4":null,"primary_ip6":null,"secrets_group":null,"cluster":null,"virtual_chassis":null,"vc_position":null,"vc_priority":null,"comments":"","local_context_schema":null,"local_context_data":null,"tags":[],"custom_fields":{"test_bug_filter_default_value":true},"config_context":{"cdp":true,"ntp":[{"ip":"10.1.1.1","prefer":false},{"ip":"10.2.2.2","prefer":true}],"lldp":true,"snmp":{"host":[{"ip":"10.1.1.1","version":"2c","community":"networktocode"}],"contact":"John Smith","location":"Network to Code - NYC | NY","community":[{"name":"ntc-public","role":"RO"},{"name":"ntc-private","role":"RW"},{"name":"networktocode","role":"RO"},{"name":"secure","role":"RW"}]},"aaa-new-model":false,"acl":{"definitions":{"named":{"PERMIT_ROUTES":["10 permit ip any any"]}}},"route-maps":{"PERMIT_CONN_ROUTES":{"seq":10,"type":"permit","statements":["match ip address PERMIT_ROUTES"]}}},"created":"2022-05-25","last_updated":"2022-05-25T20:52:14.105314Z","display":"ams01-ceos-02"}]}
```

Creates [device.json](device.json):

```json
{
    "name": "ams01-ceos-02",
    "device_type": {
        "slug": "ceos"
    },
    "device_role": {
        "slug": "router"
    },
    "site": {
        "slug": "ams01"
    }
}
```

### Creating codegen config

```yaml
$ oapi-codegen --output-config --old-config-style -generate "client,types" -o nautobot.go -package nautobot swagger.yaml > config.yaml
package: nautobot
generate:
  client: true
  models: true
output: nautobot.go
```