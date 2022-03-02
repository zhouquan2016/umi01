// @ts-ignore
/* eslint-disable */
import {request} from 'umi';
import {MenuDataItem} from "@ant-design/pro-layout";
/** Create user This can only be done by the logged in user. POST /user */
export async function createUser(body: API.User, options?: { [key: string]: any }) {
  return request<any>('/user', {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/** Creates list of users with given input array POST /user/createWithArray */
export async function createUsersWithArrayInput(
  body: API.User[],
  options?: { [key: string]: any },
) {
  return request<any>('/user/createWithArray', {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/** Creates list of users with given input array POST /user/createWithList */
export async function createUsersWithListInput(body: API.User[], options?: { [key: string]: any }) {
  return request<any>('/user/createWithList', {
    method: 'POST',
    data: body,
    ...(options || {}),
  });
}

/** Logs user into the system GET /user/login */
export async function loginUser(
  params: {
    // query
    /** The user name for login */
    username: string;
    /** The password for login in clear text */
    password: string;
  },
  options?: { [key: string]: any },
) {
  return request<string>('/user/login', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** Logs out current logged in user session GET /user/logout */
export async function logoutUser(options?: { [key: string]: any }) {
  return request<any>('/user/logout', {
    method: 'GET',
    ...(options || {}),
  });
}

/** Get user by user name GET /user/${param0} */
export async function getUserByName(
  params: {
    // path
    /** The name that needs to be fetched. Use user1 for testing.  */
    username: string;
  },
  options?: { [key: string]: any },
) {
  const {username: param0} = params;
  return request<API.User>(`/user/${param0}`, {
    method: 'GET',
    params: {...params},
    ...(options || {}),
  });
}

/** Updated user This can only be done by the logged in user. PUT /user/${param0} */
export async function updateUser(
  params: {
    // path
    /** name that need to be updated */
    username: string;
  },
  body: API.User,
  options?: { [key: string]: any },
) {
  const {username: param0} = params;
  return request<any>(`/user/${param0}`, {
    method: 'PUT',
    params: {...params},
    data: body,
    ...(options || {}),
  });
}

/** Delete user This can only be done by the logged in user. DELETE /user/${param0} */
export async function deleteUser(
  params: {
    // path
    /** The name that needs to be deleted */
    username: string;
  },
  options?: { [key: string]: any },
) {
  const {username: param0} = params;
  return request<any>(`/user/${param0}`, {
    method: 'DELETE',
    params: {...params},
    ...(options || {}),
  });
}

export async function menuTree(
  options?: { [key: string]: any }) {
  return request<MenuDataItem>(`/api/menu/tree`, {
    method: 'GET',
    ...(options || {}),
  });
}
export interface Resource {
  id: number;
  code: string;
  name: string;
  path: string;
  menuId: number;
  isSysDefault: boolean;
}
export interface MenuData {
  id: number;
  name: string;
  isLeaf: boolean;
  isSysDefault: boolean;
  path: string;
  parentId: number;
  dept: string;
  deptName: string;
  resources: Resource[];
  children: MenuData[];
}
export function menuList(
  params: {
    sortField: string;
    sortOrder: string;
    pageSize: number;
    current: number;
  },
  options?: { [key: string]: any }) {
  return request<MenuData>(`/api/admin/menu/list`, {
    method: 'POST',
    data: params,
    ...(options || {}),
  });
}

export const user = {
  getById(id: number) {
    return request<any>("/api/admin/user/getById?id=" + id, {
      method: "GET"
    })
  },
  add(formData: any) {
    return request<any>("/api/admin/user/save", {method:"POST", data:formData})
  },
  edit(formData: any) {
    return request<any>("/api/admin/user/edit", {method:"POST", data:formData})
  },
  deleteById(id: number) {
    return request<any>("/api/admin/user/deleteById?id=" + id, {method:"GET"})
  }
}
