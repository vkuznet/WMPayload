// handlers.go
package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func insertHandler(ctx *fasthttp.RequestCtx) {
	//     var records []map[string]interface{}
	var records []interface{}
	// Determine the MIME type
	contentType := string(ctx.Request.Header.ContentType())

	if contentType == "application/ndjson" {
		records = parseNDJSON(ctx.PostBody())
	} else {
		if err := json.Unmarshal(ctx.PostBody(), &records); err != nil {
			ctx.Error("Failed to parse JSON", fasthttp.StatusBadRequest)
			return
		}
	}

	collection := mongoDB.Collection(config.MongoCollection)
	_, err := collection.InsertMany(context.Background(), records)
	if err != nil {
		ctx.Error("Failed to insert documents", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Documents inserted successfully")
}

func searchHandler(ctx *fasthttp.RequestCtx) {
	idx := ctx.QueryArgs().GetUintOrZero("idx")
	limit := ctx.QueryArgs().GetUintOrZero("limit")
	uuid := string(ctx.QueryArgs().Peek("id"))

	collection := mongoDB.Collection(config.MongoCollection)
	spec := bson.M{"id": uuid}
	cursor, err := collection.Find(context.Background(), spec, options.Find().SetSkip(int64(idx)).SetLimit(int64(limit)))
	if err != nil {
		ctx.Error("Failed to retrieve documents", fasthttp.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var records []map[string]interface{}
	if err := cursor.All(context.Background(), &records); err != nil {
		ctx.Error("Failed to parse documents", fasthttp.StatusInternalServerError)
		return
	}

	/*
		responseBody, err := json.Marshal(records)
		if err != nil {
			ctx.Error("Failed to serialize response", fasthttp.StatusInternalServerError)
			return
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody(responseBody)
	*/

	contentType := string(ctx.Request.Header.Peek("Accept"))
	if contentType == "application/ndjson" {
		ctx.SetContentType("application/ndjson")
		for _, record := range records {
			line, _ := json.Marshal(record)
			ctx.Write(line)
			ctx.Write([]byte("\n"))
		}
	} else {
		ctx.SetContentType("application/json")
		json.NewEncoder(ctx).Encode(records)
	}
}

func parseNDJSON(data []byte) []interface{} {
	var records []interface{}
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var record map[string]interface{}
		if err := json.Unmarshal(line, &record); err == nil {
			records = append(records, record)
		}
	}
	return records
}
