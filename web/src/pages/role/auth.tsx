import {TreeSelect} from "antd";
import {useEffect, useState} from "react";
import {MenuData, Resource} from "@/services/swagger/user";
import {menuTree} from "@/services/swagger/menu";

interface Props {
  value?: any;
  onChange?: any;
}
export default ({value, onChange}: Props) => {
  const [treeData, setTreeData] = useState<MenuData[]>([])

  const resources2menus = (resources: Resource[]) => {
    return resources?.map(r => {
      return {
        value: r.id,
        label: r.name,
        isLeaf: true
      }
    })
  }

  // @ts-ignore
  const transferTreeData = (menus: MenuData[]) => {
    return menus?.map(m => {
      return {
        value: "menu" + m.id,
        label: m.name,
        isLeaf: false,
        children: m.isLeaf ? resources2menus(m.resources) : transferTreeData(m.children)
      }
    })
  }

  useEffect(() => {
    // @ts-ignore
    menuTree().then(({data}) => {
      const menus = transferTreeData(data)
      console.log(menus)
      setTreeData(menus)
    }).catch((err) => console.error(err))
  }, [])
  return <TreeSelect
    treeData={treeData}
    value={value}
    onChange={onChange}
    treeDefaultExpandAll={true}
    multiple={true}
    treeCheckable={true}
  />
}
