# Gin-Swisspass


Gin-Swisspass is specially made for [Gin Framework](https://github.com/gin-gonic/gin)

## Project Context and Features

When it comes to choosing a Go framework, there's a lot of confusion
about what to use. The scene is very fragmented, and detailed
comparisons of different frameworks are still somewhat rare. Meantime,
how to handle dependencies and structure projects are big topics in
the Go community. We've liked using Gin for its speed,
accessibility, and usefulness in developing microservice
architectures. In creating Gin-Swisspass, we wanted to take fuller
advantage of Gin's capabilities and help other devs do likewise.

Gin-Swisspass is expressive, flexible, and very easy to use. It allows you to:
- do authentication and authorization based on the Swisspass Token
- create router groups to place Swisspass authorization on top, using HTTP verbs and passing them
- more easily decouple services by promoting a "say what to do, not how to do it" approach
- configure your REST API directly in the code (see the "Usage" example below)


## Requirements

- [Gin](https://github.com/gin-gonic/gin)
- Swisspass credentials

## Installation

Assuming you've installed Go and Gin, run this, or just import it and let go modules do the rest

    go get github.com/schweizerischebundesbahnen/gin-swisspass
    

## Usage

First create a GinSwisspass instance. You have 3 Options to define where your clientId should come from:
FromString, or from a fixed HTTP Request Header Key or by a custom Header Key

    ginSwisspass := gin_swisspass.New("https://www-test.swisspass.ch", gin_swisspass.ClientIdFromString("<your fix client id >"), <timeout in seconds>)


### Swisspass Authorization

Allow all authenticated users

    privateRoute := router.Group("/api/private")
    
    privateRoute.Use(gin_swisspass.
            NewBuilder(ginSwisspass).
            Build()
    )
           
Allow all authenticated users with Role "ADMIN" 

    privateRoute.Use(gin_swisspass.
        NewBuilder(ginSwisspass).
        AllowRole("ADMIN").
        Build()
    )

Allow all authenticated users with Role "ADMIN" and User with SwisspassId 123

    privateRoute.Use(gin_swisspass.
        NewBuilder(ginSwisspass).
        AllowRole("ADMIN").
        AllowSwisspassId("123")
        Build()
    )   
   
Finally, define your routes

	privateRoute.GET("/", func(c *gin.Context) {
		....
	})



#### Access the userInfo

After authentication the userInfo struct can be accesst via the gin Context in any route

    userInfo,_ := ctx.Get(gin_swisspass.CONTEXT_USERINFO_KEY)

#### Testing

To test, you can use curl:

        curl -H "Authorization: Bearer $TOKEN" http://localhost:8081/api/private/
        {"message":"Hello from private for users to Sandor Szücs"}

