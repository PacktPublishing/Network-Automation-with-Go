# Random notes

## Generating Go binding from protobuf files

The Go generated code in `ems_grpc.pb.go` is the result of the following:

- `proto/ems`

```bash
$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    --go_opt=Mproto/ems/ems_grpc.proto=proto/ems \
    --go-grpc_opt=Mproto/ems/ems_grpc.proto=proto/ems \
    proto/ems/ems_grpc.proto
```

## Req

Temporary

```
telemetry model-driven
 sensor-group BGPNeighbor
  sensor-path Cisco-IOS-XR-ipv4-bgp-oper:bgp/instances/instance/instance-active/default-vrf/afs/af/neighbor-af-table/neighbor
 !
 subscription BGP
  sensor-group-id BGPNeighbor sample-interval 5000
 !
!
```

## Config

```json
go run grpc
{
  "openconfig-network-instance:network-instances": {
    "network-instance": [
      {
        "config": {
          "name": "default"
        },
        "name": "default",
        "protocols": {
          "protocol": [
            {
              "bgp": {
                "global": {
                  "afi-safis": {
                    "afi-safi": [
                      {
                        "afi-safi-name": "openconfig-bgp-types:IPV4_UNICAST",
                        "config": {
                          "afi-safi-name": "openconfig-bgp-types:IPV4_UNICAST",
                          "enabled": true
                        }
                      }
                    ]
                  },
                  "config": {
                    "as": 64512,
                    "router-id": "198.51.100.0"
                  }
                }
              },
              "config": {
                "identifier": "openconfig-policy-types:BGP",
                "name": "default"
              },
              "identifier": "openconfig-policy-types:BGP",
              "name": "default"
            }
          ]
        }
      }
    ]
  }
}



BGP config applied on sandbox-iosxr-1.cisco.com:57777


Config from sandbox-iosxr-1.cisco.com:
{
 "openconfig-network-instance:network-instances": {
  "network-instance": [
   {
    "name": "default",
    "protocols": {
     "protocol": [
      {
       "identifier": "openconfig-policy-types:BGP",
       "name": "default",
       "bgp": {
        "global": {
         "config": {
          "as": 64512,
          "router-id": "198.51.100.0"
         },
         "afi-safis": {
          "afi-safi": [
           {
            "afi-safi-name": "openconfig-bgp-types:IPV4_UNICAST",
            "config": {
             "afi-safi-name": "openconfig-bgp-types:IPV4_UNICAST",
             "enabled": true
            }
           },
           {
            "afi-safi-name": "openconfig-bgp-types:IPV6_UNICAST",
            "config": {
             "afi-safi-name": "openconfig-bgp-types:IPV6_UNICAST",
             "enabled": true
            }
           }
          ]
         }
        },
        "neighbors": {
         "neighbor": [
          {
           "neighbor-address": "2001:db8:cafe::2",
           "config": {
            "neighbor-address": "2001:db8:cafe::2",
            "peer-as": 64512,
            "description": "iBGP session"
           },
           "afi-safis": {
            "afi-safi": [
             {
              "afi-safi-name": "openconfig-bgp-types:IPV6_UNICAST",
              "config": {
               "afi-safi-name": "openconfig-bgp-types:IPV6_UNICAST",
               "enabled": true
              }
             }
            ]
           }
          }
         ]
        }
       }
      }
     ]
    }
   }
  ]
 }
}
```

## OpenConfig models

[http://ops.openconfig.net](http://ops.openconfig.net/branches/models/telemetry-version/)

- openconfig-acl
- openconfig-aft
- openconfig-bfd-ni
- openconfig-bfd
- openconfig-bgp-rib
- openconfig-bgp
- openconfig-catalog
- openconfig-interfaces
- openconfig-isis
- openconfig-lacp
- openconfig-lldp
- openconfig-local-routing
- openconfig-mpls
- openconfig-multicast
- openconfig-network-instance-sr-rsvp-coexistence
- openconfig-network-instance-sr
- openconfig-network-instance-srte-policy
- openconfig-network-instance
- openconfig-openflow
- openconfig-optical-amplifier
- openconfig-ospf
- openconfig-platform
- openconfig-probes
- openconfig-qos
- openconfig-relay-agent
- openconfig-routing-policy
- openconfig-stp
- openconfig-system
- openconfig-telemetry
- openconfig-terminal-device
- openconfig-transport-line-protection
- openconfig-types
- openconfig-vlan
- openconfig-wavelength-router

## XR Capabilities

```bash
RP/0/RP0/CPU0:iosxr#show netconf-yang capabilities | utility egrep ".*openconfig"
Sat Jan 15 04:01:57.543 UTC
http://cisco.com/ns/yang/cisco-xr-openconfig-acl-deviations                       |2018-02-07| 
http://cisco.com/ns/yang/cisco-xr-openconfig-alarms-deviations                    |2019-11-14| 
http://cisco.com/ns/yang/cisco-xr-openconfig-bfd-deviations                       |2019-05-06| 
http://cisco.com/ns/yang/cisco-xr-openconfig-bgp-policy-deviations                |2019-08-31| 
http://cisco.com/ns/yang/cisco-xr-openconfig-if-ethernet-deviations               |2016-05-16| 
http://cisco.com/ns/yang/cisco-xr-openconfig-if-ip-deviations                     |2017-02-07| 
http://cisco.com/ns/yang/cisco-xr-openconfig-if-ip-ext-deviations                 |2019-09-13| 
http://cisco.com/ns/yang/cisco-xr-openconfig-if-tunnel-deviations                 |2020-03-20| 
http://cisco.com/ns/yang/cisco-xr-openconfig-interfaces-deviations                |2016-05-16| 
http://cisco.com/ns/yang/cisco-xr-openconfig-isis-policy-deviations               |2019-06-10| 
http://cisco.com/ns/yang/cisco-xr-openconfig-lacp-deviations                      |2019-09-10| 
http://cisco.com/ns/yang/cisco-xr-openconfig-lldp-deviations                      |2017-03-08| 
http://cisco.com/ns/yang/cisco-xr-openconfig-local-routing-deviations             |2017-06-28| 
http://cisco.com/ns/yang/cisco-xr-openconfig-network-instance-deviations          |2019-10-23| 
http://cisco.com/ns/yang/cisco-xr-openconfig-platform-cpu-deviations              |2018-12-10| 
http://cisco.com/ns/yang/cisco-xr-openconfig-platform-deviations                  |2018-12-05| 
http://cisco.com/ns/yang/cisco-xr-openconfig-platform-port-deviations             |2018-12-10| 
http://cisco.com/ns/yang/cisco-xr-openconfig-platform-psu-deviations              |2018-12-10| 
http://cisco.com/ns/yang/cisco-xr-openconfig-platform-transceiver-deviations      |2019-09-05| 
http://cisco.com/ns/yang/cisco-xr-openconfig-rib-bgp-deviations                   |2016-10-16| 
http://cisco.com/ns/yang/cisco-xr-openconfig-routing-policy-deviations            |2019-02-06| 
http://cisco.com/ns/yang/cisco-xr-openconfig-rsvp-sr-ext-deviations               |2019-05-30| 
http://cisco.com/ns/yang/cisco-xr-openconfig-system-deviations                    |2018-12-12| 
http://cisco.com/ns/yang/cisco-xr-openconfig-telemetry-deviations                 |2017-03-09| 
http://cisco.com/ns/yang/cisco-xr-openconfig-vlan-deviations                      |2019-04-18| 
http://openconfig.net/yang/aaa                                                    |2018-04-12| 
http://openconfig.net/yang/aaa/types                                              |2018-04-12| 
http://openconfig.net/yang/acl                                                    |2017-05-26|*
http://openconfig.net/yang/aft                                                    |2017-05-10| 
http://openconfig.net/yang/aft/ni                                                 |2017-01-13| 
http://openconfig.net/yang/alarms                                                 |2018-01-16| 
http://openconfig.net/yang/alarms/types                                           |2018-01-16| 
http://openconfig.net/yang/bfd                                                    |2018-11-21| 
http://openconfig.net/yang/bgp-policy                                             |2017-07-30|*
http://openconfig.net/yang/bgp-types                                              |2017-02-02| 
http://openconfig.net/yang/fib-types                                              |2017-05-10| 
http://openconfig.net/yang/header-fields                                          |2017-05-26| 
http://openconfig.net/yang/interfaces                                             |2019-11-19|*
http://openconfig.net/yang/interfaces/aggregate                                   |2016-05-26| 
http://openconfig.net/yang/interfaces/ethernet                                    |2020-05-06|*
http://openconfig.net/yang/interfaces/ethernet-ext                                |2018-11-21| 
http://openconfig.net/yang/interfaces/ip                                          |2019-01-08|*
http://openconfig.net/yang/interfaces/ip-ext                                      |2018-11-21|*
http://openconfig.net/yang/interfaces/tunnel                                      |2018-11-21| 
http://openconfig.net/yang/isis-lsdb-types                                        |2018-11-21| 
http://openconfig.net/yang/isis-types                                             |2018-11-21| 
http://openconfig.net/yang/lacp                                                   |2017-05-05|*
http://openconfig.net/yang/lldp                                                   |2016-05-16|*
http://openconfig.net/yang/lldp/types                                             |2016-05-16| 
http://openconfig.net/yang/network-instance                                       |2017-01-13|*
http://openconfig.net/yang/network-instance-l3                                    |2017-01-13| 
http://openconfig.net/yang/network-instance-types                                 |2016-12-15| 
http://openconfig.net/yang/openconfig-ext                                         |2018-10-17| 
http://openconfig.net/yang/openconfig-if-types                                    |2018-01-05| 
http://openconfig.net/yang/openconfig-isis                                        |2018-11-21| 
http://openconfig.net/yang/openconfig-isis-policy                                 |2018-11-21| 
http://openconfig.net/yang/openconfig-types                                       |2019-04-16| 
http://openconfig.net/yang/packet-match-types                                     |2017-05-26| 
http://openconfig.net/yang/platform                                               |2019-04-16|*
http://openconfig.net/yang/platform-types                                         |2019-06-03| 
http://openconfig.net/yang/platform/cpu                                           |2018-01-30|*
http://openconfig.net/yang/platform/extension                                     |2018-11-21| 
http://openconfig.net/yang/platform/fan                                           |2018-11-21| 
http://openconfig.net/yang/platform/linecard                                      |2018-11-21| 
http://openconfig.net/yang/platform/port                                          |2018-01-20|*
http://openconfig.net/yang/platform/psu                                           |2018-11-21|*
http://openconfig.net/yang/platform/transceiver                                   |2018-11-25|*
http://openconfig.net/yang/policy-types                                           |2017-07-14| 
http://openconfig.net/yang/rib/bgp                                                |2016-04-11|*
http://openconfig.net/yang/rib/bgp-types                                          |2016-04-11| 
http://openconfig.net/yang/routing-policy                                         |2018-06-05| 
http://openconfig.net/yang/rsvp-sr-ext                                            |2017-03-06|*
http://openconfig.net/yang/system                                                 |2018-07-17|*
http://openconfig.net/yang/system/logging                                         |2017-09-18| 
http://openconfig.net/yang/system/management                                      |2018-08-28| 
http://openconfig.net/yang/system/procmon                                         |2017-09-18| 
http://openconfig.net/yang/system/terminal                                        |2017-09-18| 
http://openconfig.net/yang/telemetry                                              |2016-02-04|*
http://openconfig.net/yang/transport-types                                        |2019-06-27| 
http://openconfig.net/yang/types/inet                                             |2017-08-24| 
http://openconfig.net/yang/types/yang                                             |2018-04-24| 
http://openconfig.net/yang/vlan                                                   |2016-05-26|*
http://openconfig.net/yang/vlan-types                                             |2016-05-26| 

```

## XR XML Config

```xml
RP/0/RP0/CPU0:iosxr# show running-config | xml openconfig
Sat Jan 15 04:05:12.617 UTC
Building configuration...
<data>
 <system xmlns="http://openconfig.net/yang/system">
  <aaa>
   <authentication>
    <users>
     <user>
      <username>root</username>
      <config>
       <username>root</username>
       <role>root-lr</role>
       <password-hashed>$6$XvYiZ/CdPNuK4Z/.$3/15yGC1Br2nlIy/AwZVNsl0BbD.XLbqAL2h8hR4CpBxM.ir4ZilYykiaTqMwe/EB6UySyH7ea/x09ajR6NXz.</password-hashed>
      </config>
     </user>
     <user>
      <username>admin</username>
      <config>
       <username>admin</username>
       <role>root-lr</role>
       <password-hashed>$6$vEaDc/Yt1OyU4c/.$v0lze75JluVDfcM6rgDlsFY3oMB6ODv6l5rgRnk3bFrvnzSFnoEIF.hcc1O/2.YxnAuRLSy7VQGmGedvoBlOp.</password-hashed>
      </config>
     </user>
    </users>
    <config>
     <authentication-method xmlns:idx="http://openconfig.net/yang/aaa/types">idx:LOCAL</authentication-method>
    </config>
   </authentication>
   <authorization>
    <config>
     <authorization-method xmlns:idx="http://openconfig.net/yang/aaa/types">idx:LOCAL</authorization-method>
    </config>
    <events>
     <event>
      <event-type xmlns:idx="http://openconfig.net/yang/aaa/types">idx:AAA_AUTHORIZATION_EVENT_CONFIG</event-type>
      <config>
       <event-type xmlns:idx="http://openconfig.net/yang/aaa/types">idx:AAA_AUTHORIZATION_EVENT_CONFIG</event-type>
      </config>
     </event>
    </events>
   </authorization>
  </aaa>
  <grpc-server>
   <config>
    <port>57777</port>
    <enable>true</enable>
   </config>
  </grpc-server>
  <config>
   <hostname>iosxr</hostname>
  </config>
  <ssh-server>
   <config>
    <enable>true</enable>
    <protocol-version>V2</protocol-version>
   </config>
  </ssh-server>
 </system>
 <network-instances xmlns="http://openconfig.net/yang/network-instance">
  <network-instance>
   <name>default</name>
   <protocols>
    <protocol>
     <identifier xmlns:idx="http://openconfig.net/yang/policy-types">idx:BGP</identifier>
     <name>default</name>
     <config>
      <identifier xmlns:idx="http://openconfig.net/yang/policy-types">idx:BGP</identifier>
      <name>default</name>
     </config>
     <bgp>
      <global>
       <config>
        <as>65503</as>
        <router-id>1.1.1.1</router-id>
       </config>
       <afi-safis>
        <afi-safi>
         <afi-safi-name xmlns:idx="http://openconfig.net/yang/bgp-types">idx:IPV4_UNICAST</afi-safi-name>
         <config>
          <afi-safi-name xmlns:idx="http://openconfig.net/yang/bgp-types">idx:IPV4_UNICAST</afi-safi-name>
          <enabled>true</enabled>
         </config>
        </afi-safi>
        <afi-safi>
         <afi-safi-name xmlns:idx="http://openconfig.net/yang/bgp-types">idx:IPV6_UNICAST</afi-safi-name>
         <config>
          <afi-safi-name xmlns:idx="http://openconfig.net/yang/bgp-types">idx:IPV6_UNICAST</afi-safi-name>
          <enabled>true</enabled>
         </config>
        </afi-safi>
       </afi-safis>
      </global>
      <peer-groups>
       <peer-group>
        <peer-group-name>sachin</peer-group-name>
        <config>
         <peer-group-name>sachin</peer-group-name>
        </config>
        <afi-safis>
         <afi-safi>
          <afi-safi-name xmlns:idx="http://openconfig.net/yang/bgp-types">idx:IPV4_UNICAST</afi-safi-name>
          <config>
           <afi-safi-name xmlns:idx="http://openconfig.net/yang/bgp-types">idx:IPV4_UNICAST</afi-safi-name>
           <enabled>true</enabled>
          </config>
         </afi-safi>
         <afi-safi>
          <afi-safi-name xmlns:idx="http://openconfig.net/yang/bgp-types">idx:IPV6_UNICAST</afi-safi-name>
          <config>
           <afi-safi-name xmlns:idx="http://openconfig.net/yang/bgp-types">idx:IPV6_UNICAST</afi-safi-name>
           <enabled>true</enabled>
          </config>
         </afi-safi>
        </afi-safis>
       </peer-group>
      </peer-groups>
      <neighbors>
       <neighbor>
        <neighbor-address>5.5.5.5</neighbor-address>
        <config>
         <neighbor-address>5.5.5.5</neighbor-address>
         <peer-group>sachin</peer-group>
        </config>
       </neighbor>
      </neighbors>
     </bgp>
    </protocol>
    <protocol>
     <identifier xmlns:idx="http://openconfig.net/yang/policy-types">idx:ISIS</identifier>
     <name>default</name>
     <config>
      <identifier xmlns:idx="http://openconfig.net/yang/policy-types">idx:ISIS</identifier>
      <name>default</name>
     </config>
     <isis>
      <global>
       <lsp-bit>
        <overload-bit>
         <config>
          <set-bit-on-boot>true</set-bit-on-boot>
         </config>
         <reset-triggers>
          <reset-trigger>
           <reset-trigger xmlns:idx="http://openconfig.net/yang/isis-types">idx:WAIT_FOR_SYSTEM</reset-trigger>
           <config>
            <reset-trigger xmlns:idx="http://openconfig.net/yang/isis-types">idx:WAIT_FOR_SYSTEM</reset-trigger>
            <delay>360</delay>
           </config>
          </reset-trigger>
         </reset-triggers>
        </overload-bit>
       </lsp-bit>
       <config>
        <level-capability>LEVEL_1</level-capability>
        <net>49.0001.1111.1111.0000.00</net>
       </config>
       <graceful-restart>
        <config>
         <enabled>true</enabled>
        </config>
       </graceful-restart>
       <timers>
        <lsp-generation>
         <config>
          <lsp-max-wait-interval>20</lsp-max-wait-interval>
          <lsp-first-wait-interval>1</lsp-first-wait-interval>
          <lsp-second-wait-interval>20</lsp-second-wait-interval>
         </config>
        </lsp-generation>
        <config>
         <lsp-refresh-interval>65000</lsp-refresh-interval>
        </config>
       </timers>
      </global>
     </isis>
    </protocol>
    <protocol>
     <identifier xmlns:idx="http://openconfig.net/yang/policy-types">idx:STATIC</identifier>
     <name>DEFAULT</name>
     <config>
      <identifier xmlns:idx="http://openconfig.net/yang/policy-types">idx:STATIC</identifier>
      <name>DEFAULT</name>
     </config>
     <static-routes>
      <static>
       <prefix>0.0.0.0/0</prefix>
       <config>
        <prefix>0.0.0.0/0</prefix>
       </config>
       <next-hops>
        <next-hop>
         <index>##10.10.20.254##</index>
         <config>
          <index>##10.10.20.254##</index>
          <next-hop>10.10.20.254</next-hop>
         </config>
        </next-hop>
       </next-hops>
      </static>
      <static>
       <prefix>10.230.0.0/16</prefix>
       <config>
        <prefix>10.230.0.0/16</prefix>
       </config>
       <next-hops>
        <next-hop>
         <index>##1.1.1.14##</index>
         <config>
          <index>##1.1.1.14##</index>
          <next-hop>1.1.1.14</next-hop>
         </config>
        </next-hop>
       </next-hops>
      </static>
      <static>
       <prefix>10.230.0.0/17</prefix>
       <config>
        <prefix>10.230.0.0/17</prefix>
       </config>
       <next-hops>
        <next-hop>
         <index>##1.1.1.14##</index>
         <config>
          <index>##1.1.1.14##</index>
          <next-hop>1.1.1.14</next-hop>
         </config>
        </next-hop>
        <next-hop>
         <index>##1.1.1.15##</index>
         <config>
          <index>##1.1.1.15##</index>
          <next-hop>1.1.1.15</next-hop>
         </config>
        </next-hop>
       </next-hops>
      </static>
      <static>
       <prefix>192.0.2.16/28</prefix>
       <config>
        <prefix>192.0.2.16/28</prefix>
        <set-tag>10</set-tag>
       </config>
       <next-hops>
        <next-hop>
         <index>##FastEthernet0/0/0/1##192.0.2.10##</index>
         <config>
          <index>##FastEthernet0/0/0/1##192.0.2.10##</index>
          <next-hop>192.0.2.10</next-hop>
         </config>
         <interface-ref>
          <config>
           <interface>FastEthernet0/0/0/1</interface>
           <subinterface>0</subinterface>
          </config>
         </interface-ref>
        </next-hop>
       </next-hops>
      </static>
     </static-routes>
    </protocol>
   </protocols>
   <config>
    <name>default</name>
   </config>
  </network-instance>
  <network-instance>
   <name>OM</name>
   <config>
    <name>OM</name>
    <enabled-address-families xmlns:idx="http://openconfig.net/yang/openconfig-types">idx:IPV4</enabled-address-families>
   </config>
  </network-instance>
 </network-instances>
 <routing-policy xmlns="http://openconfig.net/yang/routing-policy">
  <defined-sets>
   <prefix-sets>
    <prefix-set>
     <name>TEST</name>
     <config>
      <name>TEST</name>
      <mode>IPV4</mode>
     </config>
     <prefixes>
      <prefix>
       <ip-prefix>10.100.100.0/24</ip-prefix>
       <config>
        <ip-prefix>10.100.100.0/24</ip-prefix>
        <masklength-range>exact</masklength-range>
       </config>
       <masklength-range>exact</masklength-range>
      </prefix>
      <prefix>
       <ip-prefix>172.19.127.200/29</ip-prefix>
       <config>
        <ip-prefix>172.19.127.200/29</ip-prefix>
        <masklength-range>exact</masklength-range>
       </config>
       <masklength-range>exact</masklength-range>
      </prefix>
     </prefixes>
    </prefix-set>
    <prefix-set>
     <name>JAN06</name>
     <config>
      <name>JAN06</name>
      <mode>IPV4</mode>
     </config>
     <prefixes>
      <prefix>
       <ip-prefix>172.16.10.0/30</ip-prefix>
       <config>
        <ip-prefix>172.16.10.0/30</ip-prefix>
        <masklength-range>exact</masklength-range>
       </config>
       <masklength-range>exact</masklength-range>
      </prefix>
      <prefix>
       <ip-prefix>10.10.10.0/24</ip-prefix>
       <config>
        <ip-prefix>10.10.10.0/24</ip-prefix>
        <masklength-range>exact</masklength-range>
       </config>
       <masklength-range>exact</masklength-range>
      </prefix>
     </prefixes>
    </prefix-set>
    <prefix-set>
     <name>vrf-red</name>
     <config>
      <name>vrf-red</name>
      <mode>IPV4</mode>
     </config>
     <prefixes>
      <prefix>
       <ip-prefix>1.1.1.1/32</ip-prefix>
       <config>
        <ip-prefix>1.1.1.1/32</ip-prefix>
        <masklength-range>exact</masklength-range>
       </config>
       <masklength-range>exact</masklength-range>
      </prefix>
     </prefixes>
    </prefix-set>
   </prefix-sets>
  </defined-sets>
 </routing-policy>
 <interfaces xmlns="http://openconfig.net/yang/interfaces">
  <interface>
   <name>BVI220</name>
   <config>
    <name>BVI220</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:propVirtual</type>
    <enabled>false</enabled>
   </config>
  </interface>
  <interface>
   <name>BVI221</name>
   <config>
    <name>BVI221</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:propVirtual</type>
    <enabled>false</enabled>
   </config>
  </interface>
  <interface>
   <name>BVI225</name>
   <config>
    <name>BVI225</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:propVirtual</type>
    <enabled>false</enabled>
   </config>
  </interface>
  <interface>
   <name>BVI991</name>
   <config>
    <name>BVI991</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:propVirtual</type>
    <description>DNP | TEST_MGMT</description>
   </config>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>172.19.127.201</ip>
        <config>
         <ip>172.19.127.201</ip>
         <prefix-length>29</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>Loopback0</name>
   <config>
    <name>Loopback0</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
   </config>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>1.1.1.1</ip>
        <config>
         <ip>1.1.1.1</ip>
         <prefix-length>32</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>Loopback10</name>
   <config>
    <name>Loopback10</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
   </config>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>192.168.10.1</ip>
        <config>
         <ip>192.168.10.1</ip>
         <prefix-length>24</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>Loopback20</name>
   <config>
    <name>Loopback20</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
    <description>another_testing</description>
   </config>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>192.168.20.1</ip>
        <config>
         <ip>192.168.20.1</ip>
         <prefix-length>24</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>Loopback25</name>
   <config>
    <name>Loopback25</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
    <description>my loopback interface</description>
   </config>
  </interface>
  <interface>
   <name>Loopback99</name>
   <config>
    <name>Loopback99</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
    <description>loopback99 on iosxr</description>
   </config>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>192.168.0.99</ip>
        <config>
         <ip>192.168.0.99</ip>
         <prefix-length>24</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>Loopback100</name>
   <config>
    <name>Loopback100</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
    <description>***MERGE LOOPBACK 100****</description>
   </config>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>1.1.1.100</ip>
        <config>
         <ip>1.1.1.100</ip>
         <prefix-length>32</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>Loopback101</name>
   <config>
    <name>Loopback101</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
   </config>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>192.168.101.1</ip>
        <config>
         <ip>192.168.101.1</ip>
         <prefix-length>24</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>Loopback103</name>
   <config>
    <name>Loopback103</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
   </config>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>192.168.103.1</ip>
        <config>
         <ip>192.168.103.1</ip>
         <prefix-length>24</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>Loopback111</name>
   <config>
    <name>Loopback111</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
    <enabled>false</enabled>
    <description>test loopback111</description>
   </config>
  </interface>
  <interface>
   <name>Loopback112</name>
   <config>
    <name>Loopback112</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
    <description>Configured by NETCONF Invalid</description>
   </config>
  </interface>
  <interface>
   <name>Loopback123</name>
   <config>
    <name>Loopback123</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
    <description>test</description>
   </config>
  </interface>
  <interface>
   <name>Loopback200</name>
   <config>
    <name>Loopback200</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
    <description>***MERGE LOOPBACK 200****</description>
   </config>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>1.1.1.200</ip>
        <config>
         <ip>1.1.1.200</ip>
         <prefix-length>32</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>Loopback555</name>
   <config>
    <name>Loopback555</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
    <description>Configured by Salt-Nornir using NETCONF</description>
   </config>
  </interface>
  <interface>
   <name>Loopback1000</name>
   <config>
    <name>Loopback1000</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
   </config>
  </interface>
  <interface>
   <name>Loopback1234</name>
   <config>
    <name>Loopback1234</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:softwareLoopback</type>
   </config>
  </interface>
  <interface>
   <name>tunnel-ip12445</name>
   <config>
    <name>tunnel-ip12445</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:tunnel</type>
    <description>Testsimple2-39533-1044</description>
   </config>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>10.220.45.41</ip>
        <config>
         <ip>10.220.45.41</ip>
         <prefix-length>30</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
    </subinterface>
   </subinterfaces>
   <tunnel xmlns="http://openconfig.net/yang/interfaces/tunnel">
    <ipv4>
     <addresses>
      <address>
       <ip>10.220.45.41</ip>
       <config>
        <ip>10.220.45.41</ip>
        <prefix-length>30</prefix-length>
       </config>
      </address>
     </addresses>
    </ipv4>
    <config>
     <src>185.121.241.14</src>
     <dst>38.135.71.23</dst>
    </config>
   </tunnel>
  </interface>
  <interface>
   <name>Bundle-Ether165</name>
   <config>
    <name>Bundle-Ether165</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:ieee8023adLag</type>
    <enabled>false</enabled>
   </config>
  </interface>
  <interface>
   <name>Bundle-Ether167</name>
   <config>
    <name>Bundle-Ether167</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:ieee8023adLag</type>
    <enabled>false</enabled>
   </config>
  </interface>
  <interface>
   <name>Bundle-Ether175</name>
   <config>
    <name>Bundle-Ether175</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:ieee8023adLag</type>
    <enabled>false</enabled>
   </config>
  </interface>
  <interface>
   <name>Bundle-Ether177</name>
   <config>
    <name>Bundle-Ether177</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:ieee8023adLag</type>
    <enabled>false</enabled>
   </config>
  </interface>
  <interface>
   <name>MgmtEth0/RP0/CPU0/0</name>
   <config>
    <name>MgmtEth0/RP0/CPU0/0</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:ethernetCsmacd</type>
   </config>
   <ethernet xmlns="http://openconfig.net/yang/interfaces/ethernet">
    <config>
     <auto-negotiate>false</auto-negotiate>
    </config>
   </ethernet>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>10.10.20.175</ip>
        <config>
         <ip>10.10.20.175</ip>
         <prefix-length>24</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>GigabitEthernet0/0/0/0</name>
   <config>
    <name>GigabitEthernet0/0/0/0</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:ethernetCsmacd</type>
    <enabled>false</enabled>
   </config>
   <ethernet xmlns="http://openconfig.net/yang/interfaces/ethernet">
    <config>
     <auto-negotiate>false</auto-negotiate>
    </config>
   </ethernet>
   <subinterfaces>
    <subinterface>
     <index>123</index>
     <config>
      <index>123</index>
     </config>
     <vlan xmlns="http://openconfig.net/yang/vlan">
      <config>
       <vlan-id>123</vlan-id>
      </config>
     </vlan>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>GigabitEthernet0/0/0/1</name>
   <config>
    <name>GigabitEthernet0/0/0/1</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:ethernetCsmacd</type>
    <description>DNP</description>
   </config>
   <ethernet xmlns="http://openconfig.net/yang/interfaces/ethernet">
    <config>
     <auto-negotiate>false</auto-negotiate>
    </config>
   </ethernet>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>10.100.100.1</ip>
        <config>
         <ip>10.100.100.1</ip>
         <prefix-length>24</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>GigabitEthernet0/0/0/2</name>
   <config>
    <name>GigabitEthernet0/0/0/2</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:ethernetCsmacd</type>
   </config>
   <ethernet xmlns="http://openconfig.net/yang/interfaces/ethernet">
    <config>
     <auto-negotiate>false</auto-negotiate>
    </config>
   </ethernet>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>192.168.40.1</ip>
        <config>
         <ip>192.168.40.1</ip>
         <prefix-length>24</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>GigabitEthernet0/0/0/3</name>
   <config>
    <name>GigabitEthernet0/0/0/3</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:ethernetCsmacd</type>
    <enabled>false</enabled>
   </config>
   <ethernet xmlns="http://openconfig.net/yang/interfaces/ethernet">
    <config>
     <auto-negotiate>false</auto-negotiate>
    </config>
   </ethernet>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>192.168.40.1</ip>
        <config>
         <ip>192.168.40.1</ip>
         <prefix-length>24</prefix-length>
        </config>
       </address>
       <address>
        <ip>192.168.2.1</ip>
        <config>
         <ip>192.168.2.1</ip>
         <prefix-length>24</prefix-length>
        </config>
       </address>
       <address>
        <ip>192.168.30.1</ip>
        <config>
         <ip>192.168.30.1</ip>
         <prefix-length>24</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>GigabitEthernet0/0/0/4</name>
   <config>
    <name>GigabitEthernet0/0/0/4</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:ethernetCsmacd</type>
    <description>test</description>
   </config>
   <ethernet xmlns="http://openconfig.net/yang/interfaces/ethernet">
    <config>
     <auto-negotiate>false</auto-negotiate>
    </config>
   </ethernet>
   <subinterfaces>
    <subinterface>
     <index>0</index>
     <ipv4 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>192.168.20.1</ip>
        <config>
         <ip>192.168.20.1</ip>
         <prefix-length>24</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv4>
     <ipv6 xmlns="http://openconfig.net/yang/interfaces/ip">
      <addresses>
       <address>
        <ip>2001:db8:ff::11</ip>
        <config>
         <ip>2001:db8:ff::11</ip>
         <prefix-length>64</prefix-length>
        </config>
       </address>
      </addresses>
     </ipv6>
    </subinterface>
   </subinterfaces>
  </interface>
  <interface>
   <name>GigabitEthernet0/0/0/6</name>
   <config>
    <name>GigabitEthernet0/0/0/6</name>
    <type xmlns:idx="urn:ietf:params:xml:ns:yang:iana-if-type">idx:ethernetCsmacd</type>
   </config>
   <ethernet xmlns="http://openconfig.net/yang/interfaces/ethernet">
    <config>
     <auto-negotiate>false</auto-negotiate>
    </config>
   </ethernet>
  </interface>
 </interfaces>
 <acl xmlns="http://openconfig.net/yang/acl">
  <acl-sets>
   <acl-set>
    <name>acl_1</name>
    <type>ACL_IPV4</type>
    <config>
     <name>acl_1</name>
     <type>ACL_IPV4</type>
    </config>
    <acl-entries>
     <acl-entry>
      <sequence-id>16</sequence-id>
      <config>
       <sequence-id>16</sequence-id>
      </config>
     </acl-entry>
     <acl-entry>
      <sequence-id>21</sequence-id>
      <config>
       <sequence-id>21</sequence-id>
      </config>
      <actions>
       <config>
        <forwarding-action>ACCEPT</forwarding-action>
       </config>
      </actions>
      <ipv4>
       <config>
        <protocol>6</protocol>
        <source-address>192.0.2.10/32</source-address>
        <destination-address>198.51.100.0/28</destination-address>
       </config>
      </ipv4>
      <transport>
       <config>
        <source-port>110..121</source-port>
        <tcp-flags xmlns:idx="http://openconfig.net/yang/packet-match-types">idx:TCP_RST</tcp-flags>
       </config>
      </transport>
     </acl-entry>
     <acl-entry>
      <sequence-id>23</sequence-id>
      <config>
       <sequence-id>23</sequence-id>
      </config>
      <actions>
       <config>
        <forwarding-action>REJECT</forwarding-action>
       </config>
      </actions>
      <ipv4>
       <config>
        <protocol>1</protocol>
        <destination-address>198.51.100.0/28</destination-address>
       </config>
      </ipv4>
     </acl-entry>
    </acl-entries>
   </acl-set>
   <acl-set>
    <name>acl_2</name>
    <type>ACL_IPV4</type>
    <config>
     <name>acl_2</name>
     <type>ACL_IPV4</type>
    </config>
    <acl-entries>
     <acl-entry>
      <sequence-id>10</sequence-id>
      <config>
       <sequence-id>10</sequence-id>
      </config>
     </acl-entry>
    </acl-entries>
   </acl-set>
  </acl-sets>
 </acl>
</data>
```