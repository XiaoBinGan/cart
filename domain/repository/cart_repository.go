package repository

import (
    "errors"
    "github.com/jinzhu/gorm"
   "github.com/XiaoBinGan/cart/domain/model"
)

/**
   ICartRepository
   1.create table by the Cart model struct
   2.Find Cart by cart id
   3.create Cart by cart model
   4.delete cart by cart id
   5.update cart by cart model
   6.Find all cart by cart model slice
 */
type ICartRepository interface {
   InitTable()error
   FindCartByID(int64)(*model.Cart,error)
   CreateCart(*model.Cart)(int64,error)
   DeleteCartByID(int64)error
   UpdateCart(*model.Cart)error
   FindAll(int64)([]model.Cart,error)

   CleanCart(int64)error
   IncrNum(int64,int64)error
   DecrNum(int64,int64)error
}

/**
  newCartRepository
  parameter is *gorm.DB
  before return you need use CartRepository implement ICartRepository
  return type ICartRepository
 */
func NewCartRepository(db *gorm.DB) ICartRepository {
      return &CartRepository{mysqldb:db}
}

type CartRepository struct {
    mysqldb *gorm.DB
}

//Init table
func(c *CartRepository)InitTable()error{
    return c.mysqldb.CreateTable(&model.Cart{}).Error
}
//find cart info by cart id
func (c *CartRepository)FindCartByID(cartID int64)(cart *model.Cart,err error)  {
    cart =&model.Cart{}
    return cart,c.mysqldb.First(cart, cartID).Error
    //First find first record that match given conditions, order by primary key
}
/**
   create cart
   FirstOrCreate find first matched record or create a new one with given conditions (only works with struct, map conditions)
   if db.err!=nil{}
   if db.RowsAffected ==0  mysql can't changed  it means cart insert failed

*/

func (c *CartRepository)CreateCart(cart *model.Cart)(int64,error) {
    // FirstOrCreate find first matched record or create a new one with given conditions (only works with struct, map conditions)
    db := c.mysqldb.FirstOrCreate(cart, &model.Cart{
        ProductId: cart.ProductId,
        Num:       cart.Num,
        SizeID:    cart.SizeID,
        UserID:    cart.UserID,
    })
    if db.Error!=nil{
        return 0,db.Error
    }
    if db.RowsAffected==0{
        return 0,errors.New("购物车插入失败")
    }
    return cart.ID,nil
}

/**
   Delete Cart By ID(int64)
   Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition WARNING If model has DeletedAt field, GORM will only set field DeletedAt's value to current time
 */
func (c *CartRepository)DeleteCartByID(cartId int64)error {
    return  c.mysqldb.Where(" ",cartId).Delete(&model.Cart{}).Error
}


/**
  update cart info
    Model指定要运行db操作的模型
    //更新所有用户名为hello
    db.Model(&User{}).Update("name", "hello")
    //如果用户的主键非空，将使用它作为条件，然后将只更新用户名为' hello '
    db.Model(&user).Update("name", "hello")
 */
func (c *CartRepository)UpdateCart(cart *model.Cart)error  {
    return c.mysqldb.Model(cart).Update(cart).Error
}
/**
  get all result collection
 */
func (c *CartRepository)FindAll(userID int64)(cartAll []model.Cart,err error){
    return cartAll,c.mysqldb.Where("user_id=?",userID).Find(cartAll).Error
}


/**
   clear cart by user id
 */
func (c *CartRepository)CleanCart(userid int64)error {
    return c.mysqldb.Where("user_id = ?",userid).Delete(&model.Cart{}).Error
}

/**
   add product number
    Expr generate raw SQL expression, for example:
    DB.Model(&product).Update("price", gorm.Expr("price * ? + ?", 2, 100))
 */
func (c *CartRepository)IncrNum(cartId  int64,num int64)error{
    cart :=&model.Cart{ID:cartId}
    return c.mysqldb.Model(cart).UpdateColumn("num", gorm.Expr("num +? ", num)).Error
}

/**
  reduce the cart product by  cart id and set number
 */
func (c *CartRepository)DecrNum(cartID int64,num int64)error  {
   cart := &model.Cart{ID: cartID}
    db := c.mysqldb.Model(cart).Where("num >=?", num).UpdateColumn("num", gorm.Expr("num - ?", num))
    if db.Error!=nil{
        return db.Error
    }
    if db.RowsAffected == 0{
        return errors.New("减少商品失败")
    }
    return nil
}