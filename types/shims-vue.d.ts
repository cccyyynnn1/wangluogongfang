declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare module '*.vue' {
  import type VueDefineOptions from 'unplugin-vue-define-options/webpack'
  const options: VueDefineOptions<{}, {}, any>
  export default options
}

declare module 'element-plus/dist/locale/zh-cn.mjs'
declare module 'element-plus/dist/locale/en.mjs'
