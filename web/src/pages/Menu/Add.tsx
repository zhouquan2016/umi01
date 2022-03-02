import {existsByPath, getMenuById, menuAdd} from "@/services/swagger/menu";
import {useRef} from "react";
import TreeList from "@/pages/Menu/TreeList";
import {RuleObject, StoreValue} from "rc-field-form/lib/interface";
import Resource from "./resources"
import {ModalForm, ProFormDependency, ProFormInstance, ProFormSwitch, ProFormText} from "@ant-design/pro-form";
import ProFormItem from "@ant-design/pro-form/es/components/FormItem";

interface CollectionCreateFormProps {
  visible: boolean;
  onCancel: () => void;
  okCallBack: () => void
}

export default ({
                  visible,
                  onCancel,
                  okCallBack,
                }: CollectionCreateFormProps) => {
  if (!visible) {
    return <div/>
  }
  const form = useRef<ProFormInstance>();


  const validatePath = (rule: RuleObject, value: StoreValue) => {
    // @ts-ignore
    return existsByPath(value).then(({data}) => {
      if (data) {
        return Promise.reject("路径已存在")
      }
      return Promise.resolve()
    })
  }
  return (
    <ModalForm
      visible={visible}
      title="新增菜单"
      formRef={form}
      onVisibleChange={v => !v && onCancel()}
      onFinish={(formData) => {
        // @ts-ignore
        return menuAdd(formData).then(() => {
          okCallBack()
        }).catch((err) => console.error(err))
      }}
      initialValues={
        {
          // isLeaf: false
        }
      }
      onValuesChange={(values) => {
        if (values?.parentId == 0) {
          form.current?.setFieldsValue({
            path: "/"
          })
        }else if (values?.parentId){
          // @ts-ignore
          getMenuById(values?.parentId).then(({data}) => {
            form.current?.setFieldsValue({
              path: data?.path
            })
          })
        }
      }}
    >
      <ProFormItem name="parentId" label="父菜单"  rules={[{required: true, message: '父菜单不能为空!'}]}>
        <TreeList />
      </ProFormItem>
      <ProFormText
        name="name"
        label="名称"
        rules={[{required: true, message: '名称不能为空!'}]}
      >
      </ProFormText>
      <ProFormText
        name="path"
        label="路径"
        rules={[{required: true, message: '路径不能为空!'}, {validator: validatePath}]}
        placeholder={"/开头数字字母_-构成的请求路径，例如/admin/user_a/a-x/1"}
      >
      </ProFormText>
      <ProFormSwitch label="叶子节点" rules={[{required:true, message:"叶子节点不能为空"}]} name="isLeaf"/>
      <ProFormDependency name={["isLeaf", "path"]}>
        {
          // @ts-ignore
          ({isLeaf, path}) => {
            if (isLeaf && path) {
              return <ProFormItem name="resources" label={"资源"}><Resource menuPath={path}></Resource></ProFormItem>
            }
          }
        }
      </ProFormDependency>
    </ModalForm>
  );
};

