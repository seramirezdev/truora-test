import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        domainData: null
    },
    mutations: {
        domainFound(state, domain) {
            state.domainData = domain
        }
    },
    actions: {

    }
})
