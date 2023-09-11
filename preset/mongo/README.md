# Gnomock MongoDB

Gnomock MongoDB is a [Gnomock](https://github.com/sandwich-go/gnomock) preset for
running tests against a real MongoDB container, without mocks.

```go
package mongo_test

import (
	"context"
	"fmt"

	"github.com/sandwich-go/gnomock"
	"github.com/sandwich-go/gnomock/preset/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
)

func ExamplePreset() {
	p := mongo.Preset(
		mongo.WithData("./testdata/"),
		mongo.WithUser("gnomock", "gnomick"),
	)
	c, err := gnomock.Start(p)

	defer func() { _ = gnomock.Stop(c) }()

	if err != nil {
		panic(err)
	}

	addr := c.DefaultAddress()
	uri := fmt.Sprintf("mongodb://%s:%s@%s", "gnomock", "gnomick", addr)
	clientOptions := mongooptions.Client().ApplyURI(uri)

	client, err := mongodb.NewClient(clientOptions)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	// see testdata folder to verify names/numbers
	fmt.Println(client.Database("db1").Collection("users").CountDocuments(ctx, bson.D{}))
	fmt.Println(client.Database("db2").Collection("customers").CountDocuments(ctx, bson.D{}))
	fmt.Println(client.Database("db2").Collection("countries").CountDocuments(ctx, bson.D{}))

	// Output:
	// 10 <nil>
	// 5 <nil>
	// 3 <nil>
}
```
