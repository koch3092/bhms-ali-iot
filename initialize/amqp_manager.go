package initialize

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"pack.ag/amqp"
	"time"
)

type AmqpManager struct {
	Address  string
	Username string
	Password string
	Client   *amqp.Client
	Session  *amqp.Session
	Receiver *amqp.Receiver
	Logger   *zap.Logger
}

func (am *AmqpManager) StartReceiveMessage(ctx context.Context, sdRcvMsg chan<- *amqp.Message, aRcvMsg chan<- *amqp.Message) {
	childCtx, childDone := context.WithCancel(ctx)
	am.Logger.Info("Start connect amqp server: " + am.Address)
	err := am.generateReceiverWithRetry(childCtx)
	if nil != err {
		childDone()
		return
	}
	defer func() {
		fmt.Printf("Defer amqp manager: %#v\n", am)
		_ = am.Receiver.Close(childCtx)
		am.Receiver = nil
		fmt.Printf("Close am.Receiver: %#v\n", am.Receiver)

		_ = am.Session.Close(childCtx)
		am.Session = nil
		fmt.Printf("Close am.Session: %#v\n", am.Session)

		_ = am.Client.Close()
		am.Client = nil
		fmt.Printf("Close am.Client: %#v\n", am.Client)
		fmt.Printf("Defer amqp manager: %#v\n", am)
	}()

	for {
		// 阻塞接受消息，如果ctx是background则不会被打断。
		message, err := am.Receiver.Receive(ctx)
		am.Logger.Info(fmt.Sprintf("data received: %s properties: %#v", string(message.GetData()), message.ApplicationProperties))

		if err == nil {
			// 收到数据后，发送到Channel中，给到另外一个线程处理
			sdRcvMsg <- message
			aRcvMsg <- message
			err := message.Accept()
			if err != nil {
				childDone()
				return
			}
		} else {
			fmt.Println("amqp receive data error:", err)

			//如果是主动取消，则退出程序。
			select {
			case <-childCtx.Done():
				fmt.Println("--- childCtx.Done()")
				childDone()
				return
			default:
			}

			// 非主动取消，则重新建立连接。
			err := am.generateReceiverWithRetry(childCtx)
			if nil != err {
				childDone()
				return
			}
		}
	}
}

func (am *AmqpManager) generateReceiverWithRetry(ctx context.Context) error {
	// 退避重连，从10ms依次x2，直到20s。
	duration := 10 * time.Millisecond
	maxDuration := 20000 * time.Millisecond
	times := 1

	// 异常情况，退避重连。
	for {
		select {
		case <-ctx.Done():
			return amqp.ErrConnClosed
		default:
		}

		err := am.generateReceiver()
		if nil != err {
			time.Sleep(duration)
			if duration < maxDuration {
				duration *= 2
			}
			am.Logger.Info(fmt.Sprintf("Amqp connect retry, times: %d, duration: %s", times, duration))
			times++
		} else {
			am.Logger.Info("Amqp connect init success")
			return nil
		}
	}
}

// 由于包不可见，无法判断Connection和Session状态，重启连接获取。
func (am *AmqpManager) generateReceiver() error {
	if am.Session != nil {
		receiver, err := am.Session.NewReceiver(
			amqp.LinkSourceAddress("/queue-name"),
			amqp.LinkCredit(20),
		)
		// 如果断网等行为发生，Connection会关闭导致Session建立失败，未关闭连接则建立成功。
		if err == nil {
			am.Receiver = receiver
			return nil
		}
	}

	// 清理上一个连接。
	if am.Client != nil {
		err := am.Client.Close()
		if err != nil {
			return err
		}
	}

	client, err := amqp.Dial(am.Address, amqp.ConnSASLPlain(am.Username, am.Password))
	if err != nil {
		return err
	}
	am.Client = client

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	am.Session = session

	receiver, err := am.Session.NewReceiver(
		amqp.LinkSourceAddress("/queue-name"),
		amqp.LinkCredit(20),
	)
	if err != nil {
		return err
	}
	am.Receiver = receiver

	return nil
}
