# pathwar
Pathwar monorepo

## API Documentation

* SwaggerHub: https://app.swaggerhub.com/apis/Pathwar/Pathwar/0
* Protobuf file: https://github.com/pathwar/pathwar/blob/master/server/server.proto
* Available both using gRPC and HTTP (via https://github.com/grpc-ecosystem/grpc-gateway)
* GoDoc: [![GoDoc](https://godoc.org/pathwar.pw?status.svg)](https://godoc.org/pathwar.pw)

## Database schemas (expected)

```
               ┌────────────┐    ┌─────────────┐┌────────┐       ┌──────┐                              
               │ AuthMethod │    │ Achievement ││ Coupon │       │ Team │                              
               └────────────┘    └─────────────┘└────────┘       └──────┘                              
                      │                 │            │               │                                 
┌────────────┐        │                 │            │               │                                 
│ Credential │────┌──────┐              ┌────────────┐      ┌────────────────┐                         
└────────────┘    │ User │──────────────│ TeamMember │──────│ TournamentTeam │─────────────────┐       
          ┌───────└──────┘              └────────────┘      └────────────────┘                 │       
          │           │                 │                   │                │                 │       
  ┌──────────────┐    │                 │                   │                │                 │       
  │ Notification │    │                 │                   │                │                 │       
  └──────────────┘    │                 │                   │                │                 │       
               ┌─────────────┐ ┌────────────────┐    ┌────────────┐┌───────────────────┐ ┌──────────┐  
               │ UserSession │ │ WhoswhoAttempt │    │ Tournament ││ LevelSubscription │ │ ShopItem │  
               └─────────────┘ └────────────────┘    └────────────┘└───────────────────┘ └──────────┘  
                                                                             │                 │       
                                                                             │                 │       
       ┌───────┐                                                             │                 │       
       │ Event │                 ┌────────────┐   ┌───────────────┐   ┌─────────────┐   ┌─────────────┐
       └───────┘                 │ Hypervisor │───│ LevelInstance │───│ LevelFlavor │───│    Level    │
                                 └────────────┘   └───────────────┘   └─────────────┘   └─────────────┘
```

## Production architecture (expected)

```
                                ┌─────────────────────────────────────┐
                                │       pathwar server cluster        │
                                │                                     │
                                │ ┌─────────────────────────────────┐ │
                                │ │┌─────────────┐                  │ │
                                │ ││             │                  │ │
                                │ ││  ssh proxy  │                  │ │
                                │ ││             │                  │ │
                                │ │└─────────────┘                  │ │
                                │ │┌─────────────┐                  │ │
                                │ ││             │                  │ │
                                │ ││     web     │                  │ │
                                │ ││             │                  │ │
                                │ │└─────────────┘    pathwar server│ │
                                │ │┌─────────────┐                  │ │
                ┌───────────┐   │ ││             │                  │ │   ┌─────────┐
┌───────────┐   │           │   │ ││ http proxy  │                  │ │   │         │
│           │   │           │   │ ││             │                  │ │   │   SQL   │
│   users   │──▶│  haproxy  │──▶│ │└─────────────┘                  │ │──▶│ cluster │
│           │   │           │   │ │┌─────────────┐                  │ │   │         │
└───────────┘   │           │   │ ││             │                  │ │   └─────────┘
                └───────────┘   │ ││     api     │                  │ │
                                │ ││             │                  │ │
                                │ │└─────────────┘                  │ │
                                │ └─────────────────────────────────┘ │
                                │ ┌─────────────────────────────────┐ │
                                │ │                                 │ │
                                │ │                   pathwar server│ │
                                │ │                                 │ │
                                │ └─────────────────────────────────┘ │
                                │ ┌─────────────────────────────────┐ │
                                │ │                                 │ │
                                │ │                              ...│ │
                                │ │                                 │ │
                                │ └─────────────────────────────────┘ │
                                └─────────────────────────────────────┘
```
