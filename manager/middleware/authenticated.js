export default async function({
  $axios,
  isServer,
  route,
  req,
  redirect,
  store
}) {
  if (isServer && !req) {
    return
  }
  if (route.path === '/login') {
    return
  }

  let accessToken = window.localStorage.getItem('accessToken')

  if (accessToken && !store.state.accessToken) {
    await $axios
      .post(
        '/v1/auth/token',
        {},
        {
          headers: {
            Authorization: 'Bearer ' + accessToken
          }
        }
      )
      .then(res => {
        window.localStorage.setItem('accessToken', res.data.access_token)
        store.commit(
          'setAccessToken',
          window.localStorage.getItem('accessToken')
        )
      })
      .catch(e => {
        redirectToLogin(redirect, route.path)
      })
  }
  if (!accessToken && !store.state.accessToken) {
    redirectToLogin(redirect, route.path)
  }
}

const redirectToLogin = (redirect, path) => {
  if (path === '/logout') {
    redirect('/login')
  }
  redirect(`/login?redirect=${path}`)
}
