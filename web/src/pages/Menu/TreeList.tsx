import {TreeSelect} from "antd";
import {menuChildren} from "@/services/swagger/menu";
import React, {useEffect, useState} from "react";
import {LegacyDataNode} from "rc-tree-select/lib/TreeSelect";

export interface TreeRef {
  reload: () => void;
}





export default ({value, onChange, treeRef}: { value?: any|undefined; onChange?: any|undefined; treeRef?: React.MutableRefObject<TreeRef>|undefined }) => {
  // @ts-ignore
  const rootNode: LegacyDataNode = {
    id: 0,
    value: 0,
    title: "全部",
    isLeaf: false
  }
  const [menus, setMenus] = useState([rootNode])
  const [loadKeys, setLoadKeys] = useState([])
  const [expandKeys, setExpandKeys] = useState([])
  const reload = () => {
    // @ts-ignore
    setExpandKeys([])
    // @ts-ignore
    setLoadKeys([])
    setMenus([rootNode])
  }
  useEffect(() => {
    if (treeRef) {
      treeRef.current = {
        reload
      }
    }

  }, [treeRef])

  const loadMenu = (dataNode: { value: number; }) => {
    return new Promise<void>(resolve => {
      // @ts-ignore
      menuChildren(dataNode.value).then(({data}) => {
        data?.forEach((node: { isLeaf: any; selectable: boolean; }) => {
          if (node.isLeaf) {
            node.selectable = false
          }
        })
        if (data?.length) {
          const newMenus = menus.concat(data)
          setMenus(newMenus)
        }
      }).catch((err) => {
        console.log(err)
      }).finally(() => {
        resolve();
      })
    })
  }
  return (
    <TreeSelect treeDataSimpleMode allowClear={true} value={value} onChange={onChange}
      treeData={menus}
      // @ts-ignore
      loadData={loadMenu}
      treeLoadedKeys={loadKeys}
      onTreeLoad={(keys) => {
        // @ts-ignore
        setLoadKeys(keys)
      }}
      treeExpandedKeys={expandKeys}
      // @ts-ignore
      onTreeExpand={(keys) => {
        // @ts-ignore
        setExpandKeys(keys)
      }}
    />
  );
};
