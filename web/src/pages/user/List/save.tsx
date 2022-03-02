import {DrawerForm, ProFormInstance, ProFormSelect, ProFormText} from "@ant-design/pro-form";
import {useEffect, useRef} from "react";
import {user} from "@/services/swagger/user";
import {role} from "@/services/swagger/role";

export type SaveType = "save"|"edit"

interface Props {
  visible: boolean
  id: number| null
  onCancel: () => void
  onOk: () => void
}
export default ({visible, id, onCancel, onOk}: Props) => {
  const userFrom = useRef<ProFormInstance>()

  useEffect(() => {
    if (!visible) {
      return
    }
    if (id == null) {
      userFrom.current?.resetFields()
    }else {
      user.getById(id).then(({data}) => {
        userFrom.current?.setFieldsValue(data)
      }).catch(err => {
        console.log(err)
      })
    }
  }, [visible])

  return <DrawerForm
    visible={visible} onVisibleChange={v => !v && onCancel()}
    formRef={userFrom}
    // @ts-ignore
    onFinish={(formData) => {
      let op: Promise<any>
      if (id == null) {
        op = user.add(formData)
      }else {
        op = user.edit(formData)
      }
      op.then(() => {
        onOk()
      }).catch(err => {
        console.log(err)
      })

    }}
  >
    {
      id != null && <ProFormText name="id" hidden={true}/>
    }
    <ProFormText name="name" label="姓名" rules={[{required:true, type:"string", min:2, max:10}]}/>
    {
      id == null && <ProFormText.Password name="password" label="密码" rules={[{required:true,type:"string", min:2, max:20, pattern:new RegExp("^[a-zA-Z0-9_-]{2,20}$")}]}/>
    }
    {
      id == null && <ProFormText name="email" label="邮箱" rules={[{required:true, pattern: new RegExp("^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")}]}/>
    }
    <ProFormText name="signature" label="个性签名" rules={[{type:"string", max:20}]}/>
    <ProFormText name="title" label="职位" rules={[{type:"string", max:10}]}/>
    <ProFormText name="address" label="地址" rules={[{type:"string", max:20}]}/>
    {
      id == null && <ProFormText name="phone" label="手机号" rules={[{type:"string", required:true, pattern:new RegExp("^1\\d{10}$")}]}/>
    }
    <ProFormSelect name="roleId" label="角色" rules={[{type:"number", required:true}]}
      // @ts-ignore
                   fieldProps={{
                     fieldNames: {
                       label: "name",
                       value: "id"
                     },
                   }}
                   // @ts-ignore
                   request={() => {
                     // @ts-ignore
                     return role.getAll().then(({data}) => data).catch(err => console.log(err))
                   }}
    />
  </DrawerForm>
}
