package test

import (
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

func (p *ProductResource) InitRoute(Container *restful.Container) {

	ws2 := new(restful.WebService)
	ws2.Path("/products")
	ws2.Consumes(restful.MIME_XML)
	ws2.Produces(restful.MIME_XML)

	ws2.Route(ws2.GET("/{id}").To(p.getOne).
		Doc("get the product by its id").
		Param(ws2.PathParameter("id", "identifier of the product").DataType("string")))

	ws2.Route(ws2.POST("").To(p.postOne).
		Doc("update or create a product").
		Param(ws2.BodyParameter("Product", "a Product (XML)").DataType("main.Product")))

	restful.Add(ws2)
}

type Product struct {
	Id, Title string
}

type ProductResource struct {
	// typically reference a DAO (data-access-object)
}

func (p ProductResource) getOne(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("id")
	log.Println("getting product with id:" + id)
	resp.WriteEntity(Product{Id: id, Title: "test"})
}

func (p ProductResource) postOne(req *restful.Request, resp *restful.Response) {
	updatedProduct := new(Product)
	err := req.ReadEntity(updatedProduct)
	if err != nil { // bad request
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	log.Println("updating product with id:" + updatedProduct.Id)
}
