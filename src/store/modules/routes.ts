/**
 * @description 路由拦截状态管理，目前两种模式：all模式与intelligence模式，其中partialRoutes是菜单暂未使用
 */
import { asyncRoutes, constantRoutes, resetRouter } from '@/router'
import { filterRoutes } from '@/utils/routes'
import { rolesControl } from '@/config'
import { OptionType, RoutesModuleType } from '/#/store'
import { VabRouteRecord } from '/#/router'
import { useUserStore } from './user'
import { getSystemMenusByRoleApi } from '@/api-ecs/system'
export const useRoutesStore = defineStore('routes', {
  state: (): RoutesModuleType => ({
    /**
     * 一级菜单值
     */
    tab: {
      data: undefined,
    },
    /**
     * 一级菜单
     */
    tabMenu: undefined,
    /**
     * 自定义激活菜单
     */
    activeMenu: {
      data: undefined,
    },
    /**
     * 一级菜单
     */
    routes: [],
  }),
  getters: {
    getTab: (state) => state.tab,
    getTabMenu: (state) =>
      state.tab.data
        ? state.routes.find((route) => route.name === state.tab.data)
        : { meta: { title: '' }, redirect: '404' },
    getActiveMenu: (state) => state.activeMenu,
    getRoutes: (state) => state.routes.filter((_route) => _route.meta.hidden !== true),
    getPartialRoutes: (state) => state.routes.find((route) => route.name === state.tab.data)?.children || [],
  },
  actions: {
    clearRoutes() {
      this.routes = []
    },
    /**
     * @description 多模式设置路由
     * @param mode
     * @returns
     */
    async setRoutes(mode = 'none') {
      const userStore = useUserStore()
      // 默认前端路由
      let routes: VabRouteRecord[] = []
      // 设置游客路由关闭路由拦截(不需要可以删除)
      const control = mode === 'visit' ? false : rolesControl
      if (userStore.token) {
        const { data } = await getSystemMenusByRoleApi()
        if (!data || data?.length === 0) {
          resetRouter(constantRoutes)
          alert('当前用户暂未配置菜单权限，请配置后重试')
          // gp.$baseMessage('当前用户暂未配置菜单，请配置后重试', 'error', 'vab-hey-message-error')
          return false
        }
        constantRoutes[0].redirect = data[0].path
        routes = filterRoutesByMenus(asyncRoutes, Array.isArray(data) ? data : [])
      }
      // 根据权限和rolesControl过滤路由
      const accessRoutes = filterRoutes(
        [
          ...constantRoutes,
          ...routes,
          {
            path: '/:pathMatch(.*)*',
            redirect: '/404',
            name: 'NotFound',
            meta: {
              hidden: true,
            },
          },
        ],
        control
      )
      // 设置菜单所需路由
      this.routes = JSON.parse(JSON.stringify(accessRoutes))
      // 根据可访问路由重置Vue Router
      resetRouter(accessRoutes)

      return true
    },
    changeMenuMeta(options: OptionType) {
      function handleRoutes(routes: VabRouteRecord[]) {
        return routes.map((route) => {
          if (route.name === options.name) Object.assign(route.meta, options.meta)
          if (route.children && route.children.length) route.children = handleRoutes(route.children)
          return route
        })
      }
      this.routes = handleRoutes(this.routes)
    },
    /**
     * @description 修改 activeName
     * @param activeMenu 当前激活菜单
     */
    changeActiveMenu(activeMenu: string) {
      this.activeMenu.data = activeMenu
    },
  },
})

const filterRoutesByMenus = (routers: VabRouteRecord[], menus: any): VabRouteRecord[] => {
  const router = []
  for (const menu of menus) {
    const { name, path, children, meta } = menu
    const curRoute = routers.find((i) => i.name === name && i.path === path && !meta.hidden)
    if (!curRoute) continue
    curRoute.meta.hidden = false
    if (Array.isArray(children) && curRoute.children) {
      router.push({
        ...curRoute,
        children: filterRoutesByMenus(curRoute.children!, children),
      })
    } else {
      router.push(curRoute)
    }
  }
  return router
}
