import EcsHookDialog from './dialog.vue'
import { createVNode, FunctionalComponent, VNode, render as vueRender } from 'vue'
import AppContext from '@/main'
export interface DialogServiceType {
  /**
   * @param 是否展示页脚  默认为 true,为 false 时需要自组件广播ecs-dialog-close事件进行关闭
   */
  showFooter?: boolean
  /**
   * @param 展示宽度 默认为 520
   */
  width?: string | number
  /**
   * @param 标题  默认为  标题
   */
  title: string
  /**
   * @param 展示内容
   */
  content: FunctionalComponent | VNode
  /**
   * @param 点击确认的回调，返回Promise会延迟关闭，请求成功后自动关闭，期间不允许手动关闭
   */
  confirm?: () => Promise<unknown> | void
  /**
   * @param 点击取消的回调，
   */
  cancel?: () => void
}
export const useEcsDialogService = (options: DialogServiceType) => {
  const container = document.createElement('div')
  container.classList.add('ecs-hook-dialog')
  document.body.appendChild(container)

  function render(props: DialogServiceType) {
    const { content, confirm, cancel, width = 520, showFooter = true, ...otherProps } = props
    const funcHandles = {
      onCancel: () => {
        cancel?.()
        handleRemove()
      },
      onConfirm: () => {
        vm.component!.props.loading = true
        const result = confirm && confirm()
        if (result instanceof Promise) {
          result
            .then((resolvedValue) => {
              Object.assign(componentInstance.component!.props, {
                loading: false,
              })
              handleRemove()
            })
            .catch((error) => {
              Object.assign(componentInstance.component!.props, {
                loading: false,
              })
            })
        } else {
          vm.component!.props.loading = false
          handleRemove()
        }
      },
    }

    const vm = createVNode(
      EcsHookDialog,
      {
        ...otherProps,
        width,
        showFooter,
        loading: false,
        modelValue: true,
        beforeClose(done: () => void) {
          const loading = componentInstance.component!.props.loading
          if (loading) return
          done()
        },
        ...funcHandles,
      },
      {
        default: () =>
          h(
            content,
            showFooter ? null : { onEcsDialogClose: funcHandles.onCancel, onEcsDialogConfirm: funcHandles.onConfirm }
          ),
      }
    )
    if (AppContext._context) {
      vm.appContext = AppContext._context
    }
    vueRender(vm, container) // 渲染组件
    return vm
  }
  const handleRemove = () => {
    if (!componentInstance.component) return
    if (componentInstance.component.props.loading || !componentInstance.component.props.modelValue) return
    Object.assign(componentInstance.component.props, {
      loading: false,
      modelValue: false,
    })
    setTimeout(() => {
      vueRender(null, container)
      document.body.removeChild(container)
    }, 800)
  }
  const componentInstance = render(options)
  return {
    destroy: handleRemove,
  }
}
export default useEcsDialogService
