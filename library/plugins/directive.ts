import type { App, DirectiveBinding } from 'vue'
import { hasPermission } from '@/utils/permission'
import { useCopy } from '@/utils'
export default {
  install(app: App<Element>) {
    /**
     * @description 自定义指令v-permissions
     */
    app.directive('permissions', {
      mounted(el: any, binding: DirectiveBinding) {
        const { value } = binding
        if (value) if (!hasPermission(value)) el.parentNode && el.parentNode.removeChild(el)
      },
    })
    /**
     * @description 自定义指令v-tableLazy,表格懒加载
     */
    app.directive('lazy', {
      mounted(el: HTMLElement, binding: DirectiveBinding) {
        const observe = new IntersectionObserver(
          ([{ isIntersecting }]) => {
            if (isIntersecting) {
              // 停止观察
              observe.unobserve(el)
              binding.value
            }
          },
          {
            threshold: 0,
          }
        )
        // 开启观察
        observe.observe(el)
      },
    })

    /**
     * @description 自定义指令v-copy,复制
     */
    app.directive('copy', (el: HTMLElement, binding: DirectiveBinding) => {
      if (el.querySelector('.directive_copy')) {
        return
      }
      const text = binding.value
      if (!text) return
      const elStyle = getComputedStyle(el)
      if (elStyle.position === 'static') el.style.position = 'relative'
      el.classList.add('directive_copy_body')
      const copyBtn = document.createElement('span')
      copyBtn.className = 'directive_copy'
      copyBtn.textContent = '复制'
      const copyDom = el.appendChild(copyBtn)
      copyDom.onclick = () => {
        useCopy(text)
      }
    })

    /**
     * @description 自定义指令v-fold,折叠展开
     */
    app.directive('fold', (el: HTMLElement, binding: DirectiveBinding) => {
      if (el.querySelector('.directive_fold')) {
        return
      }
      const list = Array.from(el.children)
      let maxHeight = 0
      list?.forEach((key: any, index: number) => {
        const height = el.children[index].getBoundingClientRect().height
        maxHeight = maxHeight >= height ? maxHeight : height
      })
      const elStyle = getComputedStyle(el)
      if (elStyle.position === 'static') el.style.position = 'relative'
      el.classList.add('my-fold-class')
      if (maxHeight >= 460) {
        const foldBtn = document.createElement('div')
        foldBtn.className = 'directive_fold'
        const foldDom = el.appendChild(foldBtn)
        foldDom.onclick = () => {
          el.classList.remove('my-fold-class')
          el.style.maxHeight = 'initial'
          el.removeChild(foldBtn)
        }
      }
    })

    /**
     * @description 自定义指令v-fold,折叠展开
     */
    app.directive('fold2', {
      created(el: HTMLElement) {
        const elStyle = getComputedStyle(el)
        if (elStyle.position === 'static') el.style.position = 'relative'
        el.classList.add('my-fold-class')
        setTimeout(() => {
          const list = Array.from(el.children)
          let maxHeight = 0
          list?.forEach((key: any, index: number) => {
            const height = el.children[index].getBoundingClientRect().height
            maxHeight += height
          })
          if (maxHeight >= 460) {
            const foldBtn = document.createElement('div')
            foldBtn.className = 'directive_fold'
            const foldDom = el.appendChild(foldBtn)
            foldDom.onclick = () => {
              el.classList.add('my-expand-class')
              el.removeChild(foldBtn)
            }
          }
        })
      },
    })

    /**
     * @description 自定义指令v-fold,折叠展开
     */
    app.directive('dialogBackTop', (el: HTMLElement) => {
      if (document.querySelector('.dialogBackTop')) return
      nextTick(() => {
        // 指令的祖父节点
        const rootParent = el.offsetParent
        const rootParentRect = rootParent?.getBoundingClientRect()
        if (rootParent) {
          // dialog滚动节点
          const dialogNode = ref(rootParent.parentNode as HTMLDivElement)
          const { y } = useScroll(dialogNode, { behavior: 'smooth' })
          let _backNode = null as HTMLDivElement | null
          watch(y, () => {
            if (y.value > 300) {
              if (_backNode) return (_backNode!.style.display = 'block')
              _backNode = document.createElement('div')
              _backNode.classList.add('dialogBackTop')
              _backNode.style.right = rootParentRect?.left! - 86 + 'px'
              _backNode.onclick = () => {
                y.value = 0
              }
              dialogNode.value.appendChild(_backNode)
            } else {
              if (_backNode) _backNode!.style.display = 'none'
            }
          })
        }
      })
    })
  },
}
