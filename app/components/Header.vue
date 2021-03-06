<template>
  <header 
    v-if="isMounted" 
    class="shadow"
  >
    <nav
      v-if="wideMode"
      class="top-nav"
    >
      <nuxt-link to="/">
        <div class="logo">
          <img
            class="icon"
            src="/img/logo.png"
          >
          明鳴的树洞
        </div>
      </nuxt-link>

      <div :class="{narrowTopNav: moderateWideMode}">
        <nuxt-link
          v-for="item in items.slice(1)"
          :key="item.name"
          :to="item.link"
          class="nav-item"
          @click.native="navigate(item.name)"
        >
          {{ item.name }}
        </nuxt-link>
      </div>
    </nav>

    <nav
      v-if="!wideMode"
      class="drawer-nav"
    >
      <div ref="drawerNavBar">
        <nuxt-link 
          to=""
          class="toggler container"
          @click.native="toggleDrawer"
        >
          <i
            v-if="!drawerShown"
            class="fas fa-bars"
          />
          <i
            v-else
            class="fas fa-times"
          />
        </nuxt-link><span v-if="!drawerShown">{{ currentRouterName }}</span>
      </div>
      <div
        :class="[drawerShown && !wideMode ? 'drawer-show' : 'drawer-hide']"
        :style="{marginTop: drawerNavBarHeight - 1 + 'px'}"
        class="drawer"
      >
        <div class="container">
          <nuxt-link
            v-for="item in items"
            :key="item.name"
            :to="item.link"
            :class="{'drawer-item-active': currentRouterName === item.name}"
            class="link"
            @click.native="navigate(item.name)"
          >
            <p>
              <i :class="item.icon" />
              {{ item.name }}
            </p>
          </nuxt-link>
        </div>
      </div>
    </nav>
  </header>
</template>

<script>
export default {
  data() {
    return {
      currentRouterName: '',
      drawerNavBarHeight: 0,
      drawerShown: false,
      wideMode: false,
      moderateWideMode: false,
      isMounted: false,
      items: [
        {
          link: '/',
          name: 'Home',
          icon: 'fas fa-home',
          wideMode: false
        },
        {
          link: '/articles',
          name: 'Articles',
          icon: 'fas fa-pencil-alt',
          wideMode: true
        },
        {
          link: '/photos',
          name: 'Photos',
          icon: 'fas fa-camera-retro',
          wideMode: true
        },
        {
          link: '/projects',
          name: 'Projects',
          icon: 'fas fa-project-diagram',
          wideMode: true
        },
        {
          link: '/reviews',
          name: 'Reviews',
          icon: 'fas fa-star-half-alt',
          wideMode: true
        }
      ]
    }
  },
  created() {
    if (process.browser) {
      window.addEventListener('resize', this.calcDrawerNavBarHeight)
      window.addEventListener('resize', this.calcIsWideMode)
      window.addEventListener('resize', this.calcIsModerateWideMode)
    }
  },
  mounted() {
    this.calcDrawerNavBarHeight()
    this.calcIsWideMode()
    this.determineCurrentRouterName()
    this.isMounted = true
  },
  beforeDestroy() {
    if (process.browser) {
      window.removeEventListener('resize', this.calcDrawerNavBarHeight)
      window.removeEventListener('resize', this.calcIsWideMode)
      window.removeEventListener('resize', this.calcIsModerateWideMode)
    }
  },
  methods: {
    calcDrawerNavBarHeight() {
      if (this.$refs.drawerNavBar) {
        this.drawerNavBarHeight = this.$refs.drawerNavBar.clientHeight
      }
    },
    calcIsWideMode() {
      this.wideMode = window.innerWidth >= 680 // $max_width_device + 1
    },
    calcIsModerateWideMode() {
      this.moderateWideMode = window.innerWidth <= 960
    },
    toggleDrawer() {
      this.drawerShown = !this.drawerShown
    },
    determineCurrentRouterName() {
      this.currentRouterName = ''
      const path = window.location.pathname
      if (path === '/') {
        this.currentRouterName = 'Home'
        return
      }
      this.items.forEach(item => {
        if (path === item.link) {
          this.currentRouterName = item.name
        }
      })
    },
    navigate(name) {
      this.currentRouterName = name
      if (!this.wideMode) {
        this.toggleDrawer()
      }
    }
  }
}
</script>

<style lang="scss" scoped>
@import '~/assets/scss/styles.scss';
.narrowTopNav {
  display: inline-block;
  width: 310px;
}
.logo {
  cursor: pointer;
  float: left;
  outline: none;
  user-select: none;
  color: $color_text;
  padding: 0 16px;
  .icon {
    vertical-align: middle;
    width: 64px;
    height: 64px;
  }
}
header {
  background-color: $color_primary_dark;
  border: none;
  @media screen and (min-width: $max_width_device + 1) {
    font-size: 1.5em;
    padding: 48px;
  }
  @media screen and (max-width: $max_width_device) {
    font-size: 1.2em;
    position: sticky;
    top: 0;
  }
}
.drawer-nav {
  .toggler {
    color: $color_text;
    display: inline-block;
    padding: 16px;
  }
  .drawer {
    height: 100%;
    width: 100%;
    top: 0;
    left: 0;
    background-color: $color_primary_dark;
    position: fixed !important;
    overflow: auto;
    &.drawer-show {
      margin-left: 0;
      transition: all 0.4s ease;
    }
    &.drawer-hide {
      margin-left: -100%;
      transition: all 0.4s ease;
    }
    .container {
      padding: 0 16px;
      i {
        margin-right: 16px;
      }
    }
    .drawer-item-active {
      color: $color_accent;
    }
  }
}
nav {
  &.top-nav {
    text-align: right;
    .nav-item {
      margin: 0 16px;
    }
  }
  .nav-item {
    border: none;
    display: inline-block;
    outline: 0;
    padding: 16px;
    vertical-align: middle;
    overflow: hidden;
    text-decoration: none;
    color: inherit;
    background-color: inherit;
    text-align: center;
    cursor: pointer;
    white-space: nowrap;
    -webkit-touch-callout: none;
    -webkit-user-select: none;
    -khtml-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
    &:hover {
      opacity: 0.8;
      color: $color_accent;
      transition: all 0.3s ease-out;
    }
  }
}
</style>
