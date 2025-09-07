<script lang="ts">
  export default {
    name: 'LargeGraph',
  }
</script>

<script setup lang="ts">
  import G6, { ComboConfig, EdgeConfig, NodeConfig, TreeGraphData } from '@antv/g6'
  import { uniqueId } from '@antv/util'
  import { ClusterData } from '@antv/algorithm/lib/types'
  import { colorSets, global } from './large-graph-register'
  import { getSiteAppInterActionApi } from '@/api-ecs/site'
  const $baseMessage: any = inject('$baseMessage')
  const props = defineProps<{
    siteId?: number | string
    isFull?: boolean
  }>()
  let data: any = {}
  const graphNodesMap = new Map()
  const modelMap = new Map()
  const NODESIZEMAPPING = 'degree'
  const SMALLGRAPHLABELMAXLENGTH = 20
  let labelMaxLength = SMALLGRAPHLABELMAXLENGTH
  const DEFAULTNODESIZE = 20
  const DEFAULTAGGREGATEDNODESIZE = 53
  const NODE_LIMIT = 100 // 页面元素最多上限，超出后只能展开一个站点
  let cachePositions: any = {}
  let descreteNodeCenter: any
  let manipulatePosition: any | undefined = undefined
  let currentUnproccessedData: {
    nodes: any[]
    edges: any[]
  } = { nodes: [], edges: [] }
  let expandArray: any[] = []
  let collapseArray: any[] = []

  let CANVAS_WIDTH = 0,
    CANVAS_HEIGHT = 0
  let shiftKeydown = false
  let nodeMap: any = {}
  let graph: null | any = null
  let hiddenItemIds: any[] = []
  let aggregatedNodeMap: any = {}

  let layout = {
    type: '',
    instance: null as any,
    destroyed: true,
  }

  const graphRef = ref<HTMLDivElement>()

  const descendCompare = (p: any) => {
    // 这是比较函数
    return function (m: any, n: any) {
      const a = m[p]
      const b = n[p]
      return b - a // 降序
    }
  }
  const labelFormatter = (text: string, minLength = 20): string => {
    if (text && text.split('').length > minLength) return `${text.substr(0, minLength)}...`
    return text
  }
  // 截断长文本。length 为文本截断后长度，elipsis 是后缀
  const formatText = (text: string, length = 10, elipsis = '...') => {
    if (!text) return ''
    if (text.length > length) {
      return `${text.substr(0, length)}${elipsis}`
    }
    return text
  }
  const processNodesEdges = (
    nodes: any[],
    edges: any[],
    width: number,
    height: number,
    largeGraphMode: any,
    edgeLabelVisible: boolean,
    isNewGraph = false
  ) => {
    if (!nodes || nodes.length === 0) return {}
    const currentNodeMap: any = {}
    let maxNodeCount = -Infinity
    const paddingRatio = 0.3
    const paddingLeft = paddingRatio * width
    const paddingTop = paddingRatio * height
    nodes.forEach(
      (node: {
        type: string
        level: number
        isReal: boolean
        label: string
        id: string
        labelLineNum: undefined
        oriLabel: any
        degree: number
        inDegree: number
        outDegree: number
        count: number
        x: any
        y: any
        new: boolean
      }) => {
        node.type = node.level === 0 ? 'real-node' : 'aggregated-node'
        node.isReal = node.level === 0 ? true : false
        node.label = `${node.label}`
        node.labelLineNum = undefined
        node.oriLabel = node.label
        node.label = formatText(node.label, labelMaxLength, '...')
        node.degree = 0
        node.inDegree = 0
        node.outDegree = 0
        if (currentNodeMap[node.id]) {
          console.warn('node exists already!', node.id)
          node.id = `${node.id}${Math.random()}`
        }
        currentNodeMap[node.id] = node
        if (node.count > maxNodeCount) maxNodeCount = node.count
        const cachePosition = cachePositions ? cachePositions[node.id] : undefined
        if (cachePosition) {
          node.x = cachePosition.x
          node.y = cachePosition.y
          node.new = false
        } else {
          node.new = isNewGraph ? false : true
          if (manipulatePosition && !node.x && !node.y) {
            node.x = manipulatePosition.x + 30 * Math.cos(Math.random() * Math.PI * 2)
            node.y = manipulatePosition.y + 30 * Math.sin(Math.random() * Math.PI * 2)
          }
        }
      }
    )

    let maxCount = -Infinity
    let minCount = Infinity
    // let maxCount = 0;
    edges.forEach((edge: { id: string; source: string | number; target: string | number; count: number }) => {
      // to avoid the dulplicated id to nodes
      if (!edge.id) edge.id = uniqueId('edge')
      else if (edge.id.split('-')[0] !== 'edge') edge.id = `edge-${edge.id}`
      // TODO: delete the following line after the queried data is correct
      if (!currentNodeMap[edge.source] || !currentNodeMap[edge.target]) {
        console.warn('edge source target does not exist', edge.source, edge.target, edge.id)
        return
      }
      const sourceNode = currentNodeMap[edge.source]
      const targetNode = currentNodeMap[edge.target]
      if (!sourceNode || !targetNode) console.warn('source or target is not defined!!!', edge, sourceNode, targetNode)

      // calculate the degree
      sourceNode.degree++
      targetNode.degree++
      sourceNode.outDegree++
      targetNode.inDegree++

      if (edge.count > maxCount) maxCount = edge.count
      if (edge.count < minCount) minCount = edge.count
    })

    nodes.sort(descendCompare(NODESIZEMAPPING))
    const maxDegree = nodes[0].degree || 1

    const descreteNodes: any[] = []
    nodes.forEach(
      (
        node: {
          count: number
          level: number
          size: any
          isReal: boolean
          labelCfg: {
            position: string
            offset: number
            style: {
              fill: any
              fontSize: number
              stroke: any
              lineWidth: number
            }
          }
          degree: any
        },
        i: any
      ) => {
        // assign the size mapping to the outDegree
        const countRatio = node.count / maxNodeCount
        const isRealNode = node.level === 0
        node.size = isRealNode ? DEFAULTNODESIZE : DEFAULTAGGREGATEDNODESIZE
        node.isReal = isRealNode
        node.labelCfg = {
          position: 'bottom',
          offset: 5,
          style: {
            fill: global.node.labelCfg.style.fill,
            fontSize: 6 + countRatio * 6 || 12,
            stroke: global.node.labelCfg.style.stroke,
            lineWidth: 3,
          },
        }

        if (!node.degree) {
          descreteNodes.push(node)
        }
      }
    )

    const countRange = maxCount - minCount
    const minEdgeSize = 1
    const maxEdgeSize = 7
    const edgeSizeRange = maxEdgeSize - minEdgeSize
    edges.forEach((edge: any) => {
      const targetNode = currentNodeMap[edge.target]
      const size = ((edge.count - minCount) / countRange) * edgeSizeRange + minEdgeSize || 1
      edge.size = size
      const arrowWidth = Math.max(size / 2 + 2, 3)
      const arrowLength = 10
      const arrowBeging = targetNode.size + arrowLength
      let arrowPath: string | undefined = `M ${arrowBeging},0 L ${arrowBeging + arrowLength},-${arrowWidth} L ${
        arrowBeging + arrowLength
      },${arrowWidth} Z`
      let d = targetNode.size / 2 + arrowLength
      if (edge.source === edge.target) {
        edge.type = 'line'
        // edge.type = 'loop'
        arrowPath = undefined
      }
      const sourceNode = currentNodeMap[edge.source]
      const isRealEdge = targetNode.isReal && sourceNode.isReal
      edge.isReal = isRealEdge
      const stroke = isRealEdge ? global.edge.style.realEdgeStroke : global.edge.style.stroke
      const opacity = isRealEdge ? global.edge.style.realEdgeOpacity : global.edge.style.strokeOpacity

      const dash = Math.max(size, 2)
      const lineDash = isRealEdge ? undefined : [dash, dash]
      // const lineDash = [dash, dash]
      edge.style = {
        stroke,
        strokeOpacity: opacity,
        cursor: 'pointer',
        lineAppendWidth: Math.max(edge.size || 5, 5),
        fillOpacity: 1,
        lineDash,
        lineWidth: 2,
        endArrow: arrowPath
          ? {
              path: arrowPath,
              d,
              fill: stroke,
              strokeOpacity: 0,
            }
          : false,
      }
      edge.labelCfg = {
        autoRotate: true,
        style: {
          stroke: global.edge.labelCfg.style.stroke,
          fill: global.edge.labelCfg.style.fill,
          lineWidth: 2,
          fontSize: 12,
          lineAppendWidth: 10,
          opacity: 1,
        },
      }
      if (!edge.oriLabel) edge.oriLabel = edge.label
      if (largeGraphMode || !edgeLabelVisible) edge.label = ''
      else {
        edge.label = labelFormatter(edge.label, labelMaxLength)
      }

      // arrange the other nodes around the hub
      const sourceDis = sourceNode.size / 2 + 20
      const targetDis = targetNode.size / 2 + 20
      if (sourceNode.x && !targetNode.x) {
        targetNode.x = sourceNode.x + sourceDis * Math.cos(Math.random() * Math.PI * 2)
      }
      if (sourceNode.y && !targetNode.y) {
        targetNode.y = sourceNode.y + sourceDis * Math.sin(Math.random() * Math.PI * 2)
      }
      if (targetNode.x && !sourceNode.x) {
        sourceNode.x = targetNode.x + targetDis * Math.cos(Math.random() * Math.PI * 2)
      }
      if (targetNode.y && !sourceNode.y) {
        sourceNode.y = targetNode.y + targetDis * Math.sin(Math.random() * Math.PI * 2)
      }

      if (!sourceNode.x && !sourceNode.y && manipulatePosition) {
        sourceNode.x = manipulatePosition.x + 30 * Math.cos(Math.random() * Math.PI * 2)
        sourceNode.y = manipulatePosition.y + 30 * Math.sin(Math.random() * Math.PI * 2)
      }
      if (!targetNode.x && !targetNode.y && manipulatePosition) {
        targetNode.x = manipulatePosition.x + 30 * Math.cos(Math.random() * Math.PI * 2)
        targetNode.y = manipulatePosition.y + 30 * Math.sin(Math.random() * Math.PI * 2)
      }
    })

    descreteNodeCenter = {
      x: width - paddingLeft,
      y: height - paddingTop,
    }
    descreteNodes.forEach((node) => {
      if (!node.x && !node.y) {
        node.x = descreteNodeCenter.x + 30 * Math.cos(Math.random() * Math.PI * 2)
        node.y = descreteNodeCenter.y + 30 * Math.sin(Math.random() * Math.PI * 2)
      }
    })

    G6.Util.processParallelEdges(edges, 12.5, 'custom-quadratic', 'custom-line')
    return {
      maxDegree,
      edges,
    }
  }
  const examAncestors = (
    model: NodeConfig | EdgeConfig | ComboConfig | TreeGraphData,
    expandedArray: any[],
    length: number,
    keepTags: any
  ) => {
    for (let i = 0; i < length; i++) {
      const expandedNode = expandedArray[i]
      if (!keepTags[i] && model.parentId === expandedNode.id) {
        keepTags[i] = true // 需要被保留
        examAncestors(expandedNode, expandedArray, length, keepTags)
        break
      }
    }
  }
  const manageExpandCollapseArray = (
    nodeNumber: number,
    model: NodeConfig | EdgeConfig | ComboConfig | TreeGraphData,
    collapseArray: { id: any; parentId: any; level: number }[],
    expandArray: any[]
  ) => {
    manipulatePosition = { x: model.x, y: model.y }

    // 维护 expandArray，若当前画布节点数高于上限，移出 expandedArray 中非 model 祖先的节点)
    if (nodeNumber > NODE_LIMIT) {
      // 若 keepTags[i] 为 true，则 expandedArray 的第 i 个节点需要被保留
      const keepTags: any = {}
      const expandLen = expandArray.length
      // 检查 X 的所有祖先并标记 keepTags
      examAncestors(model, expandArray, expandLen, keepTags)
      // 寻找 expandedArray 中第一个 keepTags 不为 true 的点
      let shiftNodeIdx = -1
      for (let i = 0; i < expandLen; i++) {
        if (!keepTags[i]) {
          shiftNodeIdx = i
          break
        }
      }
      // 如果有符合条件的节点，将其从 expandedArray 中移除
      if (shiftNodeIdx !== -1) {
        let foundNode = expandArray[shiftNodeIdx]
        if (foundNode.level === 2) {
          let foundLevel1 = false
          // 找到 expandedArray 中 parentId = foundNode.id 且 level = 1 的第一个节点
          for (let i = 0; i < expandLen; i++) {
            const eNode = expandArray[i]
            if (eNode.parentId === foundNode.id && eNode.level === 1) {
              foundLevel1 = true
              foundNode = eNode
              expandArray.splice(i, 1)
              break
            }
          }
          // 若未找到，则 foundNode 不变, 直接删去 foundNode
          if (!foundLevel1) expandArray.splice(shiftNodeIdx, 1)
        } else {
          // 直接删去 foundNode
          expandArray.splice(shiftNodeIdx, 1)
        }
        // const removedNode = expandedArray.splice(shiftNodeIdx, 1); // splice returns an array
        const idSplits = foundNode.id.split('-')
        let collapseNodeId
        // 去掉最后一个后缀
        for (let i = 0; i < idSplits.length - 1; i++) {
          const str = idSplits[i]
          if (collapseNodeId) collapseNodeId = `${collapseNodeId}-${str}`
          else collapseNodeId = str
        }
        const collapseNode = {
          id: collapseNodeId,
          parentId: foundNode.id,
          level: foundNode.level - 1,
        }
        collapseArray.push(collapseNode)
      }
    }

    const currentNode = {
      id: model.id,
      level: model.level,
      parentId: model.parentId,
    }

    // 加入当前需要展开的节点
    expandArray.push(currentNode)

    graph.get('canvas').setCursor('default')
    return { expandArray, collapseArray }
  }
  const getMixedGraph = (
    aggregatedData: ClusterData,
    originData: { nodes?: { id: string }[]; edges: any },
    nodeMap: { [x: string]: { clusterId: any } },
    aggregatedNodeMap: { [x: string]: { expanded: boolean } },
    expandArray: any[],
    collapseArray: any[]
  ) => {
    let nodes: any[] = [],
      edges: { source: any; target: any; id: string; label: string }[] = []
    /***
     *  @aggregatedData 聚合后的数据
     *  @aggregatedNodeMap 聚合节点数据Map
     *  @originData 查询的初始值
     *  @modelMap 所有站点下的应用（聚合的节点）
     *  @_modelMap 把modelMap的key取出根据展开的id取节点
     *  @expandKeys 展开的id
     *  @collapseKeys 闭合的id
     */

    const expandKeys = expandArray.map((expand) => expand.id)

    aggregatedData.clusters.forEach((cluster: { id: string | number; nodes: any }, i: any) => {
      if (expandKeys.includes(cluster.id)) {
        //聚合数据的id包含在展开字段里，把当前聚合数据的下的node添加到node数组
        nodes = nodes.concat(cluster.nodes)
        aggregatedNodeMap[cluster.id].expanded = true
      } else {
        //如果不包含就证明是其他的数据，把当前数据添加进去，标识设为false
        nodes.push(aggregatedNodeMap[cluster.id])
        aggregatedNodeMap[cluster.id].expanded = false
      }
    })
    // 原始数据links遍历
    originData.edges.forEach((edge: any) => {
      const isSourceInExpandArray = expandKeys.includes(nodeMap[edge.source].clusterId)
      const isTargetInExpandArray = expandKeys.includes(nodeMap[edge.target].clusterId)
      if (isSourceInExpandArray && isTargetInExpandArray) {
        edges.push(edge)
      } else if (isSourceInExpandArray) {
        const targetClusterId = nodeMap[edge.target].clusterId
        const vedge = {
          source: edge.source,
          target: targetClusterId,
          id: uniqueId('edge'),
          label: '',
        }
        edges.push(vedge)
      } else if (isTargetInExpandArray) {
        const sourceClusterId = nodeMap[edge.source].clusterId
        const vedge = {
          target: edge.target,
          source: sourceClusterId,
          id: uniqueId('edge'),
          label: '',
        }
        edges.push(vedge)
      }
    })
    aggregatedData.clusterEdges.forEach((edge: any) => {
      if (expandKeys.includes(edge.source) || expandKeys.includes(edge.target)) return
      else edges.push(edge)
    })
    return { nodes, edges }
  }
  const generateNeighbors = (
    centerNodeModel: {
      x?: any
      y?: any
      clusterId: any
      id?: any
      level?: number
      colorSet?: any
    },
    step: number,
    maxNeighborNumPerNode = 5
  ) => {
    if (step <= 0) return undefined
    let nodes: any[] = [],
      edges: any[] = []
    const clusterId = centerNodeModel.clusterId
    const centerId = centerNodeModel.id
    const neighborNum = Math.ceil(Math.random() * maxNeighborNumPerNode)
    for (let i = 0; i < neighborNum; i++) {
      const neighborNode = {
        id: uniqueId('node'),
        clusterId,
        level: 0,
        colorSet: centerNodeModel.colorSet,
      }
      nodes.push(neighborNode)
      const dire = Math.random() > 0.5
      const source = dire ? centerId : neighborNode.id
      const target = dire ? neighborNode.id : centerId
      const neighborEdge = {
        id: uniqueId('edge'),
        source,
        target,
        label: `${source}-${target}`,
      }
      edges.push(neighborEdge)
      const subNeighbors = generateNeighbors(neighborNode, step - 1, maxNeighborNumPerNode)
      if (subNeighbors) {
        nodes = nodes.concat(subNeighbors.nodes)
        edges = edges.concat(subNeighbors.edges)
      }
    }
    return { nodes, edges }
  }

  const showItems = (graph: any) => {
    graph.getNodes().forEach((node: any) => {
      if (!node.isVisible()) graph.showItem(node)
    })
    graph.getEdges().forEach((edge: any) => {
      if (!edge.isVisible()) edge.showItem(edge)
    })
    hiddenItemIds = []
  }

  const cacheNodePositions = (nodes: any) => {
    const positionMap: any = {}
    const nodeLength = nodes.length
    for (let i = 0; i < nodeLength; i++) {
      const node = nodes[i].getModel()
      positionMap[node.id] = {
        x: node.x,
        y: node.y,
        level: node.level,
      }
    }
    return positionMap
  }
  const clearFocusItemState = (graph: any) => {
    if (!graph) return
    clearFocusNodeState(graph)
    clearFocusEdgeState(graph)
  }
  // 清除图上所有节点的 focus 状态及相应样式
  const clearFocusNodeState = (graph: any) => {
    const focusNodes = graph.findAllByState('node', 'focus')
    focusNodes.forEach((fnode: any) => {
      graph.setItemState(fnode, 'focus', false) // false
    })
  }

  // 清除图上所有边的 focus 状态及相应样式
  const clearFocusEdgeState = (graph: any) => {
    const focusEdges = graph.findAllByState('edge', 'focus')
    focusEdges.forEach((fedge: any) => {
      graph.setItemState(fedge, 'focus', false)
    })
  }

  const hideItems = (graph: any) => {
    hiddenItemIds.forEach((id) => {
      graph.hideItem(id)
    })
  }
  const getForceLayoutConfig = (graph: any, largeGraphMode: any, configSettings?: any) => {
    let {
      linkDistance,
      edgeStrength,
      nodeStrength,
      nodeSpacing,
      preventOverlap,
      nodeSize,
      collideStrength,
      alpha,
      alphaDecay,
      alphaMin,
    } = configSettings || { preventOverlap: true }

    if (!linkDistance && linkDistance !== 0) linkDistance = 225
    if (!edgeStrength && edgeStrength !== 0) edgeStrength = 50
    if (!nodeStrength && nodeStrength !== 0) nodeStrength = 200
    if (!nodeSpacing && nodeSpacing !== 0) nodeSpacing = 5

    const config: any = {
      type: 'gForce',
      minMovement: 0.01,
      maxIteration: 5000,
      preventOverlap,
      damping: 0.99,
      gpuEnabled: largeGraphMode,
      linkDistance: (d: any) => {
        let dist = linkDistance

        const sourceNode = nodeMap[d.source] || aggregatedNodeMap[d.source]
        const targetNode = nodeMap[d.target] || aggregatedNodeMap[d.target]

        // // 两端都是聚合点
        // if (sourceNode.level && targetNode.level) dist = linkDistance * 3;
        // // 一端是聚合点，一端是真实节点
        // else if (sourceNode.level || targetNode.level) dist = linkDistance * 1.5;
        if (!sourceNode?.level && !targetNode?.level) dist = linkDistance * 0.3
        return dist
      },
      edgeStrength: (d: any) => {
        const sourceNode = nodeMap[d.source] || aggregatedNodeMap[d.source]
        const targetNode = nodeMap[d.target] || aggregatedNodeMap[d.target]
        // 聚合节点之间的引力小
        if (sourceNode?.level && targetNode?.level) return edgeStrength / 2
        // 聚合节点与真实节点之间引力大
        if (sourceNode?.level || targetNode?.level) return edgeStrength
        return edgeStrength
      },
      nodeStrength: (d: any) => {
        // 给离散点引力，让它们聚集
        if (d.degree === 0) return -10
        // 聚合点的斥力大
        if (d.level) return nodeStrength * 2
        return nodeStrength
      },
      nodeSize: (d: any) => {
        if (!nodeSize && d.size) return d.size
        return 50
      },
      nodeSpacing: (d: any) => {
        if (d.degree === 0) return nodeSpacing * 2
        if (d.level) return nodeSpacing
        return nodeSpacing
      },
      onLayoutEnd: () => {
        if (largeGraphMode) {
          graph.getEdges()?.forEach((edge: any) => {
            if (!edge.oriLabel) return
            edge.update({
              label: labelFormatter(edge.oriLabel, labelMaxLength),
            })
          })
        }
      },
      tick: () => {
        graph.refreshPositions()
      },
    }

    if (nodeSize) config['nodeSize'] = nodeSize
    if (collideStrength) config['collideStrength'] = collideStrength
    if (alpha) config['alpha'] = alpha
    if (alphaDecay) config['alphaDecay'] = alphaDecay
    if (alphaMin) config['alphaMin'] = alphaMin

    return config
  }
  const stopLayout = () => {
    layout.instance.stop()
  }
  const handleRefreshGraph = (
    graph: {
      getNodes: () => any[]
      getEdges: () => any[]
      changeData: (arg0: { nodes: any; edges: any[] }) => void
    },
    graphData: { nodes: any; edges: any },
    width: number,
    height: number,
    largeGraphMode: any,
    edgeLabelVisible: boolean,
    isNewGraph: boolean | undefined
  ) => {
    if (!graphData || !graph) return
    clearFocusItemState(graph)
    // reset the filtering
    graph.getNodes().forEach((node: { isVisible: () => any; show: () => void }) => {
      if (!node.isVisible()) node.show()
    })
    graph.getEdges().forEach((edge: { isVisible: () => any; show: () => void }) => {
      if (!edge.isVisible()) edge.show()
    })

    let nodes: any[] = [],
      edges: any[] = []

    nodes = graphData?.nodes
    const processRes = processNodesEdges(
      nodes,
      graphData.edges || [],
      width,
      height,
      largeGraphMode,
      edgeLabelVisible,
      isNewGraph
    )
    edges = processRes?.edges || []
    graph.changeData({ nodes, edges })
    hideItems(graph)
    graph.getNodes().forEach((node: { toFront: () => void }) => {
      node.toFront()
    })
    layout?.instance!.stop()
    // 在大数据量时，使用 GPU 布局
    if (nodes.length > 30) {
      layout.instance.destroy()
      const layoutConfig: any = getForceLayoutConfig(graph, true)
      layoutConfig.center = [width / 2, height / 2]
      layout.instance = new G6.Layout['gForce'](layoutConfig)
    }

    // force 需要使用不同 id 的对象才能进行全新的布局，否则会使用原来的引用。因此复制一份节点和边作为 force 的布局数据
    layout.instance.init({
      nodes: graphData.nodes,
      edges,
    })

    layout.instance.minMovement = 0.0001
    layout.instance.getMass = (d: { id: string | number }) => {
      const cachePosition = cachePositions[d.id]
      if (cachePosition) return 10
      return 1
    }

    layout.instance.execute()
    return { nodes, edges }
  }
  const bindListener = (graph: any) => {
    graph.on('node:mouseenter', (evt: any) => {
      const { item } = evt
      const model = item.getModel()
      const currentLabel = model.label
      model.oriFontSize = model.labelCfg.style.fontSize
      item.update({
        label: model.oriLabel,
      })
      model.oriLabel = currentLabel
      graph.setItemState(item, 'hover', true)
      item.toFront()
    })

    graph.on('node:mouseleave', (evt: any) => {
      const { item } = evt
      const model = item.getModel()
      const currentLabel = model.label
      item.update({
        label: model.oriLabel,
      })
      model.oriLabel = currentLabel
      graph.setItemState(item, 'hover', false)
    })

    graph.on('edge:mouseenter', (evt: any) => {
      const { item } = evt
      const model = item.getModel()
      const currentLabel = model.label
      item.update({
        label: model.oriLabel,
      })
      model.oriLabel = currentLabel
      item.toFront()
      item.getSource().toFront()
      item.getTarget().toFront()
    })

    graph.on('edge:mouseleave', (evt: any) => {
      const { item } = evt
      const model = item.getModel()
      const currentLabel = model.label
      item.update({
        label: model.oriLabel,
      })
      model.oriLabel = currentLabel
    })
    // click node to show the detail drawer
    graph.on('node:click', (evt: any) => {
      console.log(evt, 'evt_node')
      stopLayout()
      if (!shiftKeydown) clearFocusItemState(graph)
      else clearFocusEdgeState(graph)
      const { item } = evt

      // highlight the clicked node, it is down by click-select
      graph.setItemState(item, 'focus', true)

      if (!shiftKeydown) {
        // 将相关边也高亮
        const relatedEdges = item.getEdges()
        relatedEdges.forEach((edge: any) => {
          graph.setItemState(edge, 'focus', true)
        })
      }
    })

    // click edge to show the detail of integrated edge drawer
    graph.on('edge:click', (evt: any) => {
      console.log(evt, 'evt_edge')
      stopLayout()
      if (!shiftKeydown) clearFocusItemState(graph)
      const { item } = evt
      // highlight the clicked edge
      graph.setItemState(item, 'focus', true)
    })

    // click canvas to cancel all the focus state
    graph.on('canvas:click', (evt: any) => {
      clearFocusItemState(graph)
    })
  }

  const louvainHandle = (data: { edges: any[]; nodes: any[] }) => {
    const getModels = (nodes: any[]) => {
      nodes.forEach((node: any) => {
        node.id = node.appid
        node.business = node.business || '未知业务'
        const type = node.business
        graphNodesMap.set(node.id, node)
        if (modelMap.has(type)) {
          modelMap.get(type).push(node)
        } else {
          modelMap.set(type, [node])
        }
      })
      return [...modelMap.entries()].map(([model, nodes], index) => ({
        id: (index + 1).toString(),
        nodes: nodes.map((item: any) => ({
          ...item,
          clusterId: (index + 1).toString(),
        })),
        label: model,
        sumTot: nodes.length,
      }))
    }
    const clusteredNode = getModels(data.nodes)
    const getLinks = (links: any[]) => {
      const linksMap = new Map()
      for (const link of links) {
        const sourceNode = graphNodesMap.get(link.source)
        const targetNode = graphNodesMap.get(link.target)
        const clusteredSourceNode = clusteredNode.find((node: any) => node.label === sourceNode.business) as any
        const clusteredTargetNode = clusteredNode.find((node: any) => node.label === targetNode.business) as any
        const linkKey = `${clusteredSourceNode.id}-${clusteredTargetNode.id}`
        if (linksMap.has(linkKey)) continue
        linksMap.set(linkKey, {
          source: clusteredSourceNode.id,
          target: clusteredTargetNode.id,
          count: 2,
          weight: 2,
        })
      }
      return [...linksMap.values()]
    }
    const clusteredLink = getLinks(data.edges)
    return {
      clusterEdges: clusteredLink,
      clusters: clusteredNode,
    }
  }

  const initChart = () => {
    const clusteredData = louvainHandle(data)
    const aggregatedData: { nodes: any[]; edges: any[] } = {
      nodes: [],
      edges: [],
    }
    const parentNode = graphRef.value!.parentNode
    const { offsetWidth, offsetHeight } = parentNode as HTMLElement
    CANVAS_WIDTH = offsetWidth
    CANVAS_HEIGHT = offsetHeight
    clusteredData.clusters.forEach((cluster, i) => {
      cluster.nodes.forEach((node: any) => {
        node.level = 0
        node.label = node.id
        node.type = ''
        node.colorSet = colorSets[i]
        nodeMap[node.id] = node
      })
      const cnode = {
        id: cluster.id,
        type: 'aggregated-node',
        count: cluster.nodes.length,
        level: 1,
        label: cluster?.label,
        colorSet: colorSets[i],
        idx: i,
      }
      aggregatedNodeMap[cluster.id] = cnode
      aggregatedData.nodes.push(cnode)
    })
    clusteredData.clusterEdges.forEach((clusterEdge) => {
      const cedge: any = {
        ...clusterEdge,
        size: Math.log(clusterEdge.count as number),
        label: '',
        id: uniqueId('edge'),
      }
      aggregatedData.edges.push(cedge)
    })
    // 线鼠标放上去显示标签
    data.edges.forEach((edge: any) => {
      // edge.label = `${edge.source}-${edge.target}`
      edge.id = uniqueId('edge')
    })
    currentUnproccessedData = aggregatedData
    const { edges: processedEdges } = processNodesEdges(
      currentUnproccessedData.nodes,
      currentUnproccessedData.edges,
      CANVAS_WIDTH,
      CANVAS_HEIGHT,
      true,
      true,
      true
    )

    const contextMenu = new G6.Menu({
      shouldBegin(evt: any) {
        if (evt.target && evt.target.isCanvas && evt.target.isCanvas()) return true
        if (evt.item) return true
        return false
      },
      getContent(evt: any) {
        const { item } = evt
        if (evt.target && evt.target.isCanvas && evt.target.isCanvas()) {
          return `<ul>
                    <li id='show'>显示所有隐藏元素</li>
                    <li id='collapseAll'>聚合所有站点</li>
                  </ul>`
        } else if (!item) return ''

        const itemType = item?.getType()
        const model = item?.getModel()
        if (itemType && model) {
          if (itemType === 'node') {
            if (model.level !== 0) {
              return `<ul>
                        <li id='expand'>展开该站点</li>
                        <li id='hide'>隐藏该节点</li>
                      </ul>`
            } else {
              return `<ul>
                        <li id='collapse'>聚合所属站点</li>
                        <li id='hide'>隐藏该节点</li>
                      </ul>`
            }
          } else return ''
        } else return ''
      },
      handleMenuClick: (target, item) => {
        const model: any = item?.getModel()
        const liIdStrs = target.id.split('-')
        let mixedGraphData
        switch (liIdStrs[0]) {
          case 'hide':
            graph.hideItem(item)
            hiddenItemIds.push(model.id)
            break
          case 'expand': {
            const newArray = manageExpandCollapseArray(graph.getNodes().length, model, collapseArray, expandArray)
            expandArray = newArray.expandArray
            collapseArray = newArray.collapseArray
            mixedGraphData = getMixedGraph(clusteredData, data, nodeMap, aggregatedNodeMap, expandArray, collapseArray)
            break
          }
          case 'collapse': {
            const aggregatedNode = aggregatedNodeMap[model.clusterId]
            manipulatePosition = { x: aggregatedNode.x, y: aggregatedNode.y }
            collapseArray.push(aggregatedNode)
            for (let i = 0; i < expandArray.length; i++) {
              if (expandArray[i].id === model.clusterId) {
                expandArray.splice(i, 1)
                break
              }
            }
            mixedGraphData = getMixedGraph(clusteredData, data, nodeMap, aggregatedNodeMap, expandArray, collapseArray)
            break
          }
          case 'collapseAll':
            {
              expandArray = []
              collapseArray = []
              mixedGraphData = getMixedGraph(
                clusteredData,
                data,
                nodeMap,
                aggregatedNodeMap,
                expandArray,
                collapseArray
              )
            }
            break
          case 'show': {
            showItems(graph)
            break
          }
          default:
            break
        }
        if (mixedGraphData) {
          cachePositions = cacheNodePositions(graph.getNodes())
          currentUnproccessedData = mixedGraphData
          handleRefreshGraph(graph, currentUnproccessedData, CANVAS_WIDTH, CANVAS_HEIGHT, true, true, false)
        }
      },
      offsetX: 16 + 10,
      offsetY: 0,
      itemTypes: ['node', 'edge', 'canvas'],
    })

    graph = new G6.Graph({
      container: graphRef.value! as HTMLElement,
      linkCenter: true,
      width: CANVAS_WIDTH,
      height: CANVAS_HEIGHT,
      minZoom: 0.1,
      groupByTypes: false,
      modes: {
        default: [
          {
            type: 'drag-canvas',
            enableOptimize: true,
          },
          {
            type: 'zoom-canvas',
            enableOptimize: true,
            optimizeZoom: 0.01,
          },
          'drag-node',
          'shortcuts-call',
        ],
        lassoSelect: [
          {
            type: 'zoom-canvas',
            enableOptimize: true,
          },
          {
            type: 'lasso-select',
            selectedState: 'focus',
            trigger: 'drag',
          },
        ],
        fisheyeMode: [],
      },
      defaultNode: {
        type: 'aggregated-node',
        size: DEFAULTNODESIZE,
      },
      plugins: [contextMenu],
    })

    graph.get('canvas').set('localRefresh', false)

    const layoutConfig: any = getForceLayoutConfig(graph, true)
    layoutConfig.center = [CANVAS_WIDTH / 2, CANVAS_HEIGHT / 2]
    layout.instance = new G6.Layout['gForce'](layoutConfig)
    layout.instance.init({
      nodes: currentUnproccessedData.nodes,
      edges: processedEdges,
    })
    layout.instance.execute()
    bindListener(graph)
    graph.data({ nodes: aggregatedData.nodes, edges: processedEdges })
    graph.render()
  }
  const destroyChart = () => {
    graphNodesMap.clear()
    modelMap.clear()
    cachePositions = {}
    descreteNodeCenter = {}
    manipulatePosition = undefined
    currentUnproccessedData = { nodes: [], edges: [] }
    expandArray = []
    collapseArray = []

    shiftKeydown = false
    nodeMap = {}
    hiddenItemIds = []
    aggregatedNodeMap = {}

    layout = {
      type: '',
      instance: null as any,
      destroyed: true,
    }

    graph?.destroy()
    graph = null
  }

  const getSiteAppInterActionHandle = async () => {
    const { data: appData } = await getSiteAppInterActionApi(+props.siteId!)
    if (!appData || appData.links.length === 0 || appData.nodes.length === 0)
      return $baseMessage('暂无访问数据', 'warning', 'vab-hey-message-warning')
    const _links = Array.isArray(appData.links) ? appData.links : []
    const _nodes = Array.isArray(appData.nodes) ? appData.nodes : []
    data.edges = _links
    data.nodes = _nodes
    initChart()
  }

  watch(
    () => props.siteId,
    () => {
      if (props.siteId) {
        destroyChart()
        nextTick(() => {
          getSiteAppInterActionHandle()
        })
      }
    },
    {
      immediate: true,
    }
  )
  watch(
    () => props.isFull,
    () => {
      nextTick(() => {
        const { offsetWidth, offsetHeight } = graphRef.value!.parentNode as HTMLDivElement
        graph?.changeSize(offsetWidth, props.isFull ? offsetHeight - 40 : offsetHeight)
      })
    }
  )
  onUnmounted(() => {
    destroyChart()
  })
</script>

<template>
  <div ref="graphRef" class="largeGraph"></div>
</template>

<style scoped lang="scss">
  .largeGraph {
    height: 100%;
    width: 100%;
  }
</style>
