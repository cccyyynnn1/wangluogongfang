<script lang="ts">
  export default {
    name: 'AttackHighlight',
  }
</script>

<script setup lang="ts">
  import { EditorView, ViewUpdate, ViewPlugin, Decoration, type DecorationSet } from '@codemirror/view'
  import { EditorState, StateEffect } from '@codemirror/state'
  import { minimalSetup } from 'codemirror'
  import { PacketsMagic } from '@/utils/magic'
  import iconv from 'iconv-lite'
  import { CustomContextMenuDirective } from '@/components/custom-context-menu'
  import { useEcsDialogService } from '@/components/ecs-dialog'
  import AttackHighlightRule from '@/ecs/config/alert-configuration/attack-highlight-rule.vue'
  import { useUserStore } from '@/store/modules/user'
  import { PayloadDivider } from '~/src/hooks/useHttpTelegramMatch'
  const router = useRouter()
  type FormatFunType = keyof typeof formatFun
  type ConvertsType = FormatFunType | 'default'
  const vMouseMenu = CustomContextMenuDirective
  const userStore = useUserStore()
  const { attackCharacterization } = storeToRefs(userStore)
  const formatFun = {
    magic: PacketsMagic,
    gb2312: handleConvert,
    gbk: handleConvert,
  }
  const props = defineProps<{
    // searchTerm?: string[]
    attackTerm?: string[]
    attackHttpMsg: string
    attackType?: 'request' | 'response' | 'all'
  }>()
  type ScopeType = 'requestHeader' | 'requestPayload' | 'responseHeader' | 'responsePayload' | 'all'
  const currentAttack = computed(() => {
    if (!props.attackType) return null
    const { highLightConfig, highLightWhite } = attackCharacterization.value
    return {
      highLightConfig: highLightConfig
        .filter((item) => item.scope.includes(props.attackType!) || item.scope === 'all')
        .map(({ scope, content }) => ({ content, scope })),
      highLightWhite: highLightWhite
        .filter((item) => item.scope.includes(props.attackType!) || item.scope === 'all')
        .map(({ scope, content }) => ({ content, scope })),
    }
  })

  const attackWhiteRange = computed(() => {
    if (!currentAttack.value?.highLightWhite.length) return null
    const attackHttpMsg = props.attackHttpMsg.replaceAll('\r', '')
    return findKeywordPositions(attackHttpMsg, currentAttack.value.highLightWhite)
  })

  const attackDomRef = ref<HTMLDivElement>()
  const attackVal = ref('')
  const attackTerms = ref()
  const convertsType = ref<ConvertsType>('default')
  const oldConvertsType = ref<ConvertsType>('default')
  const highlightView = ref<EditorView>()
  const baseTheme = EditorView.baseTheme({
    '.cm-attackMatch': { color: '#fff', backgroundColor: 'red' },
    '.cm-searchMatch': { color: '#fff', backgroundColor: '#FF7301' },
  })
  const attackDeco = Decoration.mark({ class: 'cm-attackMatch' })
  const searchDeco = Decoration.mark({ class: 'cm-searchMatch' })
  // 右键菜单
  const options: any = {
    disabled: (params: any) => {
      return !params.keyword
    },
    params: { keyword: '' },
    menuWidth: 120,
    menuList: [
      {
        label: '添加自定义告警规则',
        tips: 'AddAlarmRules',
        callback: (menuItem: any) => {
          const resolveRouter = router.resolve({
            path: '/config/alert-configuration/custom-rules',
            query: {
              addRule: encodeURIComponent(
                JSON.stringify({
                  ruleVal: options.params.keyword,
                  ruleDirection: options.params.direction,
                })
              ),
            },
          })
          window.open(resolveRouter.href, '_blank')
        },
      },
      {
        label: '添加到高亮特征库',
        tips: 'AddHighlight',
        callback: (menuItem: any) => {
          handleAppendAttack(options.params.keyword, 'highLightConfig')
        },
      },
      {
        label: '添加到特征库白名单',
        tips: 'AddWhite',
        callback: (menuItem: any) => {
          handleAppendAttack(options.params.keyword, 'highLightWhite')
        },
      },
    ],
  }

  // 查找关键字所在的位置
  function findKeywordPositions(
    longString: string,
    keywords: {
      content: string
      scope: ScopeType
    }[]
  ) {
    const _keywords = keywords.map((i) => i.content.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'))
    // 创建正则表达式
    const regex = new RegExp(`${_keywords.join('|')}`, 'g')
    const positions = []
    let match: RegExpExecArray | null
    while ((match = regex.exec(longString)) !== null) {
      // 获取匹配关键字的位置
      const scopeIndex = keywords.findIndex((i) => i.content === match![0])
      positions.push({
        keyword: match[0],
        scope: keywords[scopeIndex].scope,
        start: match.index,
        end: match.index + match[0].length,
      })
    }
    return positions
  }

  // 格式化攻击特征
  const formatTerm = (attackTerm: string[]): string[] => {
    const funKey = convertsType.value as FormatFunType
    //  攻击特征进行格式转码
    const _baseTerm = convertsType.value === 'default' ? attackTerm : attackTerm.map((term) => formatFun[funKey](term))
    // 魔法棒工具下高亮精准匹配
    const magicTerm = _baseTerm.map((itemTerm) => {
      const includesTermIndex = attackVal.value.indexOf(itemTerm)
      if (includesTermIndex > -1) return itemTerm
      const firstUriCodeIndex = itemTerm.indexOf('%')
      const uriCode =
        itemTerm.length - firstUriCodeIndex >= 2
          ? itemTerm.slice(0, firstUriCodeIndex)
          : itemTerm.slice(firstUriCodeIndex)
      if (attackTerm.includes(uriCode)) return uriCode
      const lastUriCodeIndex = itemTerm.lastIndexOf('%')
      return uriCode.slice(0, lastUriCodeIndex + 1)
      /*
          ===============保留==========
        const includesTermIndex = attackVal.value.indexOf(itemTerm)
        if (includesTermIndex > -1) return itemTerm
        const startWithIndex = itemTerm.slice(0, 2).indexOf('%')
        const endWithIndex = itemTerm.slice(-2).indexOf('%')
        return itemTerm.slice(Math.max(startWithIndex, 0), Math.max(endWithIndex - 2, -2))
          ===============保留==========
        */
    })
    attackTerms.value = _baseTerm
    return attackTerms.value
  }
  // 格式化攻击日志
  const formatAttackValues = (valStr: string) => {
    const funKey = convertsType.value as FormatFunType
    const formatVal = convertsType.value === 'default' ? props.attackHttpMsg : formatFun[funKey](valStr)
    attackVal.value = formatVal
    highlightView.value?.dispatch({
      changes: {
        from: 0,
        to: highlightView.value.state.doc.length,
        insert: formatVal,
      },
    })
  }
  /**
   * 进行高亮显示.
   * @param attackTerm - 要高亮的攻击特征
   * @param searchTerm - 要高亮的检索词
   */
  const searchHighlightPlugin = (
    attackTerm?: string[],
    searchTerm?: {
      content: string
      scope: ScopeType
    }[]
  ) =>
    ViewPlugin.fromClass(
      class {
        decorations: DecorationSet
        constructor(view: EditorView) {
          this.decorations = this.createDeco(view, [attackTerm || [], searchTerm || []])
        }
        generateRanges(
          text: string,
          keywords: [
            string[],
            {
              content: string
              scope: ScopeType
            }[]
          ]
        ) {
          if (!text) return []
          const [attack, search] = keywords
          const ret = [] as Array<{ from: number; to: number; type: 'attack' | 'search' }>
          const _attack = formatTerm(attack)
          if (search && search.length) {
            const searchRange = findKeywordPositions(text, search)
            searchRange.forEach((_search) => {
              const { keyword, scope, start, end } = _search
              const findWhite = attackWhiteRange.value?.filter(
                (item) => [scope, 'all'].includes(item.scope) && item.keyword.includes(keyword)
              )
              if (findWhite?.length) {
                //findWhite是作用域为ALL或者是当前attackType类型范围，所以只要有一个命中就不能添加
                const hasIncludes = findWhite.find((white) => {
                  const { end: _whiteEnd, start: _whiteStart, keyword: _whiteKeyword } = white
                  return start >= _whiteStart && end <= _whiteEnd
                })
                if (!hasIncludes) ret.push({ from: start, to: end, type: 'search' })
              } else ret.push({ from: start, to: end, type: 'search' })
            })
          }
          if (_attack && _attack.length) {
            _attack.forEach((keyword) => {
              let index = 0
              while ((index = text.indexOf(keyword, index)) !== -1) {
                ret.push({ from: index, to: index + keyword.length, type: 'attack' })
                index += index + keyword.length
              }
            })
          }
          ret.sort((a, b) => {
            if (a.from === b.from) {
              return a.to - b.to
            }
            return a.from - b.from
          })
          return ret
        }
        update(update: ViewUpdate) {
          if (update.docChanged || update.viewportChanged) {
            this.decorations = this.createDeco(update.view, [attackTerm || [], searchTerm || []])
          }
        }
        createDeco(
          view: EditorView,
          keywords: [
            string[],
            {
              content: string
              scope: ScopeType
            }[]
          ] = [[], []]
        ) {
          const ranges = this.generateRanges(view.state.doc.toString(), keywords)
          return ranges.length
            ? Decoration.set(
                ranges.map((r) => {
                  return r.type === 'attack' ? attackDeco.range(r.from, r.to) : searchDeco.range(r.from, r.to)
                })
              )
            : Decoration.none
        }
      },
      {
        decorations: (v) => v.decorations,
      }
    )
  /**
   * 进行高亮显示.
   * @param searchTerm - 要高亮的攻击特征
   * @param searchTerm - 要高亮的检索词
   */
  const handleKeywordSearchInEditor = (
    attackTerm?: string[],
    searchTerm?: {
      content: string
      scope: ScopeType
    }[]
  ) => {
    if (!props.attackType) return
    highlightView.value?.dispatch({
      effects: [StateEffect.appendConfig.of([searchHighlightPlugin(attackTerm, searchTerm)])],
    })
  }
  const initAttackView = () => {
    // 处理右键菜单显示和位置设置
    const handleContextMenu = (view: any, event: any) => {
      event.preventDefault()
      const newText = window.getSelection()?.toString()
      const selection = view.state.selection
      const position = { from: selection.main.from, to: selection.main.to }
      const PayloadDividerIndex = view.state.doc.toString().indexOf(PayloadDivider)
      // const [DividerRangeStart, DividerRangeEnd] = [PayloadDividerIndex, PayloadDividerIndex + PayloadDivider.length]
      // 判断选词是否是在body还是header
      const isHeader = position.to > PayloadDividerIndex
      options.params = {
        keyword: newText,
        direction: props.attackType === 'request' ? '客户端到服务端' : '服务端到客户端',
      }
    }
    formatAttackValues(props.attackHttpMsg)
    highlightView.value = new EditorView({
      state: EditorState.create({
        doc: attackVal.value,
        extensions: [
          minimalSetup,
          baseTheme,
          EditorState.readOnly.of(true),
          EditorView.domEventHandlers({
            contextmenu: (event) => handleContextMenu(highlightView.value, event),
          }),
        ],
      }),
      parent: attackDomRef.value,
    })
    handleKeywordSearchInEditor(props.attackTerm, currentAttack.value?.highLightConfig || [])
  }
  // gbk转换
  function handleConvert(convertVal: string) {
    const currType = ['default', 'magic'].includes(oldConvertsType.value) ? 'utf-8' : oldConvertsType.value
    const targetType = ['default', 'magic'].includes(convertsType.value) ? 'utf-8' : convertsType.value
    const gbkBuffer = iconv.encode(convertVal, currType)
    const convertString = iconv.decode(gbkBuffer, targetType)
    return convertString
  }
  //  添加攻击特征到xxx
  const handleAppendAttack = (content: string, type: 'highLightConfig' | 'highLightWhite') => {
    const { destroy } = useEcsDialogService({
      title: '添加',
      showFooter: false,
      content: () => h(AttackHighlightRule, { attackHighlightData: { content, type } }),
      cancel() {
        destroy()
      },
    })
  }

  defineExpose({
    // highlight: handleKeywordSearchInEditor,
    setConvertsType: (type: ConvertsType) => (convertsType.value = convertsType.value === type ? 'default' : type),
    getAttackTerm: () => attackTerms.value,
    getAttackValues: () => attackVal.value,
  })
  onMounted(() => {
    initAttackView()
  })
  watch(
    () => convertsType.value,
    (newType, oldType) => {
      oldConvertsType.value = oldType
      attackTerms.value = undefined
      formatAttackValues(attackVal.value)
    }
  )
  watch(
    () => [props.attackTerm, attackCharacterization.value],
    () => {
      attackTerms.value = undefined
      handleKeywordSearchInEditor(props.attackTerm, currentAttack.value?.highLightConfig || [])
    }
  )
  watch(
    () => props.attackHttpMsg,
    () => {
      formatAttackValues(props.attackHttpMsg)
    }
  )
</script>
<template>
  <div ref="attackDomRef" v-mouse-menu="options" class="attack-highlight cm-readonly" />
</template>

<style scoped lang="scss">
  .attack-highlight {
    width: 100%;
    height: 100%;
    overflow-y: auto;
    border: 1px solid rgb(230, 231, 240);
    background-color: #f8f7ff;
    border-radius: 8px;
    &.cm-readonly {
      :deep(.cm-cursor) {
        visibility: hidden;
      }
    }
    :deep() {
      .cm-editor {
        padding: 10px 10px 6px;
        .cm-selectionLayer {
          display: none;
        }
        ::selection {
          background-color: var(--el-color-primary) !important;
        }
      }
    }
  }
</style>
