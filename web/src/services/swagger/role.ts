import {request} from "umi";

export interface Role {
  id: number;
  no: string;
  name: string;
  isSysDefault: boolean;
}

export function roleList(params: any) {
  return request<Role[]>("/api/admin/role/list", {
    method:"POST",
    data: params
  })
}

export function roleAdd(params: {
  no: string,
  code: string
}) {
  return request<boolean>("/api/admin/role/add", {
    method:"POST",
    data: params
  })
}

export function roleExists(no: string) {
  return request<boolean>("/api/admin/role/exists?no=" + no, {
    method:"GET",
  })
}

export function roleDelete(id: number) {
  return request<boolean>("/api/admin/role/delete?id=" + id, {
    method:"GET",
  })
}

export function roleEdit(params: {
  name: string,
  id: number
}) {
  return request<boolean>("/api/admin/role/edit", {
    method: "POST",
    data: params
  })
}
export function roleGetById(id: number) {
  return request<Role>("/api/admin/role/getById?id=" + id, {
    method:"GET"
  })
}

export const role = {
  getAll(){
    return request<Role[]>("/api/admin/role/getAll")
  }
}
