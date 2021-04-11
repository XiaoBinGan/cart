package service

import (
	"github.com/XiaoBinGan/cart/domain/model"
	"github.com/XiaoBinGan/cart/domain/repository"
)

/**
    I Cart Data Service
    1.add cart
	2.delete cart
    3.update cart
    4.find cart by cart id
  	5.find all Cart
    6.clean cart by cart id
    7.add product to the cart
	8.reduce the product on the  cart
 */
type ICartDataService interface {
     AddCart(*model.Cart)(int64,error)
     DeleteCart(int64)error
     UpdateCart(*model.Cart)error
     FindCartByID(int64)(*model.Cart,error)
     FindAllCart(int64)([]model.Cart,error)

     CleanCart(int64)error
     DecrNum(int64,int64)error
     IncrNum(int64,int64)error
}

func NewCartDataService(cartRepository repository.ICartRepository)ICartDataService  {
	return &CartDataService{CartRepository: cartRepository}
}

type CartDataService struct {
   CartRepository repository.ICartRepository
}
func(c *CartDataService)AddCart(cart *model.Cart)(int64,error){
	return  c.CartRepository.CreateCart(cart)
}
func(c *CartDataService)DeleteCart(cartId int64)error{
	return  c.CartRepository.DeleteCartByID(cartId)
}
func(c *CartDataService)UpdateCart(cart *model.Cart)error{
	return c.CartRepository.UpdateCart(cart)
}
func(c *CartDataService)FindCartByID(userId int64)(*model.Cart,error){
	return c.CartRepository.FindCartByID(userId)
}
func(c *CartDataService)FindAllCart(userId int64)([]model.Cart,error){
	return  c.CartRepository.FindAll(userId)
}
func(c *CartDataService)CleanCart(cartId int64)error{
	return c.CartRepository.CleanCart(cartId)
}
func(c *CartDataService)DecrNum(cartId int64,num int64)error{
	  return c.CartRepository.DecrNum(cartId,num)
}
func(c *CartDataService)IncrNum(cartId int64,num int64)error{
	  return c.CartRepository.DecrNum(cartId,num)
}