# WindNav 图标重新设计方案

> 当前图标问题：罗盘太通用、风线太抽象、小尺寸识别度低、与"搜索优先的轻量导航启动器"定位不匹配。

---

## 方案一：风车 (Windmill)

**设计理念：** 以"风"为直接意象，风车叶片转动代表活力与指引。

**视觉感受：** 活力、清新、直接

### icon.svg (48×48)

```svg
<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 48 48" fill="none">
  <defs>
    <linearGradient id="brandGradient" x1="0" y1="0" x2="48" y2="48" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="46" height="46" rx="11" fill="url(#brandGradient)"/>
  <circle cx="24" cy="24" r="4" fill="#ffffff" opacity="0.95"/>
  <path d="M24 24 L24 10 C28 12 30 18 24 24Z" fill="#ffffff" opacity="0.95"/>
  <path d="M24 24 L36 30 C38 26 36 22 24 24Z" fill="#ffffff" opacity="0.8"/>
  <path d="M24 24 L12 30 C10 26 12 22 24 24Z" fill="#ffffff" opacity="0.85"/>
</svg>
```

### favicon.svg (32×32)

```svg
<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32" fill="none">
  <defs>
    <linearGradient id="g" x1="0" y1="0" x2="32" y2="32" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="30" height="30" rx="7" fill="url(#g)"/>
  <circle cx="16" cy="16" r="3" fill="#ffffff" opacity="0.95"/>
  <path d="M16 16 L16 6 C19 8 21 12 16 16Z" fill="#ffffff" opacity="0.95"/>
  <path d="M16 16 L25 21 C27 18 25 14 16 16Z" fill="#ffffff" opacity="0.8"/>
  <path d="M16 16 L7 21 C5 18 7 14 16 16Z" fill="#ffffff" opacity="0.85"/>
</svg>
```

### icons.svg 中的 symbol

```svg
<symbol id="windnav-icon" viewBox="0 0 48 48">
  <defs>
    <linearGradient id="brandGradient" x1="0" y1="0" x2="48" y2="48" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="46" height="46" rx="11" fill="url(#brandGradient)"/>
  <circle cx="24" cy="24" r="4" fill="#ffffff" opacity="0.95"/>
  <path d="M24 24 L24 10 C28 12 30 18 24 24Z" fill="#ffffff" opacity="0.95"/>
  <path d="M24 24 L36 30 C38 26 36 22 24 24Z" fill="#ffffff" opacity="0.8"/>
  <path d="M24 24 L12 30 C10 26 12 22 24 24Z" fill="#ffffff" opacity="0.85"/>
</symbol>
```

---

## 方案二：纸飞机 (Paper Plane)

**设计理念：** 纸飞机代表"出发/导航/轻盈"，与"微风"概念天然契合，现代感强。

**视觉感受：** 轻盈、自由、现代

### icon.svg (48×48)

```svg
<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 48 48" fill="none">
  <defs>
    <linearGradient id="brandGradient" x1="0" y1="0" x2="48" y2="48" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="46" height="46" rx="11" fill="url(#brandGradient)"/>
  <path d="M38 14 L14 24 L22 26 L26 34 L38 14Z" fill="#ffffff" opacity="0.95"/>
  <path d="M14 24 L26 22 L38 14" stroke="#22d3ee" stroke-width="1.2" stroke-linecap="round" fill="none" opacity="0.3"/>
  <path d="M8 36 Q20 38 34 26" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" fill="none" opacity="0.45"/>
</svg>
```

### favicon.svg (32×32)

```svg
<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32" fill="none">
  <defs>
    <linearGradient id="g" x1="0" y1="0" x2="32" y2="32" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="30" height="30" rx="7" fill="url(#g)"/>
  <path d="M26 10 L10 18 L16 20 L20 26 L26 10Z" fill="#ffffff" opacity="0.95"/>
  <path d="M4 26 Q14 27 24 19" stroke="#ffffff" stroke-width="1.2" stroke-linecap="round" fill="none" opacity="0.4"/>
</svg>
```

### icons.svg 中的 symbol

```svg
<symbol id="windnav-icon" viewBox="0 0 48 48">
  <defs>
    <linearGradient id="brandGradient" x1="0" y1="0" x2="48" y2="48" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="46" height="46" rx="11" fill="url(#brandGradient)"/>
  <path d="M38 14 L14 24 L22 26 L26 34 L38 14Z" fill="#ffffff" opacity="0.95"/>
  <path d="M14 24 L26 22 L38 14" stroke="#22d3ee" stroke-width="1.2" stroke-linecap="round" fill="none" opacity="0.3"/>
  <path d="M8 36 Q20 38 34 26" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" fill="none" opacity="0.45"/>
</symbol>
```

---

## 方案三：双 W / 风纹 (Wind Mark)

**设计理念：** 从品牌名 "WindNav" 提取首字母 W，以风的流动曲线构成双 W 重叠，形成独特的品牌符号。

**视觉感受：** 抽象、优雅、品牌感强

### icon.svg (48×48)

```svg
<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 48 48" fill="none">
  <defs>
    <linearGradient id="brandGradient" x1="0" y1="0" x2="48" y2="48" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="46" height="46" rx="11" fill="url(#brandGradient)"/>
  <!-- 前景 W - 实心白色 -->
  <path d="M10 32 L18 16 L24 26 L30 16 L38 32" stroke="#ffffff" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round" fill="none" opacity="0.95"/>
  <!-- 背景 W - 流动风感，错位叠加 -->
  <path d="M6 34 Q14 14 20 26 Q24 32 28 26 Q34 14 42 34" stroke="#ffffff" stroke-width="1.8" stroke-linecap="round" fill="none" opacity="0.35"/>
  <path d="M14 34 Q20 20 24 28 Q28 20 34 34" stroke="#ffffff" stroke-width="1.2" stroke-linecap="round" fill="none" opacity="0.25"/>
</svg>
```

### favicon.svg (32×32)

```svg
<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32" fill="none">
  <defs>
    <linearGradient id="g" x1="0" y1="0" x2="32" y2="32" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="30" height="30" rx="7" fill="url(#g)"/>
  <path d="M6 22 L12 10 L16 18 L20 10 L26 22" stroke="#ffffff" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" fill="none" opacity="0.95"/>
  <path d="M4 24 Q10 12 14 18 Q16 22 18 18 Q22 12 28 24" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" fill="none" opacity="0.3"/>
</svg>
```

### icons.svg 中的 symbol

```svg
<symbol id="windnav-icon" viewBox="0 0 48 48">
  <defs>
    <linearGradient id="brandGradient" x1="0" y1="0" x2="48" y2="48" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="46" height="46" rx="11" fill="url(#brandGradient)"/>
  <path d="M10 32 L18 16 L24 26 L30 16 L38 32" stroke="#ffffff" stroke-width="3.5" stroke-linecap="round" stroke-linejoin="round" fill="none" opacity="0.95"/>
  <path d="M6 34 Q14 14 20 26 Q24 32 28 26 Q34 14 42 34" stroke="#ffffff" stroke-width="1.8" stroke-linecap="round" fill="none" opacity="0.35"/>
  <path d="M14 34 Q20 20 24 28 Q28 20 34 34" stroke="#ffffff" stroke-width="1.2" stroke-linecap="round" fill="none" opacity="0.25"/>
</symbol>
```

---

## 方案四：堆叠导航格 (Grid Launchpad)

**设计理念：** 回归产品本质——导航启动器。用堆叠的网格方块代表"应用启动/导航"。

**视觉感受：** 功能明确、科技感、产品本质

### icon.svg (48×48)

```svg
<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 48 48" fill="none">
  <defs>
    <linearGradient id="brandGradient" x1="0" y1="0" x2="48" y2="48" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="46" height="46" rx="11" fill="url(#brandGradient)"/>
  <!-- 背景层方块（半透明，错位） -->
  <rect x="10" y="10" width="14" height="14" rx="4" fill="#ffffff" opacity="0.25"/>
  <rect x="24" y="24" width="14" height="14" rx="4" fill="#ffffff" opacity="0.25"/>
  <!-- 前景层方块（白色实心，堆叠感） -->
  <rect x="14" y="14" width="14" height="14" rx="4" fill="#ffffff" opacity="0.95"/>
  <rect x="20" y="20" width="14" height="14" rx="4" fill="#ffffff" opacity="0.85"/>
</svg>
```

### favicon.svg (32×32)

```svg
<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32" fill="none">
  <defs>
    <linearGradient id="g" x1="0" y1="0" x2="32" y2="32" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="30" height="30" rx="7" fill="url(#g)"/>
  <rect x="7" y="7" width="9" height="9" rx="2.5" fill="#ffffff" opacity="0.25"/>
  <rect x="16" y="16" width="9" height="9" rx="2.5" fill="#ffffff" opacity="0.25"/>
  <rect x="10" y="10" width="9" height="9" rx="2.5" fill="#ffffff" opacity="0.95"/>
  <rect x="14" y="14" width="9" height="9" rx="2.5" fill="#ffffff" opacity="0.85"/>
</svg>
```

### icons.svg 中的 symbol

```svg
<symbol id="windnav-icon" viewBox="0 0 48 48">
  <defs>
    <linearGradient id="brandGradient" x1="0" y1="0" x2="48" y2="48" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="46" height="46" rx="11" fill="url(#brandGradient)"/>
  <rect x="10" y="10" width="14" height="14" rx="4" fill="#ffffff" opacity="0.25"/>
  <rect x="24" y="24" width="14" height="14" rx="4" fill="#ffffff" opacity="0.25"/>
  <rect x="14" y="14" width="14" height="14" rx="4" fill="#ffffff" opacity="0.95"/>
  <rect x="20" y="20" width="14" height="14" rx="4" fill="#ffffff" opacity="0.85"/>
</symbol>
```

---

## 方案五：风向标 (Wind Vane)

**设计理念：** 在罗盘基础上做减法，保留"方向"概念，但用更简洁的风向标取代。

**视觉感受：** 简洁、明确、经典

### icon.svg (48×48)

```svg
<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 48 48" fill="none">
  <defs>
    <linearGradient id="brandGradient" x1="0" y1="0" x2="48" y2="48" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="46" height="46" rx="11" fill="url(#brandGradient)"/>
  <!-- 竖杆 -->
  <rect x="23" y="12" width="2" height="24" rx="1" fill="#ffffff" opacity="0.9"/>
  <!-- 箭头（尖端朝右上） -->
  <path d="M23 14 L16 22 L20 22 L20 24 L26 24 L26 22 L30 22 Z" fill="#ffffff" opacity="0.95"/>
  <!-- 尾翼（左下十字） -->
  <path d="M24 30 L24 34" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" fill="none" opacity="0.7"/>
  <path d="M20 32 L28 32" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" fill="none" opacity="0.7"/>
  <!-- 底座圆 -->
  <circle cx="24" cy="36" r="2" fill="#ffffff" opacity="0.6"/>
</svg>
```

### favicon.svg (32×32)

```svg
<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32" fill="none">
  <defs>
    <linearGradient id="g" x1="0" y1="0" x2="32" y2="32" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="30" height="30" rx="7" fill="url(#g)"/>
  <rect x="15" y="8" width="2" height="16" rx="1" fill="#ffffff" opacity="0.9"/>
  <path d="M15 10 L10 16 L13 16 L13 18 L19 18 L19 16 L22 16 Z" fill="#ffffff" opacity="0.95"/>
  <path d="M16 21 L16 24" stroke="#ffffff" stroke-width="1.2" stroke-linecap="round" fill="none" opacity="0.7"/>
  <path d="M13 22.5 L19 22.5" stroke="#ffffff" stroke-width="1.2" stroke-linecap="round" fill="none" opacity="0.7"/>
</svg>
```

### icons.svg 中的 symbol

```svg
<symbol id="windnav-icon" viewBox="0 0 48 48">
  <defs>
    <linearGradient id="brandGradient" x1="0" y1="0" x2="48" y2="48" gradientUnits="userSpaceOnUse">
      <stop offset="0%" stop-color="#22d3ee"/>
      <stop offset="100%" stop-color="#3b82f6"/>
    </linearGradient>
  </defs>
  <rect x="1" y="1" width="46" height="46" rx="11" fill="url(#brandGradient)"/>
  <rect x="23" y="12" width="2" height="24" rx="1" fill="#ffffff" opacity="0.9"/>
  <path d="M23 14 L16 22 L20 22 L20 24 L26 24 L26 22 L30 22 Z" fill="#ffffff" opacity="0.95"/>
  <path d="M24 30 L24 34" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" fill="none" opacity="0.7"/>
  <path d="M20 32 L28 32" stroke="#ffffff" stroke-width="1.5" stroke-linecap="round" fill="none" opacity="0.7"/>
  <circle cx="24" cy="36" r="2" fill="#ffffff" opacity="0.6"/>
</symbol>
```

---

## 总结对比

| 方案 | 核心意象 | 品牌关联度 | 小尺寸识别度 | 现代感 | 推荐度 |
|------|---------|-----------|-------------|-------|-------|
| ① 风车 | 风、转动 | ★★★★ | ★★★★★ | ★★★★ | ⭐ 推荐 |
| ② 纸飞机 | 出发、轻盈 | ★★★★ | ★★★★ | ★★★★★ | ⭐ 推荐 |
| ③ 双 W 风纹 | 品牌首字母 | ★★★★★ | ★★★ | ★★★★★ | 精选 |
| ④ 堆叠导航格 | 启动器、网格 | ★★★ | ★★★★ | ★★★★ | 备选 |
| ⑤ 风向标 | 方向、导航 | ★★★ | ★★★★ | ★★★ | 备选 |

## 替换说明

选定方案后，需要替换以下 3 个文件：
1. `frontend/public/icon.svg` - 主图标
2. `frontend/public/favicon.svg` - 浏览器标签图标
3. `frontend/public/icons.svg` 中的 `#windnav-icon` `<symbol>` 定义