import ProTable, {ActionType, ProColumns} from '@ant-design/pro-table';
import {Role, roleAdd, roleDelete, roleEdit, roleGetById, roleList} from "@/services/swagger/role";
import {Button, message, Popconfirm} from "antd";
import {PlusOutlined} from "@ant-design/icons";
import React, {useRef, useState} from "react";
import Save, {EditType, roleNameRules} from "./add";
import {ProFormInstance} from "@ant-design/pro-form";

export default () => {
  const [editType, setEditType] = useState<EditType|null>(null)
  const [editValues, setEditValues] = useState<any|null>(null)
  const columns: ProColumns<Role>[] = [
    {
      dataIndex: 'index',
      valueType: 'indexBorder',
      width: 48,
    },
    {
      dataIndex: "no",
      title: "角色编号",
      editable: false
    },
    {
      dataIndex: "name",
      title: "角色名称",
      formItemProps: {
        // @ts-ignore
        rules: roleNameRules,
      },
      hideInSearch: true,
    },
    {
      dataIndex: "isSysDefault",
      title: "系统预留",
      editable: false,
      renderText(value) {
        return value ? "是" : "否"
      },
      valueType: "select",
      valueEnum: {
        true: "是",
        false: "否",
      }
    },
    {
      title: '操作',
      valueType: 'option',
      render: (text, record, _, action) => {
        const ops = []
        if (!record.isSysDefault) {
          ops.push(
            <Popconfirm title={"确认删除?"} onConfirm={() => {
              roleDelete(record.id).then(() => action?.reload()).catch(err => console.error(err))
            }}>
            <a>删除</a>
            </Popconfirm>
              )
          ops.push(<a onClick={() => {
            // @ts-ignore
            roleGetById(record.id).then(({data}) => {
              if (!data) {
                message.warn("角色未找到")
                return
              }
              setEditType("edit")
              setEditValues(data)
            }).catch(err => console.log(err))
          }
          }>修改</a>)
        }
        return ops
      }
    }
  ]

  const formRef = useRef<ProFormInstance>()
  const actionRef = useRef<ActionType>()
  return (<
    React.Fragment>
    < ProTable
      columns={columns}
      formRef={formRef}
      actionRef={actionRef}
      request={async (params = {}, sort) => {
        const queryParams = {
          ...params,
          sortField: sort.field,
          sortOrder: sort.order,
        }
        // @ts-ignore
        return roleList(queryParams)
      }
      }
      rowKey="id"
      toolBarRender={() =>
        [
          <Button key="button" icon={<PlusOutlined/>} type="primary" onClick={() => {
            setEditType("save")
          }}>
            新建
          </Button>,
        ]
      }

    />
    {
      <Save editType={editType} onCancel={() => setEditType(null)} initValues={editValues}
           onSuccess={(values) => {
             let op;
             if (editType === "save") {
               op = roleAdd(values)
             }else if (editType === "edit"){
               op = roleEdit(values)
             }else {
               message.warn("未知编辑类型")
               return
             }
             op.then(() => {
               setEditType(null);
               actionRef.current?.reload()
             }).catch((err) => {
               console.error(err)
             })

           }}/>
    }
  </React.Fragment>)
}
