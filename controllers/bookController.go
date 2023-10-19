package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gimtwi/go-library-project/types"
	"github.com/gin-gonic/gin"
)

func GetAllBooks(repo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		books, err := repo.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, books)
	}
}

func GetBookByID(repo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
			return
		}

		book, err := repo.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
			return
		}
		c.JSON(http.StatusOK, book)
	}
}

func CreateBook(bookRepo types.BookRepository, authorRepo types.AuthorRepository, genreRepo types.GenreRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		var book types.Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, authorID := range book.AuthorID {
			err := authorRepo.CheckAuthor(authorID)
			fmt.Println("cb", err)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "author not found"})
				return
			}
		}

		for _, genreID := range book.GenreID {
			err := genreRepo.CheckGenre(genreID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "genre not found"})
				return
			}
		}

		//concurrency

		// authorCheck := make(chan error)
		// genreCheck := make(chan error)

		// for _, authorID := range book.AuthorID {
		// 	go func(authorID uint) {
		// 		authorCheck <- HandleAuthorCheck(authorRepo, authorID)
		// 	}(authorID)
		// }

		// for _, genreID := range book.GenreID {
		// 	go func(genreID uint) {
		// 		genreCheck <- HandleGenreCheck(genreRepo, genreID)
		// 	}(genreID)
		// }

		// var authorWg, genreWg sync.WaitGroup
		// authorWg.Add(len(book.AuthorID))
		// genreWg.Add(len(book.GenreID))

		// go func() {
		// 	authorWg.Wait()
		// 	close(authorCheck)
		// }()

		// go func() {
		// 	genreWg.Wait()
		// 	close(genreCheck)
		// }()

		// go func() {
		// 	for err := range authorCheck {
		// 		if err != nil {
		// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 			return
		// 		}
		// 		authorWg.Done()
		// 	}
		// }()

		// go func() {
		// 	for err := range genreCheck {
		// 		if err != nil {
		// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 			return
		// 		}
		// 		genreWg.Done()
		// 	}
		// }()

		if book.Quantity >= 1 {
			book.IsAvailable = true
		} else {
			book.IsAvailable = false
		}

		if err := bookRepo.Create(&book); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("time taken for a book creation: ", time.Since(start))
		c.JSON(http.StatusCreated, book)
	}
}

func UpdateBook(repo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
			return
		}

		var book types.Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		book.ID = uint(id)

		if err := repo.Update(&book); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, book)
	}
}

func DeleteBook(repo types.BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
			return
		}

		if err := repo.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
