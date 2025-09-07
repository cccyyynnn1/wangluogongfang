import { link, node } from '~/src/types'
import _ from 'lodash'
import { gp } from '@gp'
import { DeleteRemarkApi, UpdateOrSaveRemarkApi } from '~/src/api-ecs/chase'

/**
 * @description '有环有向引力拓扑图'
 * @constructor node: node[], link: link[]
 */
export default class Sketchpad {
  allNodes: node[]
  allLinks: link[]
  nodesDefault: node[]
  linksDefault: link[]
  // 构造函数
  constructor(node: node[], link: link[]) {
    this.allNodes = node
    this.allLinks = link
    this.nodesDefault = []
    this.linksDefault = []
  }

  // 根据ip递归获取下一级数据
  public RecursionGetNextLevelDataByIP(
    ip: string,
    mode: 'one' | 'all',
    handleAddNode: Function,
    handleAddLink: Function
  ) {
    let index = 3
    this.nodesDefault.forEach((node: node) => {
      if (node.ip == ip) {
        index = node.level || 3
      }
    })
    // 找到以当前点做起点的线
    const link = this.allLinks.filter((item: link) => {
      return item.clientIp == ip
    })
    link.forEach((td: link) => {
      // 找到所有以该起点做终点的点
      const node = this.allNodes.filter((item: node) => {
        return td.serverIp == item.ip
      })
      const addnNodeArr: node[] = []
      node.forEach((e: node) => {
        // 去重
        const falg = this.nodesDefault.some((element: node) => {
          return element.ip == e.ip
        })
        // 追加
        if (!falg) {
          e['level'] = index + 1
          this.nodesDefault.push(e)
          addnNodeArr.push(e)
          // handleAddNode(_.cloneDeep(e), ip)
        }
      })
      // 补全所有的线
      const links = this.completeCurLevelLinks(this.nodesDefault, this.allLinks)
      const notDisplayedLinks = findDifferences(links, this.linksDefault, ['clientIp', 'serverIp'])
      notDisplayedLinks.forEach((i: link) => {
        this.linksDefault.push(i)
        // setTimeout(() => {
        //   handleAddLink(_.cloneDeep(i))
        // }, 0)
      })
      // 先判断终点是否有下一级，有下一级加上加号
      node.forEach((j: node) => {
        const nextSkipArr = this.allLinks.filter((k: link) => {
          return k.clientIp == j.ip
        })
        const flag = isSubset(nextSkipArr, this.linksDefault, ['clientIp', 'serverIp'])
        if (!flag) {
          this.nodesDefault = this.nodesDefault.map((warp: node) => {
            if (warp.ip == j.ip) {
              warp['isUnfold'] = mode == 'all' ? 2 : 1
            }
            return warp
          })
        }
      })
      addnNodeArr.forEach((n: node) => {
        handleAddNode(_.cloneDeep(n), ip)
        // 然后将这一级父节点的加号状态改为减号
        const relatedLinks = this.linksDefault.filter((link: link) => {
          return n.ip == link.serverIp
        })
        this.nodesDefault = this.nodesDefault.map((warp: node) => {
          const res = relatedLinks.find((relatedLinksItem) => {
            return warp.ip == relatedLinksItem.clientIp
          })
          if (warp.ip == res?.clientIp) {
            warp['isUnfold'] = 2
          }
          return warp
        })
      })

      notDisplayedLinks.forEach((i: link) => {
        setTimeout(() => {
          handleAddLink(_.cloneDeep(i))
        }, 10)
      })
      if (mode == 'all' && link.length > 0 && td.serverIp != ip) {
        setTimeout(() => {
          this.RecursionGetNextLevelDataByIP(td.serverIp, 'all', handleAddNode, handleAddLink)
        }, 10)
      }
    })
  }

  // 根据ip删除下一级数据
  public DelDataByIP(ip: string): { nodesDel: node[]; linksDel: link[] } {
    const visited = new Set<string>()
    const [descendants] = this.findDescendants(ip, visited)
    // 从links数组中移除相关的边
    // relatedLinks.forEach((link) => {
    //   this.linksDefault = this.linksDefault.filter((item) => {
    //     return item.clientIp !== link.clientIp && item.serverIp !== link.serverIp
    //   })
    // })

    // 从nodes数组中移除子节点
    descendants.forEach((node) => {
      this.nodesDefault = this.nodesDefault.filter((item) => {
        return item.ip != node.ip
      })
    })

    const array = this.completeCurLevelLinks(this.nodesDefault, this.allLinks)

    const notDisplayedLinks = findDifferences(array, this.linksDefault, ['clientIp', 'serverIp'])

    this.linksDefault = array

    return { nodesDel: descendants, linksDel: notDisplayedLinks }
  }

  /**
   * @description '根据ip修改标记颜色'
   * @params node: node
   */
  public async ChangeRemarkByNodeInfo(node: node) {
    if (!node.remark) {
      const { msg } = await DeleteRemarkApi({ ip: node.ip })
      gp.$baseMessage(msg, 'success', 'vab-hey-message-success')
    } else {
      const { msg } = await UpdateOrSaveRemarkApi({ ip: node.ip, remark: node.remark })
      gp.$baseMessage(msg, 'success', 'vab-hey-message-success')
    }
    this.allNodes = this.allNodes.map((td) => {
      if (node.ip == td.ip) {
        td.remark = node.remark
      }
      return td
    })
    this.nodesDefault = this.nodesDefault.map((td) => {
      if (node.ip == td.ip) {
        td.remark = node.remark
      }
      return td
    })
  }
  // 增加自定义数据
  public AddCustomizationData(ip: string) {}

  /**
   * @description '获取默认前三级的node和link'
   */
  public InitData(): { nodesDefault: node[]; linksDefault: link[] } {
    // 默认数据
    this.nodesDefault = []
    this.linksDefault = []
    // 找到根节点也就是一级节点
    const result = this.findRootNodesAndLoops(this.allLinks)
    const rootNode = result['rootNodes']
    // 标识+号的节点，排除掉根节点
    const isUnfoldNodes = result['isUnfoldNodes']
    // 获取2级节点
    const twoLinksMap: link[] = this.allLinks.filter((linkMap) => rootNode.includes(linkMap.clientIp))
    // 根据2级节点的serverIP获取3级节点(2级node)
    const threeNode: string[] = twoLinksMap.map(
      (linkMap) => linkMap['serverIp'] ?? '' // 使用空值合并运算符提供默认值
    )
    let threeLinksMap: link[] = this.allLinks.filter(
      (linkMap) => threeNode.includes(linkMap['clientIp'] ?? '') // 使用空值合并运算符提供默认值
    )
    // 然后计算第三级中所有的serverIp 是否有下一级，如果有的话则给他+号
    let threeFoldNodes = this.getThreeFoldNodes(threeLinksMap, this.allLinks)
    // 在准备赋值+号的节点列表中，去除掉根节点。如果存在于2级节点列表中，也要把它去除掉。
    threeFoldNodes = threeFoldNodes.filter((ip) => !rootNode.includes(ip) && !threeNode.includes(ip))
    // 将2级攻击线路与3级攻击线路汇总则得到最终的3级线路
    threeLinksMap.push(...twoLinksMap)
    // 将前3级进行去重，防止A-》B B-》A B-》C
    threeLinksMap = this.removeDuplicateMaps(threeLinksMap)
    // 将默认3级所有节点与nodesMap全部节点比对 取出默认3级相关的节点
    const threeNodesMap = this.getNodesMapNew(threeLinksMap, _.cloneDeep(this.allNodes), threeFoldNodes, 'default')
    // 补全当前level的link（可能是3指向2，3指向3，2指向...）
    threeLinksMap = this.completeCurLevelLinks(threeNodesMap, this.allLinks)
    this.linksDefault = threeLinksMap
    this.nodesDefault = threeNodesMap
    this.nodesDefault = this.nodesDefault.map((node: node) => {
      if (rootNode.includes(node.ip)) {
        node['level'] = 1
      } else if (threeNode.includes(node.ip)) {
        node['level'] = 2
      } else {
        node['level'] = 3
      }
      return node
    })
    return { nodesDefault: this.nodesDefault, linksDefault: this.linksDefault }
  }

  // 递归查找所有子节点及其直接相连的边
  private findDescendants(nodeId: string, visited: Set<string>): [node[]] {
    const descendants: node[] = []
    // const relatedLinks: link[] = []

    let index = 3
    this.nodesDefault.forEach((node: node) => {
      if (node.ip == nodeId) {
        index = node.level || 3
      }
    })
    // 查找与当前节点直接相连的边
    const directLinks = this.linksDefault.filter((link: link) => link.clientIp == nodeId)

    // 添加直接相连的边到结果集中
    // relatedLinks.push(...directLinks)

    // 标记当前节点为已访问
    visited.add(nodeId)

    // 对每个直接相连的边的目标节点进行递归查找
    directLinks.forEach((link) => {
      const targetNodeId = link.serverIp
      let _index = 3
      this.nodesDefault.forEach((node: node) => {
        if (node.ip == targetNodeId) {
          _index = node.level || 3
        }
      })

      if (!visited.has(targetNodeId) && _index > index) {
        const targetNode = this.nodesDefault.find((n) => n.ip === targetNodeId)!
        descendants.push(targetNode)
        const [childDescendants] = this.findDescendants(targetNodeId, visited)
        descendants.push(...childDescendants)
        // relatedLinks.push(...childLinks)
        // relatedLinks = this.completeCurLevelLinks(threeNodesMap, this.allLinks)
      }
    })

    return [descendants]
  }

  private completeCurLevelLinks(curNodes: node[], allLinks: link[]): link[] {
    const nodes: string[] = []
    curNodes.forEach((item: node) => {
      nodes.push(item.ip)
    })
    const links = allLinks.filter((item: link) => {
      if (nodes.includes(item.clientIp) && nodes.includes(item.serverIp)) {
        return item
      }
    })
    return links
  }

  private getNodesMapNew(linksMap: link[], nodesMap: node[], isUnfoldNodes: string[], status: string): node[] {
    let uniqueList: node[] = []

    if (status === 'all') {
      // 去重
      const idMap = new Map<string, node>()
      nodesMap.forEach((nodeMap) => {
        const id = nodeMap['id']
        if (id && !idMap.has(id)) {
          idMap.set(id, nodeMap)
        }
      })
      uniqueList = Array.from(idMap.values())
      // 标记展开节点
      uniqueList = uniqueList.map((map) => {
        if (isUnfoldNodes.includes(map['ip'])) {
          map['isUnfold'] = 1
        } else {
          map['isUnfold'] = 0
        }
        return map
      })
    } else {
      // 收集所有 clientIPs 和 serverIPs
      const allIPs = new Set<string>()
      const clientIPs = new Set(linksMap.map((linkMap) => linkMap['clientIp']))
      const serverIPs = new Set(linksMap.map((linkMap) => linkMap['serverIp']))
      clientIPs.forEach((ip) => allIPs.add(ip))
      serverIPs.forEach((ip) => allIPs.add(ip))

      // 筛选 nodesMap
      const nodesMapNew = nodesMap.filter((nodeMap) => allIPs.has(nodeMap['ip']))

      // 去重
      const idMap = new Map<string, node>()
      nodesMapNew.forEach((nodeMap) => {
        const id = nodeMap['id']
        if (id && !idMap.has(id)) {
          idMap.set(id, nodeMap)
        }
      })
      uniqueList = Array.from(idMap.values())
      // 标记展开节点
      uniqueList = uniqueList.map((map) => {
        if (isUnfoldNodes.includes(map.ip)) {
          map['isUnfold'] = 1
        } else {
          map['isUnfold'] = 0
        }
        return map
      })
    }

    return uniqueList
  }

  private distinctByKey<T>(keySelector: (item: T) => string): (item: T) => boolean {
    const seen = new Set<string>()
    return (item: T) => {
      const key = keySelector(item)
      if (seen.has(key)) {
        return false
      } else {
        seen.add(key)
        return true
      }
    }
  }

  private removeDuplicateMaps(threeLinksMap: link[]): link[] {
    // 使用 filter 方法和 distinctByKey 函数来过滤掉重复的 map 对象
    return threeLinksMap.filter(this.distinctByKey((map) => `${map['clientIp'] ?? ''}->${map['serverIp'] ?? ''}`))
  }
  // 获取第三级node
  private getThreeFoldNodes(threeMap: link[], linksMap: link[]): string[] {
    // 第一步：创建 serverIPs Set，包含 threeMap 中所有 linkMap 的 SERVER_IP_KEY 值
    const serverIPs: Set<string> = new Set(threeMap.map((linkMap) => linkMap['serverIp'] ?? ''))

    // 第二步：创建 clientIPs Set，包含 linksMap 中所有 linkMap 的 CLIENTIP_KEY 值
    const clientIPs: Set<string> = new Set(linksMap.map((linkMap) => linkMap['clientIp'] ?? ''))

    // 第三步：过滤 serverIPs，保留那些也存在于 clientIPs 中的 IP 地址
    const threeFoldNodes: string[] = Array.from(serverIPs).filter((ip) => clientIPs.has(ip))
    return threeFoldNodes
  }

  // 获取根节点
  private findRootNodesAndLoops(data: link[]): { [key: string]: string[] } {
    const result: { [key: string]: string[] } = {}
    // 保存所有clientIP
    const allClientIPs: Set<string> = new Set()

    // 遍历数据，构建allClientIPs
    for (const record of data) {
      if (record.clientIp) {
        allClientIPs.add(record.clientIp)
      }
    }

    // 获取根节点
    const rootNodes = getRootNodes(data)

    // 将Set转换为Array并存入结果
    result['rootNodes'] = rootNodes
    result['isUnfoldNodes'] = Array.from(allClientIPs).filter((ip) => !rootNodes.includes(ip))

    return result
  }
}

class Graph {
  V: number
  adj: number[][]
  constructor(V: number) {
    this.V = V
    // 创建 V 个空的邻接列表
    this.adj = new Array(V).fill(null).map(() => [])
  }
  addEdge(v: number, w: number): void {
    if (v >= 0 && v < this.V && w >= 0 && w < this.V) {
      this.adj[v].push(w)
    } else {
      console.warn(`Invalid edge: (${v}, ${w}) is out of bounds.`)
    }
  }

  getAdjList(v: number): number[] | undefined {
    if (v >= 0 && v < this.V) {
      return this.adj[v]
    } else {
      console.warn(`Invalid node: ${v} is out of bounds.`)
      return undefined
    }
  }

  getNumberOfVertices(): number {
    return this.V
  }

  getAdjacencyList(): number[][] {
    return this.adj
  }
}

class TarjanSCC {
  private V: number
  private adj: number[][]
  private index: number
  private indices: number[]
  private lowlink: number[]
  private onStack: boolean[]
  private stack: number[]
  private SCCs: number[][]

  constructor(graph: Graph) {
    this.V = graph.V
    this.adj = graph.adj
    this.index = 0
    this.indices = new Array(this.V).fill(-1)
    this.lowlink = new Array(this.V).fill(0)
    this.onStack = new Array(this.V).fill(false)
    this.stack = []
    this.SCCs = []

    // 遍历所有节点，若 indices[v] 为 -1，则调用 strongConnect
    for (let v = 0; v < this.V; v++) {
      if (this.indices[v] === -1) {
        this.strongConnect(v)
      }
    }
  }

  private strongConnect(v: number): void {
    this.indices[v] = this.index
    this.lowlink[v] = this.index
    this.index++
    this.stack.push(v)
    this.onStack[v] = true

    // 遍历 v 的邻接点
    for (const w of this.adj[v]) {
      if (this.indices[w] === -1) {
        // w 尚未访问过
        this.strongConnect(w)
        this.lowlink[v] = Math.min(this.lowlink[v], this.lowlink[w])
      } else if (this.onStack[w]) {
        // w 在栈中，说明在同一个 SCC 中
        this.lowlink[v] = Math.min(this.lowlink[v], this.indices[w])
      }
    }

    // 如果 v 是当前 SCC 的根节点
    if (this.lowlink[v] === this.indices[v]) {
      const scc: number[] = []
      while (true) {
        const w = this.stack.pop()
        // 理论上不会 undefined，但为了安全可加断言
        if (w === undefined) break

        this.onStack[w] = false
        scc.push(w)

        if (w === v) {
          break
        }
      }
      this.SCCs.push(scc)
    }
  }

  public getSCCs(): number[][] {
    return this.SCCs
  }
}

function getRootNodes(logItems: link[]): string[] {
  const ipToId = new Map<string, number>()
  const idToIp = new Map<number, string>()
  let currentId = 0
  const rootNodes: string[] = []
  const edges: [number, number][] = []

  try {
    for (const record of logItems) {
      const clientIP = record.clientIp
      const serverIP = record.serverIp

      if (clientIP && serverIP) {
        if (!ipToId.has(clientIP)) {
          ipToId.set(clientIP, currentId)
          idToIp.set(currentId, clientIP)
          currentId++
        }
        if (!ipToId.has(serverIP)) {
          ipToId.set(serverIP, currentId)
          idToIp.set(currentId, serverIP)
          currentId++
        }
        const clientId = ipToId.get(clientIP)!
        const serverId = ipToId.get(serverIP)!
        edges.push([clientId, serverId])
      }
    }

    const totalNodes = currentId
    const graph = new Graph(totalNodes)

    for (const [v, w] of edges) {
      graph.addEdge(v, w)
    }

    const tarjan = new TarjanSCC(graph)
    const SCCs = tarjan.getSCCs()
    const nodeToSCC = new Map<number, number>()
    let sccId = 1

    for (const scc of SCCs) {
      for (const node of scc) {
        nodeToSCC.set(node, sccId)
      }
      sccId++
    }

    const numSuperNodes = SCCs.length
    const superAdj: Map<number, number[]> = new Map()

    for (let i = 1; i <= numSuperNodes; i++) {
      superAdj.set(i, [])
    }

    for (let v = 0; v < totalNodes; v++) {
      const adjList = graph.getAdjList(v)
      if (adjList) {
        for (const w of adjList) {
          const sccV = nodeToSCC.get(v)!
          const sccW = nodeToSCC.get(w)!
          if (sccV !== sccW) {
            superAdj.get(sccV)?.push(sccW)
          }
        }
      }
    }

    const allSuperNodes = new Set([...Array(numSuperNodes).keys()].map((x) => x + 1))
    const nonSourceSuperNodes = new Set<number>()

    for (let i = 1; i <= numSuperNodes; i++) {
      for (const adjNode of superAdj.get(i)!) {
        nonSourceSuperNodes.add(adjNode)
      }
    }

    const sourceSuperNodes = new Set([...allSuperNodes].filter((x) => !nonSourceSuperNodes.has(x)))
    const inDegree = new Map<number, number>()
    const outDegree = new Map<number, number>()

    for (const [client, server] of edges) {
      outDegree.set(client, (outDegree.get(client) || 0) + 1)
      inDegree.set(server, (inDegree.get(server) || 0) + 1)
    }

    for (const superNode of sourceSuperNodes) {
      const scc = SCCs[superNode - 1]
      for (const node of scc) {
        const outD = outDegree.get(node) || 0
        const inD = inDegree.get(node) || 0
        rootNodes.push(idToIp.get(node)!)
      }
    }
  } catch (e) {
    console.error('计算根节点异常', e)
  }

  return rootNodes
}

/**
 * @description '根据标识判断两个对象数组的不同部分'
 */
function findDifferences<T extends Record<string, any>, K extends keyof T>(arr1: T[], arr2: T[], keys: K[]): T[] {
  // 创建一个映射以存储arr1中的对象（基于多个keys）
  const map = new Map<string, T>()
  arr1.forEach((item) => {
    // 使用keys创建一个唯一的复合键
    const compositeKey = keys.map((key) => item[key]).join('|')
    map.set(compositeKey, item)
  })

  // 找到arr2中不在arr1中的对象
  const differences = arr2.filter((item) => {
    const compositeKey = keys.map((key) => item[key]).join('|')
    return !map.has(compositeKey)
  })

  // 找到arr1中不在arr2中的对象
  arr1
    .filter((item) => {
      const compositeKey = keys.map((key) => item[key]).join('|')
      return !arr2.some((other) => keys.every((key) => other[key] === item[key]))
    })
    .forEach((item) => differences.push(item))

  return differences
}

/**
 * @description '根据标识判断一个对象数组是不是另一个对象数组的子集'
 */

// 定义一个泛型函数，T是对象类型，K是键名数组
function isSubset<T extends Record<string, any>, K extends keyof T>(
  subset: T[],
  superset: T[],
  keys: readonly K[]
): boolean {
  // 创建映射以存储superset中的对象（基于多个keys）
  const supersetMap = new Map<string, T>()
  superset.forEach((item) => {
    const compositeKey = keys.map((key) => item[key]).join('|')
    supersetMap.set(compositeKey, item)
  })

  // 检查subset中的每个对象是否都在superset中
  return subset.every((item) => {
    const compositeKey = keys.map((key) => item[key]).join('|')
    return supersetMap.has(compositeKey)
  })
}
