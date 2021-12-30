package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Article struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Intro    string `json:"intro"`
	Content  string `json:"content"`
	AuthorId int    `json:"author_id"`
}
type ArticleCRUD struct {
	articles []Article
}

func (A *ArticleCRUD) List() []Article {
	return A.articles
}
func (A *ArticleCRUD) Get(id int) (Article, error) {
	for _, article := range A.articles {
		if article.Id == id {
			return article, nil
		}
	}
	return Article{}, errors.New("Resource Not found")
}
func (A *ArticleCRUD) Post(i Article) {
	defer A.Save()
}
func (A *ArticleCRUD) Update(i Article) {
	defer A.Save()
}

func (A *ArticleCRUD) Init() {
	data, err := ioutil.ReadFile("./articles.json")
	if err != nil {
		fmt.Print(err)
	}
	var storedarticles []Article
	err = json.Unmarshal(data, &storedarticles)
	if err != nil {
		fmt.Println("error:", err)
	}
	A.articles = storedarticles

}

func (A *ArticleCRUD) Save() {
	file, _ := json.MarshalIndent(A.articles, "", " ")
	//fmt.Println(file)

	_ = ioutil.WriteFile("articles.json", file, 0644)
}
func defertest() int {
	defer fmt.Println(2)
	return 1
}
func main() {
	defertest()
	articlecrud := ArticleCRUD{}
	articlecrud.Init()
	defer articlecrud.Save()
	r := gin.Default()
	//
	r.GET("/api/v1/articles", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": articlecrud.List()})
	})
	r.GET("/api/v1/articles/:id", func(c *gin.Context) {
		aid := c.Param("id")
		if aid != "" {
			article_id, err := strconv.Atoi(aid)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"reason": "invalid id"})

			} else {
				article, err := articlecrud.Get(article_id)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"reason": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"data": article})
				return
			}

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"reason": "invalid id"})
		}

	})
	r.POST("/api/v1/articles", func(c *gin.Context) {
		var articleinput Article
		c.BindJSON(&articleinput)
		duplicate := false
		for _, article := range articlecrud.articles {
			if article.Id == articleinput.Id {
				duplicate = true
				break
			}
		}
		if !duplicate {
			articlecrud.articles = append(articlecrud.articles, articleinput)
			articlecrud.Save()
			c.JSON(http.StatusCreated, gin.H{"data": articleinput})
			return
		} else {
			c.JSON(http.StatusNotFound, gin.H{"reason": "duplicate article"})
		}
		//
		//} else {
		//	fmt.Println(err.Error())
		//}

	})
	r.PUT("/api/v1/articles/:id/update", func(c *gin.Context) {
		var articleinput Article
		c.BindJSON(&articleinput)
		aid := c.Param("id")
		if aid != "" {
			article_id, err := strconv.Atoi(aid)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"reason": "invalid id"})

			} else {
				if article_id != articleinput.Id {
					c.JSON(http.StatusBadRequest, gin.H{"reason": "unintended update"})
					return
				}
				for i, article := range articlecrud.articles {
					if article.Id == articleinput.Id {
						articlecrud.articles[i] = articleinput
						articlecrud.Save()
						c.JSON(http.StatusPartialContent, gin.H{"data": articleinput})

						return

					}
				}
			}

		} else {
			c.JSON(http.StatusNotFound, gin.H{"reason": "not found any matching article"})
		}

	})
	//
	r.Run()
	r.Run(":5002")
}
