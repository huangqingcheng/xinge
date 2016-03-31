package xinge

import (
    "encoding/json"
    "errors"
    "fmt"
    "strings"
    "time"
)

type GroupPushAPI struct {
    client    *Client // 推送客户端
    messageId uint32  // 消息编号
}

func NewGroupPushAPI(client *Client) *GroupPushAPI {
    api := GroupPushAPI{
        client:    client,
        messageId: 0,
    }
}

func (api *GroupPushAPI) CreateMessage(msgtype MessageType, msgbody interface{}, env PushEnv) error {
    raw, err := json.Marshal(msgbody)
    if err != nil {
        return err
    }

    req := api.client.NewRequest("POST", createMultiMessageUrl)
    req.SetParam("message_type", msgtype)
    req.SetParam("message", string(raw))
    req.SetParam("expire_time", 1)
    req.SetParam("multi_pkg", MultiPkg_aid)
    req.SetParam("environment", env)

    resp, err := req.Execute()
    if err != nil {
        return err
    }

    if !resp.OK() {
        return errors.New("xinge: response err: " + resp.Error())
    }

    println(fmt.Sprintf("%v", resp.Result))

    return nil
}

func (api *GroupPushAPI) Push(devices []string) error {
    if api.messageId == 0 {
        return errors.New("xinge: message id is not set")
    }

    raw, err := json.Marshal(devices)
    if err != nil {
        return err
    }

    req := api.client.NewRequest("POST", multiDeviceUrl)
    req.SetParam("device_list", string(raw))
    req.SetParam("push_id", api.messageId)

    resp, err := req.Execute()
    if err != nil {
        return err
    }

    if !resp.OK() {
        return errors.New("xinge: response err: " + resp.Error())
    }

    return nil
}
