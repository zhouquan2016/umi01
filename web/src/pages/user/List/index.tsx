import ProTable, {ProColumns} from "@ant-design/pro-table";
import {useRef, useState} from "react";
import {ActionType} from "@ant-design/pro-table/lib/typing";
import {API} from "@/services/ant-design-pro/typings";
import {userList} from "@/services/ant-design-pro/api";
import {Button, Popconfirm} from "antd";
import {PlusOutlined} from "@ant-design/icons";
import React from "react";
import Save from "./save"
import {user} from "@/services/swagger/user";
import {role} from "@/services/swagger/role";

export default () => {
  const activeRef = useRef<ActionType>()
  const [showSave, setShowSave] = useState<boolean>(false)
  const [editId, setEditId] = useState<null|number>(null)
  const columns: ProColumns<API.CurrentUser>[] = [
    {
      dataIndex: 'index',
      valueType: 'indexBorder',
      width: 48,
    },
    {
      dataIndex: "avatar",
      title: "头像",
      render(_, record) {
        return <img src={record.avatar} width={"24px"} height={"24px"}/>
      },
      hideInSearch: true
    },
    {
      dataIndex: "name",
      title: "姓名"
    },
    {
      dataIndex: "email",
      title: "邮箱"
    },
    {
      dataIndex: "phone",
      title: "手机号"
    },

    {
      dataIndex:"signature",
      title:"个性签名",
      hideInSearch: true
    },

    {
      dataIndex:"title",
      title:"职位",
      hideInSearch: true
    },
    {
      dataIndex:"roleName",
      title:"角色名称",
      hideInSearch: true
    },
    {
      dataIndex:"roleId",
      title:"角色",
      valueType:"select",
      request() {
        // @ts-ignore
        return role.getAll().then(({data}) =>{
          // @ts-ignore
          return data?.map(d => {
            return {
              label: d.name,
              value: d.id
            }
          })
        }).catch((err) => console.log(err))
      }
    },
    {
      valueType: "option",
      title: "操作",
      render: (text, record, _, action) => {
        const ops = []
        if (!record.isSysDefault) {
          ops.push(<a onClick={() => {
            setShowSave(true)
            // @ts-ignore
            setEditId(record.id)
          }}>修改</a>)
          ops.push(<Popconfirm title="确认删除?" onConfirm={() => {
            //@ts-ignore
            user.deleteById(record.id).then(() => {
              action?.reload()
            }).catch(err => console.log(err))
          }}><a >删除</a></Popconfirm>)
        }
        return ops
      }
    }
  ]
  return (
    <React.Fragment>
      <ProTable headerTitle={"用户维护"} actionRef={activeRef} columns={columns}
        rowKey="id"
        request={async (params = {}, sort) => {
          const queryParams = {
            ...params,
            sortField: sort.field,
            sortOrder: sort.order,
          }
          // @ts-ignore
          return userList(queryParams, [])
        }}
        toolBarRender={() => [
          <Button key="button" icon={<PlusOutlined/>} type="primary" onClick={() => {
            setEditId(null)
            setShowSave(true)
          }}>
            新建
          </Button>,
        ]}
      />
      <Save onCancel={() => setShowSave(false)} id={editId} visible={showSave} onOk={() => {
        setShowSave(false)
        activeRef.current?.reload()
      }}/>
    </React.Fragment>
  )
}
