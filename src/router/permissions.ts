/**
 * @description 路由守卫，目前两种模式：all模式与intelligence模式
 */
import { useUserStore } from '@/store/modules/user'
import { useRoutesStore } from '@/store/modules/routes'
import { useSettingsStore } from '@/store/modules/settings'
import VabProgress, { set } from 'nprogress'
import 'nprogress/nprogress.css'
import getPageTitle from '@/utils/pageTitle'
import { authentication, loginInterception, routesWhiteList, supportVisit } from '@/config'
import { Router } from 'vue-router'

export function setupPermissions(router: Router) {
  VabProgress.configure({
    easing: 'ease',
    speed: 500,
    trickleSpeed: 200,
    showSpinner: false,
  })
  const routesStore = useRoutesStore()
  router.beforeEach(async (to: { path: string }, from: any, next: any) => {
    const {
      getTheme: { showProgressBar },
    } = useSettingsStore()
    if (window.location.href.includes('asset-access-insights')) {
      next()
    }
    const { token, getUserInfo, setVirtualRoles, resetAll, logout } = useUserStore()

    if (showProgressBar) VabProgress.start()

    let hasToken = token

    if (!loginInterception) hasToken = true

    if (hasToken) {
      if (routesStore.routes.length) {
        // 禁止已登录用户返回登录页
        if (to.path === '/login') {
          next({ path: '/' })
          if (showProgressBar) VabProgress.done()
        } else next()
      } else {
        try {
          const status = await routesStore.setRoutes(authentication)
          if (status) {
            if (loginInterception) await getUserInfo()
            return next({ ...to, replace: true })
          }
          logout()
          next({ path: '/login', replace: true })
        } catch (err) {
          console.error('vue-admin-beautiful错误拦截:', err)
          await resetAll()
          next({ path: '/login', replace: true })
        }
      }
    } else {
      if (routesWhiteList.includes(to.path)) {
        // 设置游客路由(不需要可以删除)
        if (supportVisit && !routesStore.routes.length) {
          await routesStore.setRoutes('visit')
          next({ path: to.path, replace: true })
        } else next()
      } else next({ path: '/login', replace: true })
    }
  })
  router.afterEach((to: any) => {
    document.title = getPageTitle(to.meta.title)
    if (VabProgress.status) VabProgress.done()
  })
}
