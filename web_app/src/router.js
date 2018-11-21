import Vue from 'vue'
import Router from 'vue-router'
import LoginComponent from "./views/login.vue"
import SecureComponent from "./views/secure.vue"

Vue.use(Router);

export default new Router({
    routes: [
        {
            path: '/',
            meta: {title: 'Home | Nehalem'},
            redirect: {
                name: "login"
            }
        },
        {
            path: "/login",
            meta: {title: 'Login | Nehalem'},
            name: "login",
            component: LoginComponent
        },
        {
            path: "/secure",
            meta: {title: 'Home | Secure'},
            name: "secure",
            component: SecureComponent
        }
    ]
})