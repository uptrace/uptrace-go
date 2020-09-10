package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongotrace "go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver"
)

var upclient *uptrace.Client

func main() {
	ctx := context.Background()

	upclient = setupUptrace()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	upclient.ReportError(ctx, errors.New("hello from uptrace-go!"))

	opt := options.Client()
	opt.Monitor = mongotrace.NewMonitor("mongo-service")
	opt.ApplyURI("mongodb://mongo-server:27017")

	mdb, err := mongo.Connect(ctx, opt)
	if err != nil {
		upclient.ReportError(ctx, err)
		panic(err)
	}

	if err := mdb.Ping(ctx, nil); err != nil {
		upclient.ReportError(ctx, err)
		panic(err)
	}

	if err := run(ctx, mdb.Database("example")); err != nil {
		upclient.ReportError(ctx, err)
		panic(err)
	}
}

func setupUptrace() *uptrace.Client {
	log.Printf("using UPTRACE_DSN=%q", os.Getenv("UPTRACE_DSN"))

	hostname, _ := os.Hostname()
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		Resource: map[string]interface{}{
			"hostname": hostname,
		},
	})

	return upclient
}

// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

func run(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection("inventory_insert")

	_, err := coll.InsertOne(
		ctx,
		bson.D{
			{"item", "canvas"},
			{"qty", 100},
			{"tags", bson.A{"cotton"}},
			{"size", bson.D{
				{"h", 28},
				{"w", 35.5},
				{"uom", "cm"},
			}},
		})
	if err != nil {
		return err
	}

	_, err = coll.Find(
		ctx,
		bson.D{{"item", "canvas"}},
	)
	if err != nil {
		return err
	}

	_, err = coll.InsertMany(
		ctx,
		[]interface{}{
			bson.D{
				{"item", "journal"},
				{"qty", int32(25)},
				{"tags", bson.A{"blank", "red"}},
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
			},
			bson.D{
				{"item", "mat"},
				{"qty", int32(25)},
				{"tags", bson.A{"gray"}},
				{"size", bson.D{
					{"h", 27.9},
					{"w", 35.5},
					{"uom", "cm"},
				}},
			},
			bson.D{
				{"item", "mousepad"},
				{"qty", 25},
				{"tags", bson.A{"gel", "blue"}},
				{"size", bson.D{
					{"h", 19},
					{"w", 22.85},
					{"uom", "cm"},
				}},
			},
		})
	if err != nil {
		return err
	}

	return nil
}
