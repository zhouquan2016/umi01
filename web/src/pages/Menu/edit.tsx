import {ModalForm, ProFormSwitch, ProFormText} from "@ant-design/pro-form";
import React, {useState} from "react";
import Resource from "./resources"
import {getMenuById, menuEdit} from "@/services/swagger/menu";
import ProFormItem from "@ant-design/pro-form/es/components/FormItem";
import {MenuData} from "@/services/swagger/user";
import {message} from "antd";

interface Props {
  menuId: number | any
  clearMenuId: () => void
  onSuccess: () => void
}

export default ({menuId, clearMenuId, onSuccess}: Props) => {
  if (menuId == null) {
    return <div></div>
  }
  const [menu, setMenu] = useState<MenuData | null>(null)
  const loadMenu = () => {
    // @ts-ignore
    return getMenuById(menuId).then(({data}) => {
      setMenu(data)
      return data
    })
  }
  return (
    <React.Fragment>
      <ModalForm visible={true} onVisibleChange={visible => !visible && clearMenuId()}
        // @ts-ignore
                 request={loadMenu}
                 onFinish={formData => {
                   return menuEdit(formData).then(() => {
                     message.info("更新成功")
                     onSuccess()
                   }).catch(err => {
                     console.error(err)
                   })
                 }}
      >
        <ProFormText hidden={true} name="id"/>
        <ProFormText name="deptName" disabled label="父菜单"></ProFormText>
        <ProFormText name="name" label="菜单"></ProFormText>
        <ProFormText name="path" label="路径"></ProFormText>
        <ProFormSwitch name="isLeaf" disabled label="是否叶子节点"/>
        {
          menu?.isLeaf && <ProFormItem name="resources" label={"资源"}>
            <Resource menuPath={menu?.path}/>
          </ProFormItem>
        }
      </ModalForm>
    </React.Fragment>)
}
