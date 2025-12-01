# Mermaid 图表

VitePress 支持使用 [Mermaid](https://mermaid.js.org/) 绘制各种图表，包括流程图、时序图、甘特图等。

## 流程图

```mermaid
flowchart LR
  A[开始] --> B{判断}
  B -->|是| C[执行]
  B -->|否| D[跳过]
  C --> E[结束]
  D --> E
```

## 时序图

```mermaid
sequenceDiagram
    participant C as 客户端
    participant S as 服务器
    participant D as 数据库

    C->>S: 发送请求
    S->>D: 查询数据
    D-->>S: 返回结果
    S-->>C: 响应数据
```

## 甘特图

```mermaid
gantt
    title 项目开发计划
    dateFormat YYYY-MM-DD
    section 设计阶段
        需求分析     :a1, 2024-01-01, 7d
        原型设计     :a2, after a1, 5d
    section 开发阶段
        前端开发     :b1, after a2, 14d
        后端开发     :b2, after a2, 14d
    section 测试阶段
        集成测试     :c1, after b1, 7d
        上线部署     :c2, after c1, 3d
```

## 类图

```mermaid
classDiagram
    class Animal {
        +String name
        +int age
        +makeSound()
    }
    class Dog {
        +String breed
        +bark()
    }
    class Cat {
        +String color
        +meow()
    }
    Animal <|-- Dog
    Animal <|-- Cat
```

## 状态图

```mermaid
stateDiagram-v2
    [*] --> 待处理
    待处理 --> 处理中: 开始处理
    处理中 --> 已完成: 处理成功
    处理中 --> 失败: 处理失败
    失败 --> 待处理: 重试
    已完成 --> [*]
```

## 饼图

```mermaid
pie title 技术栈占比
    "Vue.js" : 40
    "TypeScript" : 30
    "Node.js" : 20
    "其他" : 10
```

## 更多信息

查看 [Mermaid 官方文档](https://mermaid.js.org/intro/) 了解更多图表类型和语法。
