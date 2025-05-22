package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	api "mws/gen_api"
)

type serviceImpl struct {
	mu sync.RWMutex

	// если использовать map[int]map[int]*Book
	// то в мапе будут ссылки на книги с разных кусков памяти, которые были выделены где-то в хендлерах при парсинге запросов в собственно Book{}
	// это будет приводить к куче индерекций и кэш миссов
	// здесь они хранятся в более-менее непрерывном участке памяти, так как при переаллокации мапы
	// они все будут лежать в выделенном протяженном участке
	users map[int]map[int]api.Book
}

func newServiceImpl() *serviceImpl {
	return &serviceImpl{
		users: make(map[int]map[int]api.Book),
	}
}

func err(code int, format string, args ...any) *api.Error {
	return &api.Error{
		StatusCode: code,
		Message:    fmt.Sprintf(format, args...),
	}
}

func (s *serviceImpl) GetUserBooks(ctx context.Context, params api.GetUserBooksParams) ([]api.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if books, ok := s.users[params.UserID]; ok {
		values := make([]api.Book, 0, len(books))
		for _, book := range books {
			values = append(values, book)
		}
		return values, nil
	} else { // у пользователя нет книг, либо по хорошему надо отдельно проверять есть ли такой пользователь
		return []api.Book{}, nil
	}
}

func (s *serviceImpl) AddUserBook(ctx context.Context, req *api.Book, params api.AddUserBookParams) (api.AddUserBookRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[params.UserID]; !exists {
		s.users[params.UserID] = make(map[int]api.Book)
	}

	if _, exists := s.users[params.UserID][req.ID]; exists {
		return err(http.StatusConflict, "user %d is already reading the book with id %d", params.UserID, req.ID), nil
	}

	s.users[params.UserID][req.ID] = *req
	return req, nil
}

func (s *serviceImpl) GetUserBook(ctx context.Context, params api.GetUserBookParams) (api.GetUserBookRes, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if books, ok := s.users[params.UserID]; !ok {
		return err(http.StatusNotFound, "user %d not found", params.UserID), nil
	} else if book, ok := books[params.BookID]; !ok {
		return err(http.StatusNotFound, "book %d not found for user %d", params.BookID, params.UserID), nil
	} else {
		return &book, nil
	}
}

func (s *serviceImpl) UpdateReadingProgress(ctx context.Context, req *api.UpdateReadingProgressReq, params api.UpdateReadingProgressParams) (api.UpdateReadingProgressRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if books, ok := s.users[params.UserID]; !ok {
		return err(http.StatusNotFound, "user %d not found", params.UserID), nil
	} else if book, ok := books[params.BookID]; !ok {
		return err(http.StatusNotFound, "book %d not found for user %d", params.BookID, params.UserID), nil
	} else {
		book.Page = req.Page
		books[params.BookID] = book
		return &book, nil
	}
}

func (s *serviceImpl) RemoveUserBook(ctx context.Context, params api.RemoveUserBookParams) (api.RemoveUserBookRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if books, ok := s.users[params.UserID]; !ok {
		return err(http.StatusNotFound, "user %d not found", params.UserID), nil
	} else if _, ok := books[params.BookID]; !ok {
		return err(http.StatusNotFound, "book %d not found for user %d", params.BookID, params.UserID), nil
	} else {
		delete(books, params.BookID)
		return &api.RemoveUserBookNoContent{}, nil
	}
}

// func slow(wait time.Duration) func(handler http.Handler) http.Handler {
// 	return func(handler http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			time.Sleep(wait)
// 			handler.ServeHTTP(w, r)
// 		})
// 	}
// }

func main() {
	var service api.Handler = newServiceImpl()

	controller, err := api.NewServer(service)
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":8080", controller); err != nil {
		log.Fatal(err)
	}
}
