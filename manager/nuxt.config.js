const pkg = require('./package')
import colors from 'vuetify/es5/util/colors'

module.exports = {
  mode: 'spa',

  /*
  ** Headers of the page
  */
  head: {
    title: `LMM Manager`,
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content: pkg.description }
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
      { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css?family=Roboto:300,400,500,700|Material+Icons' }
    ]
  },

  /*
  ** Customize the progress-bar color
  */
  loading: { color: colors.default.blue.darken3 },

  /*
  ** Global CSS
  */
  css: [
    '~/assets/style/app.styl'
  ],

  /*
  ** Plugins to load before mounting the App
  */
  plugins: [
    '~/plugins/vuetify',
    '~/plugins/highlight'
  ],

  /*
  ** Router config
  */
  router: {
    middleware: 'authenticated'
  },

  /*
  ** Nuxt.js modules
  */
  modules: [
    // Doc: https://github.com/nuxt-community/axios-module#usage
    '@nuxtjs/axios'
  ],
  /*
  ** Axios module configuration
  */
  axios: {
    // See https://github.com/nuxt-community/axios-module#options
    baseURL: process.env.API_URL,
    browserBaseURL: process.env.API_URL_BROWSER
  },

  /*
  ** Build configuration
  */
  build: {
    /*
    ** You can extend webpack config here
    */
    extend(config, ctx) {
      // Run ESLint on save
      if (ctx.isDev && ctx.isClient) {
        config.module.rules.push({
          enforce: 'pre',
          test: /\.(js|vue)$/,
          loader: 'eslint-loader',
          exclude: /(node_modules)/,
          options: {
            fix: true
          }
        })
      }
    }
  },

  /*
  ** Server configuration
  */
  server: {
    host: process.env.SERVER_HOST,
    port: process.env.SERVER_PORT
  },

  /*
  ** Environment configuration
  */
  env: {
    ASSET_URL: process.env.ASSET_URL
  }
}
