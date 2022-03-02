import React, {useRef, useState} from 'react';
import type {ActionType, ProColumns} from '@ant-design/pro-table';
import ProTable from '@ant-design/pro-table';
import {MenuData, menuList} from "@/services/swagger/user";
import {Button, message, Popconfirm, Tag, Tooltip} from 'antd';
import {PlusOutlined} from '@ant-design/icons';
import TreeList, {TreeRef} from "@/pages/Menu/TreeList";
import {ProFormInstance} from "@ant-design/pro-form";
import Add from "@/pages/Menu/Add";
import {menuDelete, menuEdit} from "@/services/swagger/menu";
import Edit from "./edit";


export default () => {
  const treeRef = useRef<TreeRef>()
  const actionRef = useRef<ActionType>();
  const formRef = useRef<ProFormInstance>();
  const [showAdd, setShowAdd] = useState(false)
  const [resourceMenuId, setResourceMenuId] = useState(null)

  const columns: ProColumns<MenuData>[] = [

    {
      dataIndex: 'index',
      valueType: 'indexBorder',
      width: 48,
    },
    {
      title: '名称',
      dataIndex: 'name',
      copyable: true,
      ellipsis: true,
      hideInSearch: true

    },
    {
      title: '菜单层级',
      dataIndex: 'deptName',
      copyable: true,
      ellipsis: true,
      tip: '菜单层级从第一层级到当前层级以/分割',
      hideInSearch: true,
      editable: false,

    },
    {
      title: "父菜单",
      hideInTable: true,
      // @ts-ignore
      renderFormItem: (_, {type, defaultRender, ...rest}, form) => {
        if (type === 'form') {
          return null;
        }
        const status = form.getFieldValue('state');
        if (status !== 'open') {

          return (
            // value 和 onchange 会通过 form 自动注入。
            <TreeList
              // 组件的配置
              {...rest}
              // @ts-ignore
              treeRef={treeRef}
            />
          );
        }
        return defaultRender(_);
      },
      dataIndex: "parentId",
    },
    {
      title: '路径',
      dataIndex: 'path',
    },
    {
      title: '是否叶子节点',
      editable: false,
      dataIndex: 'isLeaf',
      valueType: "select",
      render(_, entity) {
        return entity.isLeaf ? "是" : "否"
      },
      valueEnum: {
        true: {
          text: "是"
        },
        false: {
          text: "否"
        }
      }

    },
    {
      title: '系统预留',
      dataIndex: 'isSysDefault',
      valueType:"select",
      valueEnum: {
        true: {
          text: "是"
        },
        false: {
          text: "否"
        }
      },
      editable: false,
      renderText(value) {
        return value ? "是" : "否"
      }

    },
    {
      title:"资源",
      dataIndex:"resources",
      hideInSearch: true,
      render(text, record, _, action) {
        return record?.resources?.map((resource, index) =>
          <Tooltip key={index} title={resource.code + " " + resource.path}>
            <Tag >{resource.name}</Tag>
          </Tooltip>
          )
      }
    },
    {
      title: '操作',
      valueType: 'option',
      render: (text, record, _, action) => {
        const ops = []
        if (!record.isSysDefault) {
          ops.push(<Popconfirm key="delete" onConfirm={() => {
            menuDelete([record.id]).then(() => {
              action?.reload()
            }).catch((err) => {
              console.log("删除菜单err:", err)
            })
          }}  title={"确认删除?"}><a>删除</a></Popconfirm>)
          ops.push(<a key="addResource" onClick={() => {
            // @ts-ignore
            setResourceMenuId(record.id)
          }}>编辑</a>)
        }
        return ops
      }
    },
  ];

  // @ts-ignore
  return (
    <React.Fragment>
      <ProTable<MenuData>
        columns={columns}
        actionRef={actionRef}
        formRef={formRef}
        request={async (params = {}, sort) => {
          const queryParams = {
            ...params,
            sortField: sort.field,
            sortOrder: sort.order,
          }
          // @ts-ignore
          return menuList(queryParams, [])
        }}
        editable={{
          type: 'multiple',
          actionRender:(row, config, defaultDom) => {
            return [defaultDom.save, defaultDom.cancel]
          },
          // @ts-ignore
          onSave(_, record) {
            let errMsg = ""
            if (record.name == "") {
              errMsg = "名称不能为空"
            }else if (record.path == "") {
              errMsg = "路径不能为空"
            }
            if (errMsg != "") {
              message.warn(errMsg)
              return Promise.reject(errMsg)
            }
            return menuEdit({

              id: record.id,
              name: record.name?.trim(),
              path: record.path?.trim()
            }).then(() => {
              actionRef.current?.reload()
            })
            // @ts-ignore

          }
        }}
        columnsState={{
          persistenceKey: 'pro-table-singe-demos',
          persistenceType: 'localStorage',
        }}
        rowKey="id"
        search={{
          labelWidth: 'auto',
        }}
        pagination={{
          defaultPageSize: 10,
        }}
        dateFormatter="string"
        headerTitle="菜单列表"
        toolBarRender={() => [
          <Button key="button" icon={<PlusOutlined/>} type="primary" onClick={() => {setShowAdd(true)}}>
            新建
          </Button>,
        ]}
      />
      {
        showAdd && <Add visible={showAdd} onCancel={() => {setShowAdd(false)}} okCallBack={() => {
          formRef.current?.resetFields();
          treeRef.current?.reload()
          actionRef?.current?.reload()
          setShowAdd(false)
        }}/>
      }
      {
        resourceMenuId && <Edit menuId={resourceMenuId} clearMenuId={() => {setResourceMenuId(null)}}
                                onSuccess={() => {
                                  setResourceMenuId(null)
                                  actionRef.current?.reload()
                                }}
                          />
      }
    </React.Fragment>
  );
}
