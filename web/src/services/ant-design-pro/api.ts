// @ts-ignore
/* eslint-disable */
import {request} from 'umi';
import {API} from "@/services/ant-design-pro/typings";
import {MenuData} from "@/services/swagger/user";

/** 获取当前的用户 GET /api/currentUser */
export async function currentUser(options?: { [key: string]: any }) {
  return request<{
    data: API.CurrentUser;
  }>('/api/currentUser', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 退出登录接口 POST /api/login/outLogin */
export async function outLogin(options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/login/outLogin', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 登录接口 POST /api/login/account */
export async function login(body: API.LoginParams, options?: { [key: string]: any }) {
  return request<API.LoginSuccess>('/api/login/account', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

export async function getMenus(body: {
  userId: number;
}, options?: { [key: string]: any }) {
  return request<API.MenuData[]>('/api/getMenuByUser', {
    method: 'POST',
    ...(options || {}),
  });
}

export function menuList(
  params: {
    sortField: string;
    sortOrder: string;
    pageSize: number;
    current: number;
  },
  options?: { [key: string]: any }) {
  return request<MenuData>(`/api/menu/list`, {
    method: 'POST',
    data: params,
    ...(options || {}),
  });
}

export async function userList(
  params: {
    sortField: string;
    sortOrder: string;
    pageSize: number;
    current: number;
  }
) {
  return request<MenuData>(`/api/user/list`, {
    method: 'POST',
    data: params,
  });
}
