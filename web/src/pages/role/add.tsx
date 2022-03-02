import {DrawerForm, ProFormInstance, ProFormText} from "@ant-design/pro-form";
import {roleExists} from "@/services/swagger/role";
import ProFormItem from "@ant-design/pro-form/es/components/FormItem";
import Auth from "./auth"
import {useRef} from "react";

export type EditType = "edit"|"save"

interface Props {
  onCancel: () => void
  onSuccess: (values: any) => void
  initValues: any,
  editType: EditType|null
}

export const roleNameRules = [
  {required: true, type: "string", min: 2, max: 20},
]
export default ({onCancel, onSuccess, initValues, editType}: Props) => {
  const formRef = useRef<ProFormInstance>();
  if (editType == "edit") {
    formRef.current?.setFieldsValue(initValues)
  }
  let noRules;
  if (editType == "save") {
    noRules = [
      {required: true, max: 10, min: 2, type: "string"},
      {
        // @ts-ignore
        validator: (rule, value) => {
          if (!value) {
            return Promise.reject()
          }
          // @ts-ignore
          return roleExists(value).then(({data}) => {
            if (data) {
              return Promise.reject("角色编码已存在")
            }
            return Promise.resolve()
          })
        }
      }
    ]
  }
  return <DrawerForm title={editType === "save" ? "新增角色" : "编辑角色"} visible={editType != null}
                     formRef={formRef}
                     onVisibleChange={(v) => {
                       if (!v) {
                         formRef.current?.resetFields()
                         onCancel()
                       }
                     }}
                     // @ts-ignore
                     onFinish={(values) => {
                       // @ts-ignore
                       values.resourceIds = values?.resourceIds?.filter(r => {
                         return (typeof r) === "number"
                       })
                       console.log(values)
                       onSuccess(values)
                       formRef.current?.resetFields()
                     }}
  >
    <ProFormText name="id" hidden={true}/>
    <ProFormText name="no" label="角色编码" disabled={editType === "edit"}
      // @ts-ignore
                 rules={noRules}></ProFormText>
    <ProFormText name="name" label="角色名称"
      // @ts-ignore
                 rules={roleNameRules} ></ProFormText>
    <ProFormItem name="resourceIds" label="资源" >

      <Auth />
    </ProFormItem>
  </DrawerForm>
}
