import {request} from "@@/plugin-request/request";
import {MenuData} from "@/services/swagger/user";


export interface MenuTreeNode {
  id: number;
  pId: number;
  value: number;
  title: string
  isLeaf: boolean
  path: string
}

export function menuChildren(
  parentId: number,
  options?: { [key: string]: any }) {
  return request<MenuTreeNode[]>(`/api/menu/children?parentId=${parentId}`, {
    method: 'POST',
    data: parentId,
    ...(options || {}),
  });
}

export function menuAdd(
  params: {
    parentId: number;
    name: string;
    path: string;
    isLeaf: boolean;
  }
) {
  return request<Number>("/api/admin/menu/add", {
    method: "POST",
    data: params
  })
}

export function menuDelete(
  params: number[]
) {
  return request<Boolean>("/api/admin/menu/delete", {
    method: "POST",
    data: params
  })
}

export function menuEdit(
  params: {
    id: number;
    name: string;
    path: string;
  }
) {
  return request<Boolean>("/api/admin/menu/edit", {
    method: "POST",
    data: params
  })
}

export function getMenuById(
  id: number) {
  return request<MenuTreeNode>(`/api/menu/${id}`, {
    method: 'GET',
    data: id,
  });
}

export function existsByPath(
  path: string
) {
  return request<Boolean>(`/api/menu/existsByPath?path=` + path, {
    method: 'GET',
  });
}

export function menuTree() {
  return request<MenuData>("/api/admin/menu/tree", {
    method: "GET"
  })
}
