package actor

import "testing"

func TestClosedChannelSend(t *testing.T) {
	// 创建一个有缓冲的 channel，并关闭它
	ch := make(chan int, 5)
	close(ch)

	// 尝试在已关闭的 channel 上发送消息
	select {
	case ch <- 42:
		t.Error("Sending on a closed channel did not fail.")
	default:
		// 发送失败，这是预期的行为
	}

	// 尝试关闭已关闭的 channel
	close(ch)

	// 尝试在已关闭的 channel 上发送消息
	select {
	case ch <- 42:
		t.Error("Sending on a doubly closed channel did not fail.")
	default:
		// 发送失败，这是预期的行为
	}

	// 从已关闭的 channel 接收消息
	select {
	case x, ok := <-ch:
		if ok {
			t.Errorf("Received value %v from a closed channel.", x)
		}
	default:
		// 接收失败，这是预期的行为
	}

	// 再次关闭已关闭的 channel
	close(ch)

	// 从已关闭的 channel 接收消息
	select {
	case x, ok := <-ch:
		if ok {
			t.Errorf("Received value %v from a doubly closed channel.", x)
		}
	default:
		// 接收失败，这是预期的行为
	}
}
