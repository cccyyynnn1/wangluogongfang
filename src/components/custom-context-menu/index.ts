import { App, DirectiveBinding, createVNode, render, ComponentPublicInstance, ObjectDirective } from 'vue'
import CustomContextMenu from './index.vue'
CustomContextMenu.install = (app: App): void => {
  app.component(CustomContextMenu.name as string, CustomContextMenu)
}
export type MenuCallback<T = any> = (
  arg0?: any,
  arg1?: HTMLElement | null,
  arg2?: HTMLElement | null,
  arg3?: MouseEvent
) => T
export type MenuSetting = {
  menuClick?: MenuCallback
  label?: string | MenuCallback<string>
  tips?: string | MenuCallback<string>
  hidden?: boolean | MenuCallback<boolean>
  disabled?: boolean | MenuCallback<boolean>
  line?: boolean
  children?: MenuSetting[]
  customClass?: string
}

function createDom(tag: string, className: string, innerText?: string) {
  let el = document.createElement(tag)
  el.setAttribute('class', className)
  if (innerText) el.innerText = innerText
  return el
}

function CreateCustomMouseMenu(options: any) {
  const className = 'custom-context-menu-container'
  let container: HTMLElement
  if (document.querySelector(`.${className}`)) {
    container = document.querySelector(`.${className}`) as HTMLElement
  } else {
    container = createDom('div', className)
  }
  const vm = createVNode(CustomContextMenu, options)
  render(vm, container)
  document.body.appendChild(container)
  return vm.component?.proxy as ComponentPublicInstance<typeof CustomContextMenu>
}

let ContextMenuMenuCtx: ComponentPublicInstance<typeof CustomContextMenu>

// 指令封装
let contextMenuEvent: (e: MouseEvent) => void
const mounted = (el: HTMLElement, binding: DirectiveBinding) => {
  const { value: directiveValue } = binding
  if (directiveValue?.menuList.length > 0) {
    contextMenuEvent = (e: MouseEvent) => {
      if (typeof directiveValue.disabled === 'function' && directiveValue.disabled(directiveValue.params)) return
      e.preventDefault()
      ContextMenuMenuCtx = CreateCustomMouseMenu({
        el,
        ...directiveValue,
      })
      const { x, y } = e
      ContextMenuMenuCtx.show(x, y)
    }
    el.removeEventListener('contextmenu', contextMenuEvent)
    el.addEventListener('contextmenu', contextMenuEvent)
  } else {
    throw new Error('At least set one menu list!')
  }
}

const unmounted = (el: HTMLElement) => {
  el.removeEventListener('contextmenu', contextMenuEvent)
}

const CustomContextMenuDirective: ObjectDirective = {
  mounted,
  unmounted,
}
export { CustomContextMenuDirective, CreateCustomMouseMenu, ContextMenuMenuCtx }
export default CustomContextMenu
