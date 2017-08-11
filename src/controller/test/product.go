package test

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

func (p *ProductResource) InitRoute(Container *restful.Container) {

	ws2 := new(restful.WebService)
	ws2.Path("/v2")
	ws2.Consumes(restful.MIME_XML)
	ws2.Produces(restful.MIME_XML)
	ws2.Route(ws2.GET("/").To(p.handleHello))
	ws2.Route(ws2.GET("/{id}").To(p.getOne).
		Doc("get the product by its id").
		Param(ws2.PathParameter("id", "identifier of the product").DataType("string")))

	ws2.Route(ws2.POST("").To(p.postOne).
		Doc("update or create a product").
		Param(ws2.BodyParameter("Product", "a Product (XML)").DataType("main.Product")))
	Container.Add(ws2)
}

type Product struct {
	Id, Title string
}

type ProductResource struct {
	// typically reference a DAO (data-access-object)
}

func (p *ProductResource) getOne(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	log.Println("getting product with id:" + id)
	resp.WriteEntity(Product{Id: id, Title: "test"})
}

func (p *ProductResource) postOne(req *restful.Request, resp *restful.Response) {
	updatedProduct := new(Product)
	err := req.ReadEntity(updatedProduct)
	if err != nil { // bad request
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	log.Println("updating product with id:" + updatedProduct.Id)
}

//测试 test
func (p *ProductResource) handleHello(req *restful.Request, resp *restful.Response) {
	fmt.Println("1111")
	io.WriteString(resp, "hello")
}
