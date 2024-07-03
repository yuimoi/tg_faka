package _type

type SubjectStatusType int

type Proxy struct {
	EnableProxy bool   `json:"enable_proxy" desc:"是否开启网络代理"`
	Protocol    string `json:"protocol" desc:"协议"`
	Host        string `json:"host" desc:"域名"`
	Port        uint   `json:"port" desc:"端口"`
}
type Headers map[string]string

type OrderStatusType int

const (
	OrderStatusPending      OrderStatusType = 0
	OrderStatusSuccess      OrderStatusType = 1
	OrderStatusForceSuccess OrderStatusType = 2
	OrderStatusTimeout      OrderStatusType = -1
	OrderStatusForceClose   OrderStatusType = -2
)

type ProductStatusType int

const (
	ProductStatusValid   ProductStatusType = 1
	ProductStatusInvalid ProductStatusType = 0
)

type ProductItemStatusType int

const (
	ProductItemStatusValid   ProductItemStatusType = 1
	ProductItemStatusPending ProductItemStatusType = 0
	ProductItemStatusLocked  ProductItemStatusType = -1
)

type PaginationQueryDataStruct struct {
	Limit     int
	Page      int
	TotalPage int
}
