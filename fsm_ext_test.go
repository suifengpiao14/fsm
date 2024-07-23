package fsm_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/fsm"
)

const (
	Order_status_init      = "0" // 初始状态
	Order_status_created   = "1" // 已下单,待买家付款
	Order_status_paid      = "2" //已付款,待卖家发货
	Order_status_shipped   = "3" //已发货,待买家收货
	Order_status_received  = "4" //买家已收货,订单完成
	Order_status_refunded  = "5" //已退款(发生售后，并且已经退款)
	Order_status_cancelled = "6" //已取消(中途终止)
)

var orderEvents = fsm.Events{
	// 下单
	{
		Name: "order.created",
		Src: []string{
			Order_status_created, // 重入
			Order_status_init,    // 正常
		},
		Dst: Order_status_created,
	},

	// 已付款
	{
		Name: "order.paid",
		Src: []string{
			Order_status_init,    //支持乱序跳跃
			Order_status_created, //下单后付款
			Order_status_paid,    //重入
		},
		Dst: Order_status_paid,
	},
	// 已发货
	{
		Name: "order.shipped",
		Src: []string{
			Order_status_init,    //支持乱序跳跃
			Order_status_created, //支持乱序跳跃
			Order_status_paid,    //付款后发货
			Order_status_shipped, // 重入
		},
		Dst: Order_status_shipped,
	},
	// 已收货
	{
		Name: "order.shipped",
		Src: []string{
			Order_status_init,     //支持乱序跳跃
			Order_status_created,  //支持乱序跳跃
			Order_status_paid,     //支持乱序跳跃
			Order_status_shipped,  // 商家发货后，用户收货
			Order_status_received, //重入
		},
		Dst: Order_status_received,
	},
	// 已退款
	{
		Name: "order.shipped",
		Src: []string{
			Order_status_init,     //支持乱序跳跃
			Order_status_created,  //支持乱序跳跃
			Order_status_paid,     //买家支付后未发货，发生售后并完成，到已退款
			Order_status_shipped,  // 商家发货后，发生售后并完成，到已退款
			Order_status_received, //收货后，发生售后并完成 到已退款
			Order_status_refunded, //重入
		},
		Dst: Order_status_refunded,
	},
	// 已取消
	{
		Name: "order.shipped",
		Src: []string{
			Order_status_init,    //支持乱序跳跃
			Order_status_created, //创建单，未付款前取消
			//Order_status_paid,     //买家支付后未发货，取消必须走售后流程，不会到达Order_status_cancelled 终态
			//Order_status_shipped,  // 商家发货后,取消必须走售后流程，不会到达Order_status_cancelled 终态
			//Order_status_received, //收货后,不售后情况为终态，售后后只能流转到 Order_status_refunded
			//Order_status_refunded, // 已经退款是终态，不可再流转
			Order_status_cancelled, //重入
		},
		Dst: Order_status_cancelled,
	},
}

func TestIsReverseOrder(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		oldStatus := Order_status_refunded
		newStatus := Order_status_received
		var orderFsm = fsm.NewFSM(
			oldStatus,
			orderEvents,
			fsm.Callbacks{},
		)
		ok := orderFsm.IsReverseOrder(newStatus)
		require.Equal(t, true, ok)
	})

	t.Run("false", func(t *testing.T) {
		oldStatus := Order_status_received
		newStatus := Order_status_cancelled
		var orderFsm = fsm.NewFSM(
			oldStatus,
			orderEvents,
			fsm.Callbacks{},
		)
		ok := orderFsm.IsReverseOrder(newStatus)
		require.Equal(t, false, ok)
	})

}
