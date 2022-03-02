import ProTable from "@ant-design/pro-table";
import {ProColumns} from "@ant-design/pro-table/es";
import {history, History} from "@/services/swagger/history";
import TextArea from "antd/es/input/TextArea";
import {Button} from "antd";
import {useState} from "react";

export default () => {
  const columns: ProColumns<History>[] = [
    {
      dataIndex: 'index',
      valueType: 'indexBorder',
      width: 48,
    },
    {
      title:"用户id",
      dataIndex:"userId"
    },
    {
      title:"请求路径",
      dataIndex:"path"
    },
    {
      title:"请求方法",
      dataIndex:"method",
      valueEnum: {
        "GET": "GET",
        "POST": "POST"
      }
    },
    {
      title:"请求表单",
      dataIndex:"formData"
    },
    {
      title:"请求类型",
      dataIndex:"requestType"
    },
    {
      title:"请求体",
      dataIndex:"requestBody"
    },
    {
      title:"返回类型",
      dataIndex:"responseType"
    },
    {
      title:"返回体",
      dataIndex:"responseBody",
      render(dom, value){
        return <TextArea value={value.responseBody} readOnly={true}/>
      }
    },
  ]
  const [indexAllLoading, setIndexAllLoading] = useState<boolean>(false)
  const handlerIndexAll = () => {
    setIndexAllLoading(true)
    history.indexAll().catch(err => console.error(err)).finally(() => {setIndexAllLoading(false)})
  }
  return <ProTable columns={columns}
                   request={async (params = {}, sort) => {
                     const queryParams = {
                       ...params,
                       sortField: sort.field,
                       sortOrder: sort.order,
                     }
                     // @ts-ignore
                     return history.list(queryParams)
                   }
                   }
                   toolBarRender={() => [
                     <Button key="button" type="primary" onClick={() => {handlerIndexAll()}} loading={indexAllLoading}>
                       索引全部
                     </Button>,
                   ]}

  ></ProTable>
}
