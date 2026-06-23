import { createRouter, createWebHistory } from 'vue-router'
import HomeFeedView from './views/HomeFeedView.vue'
import PostDetailView from './views/PostDetailView.vue'
import ComposeView from './views/ComposeView.vue'
import ProfileView from './views/ProfileView.vue'
import LoginView from './views/LoginView.vue'
import RegisterView from './views/RegisterView.vue'
import ModerationView from './views/ModerationView.vue'
import FavoritesView from './views/FavoritesView.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: HomeFeedView },
    { path: '/posts/:id', name: 'post-detail', component: PostDetailView },
    { path: '/compose', name: 'compose', component: ComposeView },
    { path: '/login', name: 'login', component: LoginView },
    { path: '/register', name: 'register', component: RegisterView },
    { path: '/favorites', name: 'favorites', component: FavoritesView },
    { path: '/profile', name: 'profile', component: ProfileView },
    { path: '/moderation', name: 'moderation', component: ModerationView },
  ],
})
