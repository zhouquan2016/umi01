export default [
  {
    path: '/user',
    layout: false,
    routes: [
      {
        path: '/user',
        routes: [
          {
            name: 'login',
            path: '/user/login',
            component: './user/Login',
          },
          {
            name: 'resetPwd',
            path: '/user/resetPwd',
            component: './user/ResetPwd',
          },
        ],
      },
      {
        component: './404',
      },
    ],
  },
  {
    path: '/',
    name: 'welcome',
    icon: 'smile',
    component: './Welcome',
  },
  {
    path:"/admin",
    name:"admin",
    routes: [
      {
        path: "/admin/menu",
        name:"admin.menu",
        component: "./Menu"
      },
      {
        path: "/admin/role",
        name:"admin.role",
        component: "./role"
      },
      {
        path: "/admin/user",
        name:"admin.user",
        component: "./user/List"
      },
      {
        path: "/admin/history",
        name:"admin.history",
        component: "./history"
      },
      {
        component: './404',
      }
    ]
  },
  {
    path: "/product",
    routes: [
      {
        path: "list",
        component: "./product"
      },
      {
        component: './404',
      }
    ]
  },

  {
    component: './404',
  },
];
