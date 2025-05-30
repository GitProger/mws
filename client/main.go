package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	client "mws/gen_api"
)

func add(ctx context.Context, c *client.Client, example *client.Book, userID int) {
	if addedBook, err := c.AddUserBook(ctx, example, client.AddUserBookParams{UserID: userID}); err != nil {
		log.Panic(err)
	} else {
		fmt.Print("Book added: ")
		json.NewEncoder(os.Stdout).Encode(addedBook)
	}
}

func update(ctx context.Context, c *client.Client, userID, bookID, page int) {
	if book, err := c.UpdateReadingProgress(ctx,
		&client.UpdateReadingProgressReq{Page: page},
		client.UpdateReadingProgressParams{UserID: userID, BookID: bookID}); err != nil {
		log.Panic(err)
	} else {
		fmt.Printf("Page updated: %d\n", book.(*client.Book).Page)
	}
}

func get(ctx context.Context, c *client.Client, userID, bookID int) {
	if book, err := c.GetUserBook(ctx, client.GetUserBookParams{UserID: userID, BookID: bookID}); err != nil {
		log.Panic(err)
	} else {
		fmt.Printf("%d's Book %d: ", userID, bookID)
		json.NewEncoder(os.Stdout).Encode(book)
	}
}

func remove(ctx context.Context, c *client.Client, userID, bookID int) {
	if _, err := c.RemoveUserBook(ctx, client.RemoveUserBookParams{UserID: userID, BookID: bookID}); err != nil {
		log.Panic(err)
	} else {
		fmt.Println("Book removed")
	}
}

func list(ctx context.Context, c *client.Client, userID int) {
	if books, err := c.GetUserBooks(ctx, client.GetUserBooksParams{UserID: userID}); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Books:")
		for _, b := range books {
			fmt.Printf(" - '%s' (page %d)\n", b.Title, b.Page)
		}
	}
}

func test() {
	c, err := client.NewClient("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	userID := 1
	bookID := 1234

	date, _ := time.Parse(time.DateOnly, "1957-11-23")
	example := &client.Book{
		ID:        bookID,
		Title:     "Доктор Живаго",
		Author:    "Борис Пастернак",
		Published: date,
		Page:      2,
	}

	add(ctx, c, example, userID)
	list(ctx, c, userID)
	update(ctx, c, userID, bookID, 25)
	get(ctx, c, userID, bookID)
	remove(ctx, c, userID, bookID)

	if res, err := c.GetUserBook(ctx, client.GetUserBookParams{UserID: userID, BookID: bookID}); err != nil {
		log.Fatal(err)
	} else {
		json.NewEncoder(os.Stdout).Encode(res.(*client.Error)) // явно сконвертим
	}
}

func printHelp() {
	fmt.Println(`Available commands:
    help                        - show this help
    exit                        - exit program
    list <userID>               - list user's books
    get <userID> <bookID>       - get book info
    remove <userID> <bookID>    - remove book
    update <userID> <bookID> <page> - update reading progress
    add <userID> <title>        - add new book`)
}

func interactive() {
	serv, err := client.NewClient("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	parse := func(expect string, args []string, types string) ([]any, bool) {
		if len(args) < len(types) {
			fmt.Println(expect)
			return nil, false
		}
		parsed := make([]any, 0, len(args))
		for i, typ := range types {
			if typ == 'i' {
				p, err := strconv.Atoi(args[i])
				if err != nil {
					fmt.Println(expect)
					return nil, false
				}
				parsed = append(parsed, p)
			} else {
				parsed = append(parsed, args[i])
			}
		}
		return parsed, true
	}

	input := func() bool {
		fmt.Print("> ")
		return scanner.Scan()
	}

	for input() {
		line := scanner.Text()
		cmd, argStr, _ := strings.Cut(line, " ")
		args := strings.Split(argStr, " ")

		if func() (exit bool) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("error:", r)
				}
			}()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			switch cmd {
			case "help":
				printHelp()
			case "exit":
				exit = true
			case "add":
				args := strings.SplitN(argStr, " ", 4)
				if args, ok := parse("wrong format, expected: add <userID> <bookID> <book title> <author name>", args, "iiss"); ok {
					book := &client.Book{Page: 1, ID: args[1].(int), Title: args[2].(string), Author: args[2].(string), Published: time.Now()}
					add(ctx, serv, book, args[0].(int))
				}
			case "list":
				if args, ok := parse("wrong format, expected: list <userID>", args, "i"); ok {
					list(ctx, serv, args[0].(int))
				}
			case "get":
				if args, ok := parse("wrong format, expected: get <userID> <bookID>", args, "ii"); ok {
					get(ctx, serv, args[0].(int), args[1].(int))
				}
			case "remove":
				if args, ok := parse("wrong format, expected: remove <userID> <bookID>", args, "ii"); ok {
					remove(ctx, serv, args[0].(int), args[1].(int))
				}
			case "update":
				if args, ok := parse("wrong format, expected: update <userID> <bookID> <page>", args, "iii"); ok {
					update(ctx, serv, args[0].(int), args[1].(int), args[2].(int))
				}
			default:
				printHelp()
			}
			return
		}() {
			break
		}
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "i" {
		interactive()
	} else {
		test()
	}
}
