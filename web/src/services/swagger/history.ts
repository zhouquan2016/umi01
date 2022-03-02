import {request} from "umi";

export interface History {
  id: number;
  userId: number
  path: string
  method: string
  formData: string
  requestType: string
  requestBody: string
  responseType: string
  responseBody: string
}

export const history = {
  list(params: any) {
    return request<any>("/api/admin/history/list", {
      method: "POST",
      data: params
    })
  },
  indexAll() {
    return request<any>("/api/admin/history/indexAll", {
      method: "GET"
    })
  }
}
