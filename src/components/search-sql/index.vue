<script lang="ts">
  export default {
    name: 'SearchSql',
  }
  type FormSqlItem = {
    field?: string
    operator: string
    value?: string
    type: 'field' | 'operator'
  }
</script>

<script setup lang="ts">
  import { EditorView, Decoration } from '@codemirror/view'
  import { EditorState } from '@codemirror/state'
  import { Codemirror } from 'vue-codemirror'
  import { sql } from '@codemirror/lang-sql'
  import { autocompletion, snippetCompletion, CompletionContext, startCompletion } from '@codemirror/autocomplete'
  import { TableColumnItemType, TableFieldType } from '/#/store'
  import { Search } from '@element-plus/icons-vue'
  import VabDialog from '@/plugins/VabDialog/index.vue'
  import { cloneDeep, debounce } from 'lodash'
  import { useUserStore } from '@/store/modules/user'
  import { getSearchHistories } from '@/api-ecs/public'
  import { onClickOutside } from '@vueuse/core'
  const userStore = useUserStore()
  // IPV4 包含网段校验
  const ipv4Regex =
    /^(?!0)(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.(?:[0-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])(\/([1-9]|[1-2]\d|3[0-2]))?$/
  // IPV6
  const ipv6Regex =
    /^(([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$/
  // 端口
  const portRegex = /^(?:[1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$/
  const $baseMessage: any = inject('$baseMessage')
  const SqlParser = require('js-sql-parser')
  // 关系运算符
  const relationalOperator = [
    { label: '=', detail: '等于', template: '=' },
    { label: '!=', detail: '不等于', template: '!=' },
    { label: 'like', detail: '字符串包含', template: 'like' },
    { label: 'not_like', detail: '字符串不包含', template: 'not_like' },
    { label: 'exists', detail: '存在', template: 'exists' },
    { label: 'not_exists', detail: '不存在', template: 'no_exists' },
    { label: '>', detail: '大于', template: '>' },
    { label: '>=', detail: '大于等于', template: '>=' },
    { label: '<', detail: '小于', template: '<' },
    { label: '<=', detail: '小于等于', template: '<=' },
  ]
  const relationalKeys = ['=', '!=', 'like', '>', '<', '>=', '<=']
  // 逻辑运算
  const logicOperator = [
    { label: 'and', detail: '并且', template: 'and' },
    { label: 'or', detail: '或', template: 'or' },
    { label: 'not', detail: '非', template: 'not' },
    { label: '(', detail: '左括号', template: '(' },
    { label: ')', detail: '右括号', template: ')' },
  ]
  const props = withDefaults(
    defineProps<{
      modelValue: string
      placeholder?: string
      /**
       * @params 站点页面不出indexType
       */
      indexType?: number
      workspaceId?: number
      /**
       * @params siteData 只在站点出现，出现时不需要indexType
       */
      siteData?: {
        siteId?: number
        apiId?: number
      }
    }>(),
    {
      modelValue: '',
      placeholder: '请输入检索语句',
      indexType: 1,
      workspaceId: undefined,
      siteData: () => ({
        siteId: undefined,
        apiId: undefined,
      }),
    }
  )
  const emits = defineEmits<{
    (e: 'on-change', sqlVal: string): void
    (e: 'on-search', sqlVal: string): void
  }>()

  /**
   * @description 是否是推荐的，组件初始化时设为True，当选择检索历史，手动输入SQl，点击检索后设为False。
   */
  const isRecommend = ref(true)
  const searchSqlRef = ref<HTMLDivElement>()
  const editView = shallowRef<EditorView>()
  const editState = shallowRef<EditorState>()

  const sqlFieldRef = ref<HTMLDivElement>()
  const sqlOperatordRef = ref<HTMLDivElement>()

  // sql语句值
  const sqlVal = ref(props.modelValue)
  // 是否聚焦
  const isFocus = ref(false)
  // SQL语句提示
  const sqlHintPaths = ref<{ label: string; detail: string; template: string }[]>([])
  // SQL语句提示
  let allSqlHintPathsField: string[] = []
  // SQL语句是否合法
  const isValidate = ref(true)
  //表单式输入sql
  const formSqlList = ref<FormSqlItem[]>([])
  // formSQL编辑信息
  const formSqlEditData = ref<FormSqlItem>()
  const formSqlEditDialogVisible = ref(false)
  const formSqlEditDirection = ref<'left' | 'self' | 'right'>('self')
  const formSqlEditIndex = ref(-1)
  const formSqlSearchFieldStr = ref('')
  const formSqlFields = computed(() =>
    formSqlSearchFieldStr.value === ''
      ? sqlHintPaths.value
      : sqlHintPaths.value.filter((item) =>
          item.label.toLowerCase().includes(formSqlSearchFieldStr.value.toLowerCase())
        )
  )
  //  查询sql记录
  const searchSqlHistories = ref<{ id: number; searchSql: string }[]>([])
  const formatColum = (
    colums: TableColumnItemType[] | { fieldNameCn: string; fieldNameEn: string; fieldType?: TableFieldType }[]
  ) => {
    allSqlHintPathsField = []
    const arr = []
    for (const { fieldNameCn, fieldNameEn, fieldType = undefined } of colums) {
      if (fieldType === undefined || fieldType !== 'no_serach') {
        allSqlHintPathsField.push(fieldNameCn)
        arr.push({
          label: fieldNameCn,
          detail: fieldNameEn,
          template: fieldNameCn,
        })
      }
    }
    return arr
  }
  // 控制sql配置组件显示
  const sqlConfigVisible = ref(false)
  const sqlConfigTriggerRef = ref()

  const cmOption = reactive({
    tabSize: 1,
    autoDestroy: true,
    placeholder: props.placeholder,
    extensions: [
      sql(),
      autocompletion({
        activateOnTyping: true,
        override: [],
      }),
    ],
  })
  onClickOutside(searchSqlRef, (event) => (isFocus.value = false))
  const handleTriggerSql = (sqlItem: FormSqlItem, index: number, e?: MouseEvent) => {
    sqlConfigTriggerRef.value?.classList?.remove('active')
    //这是为了点击formSql时，填充值的时候不携带双引号，在确认编辑后会自动补全双引号。
    const newItem = cloneDeep(sqlItem)
    newItem.value = newItem.value?.slice(1, -1)
    formSqlEditData.value = newItem
    formSqlEditIndex.value = index

    if (!e) return
    sqlConfigTriggerRef.value = e.target
    sqlConfigTriggerRef.value?.classList?.add('active')
    sqlConfigVisible.value = true
  }
  const handleReady = ({ view, state }: { state: EditorState; view: EditorView }) => {
    editView.value = view
    editState.value = state
  }
  const handleFormSqlEdit = (direction: 'left' | 'self' | 'right') => {
    if (direction !== 'self') {
      const [left, right]: number[] = [formSqlEditIndex.value - 1, formSqlEditIndex.value]
      const siblingItem =
        direction === 'left' ? formSqlList.value[left] || formSqlList.value[right] : formSqlList.value[right]
      if (siblingItem.type === 'field') {
        formSqlEditData.value = { type: 'operator', operator: 'and' }
      } else {
        formSqlEditData.value = { type: 'field', operator: '=', value: '', field: '源IP' }
      }
    }
    sqlConfigTriggerRef.value?.classList.add('active')
    formSqlEditDirection.value = direction
    formSqlEditDialogVisible.value = true
  }
  const handleFormSqlDelete = () => {
    if (formSqlEditIndex.value < 0) return
    formSqlList.value.splice(formSqlEditIndex.value, 1)
    sqlConfigVisible.value = false
    formSqlEditIndex.value = -1
    sqlVal.value = handleSyncField2Sql()
  }
  const handleRemoveEditTag = () => {
    sqlConfigTriggerRef.value?.classList?.remove('active')
    sqlConfigVisible.value = false
  }
  const handleGetHistories = async () => {
    const { workspaceId, indexType, siteData } = props
    const { data } = await getSearchHistories({ ...siteData, workspaceId, indexType, pageNum: 1, pagesize: 10 })
    searchSqlHistories.value = data.records || []
  }
  // 给每个括号加空格（双引号里的除外）
  const handleFormatterSql = (sqlStr: string) => {
    // 中英文符号替换
    const symbolMapping: { [key: string]: string } = {
      '，': ',',
      '。': '.',
      '！': '!',
      '？': '?',
      '：': ':',
      '；': ';',
      '（': '(',
      '）': ')',
      '“': '"',
      '”': '"',
      '‘': "'",
      '’': "'",
      // 可以根据需要添加更多映射
    }
    const escapedSymbols = Object.keys(symbolMapping)
      .map((symbol) => `\\${symbol}`)
      .join('')
    const replacePattern = new RegExp(`[${escapedSymbols}]`, 'g')
    const placeholders: string[] = []
    // 对字符串先把双引号里面的值进行替换然后对字符串里的中英文符号进行替换
    const processedSQL = sqlStr
      .replace(/"[^"]*"/g, (match) => {
        placeholders.push(match)
        return `"PLACEHOLDER_${placeholders.length - 1}"`
      })
      .replace(replacePattern, (match) => symbolMapping[match] || match)
    // 对双引号添加空格
    const modifiedSQL = processedSQL
      .replace(/[\(\)]/g, (math) => ` ${math} `)
      .replace(/\s+/g, ' ')
      .trim()
    const finalSQL = modifiedSQL.replace(/"PLACEHOLDER_(\d+)"/g, (match, index) => placeholders[index])
    return finalSQL
  }
  /**
   *
   * @description 校验FormSqlItem是否符合规范
   * @param validatOperator 是否校验operator，默认true，false为不校验返回false
   */
  const handleValidateFormSql = (formSqlItem: FormSqlItem) => {
    if (formSqlItem.type === 'operator') return false
    return formSqlItem.operator.endsWith('exists') ? true : !!formSqlItem.value
  }

  function isWithinQuotes(s: string, editStart: number, editEnd: number) {
    let quoteLevel = 0 // 用来跟踪双引号的层级
    let currentQuoteStart = -1 // 当前最外层双引号的开始位置
    let escape = false // 用来标记转义符
    // 遍历字符串，查找最外层双引号
    for (let i = 0; i < s.length; i++) {
      const char = s[i]
      // 如果遇到转义字符，跳过下一个字符
      if (escape) {
        escape = false
        continue
      }
      // 检查转义字符
      if (char === '\\') {
        escape = true
        continue
      }
      // 处理双引号
      if (char === '"') {
        if (quoteLevel === 0) {
          currentQuoteStart = i // 进入最外层双引号
        } else if (quoteLevel === 1) {
          // 退出最外层双引号
          if (editStart >= currentQuoteStart && editEnd <= i) {
            return true // 如果编辑范围在这对最外层双引号内，返回 true
          }
        }
        quoteLevel = 1 - quoteLevel // 切换双引号的层级
      }
    }
    return false // 如果没有在最外层双引号内
  }

  // 针对sql语句出现Exists的处理
  const handlingExistsStr = (sql: string) => {
    // 用一个占位符替换所有的双引号内容
    const placeholders: string[] = []
    const processedSQL = sql.replace(/"[^"]*"/g, (match) => {
      placeholders.push(match)
      return `"PLACEHOLDER_${placeholders.length - 1}"`
    })
    // 进行正则替换(匹配 exists或者not_exists 前面的数据，并替换)
    const modifiedSQL = processedSQL
      .replace(/(\S+)\s+(EXISTS|NOT_EXISTS)/gi, (match, p1, p2) => {
        return `${p2.replace('_', ' ')} (SELECT * FROM ${p1})`
      })
      .replace(/not_like/gi, 'not like')
    // 恢复双引号内容
    const finalSQL = modifiedSQL.replace(/"PLACEHOLDER_(\d+)"/g, (match, index) => placeholders[index])
    return finalSQL
  }
  // sql提示
  const hintPlugin = () => {
    const hintPaths = sqlHintPaths.value
    return (context: CompletionContext) => {
      // 匹配当前输入前面的所有非空字符
      const word = context.matchBefore(/\S*/)
      const docStr = context.state.doc.toString()
      /*
            如果为空格 根据前面字段显示不同提示
            1:以冒号结尾 显示连接符
            2:以连接符结尾 显示所有字段
            3:以空格结尾 显示逻辑运算符
            4:以逻辑运算符结尾 需要给Sql新增冒号
            5:当exists或not_exists结尾时提示连接符
        */
      if (!word) return null
      if (word.from == word.to) {
        const text = docStr.slice(0, word.to - 1)
        const hasLink = ['and', 'or', 'not'].some((i) => text.endsWith(i))
        const hasLogic = relationalKeys.some((i) => text.endsWith(i))
        const hasExists = ['exists', 'not_exists', ' ', ')', '"'].some((i) => text.toLowerCase().endsWith(i))
        if (!text || hasLink) {
          return {
            from: word.from,
            to: word.to,
            filter: false,
            options: hintPaths.map(({ detail, label }) =>
              snippetCompletion(`${label} \${0}`, {
                label,
                detail,
              })
            ),
          }
        }
        if (hasExists) {
          return {
            from: word.from,
            to: word.to,
            filter: false,
            options: logicOperator.map(({ detail, label }) =>
              snippetCompletion(`${label} \${0}`, {
                label,
                detail,
              })
            ),
          }
        }
        const _relationalOperator = text.endsWith('方向')
          ? [
              { label: '=', detail: '等于', template: '=' },
              { label: '!=', detail: '不等于', template: '!=' },
            ]
          : relationalOperator
        return {
          from: word.from,
          to: word.to,
          filter: false,
          options: _relationalOperator.map(({ label, detail }) => {
            const isExists = ['exists', 'not_exists'].includes(label)
            return snippetCompletion(isExists ? `${label}\${0}` : `${label} "\${0}"`, {
              label,
              detail,
            })
          }),
        }
      }
      const inQuotes = isWithinQuotes(docStr, word.from, word.to)
      const newInput = word.text
      if (inQuotes) return null
      //用来匹配非特殊符号字符可进行语法提示
      const regex = /[A-Za-z0-9\u4e00-\u9fa5\s]+/g
      const matches = newInput.match(regex)?.[0]
      if (!matches) return null
      return {
        from: word.to - matches.length,
        to: word.to,
        filter: false,
        options: hintPaths
          .filter((item) => item.label.toLowerCase().includes(matches.toLowerCase()))
          .map(({ detail, label }) =>
            snippetCompletion(`${label} \${0}`, {
              label,
              detail,
            })
          ),
      }
    }
  }
  // 根据sql转换为form数据
  function formatterSQl2Frorm(inputSql: string, showMsg = false) {
    // 定义逻辑操作符和比较操作符
    const logicalOperators = ['and', 'or', 'not', '(', ')']
    const comparisonOperators = ['=', '!=', '>', '>=', '<', '<=', 'like', 'not_like', 'exists', 'not_exists']
    // 使用正则表达式将输入字符串拆分为令牌，保留引号中的内容
    const tokens = handleFormatterSql(inputSql).match(/"((?:[^"\\]|\\.)*)"|\S+/g)
    if (!tokens) return []
    const result: FormSqlItem[] = []
    let i = 0
    while (i < tokens.length) {
      let token = tokens[i]
      if (logicalOperators.includes(token.toLowerCase())) {
        // 如果是逻辑操作符
        result.push({
          type: 'operator',
          operator: token,
        })
        i++
      } else {
        // 处理元素
        const field = token
        let operator = ''
        let value = ''
        // 尝试匹配多词比较操作符（最长匹配原则）
        let maxOperatorLength = 2 // 最大操作符长度
        let operatorFound = false
        for (let len = maxOperatorLength; len >= 1; len--) {
          const opTokens = tokens.slice(i + 1, i + 1 + len)
          if (!opTokens.length) throw new Error(`${field}字段没有连接符`)
          const opString = opTokens.join(' ')
          if (comparisonOperators.includes(opString.toLowerCase())) {
            operator = opString
            operatorFound = true
            i += 1 + len // 移动索引，跳过字段和操作符
            // 检查操作符是否需要值
            if (operator.toLowerCase() === 'exists' || operator.toLowerCase() === 'not_exists') {
              value = ''
            } else {
              value = tokens[i]
              if (!value) {
                if (showMsg) $baseMessage(`操作符 '${operator}' 缺少值`, 'error', 'vab-hey-message-error')
                throw new Error(`操作符 '${operator}' 缺少值`)
              }
              i++ // 移动索引，跳过值
            }
            break
          } else {
            operator = tokens[i]
          }
        }
        if (!operatorFound) {
          const err_operator = tokens[i + 1]
          const text = `未知的比较操作符在位置 ${inputSql.indexOf(err_operator)} ${err_operator}`
          if (showMsg) $baseMessage(text, 'error', 'vab-hey-message-error')
          throw new Error(text)
        }
        result.push({
          type: 'field',
          value,
          field,
          operator: operator,
        })
      }
    }
    return result
  }
  // 校验SQL语句是否合法
  const handleValidateSql = (sql: string, showMsg = false) => {
    if (sql === '' || !sqlHintPaths.value.length) {
      formSqlList.value = [
        {
          field: '源IP',
          operator: '=',
          value: '',
          type: 'field',
        },
        {
          operator: 'and',
          type: 'operator',
        },
        {
          field: '源端口',
          operator: '=',
          value: '',
          type: 'field',
        },
        {
          operator: 'and',
          type: 'operator',
        },
        {
          field: '目的IP',
          operator: '=',
          value: '',
          type: 'field',
        },
        {
          operator: 'and',
          type: 'operator',
        },
        {
          field: '目的端口',
          operator: '=',
          value: '',
          type: 'field',
        },
      ]
      isRecommend.value = true
      return true
    }
    try {
      // 当输入输入符合端口或ipv4、ipv6
      if (ipv4Regex.test(sql) || ipv6Regex.test(sql) || portRegex.test(sql)) return true
      // 当推荐状态的时候不要同步
      const _formSqlList = formatterSQl2Frorm(sql, showMsg)
      if (!isRecommend.value) formSqlList.value = _formSqlList
      const SQL = `select * from user where ${handlingExistsStr(sql)}`
      const { value } = SqlParser.parse(SQL)
      const syntaxValidat = !['Number', 'Identifier', 'Prefix', 'Boolean', 'Null'].includes(value.where?.type)
      //语法校验不通过直接return
      if (!syntaxValidat) {
        if (showMsg) $baseMessage('请检查SQL是否合规', 'error', 'vab-hey-message-error')
        return false
      }
      const hanldeValidateSqlField = (fieldStr: string) => {
        const currFieldIndex = allSqlHintPathsField.findIndex((i) => i.toLowerCase() === fieldStr)
        if (currFieldIndex < 0 || fieldStr !== allSqlHintPathsField[currFieldIndex].toLowerCase()) {
          if (showMsg) $baseMessage(`${fieldStr}字段不存在`, 'error', 'vab-hey-message-error')
          throw new Error(`${fieldStr}字段不存在`)
        }
      }
      for (const formItem of _formSqlList) {
        if (formItem.type === 'operator') continue
        const myField = formItem.field!.toLowerCase()
        hanldeValidateSqlField(myField)
        if (!['源ip', '目的ip', '源端口', '目的端口'].includes(myField)) continue
        if (formItem.operator.toLowerCase().endsWith('exists')) continue
        const val = formItem.value || ''
        if (
          ['源ip', '目的ip'].includes(myField) &&
          !['like', 'not_like'].includes(formItem.operator.toLowerCase()) &&
          val
        ) {
          if (ipv4Regex.test(val.slice(1, -1)) || ipv6Regex.test(val.slice(1, -1))) continue
          if (showMsg) $baseMessage('请检查IP是否合规', 'error', 'vab-hey-message-error')
          return false
        }
        if (
          ['源端口', '目的端口'].includes(myField) &&
          !['like', 'not_like'].includes(formItem.operator.toLowerCase()) &&
          val
        ) {
          if (portRegex.test(val.slice(1, -1))) continue
          if (showMsg) $baseMessage('请检查端口是否合规', 'error', 'vab-hey-message-error')
          return false
        }
      }
      return true
    } catch (error) {
      if (error instanceof Error) {
        const ErrorMsg = error.message.split('\n')
        console.error(error)
      }
      return false
    }
  }
  //FormSQL数据同步到SQL输入框
  const handleSyncField2Sql = () => {
    let newSql: string
    if (isRecommend.value) {
      const validateFormSqlIndex = formSqlList.value
        .map((item, index) => (handleValidateFormSql(item) ? index : -1))
        .filter((i) => i > -1)
      if (validateFormSqlIndex.length === 0) return (sqlVal.value = '')
      // 拼接新的Sql
      newSql =
        validateFormSqlIndex.length > 1
          ? validateFormSqlIndex.reduce((str, validateIndex, $index) => {
              const curr = formSqlList.value[validateIndex]
              const prev = formSqlList.value[validateIndex - 1]
              return $index === 0
                ? `${curr.field} ${
                    curr.operator.endsWith('exists') ? `${curr.operator}` : `${curr.operator} ${curr.value}`
                  }`
                : `${prev.type === 'field' ? str : `${str} ${prev.operator} `}${curr.field} ${
                    curr.operator.endsWith('exists') ? `${curr.operator}` : `${curr.operator} ${curr.value}`
                  }`
            }, '')
          : `${formSqlList.value[validateFormSqlIndex[0]].field} ${
              formSqlList.value[validateFormSqlIndex[0]].operator.endsWith('exists')
                ? `${formSqlList.value[validateFormSqlIndex[0]].operator}`
                : `${formSqlList.value[validateFormSqlIndex[0]].operator} ${
                    formSqlList.value[validateFormSqlIndex[0]].value
                  }`
            }`
    } else {
      newSql = formSqlList.value
        .map((item) => {
          const { type, operator, value, field } = item
          const _value =
            type === 'operator'
              ? operator
              : operator.endsWith('exists')
              ? `${field} ${operator}`
              : `${field} ${operator} ${value}`
          return _value
        })
        .join(' ')
    }
    return newSql
  }

  // 编辑或者新增FormSql，并同步校验通过的Field
  const hanldeFormSqlEditSave = () => {
    if (!formSqlEditData.value) return
    const handleValidateItemSql = (itemSql: FormSqlItem) => {
      const { type, operator, value = '', field = '' } = itemSql

      if (operator.toLowerCase().endsWith('exists')) return true

      if (
        ['源端口', '目的端口', '源ip', '目的ip'].includes(field.toLowerCase()) &&
        operator.toLowerCase().endsWith('like')
      )
        return $baseMessage(`${field}不允许使用${operator}运算符`, 'error', 'vab-hey-message-error')

      if (type === 'operator' && !['and', 'or', 'not', '(', ')'].includes(operator.toLowerCase()))
        return $baseMessage('请检查SQL连接符是否合规', 'error', 'vab-hey-message-error')

      if (type === 'field' && (!value || !field))
        return $baseMessage('请检查SQL值是否合规', 'error', 'vab-hey-message-error')

      if (['源ip', '目的ip'].includes(field.toLowerCase())) {
        if (ipv4Regex.test(value) || ipv6Regex.test(value)) return true
        return $baseMessage('请检查IP是否合规', 'error', 'vab-hey-message-error')
      }

      if (['源端口', '目的端口'].includes(field.toLowerCase()) && !portRegex.test(value))
        return $baseMessage('请检查端口是否合规', 'error', 'vab-hey-message-error')

      return true
    }
    if (handleValidateItemSql(formSqlEditData.value) === true) {
      const { type, operator, value = '' } = formSqlEditData.value
      formSqlEditData.value.value = `"${value.replace(/(?<!\\)"/g, '\\"')}"`
      if (formSqlEditDirection.value === 'self') {
        //需要对val手动补全双引号，在合并formSql的时候不需要另外补双引号，为了是在编辑formSql时输入过程中涉及的双引号问题。
        formSqlEditData.value = formSqlEditData.value.type === 'operator' ? { type, operator } : formSqlEditData.value
        formSqlList.value[formSqlEditIndex.value] = formSqlEditData.value
      } else {
        const newIndex = formSqlEditDirection.value === 'left' ? formSqlEditIndex.value : formSqlEditIndex.value + 1
        const addFormSqlInfo = type === 'operator' ? { type, operator } : formSqlEditData.value
        formSqlList.value.splice(newIndex, 0, addFormSqlInfo)
      }
      sqlVal.value = handleSyncField2Sql()
      formSqlEditDialogVisible.value = false
    }
  }
  const handleUpdate = (value: string, viewUpdate: any) => {
    isRecommend.value = false
    handleStartCompletion()
    emits('on-change', value)
  }
  // 手动触发自动补全
  const handleStartCompletion = () => {
    if (editView.value) startCompletion(editView.value)
  }
  // 手动改变sql值，这时候需要把推荐状态改变
  const handleChangeSql = (sql: string) => {
    isRecommend.value = false
    sqlVal.value = sql
    emits('on-change', sql)
  }
  // 点击检索成功后添加历史
  const handleChangeHistories = (newHistory = '') => {
    handleGetHistories()
  }

  // 添加自定义的 SQL 关键字
  function customSQLKeywords(view: any) {
    const attackDeco = Decoration.mark({ class: 'ͼb' })
    let keywords = /not_exists|not_like/gi
    let decorations = []
    for (let { from, to } of view.visibleRanges) {
      let text = view.state.doc.sliceString(from, to)
      let match
      while ((match = keywords.exec(text))) {
        decorations.push({
          from: from + match.index,
          to: from + match.index + match[0].length,
        })
      }
    }
    return Decoration.set(decorations.map((r) => attackDeco.range(r.from, r.to)))
  }

  const formatDisplaysFiled = (baseColumn: TableColumnItemType[], baseFiled: number[]) => {
    const _baseColumn = cloneDeep(baseColumn)
    return baseFiled
      .map((id) => {
        const index = _baseColumn.findIndex((item) => item.id === id)
        const current = _baseColumn.splice(index, 1)
        return current[0]
      })
      .concat(_baseColumn)
  }
  watch(
    () => [userStore.tableColumns, props.indexType],
    debounce((newVal, oldVal) => {
      const { getTableColumn, getDisPlaysFiled } = userStore
      const baseTableColumn = getTableColumn(props.indexType)
      const baseDisPlaysFiled = getDisPlaysFiled(props.indexType)
      sqlHintPaths.value = formatColum(formatDisplaysFiled(baseTableColumn, baseDisPlaysFiled))
      cmOption.extensions = [
        sql(),
        autocompletion({
          override: [hintPlugin()],
        }),
        EditorView.decorations.of(customSQLKeywords),
        EditorView.domEventHandlers({
          paste(event, view) {
            event.preventDefault() // 阻止默认粘贴行为
            const text = event.clipboardData?.getData('text/plain')
            const from = view.state.selection.main.from
            if (text) {
              const cleanedText = text.replace(/[\r\n\t]/g, '')
              view.dispatch({
                changes: { from: from, insert: cleanedText },
                selection: {
                  anchor: cleanedText.length + from,
                },
              })
            }
            isRecommend.value = false
            return true
          },
        }),
        EditorState.transactionFilter.of((tr) => {
          if (tr.newDoc.lines > 1) {
            isFocus.value = false
            emits('on-search', sqlVal.value)
            return []
          }
          return tr
        }),
        EditorView.updateListener.of((update) => {
          if (update.focusChanged && update.view.hasFocus) {
            if (sqlVal.value) handleStartCompletion()
            isFocus.value = true
          }
        }),
      ]
    }, 500),
    {
      immediate: true,
    }
  )
  defineExpose({
    isValidate: isValidate,
    getValidateSql: () => {
      const newVal = handleFormatterSql(sqlVal.value)
      const validate = handleValidateSql(newVal, true)
      isValidate.value = validate
      sqlVal.value = newVal
      emits('on-change', newVal)
      return validate
    },
    changeSql: handleChangeSql,
    changeHistories: handleChangeHistories,
  })
  watch(
    sqlVal,
    debounce(() => {
      const newSql = handleFormatterSql(sqlVal.value)
      isValidate.value = handleValidateSql(newSql)
    }, 500)
  )

  watch(
    () => formSqlEditDialogVisible.value,
    () => {
      if (!formSqlEditDialogVisible.value || formSqlEditData.value?.type === 'operator') return
      nextTick(() => {
        sqlOperatordRef.value?.querySelector('.sql-config-operator-item.active')?.scrollIntoView({
          behavior: 'smooth',
          block: 'center',
        })
        sqlFieldRef.value?.querySelector('.sql-config-field-item.active')?.scrollIntoView({
          behavior: 'smooth',
          block: 'center',
        })
      })
    }
  )
  watchEffect(() => {
    handleGetHistories()
  })
  onUnmounted(() => {
    isFocus.value = false
  })
</script>

<template>
  <div
    ref="searchSqlRef"
    class="searchSql"
    :class="{ focus: isFocus, hasError: !isValidate }"
    @click="handleRemoveEditTag"
  >
    <codemirror
      v-model="sqlVal"
      :style="{
        width: '100%',
        height: '40px',
        backgroundColor: '#fff',
        color: '#333',
      }"
      v-bind="{ ...cmOption }"
      @change="handleUpdate"
      @ready="handleReady"
    />
    <div v-if="isFocus" class="sqlHelps">
      <div class="sqlHelps-header">
        <div class="sql-form" @mouseleave="handleRemoveEditTag">
          <div
            v-for="(sqlItem, index) in formSqlList"
            :key="index"
            class="sql-form-item"
            :class="[sqlItem.type]"
            @click="handleFormSqlEdit('self')"
            @mousemove="handleTriggerSql(sqlItem, index, $event)"
          >
            {{
              sqlItem.type === 'field'
                ? `${sqlItem.field} ${sqlItem.operator} ` +
                  `${sqlItem.operator.endsWith('exists') ? '' : `${sqlItem.value}`}`
                : `${sqlItem.operator}`
            }}
          </div>
          <div
            class="sql-form-item add"
            @click="handleFormSqlEdit('right')"
            @mousemove="handleTriggerSql({ type: 'field', operator: '=' }, formSqlList.length - 1, undefined)"
          >
            <el-icon><Plus /></el-icon>
          </div>
        </div>
      </div>
      <div class="sql-history">
        <h4>历史检索条件：</h4>
        <template v-if="searchSqlHistories.length">
          <div
            v-for="history of searchSqlHistories"
            :key="history.id"
            class="sql-item"
            @click="handleChangeSql(history.searchSql)"
          >
            {{ history.searchSql }}
          </div>
        </template>
        <el-empty v-else description="暂无历史数据" :image-size="100" style="padding-block: 20px" />
      </div>
      <div class="sql-footer">
        <span class="sql-closed" @click="isFocus = false">关闭</span>
      </div>
    </div>
    <!-- formSQL编辑菜单 -->
    <el-popover
      :append-to="searchSqlRef"
      effect="light"
      placement="bottom-start"
      :popper-style="{ minWidth: '120px', width: '120px', zIndex: 1991, marginTop: '-8px', padding: 0 }"
      :show-arrow="false"
      :virtual-ref="sqlConfigTriggerRef"
      virtual-triggering
      :visible="sqlConfigVisible"
    >
      <div class="editTools" @mouseleave="handleRemoveEditTag" @mousemove="sqlConfigVisible = true">
        <div class="toolItem" @click="handleFormSqlEdit('self')">编辑</div>
        <div class="toolItem" @click="handleFormSqlEdit('left')">在左侧添加</div>
        <div class="toolItem" @click="handleFormSqlEdit('right')">在右侧添加</div>
        <el-divider />
        <div class="toolItem" @click="handleFormSqlDelete">
          <el-icon><Delete /></el-icon>
          删除
        </div>
      </div>
    </el-popover>
    <vab-dialog
      v-model="formSqlEditDialogVisible"
      :title="formSqlEditDirection !== 'self' ? '添加' : '编辑'"
      width="775"
    >
      <div v-if="formSqlEditData" class="sql-config-content">
        <div class="sql-config-top">
          <div class="sql-config-options">
            <div
              class="sql-config-option"
              :class="{ active: formSqlEditData?.type === 'operator' }"
              @click="formSqlEditData.type = 'operator'"
            >
              逻辑符号
            </div>
            <div
              class="sql-config-option"
              :class="{ active: formSqlEditData.type === 'field' }"
              @click="formSqlEditData.type = 'field'"
            >
              检索元素
            </div>
          </div>
          <template v-if="formSqlEditData.type === 'field'">
            <div ref="sqlFieldRef" class="sql-config-field">
              <div class="sticky">
                <el-input v-model.lazy="formSqlSearchFieldStr" clearable placeholder="搜索" :suffix-icon="Search" />
              </div>
              <div
                v-for="field of formSqlFields"
                :key="field.label"
                class="sql-config-field-item"
                :class="{ active: formSqlEditData.field === field.label }"
                @click="formSqlEditData.field = field.label"
              >
                {{ field.label }}
              </div>
            </div>
            <div ref="sqlOperatordRef" class="sql-config-operator">
              <template v-for="relationalOperatorItem of relationalOperator" :key="relationalOperatorItem.label">
                <div
                  v-if="
                    (formSqlEditData.field === '方向' && ['!=', '='].includes(relationalOperatorItem.label)) ||
                    formSqlEditData.field !== '方向'
                  "
                  class="sql-config-operator-item"
                  :class="{ active: relationalOperatorItem.label === formSqlEditData.operator }"
                  @click="formSqlEditData.operator = relationalOperatorItem.label"
                >
                  {{ relationalOperatorItem.label }} {{ relationalOperatorItem.detail }}
                </div>
              </template>
            </div>
            <div class="sql-config-text">
              <el-input
                v-if="!formSqlEditData.operator.toLowerCase().endsWith('exists')"
                v-model="formSqlEditData.value"
                placeholder="请输入"
                resize="none"
                type="textarea"
              />
            </div>
          </template>
          <div v-else class="sql-config-logicOperator">
            <div
              v-for="logicOperatorItem of logicOperator"
              :key="logicOperatorItem.label"
              class="sql-config-operator-item"
              :class="{ active: logicOperatorItem.label === formSqlEditData.operator }"
              @click="formSqlEditData.operator = logicOperatorItem.label"
            >
              {{ logicOperatorItem.label }}
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button :auto-insert-space="false" type="primary" @click="hanldeFormSqlEditSave">确认</el-button>
        <el-button :auto-insert-space="false" @click="formSqlEditDialogVisible = false">取消</el-button>
      </template>
    </vab-dialog>
  </div>
</template>

<style scoped lang="scss">
  .searchSql {
    width: 100%;
    position: relative;
    &.focus {
      :deep(.cm-editor) {
        border: 1px solid var(--el-color-primary);
        .cm-activeLine {
          background: none !important;
        }
      }
    }
    &.hasError {
      :deep(.cm-editor) {
        border: 1px solid var(--el-color-danger);
      }
    }
    :deep() {
      .v-codemirror {
        position: absolute;
        inset: 0;
      }
      .cm-editor {
        border: 1px solid var(--el-border-color);
        .cm-line {
          height: 30px;
          padding: 0 11px !important;
          line-height: 30px;
        }
        .cm-cursor-primary {
          top: 11px !important;
          height: 16px !important;
        }
        .cm-tooltip {
          z-index: 9999999;
          background-color: #fff;
          ul {
            max-height: 24em;
            li {
              padding: 5px 20px;
              font-size: 16px;
              &[aria-selected] {
                background-color: var(--el-color-primary);
                .cm-completionLabel,
                .cm-completionDetail {
                  color: #fff;
                }
              }
            }
          }
          .cm-completionIcon {
            display: none;
          }
          .cm-completionLabel {
            color: var(--el-color-primary);
          }
          .cm-completionDetail {
            color: var(--el-color-info-light-5);
          }
        }
      }
    }
    .sqlHelps {
      position: absolute;
      top: 42px;
      left: 0;
      right: 0;
      height: max-content;
      z-index: 1990;
      background: #ffffff;
      box-shadow: 0px 0px 4px 4px rgba(77, 81, 86, 0.05);
      border-radius: 4px;
      border: 1px solid var(--el-border-color);
      .sqlHelps-header {
        display: flex;
        border-bottom: 1px solid var(--el-border-color);
        justify-content: space-between;
        align-items: baseline;
        padding: 10px 14px 0 4px;
        .sql-form {
          display: flex;
          flex-wrap: wrap;
          .sql-form-add {
            width: 22px;
            height: 22px;
            background: #ffffff;
            border-radius: 4px;
          }
          .sql-form-item {
            display: flex;
            margin-bottom: 8px;
            margin-inline: 10px 0;
            align-items: center;
            border-radius: 4px;
            line-height: 26px;
            padding-inline: 10px;
            font-weight: 500;
            font-size: 13px;
            cursor: pointer;
            &.operator {
              background: #edeaff;
              color: #6954f0;
              &.active,
              &:hover {
                background: #cdc4ff;
              }
            }
            &.field {
              color: #623f10;
              background: #fff8d8;
              &.active,
              &:hover {
                background: #ffeea4;
              }
            }
            &.add {
              width: 23px;
              height: 23px;
              border-radius: 4px;
              border: 1px solid #eae7f8;
              justify-content: center;
              margin-top: 2px;
              .el-icon {
                color: #9793b0;
              }
              &:hover {
                border: 1px solid var(--el-color-primary);
                .el-icon {
                  color: var(--el-color-primary);
                }
              }
            }
          }
        }
      }
      .sql-history {
        max-height: 315px;
        padding-inline: 14px;
        overflow-y: auto;
        h4 {
          font-weight: 400;
          font-size: 14px;
          color: #4a4759;
          position: sticky;
          padding-top: 10px;
          top: 0;
          background-color: #fff;
        }
        .sql-item {
          font-weight: 400;
          font-size: 13px;
          color: #635b91;
          line-height: 20px;
          background: #fafaff;
          border-radius: 4px;
          border: 1px solid #eae7f8;
          padding-inline: 10px;
          margin-bottom: 8px;
          width: fit-content;
          padding-block: 4px;
          cursor: pointer;
        }
      }
      .sql-footer {
        display: flex;
        justify-content: end;
        // margin-block: 10px;
      }
      .sql-closed {
        font-weight: 400;
        font-size: 13px;
        color: #6954f0;
        cursor: pointer;
        padding-block: 2px;
        padding-inline: 14px;
      }
    }
  }
  .sql-config-content {
    display: flex;
    flex-direction: column;
    width: 100%;
    height: 345px;
    background: #ffffff;
    border-radius: 4px;
    border: 1px solid #efecf7;
    border-top: 0;
    .sql-config-top {
      flex: 1;
      display: flex;
      overflow: hidden;
      & > div:not(:last-child) {
        border-right: 1px solid #efecf7;
      }
      .sql-config-options {
        width: 110px;
        background: #fcfcff;
        padding-inline: 10px;
        padding-top: 10px;
        .sql-config-option {
          height: 34px;
          font-weight: 500;
          font-size: 14px;
          color: #4a4759;
          line-height: 34px;
          border-radius: 4px;
          text-align: center;
          cursor: pointer;
          &.active {
            color: #fff;
            background-color: var(--el-color-primary);
          }
        }
      }
      .sql-config-field,
      .sql-config-operator,
      .sql-config-text {
        flex: 1;
        overflow-y: auto;
      }
      .sql-config-field-item,
      .sql-config-operator-item {
        width: 100%;
        height: 32px;
        border-radius: 6px;
        font-weight: 500;
        font-size: 14px;
        color: #4a4759;
        line-height: 32px;
        text-indent: 10px;
        cursor: pointer;
        &.active {
          background: #f5f3fe;
          color: var(--el-color-primary);
        }
      }
      .sql-config-operator {
        padding: 10px;
        .sql-config-operator-item {
          &:first-child {
            margin-bottom: 4px;
          }
        }
      }
      .sql-config-field {
        padding: 0 0 10px 0;
        .sticky {
          width: 100%;
          margin-bottom: 4px;
          position: sticky;
          top: 0;
          background-color: #fff;
          padding-inline: 10px;
          padding-top: 10px;
        }
        .sql-config-field-item {
          margin-inline: 10px;
          width: auto;
        }
      }
      .sql-config-logicOperator {
        flex: 1;
        padding: 10px;
        text-align: center;
      }
      .sql-config-text {
        margin: 10px;
        width: 160px;
        :deep() {
          .el-textarea,
          .el-textarea__inner {
            height: 100%;
            box-shadow: none;
          }
          .el-textarea__inner {
            padding-inline: 0;
          }
        }
      }
    }
    .sql-config-footer {
      border-top: 1px solid #efecf7;
      height: 56px;
      line-height: 56px;
      text-align: right;
      padding-inline: 14px;
    }
  }
  .editTools {
    .el-divider {
      margin-block: 8px 0;
    }
    .toolItem {
      width: 98px;
      height: 34px;
      border-radius: 4px;
      text-align: center;
      line-height: 34px;
      font-weight: 400;
      font-size: 14px;
      color: #606266;
      margin: 0 auto;
      cursor: pointer;
      .el-icon {
        vertical-align: middle;
      }
      &:last-child {
        margin-block: 4px;
      }
      &:first-child {
        margin-top: 8px;
      }
      &:hover {
        background: #f3f1fe;
        font-weight: 500;
        color: #5236ff;
      }
    }
  }
  :deep() {
    .el-dialog__body {
      padding-block: 0 10px !important;
      padding-inline: 0 !important;
    }
    .el-dialog__footer {
      .el-button {
        border-radius: 4px;
      }
    }
  }
</style>
