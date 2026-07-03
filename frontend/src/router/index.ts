import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: () => import('../views/HomeView.vue') },
    { path: '/admin/login', name: 'login', component: () => import('../views/admin/LoginView.vue') },
    {
      path: '/admin',
      component: () => import('../views/admin/AdminLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', name: 'admin-dashboard', component: () => import('../views/admin/AdminDashboard.vue') },
        { path: 'sites', name: 'admin-sites', component: () => import('../views/admin/AdminSites.vue') },
        { path: 'categories', name: 'admin-categories', component: () => import('../views/admin/AdminCategories.vue') },
        { path: 'tags', name: 'admin-tags', component: () => import('../views/admin/AdminTags.vue') },
        { path: 'settings', name: 'admin-settings', component: () => import('../views/admin/AdminSettings.vue') },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.token) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.name === 'login' && auth.token) {
    return { name: 'admin-dashboard' }
  }
  return true
})

export default router
