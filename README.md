# Ghost 👻

![Version](https://img.shields.io/badge/Release-v0.0.1-brightgreen?style=flat-square)
![License](https://img.shields.io/:License-MIT-green.svg?style=flat-square)

A simple HTTP server framework in Go, without any third-party dependencies.

Use it ONLY when you seek for simpleness, and do not extremely eager for performance and robustness.

## Usage

All you need to do, is to make an interesting ghost, then run it.

```go
package main

import (
    "github.com/deadblue/ghost"
    "github.com/deadblue/ghost/view"
)

type YourGhost struct{}

// Get will handle request "GET /"
func (g *YourGhost) Get(_ ghost.Context) (ghost.View, error) {
    return view.Text("You are here!"), nil
}

// GetIndexAsHtml will handle request "GET /index.html"
func (g *YourGhost) GetIndexAsHtml(_ ghost.Context) (ghost.View, error) {
    return view.Text("index.html"), nil
}

// GetDataById will handle request "GET /data/{id}", where the "id" is a path variable.
func (g *YourGhost) GetDataById(ctx ghost.Context) (ghost.View, error) {
    dataId := ctx.PathVar("id")
    return view.Text("Getting data whose id is " + dataId), nil
}

// PostUpdate will handle request "POST /update" 
func (g *YourGhost) PostUpdate(ctx ghost.Context) (ghost.View, error) {
    // Get post data from ctx.Request()
    return view.Text("Update done!"), nil
}

// BuyMeACoffee will handle request "BUY /me/a/coffee"
func (g *YourGhost) BuyMeACoffee(_ ghost.Context) (ghost.View, error) {
    return view.Text("Thank you!"), nil
}

// Implement ghost.Binder interface, to speed up the controller invoking.
func (g *YourGhost) Bind(v interface{}) ghost.Controller {
    if fn, ok := v.(func(*YourGhost, ghost.Context)(ghost.View, error)); ok {
        return func(ctx ghost.Context)(ghost.View, error) {
            return fn(g, ctx)
        }
    } else {
        return nil
    }
}

// Implement ghost.StartupHandler interface, to initialize your ghost before shell running.
func (g *YourGhost) OnStartup() error {
    // Initializing  ...
    return nil
}

// Implement ghost.ShutdownHandler interface, to finalize your ghost after shell shutdown.
func (g *YourGhost) OnShutdown() error { 
    // Finalizing ...
    return nil
}

func main() {
    err := ghost.Born(&YourGhost{}).Run()
    if err != nil {
        panic(err)
    }
}
```

## Mechanism

Each method on the ghost object, which is in form of the `ghost.Controller`, will be mounted as a request handler. 

The method name, will be translated as the mount path, following these rules:

* Suppose the method name is in camel-case, split it into words.
* The first word will be treated as request method.
* If there is no remain words, the method function will handle the root request.
* If there are remain words, each word will be treated as a path segment.
* In remain words, there are some special words that won't be treated as path segment:
  * `By`: The next word will be treated as a path variable.
  * `As`: Link the next word with "." instead of path separator ("/").

For examples:

| Method Name       | Handle Request
|-------------------|---------------
| Get               | GET /
| GetIndex          | GET /index
| GetUserProfile    | GET /user/profile
| PostUserProfile   | POST /user/profile
| GetDataById       | GET /data/{id}
| GetDataByTypeById | GET /data/{type}/{id}
| GetDataByIdAsJson | GET /data/{id}.json
| BuyMeACoffee      | BUY /me/a/coffee

## License

MIT
