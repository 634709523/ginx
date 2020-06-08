package ziface

/*
	路由接口，这里面路由是 使用口夹着给该连接自定的处理业务
	路由的IRequest 则包含使用该连接信息和这连接的请求数据信息
*/

type IRouter interface{
	PreHandle(request IRequest)  //在处理conn业务之前的钩子方法
    Handle(request IRequest)     //处理conn业务的方法
	PostHandle(request IRequest) //处理conn业务之后的钩子方法
}