package handler

import (
	"context"
	"github.com/XiaoBinGan/common"
	"github.com/XiaoBinGan/cart/domain/model"
	"github.com/XiaoBinGan/cart/domain/service"
	. "github.com/XiaoBinGan/cart/proto/cart"
)

type Cart struct {
	CartDataService service.ICartDataService
}
//add cart


//clean cart


//add cart number


//reduce cart number

//delete cart

//get all cart info by user id

/**
   Add Cart
   init mode struct
   json.marshal request ,  unmarshal for the mod struct
   use CartDataService.AddCart(mode)
   return the cart id and errors
 */
func(c *Cart)AddCart(ctx context.Context,req *CartInfo,res *ResponseAdd) (err error){
	cart :=&model.Cart{}
	common.SwapTo(req, cart)
	res.CartId, err = c.CartDataService.AddCart(cart)
	return err

}
/**
   CleanCart
   use c.CartDataService.CleanCart(req.UserId)
 */
func(c *Cart)CleanCart(ctx context.Context,req *Clean,res *Response) error{
	if err := c.CartDataService.CleanCart(req.UserId);err!=nil{
		return err
	}
	res.Msg="购物车清除成功"
	return nil
}
/**
    add product number on cart

 */
func(c *Cart)Incr(ctx context.Context,req *Item,res *Response) error{
     if err:=c.CartDataService.IncrNum(req.Id,req.ChangeNum);err!=nil{
     	return err
	 }
	res.Msg="商品数量添加成功"
	return nil
}
func(c *Cart)Decr(ctx context.Context,req *Item,res *Response) error{
	if err:=c.CartDataService.DecrNum(req.Id,req.ChangeNum);err!=nil{
		return err
	}
	res.Msg="商品添加成功"
	return nil
}
func(c *Cart)DeleteItemByID(ctx context.Context,req *CartID,res *Response) error{
	if err := c.CartDataService.DeleteCart(req.Id);err!=nil{
		return err
	}
	res.Msg="删除成功"
	return nil
}
func(c *Cart)GetAll(ctx context.Context,req *CartFindAll,res *CartAll) error{
	 carts, err := c.CartDataService.FindAllCart(req.UserId)
	 if err!=nil{
	 	return err
	 }
	for _, v := range carts {
		c :=&CartInfo{}
		if err := common.SwapTo(v, c);err!=nil{
			return nil
		}
        res.CartInfo=append(res.CartInfo,c)
	}
	return nil
}