package webfeature

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Username    string    `json:"username" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	Price       int       `json:"price" binding:"required"`
	Description string    `json:"description" binding:"required"`
	CreatedAt   time.Time `json:"createdAt"`
}

type productHandler struct {
	sync.RWMutex
	products map[string]Product
}

func newProductHandle() *productHandler {
	return &productHandler{products: make(map[string]Product)}
}

func (p *productHandler) Create(c *gin.Context) {
	p.Lock()
	defer p.Unlock()

	// 1.参数解析
	var product Product
	if err := c.ShouldBindHeader(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2.参数校验
	if _, ok := p.products[product.Name]; ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("product %s already exist", product.Name)})
		return
	}
	product.CreatedAt = time.Now()

	// 3.逻辑处理
	p.products[product.Name] = product
	log.Printf("Register product %s success", product.Name)

	// 4.返回结果
	c.JSON(http.StatusOK, product)
}

func (p *productHandler) Get(c *gin.Context) {
	p.Lock()
	defer p.Unlock()

	product, ok := p.products[c.Param("name")]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("can not found product %s", c.Param("name"))})
		return
	}

	c.JSON(http.StatusOK, product)
}

func router() http.Handler {
	router := gin.Default()
	productHandler := newProductHandle()
	v1 := router.Group("/v1")
	productv1 := v1.Group("/products")
	productv1.POST("", productHandler.Create)
	productv1.GET(":name", productHandler.Get)
	return router
}

func main() {
	var eg errgroup.Group

	// 一进程多端口
	insecureServer := &http.Server{Addr: ":8080", Handler: router(), ReadTimeout: 5 * time.Second, WriteTimeout: 10 * time.Second}
	secureServer := &http.Server{

		Addr:         ":8443",
		Handler:      router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	eg.Go(func() error {
		err := insecureServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	eg.Go(func() error {
		err := secureServer.ListenAndServeTLS("server.pem", "server.key")
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}
