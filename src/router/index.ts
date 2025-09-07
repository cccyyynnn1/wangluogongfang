/**
 * @description router全局配置，如有必要可分文件抽离，其中asyncRoutes只有在intelligence模式下才会用到，pro版只支持remixIcon图标，具体配置请查看vip群文档
 */
import type { RouteRecordName, RouteRecordRaw } from 'vue-router'
import type { VabRouteRecord } from '/#/router'
import { createRouter, createWebHashHistory, createWebHistory } from 'vue-router'
import Layout from '@vab/layouts/index.vue'
import { setupPermissions } from './permissions'
import { authentication, isHashRouterMode, publicPath } from '@/config'

export const constantRoutes: VabRouteRecord[] = [
  {
    path: '/',
    name: 'Root',
    meta: {
      hidden: true,
    },
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/ecs/login/index.vue'),
    meta: {
      hidden: true,
    },
  },
  {
    path: '/activate',
    name: '激活',
    component: () => import('@/ecs/public/index.vue'),
    meta: {
      hidden: true,
    },
  },
  {
    path: '/help',
    name: 'Markdown',
    component: () => import('@/ecs/public/markdown.vue'),
    meta: {
      title: '帮助文档',
      icon: 'home-2-line',
      hidden: true,
    },
  },
  {
    path: '/dashboard-screen',
    name: 'DashboardScreen',
    component: () => import('@/ecs/dashboard/index.vue'),
    meta: {
      title: '数据大屏',
      hidden: true,
    },
  },
  {
    path: '/screen-highway/liaoning',
    name: 'ScreenHighwayLiaoning',
    component: () => import('@/ecs/screen-highway/liaoning.vue'),
    meta: {
      title: '数据大屏',
      hidden: true,
    },
  },
  {
    path: '/asset-access-insights',
    name: 'AssetAccessInsights',
    component: () => import('@/ecs/asset-access-insights/index.vue'),
    meta: {
      title: '资产访问洞察',
      hidden: true,
    },
  },
  {
    path: '/403',
    name: '403',
    component: () => import('@/ecs/403.vue'),
    meta: {
      hidden: true,
    },
  },
  {
    path: '/404',
    name: '404',
    component: () => import('@/ecs/404.vue'),
    meta: {
      hidden: true,
    },
  },
]

export const asyncRoutes: VabRouteRecord[] = [
  {
    path: '/index',
    name: 'Index',
    component: Layout,
    meta: {
      title: '主页',
      icon: 'home-2-line',
    },
    children: [
      {
        path: '/dashboard',
        name: 'Workbench',
        component: () => import('@/ecs/index/workbench.vue'),
        meta: {
          title: '工作台',
          icon: 'computer-line',
        },
      },
      {
        path: '/packet-replay',
        name: 'PackeReplay',
        component: () => import('@/ecs/index/packet-replay.vue'),
        meta: {
          title: '数据包回放',
          icon: 'slideshow-4-line',
        },
      },
      {
        path: '/monitoring',
        name: 'InterfaceMonitoring',
        component: () => import('@/ecs/index/interface-monitoring.vue'),
        meta: {
          title: '任务仪表盘',
          icon: 'dashboard-2-line',
        },
      },
      {
        path: '/situation-screen',
        name: 'SituationScreen',
        component: () => import('@/ecs/index/situation-screen.vue'),
        meta: {
          title: '态势',
          icon: 'funds-line',
        },
      },
    ],
  },
  {
    path: '/site',
    name: 'Site',
    component: Layout,
    meta: {
      title: '站点',
      icon: 'dropbox-line',
    },
    children: [
      {
        path: '/site_index',
        name: 'SiteIndex',
        component: () => import('@/ecs/site/index.vue'),
        meta: {
          title: '站点',
          icon: 'home-2-line',
          noColumn: true,
        },
      },
    ],
  },
  {
    path: '/retrieve',
    name: 'Retrieve',
    component: Layout,
    meta: {
      title: '调查',
      icon: 'globe-line',
    },
    children: [
      {
        path: 'index',
        name: 'retrieveIndex',
        component: () => import('@/ecs/retrieve/index.vue'),
        meta: {
          title: '应用层会话',
          icon: 'home-2-line',
        },
      },
      {
        path: 'network',
        name: 'retrieveNetwork',
        component: () => import('@/ecs/retrieve/network-layer.vue'),
        meta: {
          title: '网络层会话',
          icon: 'home-2-line',
        },
      },
    ],
  },
  // {
  //   path: '/chase',
  //   name: 'Chase',
  //   component: Layout,
  //   meta: {
  //     title: '狩猎',
  //     icon: 'speaker-2-line',
  //   },
  //   children: [
  //     {
  //       path: 'index',
  //       name: 'Chase',
  //       component: () => import('@/ecs/chase/index.vue'),
  //       meta: {
  //         title: '狩猎',
  //         icon: 'speaker-2-line',
  //       },
  //     },
  //   ],
  // },
  {
    path: '/alerts',
    name: 'Alerts',
    component: Layout,
    meta: {
      title: '告警',
      icon: 'reactjs-line',
    },
    children: [
      // {
      //   path: 'index',
      //   name: 'AlertsIndex',
      //   component: () => import('@/ecs/alert/index.vue'),
      //   meta: {
      //     title: '原始告警',
      //     icon: 'folder-chart-line',
      //   },
      // },
      {
        path: 'alarm-attack',
        name: 'AlarmAttack',
        component: () => import('@/ecs/alert/alarm-aggregation.vue'),
        meta: {
          title: '威胁告警',
          icon: 'shape-2-line',
          // noColumn: true,
        },
      },
      {
        path: 'threat-intelligence',
        name: 'AlarmIntelligence',
        component: () => import('@/ecs/alert/alarm-intelligence-center.vue'),
        meta: {
          title: '情报中心',
          icon: 'slideshow-line',
        },
      },
      {
        path: 'mail-analysis',
        name: 'MailAnalysis',
        component: () => import('@/ecs/alert/email-analysis.vue'),
        meta: {
          title: 'AI邮件分析',
          icon: 'slideshow-line',
        },
      },
      {
        path: 'alarm-abnormal-login',
        name: 'AbnormalLogin',
        component: () => import('@/ecs/alert/alarm-abnormal-login.vue'),
        meta: {
          title: '威胁专题',
          icon: 'bug-2-line',
        },
      },
      // {
      //   path: 'tunnel-topics',
      //   name: 'TunnelTopics',
      //   component: () => import('@/ecs/alert/tunnel-topics.vue'),
      //   meta: {
      //     title: '隧道专题',
      //     icon: 'slideshow-line',
      //     hidden: true,
      //   },
      // },
    ],
  },
  {
    path: '/assets',
    name: 'Assets',
    component: Layout,
    meta: {
      title: '资产',
      icon: 'keyboard-box-line',
    },
    children: [
      {
        path: 'index',
        name: 'AssetOverview',
        component: () => import('@/ecs/assets/components/asset-overview/index.vue'),
        meta: {
          title: '资产总览',
          icon: 'chat-4-line',
          noColumn: true,
        },
      },
      {
        path: 'asset-visits',
        name: 'AssetVisits',
        component: () => import('@/ecs/assets/components/asset-visits.vue'),
        meta: {
          title: '资产访问',
          icon: 'chat-forward-line',
          noColumn: true,
        },
      },
      {
        path: 'assets-site',
        name: 'AssetsSite',
        component: () => import('@/ecs/assets/components/assets-site.vue'),
        meta: {
          title: '已知站点',
          icon: 'chat-check-line',
          noColumn: true,
        },
      },
      {
        path: 'assets-unknow',
        name: 'AssetsUNSite',
        component: () => import('@/ecs/assets/components/assets-unsite.vue'),
        meta: {
          title: '未知站点',
          icon: 'chat-delete-line',
          noColumn: true,
        },
      },
      {
        path: 'network-partition',
        name: 'NetworkPartition',
        component: () => import('@/ecs/assets/components/network-partition.vue'),
        meta: {
          title: '网络分区',
          icon: 'broadcast-fill',
          noColumn: true,
        },
      },
    ],
  },
  // {
  //   path: '/algorithms',
  //   name: 'Algorithms',
  //   component: Layout,
  //   meta: {
  //     title: '算法',
  //     icon: 'steam-line',
  //   },
  //   children: [
  //     {
  //       path: 'index',
  //       name: 'AlgorithmsIndex',
  //       component: () => import('@/ecs/algorithms/index.vue'),
  //       meta: {
  //         title: '算法',
  //         icon: 'home-2-line',
  //         noColumn: true,
  //       },
  //     },
  //   ],
  // },
  {
    path: '/knowledge',
    name: 'Knowledge',
    component: Layout,
    meta: {
      title: '知识',
      icon: 'book-read-line',
    },
    children: [
      {
        path: 'index',
        name: 'KnowledgeIndex',
        component: () => import('@/ecs/knowledge/index.vue'),
        meta: {
          title: '知识',
          icon: 'book-read-line',
          noColumn: true,
        },
      },
    ],
  },
  {
    path: '/config',
    name: 'Config',
    component: Layout,
    meta: {
      title: '配置',
      icon: 'slack-line',
    },
    children: [
      {
        path: 'data-extraction',
        name: 'DataExtraction',
        meta: {
          title: '数据提取',
          icon: 'home-2-line',
        },
        children: [
          {
            path: 'index',
            name: 'DataExtractionIndex',
            component: () => import('@/ecs/config/data-extraction/index.vue'),
            meta: {
              title: '自定义字段',
              icon: 'home-2-line',
            },
          },
          {
            path: 'extraction-rule',
            name: 'ExtractionRule',
            component: () => import('@/ecs/config/data-extraction/extraction-rule.vue'),
            meta: {
              title: '提取规则',
              icon: 'home-2-line',
            },
          },
        ],
      },
      {
        path: 'data-bus',
        name: 'DataBus',
        meta: {
          title: '数据总线',
          icon: 'home-2-line',
        },
        children: [
          {
            path: 'index',
            name: 'DataBusIndex',
            component: () => import('@/ecs/config/data-bus/index.vue'),
            meta: {
              title: 'KAFKA配置',
              icon: 'home-2-line',
            },
          },
        ],
      },
      // {
      //   path: 'data-switching',
      //   name: 'DataSwitching',
      //   meta: {
      //     title: '数据转发',
      //     icon: 'home-2-line',
      //   },
      //   children: [
      //     {
      //       path: 'index',
      //       name: 'DataSwitchingIndex',
      //       component: () => import('@/ecs/config/data-switching/index.vue'),
      //       meta: {
      //         title: '转发规则',
      //         icon: 'home-2-line',
      //       },
      //     },
      //   ],
      // },
      // todo
      {
        path: 'alert-configuration',
        name: 'AlertConfiguration',
        meta: {
          title: '告警配置',
          icon: 'macbook-line',
        },
        children: [
          {
            path: 'white-list',
            name: 'AlertWhiteList',
            component: () => import('@/ecs/config/alert-configuration/index.vue'),
            meta: {
              title: '白名单管理',
              icon: 'file-list-3-line',
            },
          },
          {
            path: 'custom-rules',
            name: 'CustomRules',
            component: () => import('@/ecs/config/alert-configuration/custom-rules.vue'),
            meta: {
              title: '自定义告警',
              icon: 'file-shield-2-line',
            },
          },
          {
            path: 'information-maintenance',
            name: 'InformationMaintenance',
            component: () => import('@/ecs/config/alert-configuration/information-maintenance.vue'),
            meta: {
              title: '信息维护',
              icon: 'file-edit-line',
            },
          },
          {
            path: 'attack-characterization',
            name: 'AttackCharacterization',
            component: () => import('@/ecs/config/alert-configuration/attack-characterization.vue'),
            meta: {
              title: '高亮特征配置',
              icon: 'shape-2-line',
            },
          },
        ],
      },
      {
        path: 'assets-configuration',
        name: 'AssetsConfiguration',
        meta: {
          title: '资产配置',
          icon: 'keyboard-box-line',
        },
        children: [
          {
            path: 'assets-dic',
            name: 'AssetsDic',
            component: () => import('@/ecs/config/assets-configuration/assets-dic.vue'),
            meta: {
              title: '资产字典',
              icon: 'list-settings-line',
            },
          },
          {
            path: 'asset-import',
            name: 'AssetImport',
            component: () => import('@/ecs/config/assets-configuration/asset-import.vue'),
            meta: {
              title: '资产导入',
              icon: 'download-line',
            },
          },
        ],
      },
      {
        path: 'init-configration',
        name: 'InitConfigration',
        meta: {
          title: '初始化配置',
          icon: 'folder-settings-line',
        },
        children: [
          {
            path: 'field-initialisation',
            name: 'FieldInitialisation',
            component: () => import('@/ecs/config/init-configration/field-initialisation.vue'),
            meta: {
              title: '字段初始化',
              icon: 'font-size-2',
            },
          },
          {
            path: 'backup-restore',
            name: 'BackupRestore',
            component: () => import('@/ecs/config/init-configration/backup-restore.vue'),
            meta: {
              title: '备份与还原',
              icon: 'refund-line',
            },
          },
        ],
      },
    ],
  },
  {
    path: '/managements',
    name: 'Management',
    component: Layout,
    meta: {
      title: '系统',
      icon: 'settings-4-line',
    },
    children: [
      {
        path: 'user_management',
        name: 'UserManagement',
        meta: {
          title: '用户管理',
          icon: 'remixicon-line',
        },
        children: [
          {
            path: 'user',
            name: 'User',
            component: () => import('@/ecs/management/user/index.vue'),
            meta: {
              title: '用户管理',
            },
          },
          {
            path: 'role',
            name: 'Role',
            component: () => import('@/ecs/management/user/role-management.vue'),
            meta: {
              title: '角色管理',
            },
          },
          {
            path: 'department',
            name: 'Department',
            component: () => import('@/ecs/management/user/department-management.vue'),
            meta: {
              title: '部门管理',
            },
          },
        ],
      },
      {
        path: 'audit',
        name: 'Audit',
        meta: {
          title: '审计管理',
          icon: 'user-3-line',
        },
        children: [
          {
            path: 'operation',
            name: 'Operation',
            component: () => import('@/ecs/management/audit/operation.vue'),
            meta: {
              title: '系统日志',
            },
          },
          {
            path: 'systemLogs',
            name: 'SystemLogs',
            component: () => import('@/ecs/management/audit/system-log.vue'),
            meta: {
              title: '应用日志',
            },
          },
          {
            path: 'logingLogs',
            name: 'LogingLogs',
            component: () => import('@/ecs/management/audit/login.vue'),
            meta: {
              title: '登录日志',
            },
          },
        ],
      },
      {
        path: 'configuration',
        name: 'Configuration',
        meta: {
          title: '配置管理',
          icon: 'table-2',
        },
        children: [
          {
            path: 'index',
            name: 'ConfigManagement',
            component: () => import('@/ecs/management/configuration/index.vue'),
            meta: {
              title: '系统信息',
            },
          },
          {
            path: 'systemConfig',
            name: 'SystemConfig',
            component: () => import('@/ecs/management/configuration/system-config.vue'),
            meta: {
              title: '系统配置',
            },
          },

          {
            path: 'authorization',
            name: 'AuthorizationManagement',
            component: () => import('@/ecs/management/configuration/authorization-management.vue'),
            meta: {
              title: '授权管理',
            },
          },
          {
            path: 'fullFlowConfig',
            name: 'FullFlowConfig',
            component: () => import('@/ecs/management/configuration/fullFlow-config.vue'),
            meta: {
              title: '系统升级',
            },
          },
          // {
          //   path: 'rules-update',
          //   name: 'RulesUpdate',
          //   component: () => import('@/ecs/management/configuration/rules-update.vue'),
          //   meta: {
          //     title: '规则更新',
          //   },
          // },
          {
            path: 'log-config',
            name: 'SyslogConfig',
            component: () => import('@/ecs/management/configuration/syslog-config.vue'),
            meta: {
              title: '日志采集',
            },
          },
          {
            path: '/dataDictionaryConfig/:id',
            name: 'DataDictionaryConfig',
            component: () => import('@/ecs/management/configuration/dictionary-config.vue'),
            meta: {
              title: '字典配置',
            },
          },

          {
            path: 'log-data-output',
            name: 'LogDataOutput',
            component: () => import('@/ecs/management/configuration/syslog-data-output.vue'),
            meta: {
              title: '日志数据输出',
            },
          },
          {
            path: 'paramsConfig',
            name: 'ParamsConfig',
            component: () => import('@/ecs/management/configuration/params-config.vue'),
            meta: {
              title: '统一参数配置',
            },
          },
          {
            path: 'api-authorization',
            name: 'APIAuthorization',
            component: () => import('@/ecs/management/configuration/api-authorization/index.vue'),
            meta: {
              title: 'API授权管理',
            },
          },
        ],
      },
      {
        path: 'systemCondition',
        name: 'SystemCondition',
        meta: {
          title: '系统状态',
          icon: 'list-check-2',
        },
        children: [
          {
            path: 'index',
            // name: 'SystemCondition', 重名
            name: 'SystemConditionStatus',
            component: () => import('@/ecs/management/system-condition/index.vue'),
            meta: {
              title: '本机系统状态',
            },
          },
        ],
      },
      {
        path: 'equipment',
        name: 'Equipment',
        meta: {
          title: '设备管理',
          icon: 'calendar-check-line',
        },
        children: [
          {
            path: '/equipment/index',
            // name: 'Equipment', 重名
            name: 'EquipmentList',
            component: () => import('@/ecs/management/equipment/index.vue'),
            meta: {
              title: '基础配置',
              icon: 'calendar-check-line',
            },
          },
          {
            path: '/equipment/collecting-device',
            name: 'CollectingDevice',
            component: () => import('@/ecs/management/equipment/collecting-device.vue'),
            meta: {
              title: '采集设备',
              icon: 'calendar-check-line',
            },
          },
          {
            path: '/equipment/traffic-link',
            name: 'TrafficLink',
            component: () => import('@/ecs/management/equipment/traffic-link.vue'),
            meta: {
              title: '流量链路',
              icon: 'calendar-check-line',
            },
          },
          {
            path: '/equipment/equipment-replacement',
            name: 'EquipmentReplacement',
            component: () => import('@/ecs/management/equipment/equipment-replacement.vue'),
            meta: {
              title: '流量设备更新',
              icon: 'calendar-check-line',
            },
          },
          {
            path: '/equipment/network-information',
            name: 'NetworkInformation',
            component: () => import('@/ecs/management/equipment/network-information.vue'),
            meta: {
              title: '设备网卡信息',
              icon: 'calendar-check-line',
            },
          },
          {
            path: '/equipment/kafka-push',
            name: 'KafkaPush',
            component: () => import('@/ecs/management/equipment/kafka-push.vue'),
            meta: {
              title: 'Kafka推送',
              icon: 'calendar-check-line',
            },
          },
          {
            path: '/equipment/system-service',
            name: 'SystemService',
            component: () => import('@/ecs/management/equipment/system-service.vue'),
            meta: {
              title: '系统服务',
              icon: 'calendar-check-line',
            },
          },
          {
            path: '/equipment/load-balancing',
            name: 'LoadBalancing',
            component: () => import('@/ecs/management/equipment/load-balancing.vue'),
            meta: {
              title: '负载均衡管理',
              icon: 'calendar-check-line',
            },
          },
          {
            path: '/equipment/VS-list',
            name: 'VSList',
            component: () => import('@/ecs/management/equipment/VS-list.vue'),
            meta: {
              title: 'VS列表',
              icon: 'calendar-check-line',
            },
          },
          {
            path: '/equipment/POOL-list',
            name: 'POOLList',
            component: () => import('@/ecs/management/equipment/POOL-list.vue'),
            meta: {
              title: 'POOL列表',
              icon: 'calendar-check-line',
            },
          },
        ],
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404',
    name: 'NotFound',
    meta: {
      hidden: true,
    },
  },
]

const router = createRouter({
  history: isHashRouterMode ? createWebHashHistory(publicPath) : createWebHistory(publicPath),
  routes: constantRoutes as RouteRecordRaw[],
})

function fatteningRoutes(routes: VabRouteRecord[]): VabRouteRecord[] {
  return routes.flatMap((route: VabRouteRecord) => {
    return route.children ? fatteningRoutes(route.children) : route
  })
}

function addRouter(routes: VabRouteRecord[]) {
  routes.forEach((route: VabRouteRecord) => {
    if (!router.hasRoute(route.name)) router.addRoute(route as RouteRecordRaw)
    if (route.children) addRouter(route.children)
  })
}

export function resetRouter(routes: VabRouteRecord[] = constantRoutes) {
  routes.map((route: VabRouteRecord) => {
    if (route.children) route.children = fatteningRoutes(route.children)
  })
  // router.getRoutes().forEach(({ name }) => {
  //   router.hasRoute(<RouteRecordName>name) && router.removeRoute(<RouteRecordName>name)
  // })
  router.getRoutes().forEach(({ name }) => {
  // 增加一个 if 判断，确保 name 不是 undefined
  if (name) {
    router.hasRoute(name) && router.removeRoute(name);
  }
  })
  
  addRouter(routes)
}

export function setupRouter(app: any) {
  if (authentication === 'intelligence') addRouter(asyncRoutes)
  setupPermissions(router)
  app.use(router)
  return router
}

export default router
