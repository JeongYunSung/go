package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"protobuf/pb"
)

func main() {
	p := &pb.Person{
		Id:    1234,
		Name:  "JeongYunSung",
		Email: "123dbstjd@gmail.com",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "010-1234-5678", Type: pb.Person_MOBILE},
		},
	}

	book := &pb.AddressBook{
		People: []*pb.Person{p},
	}

	out, err := proto.Marshal(book)
	if err != nil {
		log.Fatalln("Failed to encode address book:", err)
	}

	readBook := &pb.AddressBook{}
	if err := proto.Unmarshal(out, readBook); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	fmt.Printf("default : %v\nmarshal : %v\nunmarshal : %v\n", p, out, readBook)
}
