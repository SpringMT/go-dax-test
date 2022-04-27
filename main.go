package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-dax-go/dax"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"time"
	"log"
)

var endpoint = flag.String("endpoint", "", "dax cluster endpoint")

const (
	table      = "TestDax"
)

func main() {
	flag.Parse()
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")},
	)
	if err != nil {
		log.Fatal(err)
	}
	ddb := dynamodb.New(sess)
	cfg := dax.DefaultConfig()
	cfg.HostPorts = []string{ *endpoint }
	cfg.Region = "ap-northeast-1"
	log.Println(cfg)
	dax, _ := dax.New(cfg)

	u, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}
	now := time.Now()
	// itemをputしてDAX経由でGET
	item := map[string]*dynamodb.AttributeValue{
		"ID":    {S: aws.String(u.String())},
		"value": {S: aws.String(fmt.Sprintf("%s_%s", u.String(), now))},
	}
	log.Println(item)
	in := &dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item:      item,
	}
	out, err := ddb.PutItem(in)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(out)

	key := map[string]*dynamodb.AttributeValue{
		"pk": {S: aws.String(u.String())},
	}
	inDax := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key:       key,
	}
	// https://docs.aws.amazon.com/ja_jp/amazondynamodb/latest/developerguide/DAX.html
	outDax, errDax := dax.GetItem(inDax)
	if errDax != nil {
		log.Fatal(errDax)
	}
	log.Println(outDax)
}

