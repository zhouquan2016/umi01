import {ModalForm, ProFormInstance, ProFormText} from "@ant-design/pro-form";
import React, {useRef, useState} from "react";
import {Tag} from "antd";
import {PlusOutlined} from "@ant-design/icons";

interface Resource {
  id: number,
  code: string,
  name: string,
  path: string
}

interface Props {
  value?: Resource[],
  onChange?: (newValue: any, oldValue: any) => void
  menuPath?: string
}

export default ({value, menuPath, onChange}: Props) => {
  const addForm = useRef<ProFormInstance>()
  const [showAdd, setShowAdd] = useState(false)
  const [defaultResource, setDefaultResource] = useState<Resource|null>(null)
  const [defaultIndex, setDefaultIndex] = useState<number|null>(null)

  return <React.Fragment>
    <div>
      {
        value?.map((r, index) => <Tag key={index} closable={true} onClose={(e) => {
          e.preventDefault()
          const newRs = value.slice()
          newRs.splice(index, 1)
          if (onChange) {
            onChange(newRs, value)
          }
        }}>
          <span onClick={() => {
            setDefaultIndex(index)
            setDefaultResource(r)
            setShowAdd(true)
          }}>{r.name}</span>
        </Tag>)
      }
      <Tag className="site-tag-plus" style={{cursor: "pointer"}} onClick={() => {
        setShowAdd(true)
        setDefaultIndex(null)
        setDefaultResource(null)
      }}>
        <PlusOutlined/>
      </Tag>
    </div>
    {
      showAdd &&
      <ModalForm visible={true}
        // @ts-ignore
                 initialValues={defaultResource}
                 onVisibleChange={visible => !visible && setShowAdd(false)}
                 formRef={addForm}
                 onValuesChange={(newValues) => {
                   addForm.current?.setFieldsValue({
                     path: menuPath + "/" + (newValues?.code ? newValues?.code : "")
                   })
                 }}
        // @ts-ignore
                 onFinish={(formData) => {
                   return addForm.current?.validateFields().then(() => {
                     let newRs
                     if (defaultIndex == null) {
                       newRs = value?.slice() || []
                       newRs.push(formData)
                     }else {
                       newRs = value?.slice() || []
                       newRs[defaultIndex] = formData
                     }
                     if (onChange) {
                       onChange(newRs, value)
                     }
                     setShowAdd(false)
                   })
                 }}
      >
        <ProFormText name="name" label="名称" required rules={[
          {
            required: true,
            message: "名称不能为空"
          }
        ]}></ProFormText>
        <ProFormText name="code" label="编码" required rules={[
          {
            required: true,
            message: "编码不能为空"
          },
          {
            validator: ((rule, codeValue) => new Promise<void>((resolve, reject) => {
              if (defaultIndex != null) {
                resolve()
                return
              }
              value?.forEach(r => {
                if (r.code === codeValue) {
                  reject("编码已存在")
                }
              })
              resolve()
            }))
          }
        ]}></ProFormText>
        <ProFormText name="path" disabled label="路径" required></ProFormText>
      </ModalForm>
    }
  </React.Fragment>
}
