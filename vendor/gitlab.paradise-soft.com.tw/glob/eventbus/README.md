# eventbus

[![pipeline status](http://gitlab.paradise-soft.com.tw/glob/eventbus/badges/develop/pipeline.svg)](http://gitlab.paradise-soft.com.tw/glob/eventbus/commits/develop)
[![coverage report](http://gitlab.paradise-soft.com.tw/glob/eventbus/badges/develop/coverage.svg)](http://gitlab.paradise-soft.com.tw/glob/eventbus/commits/develop)

## How to use

Example for `xunya-apis`, `xunya-service`

```go
func main() {
	// eventbus.New 建立消息隊列工具物件
	// eventbus.WithInternalEngine 注入服務內部消息推送模組
	ev := eventbus.New(eventbus.WithInternalEngine())

	// 訂閱後回傳entryID，與該筆訂閱的callback 綁定
	entryID, err := ev.Subscribe("order", callback)
	if err != nil {
		// do something
	}

	// 可利用entryID 來退訂
	ev.Unsubscribe(entryID)

	// 發佈消息
	ev.Publish("order", "message")
}

func callback(data []byte) {
	var str string
	if err := eventbus.Unmarshal(data, &str); err != nil {
		// do something
	}

	// do something
}
```

## Architect

### IEvent

使用者可使用此功能進行訂閱，退訂，發布。
依賴於`IEngine`

* Subscribe(訂閱)
    * 需傳入callback 方法，當收到訂閱消息時執行callback
    * 訂閱後回傳EntryID，為此訂閱的唯一識別
    * 內部實作
        * 註冊callback 到同個topic
        * 啟動engine 訂閱topic
        * 取得唯一識別
    
* Unsubscribe(退訂)
    * 傳入EntryID 進行退訂
    * 內部實作
        * 刪除註冊的callback
        * 當沒有任何callback 時對engine 退訂topic
        
* Publish(發佈)
    * 推送消息
    * 內部實作
        * 呼叫engine 推送消息

### IEngine

實作對 Message queue server ，包含連線，訂閱，發布。
目前有
* RedisEngine
    * 使用redis 當作message queue
    * 相依套件`github.com/garyburd/redigo/redis`
    * 創建時需傳入`redis.pool`
* InternalEngine
    * 在engine 內部實現message queue
    * 可當作服務內部的消息推送使用

## TODO 

* [ ] unit test 
