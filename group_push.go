package xinge

import (
    "encoding/json"
    "errors"
    "reflect"
    "strconv"
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

    return &api
}

func (api *GroupPushAPI) CreateMessage(msgtype MessageType, msgbody interface{}, env PushEnv) error {
    raw, err := json.Marshal(msgbody)
    if err != nil {
        return err
    }

    req := api.client.NewRequest("GET", createMultiMessageUrl)
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

    v1, ok := resp.Result.(map[string]interface{})
    if !ok {
        return errors.New("xinge: parse response result failure")
    }

    v2 := reflect.ValueOf(v1["push_id"])
    v3, err := strconv.ParseUint(v2.String(), 10, 32)
    if err != nil {
        return err
    }

    api.messageId = uint32(v3)

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
