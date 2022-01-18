## Services

### SSH

```bash
# sh run ssh    
ssh server v2
ssh server vrf default
ssh server netconf vrf default
```

### NETCONF

```bash
# sh run ssh
ssh server v2
ssh server vrf default
ssh server netconf vrf default
!
# sh run netconf
netconf agent tty
! 
# show run netconf-yang agent 
netconf-yang agent
 ssh
!
```

### gRPC

```bash
# sh run grpc
no grpc
grpc
 port 57777
 address-family ipv4
!
```

```bash
# show grpc status 
*************************show gRPC status**********************
---------------------------------------------------------------
transport                       :     grpc
access-family                   :     tcp4
TLS                             :     enabled
trustpoint                      :     
listening-port                  :     57777
max-request-per-user            :     10
max-request-total               :     128
max-streams                     :     32
max-streams-per-user            :     32
vrf-socket-ns-path              :     global-vrf
_______________________________________________________________
*************************End of showing status*****************
```