package user

import (
	"errors"
	"fmt"
	"github.com/JoyZF/go-micro-kit/mock"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUser_GetUserInfo(t *testing.T) {
	// gomock.NewController：返回 gomock.Controller，它代表 mock 生态系统中的顶级控件。
	// 定义了 mock 对象的范围、生命周期和期待值。
	// 另外它在多个 goroutine 中是安全的
	ctl := gomock.NewController(t)
	// ctl.Finish()：进行 mock 用例的期望值断言，一般会使用 defer 延迟执行，以防止我们忘记这一操作
	defer ctl.Finish()

	var id int64 = 1
	var id2 int64 = 2
	// mock.NewMockMale：创建一个新的 mock 实例
	mockMale := mock.NewMockMale(ctl)
	gomock.InOrder(
		// mockMale.EXPECT().Get(id).Return(nil)：这里有三个步骤，EXPECT()返回一个允许调用者设置期望和返回值的对象。
		// Get(id) 是设置入参并调用 mock 实例中的方法。
		// Return(nil) 是设置先前调用的方法出参。
		// 简单来说，就是设置入参并调用，最后设置返回值
		mockMale.EXPECT().Get(id).Return(errors.New("user not found")),
		mockMale.EXPECT().Get(id2).Return(nil),
	)

	// NewUser(mockMale)：创建 User 实例，值得注意的是，在这里注入了 mock 对象，因此实际在随后的 user.GetUserInfo(id) 调用（入参：id 为 1）中。
	// 它调用的是我们事先模拟好的 mock 方法
	user := NewUser(mockMale)
	fmt.Println(user.GetUserInfo(id))
	fmt.Println(user.GetUserInfo(id2))
}
