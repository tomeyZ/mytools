<template>
  <!-- 自绘下拉。
       为什么不用原生 <select>：展开列表由 WebView2 自己绘制，圆角 / 行高 / hover / 选中态
       一个都改不了（Windows 上 <option> 只认 background 和 color），和整页的令牌体系对不上。
       这套视觉与 TimeTool 的时区下拉同源：同样的触发器、同样的 caret 旋转、同样的弹出面板。
       交互：↑↓ 移动、Enter / 空格选中、Esc 关闭、Tab 移走焦点并关闭、点击面板外关闭。 -->
  <div ref="wrap" class="sel" :class="{ 'is-open': open, up }">
    <button
        type="button"
        class="sel-trigger"
        :style="{ minWidth: minWidth }"
        role="combobox"
        aria-haspopup="listbox"
        :aria-expanded="open ? 'true' : 'false'"
        :aria-activedescendant="open ? uid + '-opt-' + active : null"
        @click="toggle"
        @keydown="onKeydown"
    >
      <span class="sel-value">{{ currentLabel }}</span>
      <svg class="sel-caret" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="6 9 12 15 18 9"/>
      </svg>
    </button>

    <!-- after-leave 不能省：见 close() 里的说明 -->
    <transition name="sel-pop" @after-leave="onAfterLeave">
      <ul v-if="open" ref="menu" class="sel-menu" role="listbox">
        <li
            v-for="(o, i) in options"
            :id="uid + '-opt-' + i"
            :key="o.value"
            class="sel-item"
            :class="{ act: i === active, on: o.value === modelValue }"
            role="option"
            :aria-selected="o.value === modelValue ? 'true' : 'false'"
            @mousedown.prevent
            @mouseenter="active = i"
            @click="pick(o)"
        >
          <span class="sel-text">{{ o.label }}</span>
          <svg v-if="o.value === modelValue" class="sel-check" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="20 6 9 17 4 12"/>
          </svg>
        </li>
      </ul>
    </transition>
  </div>
</template>

<script>
// 面板离触发器的间隙，必须与 CSS 里 top / bottom 的 calc(100% + 6px) 保持一致
const GAP = 6;

// 面板高度上限，对应 CSS 的 max-height
const MAX_H = 320;

/* 找最近的可滚动祖先。
   为什么不直接判断 window.innerHeight：本应用 body / #app 都是 overflow:hidden，
   真正滚动的是 MainLayout 里的 .content（flex:1 + overflow:auto），它比视口少一整条顶栏。
   面板「放不放得下」取决于它在 .content 里溢出多少，不是视口。 */
function scrollHost(el) {
  let p = el && el.parentElement;
  while (p && p !== document.body && p !== document.documentElement) {
    const oy = getComputedStyle(p).overflowY;
    if (oy === 'auto' || oy === 'scroll' || oy === 'overlay') return p;
    p = p.parentElement;
  }
  return null;
}

export default {
  name: 'AppSelect',
  props: {
    // 允许 Number：RSA 的「密钥长度」是数字值（1024/2048/4096），
    // 要原样传给 Go 后端，不能在中间转成字符串
    modelValue: { type: [String, Number], required: true },
    // [{ value, label }]
    options: { type: Array, required: true },
    // 固定最小宽度：选项长短不一时触发器的宽度不跟着跳
    minWidth: { type: String, default: '200px' }
  },
  emits: ['update:modelValue'],
  data() {
    return {
      open: false,
      // 面板翻到触发器上方（下方放不下时）。见 place()
      up: false,
      active: 0,
      // 选项 id 的唯一前缀。同一页面可能同时挂多个 AppSelect（RSA 页就有两个），
      // 写死 'sel-opt-0' 会产出重复 id，把 aria-activedescendant 指到别人家去
      uid: 'sel-' + Math.random().toString(36).slice(2, 8)
    };
  },
  computed: {
    currentLabel() {
      const hit = this.options.find(o => o.value === this.modelValue);
      // 兜底：值不在列表里时至少把原值显示出来，不要留空
      return hit ? hit.label : this.modelValue;
    }
  },
  mounted() {
    document.addEventListener('mousedown', this.onDocMousedown);
  },
  beforeUnmount() {
    document.removeEventListener('mousedown', this.onDocMousedown);
  },
  methods: {
    toggle() {
      this.open ? this.close() : this.openMenu();
    },
    openMenu() {
      if (this.open) return;
      this.open = true;
      // 打开时把高亮定到当前选中项，避免每次都要从头按 ↓
      const i = this.options.findIndex(o => o.value === this.modelValue);
      this.active = i < 0 ? 0 : i;
      // 先判方向再滚：place() 可能把面板翻到上方，可用高度跟着变
      this.$nextTick(() => {
        this.place();
        this.scrollActiveIntoView();
      });
    },

    /* ⚠️ up 绝对不能在这里复位。
       面板是 v-if + transition，open=false 之后它还要留在 DOM 里跑 140ms 退场动画。
       这一刻把它复位成向下，面板会当场跳回触发器下方 —— 而下方正是当初判定「放不下」
       的那一侧 → 撑出滚动条 → 内容区被挤窄 10px → 旁边的二维码区整个左移一下。
       表现就是「面板从下面一闪而过，右边闪出一条空白」。
       必须等动画真正结束（after-leave）再复位，让面板在自己那一侧原地淡出。 */
    close() {
      this.open = false;
    },

    onAfterLeave() {
      this.up = false;
    },

    /* 面板默认向下展开。贴着容器底边时（默认 1024x768 下生成页的工具条就在那儿）
       它会溢出 .content → 冒出滚动条 → 内容区被挤窄 10px，看起来就是「点一下右边晃一下」。
       所以下方放不下、且上方更宽裕时翻到上面去。 */
    place() {
      const menu = this.$refs.menu;
      const wrap = this.$refs.wrap;
      if (!menu || !wrap) return;

      const host = scrollHost(wrap);
      const bound = host
        ? host.getBoundingClientRect()
        : { top: 0, bottom: window.innerHeight };
      const r = wrap.getBoundingClientRect();
      const below = bound.bottom - r.bottom - GAP;
      const above = r.top - bound.top - GAP;
      const need = menu.offsetHeight;

      this.up = below < need && above > below;

      // 翻上去后上方也未必够（窗口很矮时），把面板收到真正放得下的高度，
      // 剩下的交给面板自己的 overflow-y
      menu.style.maxHeight = this.up
        ? Math.max(120, Math.min(MAX_H, above)) + 'px'
        : '';
    },

    pick(o) {
      this.$emit('update:modelValue', o.value);
      this.close();
    },

    scrollActiveIntoView() {
      const menu = this.$refs.menu;
      if (!menu) return;
      const el = menu.children[this.active];
      if (!el) return;
      // 手算 scrollTop，不用 el.scrollIntoView()：后者会把所有可滚祖先一起滚，
      // 菜单贴着 .content 底边打开时会把整页往上带一下
      const top = el.getBoundingClientRect().top - menu.getBoundingClientRect().top + menu.scrollTop;
      const bottom = top + el.offsetHeight;
      const pad = 5;   // 面板自身的 padding，让高亮项不贴边
      if (top < menu.scrollTop + pad) {
        menu.scrollTop = top - pad;
      } else if (bottom > menu.scrollTop + menu.clientHeight - pad) {
        menu.scrollTop = bottom - menu.clientHeight + pad;
      }
    },
    onKeydown(e) {
      const k = e.key;
      if (!this.open) {
        // 关闭状态下这些键直接展开（原生 select 的习惯）
        if (k === 'ArrowDown' || k === 'ArrowUp' || k === 'Enter' || k === ' ') {
          e.preventDefault();
          this.openMenu();
        }
        return;
      }
      if (k === 'Escape') {
        e.preventDefault();
        this.close();
      } else if (k === 'Tab') {
        this.close();                 // 不拦 Tab，让它正常移走焦点
      } else if (k === 'ArrowDown') {
        e.preventDefault();
        this.active = Math.min(this.active + 1, this.options.length - 1);
        this.scrollActiveIntoView();
      } else if (k === 'ArrowUp') {
        e.preventDefault();
        this.active = Math.max(this.active - 1, 0);
        this.scrollActiveIntoView();
      } else if (k === 'Enter' || k === ' ') {
        e.preventDefault();
        const o = this.options[this.active];
        if (o) this.pick(o);
      }
    },
    onDocMousedown(e) {
      if (!this.open) return;
      const wrap = this.$refs.wrap;
      if (wrap && !wrap.contains(e.target)) {
        this.close();
      }
    }
  }
};
</script>

<style scoped>
.sel {
  position: relative;
}

/* 底色用 --surface-2（和 .inp 输入框一致）：它是表单里的控件，
   跟旁边的输入框同色才协调；浮层面板才是 --surface */
.sel-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 36px;
  padding: 0 10px 0 12px;
  font-family: inherit;
  font-size: 13px;
  color: var(--text);
  background: var(--surface-2);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: border-color var(--tr), box-shadow var(--tr);
}

.sel-trigger:hover {
  border-color: var(--border-strong);
}

.sel-trigger:focus-visible,
.sel.is-open .sel-trigger {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--ring);
  outline: none;
}

.sel-value {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-align: left;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.sel-caret {
  flex-shrink: 0;
  color: var(--text-3);
  transition: transform var(--tr);
}

.sel.is-open .sel-caret {
  transform: rotate(180deg);
}

/* 面板宽度严格等于触发器：left/right 双约束会把 padding + border 一起算进去，
   不能用 min-width:100%（那是内容盒的下限，会额外撑出 12px 再整体左歪） */
.sel-menu {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  z-index: 60;
  max-height: 320px;
  overflow-y: auto;
  margin: 0;
  padding: 5px;
  list-style: none;
  background: var(--surface);
  /* 用 --border-strong 而不是 --border：深色主题下 --shadow-lg 是 none，
     面板与卡片又同为 --surface，只剩这条边框能把浮层从背景里勾出来 */
  border: 1px solid var(--border-strong);
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-lg);
}

/* 下方放不下时翻到触发器上方。
   判据在 place() 里按滚动容器（.content）的实际可用高度算，不是按视口 */
.sel.up .sel-menu {
  top: auto;
  bottom: calc(100% + 6px);
}

/* 两列定宽：文字吸左、勾选占固定列。
   第二列即使没有勾也保留宽度，切换选项时行内元素不会左右跳 */
.sel-item {
  display: grid;
  grid-template-columns: 1fr 16px;
  align-items: center;
  gap: 8px;
  height: 34px;
  padding: 0 9px;
  border-radius: var(--radius-xs);
  font-size: 13px;
  color: var(--text-2);
  cursor: pointer;
  transition: background var(--tr), color var(--tr);
}

/* 键盘 / 鼠标经过的高亮 */
.sel-item.act {
  background: var(--surface-2);
  color: var(--text);
}

/* 当前生效项。写在 .act 之后，两者同时命中时以它为准 */
.sel-item.on {
  background: var(--accent-soft);
  color: var(--accent-strong);
  font-weight: 600;
}

.sel-text {
  min-width: 0;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.sel-check {
  color: var(--accent-strong);
}

.sel-pop-enter-active,
.sel-pop-leave-active {
  transition: opacity .14s ease, transform .14s ease;
}

/* 退场那 140ms 元素还占着位置，别再吃点击（可能落在下方的输入框上） */
.sel-pop-leave-active {
  pointer-events: none;
}

.sel-pop-enter-from,
.sel-pop-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

/* 动画方向跟着翻转，否则上弹面板会从下方「穿」过触发器 */
.sel.up .sel-pop-enter-from,
.sel.up .sel-pop-leave-to {
  transform: translateY(4px);
}
</style>
